package fluxmeter

import (
	"context"
	"errors"
	"path"

	"github.com/prometheus/prometheus/model/labels"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
)

type fluxMeterConfigSync struct {
	policyBaseAPI  iface.PolicyBase
	fluxMeterProto *policylangv1.FluxMeter
	etcdPath       string
	agentGroupName string
	fluxmeterName  string
}

// NewFluxMeterOptions creates fx options for FluxMeter.
func NewFluxMeterOptions(
	name string,
	fluxMeterProto *policylangv1.FluxMeter,
	policyBaseAPI iface.PolicyBase,
	metricSubRegistry iface.MetricSubRegistry,
) (fx.Option, error) {
	// Get Agent Group Name from FluxMeter.Selector.AgentGroup
	selectorProto := fluxMeterProto.GetSelector()
	if selectorProto == nil {
		return nil, errors.New("FluxMeter.Selector is nil")
	}
	agentGroup := selectorProto.GetAgentGroup()
	wrapperProto := &configv1.FluxMeterWrapper{
		FluxmeterName: name,
		PolicyName:    policyBaseAPI.GetPolicyName(),
		PolicyHash:    policyBaseAPI.GetPolicyHash(),
	}

	// Register FluxMeter
	err := registerFluxMeter(fluxMeterProto, wrapperProto, metricSubRegistry)
	if err != nil {
		return nil, err
	}

	etcdPath := path.Join(paths.FluxMeterConfigPath,
		paths.FluxMeterKey(agentGroup, policyBaseAPI.GetPolicyName(), name))
	configSync := &fluxMeterConfigSync{
		fluxMeterProto: fluxMeterProto,
		policyBaseAPI:  policyBaseAPI,
		agentGroupName: agentGroup,
		etcdPath:       etcdPath,
		fluxmeterName:  name,
	}

	return fx.Options(
		fx.Invoke(
			configSync.doSync,
		),
	), nil
}

func (configSync *fluxMeterConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &configv1.FluxMeterWrapper{
				FluxmeterName: configSync.fluxmeterName,
				FluxMeter:     configSync.fluxMeterProto,
				PolicyName:    configSync.policyBaseAPI.GetPolicyName(),
				PolicyHash:    configSync.policyBaseAPI.GetPolicyHash(),
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal flux meter config")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				configSync.etcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				log.Error().Err(err).Msg("Failed to put flux meter config")
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), configSync.etcdPath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete flux meter config")
				return err
			}
			return nil
		},
	})

	return nil
}

// registerFluxMeter registers histograms for fluxmeter in controller.
func registerFluxMeter(
	fluxMeterProto *policylangv1.FluxMeter,
	fluxMeterWrapper *configv1.FluxMeterWrapper,
	metricSubRegistry iface.MetricSubRegistry,
) error {
	policyNameMatcher, err := labels.NewMatcher(labels.MatchEqual, "policy_name", fluxMeterWrapper.GetPolicyName())
	if err != nil {
		return err
	}
	fluxMeterNameMatcher, err := labels.NewMatcher(labels.MatchEqual, "flux_meter_name", fluxMeterWrapper.GetFluxmeterName())
	if err != nil {
		return err
	}
	policyHashMatcher, err := labels.NewMatcher(labels.MatchEqual, "policy_hash", fluxMeterWrapper.GetPolicyHash())
	if err != nil {
		return err
	}
	metricLabels := []*labels.Matcher{policyNameMatcher, fluxMeterNameMatcher, policyHashMatcher}
	metricSubRegistry.RegisterHistogramSub(fluxMeterWrapper.GetFluxmeterName(), "flux_meter", metricLabels)
	return nil
}
