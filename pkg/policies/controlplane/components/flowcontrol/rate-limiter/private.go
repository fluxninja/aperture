package ratelimiter

import (
	"context"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policyprivatev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/private/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime/tristate"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

type rateLimiter struct {
	policyReadAPI     iface.Policy
	rateLimiterProto  *policyprivatev1.RateLimiter
	decision          *policysyncv1.RateLimiterDecision
	etcdClient        *etcdclient.Client
	componentID       string
	decisionEtcdPaths []string
}

// Name implements runtime.Component.
func (*rateLimiter) Name() string { return "RateLimiter" }

// Type implements runtime.Component.
func (*rateLimiter) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (rl *rateLimiter) ShortDescription() string {
	return iface.GetSelectorsShortDescription(rl.rateLimiterProto.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*rateLimiter) IsActuator() bool { return true }

// NewRateLimiterAndOptions creates fx options for RateLimiter.
func NewRateLimiterAndOptions(
	rateLimiterProto *policyprivatev1.RateLimiter,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := rateLimiterProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), rateLimiterProto.GetParentComponentId())
		decisionEtcdPath := path.Join(paths.RateLimiterDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	rl := &rateLimiter{
		rateLimiterProto:  rateLimiterProto,
		decision:          &policysyncv1.RateLimiterDecision{},
		policyReadAPI:     policyReadAPI,
		decisionEtcdPaths: decisionEtcdPaths,
		componentID:       rateLimiterProto.GetParentComponentId(),
	}
	return rl, fx.Options(fx.Invoke(rl.setup)), nil
}

func (rl *rateLimiter) setup(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	rl.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			deleteEtcdPath := func(paths []string) {
				for _, path := range paths {
					etcdClient.Delete(path)
				}
			}
			deleteEtcdPath(rl.decisionEtcdPaths)
			return nil
		},
	})
}

// Execute implements runtime.Component.Execute.
func (rl *rateLimiter) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	bucketCapacity := inPortReadings.ReadSingleReadingPort("bucket_capacity")
	fillAmount := inPortReadings.ReadSingleReadingPort("fill_amount")
	passThrough := tristate.FromReading(inPortReadings.ReadSingleReadingPort("pass_through"))

	decision := &policysyncv1.RateLimiterDecision{
		PassThrough: true,
	}
	if !bucketCapacity.Valid() || !fillAmount.Valid() {
		return nil, rl.publishDecision(decision)
	}

	decision.BucketCapacity = bucketCapacity.Value()
	decision.FillAmount = fillAmount.Value()
	decision.PassThrough = passThrough.IsTrue()

	return nil, rl.publishDecision(decision)
}

func (rl *rateLimiter) publishDecision(decision *policysyncv1.RateLimiterDecision) error {
	logger := rl.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if !proto.Equal(rl.decision, decision) {
		rl.decision = decision
		// Publish decision
		logger.Debug().Msg("publishing rate limiter decision")
		wrapper := &policysyncv1.RateLimiterDecisionWrapper{
			RateLimiterDecision: decision,
			CommonAttributes: &policysyncv1.CommonAttributes{
				PolicyName:  rl.policyReadAPI.GetPolicyName(),
				PolicyHash:  rl.policyReadAPI.GetPolicyHash(),
				ComponentId: rl.componentID,
			},
		}

		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("failed to marshal rate limiter decision")
			return err
		}
		for _, decisionEtcdPath := range rl.decisionEtcdPaths {
			rl.etcdClient.Put(decisionEtcdPath, string(dat))
		}
	}
	return nil
}

// DynamicConfigUpdate is a no-op for rate limiter.
func (rl *rateLimiter) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
