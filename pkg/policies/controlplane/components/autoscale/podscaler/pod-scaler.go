package podscaler

import (
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/private/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	inputReplicasPortName            = "replicas"
	outputActualReplicasPortName     = "actual_replicas"
	outputConfiguredReplicasPortName = "configured_replicas"
)

// ParsePodScaler parses a PodScaler component and returns its nested circuit representation.
func ParsePodScaler(
	podScaler *policylangv1.PodScaler,
	componentID runtime.ComponentID,
	_ iface.Policy,
) (*policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := podScaler.GetInPorts()

	if inPorts != nil {
		replicasPort := inPorts.Replicas
		if replicasPort != nil {
			nestedInPortsMap[inputReplicasPortName] = replicasPort
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := podScaler.GetOutPorts()
	if outPorts != nil {
		actualReplicasPort := outPorts.ActualReplicas
		if actualReplicasPort != nil {
			nestedOutPortsMap[outputActualReplicasPortName] = actualReplicasPort
		}
		configuredReplicasPort := outPorts.ConfiguredReplicas
		if configuredReplicasPort != nil {
			nestedOutPortsMap[outputConfiguredReplicasPortName] = configuredReplicasPort
		}
	}

	podScaleActuatorAnyProto, err := anypb.New(
		&policyprivatev1.PodScaleActuator{
			InPorts: &policyprivatev1.PodScaleActuator_Ins{
				Replicas: &policylangv1.InPort{
					Value: &policylangv1.InPort_SignalName{
						SignalName: "REPLICAS",
					},
				},
			},
			DryRun:               podScaler.DryRun,
			DryRunConfigKey:      podScaler.DryRunConfigKey,
			PodScalerComponentId: componentID.String(),
			AgentGroup:           podScaler.KubernetesObjectSelector.AgentGroup,
		},
	)
	if err != nil {
		return nil, err
	}

	podScaleReporterAnyProto, err := anypb.New(
		&policyprivatev1.PodScaleReporter{
			OutPorts: &policyprivatev1.PodScaleReporter_Outs{
				ActualReplicas: &policylangv1.OutPort{
					SignalName: "ACTUAL_REPLICAS",
				},
				ConfiguredReplicas: &policylangv1.OutPort{
					SignalName: "CONFIGURED_REPLICAS",
				},
			},
			PodScalerComponentId: componentID.String(),
			AgentGroup:           podScaler.KubernetesObjectSelector.AgentGroup,
		},
	)
	if err != nil {
		return nil, err
	}

	kos := podScaler.KubernetesObjectSelector
	sd := fmt.Sprintf("%s/%s/%s/%s/%s",
		kos.GetAgentGroup(),
		kos.GetNamespace(),
		kos.GetApiVersion(),
		kos.GetKind(),
		kos.GetName(),
	)

	nestedCircuit := &policylangv1.NestedCircuit{
		Name:             "PodScaler",
		ShortDescription: sd,
		InPortsMap:       nestedInPortsMap,
		OutPortsMap:      nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_AutoScale{
					AutoScale: &policylangv1.AutoScale{
						Component: &policylangv1.AutoScale_Private{
							Private: podScaleActuatorAnyProto,
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_AutoScale{
					AutoScale: &policylangv1.AutoScale{
						Component: &policylangv1.AutoScale_Private{
							Private: podScaleReporterAnyProto,
						},
					},
				},
			},
		},
	}

	components.AddNestedIngress(nestedCircuit, inputReplicasPortName, "REPLICAS")
	components.AddNestedEgress(nestedCircuit, outputActualReplicasPortName, "ACTUAL_REPLICAS")
	components.AddNestedEgress(nestedCircuit, outputConfiguredReplicasPortName, "CONFIGURED_REPLICAS")

	return nestedCircuit, nil
}
