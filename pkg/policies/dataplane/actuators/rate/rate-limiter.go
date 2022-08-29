package rate

import (
	"context"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/distcache"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/rate/ratetracker"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/selectors"
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

	etcdPath := path.Join(paths.RateLimiterConfigPath, paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type rateLimiterFactory struct {
	engineAPI                 iface.Engine
	distCache                 *distcache.DistCache
	lazySyncJobGroup          *jobs.JobGroup
	rateLimitDecisionsWatcher notifiers.Watcher
	agentGroupName            string
}

// main fx app.
func setupRateLimiterFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	distCache *distcache.DistCache,
	statusRegistry *status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) error {
	agentGroupName := ai.GetAgentGroup()
	rateLimitDecisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient,
		path.Join(paths.RateLimiterDecisionsPath, paths.AgentGroupPrefix(agentGroupName)))
	if err != nil {
		return err
	}
	lazySyncJobGroup, err := jobs.NewJobGroup(rateLimiterStatusRoot+".lazy_sync_jobs", statusRegistry, 0, jobs.RescheduleMode, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create lazy sync job group")
		return err
	}

	rateLimiterFactory := &rateLimiterFactory{
		engineAPI:                 e,
		distCache:                 distCache,
		lazySyncJobGroup:          lazySyncJobGroup,
		rateLimitDecisionsWatcher: rateLimitDecisionsWatcher,
		agentGroupName:            agentGroupName,
	}

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{rateLimiterFactory.newRateLimiterOptions},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusRegistry:     statusRegistry,
		PrometheusRegistry: prometheusRegistry,
		StatusPath:         rateLimiterStatusRoot,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := lazySyncJobGroup.Start()
			if err != nil {
				return err
			}
			err = rateLimitDecisionsWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			var err, merr error
			err = rateLimitDecisionsWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = lazySyncJobGroup.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
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
	registry *status.Registry,
) (fx.Option, error) {
	registryPath := path.Join(rateLimiterStatusRoot, key.String())

	wrapperMessage := &configv1.RateLimiterWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.RateLimiter == nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal rate limiter config")
		return fx.Options(), err
	}

	rateLimiterMessage := wrapperMessage.RateLimiter

	rateLimiter := &rateLimiter{
		Component:          wrapperMessage,
		rateLimiterProto:   rateLimiterMessage,
		registryPath:       registryPath,
		rateLimiterFactory: rateLimiterFactory,
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
	rateLimiterFactory *rateLimiterFactory
	rateTracker        ratetracker.RateTracker
	rateLimitChecker   *ratetracker.BasicRateLimitChecker
	registryPath       string
	rateLimiterProto   *policylangv1.RateLimiter
	name               string
}

// Make sure rateLimiter implements iface.Limiter.
var _ iface.RateLimiter = (*rateLimiter)(nil)

func (rateLimiter *rateLimiter) setup(lifecycle fx.Lifecycle, statusRegistry *status.Registry) error {
	// decision notifier
	unmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(paths.DataplaneComponentKey(rateLimiter.rateLimiterFactory.agentGroupName,
			rateLimiter.GetPolicyName(),
			rateLimiter.GetComponentIndex())),
		unmarshaller,
		rateLimiter.decisionUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			rateLimiter.rateLimitChecker = ratetracker.NewBasicRateLimitChecker()
			// loop through overrides
			for _, override := range rateLimiter.rateLimiterProto.GetOverrides() {
				label := rateLimiter.rateLimiterProto.GetLabelKey() + ":" + override.GetLabelValue()
				rateLimiter.rateLimitChecker.AddOverride(label, override.GetLimitScaleFactor())
			}
			rateLimiter.rateTracker, err = ratetracker.NewOlricRateTracker(rateLimiter.rateLimitChecker,
				rateLimiter.rateLimiterFactory.distCache,
				rateLimiter.name,
				rateLimiter.rateLimiterProto.GetLimitResetInterval().AsDuration())
			if err != nil {
				log.Error().Err(err).Msg("Failed to create limiter")
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
						log.Error().Err(err).Msg("Failed to create lazy limiter")
						return err
					}
				}
			}
			// add decisions notifier
			err = rateLimiter.rateLimiterFactory.rateLimitDecisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				log.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}

			// add to data engine
			err = rateLimiter.rateLimiterFactory.engineAPI.RegisterRateLimiter(rateLimiter)
			if err != nil {
				log.Error().Err(err).Msg("Failed to register rate limiter")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error
			// remove from data engine
			err = rateLimiter.rateLimiterFactory.engineAPI.UnregisterRateLimiter(rateLimiter)
			if err != nil {
				log.Error().Err(err).Msg("Failed to unregister rate limiter")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = rateLimiter.rateLimiterFactory.rateLimitDecisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				log.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			rateLimiter.rateTracker.Close()
			s := status.NewStatus(nil, merr)
			err = statusRegistry.Push(rateLimiter.registryPath, s)
			if err != nil {
				log.Error().Err(err).Msg("Failed to push status")
				merr = multierr.Append(merr, err)
			}

			return merr
		},
	})
	return nil
}

// GetSelector returns the selector for the rate limiter.
func (rateLimiter *rateLimiter) GetSelector() *selectorv1.Selector {
	return rateLimiter.rateLimiterProto.GetSelector()
}

// RunLimiter runs the limiter.
func (rateLimiter *rateLimiter) RunLimiter(labels selectors.Labels) *flowcontrolv1.LimiterDecision {
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
		Details: &flowcontrolv1.LimiterDecision_RateLimiter_{
			RateLimiter: &flowcontrolv1.LimiterDecision_RateLimiter{
				Label:     label,
				Remaining: int64(remaining),
				Current:   int64(current),
			},
		},
	}
}

// TakeN takes n tokens from the limiter.
func (rateLimiter *rateLimiter) TakeN(labels selectors.Labels, n int) (label string, ok bool, remaining int, current int) {
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
	if event.Type == notifiers.Remove {
		log.Debug().Msg("Decision removed")
		rateLimiter.rateLimitChecker.SetRateLimit(-1)
		return
	}

	var wrapperMessage configv1.RateLimiterDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.RateLimiterDecision == nil {
		return
	}
	if wrapperMessage.PolicyHash != rateLimiter.GetPolicyHash() {
		return
	}
	limitDecision := wrapperMessage.RateLimiterDecision
	rateLimiter.rateLimitChecker.SetRateLimit(int(limitDecision.GetLimit()))
}

// GetLimiterID returns the limiter ID.
func (rateLimiter *rateLimiter) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:     rateLimiter.GetPolicyName(),
		ComponentIndex: rateLimiter.GetComponentIndex(),
		PolicyHash:     rateLimiter.GetPolicyHash(),
	}
}
