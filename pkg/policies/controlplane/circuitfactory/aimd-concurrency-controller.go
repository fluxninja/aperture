package circuitfactory

import (
	"time"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	aimdSignalPortName         = "signal"
	aimdSetpointPortName       = "setpoint"
	aimdIsOverloadPortName     = "is_overload"
	aimdLoadMultiplierPortName = "load_multiplier"
)

// ParseAIMDConcurrencyController parses an AIMD concurrency controller and returns the parent, leaf components, and options.
func ParseAIMDConcurrencyController(
	nestedCircuitID runtime.ComponentID,
	aimdConcurrencyController *policylangv1.AIMDConcurrencyController,
	policyReadAPI iface.Policy,
) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := aimdConcurrencyController.InPorts
	if inPorts != nil {
		signalPort := inPorts.Signal
		if signalPort != nil {
			nestedInPortsMap[aimdSignalPortName] = signalPort
		}
		setpointPort := inPorts.Setpoint
		if setpointPort != nil {
			nestedInPortsMap[aimdSetpointPortName] = setpointPort
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := aimdConcurrencyController.OutPorts
	if outPorts != nil {
		isOverloadPort := outPorts.IsOverload
		if isOverloadPort != nil {
			nestedOutPortsMap[aimdIsOverloadPortName] = isOverloadPort
		}
		loadMultiplierPort := outPorts.LoadMultiplier
		if loadMultiplierPort != nil {
			nestedOutPortsMap[aimdLoadMultiplierPortName] = loadMultiplierPort
		}
	}

	isOverloadDeciderOperator := "gt"
	// if slope is greater than 0 then we want to use less than operator
	if aimdConcurrencyController.GradientParameters.Slope > 0 {
		isOverloadDeciderOperator = "lt"
	}

	alerterLabels := aimdConcurrencyController.AlerterParameters.Labels
	if alerterLabels == nil {
		alerterLabels = make(map[string]string)
	}
	alerterLabels["type"] = "concurrency_limiter"
	alerterLabels["agent_group"] = aimdConcurrencyController.FlowSelector.ServiceSelector.GetAgentGroup()
	alerterLabels["service"] = aimdConcurrencyController.FlowSelector.ServiceSelector.GetService()
	aimdConcurrencyController.AlerterParameters.Labels = alerterLabels

	nestedCircuit := &policylangv1.NestedCircuit{
		Name:             "AIMDConcurrencyController",
		ShortDescription: iface.GetServiceShortDescription(aimdConcurrencyController.FlowSelector.ServiceSelector),
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
									SignalName: "DESIRED_CONCURRENCY",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "INCOMING_CONCURRENCY",
								},
							},
						},
						OutPorts: &policylangv1.ArithmeticCombinator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "DESIRED_CONCURRENCY_RATIO",
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
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: aimdConcurrencyController.ConcurrencyLimitMultiplier,
										},
									},
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACCEPTED_CONCURRENCY",
								},
							},
						},
						OutPorts: &policylangv1.ArithmeticCombinator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "NORMAL_CONCURRENCY_LIMIT",
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
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: aimdConcurrencyController.ConcurrencyLinearIncrement,
										},
									},
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "SQRT_CONCURRENCY_INCREMENT",
								},
							},
						},
						OutPorts: &policylangv1.ArithmeticCombinator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "CONCURRENCY_INCREMENT_SINGLE_TICK",
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
									SignalName: "CONCURRENCY_INCREMENT_SINGLE_TICK",
								},
							},
							Rhs: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "CONCURRENCY_INCREMENT",
								},
							},
						},
						OutPorts: &policylangv1.ArithmeticCombinator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "CONCURRENCY_INCREMENT_INTEGRAL",
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
										SignalName: "CONCURRENCY_INCREMENT_INTEGRAL",
									},
								},
								{
									Value: &policylangv1.InPort_SignalName{
										SignalName: "ACCEPTED_CONCURRENCY",
									},
								},
							},
						},
						OutPorts: &policylangv1.Min_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "CONCURRENCY_INCREMENT_INTEGRAL_CAPPED",
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
										SignalName: "CONCURRENCY_INCREMENT_INTEGRAL_CAPPED",
									},
								},
								{
									Value: &policylangv1.InPort_ConstantSignal{
										ConstantSignal: &policylangv1.ConstantSignal{
											Const: &policylangv1.ConstantSignal_Value{
												Value: 0,
											},
										},
									},
								},
							},
						},
						OutPorts: &policylangv1.FirstValid_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "CONCURRENCY_INCREMENT_NORMAL",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Sqrt{
					Sqrt: &policylangv1.Sqrt{
						Scale: aimdConcurrencyController.ConcurrencySqrtIncrementMultiplier,
						InPorts: &policylangv1.Sqrt_Ins{
							Input: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACCEPTED_CONCURRENCY",
								},
							},
						},
						OutPorts: &policylangv1.Sqrt_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "SQRT_CONCURRENCY_INCREMENT",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_GradientController{
					GradientController: &policylangv1.GradientController{
						Parameters: aimdConcurrencyController.GradientParameters,
						InPorts: &policylangv1.GradientController_Ins{
							ControlVariable: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "ACCEPTED_CONCURRENCY",
								},
							},
							Max: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "NORMAL_CONCURRENCY_LIMIT",
								},
							},
							Optimize: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "CONCURRENCY_INCREMENT",
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
								SignalName: "DESIRED_CONCURRENCY",
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
									SignalName: "DESIRED_CONCURRENCY_RATIO",
								},
							},
						},
						OutPorts: &policylangv1.Extrapolator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "LOAD_MULTIPLIER",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_FlowControl{
					FlowControl: &policylangv1.FlowControl{
						Component: &policylangv1.FlowControl_ConcurrencyLimiter{
							ConcurrencyLimiter: &policylangv1.ConcurrencyLimiter{
								FlowSelector: aimdConcurrencyController.FlowSelector,
								Scheduler: &policylangv1.Scheduler{
									Parameters: aimdConcurrencyController.SchedulerParameters,
									OutPorts: &policylangv1.Scheduler_Outs{
										AcceptedConcurrency: &policylangv1.OutPort{
											SignalName: "ACCEPTED_CONCURRENCY",
										},
										IncomingConcurrency: &policylangv1.OutPort{
											SignalName: "INCOMING_CONCURRENCY",
										},
									},
								},
								ActuationStrategy: &policylangv1.ConcurrencyLimiter_LoadActuator{
									LoadActuator: &policylangv1.LoadActuator{
										DynamicConfigKey: aimdConcurrencyController.DryRunDynamicConfigKey,
										InPorts: &policylangv1.LoadActuator_Ins{
											LoadMultiplier: &policylangv1.InPort{
												Value: &policylangv1.InPort_SignalName{
													SignalName: "LOAD_MULTIPLIER",
												},
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
									SignalName: "LOAD_MULTIPLIER",
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
						Parameters: aimdConcurrencyController.AlerterParameters,
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
							OnFalse: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "CONCURRENCY_INCREMENT_NORMAL",
								},
							},
							OnTrue: &policylangv1.InPort{
								Value: &policylangv1.InPort_ConstantSignal{
									ConstantSignal: &policylangv1.ConstantSignal{
										Const: &policylangv1.ConstantSignal_Value{
											Value: 0,
										},
									},
								},
							},
							Switch: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "IS_OVERLOAD",
								},
							},
						},
						OutPorts: &policylangv1.Switcher_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "CONCURRENCY_INCREMENT",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_NestedSignalIngress{
					NestedSignalIngress: &policylangv1.NestedSignalIngress{
						PortName: aimdSignalPortName,
						OutPorts: &policylangv1.NestedSignalIngress_Outs{
							Signal: &policylangv1.OutPort{
								SignalName: "SIGNAL",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_NestedSignalIngress{
					NestedSignalIngress: &policylangv1.NestedSignalIngress{
						PortName: aimdSetpointPortName,
						OutPorts: &policylangv1.NestedSignalIngress_Outs{
							Signal: &policylangv1.OutPort{
								SignalName: "SETPOINT",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_NestedSignalEgress{
					NestedSignalEgress: &policylangv1.NestedSignalEgress{
						PortName: aimdIsOverloadPortName,
						InPorts: &policylangv1.NestedSignalEgress_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "IS_OVERLOAD",
								},
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_NestedSignalEgress{
					NestedSignalEgress: &policylangv1.NestedSignalEgress{
						PortName: aimdLoadMultiplierPortName,
						InPorts: &policylangv1.NestedSignalEgress_Ins{
							Signal: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "LOAD_MULTIPLIER",
								},
							},
						},
					},
				},
			},
		},
	}

	return ParseNestedCircuit(nestedCircuitID, nestedCircuit, policyReadAPI)
}
