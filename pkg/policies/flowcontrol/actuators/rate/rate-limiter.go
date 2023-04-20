package rate

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/distcache"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/rate/ratetracker"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
)

const rateLimiterStatusRoot = "rate_limiters"

var (
	fxNameTag       = config.NameTag(rateLimiterStatusRoot)
	metricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIDLabel, metrics.DecisionTypeLabel, metrics.LimiterDroppedLabel}
)

func rateLimiterModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideWatchers,
				fx.ResultTags(fxNameTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupRateLimiterFactory,
				fx.ParamTags(fxNameTag),
			),
		),
	)
}

func provideWatchers(
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
	engineAPI            iface.Engine
	registry             status.Registry
	distCache            *distcache.DistCache
	lazySyncJobGroup     *jobs.JobGroup
	decisionsWatcher     notifiers.Watcher
	dynamicConfigWatcher notifiers.Watcher
	counterVector        *prometheus.CounterVec
	agentGroupName       string
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

	dynamicConfigWatcher, err := etcdwatcher.NewWatcher(etcdClient,
		paths.RateLimiterDynamicConfigPath)
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
		Name: metrics.RateLimiterCounterMetricName,
		Help: "A counter measuring the number of times Rate Limiter was triggered",
	}, metricLabelKeys)

	rateLimiterFactory := &rateLimiterFactory{
		engineAPI:            e,
		distCache:            distCache,
		lazySyncJobGroup:     lazySyncJobGroup,
		decisionsWatcher:     decisionsWatcher,
		dynamicConfigWatcher: dynamicConfigWatcher,
		agentGroupName:       agentGroupName,
		registry:             reg,
		counterVector:        counterVector,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{rateLimiterFactory.newRateLimiterOptions})
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
			err = dynamicConfigWatcher.Start()
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
			err = dynamicConfigWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = lazySyncJobGroup.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(rateLimiterFactory.counterVector) {
				err2 := fmt.Errorf("failed to unregister rate_limiter_counter metric")
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
func (rlFactory *rateLimiterFactory) newRateLimiterOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := rlFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.RateLimiterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.RateLimiter == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal rate limiter config")
		return fx.Options(), err
	}

	rateLimiterProto := wrapperMessage.RateLimiter
	rateLimiter := &rateLimiter{
		Component:          wrapperMessage.GetCommonAttributes(),
		rateLimiterProto:   rateLimiterProto,
		rateLimiterFactory: rlFactory,
		registry:           reg,
	}
	rateLimiter.name = iface.ComponentKey(rateLimiter)

	return fx.Options(
		fx.Invoke(
			rateLimiter.setup,
		),
	), nil
}

// rateLimiter implements rate limiter on the data plane side.
type rateLimiter struct {
	iface.Component
	registry           status.Registry
	rateLimiterFactory *rateLimiterFactory
	rateTracker        ratetracker.RateTracker
	rateLimitChecker   *ratetracker.BasicRateLimitChecker
	rateLimiterProto   *policylangv1.RateLimiter
	name               string
}

// Make sure rateLimiter implements iface.Limiter.
var _ iface.RateLimiter = (*rateLimiter)(nil)

func (rateLimiter *rateLimiter) setup(lifecycle fx.Lifecycle) error {
	logger := rateLimiter.registry.GetLogger()
	etcdKey := paths.AgentComponentKey(rateLimiter.rateLimiterFactory.agentGroupName,
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
	// dynamic config notifier
	dynamicConfigUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	dynamicConfigNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		dynamicConfigUnmarshaller,
		rateLimiter.dynamicConfigUpdateCallback,
	)
	if err != nil {
		return err
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = rateLimiter.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = rateLimiter.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = rateLimiter.GetComponentId()
	rateCounterVec := rateLimiter.rateLimiterFactory.counterVector

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			rateLimiter.rateLimitChecker = ratetracker.NewBasicRateLimitChecker()
			rateLimiter.updateDynamicConfig(rateLimiter.rateLimiterProto.GetDefaultConfig())
			rateLimiter.rateTracker, err = ratetracker.NewDistCacheRateTracker(
				rateLimiter.rateLimitChecker,
				rateLimiter.rateLimiterFactory.distCache,
				rateLimiter.name,
				rateLimiter.rateLimiterProto.Parameters.GetLimitResetInterval().AsDuration())
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			// check whether lazy limiter is enabled
			if lazySyncConfig := rateLimiter.rateLimiterProto.Parameters.GetLazySync(); lazySyncConfig != nil {
				if lazySyncConfig.GetEnabled() {
					lazySyncInterval := time.Duration(int64(rateLimiter.rateLimiterProto.Parameters.GetLimitResetInterval().AsDuration()) / int64(lazySyncConfig.GetNumSync()))
					rateLimiter.rateTracker, err = ratetracker.NewLazySyncRateTracker(rateLimiter.rateTracker,
						lazySyncInterval,
						rateLimiter.rateLimiterFactory.lazySyncJobGroup)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to create lazy limiter")
						return err
					}
				}
			}
			// add decisions notifier
			err = rateLimiter.rateLimiterFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}
			// add dynamic config notifier
			err = rateLimiter.rateLimiterFactory.dynamicConfigWatcher.AddKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add dynamic config notifier")
				return err
			}

			// add to data engine
			err = rateLimiter.rateLimiterFactory.engineAPI.RegisterRateLimiter(rateLimiter)
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
			err = rateLimiter.rateLimiterFactory.engineAPI.UnregisterRateLimiter(rateLimiter)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister rate limiter")
				merr = multierr.Append(merr, err)
			}
			// remove dynamic config notifier
			err = rateLimiter.rateLimiterFactory.dynamicConfigWatcher.RemoveKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove dynamic config notifier")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = rateLimiter.rateLimiterFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			rateLimiter.rateTracker.Close()
			rateLimiter.registry.SetStatus(status.NewStatus(nil, merr))

			return merr
		},
	})
	return nil
}

