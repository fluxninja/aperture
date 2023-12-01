package autoscale

import (
	"fmt"
	"math"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

const (
	actualReplicasPortName     = "actual_scale"
	configuredReplicasPortName = "configured_scale"
	desiredReplicasPortName    = "desired_scale"
)

// ParseAutoScaler parses a AutoScaler and returns its nested circuit representation.
func ParseAutoScaler(
	autoscaler *policylangv1.AutoScaler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	retErr := func(err error) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
		return nil, nil, err
	}

	nestedInPortsMap := make(map[string]*policylangv1.InPort)

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)

	var (
		componentsScaler   []*policylangv1.Component
		shortDescription   string
		minScale, maxScale float64
	)
	scalingBackend := autoscaler.ScalingBackend
	if scalingBackend == nil {
		return retErr(fmt.Errorf("no scaling backend specified"))
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
		}
	} else {
		return retErr(fmt.Errorf("unsupported scaler type: %T", scalingBackend))
	}

	componentsAlerters := []*policylangv1.Component{
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
								SignalName: "CONFIGURED_SCALE",
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
								SignalName: "CONFIGURED_SCALE",
							},
						},
					},
					OutPorts: &policylangv1.Decider_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "DESIRED_SCALE_LESSER",
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
								SignalName: "DESIRED_SCALE_LESSER",
							},
						},
					},
				},
			},
		},
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
		signalPortName := fmt.Sprintf("scale_out_signal_%d", scaleOutIndex)
		setpointPortName := fmt.Sprintf("scale_out_setpoint_%d", scaleOutIndex)

		const prefix string = "SCALE_OUT"

		scaleOuts = append(scaleOuts, &policylangv1.InPort{
			Value: &policylangv1.InPort_SignalName{
				SignalName: fmt.Sprintf("%s_%d", prefix, scaleOutIndex),
			},
		})

		controller := scaleOutController.GetController()
		if controller == nil {
			return retErr(fmt.Errorf("scale out controller is nil"))
		}

		alerter := scaleOutController.GetAlerter()
		if alerter == nil {
			return retErr(fmt.Errorf("alerter is nil"))
		}

		switch controllerType := controller.Controller.(type) {
		case *policylangv1.ScaleOutController_Controller_Gradient:
			gradient := controllerType.Gradient

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

			components := createGradientControllerComponents(prefix, scaleOutIndex, signalPortName, setpointPortName, parameters, alerter)

			componentsScaleOut = append(componentsScaleOut, components...)
		default:
			return retErr(fmt.Errorf("scale out controller is not defined or of unexpected type"))
		}
	}

	for scaleInIndex, scaleInController := range autoscaler.ScaleInControllers {
		signalPortName := fmt.Sprintf("scale_in_signal_%d", scaleInIndex)
		setpointPortName := fmt.Sprintf("scale_in_setpoint_%d", scaleInIndex)

		const prefix string = "SCALE_IN"

		scaleIns = append(scaleIns, &policylangv1.InPort{
			Value: &policylangv1.InPort_SignalName{
				SignalName: fmt.Sprintf("%s_%d", prefix, scaleInIndex),
			},
		})

		controller := scaleInController.GetController()
		if controller == nil {
			return retErr(fmt.Errorf("scale in controller is nil"))
		}

		alerter := scaleInController.GetAlerter()
		if alerter == nil {
			return retErr(fmt.Errorf("alerter is nil"))
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

			components := createGradientControllerComponents(prefix, scaleInIndex, signalPortName, setpointPortName, parameters, alerter)

			componentsScaleIn = append(componentsScaleIn, components...)
		case *policylangv1.ScaleInController_Controller_Periodic:
			parameters := controllerType.Periodic

			components := createPeriodicControllerComponents(prefix, scaleInIndex, policyReadAPI.GetEvaluationInterval(), parameters, alerter)

			componentsScaleIn = append(componentsScaleIn, components...)
		default:
			return retErr(fmt.Errorf("scale in controller is not defined or of unexpected type"))
		}
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
							SignalName: "SCALE_OUT_MAX",
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
								SignalName: "SCALE_OUT_MAX",
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
							SignalName: "IS_SCALE_OUT_INTENT",
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
								SignalName: "IS_SCALE_OUT_INTENT",
							},
						},
						OnSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT_MAX",
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_OUT_INTENT",
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
							SignalName: "SCALE_IN_MAX",
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
								SignalName: "SCALE_IN_MAX",
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
							SignalName: "IS_SCALE_IN_INTENT",
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
								SignalName: "IS_SCALE_IN_INTENT",
							},
						},
						OnSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_IN_MAX",
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN_INTENT",
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
								SignalName: "SCALE_OUT_INTENT_HOLD",
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
								SignalName: "SCALE_OUT_INTENT",
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
					HoldFor: autoscaler.ScalingParameters.ScaleOutCooldown,
					InPorts: &policylangv1.Holder_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT_INTENT",
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
							SignalName: "SCALE_OUT_INTENT_HOLD",
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
								SignalName: "SCALE_OUT_INTENT_HOLD",
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
								SignalName: "SCALE_IN_INTENT",
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN_ONLY_INTENT",
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
								SignalName: "SCALE_IN_ONLY_INTENT",
							},
						},
						Reset_: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT_INTENT_HOLD",
							},
						},
					},
					OutPorts: &policylangv1.Holder_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN_ONLY_INTENT_HOLD",
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
								SignalName: "SCALE_OUT_INTENT_HOLD",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "CONFIGURED_SCALE",
							},
						},
					},
					OutPorts: &policylangv1.Decider_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "IS_SCALE_OUT",
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
								SignalName: "IS_SCALE_OUT",
							},
						},
						OnSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_OUT_INTENT_HOLD",
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_OUT",
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
								SignalName: "SCALE_IN_ONLY_INTENT_HOLD",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "CONFIGURED_SCALE",
							},
						},
					},
					OutPorts: &policylangv1.Decider_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "IS_SCALE_IN_ONLY",
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
								SignalName: "IS_SCALE_IN_ONLY",
							},
						},
						OnSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SCALE_IN_ONLY_INTENT_HOLD",
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "SCALE_IN_ONLY",
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
									SignalName: "SCALE_OUT",
								},
							},
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SCALE_IN_ONLY",
								},
							},
						},
					},
					OutPorts: &policylangv1.FirstValid_Outs{
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
	components = append(components, componentsAlerters...)

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("AutoScaler",
			shortDescription,
			runtime.ComponentTypeSignalProcessor),
		autoscaler,
		componentID,
		false,
	)
	if err != nil {
		return retErr(err)
	}

	// Construct nested circuit.
	nestedCircuit := &policylangv1.NestedCircuit{
		InPortsMap:  nestedInPortsMap,
		OutPortsMap: nestedOutPortsMap,
		Components:  components,
	}

	return configuredComponent, nestedCircuit, nil
}

