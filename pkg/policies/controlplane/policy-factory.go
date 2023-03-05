package controlplane

import (
	"context"
	"sync"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	prom "github.com/fluxninja/aperture/pkg/prometheus"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/prometheus/client_golang/prometheus"
)

// policyFactoryModule module for policy factory.
func policyFactoryModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				providePolicyFactory,
				fx.ParamTags(
					config.NameTag(policiesFxTag),
					config.NameTag(policiesDynamicConfigFxTag),
					iface.FxOptionsFuncTag,
					alerts.AlertsFxTag,
				),
			),
		),
		grpcgateway.RegisterHandler{Handler: policylangv1.RegisterPolicyServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(RegisterPolicyService),
		prom.Module(),
		policyModule(),
	)
}

// PolicyFactory factory for policies.
type PolicyFactory struct {
	lock                 sync.RWMutex
	circuitJobGroup      *jobs.JobGroup
	etcdClient           *etcdclient.Client
	alerterIface         alerts.Alerter
	registry             status.Registry
	dynamicConfigWatcher notifiers.Watcher
	policyTracker        map[string]*policysyncv1.PolicyWrapper
}

// Main fx app.
func providePolicyFactory(
	etcdWatcher notifiers.Watcher,
	dynamicConfigWatcher notifiers.Watcher,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	alerterIface alerts.Alerter,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
	registry status.Registry,
	prometheusRegistry *prometheus.Registry,
) (*PolicyFactory, error) {
	policiesStatusRegistry := registry.Child("system", iface.PoliciesRoot)
	logger := policiesStatusRegistry.GetLogger()

	circuitJobGroup, err := jobs.NewJobGroup(policiesStatusRegistry.Child("jg", "circuit_jobs"), jobs.JobGroupConfig{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create job group")
		return nil, err
	}

	factory := &PolicyFactory{
		registry:             policiesStatusRegistry,
		circuitJobGroup:      circuitJobGroup,
		etcdClient:           etcdClient,
		alerterIface:         alerterIface,
		dynamicConfigWatcher: dynamicConfigWatcher,
		policyTracker:        make(map[string]*policysyncv1.PolicyWrapper),
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

	notifiers.NotifierLifecycle(lifecycle, etcdWatcher, fxDriver)
	return factory, nil
}

// provideControllerPolicyFxOptions Per policy fx app in controller.
func (factory *PolicyFactory) provideControllerPolicyFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry status.Registry,
) (fx.Option, error) {
	policyMessage := &policylangv1.Policy{}
	err := unmarshaller.Unmarshal(policyMessage)
	if err != nil {
		registry.SetStatus(status.NewStatus(nil, err))
		registry.GetLogger().Error().Err(err).Msg("Failed to unmarshal policy")
		return fx.Options(), err
	}

	wrapperMessage, err := hashAndPolicyWrap(policyMessage, string(key))
	if err != nil {
		registry.SetStatus(status.NewStatus(nil, err))
		registry.GetLogger().Error().Err(err).Msg("Failed to wrap message in config properties")
		return fx.Options(), err
	}

	policyFxOptions, err := newPolicyOptions(
		wrapperMessage,
		registry,
	)
	if err != nil {
		registry.SetStatus(status.NewStatus(nil, err))
		registry.GetLogger().Warn().Err(err).Msg("Failed to create policy options")
		return fx.Options(), err
	}
	return fx.Options(
		fx.Supply(
			fx.Annotate(factory.dynamicConfigWatcher, fx.As(new(notifiers.Watcher))),
			factory.circuitJobGroup,
			factory.etcdClient,
			factory.alerterIface,
			wrapperMessage,
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
