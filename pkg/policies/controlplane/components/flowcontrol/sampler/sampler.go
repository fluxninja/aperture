package sampler

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

type samplerSync struct {
	policyReadAPI          iface.Policy
	SamplerProto           *policylangv1.Sampler
	decision               *policysyncv1.SamplerDecision
	decisionWriter         *etcdwriter.Writer
	componentID            string
	configEtcdPaths        []string
	decisionEtcdPaths      []string
	passThroughLabelValues []string
}

// Name implements runtime.Component.
func (*samplerSync) Name() string { return "Sampler" }

// Type implements runtime.Component.
func (*samplerSync) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (samplerSync *samplerSync) ShortDescription() string {
	return iface.GetSelectorsShortDescription(samplerSync.SamplerProto.Parameters.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*samplerSync) IsActuator() bool { return true }

// NewSamplerAndOptions creates fx options for Sampler and also returns agent group name associated with it.
func NewSamplerAndOptions(
	samplerProto *policylangv1.Sampler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := samplerProto.Parameters.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths, decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.SamplerConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
		decisionEtcdPath := path.Join(paths.SamplerDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	samplerSync := &samplerSync{
		SamplerProto:      samplerProto,
		decision:          &policysyncv1.SamplerDecision{},
		policyReadAPI:     policyReadAPI,
		configEtcdPaths:   configEtcdPaths,
		decisionEtcdPaths: decisionEtcdPaths,
		componentID:       componentID.String(),
	}
	return samplerSync, fx.Options(
		fx.Invoke(
			samplerSync.setupSync,
		),
	), nil
}

func (samplerSync *samplerSync) setupSync(scopedKV *etcdclient.SessionScopedKV, lifecycle fx.Lifecycle) error {
	logger := samplerSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.SamplerWrapper{
				Sampler: samplerSync.SamplerProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  samplerSync.policyReadAPI.GetPolicyName(),
					PolicyHash:  samplerSync.policyReadAPI.GetPolicyHash(),
					ComponentId: samplerSync.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal  sampler config")
				return err
			}
			var merr error
			for _, configEtcdPath := range samplerSync.configEtcdPaths {
				_, err = scopedKV.Put(clientv3.WithRequireLeader(ctx), configEtcdPath, string(dat))
				if err != nil {
					logger.Error().Err(err).Msg("failed to put sampler config")
					merr = multierr.Append(merr, err)
				}
			}
			samplerSync.decisionWriter = etcdwriter.NewWriter(&scopedKV.KVWrapper)
			return merr
		},
		OnStop: func(ctx context.Context) error {
			samplerSync.decisionWriter.Close()
			deleteEtcdPath := func(paths []string) error {
				var merr error
				for _, path := range paths {
					_, err := scopedKV.Delete(clientv3.WithRequireLeader(ctx), path)
					if err != nil {
						logger.Error().Err(err).Msgf("failed to delete etcd path %s", path)
						merr = multierr.Append(merr, err)
					}
				}
				return merr
			}

			merr := deleteEtcdPath(samplerSync.configEtcdPaths)
			merr = multierr.Append(merr, deleteEtcdPath(samplerSync.decisionEtcdPaths))
			return merr
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (samplerSync *samplerSync) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
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
	return nil, samplerSync.publishAcceptPercentage(acceptPercentageValue)
}

func (samplerSync *samplerSync) publishAcceptPercentage(acceptPercentageValue float64) error {
	logger := samplerSync.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish decision
	logger.Debug().Float64("flux", acceptPercentageValue).Msg("publishing flux sampler decision")
	wrapper := &policysyncv1.SamplerDecisionWrapper{
		SamplerDecision: &policysyncv1.SamplerDecision{
			AcceptPercentage:       acceptPercentageValue,
			PassThroughLabelValues: samplerSync.passThroughLabelValues,
		},
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:  samplerSync.policyReadAPI.GetPolicyName(),
			PolicyHash:  samplerSync.policyReadAPI.GetPolicyHash(),
			ComponentId: samplerSync.componentID,
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal flux sampler decision")
		return err
	}
	if samplerSync.decisionWriter == nil {
		logger.Panic().Msg("decision writer is nil")
	}
	for _, decisionEtcdPath := range samplerSync.decisionEtcdPaths {
		samplerSync.decisionWriter.Write(decisionEtcdPath, dat)
	}

	return nil
}

// DynamicConfigUpdate handles overrides.
func (samplerSync *samplerSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := samplerSync.policyReadAPI.GetStatusRegistry().GetLogger()
	key := samplerSync.SamplerProto.GetPassThroughLabelValuesConfigKey()
	passThroughLabelValues := []string{}
	if !unmarshaller.IsSet(key) {
		samplerSync.passThroughLabelValues = samplerSync.SamplerProto.GetPassThroughLabelValues()
		return
	}
	err := unmarshaller.UnmarshalKey(key, passThroughLabelValues)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal dynamic config")
		samplerSync.passThroughLabelValues = samplerSync.SamplerProto.GetPassThroughLabelValues()
		return
	}
	samplerSync.passThroughLabelValues = passThroughLabelValues
}
