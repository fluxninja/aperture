package regulator

import (
	"context"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/v2/pkg/etcd/writer"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

type regulatorSync struct {
	policyReadAPI          iface.Policy
	RegulatorProto         *policylangv1.Regulator
	decision               *policysyncv1.RegulatorDecision
	decisionWriter         *etcdwriter.Writer
	componentID            string
	configEtcdPaths        []string
	decisionEtcdPaths      []string
	passThroughLabelValues []string
}

// Name implements runtime.Component.
func (*regulatorSync) Name() string { return "Regulator" }

// Type implements runtime.Component.
func (*regulatorSync) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (regulatorSync *regulatorSync) ShortDescription() string {
	return iface.GetSelectorsShortDescription(regulatorSync.RegulatorProto.Parameters.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*regulatorSync) IsActuator() bool { return true }

// NewRegulatorAndOptions creates fx options for Regulator and also returns agent group name associated with it.
func NewRegulatorAndOptions(
	regulatorProto *policylangv1.Regulator,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := regulatorProto.Parameters.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths, decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.RegulatorConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
		decisionEtcdPath := path.Join(paths.RegulatorDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	regulatorSync := &regulatorSync{
		RegulatorProto:    regulatorProto,
		decision:          &policysyncv1.RegulatorDecision{},
		policyReadAPI:     policyReadAPI,
		configEtcdPaths:   configEtcdPaths,
		decisionEtcdPaths: decisionEtcdPaths,
		componentID:       componentID.String(),
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
			var merr error
			for _, configEtcdPath := range regulatorSync.configEtcdPaths {
				_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
					configEtcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
				if err != nil {
					logger.Error().Err(err).Msg("failed to put regulator config")
					merr = multierr.Append(merr, err)
				}
			}
			regulatorSync.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			return merr
		},
		OnStop: func(ctx context.Context) error {
			regulatorSync.decisionWriter.Close()
			deleteEtcdPath := func(paths []string) error {
				var merr error
				for _, path := range paths {
					_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), path)
					if err != nil {
						logger.Error().Err(err).Msgf("failed to delete etcd path %s", path)
						merr = multierr.Append(merr, err)
					}
				}
				return merr
			}

			merr := deleteEtcdPath(regulatorSync.configEtcdPaths)
			merr = multierr.Append(merr, deleteEtcdPath(regulatorSync.decisionEtcdPaths))
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
			AcceptPercentage:       acceptPercentageValue,
			PassThroughLabelValues: regulatorSync.passThroughLabelValues,
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
	for _, decisionEtcdPath := range regulatorSync.decisionEtcdPaths {
		regulatorSync.decisionWriter.Write(decisionEtcdPath, dat)
	}

	return nil
}

// DynamicConfigUpdate handles overrides.
func (regulatorSync *regulatorSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := regulatorSync.policyReadAPI.GetStatusRegistry().GetLogger()
	key := regulatorSync.RegulatorProto.GetPassThroughLabelValuesConfigKey()
	passThroughLabelValues := []string{}
	if !unmarshaller.IsSet(key) {
		regulatorSync.passThroughLabelValues = regulatorSync.RegulatorProto.GetPassThroughLabelValues()
		return
	}
	err := unmarshaller.UnmarshalKey(key, passThroughLabelValues)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal dynamic config")
		regulatorSync.passThroughLabelValues = regulatorSync.RegulatorProto.GetPassThroughLabelValues()
		return
	}
	regulatorSync.passThroughLabelValues = passThroughLabelValues
}
