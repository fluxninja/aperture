package loadscheduler

import (
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/private/v1"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	inputLoadMultiplierPortName          = "load_multiplier"
	outputObservedLoadMultiplierPortName = "observed_load_multiplier"
)

// ParseLoadScheduler parses a load scheduler from a spec.
func ParseLoadScheduler(
	loadScheduler *policylangv1.LoadScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := loadScheduler.GetInPorts()

	if inPorts != nil {
		loadMultiplierPort := inPorts.LoadMultiplier
		if loadMultiplierPort != nil {
			nestedInPortsMap[inputLoadMultiplierPortName] = loadMultiplierPort
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := loadScheduler.GetOutPorts()
	if outPorts != nil {
		observedLoadMultiplierPort := outPorts.ObservedLoadMultiplier
		if observedLoadMultiplierPort != nil {
			nestedOutPortsMap[outputObservedLoadMultiplierPortName] = observedLoadMultiplierPort
		}
	}

	// Prepare parameters for prometheus queries
	policyParams := fmt.Sprintf("%s=\"%s\",%s=\"%s\",%s=\"%s\"",
		metrics.PolicyNameLabel,
		policyReadAPI.GetPolicyName(),
		metrics.PolicyHashLabel,
		policyReadAPI.GetPolicyHash(),
		metrics.ComponentIDLabel,
		componentID,
	)

	acceptedTokensQuery := fmt.Sprintf("sum(rate(%s{%s}[10s]))",
		metrics.AcceptedTokensMetricName,
		policyParams)

	incomingTokenRate := fmt.Sprintf("sum(rate(%s{%s}[10s]))",
		metrics.IncomingTokensMetricName,
		policyParams)

	loadActuatorAnyProto, err := anypb.New(
		&policyprivatev1.LoadActuator{
			InPorts: &policyprivatev1.LoadActuator_Ins{
				LoadMultiplier: &policylangv1.InPort{
					Value: &policylangv1.InPort_SignalName{
						SignalName: "LOAD_MULTIPLIER",
					},
				},
			},
			LoadSchedulerComponentId:   componentID.String(),
			DefaultConfig:              loadScheduler.GetDefaultConfig(),
			DynamicConfigKey:           loadScheduler.GetDynamicConfigKey(),
			WorkloadLatencyBasedTokens: loadScheduler.Parameters.GetWorkloadLatencyBasedTokens(),
			Selectors:                  loadScheduler.Parameters.GetSelectors(),
		})
	if err != nil {
		return nil, err
	}

	nestedCircuit := &policylangv1.NestedCircuit{
		Name:             "LoadScheduler",
		ShortDescription: iface.GetSelectorsShortDescription(loadScheduler.Parameters.GetSelectors()),
		InPortsMap:       nestedInPortsMap,
		OutPortsMap:      nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_FlowControl{
					FlowControl: &policylangv1.FlowControl{
						Component: &policylangv1.FlowControl_Private{
							Private: loadActuatorAnyProto,
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Query{
					Query: &policylangv1.Query{
						Component: &policylangv1.Query_Promql{
							Promql: &policylangv1.PromQL{
								OutPorts: &policylangv1.PromQL_Outs{
									Output: &policylangv1.OutPort{
										SignalName: "ACCEPTED_TOKEN_RATE",
									},
								},
								QueryString:        acceptedTokensQuery,
								EvaluationInterval: durationpb.New(policyReadAPI.GetEvaluationInterval()),
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_Query{
					Query: &policylangv1.Query{
						Component: &policylangv1.Query_Promql{
							Promql: &policylangv1.PromQL{
								OutPorts: &policylangv1.PromQL_Outs{
									Output: &policylangv1.OutPort{
										SignalName: "INCOMING_TOKEN_RATE",
									},
								},
								QueryString:        incomingTokenRate,
								EvaluationInterval: durationpb.New(policyReadAPI.GetEvaluationInterval()),
							},
						},
					},
				},
			},
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
		},
	}

	components.AddNestedIngress(nestedCircuit, inputLoadMultiplierPortName, "LOAD_MULTIPLIER")
	components.AddNestedEgress(nestedCircuit, outputObservedLoadMultiplierPortName, "OBSERVED_LOAD_MULTIPLIER")

	return nestedCircuit, nil
}
