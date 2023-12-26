package concurrencyscheduler

import (
	"fmt"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/private/v1"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

const (
	maxConcurrencyPortName = "max_capacity"
	passThroughPortName    = "pass_through"
	acceptPercentageName   = "accept_percentage"
)

// ParseConcurrencyScheduler parses a ConcurrencyScheduler component and returns a configured component and a nested circuit.
func ParseConcurrencyScheduler(
	concurrencyScheduler *policylangv1.ConcurrencyScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := concurrencyScheduler.InPorts
	if inPorts != nil {
		maxConcurrencyPort := inPorts.MaxConcurrency
		if maxConcurrencyPort != nil {
			nestedInPortsMap[maxConcurrencyPortName] = maxConcurrencyPort
		}
		passThroughPort := inPorts.PassThrough
		if passThroughPort != nil {
			nestedInPortsMap[passThroughPortName] = passThroughPort
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := concurrencyScheduler.OutPorts
	if outPorts != nil {
		acceptPercentagePort := outPorts.AcceptPercentage
		if acceptPercentagePort != nil {
			nestedOutPortsMap[acceptPercentageName] = acceptPercentagePort
		}
	}

	// Prepare parameters for prometheus queries
	policyParams := fmt.Sprintf(
		"%s=\"%s\",%s=\"%s\"",
		metrics.PolicyNameLabel,
		policyReadAPI.GetPolicyName(),
		metrics.ComponentIDLabel,
		componentID.String(),
	)

	policyParamsAccepted := fmt.Sprintf("%s,%s=\"%s\"", policyParams, metrics.DecisionTypeLabel, metrics.DecisionTypeAccepted)

	acceptedPercentageQuery := fmt.Sprintf(
		"(sum(rate(%s{%s}[30s])) / sum(rate(%s{%s}[30s]))) * 100",
		metrics.WorkloadCounterMetricName,
		policyParamsAccepted,
		metrics.WorkloadCounterMetricName,
		policyParams,
	)

	concurrencySchedulerAnyProto, err := anypb.New(
		&policyprivatev1.ConcurrencyScheduler{
			InPorts: &policylangv1.ConcurrencyLimiter_Ins{
				MaxConcurrency: &policylangv1.InPort{
					Value: &policylangv1.InPort_SignalName{
						SignalName: "MAX_CONCURRENCY",
					},
				},
				PassThrough: &policylangv1.InPort{
					Value: &policylangv1.InPort_SignalName{
						SignalName: "PASS_THROUGH",
					},
				},
			},
			Selectors:         concurrencyScheduler.GetSelectors(),
			ParentComponentId: componentID.String(),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	nestedCircuit := &policylangv1.NestedCircuit{
		InPortsMap:  nestedInPortsMap,
		OutPortsMap: nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_Query{
					Query: &policylangv1.Query{
						Component: &policylangv1.Query_Promql{
							Promql: &policylangv1.PromQL{
								OutPorts: &policylangv1.PromQL_Outs{
									Output: &policylangv1.OutPort{
										SignalName: "ACCEPT_PERCENTAGE",
									},
								},
								QueryString:        acceptedPercentageQuery,
								EvaluationInterval: durationpb.New(metrics.ScrapeInterval),
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_FlowControl{
					FlowControl: &policylangv1.FlowControl{
						Component: &policylangv1.FlowControl_Private{
							Private: concurrencySchedulerAnyProto,
						},
					},
				},
			},
		},
	}

	components.AddNestedIngress(nestedCircuit, maxConcurrencyPortName, "MAX_CONCURRENCY")
	components.AddNestedIngress(nestedCircuit, passThroughPortName, "PASS_THROUGH")
	components.AddNestedEgress(nestedCircuit, acceptPercentageName, "ACCEPT_PERCENTAGE")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("ConcurrencyScheduler",
			iface.GetSelectorsShortDescription(concurrencyScheduler.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		concurrencyScheduler,
		componentID,
		false,
	)
	if err != nil {
		return nil, nil, err
	}

	return configuredComponent, nestedCircuit, nil
}
