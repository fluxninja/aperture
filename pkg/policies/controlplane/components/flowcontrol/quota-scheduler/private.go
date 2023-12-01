package quotascheduler

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

type quotaScheduler struct {
	policyReadAPI       iface.Policy
	quotaSchedulerProto *policyprivatev1.QuotaScheduler
	decision            *policysyncv1.RateLimiterDecision
	etcdClient          *etcdclient.Client
	componentID         string
	decisionEtcdPaths   []string
}

// Name implements runtime.Component.
func (*quotaScheduler) Name() string { return "QuotaScheduler" }

// Type implements runtime.Component.
func (*quotaScheduler) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (qs *quotaScheduler) ShortDescription() string {
	return iface.GetSelectorsShortDescription(qs.quotaSchedulerProto.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*quotaScheduler) IsActuator() bool { return true }

// NewQuotaSchedulerAndOptions creates fx options for QuotaScheduler.
func NewQuotaSchedulerAndOptions(
	quotaSchedulerProto *policyprivatev1.QuotaScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := quotaSchedulerProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), quotaSchedulerProto.GetParentComponentId())
		decisionEtcdPath := path.Join(paths.QuotaSchedulerDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	qs := &quotaScheduler{
		quotaSchedulerProto: quotaSchedulerProto,
		decision:            &policysyncv1.RateLimiterDecision{},
		policyReadAPI:       policyReadAPI,
		decisionEtcdPaths:   decisionEtcdPaths,
		componentID:         quotaSchedulerProto.GetParentComponentId(),
	}
	return qs, fx.Options(fx.Invoke(qs.setup)), nil
}

func (qs *quotaScheduler) setup(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	qs.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			deleteEtcdPath := func(paths []string) {
				for _, path := range paths {
					etcdClient.Delete(path)
				}
			}
			deleteEtcdPath(qs.decisionEtcdPaths)
			return nil
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (qs *quotaScheduler) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	bucketCapacity := inPortReadings.ReadSingleReadingPort("bucket_capacity")
	fillAmount := inPortReadings.ReadSingleReadingPort("fill_amount")
	ptBool := tristate.FromReading(inPortReadings.ReadSingleReadingPort("pass_through"))

	decision := &policysyncv1.RateLimiterDecision{
		PassThrough: true,
	}
	if !bucketCapacity.Valid() || !fillAmount.Valid() {
		return nil, qs.publishDecision(decision)
	}

	decision.BucketCapacity = bucketCapacity.Value()
	decision.FillAmount = fillAmount.Value()
	decision.PassThrough = ptBool.IsTrue()

	return nil, qs.publishDecision(decision)
}

func (qs *quotaScheduler) publishDecision(decision *policysyncv1.RateLimiterDecision) error {
	logger := qs.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if !proto.Equal(qs.decision, decision) {
		qs.decision = decision
		// Publish decision
		logger.Debug().Msg("publishing rate limiter decision")
		wrapper := &policysyncv1.RateLimiterDecisionWrapper{
			RateLimiterDecision: decision,
			CommonAttributes: &policysyncv1.CommonAttributes{
				PolicyName:  qs.policyReadAPI.GetPolicyName(),
				PolicyHash:  qs.policyReadAPI.GetPolicyHash(),
				ComponentId: qs.componentID,
			},
		}

		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("failed to marshal rate limiter decision")
			return err
		}
		for _, decisionEtcdPath := range qs.decisionEtcdPaths {
			qs.etcdClient.Put(decisionEtcdPath, string(dat))
		}
	}
	return nil
}

// DynamicConfigUpdate is a no-op for rate limiter.
func (qs *quotaScheduler) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
