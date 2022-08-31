package controlplane

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
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

// Fx tag to match etcd watcher name.
var policiesDriverFxTag = "policies-driver"

// PolicyFactoryModule module for policy factory.
func PolicyFactoryModule() fx.Option {
	return fx.Options(
		fx.Provide(ProvidePolicyFactory),
		etcdwatcher.Constructor{Name: policiesDriverFxTag, EtcdPath: paths.PoliciesConfigPath}.Annotate(),
		fx.Invoke(
			fx.Annotate(
				setupPolicyFxDriver,
				fx.ParamTags(
					config.NameTag(policiesDriverFxTag),
					common.FxOptionsFuncTag,
				),
			),
			RegisterPolicyService,
		),
		prometheus.Module(),
		policyModule(),
	)
}

// PolicyFactory creates policies fx app per each policy.
type PolicyFactory struct {
	circuitJobGroup   *jobs.JobGroup
	etcdClient        *etcdclient.Client
	registry          status.Registry
	policyWrapperList []*configv1.PolicyWrapper
}

// ProvidePolicyFactory returns PolicyFactory pointer.
func ProvidePolicyFactory(
	etcdClient *etcdclient.Client,
	registry status.Registry,
) (*PolicyFactory, error) {
	policiesStatusRegistry := registry.Child(iface.PoliciesRoot)

	circuitJobGroup, err := jobs.NewJobGroup(policiesStatusRegistry.Child("circuit_jobs"), 0, jobs.RescheduleMode, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create job group")
		return nil, err
	}

	return &PolicyFactory{
		registry:        policiesStatusRegistry,
		circuitJobGroup: circuitJobGroup,
		etcdClient:      etcdClient,
	}, nil
}

// Main fx app.
func setupPolicyFxDriver(
	etcdWatcher notifiers.Watcher,
	fxOptionsFuncs []notifiers.FxOptionsFunc,
	lifecycle fx.Lifecycle,
	factory *PolicyFactory,
) error {
	optionsFunc := []notifiers.FxOptionsFunc{factory.provideControllerPolicyFxOptions}
	if len(fxOptionsFuncs) > 0 {
		optionsFunc = append(optionsFunc, fxOptionsFuncs...)
	}

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: optionsFunc,
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
		},
		StatusRegistry: factory.registry,
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
func (factory *PolicyFactory) provideControllerPolicyFxOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	var wrapperMessage configv1.PolicyWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Policy == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		log.Warn().Err(err).Msg("Failed to unmarshal policy config wrapper")
		return fx.Options(), err
	}

	// save policy wrapper proto in status registry
	reg.Child("policy_config").SetStatus(status.NewStatus(&wrapperMessage, nil))
	factory.policyWrapperList = append(factory.policyWrapperList, &wrapperMessage)

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

// GetPoliciesMap wraps AllPolicies rpc function.
func (factory *PolicyFactory) GetPoliciesMap() *policylangv1.AllPoliciesResponse {
	response, _ := factory.AllPolicies(context.Background(), nil)
	return response
}

// AllPolicies returns all of the policies.
func (factory *PolicyFactory) AllPolicies(ctx context.Context, _ *emptypb.Empty) (*policylangv1.AllPoliciesResponse, error) {
	allPoliciesMap := make(map[string]*policylangv1.Policy)
	for _, policy := range factory.policyWrapperList {
		allPoliciesMap[policy.PolicyName] = policy.Policy
	}

	return &policylangv1.AllPoliciesResponse{
		AllPolicies: allPoliciesMap,
	}, nil
}

// RegisterPolicyService registers a service for policies.
func RegisterPolicyService(server *grpc.Server, factory *PolicyFactory) {
	policylangv1.RegisterPolicyServiceServer(server, factory)
}