func createGradientControllerComponents(
	prefix string,
	index int,
	signalPortName string,
	setpointPortName string,
	parameters *policylangv1.GradientController_Parameters,
	alerter *policylangv1.Alerter_Parameters,
) []*policylangv1.Component {
	scaleXSignal := fmt.Sprintf("%s_SIGNAL_%d", prefix, index)
	scaleXSetpoint := fmt.Sprintf("%s_SETPOINT_%d", prefix, index)
	scaleXPreCeil := fmt.Sprintf("%s_PRE_CEIL_%d", prefix, index)
	scaleXPreCeilFirstValid := fmt.Sprintf("%s_PRE_CEIL_FIRST_VALID_%d", prefix, index)
	scaleX := fmt.Sprintf("%s_%d", prefix, index)

	components := []*policylangv1.Component{}

	gradientComponents := []*policylangv1.Component{
		{
			Component: &policylangv1.Component_NestedSignalIngress{
				NestedSignalIngress: &policylangv1.NestedSignalIngress{
					PortName: signalPortName,
					OutPorts: &policylangv1.NestedSignalIngress_Outs{
						Signal: &policylangv1.OutPort{
							SignalName: scaleXSignal,
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
							SignalName: scaleXSetpoint,
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
								SignalName: scaleXSignal,
							},
						},
						Setpoint: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: scaleXSetpoint,
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
							SignalName: scaleXPreCeil,
						},
					},
				},
			},
		},
	}
	components = append(components, gradientComponents...)

	commonComponents := createCommonComponents(scaleXPreCeil, scaleXPreCeilFirstValid, scaleX, alerter)
	components = append(components, commonComponents...)

	return components
}

