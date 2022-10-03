package rate

import (
	"context"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/distcache"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/rate/ratetracker"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/status"
)

const rateLimiterStatusRoot = "rate_limiters"

var fxNameTag = config.NameTag(rateLimiterStatusRoot)

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

	etcdPath := path.Join(common.RateLimiterConfigPath, common.AgentGroupPrefix(agentGroupName))
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
	etcdPath := path.Join(common.RateLimiterDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	dynamicConfigWatcher, err := etcdwatcher.NewWatcher(etcdClient, common.RateLimiterDynamicConfigPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child(rateLimiterStatusRoot)
	logger := reg.GetLogger()

	lazySyncJobGroup, err := jobs.NewJobGroup(reg.Child("lazy_sync_jobs"), 0, jobs.RescheduleMode, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create lazy sync job group")
		return err
	}

	rateLimiterFactory := &rateLimiterFactory{
		engineAPI:            e,
		distCache:            distCache,
		lazySyncJobGroup:     lazySyncJobGroup,
		decisionsWatcher:     decisionsWatcher,
		dynamicConfigWatcher: dynamicConfigWatcher,
		agentGroupName:       agentGroupName,
		registry:             reg,
	}

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{
			rateLimiterFactory.newRateLimiterOptions,
		},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusRegistry:     reg,
		PrometheusRegistry: prometheusRegistry,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := lazySyncJobGroup.Start()
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
			reg.Detach()
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// per policy component.
func (rateLimiterFactory *rateLimiterFactory) newRateLimiterOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := rateLimiterFactory.registry.GetLogger()
	wrapperMessage := &wrappersv1.RateLimiterWrapper{}
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
		rateLimiterFactory: rateLimiterFactory,
		registry:           reg,
	}
	rateLimiter.name = iface.ComponentID(rateLimiter)

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
	etcdKey := common.DataplaneComponentKey(rateLimiter.rateLimiterFactory.agentGroupName, rateLimiter.GetPolicyName(), rateLimiter.GetComponentIndex())
	// decision notifier
	decisionUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		rateLimiter.decisionUpdateCallback,
	)
	// dynamic config notifier
	dynamicConfigUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	dynamicConfigNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		dynamicConfigUnmarshaller,
		rateLimiter.dynamicConfigUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			rateLimiter.rateLimitChecker = ratetracker.NewBasicRateLimitChecker()
			rateLimiter.updateDynamicConfig(rateLimiter.rateLimiterProto.GetInitConfig())
			rateLimiter.rateTracker, err = ratetracker.NewDistCacheRateTracker(
				rateLimiter.rateLimitChecker,
				rateLimiter.rateLimiterFactory.distCache,
				rateLimiter.name,
				rateLimiter.rateLimiterProto.GetLimitResetInterval().AsDuration())
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			// check whether lazy limiter is enabled
			if lazySyncConfig := rateLimiter.rateLimiterProto.GetLazySync(); lazySyncConfig != nil {
				if lazySyncConfig.GetEnabled() {
					lazySyncInterval := time.Duration(int64(rateLimiter.rateLimiterProto.GetLimitResetInterval().AsDuration()) / int64(lazySyncConfig.GetNumSync()))
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
		label := rateLimiter.rateLimiterProto.GetLabelKey() + ":" + override.GetLabelValue()
		overrides[label] = override.GetLimitScaleFactor()
	}

	logger.Debug().Interface("overrides", overrides).Str("name", rateLimiter.name).Msgf("Updating dynamic config for rate limiter")

	rateLimiter.rateLimitChecker.SetOverrides(overrides)
}

// GetSelector returns the selector for the rate limiter.
func (rateLimiter *rateLimiter) GetSelector() *selectorv1.Selector {
	return rateLimiter.rateLimiterProto.GetSelector()
}

// RunLimiter runs the limiter.
func (rateLimiter *rateLimiter) RunLimiter(labels map[string]string) *flowcontrolv1.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED

	label, ok, remaining, current := rateLimiter.TakeN(labels, 1)

	if label == "" {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:     rateLimiter.GetPolicyName(),
		PolicyHash:     rateLimiter.GetPolicyHash(),
		ComponentIndex: rateLimiter.GetComponentIndex(),
		Dropped:        !ok,
		Reason:         reason,
		Details: &flowcontrolv1.LimiterDecision_RateLimiterInfo_{
			RateLimiterInfo: &flowcontrolv1.LimiterDecision_RateLimiterInfo{
				Label:     label,
				Remaining: int64(remaining),
				Current:   int64(current),
			},
		},
	}
}

// TakeN takes n tokens from the limiter.
func (rateLimiter *rateLimiter) TakeN(labels map[string]string, n int) (label string, ok bool, remaining int, current int) {
	labelKey := rateLimiter.rateLimiterProto.GetLabelKey()
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

	var wrapperMessage wrappersv1.RateLimiterDecisionWrapper
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

	var wrapperMessage wrappersv1.RateLimiterDynamicConfigWrapper
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
		PolicyName:     rateLimiter.GetPolicyName(),
		PolicyHash:     rateLimiter.GetPolicyHash(),
		ComponentIndex: rateLimiter.GetComponentIndex(),
	}
}

// GetObserver is there to satisfy Limiter interface.
func (rateLimiter *rateLimiter) GetObserver(map[string]string) prometheus.Observer {
	// Not implemented for rate limiter
	return nil
}
