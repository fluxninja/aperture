package podscaler

import (
	"context"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policyprivatev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/private/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/v2/pkg/etcd/writer"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

// ScaleActuator struct.
type ScaleActuator struct {
	policyReadAPI        iface.Policy
	decisionWriter       *etcdwriter.Writer
	scaleActuatorProto   *policyprivatev1.PodScaleActuator
	decisionsEtcdPath    string
	agentGroupName       string
	podScalerComponentID string
}

// Name implements runtime.Component.
func (*ScaleActuator) Name() string { return "ScaleActuator" }

// Type implements runtime.Component.
func (*ScaleActuator) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (sa *ScaleActuator) ShortDescription() string { return sa.agentGroupName }

// IsActuator implements runtime.Component.
func (*ScaleActuator) IsActuator() bool { return true }

// NewScaleActuatorAndOptions creates scale actuator and its fx options.
func NewScaleActuatorAndOptions(
	scaleActuatorProto *policyprivatev1.PodScaleActuator,
	_ runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	agentGroup := scaleActuatorProto.GetAgentGroup()
	podScalerComponentID := scaleActuatorProto.GetPodScalerComponentId()
	etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), podScalerComponentID)
	decisionsEtcdPath := path.Join(paths.PodScalerDecisionsPath, etcdKey)
	sa := &ScaleActuator{
		policyReadAPI:        policyReadAPI,
		agentGroupName:       agentGroup,
		podScalerComponentID: podScalerComponentID,
		decisionsEtcdPath:    decisionsEtcdPath,
		scaleActuatorProto:   scaleActuatorProto,
	}

	return sa, fx.Options(
		fx.Invoke(sa.setupWriter),
	), nil
}

func (sa *ScaleActuator) setupWriter(scopedKV *etcdclient.SessionScopedKV, lifecycle fx.Lifecycle) error {
	logger := sa.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			sa.decisionWriter = etcdwriter.NewWriter(&scopedKV.KVWrapper)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			sa.decisionWriter.Close()
			_, err = scopedKV.Delete(clientv3.WithRequireLeader(ctx), sa.decisionsEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete scale decisions")
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (sa *ScaleActuator) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	// Get the decision from the port
	replicasReading := inPortReadings.ReadSingleReadingPort("replicas")
	var replicasValue float64
	if replicasReading.Valid() {
		if replicasReading.Value() <= 0 {
			replicasValue = 0
		} else {
			replicasValue = replicasReading.Value()
		}

		return nil, sa.publishDecision(replicasValue)
	}
	sa.publishDefaultDecision()
	return nil, nil
}

// DynamicConfigUpdate implements runtime.Component.DynamicConfigUpdate.
func (sa *ScaleActuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (sa *ScaleActuator) publishDefaultDecision() {
	// delete the decision
	sa.decisionWriter.Delete(sa.decisionsEtcdPath)
}

func (sa *ScaleActuator) publishDecision(desiredReplicas float64) error {
	logger := sa.policyReadAPI.GetStatusRegistry().GetLogger()
	// Save desired replicas in decision message
	decision := &policysyncv1.ScaleDecision{
		DesiredReplicas: int32(desiredReplicas),
	}
	// Publish decision
	logger.Autosample().Debug().Float64("desiredReplicas", desiredReplicas).Msg("Publish scale decision")
	wrapper := &policysyncv1.ScaleDecisionWrapper{
		ScaleDecision: decision,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:  sa.policyReadAPI.GetPolicyName(),
			PolicyHash:  sa.policyReadAPI.GetPolicyHash(),
			ComponentId: sa.podScalerComponentID,
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal policy decision")
		return err
	}
	sa.decisionWriter.Write(sa.decisionsEtcdPath, dat)
	return nil
}
