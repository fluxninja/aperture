package loadscheduler

import (
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

const (
	aimdSignalPortName                 = "signal"
	aimdSetpointPortName               = "setpoint"
	aimdOverloadConfirmationPortName   = "overload_confirmation"
	aimdIsOverloadPortName             = "is_overload"
	aimdDesiredLoadMultiplierPortName  = "desired_load_multiplier"
	aimdObservedLoadMultiplierPortName = "observed_load_multiplier"
)

// ParseAIMDLoadScheduler parses a AIMDLoadScheduler component and returns a configured component and a nested circuit.
func ParseAIMDLoadScheduler(
	aimdLoadScheduler *policylangv1.AIMDLoadScheduler,
	componentID runtime.ComponentID,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := aimdLoadScheduler.InPorts
	if inPorts != nil {
		signalPort := inPorts.Signal
		if signalPort != nil {
			nestedInPortsMap[aimdSignalPortName] = signalPort
		}
		setpointPort := inPorts.Setpoint
		if setpointPort != nil {
			nestedInPortsMap[aimdSetpointPortName] = setpointPort
		}
		overloadConfirmation := inPorts.OverloadConfirmation
		if overloadConfirmation != nil {
			nestedInPortsMap[aimdOverloadConfirmationPortName] = overloadConfirmation
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := aimdLoadScheduler.OutPorts
	if outPorts != nil {
		isOverloadPort := outPorts.IsOverload
		if isOverloadPort != nil {
			nestedOutPortsMap[aimdIsOverloadPortName] = isOverloadPort
		}
		desiredLoadMultiplierPort := outPorts.DesiredLoadMultiplier
		if desiredLoadMultiplierPort != nil {
			nestedOutPortsMap[aimdDesiredLoadMultiplierPortName] = desiredLoadMultiplierPort
		}
		observedLoadMultiplierPort := outPorts.ObservedLoadMultiplier
		if observedLoadMultiplierPort != nil {
			nestedOutPortsMap[aimdObservedLoadMultiplierPortName] = observedLoadMultiplierPort
		}
	}

	alerterLabels := aimdLoadScheduler.Parameters.Alerter.Labels
	if alerterLabels == nil {
		alerterLabels = make(map[string]string)
	}
	alerterLabels["type"] = "aimd_load_scheduler"
	aimdLoadScheduler.Parameters.Alerter.Labels = alerterLabels

	nestedCircuit := prepareLoadSchedulerCommonComponents()
	nestedCircuit.InPortsMap = nestedInPortsMap
	nestedCircuit.OutPortsMap = nestedOutPortsMap

	overloadDeciderOperator := components.GT.String()
	// if slope is greater than 0 then we want to use less than operator
	if aimdLoadScheduler.Parameters.Gradient.Slope > 0 {
		overloadDeciderOperator = components.LT.String()
	}

	nestedCircuit.Components = append(nestedCircuit.Components,
		&policylangv1.Component{
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
							DryRunConfigKey: aimdLoadScheduler.DryRunConfigKey,
							DryRun:          aimdLoadScheduler.DryRun,
							Parameters:      aimdLoadScheduler.Parameters.LoadScheduler,
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_Alerter{
				Alerter: &policylangv1.Alerter{
					Parameters: aimdLoadScheduler.Parameters.Alerter,
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
					Parameters: aimdLoadScheduler.Parameters.Gradient,
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
					InitialValue:       aimdLoadScheduler.Parameters.MaxLoadMultiplier,
					EvaluationInterval: durationpb.New(metricScrapeInterval),
					InPorts: &policylangv1.Integrator_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aimdLoadScheduler.Parameters.LoadMultiplierLinearIncrement,
									},
								},
							},
						},
						Max: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aimdLoadScheduler.Parameters.MaxLoadMultiplier,
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
										Value: aimdLoadScheduler.Parameters.MaxLoadMultiplier,
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

	components.AddNestedIngress(nestedCircuit, aimdSignalPortName, "SIGNAL")
	components.AddNestedIngress(nestedCircuit, aimdSetpointPortName, "SETPOINT")
	components.AddNestedIngress(nestedCircuit, aimdOverloadConfirmationPortName, "OVERLOAD_CONFIRMATION")
	components.AddNestedEgress(nestedCircuit, aimdIsOverloadPortName, "IS_OVERLOAD")
	components.AddNestedEgress(nestedCircuit, aimdDesiredLoadMultiplierPortName, "DESIRED_LOAD_MULTIPLIER")
	components.AddNestedEgress(nestedCircuit, aimdObservedLoadMultiplierPortName, "OBSERVED_LOAD_MULTIPLIER")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("AIMDLoadScheduler",
			iface.GetSelectorsShortDescription(aimdLoadScheduler.Parameters.LoadScheduler.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		aimdLoadScheduler,
		componentID,
		false,
	)
	if err != nil {
		return nil, nil, err
	}

	return configuredComponent, nestedCircuit, nil
}
