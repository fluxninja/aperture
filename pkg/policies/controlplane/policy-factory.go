package controlplane

import (
	"context"

	"go.uber.org/fx"

	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/prometheus"
	"github.com/fluxninja/aperture/pkg/status"
)

// Fx tag to match etcd watcher name.
var policiesDriverFxTag = "policies-driver"

// policyFactoryModule module for policy factory.
func policyFactoryModule() fx.Option {
	return fx.Options(
		etcdwatcher.Constructor{Name: policiesDriverFxTag, EtcdPath: common.PoliciesConfigPath}.Annotate(),
		fx.Invoke(
			fx.Annotate(
				setupPolicyFxDriver,
				fx.ParamTags(
					config.NameTag(policiesDriverFxTag),
					common.FxOptionsFuncTag,
				),
			),
		),
		prometheus.Module(),
		policyModule(),
	)
}

type policyFactory struct {
	circuitJobGroup *jobs.JobGroup
	etcdClient      *etcdclient.Client
	registry        status.Registry
}

// Main fx app.
func setupPolicyFxDriver(
	etcdWatcher notifiers.Watcher,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
	registry status.Registry,
) error {
	policiesStatusRegistry := registry.Child(iface.PoliciesRoot)

	circuitJobGroup, err := jobs.NewJobGroup(policiesStatusRegistry.Child("circuit_jobs"), 0, jobs.RescheduleMode, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create job group")
		return err
	}

	factory := &policyFactory{
		registry:        policiesStatusRegistry,
		circuitJobGroup: circuitJobGroup,
		etcdClient:      etcdClient,
	}

	optionsFunc := []notifiers.FxOptionsFunc{factory.provideControllerPolicyFxOptions}
	if len(fxOptionsFuncs) > 0 {
		optionsFunc = append(optionsFunc, fxOptionsFuncs...)
	}

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: optionsFunc,
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
		},
		StatusRegistry: policiesStatusRegistry,
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
			err := factory.circuitJobGroup.Stop()
			if err != nil {
				return err
			}
			return nil
		},
	})

	notifiers.NotifierLifecycle(lifecycle, etcdWatcher, fxDriver)
	return nil
}

// provideControllerPolicyFxOptions Per policy fx app in controller.
func (factory *policyFactory) provideControllerPolicyFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	var wrapperMessage wrappersv1.PolicyWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Policy == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		log.Warn().Err(err).Msg("Failed to unmarshal policy config wrapper")
		return fx.Options(), err
	}

	// save policy wrapper proto in status registry
	reg.Child("policy_config").SetStatus(status.NewStatus(&wrapperMessage, nil))

	policyFxOptions, err := newPolicyOptions(
		&wrapperMessage,
	)
	if err != nil {
		reg.SetStatus(status.NewStatus(nil, err))
		log.Warn().Err(err).Msg("Failed to create policy options")
		return fx.Options(), err
	}
	return fx.Options(
		policyFxOptions,
		fx.Supply(
			factory.circuitJobGroup,
			factory.etcdClient,
		),
	), nil
}
