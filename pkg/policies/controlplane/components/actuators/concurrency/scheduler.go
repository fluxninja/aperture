package concurrency

import (
	"context"
	"fmt"
	"math"
	"path"
	"time"

	prometheusmodel "github.com/prometheus/common/model"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/promql"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/policies/paths"
)

var (
	concurrencyQueryInterval = time.Second * 1
	tokensQueryInterval      = time.Second * 10
)

// Scheduler is part of the concurrency control component stack.
type Scheduler struct {
	policyReadAPI iface.Policy
	// Prometheus query for accepted concurrency
	acceptedQuery *promql.ScalarQuery
	// Prometheus query for incoming concurrency
	incomingQuery *promql.ScalarQuery
	// Prometheus query for tokens based on ms latency
	tokensQuery *promql.TaggedQuery

	// saves tokens value per workload read from prometheus
	tokensByWorkload *policysyncv1.TokensDecision
	writer           *etcdwriter.Writer
	agentGroupName   string
	etcdPath         string
	componentIndex   int
}

// Name implements runtime.Component.
func (*Scheduler) Name() string { return "Scheduler" }

// Type implements runtime.Component.
func (*Scheduler) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// NewSchedulerAndOptions creates scheduler and its fx options.
func NewSchedulerAndOptions(
	schedulerProto *policylangv1.Scheduler,
	componentIndex int,
	policyReadAPI iface.Policy,
	agentGroupName string,
) (runtime.Component, fx.Option, error) {
	etcdPath := path.Join(paths.AutoTokenResultsPath,
		paths.AgentComponentKey(agentGroupName, policyReadAPI.GetPolicyName(), int64(componentIndex)))

	scheduler := &Scheduler{
		policyReadAPI: policyReadAPI,
		tokensByWorkload: &policysyncv1.TokensDecision{
			TokensByWorkloadIndex: make(map[string]uint64),
		},
		agentGroupName: agentGroupName,
		componentIndex: componentIndex,
		etcdPath:       etcdPath,
	}

	// Prepare parameters for prometheus queries
	policyParams := fmt.Sprintf("%s=\"%s\",%s=\"%s\",%s=\"%d\"",
		metrics.PolicyNameLabel,
		policyReadAPI.GetPolicyName(),
		metrics.PolicyHashLabel,
		policyReadAPI.GetPolicyHash(),
		metrics.ComponentIndexLabel,
		componentIndex,
	)

	acceptedQuery, acceptedQueryOptions, acceptedQueryErr := promql.NewScalarQueryAndOptions(
		fmt.Sprintf("sum(rate(%s{%s}[10s]))",
			metrics.AcceptedConcurrencyMetricName,
			policyParams),
		concurrencyQueryInterval,
		componentIndex,
		policyReadAPI,
		"AcceptedConcurrency",
	)
	if acceptedQueryErr != nil {
		return nil, fx.Options(), acceptedQueryErr
	}
	scheduler.acceptedQuery = acceptedQuery

	incomingQuery, incomingQueryOptions, incomingQueryErr := promql.NewScalarQueryAndOptions(
		fmt.Sprintf("sum(rate(%s{%s}[10s]))",
			metrics.IncomingConcurrencyMetricName,
			policyParams),
		concurrencyQueryInterval,
		componentIndex,
		policyReadAPI,
		"IncomingConcurrency",
	)
	if incomingQueryErr != nil {
		return nil, nil, incomingQueryErr
	}
	scheduler.incomingQuery = incomingQuery

	// add decision_type filter to the params
	autoTokensPolicyParams := policyParams + ",decision_type!=\"DECISION_TYPE_REJECTED\""
	if schedulerProto.AutoTokens {
		tokensQuery, tokensQueryOptions, tokensQueryErr := promql.NewTaggedQueryAndOptions(
			fmt.Sprintf("sum by (%s) (increase(%s{%s}[30m])) / sum by (%s) (increase(%s{%s}[30m]))",
				metrics.WorkloadIndexLabel,
				metrics.WorkloadLatencySumMetricName,
				autoTokensPolicyParams,
				metrics.WorkloadIndexLabel,
				metrics.WorkloadLatencyCountMetricName,
				autoTokensPolicyParams),
			tokensQueryInterval,
			componentIndex,
			policyReadAPI,
			"Tokens",
		)
		if tokensQueryErr != nil {
			return nil, nil, tokensQueryErr
		}
		scheduler.tokensQuery = tokensQuery
		return scheduler,
			fx.Options(
				acceptedQueryOptions,
				incomingQueryOptions,
				tokensQueryOptions,
				fx.Invoke(scheduler.setupWriter),
			),
			nil
	} else {
		return scheduler,
			fx.Options(
				acceptedQueryOptions,
				incomingQueryOptions,
				fx.Invoke(scheduler.setupWriter),
			),
			nil
	}
}

