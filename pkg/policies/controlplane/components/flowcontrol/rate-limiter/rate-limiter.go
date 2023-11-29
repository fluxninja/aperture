package ratelimiter

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

// ParseRateLimiter parses a RateLimiter component and returns a configured component and a nested circuit.
func ParseRateLimiter(
	rateLimiter *policylangv1.RateLimiter,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := rateLimiter.InPorts
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
	outPorts := rateLimiter.OutPorts
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
		metrics.RateLimiterCounterTotalMetricName,
		policyParamsAccepted,
		metrics.RateLimiterCounterTotalMetricName,
		policyParams,
	)

	labelKey := rateLimiter.Parameters.GetLimitByLabelKey()
	if labelKey == "" {
		// Deprecated: Remove in v3.0.0
		labelKey = rateLimiter.Parameters.GetLabelKey()
	}

	rateLimiterAnyProto, err := anypb.New(
		&policyprivatev1.RateLimiter{
			InPorts: &policyprivatev1.RateLimiter_Ins{
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
			Selectors: rateLimiter.GetSelectors(),
			Parameters: &policyprivatev1.RateLimiter_Parameters{
				LimitByLabelKey: labelKey,
				Interval:        rateLimiter.Parameters.GetInterval(),
				ContinuousFill:  rateLimiter.Parameters.GetContinuousFill(),
				MaxIdleTime:     rateLimiter.Parameters.GetMaxIdleTime(),
				LazySync: &policyprivatev1.RateLimiter_Parameters_LazySync{
					Enabled: rateLimiter.Parameters.GetLazySync().GetEnabled(),
					NumSync: rateLimiter.Parameters.GetLazySync().GetNumSync(),
				},
				DelayInitialFill: rateLimiter.Parameters.GetDelayInitialFill(),
			},
			RequestParameters: &policyprivatev1.RateLimiter_RequestParameters{
				TokensLabelKey:           rateLimiter.RequestParameters.GetTokensLabelKey(),
				DeniedResponseStatusCode: rateLimiter.RequestParameters.GetDeniedResponseStatusCode(),
			},
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
							Private: rateLimiterAnyProto,
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
		runtime.NewDummyComponent("RateLimiter",
			iface.GetSelectorsShortDescription(rateLimiter.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		rateLimiter,
		componentID,
		false,
	)
	if err != nil {
		return nil, nil, err
	}

	return configuredComponent, nestedCircuit, nil
}
