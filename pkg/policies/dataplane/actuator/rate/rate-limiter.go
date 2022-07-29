package rate

import (
	"context"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policydecisionsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/distcache"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/flowcontrol/ratelimiter"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/component"
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
		jobs.JobGroupConstructor{Group: rateLimiterStatusRoot}.Annotate(),
		fx.Invoke(
			fx.Annotate(
				setupRateLimiterFactory,
				fx.ParamTags(fxNameTag, fxNameTag),
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
	engineAPI                 iface.EngineAPI
	distCache                 *distcache.DistCache
	lazySyncJobGroup          *jobs.JobGroup
	rateLimitDecisionsWatcher notifiers.Watcher
	agentGroupName            string
}

// main fx app.
func setupRateLimiterFactory(
	watcher notifiers.Watcher,
	lazySyncJobGroup *jobs.JobGroup,
	lifecycle fx.Lifecycle,
	e iface.EngineAPI,
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
			err := rateLimitDecisionsWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			err := rateLimitDecisionsWatcher.Stop()
			if err != nil {
				return err
			}
			return nil
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

	wrapperMessage := &configv1.ConfigPropertiesWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.Config == nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal rate limiter config")
		return fx.Options(), err
	}

	rateLimiterProto := &policylangv1.RateLimiter{}
	err = wrapperMessage.Config.UnmarshalTo(rateLimiterProto)
	if err != nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal rate limiter config")
		return fx.Options(), err
	}

	rateLimiter := &rateLimiter{
		ComponentAPI:       component.NewComponent(wrapperMessage),
		rateLimiterProto:   rateLimiterProto,
		registryPath:       registryPath,
		rateLimiterFactory: rateLimiterFactory,
	}
	rateLimiter.metricID = paths.MetricIDForComponent(rateLimiter)

	return fx.Options(
		fx.Invoke(
			rateLimiter.setup,
		),
	), nil
}

// rateLimiter implements rate limiter on the data plane side.
type rateLimiter struct {
	component.ComponentAPI
	rateLimiterFactory *rateLimiterFactory
	limiter            ratelimiter.RateLimiter
	limitCheck         *ratelimiter.BasicRateLimitCheck
	registryPath       string
	rateLimiterProto   *policylangv1.RateLimiter
	metricID           string
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
		notifiers.Key(paths.IdentifierForComponent(rateLimiter.rateLimiterFactory.agentGroupName,
			rateLimiter.ComponentAPI.GetPolicyName(),
			rateLimiter.ComponentAPI.GetComponentIndex())),
		unmarshaller,
		rateLimiter.decisionUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			rateLimiter.limitCheck = ratelimiter.NewBasicRateLimitCheck()
			// loop through overrides
			for _, override := range rateLimiter.rateLimiterProto.GetOverrides() {
				label := rateLimiter.rateLimiterProto.GetLabelKey() + ":" + override.GetLabelValue()
				rateLimiter.limitCheck.AddOverride(label, override.GetLimitScaleFactor())
			}
			rateLimiter.limiter, err = ratelimiter.NewOlricRateLimiter(rateLimiter.limitCheck,
				rateLimiter.rateLimiterFactory.distCache,
				rateLimiter.metricID,
				rateLimiter.rateLimiterProto.GetLimitResetInterval().AsDuration())
			if err != nil {
				log.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			// check whether lazy limiter is enabled
			if lazySyncConfig := rateLimiter.rateLimiterProto.GetLazySyncConfig(); lazySyncConfig != nil {
				if lazySyncConfig.GetEnabled() {
					lazySyncInterval := time.Duration(int64(rateLimiter.rateLimiterProto.GetLimitResetInterval().AsDuration()) / int64(lazySyncConfig.GetNumSync()))
					rateLimiter.limiter, err = ratelimiter.NewLazySyncRateLimiter(rateLimiter.limiter,
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
			rateLimiter.limiter.Close()
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
func (rateLimiter *rateLimiter) GetSelector() *policylangv1.Selector {
	return rateLimiter.rateLimiterProto.GetSelector()
}

// RunLimiter runs the limiter.
func (rateLimiter *rateLimiter) RunLimiter(labels selectors.Labels) *flowcontrolv1.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED

	label, ok, _, _ := rateLimiter.TakeN(labels, 1)

	if label == "" {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
	}

	return &flowcontrolv1.LimiterDecision{
		Decision: &flowcontrolv1.LimiterDecision_RateLimiterDecision_{
			RateLimiterDecision: &flowcontrolv1.LimiterDecision_RateLimiterDecision{
				PolicyName:     rateLimiter.GetPolicyName(),
				PolicyHash:     rateLimiter.GetPolicyHash(),
				ComponentIndex: rateLimiter.GetComponentIndex(),
			},
		},
		Dropped: !ok,
		Reason:  reason,
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

	ok, remaining, current = rateLimiter.limiter.TakeN(label, n)

	return
}

func (rateLimiter *rateLimiter) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	if event.Type == notifiers.Remove {
		log.Debug().Msg("Decision removed")
		rateLimiter.limitCheck.SetRateLimit(-1)
		return
	}

	var wrapperMessage configv1.ConfigPropertiesWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Config == nil {
		return
	}
	if wrapperMessage.PolicyHash != rateLimiter.ComponentAPI.GetPolicyHash() {
		return
	}
	limitDecision := &policydecisionsv1.RateLimiterDecision{}
	err = wrapperMessage.Config.UnmarshalTo(limitDecision)
	if err != nil {
		return
	}
	rateLimiter.limitCheck.SetRateLimit(int(limitDecision.GetLimit()))
}
