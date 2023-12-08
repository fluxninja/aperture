package concurrencyscheduler

import (
	"context"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policyprivatev1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/private/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime/tristate"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

type concurrencyScheduler struct {
	policyReadAPI             iface.Policy
	concurrencySchedulerProto *policyprivatev1.ConcurrencyScheduler
	decision                  *policysyncv1.ConcurrencyLimiterDecision
	etcdClient                *etcdclient.Client
	componentID               string
	decisionEtcdPaths         []string
}

// Name implements runtime.Component.
func (*concurrencyScheduler) Name() string { return "ConcurrencyScheduler" }

// Type implements runtime.Component.
func (*concurrencyScheduler) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (cs *concurrencyScheduler) ShortDescription() string {
	return iface.GetSelectorsShortDescription(cs.concurrencySchedulerProto.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*concurrencyScheduler) IsActuator() bool { return true }

// NewConcurrencySchedulerAndOptions creates fx options for ConcurrencyScheduler.
func NewConcurrencySchedulerAndOptions(
	concurrencySchedulerProto *policyprivatev1.ConcurrencyScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := concurrencySchedulerProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), concurrencySchedulerProto.GetParentComponentId())
		decisionEtcdPath := path.Join(paths.ConcurrencySchedulerDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	cs := &concurrencyScheduler{
		concurrencySchedulerProto: concurrencySchedulerProto,
		decision:                  &policysyncv1.ConcurrencyLimiterDecision{},
		policyReadAPI:             policyReadAPI,
		decisionEtcdPaths:         decisionEtcdPaths,
		componentID:               concurrencySchedulerProto.GetParentComponentId(),
	}
	return cs, fx.Options(fx.Invoke(cs.setup)), nil
}

func (cs *concurrencyScheduler) setup(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	cs.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			deleteEtcdPath := func(paths []string) {
				for _, path := range paths {
					etcdClient.Delete(path)
				}
			}
			deleteEtcdPath(cs.decisionEtcdPaths)
			return nil
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (cs *concurrencyScheduler) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	maxConcurrency := inPortReadings.ReadSingleReadingPort("max_concurrency")
	ptBool := tristate.FromReading(inPortReadings.ReadSingleReadingPort("pass_through"))

	decision := &policysyncv1.ConcurrencyLimiterDecision{
		PassThrough: true,
	}
	if !maxConcurrency.Valid() {
		return nil, cs.publishDecision(decision)
	}

	decision.MaxConcurrency = maxConcurrency.Value()
	decision.PassThrough = ptBool.IsTrue()

	return nil, cs.publishDecision(decision)
}

func (cs *concurrencyScheduler) publishDecision(decision *policysyncv1.ConcurrencyLimiterDecision) error {
	logger := cs.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if !proto.Equal(cs.decision, decision) {
		cs.decision = decision
		// Publish decision
		logger.Debug().Msg("publishing concurrency limiter decision")
		wrapper := &policysyncv1.ConcurrencyLimiterDecisionWrapper{
			ConcurrencyLimiterDecision: decision,
			CommonAttributes: &policysyncv1.CommonAttributes{
				PolicyName:  cs.policyReadAPI.GetPolicyName(),
				PolicyHash:  cs.policyReadAPI.GetPolicyHash(),
				ComponentId: cs.componentID,
			},
		}

		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("failed to marshal concurrency limiter decision")
			return err
		}
		for _, decisionEtcdPath := range cs.decisionEtcdPaths {
			cs.etcdClient.Put(decisionEtcdPath, string(dat))
		}
	}
	return nil
}

// DynamicConfigUpdate is a no-op for concurrency limiter.
func (cs *concurrencyScheduler) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
