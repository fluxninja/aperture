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
	"github.com/fluxninja/aperture/pkg/policies/apis/policyapi"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/component"
	"github.com/fluxninja/aperture/pkg/utils"
)

type fluxMeterConfigSync struct {
	policyBaseAPI  policyapi.PolicyBaseAPI
	fluxMeterProto *policylangv1.FluxMeter
	etcdPath       string
	agentGroupName string
}

// NewFluxMeterOptions creates fx options for FluxMeter.
func NewFluxMeterOptions(
	fluxMeterProto *policylangv1.FluxMeter,
	policyBaseAPI policyapi.PolicyBaseAPI,
	metricSubRegistry policyapi.MetricSubRegistry,
) (fx.Option, error) {
	// Get Agent Group Name from FluxMeter.Selector.AgentGroup
	selectorProto := fluxMeterProto.GetSelector()
	if selectorProto == nil {
		return nil, errors.New("FluxMeter.Selector is nil")
	}
	agentGroupName := selectorProto.GetAgentGroup()

	wrapperProto := &configv1.ConfigPropertiesWrapper{
		AgentGroupName: agentGroupName,
		PolicyName:     policyBaseAPI.GetPolicyName(),
		PolicyHash:     policyBaseAPI.GetPolicyHash(),
	}

	// Register FluxMeter
	err := registerFluxMeter(fluxMeterProto, wrapperProto, metricSubRegistry)
	if err != nil {
		return nil, err
	}

	etcdPath := path.Join(paths.FluxMeterConfigPath,
		paths.IdentifierForFluxMeter(agentGroupName, policyBaseAPI.GetPolicyName(), fluxMeterProto.GetName()))
	configSync := &fluxMeterConfigSync{
		fluxMeterProto: fluxMeterProto,
		policyBaseAPI:  policyBaseAPI,
		agentGroupName: agentGroupName,
		etcdPath:       etcdPath,
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
			wrapper, err := utils.WrapWithConfProps(
				configSync.fluxMeterProto,
				configSync.agentGroupName,
				configSync.policyBaseAPI.GetPolicyName(),
				configSync.policyBaseAPI.GetPolicyHash(),
				0,
			)
			if err != nil {
				log.Error().Err(err).Msg("Failed to wrap flux meter config in config properties")
				return err
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

// TODO (hasit): rename fluxmeter metric name to a static one 'flux_meter'

// registerFluxMeter registers histograms for fluxmeter in controller.
func registerFluxMeter(fluxMeterProto *policylangv1.FluxMeter, componentAPI component.ComponentAPI, metricSubRegistry policyapi.MetricSubRegistry) error {
	// Original metric name
	metricName := fluxMeterProto.Name

	metricID := paths.MetricIDForFluxMeter(componentAPI, metricName)
	matcher, err := labels.NewMatcher(labels.MatchEqual, "metric_id", metricID)
	if err != nil {
		return err
	}
	metricLabels := []*labels.Matcher{matcher}
	metricSubRegistry.RegisterHistogramSub(metricName, metricName, metricLabels)
	return nil
}
