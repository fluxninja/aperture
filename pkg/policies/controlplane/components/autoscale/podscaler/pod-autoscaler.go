package podscaler

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

const (
	autoscalerActualReplicasPortName     = "actual_replicas"
	autoscalerConfiguredReplicasPortName = "configured_replicas"
	autoscalerDesiredReplicasPortName    = "desired_replicas"
)

// ParsePodAutoscaler parses a PodAutoscaler and returns its nested circuit representation.
func ParsePodAutoscaler(
	podAutoscaler *policylangv1.PodAutoscaler,
) (*policylangv1.NestedCircuit, error) {
	// nestedInPortsMap := make(map[string]*policylangv1.InPort)

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := podAutoscaler.OutPorts
	if outPorts != nil {
		actualReplicasPort := outPorts.ActualReplicas
		if actualReplicasPort != nil {
			nestedOutPortsMap[autoscalerActualReplicasPortName] = actualReplicasPort
		}
		configuredReplicasPort := outPorts.ConfiguredReplicas
		if configuredReplicasPort != nil {
			nestedOutPortsMap[autoscalerConfiguredReplicasPortName] = configuredReplicasPort
		}
		desiredReplicasPort := outPorts.DesiredReplicas
		if desiredReplicasPort != nil {
			nestedOutPortsMap[autoscalerDesiredReplicasPortName] = desiredReplicasPort
		}
	}

	// Components to find the min and max values for the desired replicas.
	/*minMaxComponents := []*policylangv1.Component{
			{
				Component: &policylangv1.Component_ArithmeticCombinator{
	        ArithmeticCombinator: &policylangv1.ArithmeticCombinator{,
	          Operator: "mul",
	          InPorts: &policylangv1.ArithmeticCombinator_Ins{
	            Lhs: &policylangv1.InPort{
	              Value: &policylangv1.InPort_ConstantSignal{
	                Const: &policylangv1.ConstantSignal_Value{
	                  Value: float64(podAutoscaler.ScaleInMaxPercentage) / 100.0,
	                },
	              },
	            },
	            Rhs: &policylangv1.InPort{
	              Value: &policylangv1.InPort_SignalName{
	                SignalName: "ACTUAL_REPLICAS",
	              },
	            },
	          },
	          OutPorts: &policylangv1.ArithmeticCombinator_Outs{
	            Output: &policylangv1.OutPort{
	              SignalName: "MAX_SCALE_IN_PRE_FLOOR"
	            },
	          },
	        },
				},
			},
		}*/

	return nil, nil
}
