package quotascheduler

import (
	"fmt"

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
	bucketCapacityPortName = "bucket_capacity"
	fillAmountPortName     = "fill_amount"
	passThroughPortName    = "pass_through"
	acceptPercentageName   = "accept_percentage"
)

// ParseQuotaScheduler parses a QuotaScheduler component and returns a configured component and a nested circuit.
func ParseQuotaScheduler(
	quotaScheduler *policylangv1.QuotaScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := quotaScheduler.InPorts
	if inPorts != nil {
		bucketCapacityPort := inPorts.BucketCapacity
		if bucketCapacityPort != nil {
			nestedInPortsMap[bucketCapacityPortName] = bucketCapacityPort
		}
		fillAmountPort := inPorts.FillAmount
		if fillAmountPort != nil {
			nestedInPortsMap[fillAmountPortName] = fillAmountPort
		}
		passThroughPort := inPorts.PassThrough
		if passThroughPort != nil {
			nestedInPortsMap[passThroughPortName] = passThroughPort
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := quotaScheduler.OutPorts
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

	labelKey := quotaScheduler.RateLimiter.GetLimitByLabelKey()
	if labelKey == "" {
		// Deprecated: Remove in v3.0.0
		labelKey = quotaScheduler.RateLimiter.GetLabelKey()
	}

	quotaSchedulerAnyProto, err := anypb.New(
		&policyprivatev1.QuotaScheduler{
			InPorts: &policylangv1.RateLimiter_Ins{
				BucketCapacity: &policylangv1.InPort{
					Value: &policylangv1.InPort_SignalName{
						SignalName: "BUCKET_CAPACITY",
					},
				},
				FillAmount: &policylangv1.InPort{
					Value: &policylangv1.InPort_SignalName{
						SignalName: "FILL_AMOUNT",
					},
				},
				PassThrough: &policylangv1.InPort{
					Value: &policylangv1.InPort_SignalName{
						SignalName: "PASS_THROUGH",
					},
				},
			},
			Selectors: quotaScheduler.GetSelectors(),
			RateLimiter: &policylangv1.RateLimiter_Parameters{
				LabelKey:       labelKey,
				Interval:       quotaScheduler.RateLimiter.GetInterval(),
				ContinuousFill: quotaScheduler.RateLimiter.GetContinuousFill(),
				MaxIdleTime:    quotaScheduler.RateLimiter.GetMaxIdleTime(),
				LazySync: &policylangv1.RateLimiter_Parameters_LazySync{
					Enabled: quotaScheduler.RateLimiter.GetLazySync().GetEnabled(),
					NumSync: quotaScheduler.RateLimiter.GetLazySync().GetNumSync(),
				},
				DelayInitialFill: quotaScheduler.RateLimiter.GetDelayInitialFill(),
			},
			Scheduler:         quotaScheduler.GetScheduler(),
			ParentComponentId: componentID.String(),
		})
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
							Private: quotaSchedulerAnyProto,
						},
					},
				},
			},
		},
	}

	components.AddNestedIngress(nestedCircuit, bucketCapacityPortName, "BUCKET_CAPACITY")
	components.AddNestedIngress(nestedCircuit, fillAmountPortName, "FILL_AMOUNT")
	components.AddNestedIngress(nestedCircuit, passThroughPortName, "PASS_THROUGH")
	components.AddNestedEgress(nestedCircuit, acceptPercentageName, "ACCEPT_PERCENTAGE")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("QuotaScheduler",
			iface.GetSelectorsShortDescription(quotaScheduler.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		quotaScheduler,
		componentID,
		false,
	)
	if err != nil {
		return nil, nil, err
	}

	return configuredComponent, nestedCircuit, nil
}
