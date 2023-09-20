package loadscheduler

import (
	"time"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	alsSignalPortName                 = "signal"
	alsSetpointPortName               = "setpoint"
	alsOverloadConfirmationPortName   = "overload_confirmation"
	alsIsOverloadPortName             = "is_overload"
	alsDesiredLoadMultiplierPortName  = "desired_load_multiplier"
	alsObservedLoadMultiplierPortName = "observed_load_multiplier"
)

// ParseAdaptiveLoadScheduler parses and returns nested circuit representation of AdaptiveLoadScheduler.
func ParseAdaptiveLoadScheduler(
	adaptiveLoadScheduler *policylangv1.AdaptiveLoadScheduler,
	componentID runtime.ComponentID,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	retErr := func(err error) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
		return nil, nil, err
	}

	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := adaptiveLoadScheduler.InPorts
	if inPorts != nil {
		signalPort := inPorts.Signal
		if signalPort != nil {
			nestedInPortsMap[alsSignalPortName] = signalPort
		}
		setpointPort := inPorts.Setpoint
		if setpointPort != nil {
			nestedInPortsMap[alsSetpointPortName] = setpointPort
		}
		overloadConfirmation := inPorts.OverloadConfirmation
		if overloadConfirmation != nil {
			nestedInPortsMap[alsOverloadConfirmationPortName] = overloadConfirmation
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := adaptiveLoadScheduler.OutPorts
	if outPorts != nil {
		isOverloadPort := outPorts.IsOverload
		if isOverloadPort != nil {
			nestedOutPortsMap[alsIsOverloadPortName] = isOverloadPort
		}
		desiredLoadMultiplierPort := outPorts.DesiredLoadMultiplier
		if desiredLoadMultiplierPort != nil {
			nestedOutPortsMap[alsDesiredLoadMultiplierPortName] = desiredLoadMultiplierPort
		}
		observedLoadMultiplierPort := outPorts.ObservedLoadMultiplier
		if observedLoadMultiplierPort != nil {
			nestedOutPortsMap[alsObservedLoadMultiplierPortName] = observedLoadMultiplierPort
		}
	}

	overloadDeciderOperator := components.GT.String()
	// if slope is greater than 0 then we want to use less than operator
	if adaptiveLoadScheduler.Parameters.Gradient.Slope > 0 {
		overloadDeciderOperator = components.LT.String()
	}

	// Make a copy, because we shouldn't mutate the input proto.
	alerterParameters := proto.Clone(adaptiveLoadScheduler.Parameters.Alerter).(*policylangv1.Alerter_Parameters)
	if alerterParameters.Labels == nil {
		alerterParameters.Labels = make(map[string]string)
	}
	alerterParameters.Labels["type"] = "load_scheduler"

	nestedCircuit := &policylangv1.NestedCircuit{
		InPortsMap:  nestedInPortsMap,
		OutPortsMap: nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_FirstValid{
					FirstValid: &policylangv1.FirstValid{
						InPorts: &policylangv1.FirstValid_Ins{
							Inputs: []*policylangv1.InPort{
								{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "OVERLOAD_CONFIRMATION",
									},
								},
								{
									Value: &policylangv1.InPort_ConstantSignal{
										ConstantSignal: &policylangv1.ConstantSignal{
											Const: &policylangv1.ConstantSignal_Value{
												Value: 1, // Overload confirmation is assumed true by default. This makes the same circuit work in case overload confirmation is not provided. If the required behavior is to assume false by default then the policy needs to make sure to provide a valid signal with desired defaults.
											},
										},
									},
								},
							},
						},
						OutPorts: &policylangv1.FirstValid_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "OVERLOAD_CONFIRMATION_WITH_DEFAULT",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Decider{
					Decider: &policylangv1.Decider{
						Operator: overloadDeciderOperator,
						TrueFor:  durationpb.New(0),
						FalseFor: durationpb.New(0),
						InPorts: &policylangv1.Decider_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SIGNAL",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SETPOINT",
								},
							},
						},
						OutPorts: &policylangv1.Decider_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "OVERLOAD_BASED_ON_SIGNAL",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_And{
					And: &policylangv1.And{
						InPorts: &policylangv1.And_Ins{
							Inputs: []*policylangv1.InPort{
								{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "OVERLOAD_BASED_ON_SIGNAL",
									},
								},
								{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "OVERLOAD_CONFIRMATION_WITH_DEFAULT",
									},
								},
							},
						},
						OutPorts: &policylangv1.And_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "IS_OVERLOAD",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_GradientController{
					GradientController: &policylangv1.GradientController{
						Parameters: adaptiveLoadScheduler.Parameters.Gradient,
						InPorts: &policylangv1.GradientController_Ins{
							ControlVariable: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "OBSERVED_LOAD_MULTIPLIER",
								},
							},
							Setpoint: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SETPOINT",
								},
							},
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SIGNAL",
								},
							},
						},
						OutPorts: &policylangv1.GradientController_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "CONTROLLER_ADJUSTED_LOAD_MULTIPLIER",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Extrapolator{
					Extrapolator: &policylangv1.Extrapolator{
						Parameters: &policylangv1.Extrapolator_Parameters{
							MaxExtrapolationInterval: durationpb.New(time.Second * 5),
						},
						InPorts: &policylangv1.Extrapolator_Ins{
							Input: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "CONTROLLER_ADJUSTED_LOAD_MULTIPLIER",
								},
							},
						},
						OutPorts: &policylangv1.Extrapolator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "LOAD_MULTIPLIER_IF_OVERLOAD",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Integrator{
					Integrator: &policylangv1.Integrator{
						InitialValue:       adaptiveLoadScheduler.Parameters.MaxLoadMultiplier,
						EvaluationInterval: durationpb.New(metricScrapeInterval),
						InPorts: &policylangv1.Integrator_Ins{
							Input: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: adaptiveLoadScheduler.Parameters.LoadMultiplierLinearIncrement,
										},
									},
								},
							},
							Max: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: adaptiveLoadScheduler.Parameters.MaxLoadMultiplier,
										},
									},
								},
							},
							Reset_: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "IS_OVERLOAD",
								},
							},
						},
						OutPorts: &policylangv1.Integrator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "LOAD_MULTIPLIER_OPTIMIZATION",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_ArithmeticCombinator{
					ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
						Operator: components.Add.String(),
						InPorts: &policylangv1.ArithmeticCombinator_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "OBSERVED_LOAD_MULTIPLIER",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "LOAD_MULTIPLIER_OPTIMIZATION",
								},
							},
						},
						OutPorts: &policylangv1.ArithmeticCombinator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "LOAD_MULTIPLIER_IF_NORMAL",
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
									SignalName: "IS_OVERLOAD",
								},
							},
							OnSignal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "LOAD_MULTIPLIER_IF_OVERLOAD",
								},
							},
							OffSignal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "LOAD_MULTIPLIER_IF_NORMAL",
								},
							},
						},
						OutPorts: &policylangv1.Switcher_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "LOAD_MULTIPLIER_MUX",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Decider{
					Decider: &policylangv1.Decider{
						Operator: components.GTE.String(),
						InPorts: &policylangv1.Decider_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "LOAD_MULTIPLIER_MUX",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: adaptiveLoadScheduler.Parameters.MaxLoadMultiplier,
										},
									},
								},
							},
						},
						OutPorts: &policylangv1.Decider_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "PASS_THROUGH",
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
									SignalName: "PASS_THROUGH",
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
									SignalName: "LOAD_MULTIPLIER_MUX",
								},
							},
						},
						OutPorts: &policylangv1.Switcher_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "DESIRED_LOAD_MULTIPLIER",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_FlowControl{
					FlowControl: &policylangv1.FlowControl{
						Component: &policylangv1.FlowControl_LoadScheduler{
							LoadScheduler: &policylangv1.LoadScheduler{
								InPorts: &policylangv1.LoadScheduler_Ins{
									LoadMultiplier: &policylangv1.InPort{
										Value: &policylangv1.InPort_SignalName{
											SignalName: "DESIRED_LOAD_MULTIPLIER",
										},
									},
								},
								OutPorts: &policylangv1.LoadScheduler_Outs{
									ObservedLoadMultiplier: &policylangv1.OutPort{
										SignalName: "OBSERVED_LOAD_MULTIPLIER",
									},
								},
								DryRunConfigKey: adaptiveLoadScheduler.DryRunConfigKey,
								DryRun:          adaptiveLoadScheduler.DryRun,
								Parameters:      adaptiveLoadScheduler.Parameters.LoadScheduler,
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Decider{
					Decider: &policylangv1.Decider{
						InPorts: &policylangv1.Decider_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "DESIRED_LOAD_MULTIPLIER",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: 1,
										},
									},
								},
							},
						},
						Operator: components.LT.String(),
						OutPorts: &policylangv1.Decider_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "LOAD_MULTIPLIER_ALERT",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Alerter{
					Alerter: &policylangv1.Alerter{
						Parameters: alerterParameters,
						InPorts: &policylangv1.Alerter_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "LOAD_MULTIPLIER_ALERT",
								},
							},
						},
					},
				},
			},
		},
	}

	components.AddNestedIngress(nestedCircuit, alsSignalPortName, "SIGNAL")
	components.AddNestedIngress(nestedCircuit, alsSetpointPortName, "SETPOINT")
	components.AddNestedIngress(nestedCircuit, alsOverloadConfirmationPortName, "OVERLOAD_CONFIRMATION")
	components.AddNestedEgress(nestedCircuit, alsIsOverloadPortName, "IS_OVERLOAD")
	components.AddNestedEgress(nestedCircuit, alsDesiredLoadMultiplierPortName, "DESIRED_LOAD_MULTIPLIER")
	components.AddNestedEgress(nestedCircuit, alsObservedLoadMultiplierPortName, "OBSERVED_LOAD_MULTIPLIER")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("AdaptiveLoadScheduler",
			iface.GetSelectorsShortDescription(adaptiveLoadScheduler.Parameters.LoadScheduler.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		adaptiveLoadScheduler,
		componentID,
		false,
	)
	if err != nil {
		return retErr(err)
	}

	return configuredComponent, nestedCircuit, nil
}
