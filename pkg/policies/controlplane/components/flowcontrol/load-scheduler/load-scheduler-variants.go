package loadscheduler

import (
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
)

func prepareLoadSchedulerCommonComponents() *policylangv1.NestedCircuit {
	nestedCircuit := &policylangv1.NestedCircuit{
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
												// Overload confirmation is assumed true by default. This allows the circuit to function even if
												// overload confirmation is not provided. If the required behavior is to assume false by default,
												// the policy needs to ensure a valid signal with the desired defaults is provided.
												Value: 1,
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
		},
	}

	return nestedCircuit
}
