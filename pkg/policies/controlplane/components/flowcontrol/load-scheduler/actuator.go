package loadscheduler

import (
	"context"
	"fmt"
	"math"
	"path"
	"time"

	prometheusmodel "github.com/prometheus/common/model"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	policyprivatev1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/private/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/query/promql"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

// Actuator struct.
type Actuator struct {
	policyReadAPI            iface.Policy
	etcdClient               *etcdclient.Client
	actuatorProto            *policyprivatev1.LoadActuator
	tokensQuery              *promql.TaggedQuery
	loadSchedulerComponentID string
	cpID                     string
	etcdPaths                []string
	doActuate                bool
	ticksPerExecution        int
}

// Make sure Actuator complies with Component interface.
var _ runtime.Component = (*Actuator)(nil)

// Make sure Actuator implements background job.
var _ runtime.BackgroundJob = (*Actuator)(nil)

// Name implements runtime.Component.
func (*Actuator) Name() string { return "LoadActuator" }

// Type implements runtime.Component.
func (*Actuator) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (la *Actuator) ShortDescription() string {
	return fmt.Sprintf("%d agent groups", len(la.etcdPaths))
}

// IsActuator implements runtime.Component.
func (*Actuator) IsActuator() bool { return true }

// NewActuatorAndOptions creates load actuator and its fx options.
func NewActuatorAndOptions(
	actuatorProto *policyprivatev1.LoadActuator,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	var (
		etcdPaths []string
		options   []fx.Option
	)
	loadSchedulerComponentID := actuatorProto.LoadSchedulerComponentId

	s := actuatorProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), loadSchedulerComponentID)
		etcdPath := path.Join(paths.LoadSchedulerDecisionsPath, etcdKey)
		etcdPaths = append(etcdPaths, etcdPath)
	}

	lsa := &Actuator{
		policyReadAPI:            policyReadAPI,
		loadSchedulerComponentID: loadSchedulerComponentID,
		etcdPaths:                etcdPaths,
		actuatorProto:            actuatorProto,
		cpID:                     componentID.String(),
		ticksPerExecution:        policyReadAPI.TicksInDuration(metrics.ScrapeInterval),
	}

	// Prepare parameters for prometheus queries
	policyParams := fmt.Sprintf("%s=\"%s\",%s=\"%s\",%s=\"%s\"",
		metrics.PolicyNameLabel,
		policyReadAPI.GetPolicyName(),
		metrics.PolicyHashLabel,
		policyReadAPI.GetPolicyHash(),
		metrics.ComponentIDLabel,
		lsa.loadSchedulerComponentID,
	)
	if actuatorProto.WorkloadLatencyBasedTokens {
		tokensQuery, tokensQueryOptions, tokensQueryErr := promql.NewTaggedQueryAndOptions(
			fmt.Sprintf("sum by (%s) (increase(%s{%s}[30m])) / sum by (%s) (increase(%s{%s}[30m]))",
				metrics.WorkloadIndexLabel,
				metrics.WorkloadLatencySumMetricName,
				policyParams,
				metrics.WorkloadIndexLabel,
				metrics.WorkloadLatencyCountMetricName,
				policyParams),
			5*policyReadAPI.GetEvaluationInterval(),
			runtime.NewComponentID(loadSchedulerComponentID),
			policyReadAPI,
			"Tokens",
		)
		if tokensQueryErr != nil {
			return nil, nil, tokensQueryErr
		}
		lsa.tokensQuery = tokensQuery
		options = append(options, tokensQueryOptions)
	}

	options = append(options, fx.Invoke(lsa.setup))

	return lsa, fx.Options(options...), nil
}

func (la *Actuator) setup(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	la.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			for _, etcdPath := range la.etcdPaths {
				etcdClient.Delete(etcdPath)
			}
			return nil
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (la *Actuator) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	circuitAPI.ScheduleConditionalBackgroundJob(la, la.ticksPerExecution)
	retErr := func(err error) (runtime.PortToReading, error) {
		var errMulti error
		pErr := la.publishDefaultDecision(circuitAPI.GetTickInfo())
		if pErr != nil {
			errMulti = multierr.Append(err, pErr)
		}
		return nil, errMulti
	}

	tokensByWorkload := make(map[string]float64)
	if la.tokensQuery != nil {
		taggedResult, err := la.tokensQuery.ExecuteTaggedQuery(circuitAPI)
		if err != nil {
			if err != promql.ErrNoQueriesReturned {
				return retErr(err)
			}
		}
		promValue := taggedResult.Value
		if promValue != nil {
			vector, ok := promValue.(prometheusmodel.Vector)
			if !ok {
				err = fmt.Errorf("tokens query returned a non-vector value")
				return retErr(err)
			}
			for _, sample := range vector {
				for k, v := range sample.Metric {
					if k == metrics.WorkloadIndexLabel {
						// if sample.Value is NaN, continue
						if math.IsNaN(float64(sample.Value)) {
							continue
						}
						workloadIndex := string(v)
						sampleValue := float64(sample.Value)
						tokensByWorkload[workloadIndex] = sampleValue
						break
					}
				}
			}
		}
	}

	var lm float64
	var pt bool
	lmValue := inPortReadings.ReadSingleReadingPort("load_multiplier")
	if !lmValue.Valid() {
		pt = true
	} else {
		lm = lmValue.Value()
		if lm <= 0 {
			lm = 0
		}
	}

	return nil, la.publishDecision(circuitAPI.GetTickInfo(), lm, pt, tokensByWorkload)
}

// DynamicConfigUpdate implements runtime.Component.DynamicConfigUpdate.
func (la *Actuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (la *Actuator) publishDefaultDecision(tickInfo runtime.TickInfo) error {
	return la.publishDecision(tickInfo, 1.0, true, nil)
}

func (la *Actuator) publishDecision(tickInfo runtime.TickInfo, loadMultiplier float64, passThrough bool, tokensByWorkload map[string]float64) error {
	if !la.doActuate {
		return nil
	}
	la.doActuate = false
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	// validUntil = time.Now() + tickInfo.Interval() * la.ticksPerExecution
	validUntil := tickInfo.Timestamp().Add(tickInfo.Interval() * time.Duration(la.ticksPerExecution*5))
	// Save load multiplier in decision message
	decision := &policysyncv1.LoadDecision{
		LoadMultiplier:        loadMultiplier,
		PassThrough:           passThrough,
		TickInfo:              tickInfo.Serialize(),
		TokensByWorkloadIndex: tokensByWorkload,
		ValidUntil:            timestamppb.New(validUntil),
	}
	// Publish decision
	logger.Autosample().Debug().Float64("loadMultiplier", loadMultiplier).Bool("passThrough", passThrough).Msg("Publish load decision")
	wrapper := &policysyncv1.LoadDecisionWrapper{
		LoadDecision: decision,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:  la.policyReadAPI.GetPolicyName(),
			PolicyHash:  la.policyReadAPI.GetPolicyHash(),
			ComponentId: la.loadSchedulerComponentID,
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal policy decision")
		return err
	}
	for _, etcdPath := range la.etcdPaths {
		la.etcdClient.Put(etcdPath, string(dat))
	}

	return nil
}

// GetJob implements runtime.BackgroundJob.GetJob.
func (la *Actuator) GetJob() jobs.Job {
	return jobs.NewNoOpJob(la.cpID)
}

// NotifyCompletion implements runtime.BackgroundJob.NotifyCompletion.
func (la *Actuator) NotifyCompletion() {
	la.doActuate = true
}
