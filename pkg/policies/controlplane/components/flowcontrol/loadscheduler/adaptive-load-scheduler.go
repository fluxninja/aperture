package loadscheduler

import (
	"time"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	alsSignalPortName                 = "signal"
	alsSetpointPortName               = "setpoint"
	alsOverloadConfirmationPortName   = "overload_confirmation"
	alsIsOverloadPortName             = "is_overload"
	alsDesiredLoadMultiplierPortName  = "desired_load_multiplier"
	alsObservedLoadMultiplierPortName = "observed_load_multiplier"
	alsAcceptedTokenRatePortName      = "accepted_token_rate"
	alsIncomingTokenRatePortName      = "incoming_token_rate"
)

// ParseAdaptiveLoadScheduler parses and returns nested circuit representation of AdaptiveLoadScheduler.
func ParseAdaptiveLoadScheduler(
	adaptiveLoadScheduler *policylangv1.AdaptiveLoadScheduler,
) (*policylangv1.NestedCircuit, error) {
	// Deprecated 1.8.0
	flowSelectorProto := adaptiveLoadScheduler.GetFlowSelector()
	if flowSelectorProto != nil {
		selector := &policylangv1.Selector{
			ControlPoint: flowSelectorProto.FlowMatcher.ControlPoint,
			LabelMatcher: flowSelectorProto.FlowMatcher.LabelMatcher,
			Service:      flowSelectorProto.ServiceSelector.Service,
			AgentGroup:   flowSelectorProto.ServiceSelector.AgentGroup,
		}
		adaptiveLoadScheduler.Selectors = append(adaptiveLoadScheduler.Selectors, selector)
		adaptiveLoadScheduler.FlowSelector = nil
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
		acceptedTokenRatePort := outPorts.AcceptedTokenRate
		if acceptedTokenRatePort != nil {
			nestedOutPortsMap[alsAcceptedTokenRatePortName] = acceptedTokenRatePort
		}
		incomingTokenRatePort := outPorts.IncomingTokenRate
		if incomingTokenRatePort != nil {
			nestedOutPortsMap[alsIncomingTokenRatePortName] = incomingTokenRatePort
		}
	}

	isOverloadDeciderOperator := "gt"
	// if slope is greater than 0 then we want to use less than operator
	if adaptiveLoadScheduler.GradientParameters.Slope > 0 {
		isOverloadDeciderOperator = "lt"
	}

	alerterLabels := adaptiveLoadScheduler.AlerterParameters.Labels
	if alerterLabels == nil {
		alerterLabels = make(map[string]string)
	}
	alerterLabels["type"] = "load_scheduler"
	adaptiveLoadScheduler.AlerterParameters.Labels = alerterLabels

	nestedCircuit := &policylangv1.NestedCircuit{
		Name:             "AdaptiveLoadScheduler",
		ShortDescription: iface.GetSelectorsShortDescription(adaptiveLoadScheduler.GetSelectors()),
		InPortsMap:       nestedInPortsMap,
		OutPortsMap:      nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_ArithmeticCombinator{
					ArithmeticCombinator: &policylangv1.ArithmeticCombinator{
						Operator: "div",
						InPorts: &policylangv1.ArithmeticCombinator_Ins{
							Lhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACCEPTED_TOKEN_RATE",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "INCOMING_TOKEN_RATE",
								},
							},
						},
						OutPorts: &policylangv1.ArithmeticCombinator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "OBSERVED_LOAD_MULTIPLIER",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_GradientController{
					GradientController: &policylangv1.GradientController{
						Parameters: adaptiveLoadScheduler.GradientParameters,
						InPorts: &policylangv1.GradientController_Ins{
							ControlVariable: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "OBSERVED_LOAD_MULTIPLIER",
								},
							},
							Max: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: adaptiveLoadScheduler.MaxLoadMultiplier,
										},
									},
								},
							},
							Optimize: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "LOAD_MULTIPLIER_INCREMENT",
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
								SignalName: "DESIRED_LOAD_MULTIPLIER",
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
										SignalName: "OVERLOAD_CONFIRMATION",
									},
								},
								{
									Value: &policylangv1.InPort_ConstantSignal{
										ConstantSignal: &policylangv1.ConstantSignal{
											Const: &policylangv1.ConstantSignal_Value{
												Value: 1, // OVERLOAD_CONFIRMATION is true by default
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
				Component: &policylangv1.Component_FlowControl{
					FlowControl: &policylangv1.FlowControl{
						Component: &policylangv1.FlowControl_LoadScheduler{
							LoadScheduler: &policylangv1.LoadScheduler{
								Selectors: adaptiveLoadScheduler.Selectors,
								Scheduler: &policylangv1.LoadScheduler_Scheduler{
									Parameters: adaptiveLoadScheduler.SchedulerParameters,
									OutPorts: &policylangv1.LoadScheduler_Scheduler_Outs{
										AcceptedTokenRate: &policylangv1.OutPort{
											SignalName: "ACCEPTED_TOKEN_RATE",
										},
										IncomingTokenRate: &policylangv1.OutPort{
											SignalName: "INCOMING_TOKEN_RATE",
										},
									},
								},
								Actuator: &policylangv1.LoadScheduler_Actuator{
									DynamicConfigKey: adaptiveLoadScheduler.DynamicConfigKey,
									DefaultConfig:    adaptiveLoadScheduler.DefaultConfig,
									InPorts: &policylangv1.LoadScheduler_Actuator_Ins{
										LoadMultiplier: &policylangv1.InPort{
											Value: &policylangv1.InPort_SignalName{
												SignalName: "DESIRED_LOAD_MULTIPLIER",
											},
										},
									},
								},
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
						Operator: "lt",
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
						Parameters: adaptiveLoadScheduler.AlerterParameters,
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
			{
				Component: &policylangv1.Component_Decider{
					Decider: &policylangv1.Decider{
						Operator: isOverloadDeciderOperator,
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
				Component: &policylangv1.Component_Integrator{
					Integrator: &policylangv1.Integrator{
						InPorts: &policylangv1.Integrator_Ins{
							Input: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: adaptiveLoadScheduler.LoadMultiplierLinearIncrement,
										},
									},
								},
							},
							Max: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: adaptiveLoadScheduler.MaxLoadMultiplier,
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
								SignalName: "LOAD_MULTIPLIER_INCREMENT",
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
	components.AddNestedEgress(nestedCircuit, alsAcceptedTokenRatePortName, "ACCEPTED_TOKEN_RATE")
	components.AddNestedEgress(nestedCircuit, alsIncomingTokenRatePortName, "INCOMING_TOKEN_RATE")

	return nestedCircuit, nil
}
