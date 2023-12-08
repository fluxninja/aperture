package concurrencylimiter

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

type concurrencyLimiter struct {
	policyReadAPI           iface.Policy
	concurrencyLimiterProto *policyprivatev1.ConcurrencyLimiter
	decision                *policysyncv1.ConcurrencyLimiterDecision
	etcdClient              *etcdclient.Client
	componentID             string
	decisionEtcdPaths       []string
}

// Name implements runtime.Component.
func (*concurrencyLimiter) Name() string { return "ConcurrencyLimiter" }

// Type implements runtime.Component.
func (*concurrencyLimiter) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (cl *concurrencyLimiter) ShortDescription() string {
	return iface.GetSelectorsShortDescription(cl.concurrencyLimiterProto.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*concurrencyLimiter) IsActuator() bool { return true }

// NewConcurrencyLimiterAndOptions creates fx options for ConcurrencyLimiter.
func NewConcurrencyLimiterAndOptions(
	concurrencyLimiterProto *policyprivatev1.ConcurrencyLimiter,
	_ runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := concurrencyLimiterProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), concurrencyLimiterProto.GetParentComponentId())
		decisionEtcdPath := path.Join(paths.ConcurrencyLimiterDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	cl := &concurrencyLimiter{
		concurrencyLimiterProto: concurrencyLimiterProto,
		decision:                &policysyncv1.ConcurrencyLimiterDecision{},
		policyReadAPI:           policyReadAPI,
		decisionEtcdPaths:       decisionEtcdPaths,
		componentID:             concurrencyLimiterProto.GetParentComponentId(),
	}
	return cl, fx.Options(fx.Invoke(cl.setup)), nil
}

func (cl *concurrencyLimiter) setup(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	cl.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			deleteEtcdPath := func(paths []string) {
				for _, path := range paths {
					etcdClient.Delete(path)
				}
			}
			deleteEtcdPath(cl.decisionEtcdPaths)
			return nil
		},
	})
}

// Execute implements runtime.Component.Execute.
func (cl *concurrencyLimiter) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	maxConcurrency := inPortReadings.ReadSingleReadingPort("max_concurrency")
	passThrough := tristate.FromReading(inPortReadings.ReadSingleReadingPort("pass_through"))

	decision := &policysyncv1.ConcurrencyLimiterDecision{
		PassThrough: true,
	}
	if !maxConcurrency.Valid() {
		return nil, cl.publishDecision(decision)
	}

	decision.MaxConcurrency = maxConcurrency.Value()
	decision.PassThrough = passThrough.IsTrue()

	return nil, cl.publishDecision(decision)
}

func (cl *concurrencyLimiter) publishDecision(decision *policysyncv1.ConcurrencyLimiterDecision) error {
	logger := cl.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if !proto.Equal(cl.decision, decision) {
		cl.decision = decision
		// Publish decision
		logger.Debug().Msg("publishing concurrency limiter decision")
		wrapper := &policysyncv1.ConcurrencyLimiterDecisionWrapper{
			ConcurrencyLimiterDecision: decision,
			CommonAttributes: &policysyncv1.CommonAttributes{
				PolicyName:  cl.policyReadAPI.GetPolicyName(),
				PolicyHash:  cl.policyReadAPI.GetPolicyHash(),
				ComponentId: cl.componentID,
			},
		}

		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("failed to marshal concurrency limiter decision")
			return err
		}
		for _, decisionEtcdPath := range cl.decisionEtcdPaths {
			cl.etcdClient.Put(decisionEtcdPath, string(dat))
		}
	}
	return nil
}

// DynamicConfigUpdate is a no-op for concurrency limiter.
func (cl *concurrencyLimiter) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
