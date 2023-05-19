package ratelimiter

import (
	"context"
	"fmt"
	"path"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/agentinfo"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/distcache"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	leakybucket "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/rate-limiter/leaky-bucket"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const leakyBucketRateLimiterStatusRoot = "leaky_bucket_rate_limiters"

var (
	lbFxTag         = config.NameTag(leakyBucketRateLimiterStatusRoot)
	metricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIDLabel, metrics.DecisionTypeLabel, metrics.LimiterDroppedLabel}
)

func leakyBucketRateLimiterModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideLeakyBucketRateLimiterWatchers,
				fx.ResultTags(lbFxTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupLeakyBucketFactory,
				fx.ParamTags(lbFxTag),
			),
		),
	)
}

func provideLeakyBucketRateLimiterWatchers(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroupName := ai.GetAgentGroup()

	etcdPath := path.Join(paths.LeakyBucketRateLimiterConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type leakyBucketFactory struct {
	engineAPI        iface.Engine
	registry         status.Registry
	distCache        *distcache.DistCache
	decisionsWatcher notifiers.Watcher
	counterVector    *prometheus.CounterVec
	agentGroupName   string
}

// main fx app.
func setupLeakyBucketFactory(
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
	etcdPath := path.Join(paths.LeakyBucketRateLimiterDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", leakyBucketRateLimiterStatusRoot)

	counterVector := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.RateLimiterCounterTotalMetricName,
		Help: "A counter measuring the number of times Rate Limiter was triggered",
	}, metricLabelKeys)

	rateLimiterFactory := &leakyBucketFactory{
		engineAPI:        e,
		distCache:        distCache,
		decisionsWatcher: decisionsWatcher,
		agentGroupName:   agentGroupName,
		registry:         reg,
		counterVector:    counterVector,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{rateLimiterFactory.newLeakyBucketOptions})
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := prometheusRegistry.Register(rateLimiterFactory.counterVector)
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
			if !prometheusRegistry.Unregister(rateLimiterFactory.counterVector) {
				err2 := fmt.Errorf("failed to unregister rate_limiter_counter_total metric")
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
func (rlFactory *leakyBucketFactory) newLeakyBucketOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := rlFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.LeakyBucketRateLimiterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.RateLimiter == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal rate limiter config")
		return fx.Options(), err
	}

	lbProto := wrapperMessage.RateLimiter
	lb := &leakyBucket{
		Component: wrapperMessage.GetCommonAttributes(),
		lbProto:   lbProto,
		lbFactory: rlFactory,
		registry:  reg,
	}
	lb.name = iface.ComponentKey(lb)

	return fx.Options(
		fx.Invoke(
			lb.setup,
		),
	), nil
}

// leakyBucket implements rate limiter on the data plane side.
type leakyBucket struct {
	iface.Component
	registry  status.Registry
	lbFactory *leakyBucketFactory
	lb        *leakybucket.LeakyBucketRateLimiter
	lbProto   *policylangv1.LeakyBucketRateLimiter
	name      string
}

// Make sure rateLimiter implements iface.Limiter.
var _ iface.RateLimiter = (*leakyBucket)(nil)

func (rateLimiter *leakyBucket) setup(lifecycle fx.Lifecycle) error {
	logger := rateLimiter.registry.GetLogger()
	etcdKey := paths.AgentComponentKey(rateLimiter.lbFactory.agentGroupName,
		rateLimiter.GetPolicyName(),
		rateLimiter.GetComponentId())
	// decision notifier
	decisionUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		rateLimiter.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = rateLimiter.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = rateLimiter.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = rateLimiter.GetComponentId()
	rateCounterVec := rateLimiter.lbFactory.counterVector

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			rateLimiter.lb, err = leakybucket.NewLeakyBucket(
				rateLimiter.lbFactory.distCache,
				rateLimiter.name,
				rateLimiter.lbProto.Parameters.GetMaxIdleTime().AsDuration(),
			)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			// add decisions notifier
			err = rateLimiter.lbFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}

			// add to data engine
			err = rateLimiter.lbFactory.engineAPI.RegisterRateLimiter(rateLimiter)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register rate limiter")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error
			deleted := rateCounterVec.DeletePartialMatch(metricLabels)
			if deleted == 0 {
				logger.Warn().Msg("Could not delete rate limiter counter from its metric vector. No traffic to generate metrics?")
			}
			// remove from data engine
			err = rateLimiter.lbFactory.engineAPI.UnregisterRateLimiter(rateLimiter)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister rate limiter")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = rateLimiter.lbFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			rateLimiter.lb.Close()
			rateLimiter.registry.SetStatus(status.NewStatus(nil, merr))

			return merr
		},
	})
	return nil
}

