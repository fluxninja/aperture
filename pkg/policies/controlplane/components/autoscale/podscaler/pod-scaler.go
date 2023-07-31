package podscaler

import (
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/private/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
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
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	retErr := func(err error) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
		return nil, nil, err
	}

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
			PodScalerComponentId: componentID.String(),
			AgentGroup:           podScaler.KubernetesObjectSelector.AgentGroup,
		},
	)
	if err != nil {
		return retErr(err)
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
		return retErr(err)
	}

	nestedCircuit := &policylangv1.NestedCircuit{
		InPortsMap:  nestedInPortsMap,
		OutPortsMap: nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_BoolVariable{
					BoolVariable: &policylangv1.BoolVariable{
						ConstantOutput: podScaler.GetDryRun(),
						ConfigKey:      podScaler.GetDryRunConfigKey(),
						OutPorts: &policylangv1.BoolVariable_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "DRY_RUN",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Switcher{
					Switcher: &policylangv1.Switcher{
						InPorts: &policylangv1.Switcher_Ins{
							Switch: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DRY_RUN",
								},
							},
							OnSignal: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_SpecialValue{
											SpecialValue: "NaN",
										},
									},
								},
							},
							OffSignal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "REPLICAS_INPUT",
								},
							},
						},
						OutPorts: &policylangv1.Switcher_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "REPLICAS",
							},
						},
					},
				},
			},

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

	components.AddNestedIngress(nestedCircuit, inputReplicasPortName, "REPLICAS_INPUT")
	components.AddNestedEgress(nestedCircuit, outputActualReplicasPortName, "ACTUAL_REPLICAS")
	components.AddNestedEgress(nestedCircuit, outputConfiguredReplicasPortName, "CONFIGURED_REPLICAS")

	kos := podScaler.KubernetesObjectSelector
	sd := fmt.Sprintf("%s/%s/%s/%s/%s",
		kos.GetAgentGroup(),
		kos.GetNamespace(),
		kos.GetApiVersion(),
		kos.GetKind(),
		kos.GetName(),
	)
	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("PodScaler",
			sd,
			runtime.ComponentTypeSignalProcessor),
		podScaler,
		componentID,
		false,
	)
	if err != nil {
		return retErr(err)
	}
	return configuredComponent, nestedCircuit, nil
}
