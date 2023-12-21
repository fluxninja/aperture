package ratelimiter

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"time"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/rate-limiter"
	globaltokenbucket "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/rate-limiter/global-token-bucket"
	lazysync "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/rate-limiter/lazy-sync"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/status"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/types/known/durationpb"
)

const rateLimiterStatusRoot = "rate_limiters"

var (
	fxTag           = config.NameTag(rateLimiterStatusRoot)
	metricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIDLabel, metrics.DecisionTypeLabel, metrics.LimiterDroppedLabel}
)

func rateLimiterModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideRateLimiterWatchers,
				fx.ResultTags(fxTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupRateLimiterFactory,
				fx.ParamTags(fxTag),
			),
		),
	)
}

func provideRateLimiterWatchers(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroupName := ai.GetAgentGroup()

	etcdPath := path.Join(paths.RateLimiterConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type rateLimiterFactory struct {
	engineAPI        iface.Engine
	registry         status.Registry
	distCache        *distcache.DistCache
	lazySyncJobGroup *jobs.JobGroup
	decisionsWatcher notifiers.Watcher
	counterVector    *prometheus.CounterVec
	agentGroupName   string
}

// main fx app.
func setupRateLimiterFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	distCache *distcache.DistCache,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) error {
	agentGroupName := ai.GetAgentGroup()
	etcdPath := path.Join(paths.RateLimiterDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", rateLimiterStatusRoot)

	logger := reg.GetLogger()

	lazySyncJobGroup, err := jobs.NewJobGroup(reg.Child("sync", "lazy_sync_jobs"), jobs.JobGroupConfig{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create lazy sync job group")
		return err
	}

	counterVector := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.RateLimiterCounterTotalMetricName,
		Help: "A counter measuring the number of times Rate Limiter was triggered",
	}, metricLabelKeys)

	rateLimiterFactory := &rateLimiterFactory{
		engineAPI:        e,
		distCache:        distCache,
		lazySyncJobGroup: lazySyncJobGroup,
		decisionsWatcher: decisionsWatcher,
		agentGroupName:   agentGroupName,
		registry:         reg,
		counterVector:    counterVector,
	}

	fxDriver, err := notifiers.NewFxDriver(
		reg,
		prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{rateLimiterFactory.newRateLimiterOptions},
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := prometheusRegistry.Register(rateLimiterFactory.counterVector)
			if err != nil {
				return err
			}
			err = lazySyncJobGroup.Start()
			if err != nil {
				return err
			}
			err = decisionsWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			var err, merr error
			err = decisionsWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = lazySyncJobGroup.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(rateLimiterFactory.counterVector) {
				err2 := fmt.Errorf("failed to unregister metric")
				merr = multierr.Append(merr, err2)
			}
			reg.Detach()
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// per component fx app.
func (rlFactory *rateLimiterFactory) newRateLimiterOptions(key notifiers.Key, unmarshaller config.Unmarshaller, reg status.Registry) (fx.Option, error) {
	logger := rlFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.RateLimiterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.RateLimiter == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal rate limiter config")
		return fx.Options(), err
	}

	rlProto := wrapperMessage.RateLimiter
	rl := &rateLimiter{
		Component: wrapperMessage.GetCommonAttributes(),
		rlProto:   rlProto,
		rlFactory: rlFactory,
		registry:  reg,
	}
	rl.name = iface.ComponentKey(rl)

	return fx.Options(
		fx.Invoke(
			rl.setup,
		),
	), nil
}

// rateLimiter implements rate limiter on the data plane side.
type rateLimiter struct {
	iface.Component
	registry  status.Registry
	rlFactory *rateLimiterFactory
	limiter   ratelimiter.RateLimiter
	inner     *globaltokenbucket.GlobalTokenBucket
	rlProto   *policylangv1.RateLimiter
	name      string
}

// Make sure rateLimiter implements iface.Limiter.
var _ iface.Limiter = (*rateLimiter)(nil)

func (rl *rateLimiter) setup(lifecycle fx.Lifecycle) error {
	logger := rl.registry.GetLogger()
	etcdKey := paths.AgentComponentKey(
		rl.rlFactory.agentGroupName,
		rl.GetPolicyName(),
		rl.GetComponentId(),
	)
	// decision notifier
	decisionUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		rl.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = rl.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = rl.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = rl.GetComponentId()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			rl.inner, err = globaltokenbucket.NewGlobalTokenBucket(
				rl.rlFactory.distCache,
				rl.name,
				rl.rlProto.Parameters.GetInterval().AsDuration(),
				rl.rlProto.Parameters.GetMaxIdleTime().AsDuration(),
				rl.rlProto.Parameters.GetContinuousFill(),
				rl.rlProto.Parameters.GetDelayInitialFill(),
			)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			rl.limiter = rl.inner
			// check whether lazy limiter is enabled
			if lazySyncConfig := rl.rlProto.Parameters.GetLazySync(); lazySyncConfig != nil {
				if lazySyncConfig.GetEnabled() {
					rl.limiter, err = lazysync.NewLazySyncRateLimiter(rl.limiter,
						rl.rlProto.Parameters.GetInterval().AsDuration(),
						lazySyncConfig.GetNumSync(),
						rl.rlFactory.lazySyncJobGroup)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to create lazy limiter")
						return err
					}
				}
			}

			// add decisions notifier
			err = rl.rlFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}

			// add to data engine
			err = rl.rlFactory.engineAPI.RegisterRateLimiter(rl)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register rate limiter")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error
			deleted := rl.rlFactory.counterVector.DeletePartialMatch(metricLabels)
			if deleted == 0 {
				logger.Warn().Msg("Could not delete rate limiter counter from its metric vector. No traffic to generate metrics?")
			}
			// remove from data engine
			err = rl.rlFactory.engineAPI.UnregisterRateLimiter(rl)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister rate limiter")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = rl.rlFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			rl.limiter.Close()
			rl.registry.SetStatus(status.NewStatus(nil, merr))

			return merr
		},
	})
	return nil
}

