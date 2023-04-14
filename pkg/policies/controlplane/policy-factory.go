package controlplane

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	prom "github.com/fluxninja/aperture/pkg/prometheus"
	"github.com/fluxninja/aperture/pkg/status"
)

// Fx tag to match etcd watcher name.
var (
	policiesEtcdWatcherFxTag              = "policies-driver"
	policiesDynamicConfigEtcdWarcherFxTag = "policies-dynamic-config-driver"
)

// policyFactoryModule module for policy factory.
func policyFactoryModule() fx.Option {
	return fx.Options(
		etcdwatcher.Constructor{Name: policiesEtcdWatcherFxTag, EtcdPath: paths.PoliciesConfigPath}.Annotate(),
		etcdwatcher.Constructor{Name: policiesDynamicConfigEtcdWarcherFxTag, EtcdPath: paths.PoliciesDynamicConfigPath}.Annotate(),
		fx.Provide(
			fx.Annotate(
				providePolicyFactory,
				fx.ParamTags(
					config.NameTag(policiesEtcdWatcherFxTag),
					config.NameTag(policiesDynamicConfigEtcdWarcherFxTag),
					iface.FxOptionsFuncTag,
					alerts.AlertsFxTag,
				),
			),
		),
		grpcgateway.RegisterHandler{Handler: policylangv1.RegisterPolicyServiceHandlerFromEndpoint}.Annotate(),
		fx.Provide(
			fx.Annotate(
				RegisterPolicyService,
				fx.ParamTags(
					config.NameTag(policiesTrackerFxTag),
					config.NameTag(policiesDynamicConfigTrackerFxTag),
				),
			),
		),
		prom.Module(),
		policyModule(),
	)
}

// PolicyFactory factory for policies.
type PolicyFactory struct {
	lock                             sync.RWMutex
	circuitJobGroup                  *jobs.JobGroup
	etcdClient                       *etcdclient.Client
	alerterIface                     alerts.Alerter
	registry                         status.Registry
	policiesDynamicConfigEtcdWatcher notifiers.Watcher
	policyTracker                    map[string]*policysyncv1.PolicyWrapper
}

// Main fx app.
func providePolicyFactory(
	policiesEtcdWatcher notifiers.Watcher,
	policiesDynamicConfigEtcdWatcher notifiers.Watcher,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	alerterIface alerts.Alerter,
	etcdClient *etcdclient.Client,
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
			fx.Annotate(factory.policiesDynamicConfigEtcdWatcher, fx.As(new(notifiers.Watcher))),
			factory.circuitJobGroup,
			factory.etcdClient,
			factory.alerterIface,
			&wrapperMessage,
		),
		policyFxOptions,
		fx.Invoke(factory.trackPolicy),
	), nil
}

func (factory *PolicyFactory) trackPolicy(wrapperMessage *policysyncv1.PolicyWrapper, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			factory.lock.Lock()
			defer factory.lock.Unlock()
			factory.policyTracker[wrapperMessage.GetCommonAttributes().GetPolicyName()] = wrapperMessage
			return nil
		},
		OnStop: func(context.Context) error {
			factory.lock.Lock()
			defer factory.lock.Unlock()
			delete(factory.policyTracker, wrapperMessage.GetCommonAttributes().GetPolicyName())
			return nil
		},
	})
}

// GetPolicyWrappers returns all policy wrappers.
func (factory *PolicyFactory) GetPolicyWrappers() map[string]*policysyncv1.PolicyWrapper {
	factory.lock.RLock()
	defer factory.lock.RUnlock()
	// deepcopy wrappers
	policyWrappers := make(map[string]*policysyncv1.PolicyWrapper)
	for k, v := range factory.policyTracker {
		policyWrappers[k] = v.DeepCopy()
	}
	return policyWrappers
}

// GetPolicies returns all policies.
func (factory *PolicyFactory) GetPolicies() *policylangv1.Policies {
	policyWrappers := factory.GetPolicyWrappers()
	policies := make(map[string]*policylangv1.Policy)
	for _, v := range policyWrappers {
		policies[v.GetCommonAttributes().GetPolicyName()] = v.GetPolicy().DeepCopy()
	}
	return &policylangv1.Policies{
		Policies: policies,
	}
}

// GetPolicy returns policy matching given name.
func (factory *PolicyFactory) GetPolicy(name string) *policylangv1.Policy {
	policyWrappers := factory.GetPolicyWrappers()
	var policy *policylangv1.Policy
	for _, v := range policyWrappers {
		if v.GetCommonAttributes().GetPolicyName() == name {
			policy = v.GetPolicy().DeepCopy()
			break
		}
	}
	return policy
}