func createPeriodicControllerComponents(
	prefix string,
	index int,
	policyEvaluationInterval time.Duration,
	parameters *policylangv1.PeriodicDecrease,
	alerter *policylangv1.Alerter_Parameters,
) []*policylangv1.Component {
	scaleXReductionPreCeil := fmt.Sprintf("%s_REDUCTION_PRE_CEIL_%d", prefix, index)
	scaleXReductionPreCeilAdjusted := fmt.Sprintf("%s_REDUCTION_PRE_CEIL_ADJUSTED_%d", prefix, index)
	scaleXReduction := fmt.Sprintf("%s_REDUCTION_%d", prefix, index)
	scaleXPropsedScale := fmt.Sprintf("%s_PROPOSED_SCALE_%d", prefix, index)
	scaleXPeriodicPulse := fmt.Sprintf("%s_PERIODIC_PULSE_%d", prefix, index)
	scaleXPreCeil := fmt.Sprintf("%s_PRE_CEIL_%d", prefix, index)
	scaleXPreCeilFirstValid := fmt.Sprintf("%s_PRE_CEIL_FIRST_VALID_%d", prefix, index)
	scaleX := fmt.Sprintf("%s_%d", prefix, index)

	components := []*policylangv1.Component{}

	periodicComponents := []*policylangv1.Component{
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
							SignalName: scaleXReductionPreCeil,
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
									SignalName: scaleXReductionPreCeil,
								},
							},
						},
					},
					OutPorts: &policylangv1.Max_Outs{
						Output: &policylangv1.OutPort{
							SignalName: scaleXReductionPreCeilAdjusted,
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
								SignalName: scaleXReductionPreCeilAdjusted,
							},
						},
					},
					OutPorts: &policylangv1.UnaryOperator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: scaleXReduction,
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
								SignalName: scaleXReduction,
							},
						},
					},
					OutPorts: &policylangv1.ArithmeticCombinator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: scaleXPropsedScale,
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_PulseGenerator{
				PulseGenerator: &policylangv1.PulseGenerator{
					FalseFor: parameters.Period,
					TrueFor:  durationpb.New(policyEvaluationInterval),
					OutPorts: &policylangv1.PulseGenerator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: scaleXPeriodicPulse,
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
								SignalName: scaleXPeriodicPulse,
							},
						},
						OnSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: scaleXPropsedScale,
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: scaleXPreCeil,
						},
					},
				},
			},
		},
	}
	components = append(components, periodicComponents...)

	commonComponents := createCommonComponents(scaleXPreCeil, scaleXPreCeilFirstValid, scaleX, alerter)
	components = append(components, commonComponents...)

	return components
}

func createCommonComponents(
	scalePreCeil string,
	scalePreCeilFirstValid string,
	scale string,
	alerter *policylangv1.Alerter_Parameters,
) []*policylangv1.Component {
	components := []*policylangv1.Component{
		{
			Component: &policylangv1.Component_FirstValid{
				FirstValid: &policylangv1.FirstValid{
					InPorts: &policylangv1.FirstValid_Ins{
						Inputs: []*policylangv1.InPort{
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: scalePreCeil,
								},
							},
							{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACTUAL_SCALE",
								},
							},
						},
					},
					OutPorts: &policylangv1.FirstValid_Outs{
						Output: &policylangv1.OutPort{
							SignalName: scalePreCeilFirstValid,
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
								SignalName: scalePreCeilFirstValid,
							},
						},
					},
					OutPorts: &policylangv1.UnaryOperator_Outs{
						Output: &policylangv1.OutPort{
							SignalName: scale,
						},
					},
				},
			},
		},
		{
			Component: &policylangv1.Component_Alerter{
				Alerter: &policylangv1.Alerter{
					Parameters: alerter,
					InPorts: &policylangv1.Alerter_Ins{
						Signal: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: scale,
							},
						},
					},
				},
			},
		},
	}
	return components
}
