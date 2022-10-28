package concurrency

import (
	"context"
	"path"
	"sync"

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
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/rs/zerolog"
)

// LoadActuator struct.
type LoadActuator struct {
	dryRunLock        sync.RWMutex
	policyReadAPI     iface.Policy
	decisionWriter    *etcdwriter.Writer
	loadActuatorProto *policylangv1.LoadActuator
	decisionsEtcdPath string
	agentGroupName    string
	componentIndex    int
	dryRun            bool
}

// NewLoadActuatorAndOptions creates load actuator and its fx options.
func NewLoadActuatorAndOptions(
	loadActuatorProto *policylangv1.LoadActuator,
	componentIndex int,
	policyReadAPI iface.Policy,
	agentGroup string,
) (runtime.Component, fx.Option, error) {
	componentID := common.DataplaneComponentKey(agentGroup, policyReadAPI.GetPolicyName(), int64(componentIndex))
	decisionsEtcdPath := path.Join(common.LoadActuatorDecisionsPath, componentID)
	dryRun := false
	if loadActuatorProto.GetDefaultConfig() != nil {
		dryRun = loadActuatorProto.GetDefaultConfig().GetDryRun()
	}
	lsa := &LoadActuator{
		policyReadAPI:     policyReadAPI,
		agentGroupName:    agentGroup,
		componentIndex:    componentIndex,
		decisionsEtcdPath: decisionsEtcdPath,
		loadActuatorProto: loadActuatorProto,
		dryRun:            dryRun,
	}

	return lsa, fx.Options(
		fx.Invoke(lsa.setupWriter),
	), nil
}

func (la *LoadActuator) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			la.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			la.decisionWriter.Close()
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), la.decisionsEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete load decisions")
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (la *LoadActuator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
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
				return nil, la.publishDecision(lmValue, false)
			} else {
				logger.Sample(zerolog.Often).Info().Msg("Invalid load multiplier data")
			}
		} else {
			logger.Sample(zerolog.Often).Info().Msg("load_multiplier port has no reading")
		}
	} else {
		logger.Sample(zerolog.Often).Info().Msg("load_multiplier port not found")
	}
	return nil, la.publishDefaultDecision()
}

// DynamicConfigUpdate finds the dynamic config and syncs the decision to agent.
func (la *LoadActuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	la.dryRunLock.Lock()
	defer la.dryRunLock.Unlock()
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	key := la.loadActuatorProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		dynamicConfig := &policylangv1.LoadActuator_DynamicConfig{}
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal dynamic config")
			return
		}
		la.setConfig(dynamicConfig)
	} else {
		la.setConfig(la.loadActuatorProto.GetDefaultConfig())
	}
}

func (la *LoadActuator) setConfig(config *policylangv1.LoadActuator_DynamicConfig) {
	if config != nil {
		la.dryRun = config.GetDryRun()
	} else {
		la.dryRun = false
	}
}

func (la *LoadActuator) publishDefaultDecision() error {
	return la.publishDecision(1.0, true)
}

func (la *LoadActuator) publishDecision(loadMultiplier float64, passThrough bool) error {
	la.dryRunLock.RLock()
	defer la.dryRunLock.RUnlock()
	if la.dryRun {
		passThrough = true
	}
	logger := la.policyReadAPI.GetStatusRegistry().GetLogger()
	// Save load multiplier in decision message
	decision := &policysyncv1.LoadDecision{
		LoadMultiplier: loadMultiplier,
		PassThrough:    passThrough,
	}
	// Publish decision
	logger.Sample(zerolog.Often).Debug().Float64("loadMultiplier", loadMultiplier).Bool("passThrough", passThrough).Msg("Publish load decision")
	wrapper := &policysyncv1.LoadDecisionWrapper{
		LoadDecision: decision,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:     la.policyReadAPI.GetPolicyName(),
			PolicyHash:     la.policyReadAPI.GetPolicyHash(),
			ComponentIndex: int64(la.componentIndex),
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal policy decision")
		return err
	}
	la.decisionWriter.Write(la.decisionsEtcdPath, dat)
	return nil
}
