package concurrency

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"time"

	prometheusmodel "github.com/prometheus/common/model"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policydecisionsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

var (
	concurrencyQueryInterval = time.Second * 1
	tokensQueryInterval      = time.Second * 10
)

// Scheduler is part of the concurrency control component stack.
type Scheduler struct {
	policyReadAPI iface.Policy
	// saves promValue result from tokens query to check if anything changed
	tokensPromValue prometheusmodel.Value
	// Prometheus query for accepted concurrency
	acceptedQuery *components.ScalarQuery
	// Prometheus query for incoming concurrency
	incomingQuery *components.ScalarQuery
	// Prometheus query for tokens based on ms latency
	tokensQuery *components.TaggedQuery

	// saves tokens value per workload read from prometheus
	tokensByWorkload *policydecisionsv1.TokensDecision
	writer           *etcdwriter.Writer
	agentGroupName   string
	etcdPath         string
	componentIndex   int
}

// NewSchedulerAndOptions creates scheduler and its fx options.
func NewSchedulerAndOptions(
	schedulerProto *policylangv1.Scheduler,
	componentIndex int,
	policyReadAPI iface.Policy,
	agentGroupName string,
) (runtime.Component, fx.Option, error) {
	etcdPath := path.Join(common.AutoTokenResultsPath,
		common.DataplaneComponentKey(agentGroupName, policyReadAPI.GetPolicyName(), int64(componentIndex)))

	scheduler := &Scheduler{
		policyReadAPI: policyReadAPI,
		tokensByWorkload: &policydecisionsv1.TokensDecision{
			TokensByWorkloadIndex: make(map[string]uint64),
		},
		agentGroupName:  agentGroupName,
		componentIndex:  componentIndex,
		tokensPromValue: nil,
		etcdPath:        etcdPath,
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

	acceptedQuery, acceptedQueryOptions, acceptedQueryErr := components.NewScalarQueryAndOptions(
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

	incomingQuery, incomingQueryOptions, incomingQueryErr := components.NewScalarQueryAndOptions(
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
		tokensQuery, tokensQueryOptions, tokensQueryErr := components.NewTaggedQueryAndOptions(
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
		promValue, err := s.tokensQuery.ExecutePromQuery(tickInfo)
		if err != nil {
			if err != components.ErrNoQueriesReturned {
				logger.Error().Err(err).Msg("could not read tokens query from prometheus")
				errMulti = multierr.Append(errMulti, err)
			}
		} else if promValue != nil && !reflect.DeepEqual(promValue, s.tokensPromValue) {
			// update only if something changed
			s.tokensPromValue = promValue

			if vector, ok := promValue.(prometheusmodel.Vector); ok {
				tokensDecision := &policydecisionsv1.TokensDecision{
					TokensByWorkloadIndex: make(map[string]uint64),
				}
				for _, sample := range vector {
					for k, v := range sample.Metric {
						if k == metrics.WorkloadIndexLabel {
							workloadIndex := string(v)
							sampleValue := uint64(sample.Value)
							tokensDecision.TokensByWorkloadIndex[workloadIndex] = sampleValue
							break
						}
					}
				}
				err = s.publishQueryTokens(tokensDecision)
				if err != nil {
					logger.Error().Err(err).Msg("failed to publish tokens")
				}
			} else {
				err = fmt.Errorf("tokens query returned a non-vector value")
				logger.Error().Err(err).Msg("Failed to parse tokens")
			}
		}
	}

	var acceptedReading, incomingReading runtime.Reading

	outPortReadings := make(runtime.PortToValue)

	acceptedValue, err := s.acceptedQuery.ExecuteScalarQuery(tickInfo)
	if err != nil {
		acceptedReading = runtime.InvalidReading()
		if err != components.ErrNoQueriesReturned {
			errMulti = multierr.Append(errMulti, err)
		}
	} else {
		acceptedReading = runtime.NewReading(acceptedValue)
	}
	outPortReadings["accepted_concurrency"] = []runtime.Reading{acceptedReading}

	incomingValue, err := s.incomingQuery.ExecuteScalarQuery(tickInfo)
	if err != nil {
		incomingReading = runtime.InvalidReading()
		if err != components.ErrNoQueriesReturned {
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

func (s *Scheduler) publishQueryTokens(tokens *policydecisionsv1.TokensDecision) error {
	logger := s.policyReadAPI.GetStatusRegistry().GetLogger()
	// TODO: publish only on change
	s.tokensByWorkload = tokens
	policyName := s.policyReadAPI.GetPolicyName()
	policyHash := s.policyReadAPI.GetPolicyHash()

	wrapper := &wrappersv1.TokensDecisionWrapper{
		TokensDecision: tokens,
		CommonAttributes: &wrappersv1.CommonAttributes{
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
