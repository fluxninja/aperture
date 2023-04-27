package regulator

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

type regulatorSync struct {
	policyReadAPI         iface.Policy
	RegulatorProto        *policylangv1.Regulator
	decision              *policysyncv1.RegulatorDecision
	decisionWriter        *etcdwriter.Writer
	dynamicConfigWriter   *etcdwriter.Writer
	configEtcdPath        string
	decisionsEtcdPath     string
	dynamicConfigEtcdPath string
	componentID           string
}

// Name implements runtime.Component.
func (*regulatorSync) Name() string { return "Regulator" }

// Type implements runtime.Component.
func (*regulatorSync) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (regulatorSync *regulatorSync) ShortDescription() string {
	return iface.GetServiceShortDescription(regulatorSync.RegulatorProto.Parameters.FlowSelector.ServiceSelector)
}

// IsActuator implements runtime.Component.
func (*regulatorSync) IsActuator() bool { return true }

// NewRegulatorAndOptions creates fx options for Regulator and also returns agent group name associated with it.
func NewRegulatorAndOptions(
	RegulatorProto *policylangv1.Regulator,
	componentID string,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	agentGroup := RegulatorProto.Parameters.FlowSelector.ServiceSelector.AgentGroup
	etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID)
	configEtcdPath := path.Join(paths.RegulatorConfigPath, etcdKey)
	decisionsEtcdPath := path.Join(paths.RegulatorDecisionsPath, etcdKey)
	dynamicConfigEtcdPath := path.Join(paths.RegulatorDynamicConfigPath, etcdKey)

	regulatorSync := &regulatorSync{
		RegulatorProto:        RegulatorProto,
		decision:              &policysyncv1.RegulatorDecision{},
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

func (regulatorSync *regulatorSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := regulatorSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.RegulatorWrapper{
				Regulator: regulatorSync.RegulatorProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  regulatorSync.policyReadAPI.GetPolicyName(),
					PolicyHash:  regulatorSync.policyReadAPI.GetPolicyHash(),
					ComponentId: regulatorSync.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal  regulator config")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				regulatorSync.configEtcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				logger.Error().Err(err).Msg("failed to put  regulator config")
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
				logger.Error().Err(err).Msg("failed to delete  regulator config")
				merr = multierr.Append(merr, err)
			}
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), regulatorSync.decisionsEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete  regulator decisions")
				merr = multierr.Append(merr, err)
			}
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), regulatorSync.dynamicConfigEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete  regulator dynamic config")
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (regulatorSync *regulatorSync) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
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

func (regulatorSync *regulatorSync) publishAcceptPercentage(acceptPercentageValue float64) error {
	logger := regulatorSync.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish decision
	logger.Debug().Float64("flux", acceptPercentageValue).Msg("publishing flux regulator decision")
	wrapper := &policysyncv1.RegulatorDecisionWrapper{
		RegulatorDecision: &policysyncv1.RegulatorDecision{
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
func (regulatorSync *regulatorSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := regulatorSync.policyReadAPI.GetStatusRegistry().GetLogger()
	publishDynamicConfig := func(dynamicConfig *policylangv1.Regulator_DynamicConfig) {
		wrapper := &policysyncv1.RegulatorDynamicConfigWrapper{
			RegulatorDynamicConfig: dynamicConfig,
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
		logger.Info().Msg("regulator dynamic config updated")
	}
	dynamicConfig := &policylangv1.Regulator_DynamicConfig{}
	key := regulatorSync.RegulatorProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal dynamic config")
			return
		}
		publishDynamicConfig(dynamicConfig)
	} else {
		publishDynamicConfig(regulatorSync.RegulatorProto.GetDefaultConfig())
	}
}
