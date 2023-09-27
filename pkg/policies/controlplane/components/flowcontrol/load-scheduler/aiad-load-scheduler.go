package loadscheduler

import (
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

const (
	aiadSignalPortName                 = "signal"
	aiadSetpointPortName               = "setpoint"
	aiadOverloadConfirmationPortName   = "overload_confirmation"
	aiadIsOverloadPortName             = "is_overload"
	aiadDesiredLoadMultiplierPortName  = "desired_load_multiplier"
	aiadObservedLoadMultiplierPortName = "observed_load_multiplier"
)

// ParseAIADLoadScheduler parses a AIADLoadScheduler component and returns a configured component and a nested circuit.
func ParseAIADLoadScheduler(
	aiadLoadScheduler *policylangv1.AIADLoadScheduler,
	componentID runtime.ComponentID,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := aiadLoadScheduler.InPorts
	if inPorts != nil {
		signalPort := inPorts.Signal
		if signalPort != nil {
			nestedInPortsMap[aiadSignalPortName] = signalPort
		}
		setpointPort := inPorts.Setpoint
		if setpointPort != nil {
			nestedInPortsMap[aiadSetpointPortName] = setpointPort
		}
		overloadConfirmation := inPorts.OverloadConfirmation
		if overloadConfirmation != nil {
			nestedInPortsMap[aiadOverloadConfirmationPortName] = overloadConfirmation
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := aiadLoadScheduler.OutPorts
	if outPorts != nil {
		isOverloadPort := outPorts.IsOverload
		if isOverloadPort != nil {
			nestedOutPortsMap[aiadIsOverloadPortName] = isOverloadPort
		}
		desiredLoadMultiplierPort := outPorts.DesiredLoadMultiplier
		if desiredLoadMultiplierPort != nil {
			nestedOutPortsMap[aiadDesiredLoadMultiplierPortName] = desiredLoadMultiplierPort
		}
		observedLoadMultiplierPort := outPorts.ObservedLoadMultiplier
		if observedLoadMultiplierPort != nil {
			nestedOutPortsMap[aiadObservedLoadMultiplierPortName] = observedLoadMultiplierPort
		}
	}

	alerterLabels := aiadLoadScheduler.Parameters.Alerter.Labels
	if alerterLabels == nil {
		alerterLabels = make(map[string]string)
	}
	alerterLabels["type"] = "aiad_load_scheduler"
	aiadLoadScheduler.Parameters.Alerter.Labels = alerterLabels

	nestedCircuit := prepareLoadSchedulerCommonComponents()
	nestedCircuit.InPortsMap = nestedInPortsMap
	nestedCircuit.OutPortsMap = nestedOutPortsMap

	overloadDeciderOperator := aiadLoadScheduler.OverloadCondition

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
							DryRunConfigKey: aiadLoadScheduler.DryRunConfigKey,
							DryRun:          aiadLoadScheduler.DryRun,
							Parameters:      aiadLoadScheduler.Parameters.LoadScheduler,
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_Alerter{
				Alerter: &policylangv1.Alerter{
					Parameters: aiadLoadScheduler.Parameters.Alerter,
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
			Component: &policylangv1.Component_Switcher{
				Switcher: &policylangv1.Switcher{
					InPorts: &policylangv1.Switcher_Ins{
						Switch: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "IS_OVERLOAD",
							},
						},
						OnSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: (-1.0 * aiadLoadScheduler.Parameters.LoadMultiplierLinearDecrement),
									},
								},
							},
						},
						OffSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aiadLoadScheduler.Parameters.LoadMultiplierLinearIncrement,
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "LOAD_MULTIPLIER_DELTA",
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
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: (1.0 - aiadLoadScheduler.Parameters.LoadMultiplierLinearDecrement),
									},
								},
							},
						},
						OffSignal: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aiadLoadScheduler.Parameters.MaxLoadMultiplier,
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Switcher_Outs{
						Output: &policylangv1.OutPort{
							SignalName: "INTEGRATOR_MAX",
						},
					},
				},
			},
		},
		&policylangv1.Component{
			Component: &policylangv1.Component_Integrator{
				Integrator: &policylangv1.Integrator{
					InitialValue:       aiadLoadScheduler.Parameters.MaxLoadMultiplier,
					EvaluationInterval: durationpb.New(metricScrapeInterval),
					InPorts: &policylangv1.Integrator_Ins{
						Input: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "LOAD_MULTIPLIER_DELTA",
							},
						},
						Max: &policylangv1.InPort{
							Value: &policylangv1.InPort_SignalName{
								SignalName: "INTEGRATOR_MAX",
							},
						},
						Min: &policylangv1.InPort{
							Value: &policylangv1.InPort_ConstantSignal{
								ConstantSignal: &policylangv1.ConstantSignal{
									Const: &policylangv1.ConstantSignal_Value{
										Value: aiadLoadScheduler.Parameters.MinLoadMultiplier,
									},
								},
							},
						},
					},
					OutPorts: &policylangv1.Integrator_Outs{
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
										Value: aiadLoadScheduler.Parameters.MaxLoadMultiplier,
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

	components.AddNestedIngress(nestedCircuit, aiadSignalPortName, "SIGNAL")
	components.AddNestedIngress(nestedCircuit, aiadSetpointPortName, "SETPOINT")
	components.AddNestedIngress(nestedCircuit, aiadOverloadConfirmationPortName, "OVERLOAD_CONFIRMATION")
	components.AddNestedEgress(nestedCircuit, aiadIsOverloadPortName, "IS_OVERLOAD")
	components.AddNestedEgress(nestedCircuit, aiadDesiredLoadMultiplierPortName, "DESIRED_LOAD_MULTIPLIER")
	components.AddNestedEgress(nestedCircuit, aiadObservedLoadMultiplierPortName, "OBSERVED_LOAD_MULTIPLIER")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("AIADLoadScheduler",
			iface.GetSelectorsShortDescription(aiadLoadScheduler.Parameters.LoadScheduler.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		aiadLoadScheduler,
		componentID,
		false,
	)
	if err != nil {
		return nil, nil, err
	}

	return configuredComponent, nestedCircuit, nil
}
