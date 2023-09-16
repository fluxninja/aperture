package loadscheduler

import (
	"errors"

	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
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

	// Convert pre-throttling strategy fields to the new throttling strategy
	convertPreThrottlingStrategyFields(adaptiveLoadScheduler)

	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := adaptiveLoadScheduler.InPorts
	if inPorts != nil {
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

	alerterLabels := adaptiveLoadScheduler.Parameters.Alerter.Labels
	if alerterLabels == nil {
		alerterLabels = make(map[string]string)
	}
	alerterLabels["type"] = "load_scheduler"
	adaptiveLoadScheduler.Parameters.Alerter.Labels = alerterLabels

	// Needs PASS_THROUGH_FROM_STRATEGY, LOAD_MULTIPLIER_FROM_STRATEGY and OVERLOAD_FROM_STRATEGY
	// Provides IS_OVERLOAD, DESIRED_LOAD_MULTIPLIER, OBSERVED_LOAD_MULTIPLIER
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
				Component: &policylangv1.Component_And{
					And: &policylangv1.And{
						InPorts: &policylangv1.And_Ins{
							Inputs: []*policylangv1.InPort{
								{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "OVERLOAD_FROM_STRATEGY",
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
				Component: &policylangv1.Component_Switcher{
					Switcher: &policylangv1.Switcher{
						InPorts: &policylangv1.Switcher_Ins{
							Switch: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "PASS_THROUGH_FROM_STRATEGY",
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
									SignalName: "LOAD_MULTIPLIER_FROM_STRATEGY",
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
						Parameters: adaptiveLoadScheduler.Parameters.Alerter,
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

	if adaptiveLoadScheduler.GetAimdThrottlingStrategy() != nil {
		parseAimdThrottlingStrategy(adaptiveLoadScheduler.GetAimdThrottlingStrategy(), nestedCircuit, nestedInPortsMap)
	} else if adaptiveLoadScheduler.GetRangeThrottlingStrategy() != nil {
		parseRangeThrottlingStrategy(adaptiveLoadScheduler.GetRangeThrottlingStrategy(), nestedCircuit, nestedInPortsMap)
	} else {
		return retErr(errors.New("no throttling strategy provided"))
	}

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

func parseAimdThrottlingStrategy(
	aimdThrottling *policylangv1.AdaptiveLoadScheduler_AIMDThrottlingStrategy,
	nestedCircuit *policylangv1.NestedCircuit,
	nestedInPortsMap map[string]*policylangv1.InPort,
) {
	inPorts := aimdThrottling.InPorts
	if inPorts != nil {
		signalPort := inPorts.Signal
		if signalPort != nil {
			nestedInPortsMap[alsSignalPortName] = signalPort
		}
		setpointPort := inPorts.Setpoint
		if setpointPort != nil {
			nestedInPortsMap[alsSetpointPortName] = setpointPort
		}
	}
	components.AddNestedIngress(nestedCircuit, alsSignalPortName, "SIGNAL")
	components.AddNestedIngress(nestedCircuit, alsSetpointPortName, "SETPOINT")

	overloadDeciderOperator := components.GT.String()
	// if slope is greater than 0 then we want to use less than operator
	if aimdThrottling.Gradient.Slope > 0 {
		overloadDeciderOperator = components.LT.String()
	}

	nestedCircuit.Components = append(nestedCircuit.Components,
		&policylangv1.Component{
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
							SignalName: "OVERLOAD_FROM_STRATEGY",
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_GradientController{
				GradientController: &policylangv1.GradientController{
					Parameters: aimdThrottling.Gradient,
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
							SignalName: "LOAD_MULTIPLIER_IF_OVERLOAD",
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_Integrator{
				Integrator: &policylangv1.Integrator{
					InitialValue:       aimdThrottling.MaxLoadMultiplier,
					EvaluationInterval: durationpb.New(metricScrapeInterval),
					InPorts: &policylangv1.Integrator_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aimdThrottling.LoadMultiplierLinearIncrement,
									},
								},
							},
						},
						Max: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aimdThrottling.MaxLoadMultiplier,
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
		&policylangv1.Component{
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
		&policylangv1.Component{
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
							SignalName: "LOAD_MULTIPLIER_FROM_STRATEGY",
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_Decider{
				Decider: &policylangv1.Decider{
					Operator: components.GTE.String(),
					InPorts: &policylangv1.Decider_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "LOAD_MULTIPLIER_FROM_STRATEGY",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aimdThrottling.MaxLoadMultiplier,
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Decider_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "PASS_THROUGH_FROM_STRATEGY",
						},
					},
				},
			},
		},
	)
}

func parseRangeThrottlingStrategy(rangeThrottling *policylangv1.AdaptiveLoadScheduler_RangeThrottlingStrategy, nestedCircuit *policylangv1.NestedCircuit, nestedInPortsMap map[string]*policylangv1.InPort) {
	inPorts := rangeThrottling.InPorts
	if inPorts != nil {
		signalPort := inPorts.Signal
		if signalPort != nil {
			nestedInPortsMap[alsSignalPortName] = signalPort
		}
	}
	components.AddNestedIngress(nestedCircuit, alsSignalPortName, "SIGNAL")

	parameters := rangeThrottling.Parameters
	preStart := 1.0
	postEnd := parameters.End.LoadMultiplier
	if parameters.Start.LoadMultiplier < parameters.End.LoadMultiplier {
		preStart = parameters.Start.LoadMultiplier
		postEnd = 1.0
	}

	nestedCircuit.Components = append(nestedCircuit.Components,
		&policylangv1.Component{
			Component: &policylangv1.Component_PolynomialRangeFunction{
				PolynomialRangeFunction: &policylangv1.PolynomialRangeFunction{
					Parameters: &policylangv1.PolynomialRangeFunction_Parameters{
						Start: &policylangv1.PolynomialRangeFunction_Parameters_Datapoint{
							Input:  parameters.Start.Threshold,
							Output: parameters.Start.LoadMultiplier,
						},
						End: &policylangv1.PolynomialRangeFunction_Parameters_Datapoint{
							Input:  parameters.End.Threshold,
							Output: parameters.End.LoadMultiplier,
						},
						OutsideRange: &policylangv1.PolynomialRangeFunction_Parameters_ClampToCustomValues_{
							ClampToCustomValues: &policylangv1.PolynomialRangeFunction_Parameters_ClampToCustomValues{
								PreStart: preStart,
								PostEnd:  postEnd,
							},
						},
					},
					InPorts: &policylangv1.PolynomialRangeFunction_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "SIGNAL",
							},
						},
					},
					OutPorts: &policylangv1.PolynomialRangeFunction_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "LOAD_MULTIPLIER_FROM_STRATEGY",
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_Decider{
				Decider: &policylangv1.Decider{
					Operator: components.GTE.String(),
					InPorts: &policylangv1.Decider_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "LOAD_MULTIPLIER_FROM_STRATEGY",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: 1.0,
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Decider_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "PASS_THROUGH_FROM_STRATEGY",
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_Decider{
				Decider: &policylangv1.Decider{
					Operator: components.LT.String(),
					InPorts: &policylangv1.Decider_Ins{
						Lhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "LOAD_MULTIPLIER_FROM_STRATEGY",
							},
						},
						Rhs: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: 1.0,
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Decider_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "OVERLOAD_FROM_STRATEGY",
						},
					},
				},
			},
		},
	)
}

func convertPreThrottlingStrategyFields(adaptiveLoadScheduler *policylangv1.AdaptiveLoadScheduler) {
	// If a throttling strategy is provided, then we are using the latest version of this component
	if adaptiveLoadScheduler.GetThrottlingStrategy() != nil {
		return
	}
	// Set Gradient Throttling Strategy
	adaptiveLoadScheduler.ThrottlingStrategy = &policylangv1.AdaptiveLoadScheduler_AimdThrottlingStrategy{
		AimdThrottlingStrategy: &policylangv1.AdaptiveLoadScheduler_AIMDThrottlingStrategy{
			Gradient:                      adaptiveLoadScheduler.Parameters.Gradient,
			MaxLoadMultiplier:             adaptiveLoadScheduler.Parameters.MaxLoadMultiplier,
			LoadMultiplierLinearIncrement: adaptiveLoadScheduler.Parameters.LoadMultiplierLinearIncrement,
			InPorts:                       &policylangv1.AdaptiveLoadScheduler_AIMDThrottlingStrategy_Ins{},
		},
	}
	convertPreThrottlingStrategyInPorts(adaptiveLoadScheduler.InPorts, adaptiveLoadScheduler.ThrottlingStrategy.(*policylangv1.AdaptiveLoadScheduler_AimdThrottlingStrategy).AimdThrottlingStrategy.InPorts)
}

func convertPreThrottlingStrategyInPorts(adaptiveLoadSchedulerInPorts *policylangv1.AdaptiveLoadScheduler_Ins, gradientThrottlingStrategyInPorts *policylangv1.AdaptiveLoadScheduler_AIMDThrottlingStrategy_Ins) {
	if adaptiveLoadSchedulerInPorts == nil {
		return
	}
	gradientThrottlingStrategyInPorts.Signal = adaptiveLoadSchedulerInPorts.Signal
	gradientThrottlingStrategyInPorts.Setpoint = adaptiveLoadSchedulerInPorts.Setpoint
}
