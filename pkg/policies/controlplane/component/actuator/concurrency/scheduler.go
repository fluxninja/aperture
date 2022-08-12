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
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/apis/policyapi"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/reading"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/utils"
)

var (
	concurrencyQueryInterval = time.Second * 1
	tokensQueryInterval      = time.Minute * 5
)

const (
	workloadIndexLabel = "workload_index"
)

// Scheduler is part of the concurrency control component stack.
type Scheduler struct {
	policyReadAPI policyapi.PolicyReadAPI
	// saves promValue result from tokens query to check if anything changed
	tokensPromValue prometheusmodel.Value
	// Prometheus query for accepted concurrency
	acceptedQuery *component.ScalarQuery
	// Prometheus query for incoming concurrency
	incomingQuery *component.ScalarQuery
	// Prometheus query for tokens based on ms latency
	tokensQuery *component.TaggedQuery

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
	policyReadAPI policyapi.PolicyReadAPI,
	agentGroupName string,
) (runtime.Component, fx.Option, error) {
	etcdPath := path.Join(paths.AutoTokenResultsPath,
		paths.IdentifierForComponent(agentGroupName, policyReadAPI.GetPolicyName(), int64(componentIndex)))

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

	acceptedQuery, acceptedQueryOptions, acceptedQueryErr := component.NewScalarQueryAndOptions(
		fmt.Sprintf("sum(rate(accepted_concurrency{policy_name=\"%s\",policy_hash=\"%s\",component_index=\"%d\"}[10s]))", policyReadAPI.GetPolicyName(), policyReadAPI.GetPolicyHash(), componentIndex),
		concurrencyQueryInterval,
		componentIndex,
		policyReadAPI,
		"AcceptedConcurrency",
	)
	if acceptedQueryErr != nil {
		return nil, fx.Options(), acceptedQueryErr
	}
	scheduler.acceptedQuery = acceptedQuery

	incomingQuery, incomingQueryOptions, incomingQueryErr := component.NewScalarQueryAndOptions(
		fmt.Sprintf("sum(rate(incoming_concurrency{policy_name=\"%s\",policy_hash=\"%s\",component_index=\"%d\"}[10s]))", policyReadAPI.GetPolicyName(), policyReadAPI.GetPolicyHash(), componentIndex),
		concurrencyQueryInterval,
		componentIndex,
		policyReadAPI,
		"IncomingConcurrency",
	)
	if incomingQueryErr != nil {
		return nil, nil, incomingQueryErr
	}
	scheduler.incomingQuery = incomingQuery

	if schedulerProto.AutoTokens {
		tokensQuery, tokensQueryOptions, tokensQueryErr := component.NewTaggedQueryAndOptions(
			fmt.Sprintf("sum by %s (increase(workload_latency_ms_sum{policy_name=\"%s\",policy_hash=\"%s\",component_index=\"%d\"}[30m])) / sum by %s (increase(workload_latency_ms_count{policy_name=\"%s\",policy_hash=\"%s\",component_index=\"%d\"}[30m]))",
				workloadIndexLabel, policyReadAPI.GetPolicyName(), policyReadAPI.GetPolicyHash(), componentIndex,
				workloadIndexLabel, policyReadAPI.GetPolicyName(), policyReadAPI.GetPolicyHash(), componentIndex),
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
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s.writer = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), s.etcdPath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete tokens decision config")
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
	var errMulti error

	if s.tokensQuery != nil {
		promValue, err := s.tokensQuery.ExecutePromQuery(tickInfo)
		if err != nil {
			log.Error().Err(err).Msg("could not read tokens query from prometheus")
			errMulti = multierr.Append(errMulti, err)
		} else if promValue != nil && !reflect.DeepEqual(promValue, s.tokensPromValue) {
			// update only if something changed
			s.tokensPromValue = promValue

			if vector, ok := promValue.(prometheusmodel.Vector); ok {
				tokensDecision := &policydecisionsv1.TokensDecision{
					TokensByWorkloadIndex: make(map[string]uint64),
				}
				for _, sample := range vector {
					for k, v := range sample.Metric {
						if k == workloadIndexLabel {
							workloadIndex := string(v)
							sampleValue := uint64(sample.Value)
							tokensDecision.TokensByWorkloadIndex[workloadIndex] = sampleValue
							break
						}
					}
				}
				err = s.publishQueryTokens(tokensDecision)
				if err != nil {
					log.Error().Err(err).Msg("failed to publish tokens")
				}
			} else {
				err = fmt.Errorf("tokens query returned a non-vector value")
				log.Error().Err(err).Msg("Failed to parse tokens")
			}
		}
	}

	var acceptedReading, incomingReading reading.Reading

	outPortReadings := make(runtime.PortToValue)

	acceptedValue, err := s.acceptedQuery.ExecuteScalarQuery(tickInfo)
	if err != nil {
		acceptedReading = reading.NewInvalid()
		errMulti = multierr.Append(errMulti, err)
	} else {
		acceptedReading = reading.New(acceptedValue)
	}
	outPortReadings["accepted_concurrency"] = []reading.Reading{acceptedReading}

	incomingValue, err := s.incomingQuery.ExecuteScalarQuery(tickInfo)
	if err != nil {
		incomingReading = reading.NewInvalid()
		errMulti = multierr.Append(errMulti, err)
	} else {
		incomingReading = reading.New(incomingValue)
	}
	outPortReadings["incoming_concurrency"] = []reading.Reading{incomingReading}

	return outPortReadings, errMulti
}

func (s *Scheduler) publishQueryTokens(tokens *policydecisionsv1.TokensDecision) error {
	s.tokensByWorkload = tokens
	policyName := s.policyReadAPI.GetPolicyName()
	policyHash := s.policyReadAPI.GetPolicyHash()

	wrapper, err := utils.WrapWithConfProps(tokens, s.agentGroupName, policyName, policyHash, s.componentIndex)
	if err != nil {
		log.Error().Err(err).Msg("Failed to wrap tokens in config properties")
		return err
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal tokens")
		return err
	}
	s.writer.Write(s.etcdPath, dat)
	return nil
}
