package concurrency

import (
	"context"
	"fmt"
	"path"
	"time"

	policydecisionsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	prometheusmodel "github.com/prometheus/common/model"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
)

var demandQueryInterval = time.Second * 10

// ConcurrencyActuator is a component that computes the concurrency multiplier and concurrency demand decisions.
type ConcurrencyActuator struct {
	// demandResult from demand query
	demandResult          components.TaggedResult
	policyReadAPI         iface.Policy
	concurrencyResultTick runtime.TickInfo
	// A map from Agent ID to it's concurrency decision
	concurrencyDemandDecisions    map[string]*policydecisionsv1.ConcurrencyDemandDecision
	concurrencyMultiplierDecision *policydecisionsv1.ConcurrencyMultiplierDecision
	// Prometheus query for demand concurrency
	demandQuery              *components.TaggedQuery
	writer                   *etcdwriter.Writer
	agentGroup               string
	multiplierEtcdPath       string
	demandDefaultEtcdPath    string
	concurrencyDemandTotal   float64
	concurrencyDemandAverage float64
	componentIndex           int
}

// NewConcurrencyActuatorAndOptions creates concurrency actuator and its fx options.
func NewConcurrencyActuatorAndOptions(
	_ *policylangv1.ConcurrencyActuator,
	componentIndex int,
	policyReadAPI iface.Policy,
	agentGroup string,
) (runtime.Component, fx.Option, error) {
	multiplierEtcdPath := path.Join(common.ConcurrencyMultiplierDecisionsPath,
		common.DataplaneComponentKey(agentGroup, policyReadAPI.GetPolicyName(), int64(componentIndex)))
	demandDefaultEtcdPath := path.Join(common.ConcurrencyDemandDefaultDecisionsPath,
		common.DataplaneComponentKey(agentGroup, policyReadAPI.GetPolicyName(), int64(componentIndex)))
	ca := &ConcurrencyActuator{
		policyReadAPI:              policyReadAPI,
		componentIndex:             componentIndex,
		agentGroup:                 agentGroup,
		multiplierEtcdPath:         multiplierEtcdPath,
		demandDefaultEtcdPath:      demandDefaultEtcdPath,
		concurrencyDemandDecisions: make(map[string]*policydecisionsv1.ConcurrencyDemandDecision),
	}
	ca.concurrencyMultiplierDecision = &policydecisionsv1.ConcurrencyMultiplierDecision{}

	// Prepare parameters for prometheus queries
	policyParams := fmt.Sprintf("%s=\"%s\",%s=\"%s\",%s=\"%d\"",
		metrics.PolicyNameLabel,
		policyReadAPI.GetPolicyName(),
		metrics.PolicyHashLabel,
		policyReadAPI.GetPolicyHash(),
		metrics.ComponentIndexLabel,
		componentIndex,
	)

	demandQuery, demandQueryOptions, demandQueryErr := components.NewTaggedQueryAndOptions(
		fmt.Sprintf("sum by (%s) (rate(%s{%s}[1m]))",
			metrics.ProcessUUIDLabel,
			metrics.IncomingConcurrencyMetricName,
			policyParams),
		demandQueryInterval,
		componentIndex,
		policyReadAPI,
		"DemandConcurrency",
	)
	if demandQueryErr != nil {
		return nil, nil, demandQueryErr
	}
	ca.demandQuery = demandQuery

	return ca, fx.Options(
		fx.Invoke(ca.setupWriter),
		demandQueryOptions,
	), nil
}

func (ca *ConcurrencyActuator) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := ca.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ca.writer = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), ca.multiplierEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete concurrency multiplier decision config")
				return err
			}
			ca.writer.Close()
			return nil
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (ca *ConcurrencyActuator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	log.Info().Msg("Executing Concurrency Actuator")
	logger := ca.policyReadAPI.GetStatusRegistry().GetLogger()
	var errMulti error

	demandResult, err := ca.demandQuery.ExecuteTaggedQuery(tickInfo)
	if err != nil {
		if err != components.ErrNoQueriesReturned {
			logger.Error().Err(err).Msg("could not read demand query from prometheus")
			errMulti = multierr.Append(errMulti, err)
		}
	} else if (ca.demandResult == components.TaggedResult{}) ||
		(ca.demandResult.TickInfo.Tick() != demandResult.TickInfo.Tick()) {
		// Save the new demand result
		ca.demandResult = demandResult
		promValue := demandResult.Value
		if promValue == nil {
			logger.Error().Msg("demand query returned nil value")
			errMulti = multierr.Append(errMulti, fmt.Errorf("demand query returned nil value"))
		} else {
			if vector, ok := promValue.(prometheusmodel.Vector); !ok {
				err = fmt.Errorf("demand query returned a non-vector value")
				logger.Error().Err(err).Msg("Failed to parse concurrency demand query")
				errMulti = multierr.Append(errMulti, err)
			} else {
				concurrencyDemandDecisions := make(map[string]*policydecisionsv1.ConcurrencyDemandDecision)
				concurrencyDemandTotal := 0.0
				for _, sample := range vector {
					for k, v := range sample.Metric {
						if k == metrics.ProcessUUIDLabel {
							processUUID := string(v)
							sampleValue := float64(sample.Value)
							concurrencyDemandDecisions[processUUID] = &policydecisionsv1.ConcurrencyDemandDecision{
								Tick:   int64(demandResult.TickInfo.Tick()),
								Demand: sampleValue,
							}
							concurrencyDemandTotal += sampleValue
							break
						}
					}
				}
				// Save the new concurrencyDenandTotal
				ca.concurrencyDemandTotal = concurrencyDemandTotal
				if concurrencyDemandTotal > 0 {
					ca.concurrencyDemandAverage = concurrencyDemandTotal / float64(len(concurrencyDemandDecisions))
				} else {
					ca.concurrencyDemandAverage = 0
				}
				// Save the tick of demand result
				ca.concurrencyResultTick = demandResult.TickInfo
				concurrencyDemandDefaultDecision := &policydecisionsv1.ConcurrencyDemandDecision{
					Tick:   int64(demandResult.TickInfo.Tick()),
					Demand: ca.concurrencyDemandAverage,
				}
				// Publish the new concurrencyDemandDecisions
				err = ca.publishDemandDecisions(concurrencyDemandDecisions, concurrencyDemandDefaultDecision)
				if err != nil {
					logger.Error().Err(err).Msg("failed to publish demand decisions")
					errMulti = multierr.Append(errMulti, err)
				}
			}
		}
	}

	// Get the decision from the port
	desiredConcurrency, ok := inPortReadings["desired_concurrency"]
	if ok {
		if len(desiredConcurrency) > 0 {
			desiredConcurrencyReading := desiredConcurrency[0]
			if !desiredConcurrencyReading.Valid() {
				// remove the multiplier decision
				err := ca.deleteMultiplier()
				if err != nil {
					errMulti = multierr.Append(errMulti, err)
				}
			} else {
				multiplier := ca.computeMultiplier(desiredConcurrencyReading.Value())
				err := ca.publishMultiplier(multiplier)
				if err != nil {
					errMulti = multierr.Append(errMulti, err)
				}
			}
		}
	}
	return nil, errMulti
}

