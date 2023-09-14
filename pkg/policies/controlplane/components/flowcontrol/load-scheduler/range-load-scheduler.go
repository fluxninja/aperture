package loadscheduler

import (
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

const (
	rangeSignalPortName                 = "signal"
	rangeOverloadConfirmationPortName   = "overload_confirmation"
	rangeIsOverloadPortName             = "is_overload"
	rangeDesiredLoadMultiplierPortName  = "desired_load_multiplier"
	rangeObservedLoadMultiplierPortName = "observed_load_multiplier"
)

// ParseRangeDrivenLoadScheduler parses a range driven load scheduler component.
func ParseRangeDrivenLoadScheduler(
	rangeDrivenLoadScheduler *policylangv1.RangeDrivenLoadScheduler,
	componentID runtime.ComponentID,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	retErr := func(err error) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
		return nil, nil, err
	}

	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := rangeDrivenLoadScheduler.InPorts
	if inPorts != nil {
		signalPort := inPorts.Signal
		if signalPort != nil {
			nestedInPortsMap[rangeSignalPortName] = signalPort
		}
		overloadConfirmation := inPorts.OverloadConfirmation
		if overloadConfirmation != nil {
			nestedInPortsMap[rangeOverloadConfirmationPortName] = overloadConfirmation
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := rangeDrivenLoadScheduler.OutPorts
	if outPorts != nil {
		isOverloadPort := outPorts.IsOverload
		if isOverloadPort != nil {
			nestedOutPortsMap[rangeIsOverloadPortName] = isOverloadPort
		}
		desiredLoadMultiplierPort := outPorts.DesiredLoadMultiplier
		if desiredLoadMultiplierPort != nil {
			nestedOutPortsMap[rangeDesiredLoadMultiplierPortName] = desiredLoadMultiplierPort
		}
		observedLoadMultiplierPort := outPorts.ObservedLoadMultiplier
		if observedLoadMultiplierPort != nil {
			nestedOutPortsMap[rangeObservedLoadMultiplierPortName] = observedLoadMultiplierPort
		}
	}

	alerterLabels := rangeDrivenLoadScheduler.Parameters.Alerter.Labels
	if alerterLabels == nil {
		alerterLabels = make(map[string]string)
	}
	alerterLabels["type"] = "range_driven_load_scheduler"
	rangeDrivenLoadScheduler.Parameters.Alerter.Labels = alerterLabels

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
								DryRunConfigKey: rangeDrivenLoadScheduler.DryRunConfigKey,
								DryRun:          rangeDrivenLoadScheduler.DryRun,
								Parameters:      rangeDrivenLoadScheduler.Parameters.LoadScheduler,
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
						Parameters: rangeDrivenLoadScheduler.Parameters.Alerter,
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

	parameters := rangeDrivenLoadScheduler.Parameters
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

	components.AddNestedIngress(nestedCircuit, rangeSignalPortName, "SIGNAL")
	components.AddNestedIngress(nestedCircuit, rangeOverloadConfirmationPortName, "OVERLOAD_CONFIRMATION")
	components.AddNestedEgress(nestedCircuit, rangeIsOverloadPortName, "IS_OVERLOAD")
	components.AddNestedEgress(nestedCircuit, rangeDesiredLoadMultiplierPortName, "DESIRED_LOAD_MULTIPLIER")
	components.AddNestedEgress(nestedCircuit, rangeObservedLoadMultiplierPortName, "OBSERVED_LOAD_MULTIPLIER")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("RangeDrivenLoadScheduler",
			iface.GetSelectorsShortDescription(rangeDrivenLoadScheduler.Parameters.LoadScheduler.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		rangeDrivenLoadScheduler,
		componentID,
		false,
	)
	if err != nil {
		return retErr(err)
	}

	return configuredComponent, nestedCircuit, nil
}
