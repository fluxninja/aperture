package flowregulator

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

type flowRegulatorSync struct {
	policyReadAPI         iface.Policy
	flowRegulatorProto    *policylangv1.FlowRegulator
	decision              *policysyncv1.FlowRegulatorDecision
	decisionWriter        *etcdwriter.Writer
	dynamicConfigWriter   *etcdwriter.Writer
	configEtcdPath        string
	decisionsEtcdPath     string
	dynamicConfigEtcdPath string
	componentID           string
}

// Name implements runtime.Component.
func (*flowRegulatorSync) Name() string { return "FlowRegulator" }

// Type implements runtime.Component.
func (*flowRegulatorSync) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (regulatorSync *flowRegulatorSync) ShortDescription() string {
	return iface.GetServiceShortDescription(regulatorSync.flowRegulatorProto.Parameters.FlowSelector.ServiceSelector)
}

// IsActuator implements runtime.Component.
func (*flowRegulatorSync) IsActuator() bool { return true }

// NewFlowRegulatorAndOptions creates fx options for FlowRegulator and also returns agent group name associated with it.
func NewFlowRegulatorAndOptions(
	flowRegulatorProto *policylangv1.FlowRegulator,
	componentID string,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	agentGroup := flowRegulatorProto.Parameters.FlowSelector.ServiceSelector.AgentGroup
	etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID)
	configEtcdPath := path.Join(paths.FlowRegulatorConfigPath, etcdKey)
	decisionsEtcdPath := path.Join(paths.FlowRegulatorDecisionsPath, etcdKey)
	dynamicConfigEtcdPath := path.Join(paths.FlowRegulatorDynamicConfigPath, etcdKey)

	regulatorSync := &flowRegulatorSync{
		flowRegulatorProto:    flowRegulatorProto,
		decision:              &policysyncv1.FlowRegulatorDecision{},
		policyReadAPI:         policyReadAPI,
		configEtcdPath:        configEtcdPath,
		decisionsEtcdPath:     decisionsEtcdPath,
		dynamicConfigEtcdPath: dynamicConfigEtcdPath,
		componentID:           componentID,
	}
	return regulatorSync, fx.Options(
		fx.Invoke(
			regulatorSync.setupSync,
		),
	), nil
}

func (regulatorSync *flowRegulatorSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := regulatorSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.FlowRegulatorWrapper{
				FlowRegulator: regulatorSync.flowRegulatorProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  regulatorSync.policyReadAPI.GetPolicyName(),
					PolicyHash:  regulatorSync.policyReadAPI.GetPolicyHash(),
					ComponentId: regulatorSync.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal flow regulator config")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				regulatorSync.configEtcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				logger.Error().Err(err).Msg("failed to put flow regulator config")
				return err
			}
			regulatorSync.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			regulatorSync.dynamicConfigWriter = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			regulatorSync.dynamicConfigWriter.Close()
			regulatorSync.decisionWriter.Close()
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), regulatorSync.configEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete flow regulator config")
				merr = multierr.Append(merr, err)
			}
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), regulatorSync.decisionsEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete flow regulator decisions")
				merr = multierr.Append(merr, err)
			}
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), regulatorSync.dynamicConfigEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete flow regulator dynamic config")
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (regulatorSync *flowRegulatorSync) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	acceptPercentage, ok := inPortReadings["accept_percentage"]
	if !ok {
		return nil, nil
	}

	if len(acceptPercentage) == 0 {
		return nil, nil
	}

	acceptPercentageReading := inPortReadings.ReadSingleReadingPort("accept_percentage")
	var acceptPercentageValue float64
	if !acceptPercentageReading.Valid() {
		acceptPercentageValue = 100 // default to 100%
	} else {
		acceptPercentageValue = acceptPercentageReading.Value()
	}
	return nil, regulatorSync.publishAcceptPercentage(acceptPercentageValue)
}

func (regulatorSync *flowRegulatorSync) publishAcceptPercentage(acceptPercentageValue float64) error {
	logger := regulatorSync.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish decision
	logger.Debug().Float64("flux", acceptPercentageValue).Msg("publishing flux regulator decision")
	wrapper := &policysyncv1.FlowRegulatorDecisionWrapper{
		FlowRegulatorDecision: &policysyncv1.FlowRegulatorDecision{
			AcceptPercentage: acceptPercentageValue,
		},
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:  regulatorSync.policyReadAPI.GetPolicyName(),
			PolicyHash:  regulatorSync.policyReadAPI.GetPolicyHash(),
			ComponentId: regulatorSync.componentID,
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal flux regulator decision")
		return err
	}
	if regulatorSync.decisionWriter == nil {
		logger.Panic().Msg("decision writer is nil")
	}
	regulatorSync.decisionWriter.Write(regulatorSync.decisionsEtcdPath, dat)

	return nil
}

// DynamicConfigUpdate handles overrides.
func (regulatorSync *flowRegulatorSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := regulatorSync.policyReadAPI.GetStatusRegistry().GetLogger()
	publishDynamicConfig := func(dynamicConfig *policylangv1.FlowRegulator_DynamicConfig) {
		wrapper := &policysyncv1.FlowRegulatorDynamicConfigWrapper{
			FlowRegulatorDynamicConfig: dynamicConfig,
			CommonAttributes: &policysyncv1.CommonAttributes{
				PolicyName:  regulatorSync.policyReadAPI.GetPolicyName(),
				PolicyHash:  regulatorSync.policyReadAPI.GetPolicyHash(),
				ComponentId: regulatorSync.componentID,
			},
		}
		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("failed to marshal dynamic config")
			return
		}
		if regulatorSync.dynamicConfigWriter == nil {
			logger.Panic().Msg("dynamic config writer is nil")
		}
		regulatorSync.dynamicConfigWriter.Write(regulatorSync.dynamicConfigEtcdPath, dat)
		logger.Info().Msg("Flow Regulator dynamic config updated")
	}
	dynamicConfig := &policylangv1.FlowRegulator_DynamicConfig{}
	key := regulatorSync.flowRegulatorProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal dynamic config")
			return
		}
		publishDynamicConfig(dynamicConfig)
	} else {
		publishDynamicConfig(regulatorSync.flowRegulatorProto.GetDefaultConfig())
	}
}
