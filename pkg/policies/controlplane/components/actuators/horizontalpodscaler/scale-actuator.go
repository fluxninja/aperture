package horizontalpodscaler

import (
	"context"
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

// ScaleActuator struct.
type ScaleActuator struct {
	policyReadAPI      iface.Policy
	decisionWriter     *etcdwriter.Writer
	scaleActuatorProto *policylangv1.HorizontalPodScaler_ScaleActuator
	decisionsEtcdPath  string
	agentGroupName     string
	componentIndex     int
	dryRun             bool
}

// Name implements runtime.Component.
func (*ScaleActuator) Name() string { return "ScaleActuator" }

// Type implements runtime.Component.
func (*ScaleActuator) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// NewScaleActuatorAndOptions creates scale actuator and its fx options.
func NewScaleActuatorAndOptions(
	scaleActuatorProto *policylangv1.HorizontalPodScaler_ScaleActuator,
	componentIndex int,
	policyReadAPI iface.Policy,
	agentGroup string,
) (runtime.Component, fx.Option, error) {
	componentID := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), int64(componentIndex))
	decisionsEtcdPath := path.Join(paths.HorizontalPodScalerDecisionsPath, componentID)
	dryRun := false
	if scaleActuatorProto.GetDefaultConfig() != nil {
		dryRun = scaleActuatorProto.GetDefaultConfig().GetDryRun()
	}
	sa := &ScaleActuator{
		policyReadAPI:      policyReadAPI,
		agentGroupName:     agentGroup,
		componentIndex:     componentIndex,
		decisionsEtcdPath:  decisionsEtcdPath,
		scaleActuatorProto: scaleActuatorProto,
		dryRun:             dryRun,
	}

	return sa, fx.Options(
		fx.Invoke(sa.setupWriter),
	), nil
}

func (sa *ScaleActuator) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := sa.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			sa.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			sa.decisionWriter.Close()
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), sa.decisionsEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete scale decisions")
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (sa *ScaleActuator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	// Get the decision from the port
	replicasReading := inPortReadings.ReadSingleValuePort("desired_replicas")
	var replicasValue float64
	if replicasReading.Valid() {
		if replicasReading.Value() <= 0 {
			replicasValue = 0
		} else {
			replicasValue = replicasReading.Value()
		}

		return nil, sa.publishDecision(replicasValue)
	}
	sa.publishDefaultDecision()
	return nil, nil
}

// DynamicConfigUpdate finds the dynamic config and syncs the decision to agent.
func (sa *ScaleActuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := sa.policyReadAPI.GetStatusRegistry().GetLogger()
	key := sa.scaleActuatorProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		dynamicConfig := &policylangv1.HorizontalPodScaler_ScaleActuator_DynamicConfig{}
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal dynamic config")
			return
		}
		sa.setConfig(dynamicConfig)
	} else {
		sa.setConfig(sa.scaleActuatorProto.GetDefaultConfig())
	}
}

func (sa *ScaleActuator) setConfig(config *policylangv1.HorizontalPodScaler_ScaleActuator_DynamicConfig) {
	if config != nil {
		sa.dryRun = config.GetDryRun()
	} else {
		sa.dryRun = false
	}
}

func (sa *ScaleActuator) publishDefaultDecision() {
	// delete the decision
	sa.decisionWriter.Delete(sa.decisionsEtcdPath)
}

func (sa *ScaleActuator) publishDecision(desiredReplicas float64) error {
	if sa.dryRun {
		sa.publishDefaultDecision()
		return nil
	}
	logger := sa.policyReadAPI.GetStatusRegistry().GetLogger()
	// Save desired replicas in decision message
	decision := &policysyncv1.ScaleDecision{
		DesiredReplicas: int32(desiredReplicas),
	}
	// Publish decision
	logger.Autosample().Debug().Float64("desiredReplicas", desiredReplicas).Msg("Publish scale decision")
	wrapper := &policysyncv1.ScaleDecisionWrapper{
		ScaleDecision: decision,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:     sa.policyReadAPI.GetPolicyName(),
			PolicyHash:     sa.policyReadAPI.GetPolicyHash(),
			ComponentIndex: int64(sa.componentIndex),
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal policy decision")
		return err
	}
	sa.decisionWriter.Write(sa.decisionsEtcdPath, dat)
	return nil
}