// DynamicConfigUpdate is a no-op for this component.
func (ca *ConcurrencyActuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (ca *ConcurrencyActuator) computeMultiplier(desiredConcurrency float64) float64 {
	if ca.concurrencyDemandTotal == 0 {
		return 1
	}
	return desiredConcurrency / ca.concurrencyDemandTotal
}

func (ca *ConcurrencyActuator) publishMultiplier(multiplier float64) error {
	logger := ca.policyReadAPI.GetStatusRegistry().GetLogger()
	logger.Debug().Msgf("Publishing concurrency multiplier: %f", multiplier)
	policyName := ca.policyReadAPI.GetPolicyName()
	policyHash := ca.policyReadAPI.GetPolicyHash()
	wrapper := &wrappersv1.ConcurrencyMultiplierDecisionWrapper{
		CommonAttributes: &wrappersv1.CommonAttributes{
			PolicyName:     policyName,
			PolicyHash:     policyHash,
			ComponentIndex: int64(ca.componentIndex),
		},
		ConcurrencyMultiplierDecision: &policydecisionsv1.ConcurrencyMultiplierDecision{
			Tick:       int64(ca.concurrencyResultTick.Tick()),
			Multiplier: multiplier,
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal multiplier decision")
		return err
	}
	// Write the multiplier to etcd
	ca.writer.Write(ca.multiplierEtcdPath, dat)
	return nil
}

func (ca *ConcurrencyActuator) deleteMultiplier() error {
	logger := ca.policyReadAPI.GetStatusRegistry().GetLogger()
	logger.Debug().Msg("Deleting concurrency multiplier decision")
	// Delete the multiplier from etcd
	ca.writer.Delete(ca.multiplierEtcdPath)
	return nil
}

func (ca *ConcurrencyActuator) publishDemandDecisions(concurrencyDemandDecisions map[string]*policydecisionsv1.ConcurrencyDemandDecision, concurrencyDemandDefaultDecision *policydecisionsv1.ConcurrencyDemandDecision) error {
	logger := ca.policyReadAPI.GetStatusRegistry().GetLogger()
	logger.Debug().Msg("Publishing concurrency demand decisions")
	policyName := ca.policyReadAPI.GetPolicyName()
	policyHash := ca.policyReadAPI.GetPolicyHash()
	publishDemandDecision := func(etcdPath string, concurrencyDemandDecision *policydecisionsv1.ConcurrencyDemandDecision) error {
		wrapper := &wrappersv1.ConcurrencyDemandDecisionWrapper{
			CommonAttributes: &wrappersv1.CommonAttributes{
				PolicyName:     policyName,
				PolicyHash:     policyHash,
				ComponentIndex: int64(ca.componentIndex),
			},
			ConcurrencyDemandDecision: concurrencyDemandDecision,
		}
		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal demand decision")
			return err
		}
		// Write the demand decision to etcd
		ca.writer.Write(etcdPath, dat)
		return nil
	}

	// Publish for each processUUID
	for processUUID, concurrencyDemandDecision := range concurrencyDemandDecisions {
		err := publishDemandDecision(ca.demandEtcdPath(processUUID), concurrencyDemandDecision)
		if err != nil {
			return err
		}
	}
	// Delete any demand decisions that are no longer needed
	for processUUID := range ca.concurrencyDemandDecisions {
		if _, ok := concurrencyDemandDecisions[processUUID]; !ok {
			ca.writer.Delete(ca.demandEtcdPath(processUUID))
		}
	}

	// Publish the default demand decision
	return publishDemandDecision(ca.demandDefaultEtcdPath, concurrencyDemandDefaultDecision)
}

func (ca *ConcurrencyActuator) demandEtcdPath(processUUID string) string {
	return path.Join(common.ConcurrencyDemandDecisionsPath,
		common.DataplaneComponentAgentKey(processUUID, ca.agentGroup, ca.policyReadAPI.GetPolicyName(), int64(ca.componentIndex)))
}
