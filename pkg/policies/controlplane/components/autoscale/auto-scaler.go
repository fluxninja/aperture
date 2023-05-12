package autoscale

import (
	"fmt"
	"math"

	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
)

const (
	actualReplicasPortName                     = "actual_scale"
	configuredReplicasPortName                 = "configured_scale"
	desiredReplicasPortName                    = "desired_scale"
	autoscalerScaleOutSignalPortNameTemplate   = "scale_out_signal_%d"
	autoscalerScaleOutSetpointPortNameTemplate = "scale_out_setpoint_%d"
	autoscalerScaleInSignalPortNameTemplate    = "scale_in_signal_%d"
	autoscalerScaleInSetpointPortNameTemplate  = "scale_in_setpoint_%d"
)

// ParseAutoScaler parses a AutoScaler and returns its nested circuit representation.
func ParseAutoScaler(
	autoscaler *policylangv1.AutoScaler,
	policyReadAPI iface.Policy,
) (*policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)

	var (
		componentsScaler   []*policylangv1.Component
		shortDescription   string
		minScale, maxScale float64
	)
	scalingBackend := autoscaler.ScalingBackend
	if scalingBackend == nil {
		return nil, fmt.Errorf("no scaling backend specified")
	}
	if kubernetesReplicas := scalingBackend.GetKubernetesReplicas(); kubernetesReplicas != nil {
		minScale = float64(kubernetesReplicas.MinReplicas)
		maxScale = float64(kubernetesReplicas.MaxReplicas)
		shortDescription = kubernetesReplicas.KubernetesObjectSelector.GetNamespace() + ":" + kubernetesReplicas.KubernetesObjectSelector.GetKind() + "/" + kubernetesReplicas.KubernetesObjectSelector.Name
		outPorts := kubernetesReplicas.OutPorts
		if outPorts != nil {
			actualReplicasPort := outPorts.ActualReplicas
			if actualReplicasPort != nil {
				nestedOutPortsMap[actualReplicasPortName] = actualReplicasPort
			}
			configuredReplicasPort := outPorts.ConfiguredReplicas
			if configuredReplicasPort != nil {
				nestedOutPortsMap[configuredReplicasPortName] = configuredReplicasPort
			}
			desiredReplicasPort := outPorts.DesiredReplicas
			if desiredReplicasPort != nil {
				nestedOutPortsMap[desiredReplicasPortName] = desiredReplicasPort
			}
		}

		componentsScaler = []*policylangv1.Component{
			{
				Component: &policylangv1.Component_AutoScale{
					AutoScale: &policylangv1.AutoScale{
						Component: &policylangv1.AutoScale_PodScaler{
							PodScaler: &policylangv1.PodScaler{
								KubernetesObjectSelector: kubernetesReplicas.KubernetesObjectSelector,
								InPorts: &policylangv1.PodScaler_Ins{
									Replicas: &policylangv1.InPort{
										Value: &policylangv1.InPort_SignalName{
											SignalName: "DESIRED_SCALE",
										},
									},
								},
								OutPorts: &policylangv1.PodScaler_Outs{
									ActualReplicas: &policylangv1.OutPort{
										SignalName: "ACTUAL_SCALE",
									},
									ConfiguredReplicas: &policylangv1.OutPort{
										SignalName: "CONFIGURED_SCALE",
									},
								},
								DryRunConfigKey: autoscaler.DryRunConfigKey,
								DryRun:          autoscaler.DryRun,
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_NestedSignalEgress{
					NestedSignalEgress: &policylangv1.NestedSignalEgress{
						PortName: actualReplicasPortName,
						InPorts: &policylangv1.NestedSignalEgress_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACTUAL_SCALE",
								},
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_NestedSignalEgress{
					NestedSignalEgress: &policylangv1.NestedSignalEgress{
						PortName: configuredReplicasPortName,
						InPorts: &policylangv1.NestedSignalEgress_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "CONFIGURED_SCALE",
								},
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_NestedSignalEgress{
					NestedSignalEgress: &policylangv1.NestedSignalEgress{
						PortName: desiredReplicasPortName,
						InPorts: &policylangv1.NestedSignalEgress_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_SCALE",
								},
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Decider{
					Decider: &policylangv1.Decider{
						Operator: "gt",
						InPorts: &policylangv1.Decider_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_SCALE",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACTUAL_SCALE",
								},
							},
						},
						OutPorts: &policylangv1.Decider_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "DESIRED_SCALE_GREATER",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Alerter{
					Alerter: &policylangv1.Alerter{
						Parameters: autoscaler.ScalingParameters.ScaleOutAlerter,
						InPorts: &policylangv1.Alerter_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_SCALE_GREATER",
								},
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Decider{
					Decider: &policylangv1.Decider{
						Operator: "lt",
						InPorts: &policylangv1.Decider_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_SCALE",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACTUAL_SCALE",
								},
							},
						},
						OutPorts: &policylangv1.Decider_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "DESIRED_SCALE_LESS",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Alerter{
					Alerter: &policylangv1.Alerter{
						Parameters: autoscaler.ScalingParameters.ScaleInAlerter,
						InPorts: &policylangv1.Alerter_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_SCALE_LESS",
								},
							},
						},
					},
				},
			},
		}
	} else {
		return nil, fmt.Errorf("unsupported scaler type: %T", scalingBackend)
	}

	// Components to find the min and max values for the desired scale.
	componentsMinDesiredScale := []*policylangv1.Component{
		{
			Component: &policylangv1.Component_ArithmeticCombinator{
				ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
					Operator: "mul",
					InPorts: &policylangv1.ArithmeticCombinator_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: float64(autoscaler.ScalingParameters.MaxScaleInPercentage) / 100.0,
									},
								},
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "ACTUAL_SCALE",
							},
						},
					},
					OutPorts: &policylangv1.ArithmeticCombinator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_SCALE_IN_PRE_FLOOR",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_UnaryOperator{
				UnaryOperator: &policylangv1.UnaryOperator{
					Operator: "floor",
					InPorts: &policylangv1.UnaryOperator_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "MAX_SCALE_IN_PRE_FLOOR",
							},
						},
					},
					OutPorts: &policylangv1.UnaryOperator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_SCALE_IN_PRE_CONSTRAINT",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MAX_SCALE_IN_PRE_CONSTRAINT",
								},
							},
							{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: 1.0,
										},
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_SCALE_IN",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_ArithmeticCombinator{
				ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
					Operator: "sub",
					InPorts: &policylangv1.ArithmeticCombinator_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "ACTUAL_SCALE",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "MAX_SCALE_IN",
							},
						},
					},
					OutPorts: &policylangv1.ArithmeticCombinator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MIN_DESIRED_SCALE_PRE_CONSTRAINT",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MIN_DESIRED_SCALE_PRE_CONSTRAINT",
								},
							},
							{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: float64(minScale),
										},
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MIN_DESIRED_SCALE_POST_MIN_SCALE",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Min{
				Min: &policylangv1.Min{
					InPorts: &policylangv1.Min_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MIN_DESIRED_SCALE_POST_MIN_SCALE",
								},
							},
							{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: float64(maxScale),
										},
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Min_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MIN_DESIRED_SCALE",
						},
					},
				},
			},
		},
	}

	componentsMaxDesiredScale := []*policylangv1.Component{
		{
			Component: &policylangv1.Component_ArithmeticCombinator{
				ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
					Operator: "mul",
					InPorts: &policylangv1.ArithmeticCombinator_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: float64(autoscaler.ScalingParameters.MaxScaleOutPercentage) / 100.0,
									},
								},
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "ACTUAL_SCALE",
							},
						},
					},
					OutPorts: &policylangv1.ArithmeticCombinator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_SCALE_OUT_PRE_FLOOR",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_UnaryOperator{
				UnaryOperator: &policylangv1.UnaryOperator{
					Operator: "floor",
					InPorts: &policylangv1.UnaryOperator_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "MAX_SCALE_OUT_PRE_FLOOR",
							},
						},
					},
					OutPorts: &policylangv1.UnaryOperator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_SCALE_OUT_PRE_CONSTRAINT",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MAX_SCALE_OUT_PRE_CONSTRAINT",
								},
							},
							{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: 1.0,
										},
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_SCALE_OUT",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_ArithmeticCombinator{
				ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
					Operator: "add",
					InPorts: &policylangv1.ArithmeticCombinator_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "ACTUAL_SCALE",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "MAX_SCALE_OUT",
							},
						},
					},
					OutPorts: &policylangv1.ArithmeticCombinator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_DESIRED_SCALE_PRE_CONSTRAINT",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Min{
				Min: &policylangv1.Min{
					InPorts: &policylangv1.Min_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MAX_DESIRED_SCALE_PRE_CONSTRAINT",
								},
							},
							{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: float64(maxScale),
										},
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Min_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_DESIRED_SCALE_POST_MAX_SCALE",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MAX_DESIRED_SCALE_POST_MAX_SCALE",
								},
							},
							{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: float64(minScale),
										},
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "MAX_DESIRED_SCALE",
						},
					},
				},
			},
		},
	}

	var (
		scaleOuts, scaleIns                   []*policylangv1.InPort
		componentsScaleOut, componentsScaleIn []*policylangv1.Component
	)

	for scaleOutIndex, scaleOutController := range autoscaler.ScaleOutControllers {
		signalPortName := fmt.Sprintf(autoscalerScaleOutSignalPortNameTemplate, scaleOutIndex)
		setpointPortName := fmt.Sprintf(autoscalerScaleOutSetpointPortNameTemplate, scaleOutIndex)

		scaleOutSignal := fmt.Sprintf("SCALE_OUT_SIGNAL_%d", scaleOutIndex)
		scaleOutSetpoint := fmt.Sprintf("SCALE_OUT_SETPOINT_%d", scaleOutIndex)
		scaleOutOutputPreCeil := fmt.Sprintf("SCALE_OUT_OUTPUT_PRE_CEIL_%d", scaleOutIndex)
		scaleOutOutput := fmt.Sprintf("SCALE_OUT_OUTPUT_%d", scaleOutIndex)
		isScaleOut := fmt.Sprintf("IS_SCALE_OUT_%d", scaleOutIndex)
		scaleOut := fmt.Sprintf("SCALE_OUT_%d", scaleOutIndex)

		scaleOuts = append(scaleOuts, &policylangv1.InPort{
			Value: &policylangv1.InPort_SignalName{
				SignalName: scaleOut,
			},
		})

		controller := scaleOutController.GetController()
		if controller == nil {
			return nil, fmt.Errorf("scale out controller is nil")
		}
		if gradient := controller.GetGradient(); gradient != nil {
			inPorts := gradient.GetInPorts()
			if inPorts == nil {
				inPorts = &policylangv1.IncreasingGradient_Ins{}
			}
			signalPort := inPorts.GetSignal()
			if signalPort != nil {
				nestedInPortsMap[signalPortName] = signalPort
			}
			setpointPort := inPorts.GetSetpoint()
			if setpointPort != nil {
				nestedInPortsMap[setpointPortName] = setpointPort
			}
			parameters := &policylangv1.GradientController_Parameters{
				Slope:       1.0,
				MinGradient: 1.0,
				MaxGradient: math.Inf(1),
			}
			if params := gradient.GetParameters(); params != nil {
				parameters.Slope = params.Slope
				parameters.MaxGradient = params.MaxGradient
			}

			components := []*policylangv1.Component{
				{
					Component: &policylangv1.Component_Alerter{
						Alerter: &policylangv1.Alerter{
							Parameters: scaleOutController.Alerter,
							InPorts: &policylangv1.Alerter_Ins{
								Signal: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleOutSetpoint,
									},
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_NestedSignalIngress{
						NestedSignalIngress: &policylangv1.NestedSignalIngress{
							PortName: signalPortName,
							OutPorts: &policylangv1.NestedSignalIngress_Outs{
								Signal: &policylangv1.OutPort{
									SignalName: scaleOutSignal,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_NestedSignalIngress{
						NestedSignalIngress: &policylangv1.NestedSignalIngress{
							PortName: setpointPortName,
							OutPorts: &policylangv1.NestedSignalIngress_Outs{
								Signal: &policylangv1.OutPort{
									SignalName: scaleOutSetpoint,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_GradientController{
						GradientController: &policylangv1.GradientController{
							Parameters: parameters,
							InPorts: &policylangv1.GradientController_Ins{
								Signal: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleOutSignal,
									},
								},
								Setpoint: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleOutSetpoint,
									},
								},
								ControlVariable: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "ACTUAL_SCALE",
									},
								},
							},
							OutPorts: &policylangv1.GradientController_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleOutOutputPreCeil,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_UnaryOperator{
						UnaryOperator: &policylangv1.UnaryOperator{
							Operator: "ceil",
							InPorts: &policylangv1.UnaryOperator_Ins{
								Input: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleOutOutputPreCeil,
									},
								},
							},
							OutPorts: &policylangv1.UnaryOperator_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleOutOutput,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_Decider{
						Decider: &policylangv1.Decider{
							Operator: "gt",
							TrueFor:  durationpb.New(0),
							FalseFor: durationpb.New(0),
							InPorts: &policylangv1.Decider_Ins{
								Lhs: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleOutOutput,
									},
								},
								Rhs: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "ACTUAL_SCALE",
									},
								},
							},
							OutPorts: &policylangv1.Decider_Outs{
								Output: &policylangv1.OutPort{
									SignalName: isScaleOut,
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
										SignalName: isScaleOut,
									},
								},
								OnSignal: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleOutOutput,
									},
								},
							},
							OutPorts: &policylangv1.Switcher_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleOut,
								},
							},
						},
					},
				},
			}

			componentsScaleOut = append(componentsScaleOut, components...)
		} else {
			return nil, fmt.Errorf("scale out controller is not defined or of unexpected type")
		}
	}

	for scaleInIndex, scaleInController := range autoscaler.ScaleInControllers {
		signalPortName := fmt.Sprintf(autoscalerScaleInSignalPortNameTemplate, scaleInIndex)
		setpointPortName := fmt.Sprintf(autoscalerScaleInSetpointPortNameTemplate, scaleInIndex)

		scaleInReductionPreCeil := fmt.Sprintf("SCALE_IN_REDUCTION_PRE_CEIL_%d", scaleInIndex)
		scaleInReductionPreCeilAdjusted := fmt.Sprintf("SCALE_IN_REDUCTION_PRE_CEIL_ADJUSTED_%d", scaleInIndex)
		scaleInReduction := fmt.Sprintf("SCALE_IN_REDUCTION_%d", scaleInIndex)
		scaleInPropsedScale := fmt.Sprintf("SCALE_IN_PROPOSED_SCALE_%d", scaleInIndex)
		scaleInPeriodicPulse := fmt.Sprintf("SCALE_IN_PERIODIC_PULSE_%d", scaleInIndex)
		scaleInSignal := fmt.Sprintf("SCALE_IN_SIGNAL_%d", scaleInIndex)
		scaleInSetpoint := fmt.Sprintf("SCALE_IN_SETPOINT_%d", scaleInIndex)
		scaleInOutputPreCeil := fmt.Sprintf("SCALE_IN_OUTPUT_PRE_CEIL_%d", scaleInIndex)
		scaleInOutput := fmt.Sprintf("SCALE_IN_OUTPUT_%d", scaleInIndex)
		isScaleIn := fmt.Sprintf("IS_SCALE_IN_%d", scaleInIndex)
		scaleIn := fmt.Sprintf("SCALE_IN_%d", scaleInIndex)

		scaleIns = append(scaleIns, &policylangv1.InPort{
			Value: &policylangv1.InPort_SignalName{
				SignalName: scaleIn,
			},
		})

		controller := scaleInController.GetController()
		if controller == nil {
			return nil, fmt.Errorf("scale in controller is nil")
		}

		switch controllerType := controller.Controller.(type) {
		case *policylangv1.ScaleInController_Controller_Gradient:
			gradient := controllerType.Gradient

			inPorts := gradient.GetInPorts()
			if inPorts == nil {
				inPorts = &policylangv1.DecreasingGradient_Ins{}
			}
			signalPort := inPorts.GetSignal()
			if signalPort != nil {
				nestedInPortsMap[signalPortName] = signalPort
			}
			setpointPort := inPorts.GetSetpoint()
			if setpointPort != nil {
				nestedInPortsMap[setpointPortName] = setpointPort
			}
			parameters := &policylangv1.GradientController_Parameters{
				Slope:       1.0,
				MinGradient: math.Inf(-1),
				MaxGradient: 1.0,
			}
			if params := gradient.GetParameters(); params != nil {
				parameters.Slope = params.GetSlope()
				parameters.MinGradient = params.GetMinGradient()
			}

			components := []*policylangv1.Component{
				{
					Component: &policylangv1.Component_Alerter{
						Alerter: &policylangv1.Alerter{
							Parameters: scaleInController.Alerter,
							InPorts: &policylangv1.Alerter_Ins{
								Signal: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleInSetpoint,
									},
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_NestedSignalIngress{
						NestedSignalIngress: &policylangv1.NestedSignalIngress{
							PortName: signalPortName,
							OutPorts: &policylangv1.NestedSignalIngress_Outs{
								Signal: &policylangv1.OutPort{
									SignalName: scaleInSignal,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_NestedSignalIngress{
						NestedSignalIngress: &policylangv1.NestedSignalIngress{
							PortName: setpointPortName,
							OutPorts: &policylangv1.NestedSignalIngress_Outs{
								Signal: &policylangv1.OutPort{
									SignalName: scaleInSetpoint,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_GradientController{
						GradientController: &policylangv1.GradientController{
							Parameters: parameters,
							InPorts: &policylangv1.GradientController_Ins{
								Signal: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleInSignal,
									},
								},
								Setpoint: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleInSetpoint,
									},
								},
								ControlVariable: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "ACTUAL_SCALE",
									},
								},
							},
							OutPorts: &policylangv1.GradientController_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleInOutputPreCeil,
								},
							},
						},
					},
				},
			}

			componentsScaleIn = append(componentsScaleIn, components...)
		case *policylangv1.ScaleInController_Controller_Periodic:
			periodic := controllerType.Periodic
			parameters := periodic.GetParameters()

			components := []*policylangv1.Component{
				{
					Component: &policylangv1.Component_ArithmeticCombinator{
						ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
							Operator: "mul",
							InPorts: &policylangv1.ArithmeticCombinator_Ins{
								Lhs: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "ACTUAL_SCALE",
									},
								},
								Rhs: &policylangv1.InPort{
									Value: &policylangv1.InPort_ConstantSignal{
										ConstantSignal: &policylangv1.ConstantSignal{
											Const: &policylangv1.ConstantSignal_Value{
												Value: (float64(parameters.ScaleInPercentage) / 100.0),
											},
										},
									},
								},
							},
							OutPorts: &policylangv1.ArithmeticCombinator_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleInReductionPreCeil,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_Max{
						Max: &policylangv1.Max{
							InPorts: &policylangv1.Max_Ins{
								Inputs: []*policylangv1.InPort{
									{
										Value: &policylangv1.InPort_ConstantSignal{
											ConstantSignal: &policylangv1.ConstantSignal{
												Const: &policylangv1.ConstantSignal_Value{
													Value: 1.0,
												},
											},
										},
									},
									{
										Value: &policylangv1.InPort_SignalName{
											SignalName: scaleInReductionPreCeil,
										},
									},
								},
							},
							OutPorts: &policylangv1.Max_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleInReductionPreCeilAdjusted,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_UnaryOperator{
						UnaryOperator: &policylangv1.UnaryOperator{
							Operator: "ceil",
							InPorts: &policylangv1.UnaryOperator_Ins{
								Input: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleInReductionPreCeilAdjusted,
									},
								},
							},
							OutPorts: &policylangv1.UnaryOperator_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleInReduction,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_ArithmeticCombinator{
						ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
							Operator: "sub",
							InPorts: &policylangv1.ArithmeticCombinator_Ins{
								Lhs: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "ACTUAL_SCALE",
									},
								},
								Rhs: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleInReduction,
									},
								},
							},
							OutPorts: &policylangv1.ArithmeticCombinator_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleInPropsedScale,
								},
							},
						},
					},
				},
				{
					Component: &policylangv1.Component_PulseGenerator{
						PulseGenerator: &policylangv1.PulseGenerator{
							FalseFor: parameters.Period,
							TrueFor:  durationpb.New(policyReadAPI.GetEvaluationInterval()),
							OutPorts: &policylangv1.PulseGenerator_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleInPeriodicPulse,
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
										SignalName: scaleInPeriodicPulse,
									},
								},
								OnSignal: &policylangv1.InPort{
									Value: &policylangv1.InPort_SignalName{
										SignalName: scaleInPropsedScale,
									},
								},
							},
							OutPorts: &policylangv1.Switcher_Outs{
								Output: &policylangv1.OutPort{
									SignalName: scaleInOutputPreCeil,
								},
							},
						},
					},
				},
			}

			componentsScaleIn = append(componentsScaleIn, components...)
		default:
			return nil, fmt.Errorf("scale in controller is not defined or of unexpected type")
		}

		// common components
		commonComponents := []*policylangv1.Component{
			{
				Component: &policylangv1.Component_UnaryOperator{
					UnaryOperator: &policylangv1.UnaryOperator{
						Operator: "ceil",
						InPorts: &policylangv1.UnaryOperator_Ins{
							Input: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: scaleInOutputPreCeil,
								},
							},
						},
						OutPorts: &policylangv1.UnaryOperator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: scaleInOutput,
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Decider{
					Decider: &policylangv1.Decider{
						Operator: "lt",
						TrueFor:  durationpb.New(0),
						FalseFor: durationpb.New(0),
						InPorts: &policylangv1.Decider_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: scaleInOutput,
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACTUAL_SCALE",
								},
							},
						},
						OutPorts: &policylangv1.Decider_Outs{
							Output: &policylangv1.OutPort{
								SignalName: isScaleIn,
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
									SignalName: isScaleIn,
								},
							},
							OnSignal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: scaleInOutput,
								},
							},
						},
						OutPorts: &policylangv1.Switcher_Outs{
							Output: &policylangv1.OutPort{
								SignalName: scaleIn,
							},
						},
					},
				},
			},
		}
		// append common components to scale in components
		componentsScaleIn = append(componentsScaleIn, commonComponents...)
	}

	// Process scale in and scale out signals to scale the pods.
	componentsPreScaler := []*policylangv1.Component{
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: scaleOuts,
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_OUT",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: scaleIns,
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Holder{
				Holder: &policylangv1.Holder{
					HoldFor: autoscaler.ScalingParameters.ScaleOutCooldown,
					InPorts: &policylangv1.Holder_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT",
							},
						},
						Reset_: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "IS_COOLDOWN_OVERRIDE",
							},
						},
					},
					OutPorts: &policylangv1.Holder_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_OUT_HOLD",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Holder{
				Holder: &policylangv1.Holder{
					HoldFor: autoscaler.ScalingParameters.ScaleInCooldown,
					InPorts: &policylangv1.Holder_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT",
							},
						},
						Reset_: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "IS_COOLDOWN_OVERRIDE",
							},
						},
					},
					OutPorts: &policylangv1.Holder_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN_COOLDOWN_OVERRIDE",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_ArithmeticCombinator{
				ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
					Operator: "mul",
					InPorts: &policylangv1.ArithmeticCombinator_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT_HOLD",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: (float64(autoscaler.ScalingParameters.CooldownOverridePercentage) / 100.0) + 1.0,
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.ArithmeticCombinator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "COOLDOWN_OVERRIDE_THRESHOLD",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Decider{
				Decider: &policylangv1.Decider{
					Operator: "gt",
					InPorts: &policylangv1.Decider_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "COOLDOWN_OVERRIDE_THRESHOLD",
							},
						},
					},
					OutPorts: &policylangv1.Decider_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "IS_COOLDOWN_OVERRIDE",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Holder{
				Holder: &policylangv1.Holder{
					HoldFor: autoscaler.ScalingParameters.ScaleInCooldown,
					InPorts: &policylangv1.Holder_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_IN",
							},
						},
					},
					OutPorts: &policylangv1.Holder_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN_HOLD",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_FirstValid{
				FirstValid: &policylangv1.FirstValid{
					InPorts: &policylangv1.FirstValid_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SCALE_IN_COOLDOWN_OVERRIDE",
								},
							},
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SCALE_IN_HOLD",
								},
							},
						},
					},
					OutPorts: &policylangv1.FirstValid_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN_OR_OUT_HOLD",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_FirstValid{
				FirstValid: &policylangv1.FirstValid{
					InPorts: &policylangv1.FirstValid_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SCALE_OUT_HOLD",
								},
							},
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SCALE_IN_HOLD",
								},
							},
						},
					},
					OutPorts: &policylangv1.FirstValid_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_OUT_OR_IN_HOLD",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SCALE_OUT_OR_IN_HOLD",
								},
							},
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SCALE_IN_OR_OUT_HOLD",
								},
							},
						},
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "DESIRED_SCALE_PRE_MAX",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Max{
				Max: &policylangv1.Max{
					InPorts: &policylangv1.Max_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_SCALE_PRE_MAX",
								},
							},
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MIN_DESIRED_SCALE",
								},
							},
						},
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "DESIRED_SCALE_PRE_MIN",
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Min{
				Min: &policylangv1.Min{
					InPorts: &policylangv1.Min_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_SCALE_PRE_MIN",
								},
							},
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "MAX_DESIRED_SCALE",
								},
							},
						},
					},
					OutPorts: &policylangv1.Min_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "DESIRED_SCALE",
						},
					},
				},
			},
		},
	}

	// Concatenate the components into a single components list.
	var components []*policylangv1.Component
	components = append(components, componentsMinDesiredScale...)
	components = append(components, componentsMaxDesiredScale...)
	components = append(components, componentsScaleIn...)
	components = append(components, componentsScaleOut...)
	components = append(components, componentsPreScaler...)
	components = append(components, componentsScaler...)

	// Construct nested circuit.
	nestedCircuit := &policylangv1.NestedCircuit{
		Name:             "AutoScaler",
		ShortDescription: shortDescription,
		InPortsMap:       nestedInPortsMap,
		OutPortsMap:      nestedOutPortsMap,
		Components:       components,
	}

	return nestedCircuit, nil
}