func (s *Scheduler) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := s.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s.writer = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), s.etcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete tokens decision config")
				return err
			}
			s.writer.Close()
			return nil
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (s *Scheduler) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	logger := s.policyReadAPI.GetStatusRegistry().GetLogger()
	var errMulti error

	if s.tokensQuery != nil {
		taggedResult, err := s.tokensQuery.ExecuteTaggedQuery(tickInfo)
		promValue := taggedResult.Value
		if err != nil {
			if err != promql.ErrNoQueriesReturned {
				logger.Error().Err(err).Msg("could not read tokens query from prometheus")
				errMulti = multierr.Append(errMulti, err)
			}
		} else if promValue != nil {
			if vector, ok := promValue.(prometheusmodel.Vector); ok {
				tokensDecision := &policysyncv1.TokensDecision{
					TokensByWorkloadIndex: make(map[string]uint64),
				}
				for _, sample := range vector {
					for k, v := range sample.Metric {
						if k == metrics.WorkloadIndexLabel {
							// if sample.Value is NaN, continue
							if math.IsNaN(float64(sample.Value)) {
								continue
							}
							workloadIndex := string(v)
							sampleValue := uint64(sample.Value)
							tokensDecision.TokensByWorkloadIndex[workloadIndex] = sampleValue
							break
						}
					}
				}
				err = s.publishQueryTokens(tokensDecision)
				if err != nil {
					errMulti = multierr.Append(errMulti, err)
					logger.Error().Err(err).Msg("failed to publish tokens")
				}
			} else {
				err = fmt.Errorf("tokens query returned a non-vector value")
				errMulti = multierr.Append(errMulti, err)
				logger.Error().Err(err).Msg("Failed to parse tokens")
			}
		}
	}

	var acceptedReading, incomingReading runtime.Reading

	outPortReadings := make(runtime.PortToValue)

	acceptedScalarResult, err := s.acceptedQuery.ExecuteScalarQuery(tickInfo)
	acceptedValue := acceptedScalarResult.Value
	if err != nil {
		acceptedReading = runtime.InvalidReading()
		if err != promql.ErrNoQueriesReturned {
			errMulti = multierr.Append(errMulti, err)
		}
	} else {
		acceptedReading = runtime.NewReading(acceptedValue)
	}
	outPortReadings["accepted_concurrency"] = []runtime.Reading{acceptedReading}

	incomingScalarResult, err := s.incomingQuery.ExecuteScalarQuery(tickInfo)
	incomingValue := incomingScalarResult.Value
	if err != nil {
		incomingReading = runtime.InvalidReading()
		if err != promql.ErrNoQueriesReturned {
			errMulti = multierr.Append(errMulti, err)
		}
	} else {
		incomingReading = runtime.NewReading(incomingValue)
	}
	outPortReadings["incoming_concurrency"] = []runtime.Reading{incomingReading}

	return outPortReadings, errMulti
}

// DynamicConfigUpdate is a no-op for this component.
func (s *Scheduler) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}

func (s *Scheduler) publishQueryTokens(tokens *policysyncv1.TokensDecision) error {
	logger := s.policyReadAPI.GetStatusRegistry().GetLogger()
	s.tokensByWorkload = tokens
	policyName := s.policyReadAPI.GetPolicyName()
	policyHash := s.policyReadAPI.GetPolicyHash()

	wrapper := &policysyncv1.TokensDecisionWrapper{
		TokensDecision: tokens,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:     policyName,
			PolicyHash:     policyHash,
			ComponentIndex: int64(s.componentIndex),
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal tokens")
		return err
	}
	s.writer.Write(s.etcdPath, dat)
	return nil
}
