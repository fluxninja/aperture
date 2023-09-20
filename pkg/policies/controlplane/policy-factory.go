package controlplane

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	googletoken "github.com/fluxninja/aperture/v2/pkg/google"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	prom "github.com/fluxninja/aperture/v2/pkg/prometheus"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// Fx tag to match etcd watcher name.
var (
	policiesEtcdWatcherFxTag              = "policies-driver"
	policiesDynamicConfigEtcdWatcherFxTag = "policies-dynamic-config-driver"
)

// policyFactoryModule module for policy factory.
func policyFactoryModule() fx.Option {
	return fx.Options(
		etcdwatcher.Constructor{Name: policiesEtcdWatcherFxTag, EtcdPath: paths.PoliciesConfigPath}.Annotate(),
		etcdwatcher.Constructor{Name: policiesDynamicConfigEtcdWatcherFxTag, EtcdPath: paths.PoliciesDynamicConfigPath}.Annotate(),
		fx.Provide(
			fx.Annotate(
				providePolicyFactory,
				fx.ParamTags(
					config.NameTag(policiesEtcdWatcherFxTag),
					config.NameTag(policiesDynamicConfigEtcdWatcherFxTag),
					iface.FxOptionsFuncTag,
					alerts.AlertsFxTag,
				),
			),
		),
		grpcgateway.RegisterHandler{Handler: policylangv1.RegisterPolicyServiceHandlerFromEndpoint}.Annotate(),
		fx.Provide(
			fx.Annotate(
				RegisterPolicyService,
			),
		),
		prom.Module(),
		googletoken.Module(),
		policyModule(),
	)
}

// PolicyFactory factory for policies.
type PolicyFactory struct {
	lock                             sync.RWMutex
	circuitJobGroup                  *jobs.JobGroup
	etcdClient                       *etcdclient.Client
	sessionScopedKV                  *etcdclient.SessionScopedKV
	prometheusEnforcer               *prom.PrometheusEnforcer
	alerterIface                     alerts.Alerter
	registry                         status.Registry
	policiesDynamicConfigEtcdWatcher notifiers.Watcher
	policyTracker                    map[string]*policysyncv1.PolicyWrapper // keyed by wrapper.CommonAttributes.PolicyName
}

// Main fx app.
func providePolicyFactory(
	policiesEtcdWatcher notifiers.Watcher,
	policiesDynamicConfigEtcdWatcher notifiers.Watcher,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	alerterIface alerts.Alerter,
	etcdClient *etcdclient.Client,
	sessionScopedKV *etcdclient.SessionScopedKV,
	enforcer *prom.PrometheusEnforcer,
	lifecycle fx.Lifecycle,
	registry status.Registry,
	prometheusRegistry *prometheus.Registry,
) (*PolicyFactory, error) {
	policiesStatusRegistry := registry.Child("system", iface.PoliciesRoot)
	logger := policiesStatusRegistry.GetLogger()

	circuitJobGroup, err := jobs.NewJobGroup(policiesStatusRegistry.Child("job-group", "circuit_jobs"), jobs.JobGroupConfig{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create job group")
		return nil, err
	}

	factory := &PolicyFactory{
		registry:                         policiesStatusRegistry,
		circuitJobGroup:                  circuitJobGroup,
		etcdClient:                       etcdClient,
		sessionScopedKV:                  sessionScopedKV,
		prometheusEnforcer:               enforcer,
		alerterIface:                     alerterIface,
		policiesDynamicConfigEtcdWatcher: policiesDynamicConfigEtcdWatcher,
		policyTracker:                    make(map[string]*policysyncv1.PolicyWrapper),
	}

	optionsFunc := []notifiers.FxOptionsFunc{factory.provideControllerPolicyFxOptions}
	if len(fxOptionsFuncs) > 0 {
		optionsFunc = append(optionsFunc, fxOptionsFuncs...)
	}

	fxDriver, err := notifiers.NewFxDriver(
		policiesStatusRegistry,
		prometheusRegistry,
		config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
		optionsFunc,
	)
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := factory.circuitJobGroup.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			defer policiesStatusRegistry.Detach()
			err := factory.circuitJobGroup.Stop()
			if err != nil {
				return err
			}
			return nil
		},
	})

	notifiers.NotifierLifecycle(lifecycle, policiesEtcdWatcher, fxDriver)
	return factory, nil
}

// provideControllerPolicyFxOptions Per policy fx app in controller.
func (factory *PolicyFactory) provideControllerPolicyFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry status.Registry,
) (fx.Option, error) {
	var wrapperMessage policysyncv1.PolicyWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Policy == nil {
		registry.SetStatus(status.NewStatus(nil, err))
		registry.GetLogger().Error().Err(err).Msg("Failed to unmarshal policy config wrapper")
		return fx.Options(), err
	}

	policyFxOptions, err := newPolicyOptions(
		&wrapperMessage,
		registry,
	)
	if err != nil {
		registry.SetStatus(status.NewStatus(nil, err))
		registry.GetLogger().Warn().Err(err).Msg("Failed to create policy options")
		return fx.Options(), err
	}
	return fx.Options(
		fx.Supply(
			fx.Annotate(
				factory.policiesDynamicConfigEtcdWatcher,
				fx.As(new(notifiers.Watcher)),
			),
			factory.circuitJobGroup,
			factory.etcdClient,
			factory.sessionScopedKV,
			factory.prometheusEnforcer,
			factory.alerterIface,
			&wrapperMessage,
		),
		policyFxOptions,
		fx.Invoke(factory.trackPolicy),
	), nil
}

func (factory *PolicyFactory) trackPolicy(wrapperMessage *policysyncv1.PolicyWrapper, lifecycle fx.Lifecycle) {
	policyName := wrapperMessage.GetCommonAttributes().GetPolicyName()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			factory.lock.Lock()
			defer factory.lock.Unlock()
			factory.policyTracker[policyName] = wrapperMessage
			return nil
		},
		OnStop: func(context.Context) error {
			factory.lock.Lock()
			defer factory.lock.Unlock()
			delete(factory.policyTracker, policyName)
			return nil
		},
	})
	lifecycle.Append(fx.StartStopHook(
		func() { log.Info().Str("policy", policyName).Msg("Policy loaded to controller") },
		func() { log.Info().Str("policy", policyName).Msg("Unloading policy from controller") },
	))
}

// GetPolicyWrappers returns all policy wrappers.
func (factory *PolicyFactory) GetPolicyWrappers() map[string]*policysyncv1.PolicyWrapper {
	factory.lock.RLock()
	defer factory.lock.RUnlock()
	// deepcopy wrappers
	policyWrappers := make(map[string]*policysyncv1.PolicyWrapper, len(factory.policyTracker))
	for k, v := range factory.policyTracker {
		policyWrappers[k] = proto.Clone(v).(*policysyncv1.PolicyWrapper)
	}
	return policyWrappers
}

// GetPolicyWrapper returns policy wrapper matching given name.
func (factory *PolicyFactory) GetPolicyWrapper(name string) *policysyncv1.PolicyWrapper {
	factory.lock.RLock()
	defer factory.lock.RUnlock()
	policyWrapper, exists := factory.policyTracker[name]
	if !exists {
		return nil
	}
	return proto.Clone(policyWrapper).(*policysyncv1.PolicyWrapper)
}
