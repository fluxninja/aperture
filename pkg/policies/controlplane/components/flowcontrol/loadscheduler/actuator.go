package loadscheduler

import (
	"context"
	"fmt"
	"math"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/private/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/query/promql"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	prometheusmodel "github.com/prometheus/common/model"
)

// Actuator struct.
type Actuator struct {
	policyReadAPI            iface.Policy
	decisionWriter           *etcdwriter.Writer
	actuatorProto            *policyprivatev1.LoadActuator
	tokensQuery              *promql.TaggedQuery
	loadSchedulerComponentID string
	etcdPaths                []string
	dryRun                   bool
}

// Name implements runtime.Component.
func (*Actuator) Name() string { return "Actuator" }

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
	_ runtime.ComponentID,
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

	dryRun := false
	if actuatorProto.GetDefaultConfig() != nil {
		dryRun = actuatorProto.GetDefaultConfig().GetDryRun()
	}

	lsa := &Actuator{
		policyReadAPI:            policyReadAPI,
		loadSchedulerComponentID: loadSchedulerComponentID,
		etcdPaths:                etcdPaths,
		actuatorProto:            actuatorProto,
		dryRun:                   dryRun,
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
			10*policyReadAPI.GetEvaluationInterval(),
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

	options = append(options, fx.Invoke(lsa.setupWriter))

	return lsa, fx.Options(options...), nil
}

func (la *Actuator) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			la.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			la.decisionWriter.Close()
			for _, etcdPath := range la.etcdPaths {
				_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), etcdPath)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to delete load decisions")
					merr = multierr.Append(merr, err)
				}
			}
			return merr
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (la *Actuator) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	retErr := func(err error) (runtime.PortToReading, error) {
		var errMulti error
		pErr := la.publishDefaultDecision(tickInfo)
		if pErr != nil {
			errMulti = multierr.Append(err, pErr)
		}
		return nil, errMulti
	}

	tokensByWorkload := make(map[string]uint64)

	if la.tokensQuery != nil {
		taggedResult, err := la.tokensQuery.ExecuteTaggedQuery(tickInfo)
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
						sampleValue := uint64(sample.Value)
						tokensByWorkload[workloadIndex] = sampleValue
						break
					}
				}
			}
		}
	}

	// Get the decision from the port
	lm, ok := inPortReadings["load_multiplier"]
	if !ok {
		return retErr(fmt.Errorf("load_multiplier port not found"))
	}
	if len(lm) == 0 {
		return retErr(fmt.Errorf("load_multiplier port has no reading"))
	}
	lmReading := lm[0]
	var lmValue float64
	if !lmReading.Valid() {
		return retErr(fmt.Errorf("invalid load_multiplier reading"))
	}
	if lmReading.Value() <= 0 {
		lmValue = 0
	} else {
		lmValue = lmReading.Value()
	}
	return nil, la.publishDecision(tickInfo, lmValue, false, tokensByWorkload)
}

// DynamicConfigUpdate finds the dynamic config and syncs the decision to agent.
func (la *Actuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	key := la.actuatorProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		dynamicConfig := &policylangv1.LoadScheduler_DynamicConfig{}
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal dynamic config")
			return
		}
		la.setConfig(dynamicConfig)
	} else {
		la.setConfig(la.actuatorProto.GetDefaultConfig())
	}
}

func (la *Actuator) setConfig(config *policylangv1.LoadScheduler_DynamicConfig) {
	if config != nil {
		la.dryRun = config.GetDryRun()
	} else {
		la.dryRun = false
	}
}

func (la *Actuator) publishDefaultDecision(tickInfo runtime.TickInfo) error {
	return la.publishDecision(tickInfo, 1.0, true, nil)
}

func (la *Actuator) publishDecision(tickInfo runtime.TickInfo, loadMultiplier float64, passThrough bool, tokensByWorkload map[string]uint64) error {
	if la.dryRun {
		passThrough = true
	}
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	// Save load multiplier in decision message
	decision := &policysyncv1.LoadDecision{
		LoadMultiplier:        loadMultiplier,
		PassThrough:           passThrough,
		TickInfo:              tickInfo.Serialize(),
		TokensByWorkloadIndex: tokensByWorkload,
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
		la.decisionWriter.Write(etcdPath, dat)
	}

	return nil
}
