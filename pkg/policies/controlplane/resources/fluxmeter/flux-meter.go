package fluxmeter

import (
	"context"
	"errors"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
)

type fluxMeterConfigSync struct {
	policyReadAPI  iface.Policy
	fluxMeterProto *policylangv1.FluxMeter
	etcdPath       string
	agentGroupName string
	fluxMeterName  string
}

// NewFluxMeterOptions creates fx options for FluxMeter.
func NewFluxMeterOptions(
	name string,
	fluxMeterProto *policylangv1.FluxMeter,
	policyBaseAPI iface.Policy,
) (fx.Option, error) {
	// Get Agent Group Name from FluxMeter.Selector.AgentGroup
	selectorProto := fluxMeterProto.GetSelector()
	if selectorProto == nil {
		return nil, errors.New("FluxMeter.Selector is nil")
	}
	agentGroup := selectorProto.ServiceSelector.GetAgentGroup()

	etcdPath := path.Join(paths.FluxMeterConfigPath,
		paths.FluxMeterKey(agentGroup, name))
	configSync := &fluxMeterConfigSync{
		fluxMeterProto: fluxMeterProto,
		policyReadAPI:  policyBaseAPI,
		agentGroupName: agentGroup,
		etcdPath:       etcdPath,
		fluxMeterName:  name,
	}

	return fx.Options(
		fx.Invoke(
			configSync.doSync,
		),
	), nil
}

func (configSync *fluxMeterConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.FluxMeterWrapper{
				FluxMeterName: configSync.fluxMeterName,
				FluxMeter:     configSync.fluxMeterProto,
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal flux meter config")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				configSync.etcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				logger.Error().Err(err).Msg("Failed to put flux meter config")
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), configSync.etcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete flux meter config")
				return err
			}
			return nil
		},
	})

	return nil
}