// GetSelectors returns the selectors for the rate limiter.
func (rateLimiter *leakyBucket) GetSelectors() []*policylangv1.Selector {
	return rateLimiter.lbProto.GetSelectors()
}

// Decide runs the limiter.
func (rateLimiter *leakyBucket) Decide(ctx context.Context, labels map[string]string) *flowcontrolv1.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED

	tokens := float64(1)
	// get tokens from labels
	if rateLimiter.lbProto.Parameters.TokensLabelKey != "" {
		if val, ok := labels[rateLimiter.lbProto.Parameters.TokensLabelKey]; ok {
			if parsedTokens, err := strconv.ParseFloat(val, 64); err == nil {
				tokens = parsedTokens
			}
		}
	}

	label, ok, remaining, current := rateLimiter.TakeIfAvailable(labels, tokens)

	tokensConsumed := float64(0)
	if ok {
		tokensConsumed = tokens
	}

	if label == "" {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:  rateLimiter.GetPolicyName(),
		PolicyHash:  rateLimiter.GetPolicyHash(),
		ComponentId: rateLimiter.GetComponentId(),
		Dropped:     !ok,
		Reason:      reason,
		Details: &flowcontrolv1.LimiterDecision_RateLimiterInfo_{
			RateLimiterInfo: &flowcontrolv1.LimiterDecision_RateLimiterInfo{
				Label:          label,
				Remaining:      remaining,
				Current:        current,
				TokensConsumed: tokensConsumed,
			},
		},
	}
}

// Revert returns the tokens to the limiter.
func (rateLimiter *leakyBucket) Revert(labels map[string]string, decision *flowcontrolv1.LimiterDecision) {
	if rateLimiterDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_RateLimiterInfo_); ok {
		tokens := rateLimiterDecision.RateLimiterInfo.TokensConsumed
		if tokens > 0 {
			rateLimiter.TakeIfAvailable(labels, -tokens)
		}
	}
}

// TakeIfAvailable takes n tokens from the limiter.
func (rateLimiter *leakyBucket) TakeIfAvailable(labels map[string]string, n float64) (label string, ok bool, remaining float64, current float64) {
	labelKey := rateLimiter.lbProto.Parameters.GetLabelKey()
	var labelValue string
	if val, found := labels[labelKey]; found {
		labelValue = val
	} else {
		return "", true, -1, -1
	}

	label = labelKey + ":" + labelValue

	if rateLimiter.lb.GetParameters().BucketCapacity < 0 {
		return label, true, -1, -1
	}

	ok, remaining, current = rateLimiter.lb.TakeIfAvailable(label, n)
	return
}

func (rateLimiter *leakyBucket) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := rateLimiter.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		rateLimiter.lb.SetParameters(
			leakybucket.Parameters{
				BucketCapacity: -1,
			})
		return
	}

	var wrapperMessage policysyncv1.LeakyBucketRateLimiterDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.RateLimiterDecision == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != rateLimiter.GetPolicyHash() {
		return
	}
	limitDecision := wrapperMessage.RateLimiterDecision
	rateLimiter.lb.SetParameters(
		leakybucket.Parameters{
			BucketCapacity: limitDecision.BucketCapacity,
			LeakAmount:     limitDecision.LeakAmount,
			LeakInterval:   limitDecision.LeakInterval.AsDuration(),
		})
}

// GetLimiterID returns the limiter ID.
func (rateLimiter *leakyBucket) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  rateLimiter.GetPolicyName(),
		PolicyHash:  rateLimiter.GetPolicyHash(),
		ComponentID: rateLimiter.GetComponentId(),
	}
}

// GetRequestCounter returns counter for tracking number of times rateLimiter was triggered.
func (rateLimiter *leakyBucket) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := rateLimiter.lbFactory.counterVector.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get counter")
		return nil
	}

	return counter
}
