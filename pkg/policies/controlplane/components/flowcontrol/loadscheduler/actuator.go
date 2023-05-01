package loadscheduler

import (
	"context"
	"fmt"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/policies/paths"
)

// Actuator struct.
type Actuator struct {
	policyReadAPI  iface.Policy
	decisionWriter *etcdwriter.Writer
	actuatorProto  *policylangv1.LoadScheduler_Actuator
	componentID    string
	etcdPaths      []string
	dryRun         bool
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
	actuatorProto *policylangv1.LoadScheduler_Actuator,
	componentID string,
	policyReadAPI iface.Policy,
	agentGroups []string,
) (runtime.Component, fx.Option, error) {
	var etcdPaths []string
	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID)
		etcdPath := path.Join(paths.LoadSchedulerDecisionsPath, etcdKey)
		etcdPaths = append(etcdPaths, etcdPath)
	}

	dryRun := false
	if actuatorProto.GetDefaultConfig() != nil {
		dryRun = actuatorProto.GetDefaultConfig().GetDryRun()
	}

	lsa := &Actuator{
		policyReadAPI: policyReadAPI,
		componentID:   componentID,
		etcdPaths:     etcdPaths,
		actuatorProto: actuatorProto,
		dryRun:        dryRun,
	}

	return lsa, fx.Options(
		fx.Invoke(lsa.setupWriter),
	), nil
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
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	// Get the decision from the port
	lm, ok := inPortReadings["load_multiplier"]
	if ok {
		if len(lm) > 0 {
			lmReading := lm[0]
			var lmValue float64
			if lmReading.Valid() {
				if lmReading.Value() <= 0 {
					lmValue = 0
				} else {
					lmValue = lmReading.Value()
				}

				return nil, la.publishDecision(tickInfo, lmValue, false)
			} else {
				logger.Autosample().Info().Msg("Invalid load multiplier data")
			}
		} else {
			logger.Autosample().Info().Msg("load_multiplier port has no reading")
		}
	} else {
		logger.Autosample().Info().Msg("load_multiplier port not found")
	}
	return nil, la.publishDefaultDecision(tickInfo)
}

// DynamicConfigUpdate finds the dynamic config and syncs the decision to agent.
func (la *Actuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	key := la.actuatorProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		dynamicConfig := &policylangv1.LoadScheduler_Actuator_DynamicConfig{}
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal dynamic config")
			return
		}
		la.setConfig(dynamicConfig)
	} else {
		la.setConfig(la.actuatorProto.GetDefaultConfig())
	}
}

func (la *Actuator) setConfig(config *policylangv1.LoadScheduler_Actuator_DynamicConfig) {
	if config != nil {
		la.dryRun = config.GetDryRun()
	} else {
		la.dryRun = false
	}
}

func (la *Actuator) publishDefaultDecision(tickInfo runtime.TickInfo) error {
	return la.publishDecision(tickInfo, 1.0, true)
}

func (la *Actuator) publishDecision(tickInfo runtime.TickInfo, loadMultiplier float64, passThrough bool) error {
	if la.dryRun {
		passThrough = true
	}
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	// Save load multiplier in decision message
	decision := &policysyncv1.LoadDecision{
		LoadMultiplier: loadMultiplier,
		PassThrough:    passThrough,
		TickInfo:       tickInfo.Serialize(),
	}
	// Publish decision
	logger.Autosample().Debug().Float64("loadMultiplier", loadMultiplier).Bool("passThrough", passThrough).Msg("Publish load decision")
	wrapper := &policysyncv1.LoadDecisionWrapper{
		LoadDecision: decision,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:  la.policyReadAPI.GetPolicyName(),
			PolicyHash:  la.policyReadAPI.GetPolicyHash(),
			ComponentId: la.componentID,
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
