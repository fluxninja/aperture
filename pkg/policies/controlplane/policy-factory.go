package controlplane

import (
	"context"
	"encoding/json"

	"go.uber.org/fx"
	"sigs.k8s.io/yaml"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/prometheus"
	"github.com/fluxninja/aperture/pkg/status"
)

// policyFactoryModule module for policy factory.
func policyFactoryModule() fx.Option {
	return fx.Options(
		fx.Invoke(
			fx.Annotate(
				setupPolicyFxDriver,
				fx.ParamTags(
					config.NameTag(policiesFxTag),
					config.NameTag(policiesDynamicConfigFxTag),
					common.FxOptionsFuncTag,
				),
			),
		),
		prometheus.Module(),
		policyModule(),
	)
}

type policyFactory struct {
	circuitJobGroup      *jobs.JobGroup
	etcdClient           *etcdclient.Client
	registry             status.Registry
	dynamicConfigWatcher notifiers.Watcher
}

// Main fx app.
func setupPolicyFxDriver(
	etcdWatcher notifiers.Watcher,
	dynamicConfigWatcher notifiers.Watcher,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
	registry status.Registry,
) error {
	wrapPolicy := func(key notifiers.Key, bytes []byte, etype notifiers.EventType) (notifiers.Key, []byte, error) {
		var dat []byte
		switch etype {
		case notifiers.Write:
			policyMessage := &policylangv1.Policy{}
			unmarshalErr := config.UnmarshalYAML(bytes, policyMessage)
			if unmarshalErr != nil {
				log.Warn().Err(unmarshalErr).Msg("Failed to unmarshal policy")
				return key, nil, unmarshalErr
			}

			wrapper, wrapErr := hashAndPolicyWrap(policyMessage, string(key))
			if wrapErr != nil {
				log.Warn().Err(wrapErr).Msg("Failed to wrap message in config properties")
				return key, nil, wrapErr
			}
			var marshalWrapErr error
			jsonDat, marshalWrapErr := json.Marshal(wrapper)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return key, nil, marshalWrapErr
			}
			// convert to yaml
			dat, marshalWrapErr = yaml.JSONToYAML(jsonDat)
			if marshalWrapErr != nil {
				log.Warn().Err(marshalWrapErr).Msgf("Failed to marshal config wrapper for proto message %+v", &wrapper)
				return key, nil, marshalWrapErr
			}
		}
		return key, dat, nil
	}

	policiesStatusRegistry := registry.Child(iface.PoliciesRoot)
	logger := policiesStatusRegistry.GetLogger()

	circuitJobGroup, err := jobs.NewJobGroup(policiesStatusRegistry.Child("circuit_jobs"), 0, jobs.RescheduleMode, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create job group")
		return err
	}

	factory := &policyFactory{
		registry:             policiesStatusRegistry,
		circuitJobGroup:      circuitJobGroup,
		etcdClient:           etcdClient,
		dynamicConfigWatcher: dynamicConfigWatcher,
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

	// content transform callback to wrap policy in config properties wrapper
	fxDriver.SetTransformFunc(wrapPolicy)

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
	return nil
}

// provideControllerPolicyFxOptions Per policy fx app in controller.
func (factory *policyFactory) provideControllerPolicyFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry status.Registry,
) (fx.Option, error) {
	var wrapperMessage wrappersv1.PolicyWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Policy == nil {
		registry.SetStatus(status.NewStatus(nil, err))
		registry.GetLogger().Warn().Err(err).Msg("Failed to unmarshal policy config wrapper")
		return fx.Options(), err
	}

	// save policy wrapper proto in status registry
	registry.Child("policy_config").SetStatus(status.NewStatus(&wrapperMessage, nil))

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
			fx.Annotate(factory.dynamicConfigWatcher, fx.As(new(notifiers.Watcher))),
			factory.circuitJobGroup,
			factory.etcdClient,
		),
		policyFxOptions,
	), nil
}
