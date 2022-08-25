package controlplane

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/fx"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/prometheus"
	"github.com/fluxninja/aperture/pkg/status"
)

var (
	// Path in status registry for policies results.
	policiesStatusRoot = "policies"
	// Fx tag to match etcd watcher name.
	policiesDriverFxTag = "policies-driver"
)

// PolicyFactoryModule module for policy factory.
func PolicyFactoryModule() fx.Option {
	return fx.Options(
		etcdwatcher.Constructor{Name: policiesDriverFxTag, EtcdPath: paths.Policies}.Annotate(),
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
	registryPath    string
}

// Main fx app.
func setupPolicyFxDriver(
	etcdWatcher notifiers.Watcher,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
	registry *status.Registry,
) error {
	circuitJobGroup, err := jobs.NewJobGroup(iface.PoliciesRoot+".circuit_jobs", registry, 0, jobs.RescheduleMode, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create job group")
		return err
	}

	factory := &policyFactory{
		registryPath:    policiesStatusRoot,
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
		StatusRegistry: registry,
		StatusPath:     policiesStatusRoot,
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
	registry *status.Registry,
) (fx.Option, error) {
	var wrapperMessage configv1.PolicyWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Policy == nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(factory.registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal policy config wrapper")
		return fx.Options(), err
	}
	policyFxOptions, err := newPolicyOptions(
		&wrapperMessage,
	)
	if err != nil {
		s := status.NewStatus(nil, err)
		rPErr := registry.Push(factory.registryPath, s)
		if rPErr != nil {
			// Wrap errors
			err = errors.Wrap(err, rPErr.Error())
		}
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
