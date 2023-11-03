package loadscheduler

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/private/v1"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

const (
	inputLoadMultiplierPortName          = "load_multiplier"
	outputObservedLoadMultiplierPortName = "observed_load_multiplier"
	// TODO: move to a common area so that it is shared with OTEL collector.
	metricScrapeInterval = time.Second * 10
)

// ParseLoadScheduler parses a load scheduler from a spec.
func ParseLoadScheduler(
	loadScheduler *policylangv1.LoadScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	retErr := func(err error) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
		return nil, nil, err
	}

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
	policyParams := fmt.Sprintf(
		"%s=\"%s\",%s=\"%s\",%s=\"%s\"",
		metrics.PolicyNameLabel,
		policyReadAPI.GetPolicyName(),
		metrics.PolicyHashLabel,
		policyReadAPI.GetPolicyHash(),
		metrics.ComponentIDLabel,
		componentID,
	)

	// TODO: 30s is derived from Agent's load multiplier token bucket's window size. Define the constant at a common location.
	observedLoadMultiplierQuery := fmt.Sprintf(
		"sum(increase(%s{%s}[30s])) / (sum(increase(%s{%s}[30s])) + sum(increase(%s{%s}[30s])))",
		metrics.AcceptedTokensMetricName,
		policyParams,
		metrics.AcceptedTokensMetricName,
		policyParams,
		metrics.RejectedTokensMetricName,
		policyParams,
	)

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
			WorkloadLatencyBasedTokens: loadScheduler.Parameters.GetWorkloadLatencyBasedTokens(),
			Selectors:                  loadScheduler.Parameters.GetSelectors(),
		})
	if err != nil {
		return retErr(err)
	}

	nestedCircuit := &policylangv1.NestedCircuit{
		InPortsMap:  nestedInPortsMap,
		OutPortsMap: nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_BoolVariable{
					BoolVariable: &policylangv1.BoolVariable{
						ConstantOutput: loadScheduler.GetDryRun(),
						ConfigKey:      loadScheduler.GetDryRunConfigKey(),
						OutPorts: &policylangv1.BoolVariable_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "DRY_RUN",
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
									SignalName: "DRY_RUN",
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
									SignalName: "LOAD_MULTIPLIER_INPUT",
								},
							},
						},
						OutPorts: &policylangv1.Switcher_Outs{
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
										SignalName: "OBSERVED_LOAD_MULTIPLIER",
									},
								},
								QueryString:        observedLoadMultiplierQuery,
								EvaluationInterval: durationpb.New(metricScrapeInterval),
							},
						},
					},
				},
			},
		},
	}

	components.AddNestedIngress(nestedCircuit, inputLoadMultiplierPortName, "LOAD_MULTIPLIER_INPUT")
	components.AddNestedEgress(nestedCircuit, outputObservedLoadMultiplierPortName, "OBSERVED_LOAD_MULTIPLIER")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("LoadScheduler",
			iface.GetSelectorsShortDescription(loadScheduler.Parameters.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		loadScheduler,
		componentID,
		false,
	)
	if err != nil {
		return retErr(err)
	}

	return configuredComponent, nestedCircuit, nil
}