func (rateLimiter *rateLimiter) updateDynamicConfig(dynamicConfig *policylangv1.RateLimiter_DynamicConfig) {
	logger := rateLimiter.registry.GetLogger()

	if dynamicConfig == nil {
		return
	}
	overrides := ratetracker.Overrides{}
	// loop through overrides
	for _, override := range dynamicConfig.GetOverrides() {
		label := rateLimiter.rateLimiterProto.Parameters.GetLabelKey() + ":" + override.GetLabelValue()
		overrides[label] = override.GetLimitScaleFactor()
	}

	logger.Debug().Interface("overrides", overrides).Str("name", rateLimiter.name).Msgf("Updating dynamic config for rate limiter")

	rateLimiter.rateLimitChecker.SetOverrides(overrides)
}

// GetFlowSelector returns the selector for the rate limiter.
func (rateLimiter *rateLimiter) GetFlowSelector() *policylangv1.FlowSelector {
	return rateLimiter.rateLimiterProto.GetFlowSelector()
}

// Decide runs the limiter.
func (rateLimiter *rateLimiter) Decide(ctx context.Context,
	labels map[string]string,
) *flowcontrolv1.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED

	tokens := uint64(1)
	// get tokens from labels
	if rateLimiter.rateLimiterProto.Parameters.TokensLabelKey != "" {
		if val, ok := labels[rateLimiter.rateLimiterProto.Parameters.TokensLabelKey]; ok {
			if parsedTokens, err := strconv.ParseUint(val, 10, 64); err == nil {
				tokens = parsedTokens
			}
		}
	}

	label, ok, remaining, current := rateLimiter.TakeN(labels, int(tokens))

	tokensConsumed := uint64(0)
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
				Remaining:      int64(remaining),
				Current:        int64(current),
				TokensConsumed: tokensConsumed,
			},
		},
	}
}

// Revert returns the tokens to the limiter.
func (rateLimiter *rateLimiter) Revert(labels map[string]string, decision *flowcontrolv1.LimiterDecision) {
	if rateLimiterDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_RateLimiterInfo_); ok {
		tokens := rateLimiterDecision.RateLimiterInfo.TokensConsumed
		if tokens > 0 {
			rateLimiter.TakeN(labels, -int(tokens))
		}
	}
}

// TakeN takes n tokens from the limiter.
func (rateLimiter *rateLimiter) TakeN(labels map[string]string, n int) (label string, ok bool, remaining int, current int) {
	labelKey := rateLimiter.rateLimiterProto.Parameters.GetLabelKey()
	var labelValue string
	if val, found := labels[labelKey]; found {
		labelValue = val
	} else {
		return "", true, -1, -1
	}

	label = labelKey + ":" + labelValue

	ok, remaining, current = rateLimiter.rateTracker.TakeN(label, n)
	return
}

func (rateLimiter *rateLimiter) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := rateLimiter.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		rateLimiter.rateLimitChecker.SetRateLimit(-1)
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
	if commonAttributes.PolicyHash != rateLimiter.GetPolicyHash() {
		return
	}
	limitDecision := wrapperMessage.RateLimiterDecision
	rateLimiter.rateLimitChecker.SetRateLimit(int(limitDecision.GetLimit()))
}

func (rateLimiter *rateLimiter) dynamicConfigUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := rateLimiter.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Dynamic config removed")
		rateLimiter.rateLimitChecker.SetOverrides(ratetracker.Overrides{})
		return
	}

	var wrapperMessage policysyncv1.RateLimiterDynamicConfigWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.RateLimiterDynamicConfig == nil {
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
	dynamicConfig := wrapperMessage.RateLimiterDynamicConfig
	rateLimiter.updateDynamicConfig(dynamicConfig)
}

// GetLimiterID returns the limiter ID.
func (rateLimiter *rateLimiter) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  rateLimiter.GetPolicyName(),
		PolicyHash:  rateLimiter.GetPolicyHash(),
		ComponentID: rateLimiter.GetComponentId(),
	}
}

// GetRequestCounter returns counter for tracking number of times rateLimiter was triggered.
func (rateLimiter *rateLimiter) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := rateLimiter.rateLimiterFactory.counterVector.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get counter")
		return nil
	}

	return counter
}