// GetSelectors returns the selectors for the rate limiter.
func (rl *rateLimiter) GetSelectors() []*policylangv1.Selector {
	return rl.rlProto.GetSelectors()
}

// Decide runs the limiter.
func (rl *rateLimiter) Decide(ctx context.Context, labels labels.Labels) *flowcontrolv1.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED

	tokens := float64(1)
	// get tokens from labels
	rParams := rl.rlProto.GetRequestParameters()
	var deniedResponseStatusCode flowcontrolv1.StatusCode
	if rParams != nil {
		deniedResponseStatusCode = rParams.GetDeniedResponseStatusCode()
		tokensLabelKey := rParams.GetTokensLabelKey()
		if tokensLabelKey != "" {
			if val, ok := labels.Get(tokensLabelKey); ok {
				if parsedTokens, err := strconv.ParseFloat(val, 64); err == nil {
					tokens = parsedTokens
				}
			}
		}
	}

	label, ok, waitTime, remaining, current := rl.takeIfAvailable(ctx, labels, tokens)

	tokensConsumed := float64(0)
	if ok {
		tokensConsumed = tokens
	}

	if label == "" {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:               rl.GetPolicyName(),
		PolicyHash:               rl.GetPolicyHash(),
		ComponentId:              rl.GetComponentId(),
		Dropped:                  !ok,
		DeniedResponseStatusCode: deniedResponseStatusCode,
		Reason:                   reason,
		WaitTime:                 durationpb.New(waitTime),
		Details: &flowcontrolv1.LimiterDecision_RateLimiterInfo_{
			RateLimiterInfo: &flowcontrolv1.LimiterDecision_RateLimiterInfo{
				Label: label,
				TokensInfo: &flowcontrolv1.LimiterDecision_TokensInfo{
					Remaining: remaining,
					Current:   current,
					Consumed:  tokensConsumed,
				},
			},
		},
	}
}

// Revert returns the tokens to the limiter.
func (rl *rateLimiter) Revert(ctx context.Context, labels labels.Labels, decision *flowcontrolv1.LimiterDecision) {
	if rl.limiter.GetPassThrough() {
		return
	}

	if rateLimiterDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_RateLimiterInfo_); ok {
		tokens := rateLimiterDecision.RateLimiterInfo.TokensInfo.Consumed
		if tokens > 0 {
			rl.takeIfAvailable(ctx, labels, -tokens)
		}
	}
}

// takeIfAvailable takes n tokens from the limiter.
func (rl *rateLimiter) takeIfAvailable(
	ctx context.Context,
	labels labels.Labels,
	n float64,
) (label string, ok bool, waitTime time.Duration, remaining float64, current float64) {
	if rl.limiter.GetPassThrough() {
		return label, true, 0, 0, 0
	}

	labelKey := rl.rlProto.Parameters.GetLimitByLabelKey()
	if labelKey == "" {
		// Deprecated: Remove in v3.0.0
		labelKey = rl.rlProto.Parameters.GetLabelKey()
	}
	if labelKey == "" {
		label = "default"
	} else {
		labelValue, found := labels.Get(labelKey)
		if !found {
			return "", true, 0, 0, 0
		}
		label = labelKey + ":" + labelValue
	}

	ok, waitTime, remaining, current = rl.limiter.TakeIfAvailable(ctx, label, n)
	return
}

func (rl *rateLimiter) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := rl.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		rl.limiter.SetPassThrough(true)
		return
	}

	var wrapperMessage policysyncv1.RateLimiterDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.RateLimiterDecision == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != rl.GetPolicyHash() {
		return
	}
	limitDecision := wrapperMessage.RateLimiterDecision
	rl.inner.SetBucketCapacity(limitDecision.BucketCapacity)
	rl.inner.SetFillAmount(limitDecision.FillAmount)
	rl.inner.SetPassThrough(limitDecision.PassThrough)
}

// GetLimiterID returns the limiter ID.
func (rl *rateLimiter) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  rl.GetPolicyName(),
		PolicyHash:  rl.GetPolicyHash(),
		ComponentID: rl.GetComponentId(),
	}
}

// GetRequestCounter returns counter for tracking number of times rateLimiter was triggered.
func (rl *rateLimiter) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := rl.rlFactory.counterVector.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get counter")
		return nil
	}
	return counter
}

// GetRampMode is always false for rateLimiters.
func (rl *rateLimiter) GetRampMode() bool {
	return false
}
