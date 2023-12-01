package podscaler

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
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

// ScaleActuator struct.
type ScaleActuator struct {
	policyReadAPI        iface.Policy
	etcdClient           *etcdclient.Client
	scaleActuatorProto   *policyprivatev1.PodScaleActuator
	decisionsEtcdPath    string
	agentGroupName       string
	podScalerComponentID string
}

// Name implements runtime.Component.
func (*ScaleActuator) Name() string { return "PodScaleActuator" }

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

func (sa *ScaleActuator) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	sa.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(_ context.Context) error {
			etcdClient.Delete(sa.decisionsEtcdPath)
			return nil
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (sa *ScaleActuator) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
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
	sa.etcdClient.Delete(sa.decisionsEtcdPath)
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
	sa.etcdClient.Put(sa.decisionsEtcdPath, string(dat))
	return nil
}
