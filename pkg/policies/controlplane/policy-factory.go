package controlplane

import (
	"github.com/pkg/errors"
	"go.uber.org/fx"

	configv1 "aperture.tech/aperture/api/gen/proto/go/aperture/common/config/v1"
	policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"aperture.tech/aperture/pkg/config"
	etcdclient "aperture.tech/aperture/pkg/etcd/client"
	etcdwatcher "aperture.tech/aperture/pkg/etcd/watcher"
	"aperture.tech/aperture/pkg/jobs"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/notifiers"
	"aperture.tech/aperture/pkg/paths"
	"aperture.tech/aperture/pkg/policies/controlplane/common"
	"aperture.tech/aperture/pkg/prometheus"
	"aperture.tech/aperture/pkg/status"
)

var (
	// Path in status registry for policies results.
	policiesStatusRoot = "policies"
	// Fx tag to match etcd watcher name.
	policiesDriverFxTag = "policies-driver"
)

const (
	circuitJobGroupTag = "circuit-job-group"
)

// PolicyFactoryModule module for policy factory.
func PolicyFactoryModule() fx.Option {
	return fx.Options(
		jobs.JobGroupConstructor{Group: circuitJobGroupTag}.Annotate(),
		etcdwatcher.Constructor{Name: policiesDriverFxTag, EtcdPath: paths.Policies}.Annotate(),
		fx.Invoke(
			fx.Annotate(
				setupPolicyFxDriver,
				fx.ParamTags(
					config.NameTag(policiesDriverFxTag),
					config.NameTag(circuitJobGroupTag),
					common.FxOptionsFuncTag,
				),
			),
		),
		prometheus.Module(),
		PolicyModule(),
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
	circuitJobGroup *jobs.JobGroup,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
	registry *status.Registry,
) {
	factory := &policyFactory{
		registryPath:    policiesStatusRoot,
		circuitJobGroup: circuitJobGroup,
		etcdClient:      etcdClient,
	}

	optionsFunc := []notifiers.FxOptionsFunc{factory.ProvideControllerPolicyFxOptions}
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

	notifiers.NotifierLifecycle(lifecycle, etcdWatcher, fxDriver)
}

// ProvideControllerPolicyFxOptions Per policy fx app in controller.
func (factory *policyFactory) ProvideControllerPolicyFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry *status.Registry,
) (fx.Option, error) {
	var wrapperMessage configv1.ConfigPropertiesWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Config == nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(factory.registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal policy config wrapper")
		return fx.Options(), err
	}
	var policyMessage policylangv1.Policy
	err = wrapperMessage.Config.UnmarshalTo(&policyMessage)
	if err != nil {
		s := status.NewStatus(nil, err)
		_ = registry.Push(factory.registryPath, s)
		log.Warn().Err(err).Msg("Failed to unmarshal policy")
		return fx.Options(), err
	}
	policyFxOptions, err := NewPolicyOptions(
		factory.circuitJobGroup,
		factory.etcdClient,
		&wrapperMessage,
		&policyMessage,
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
	return policyFxOptions, nil
}
