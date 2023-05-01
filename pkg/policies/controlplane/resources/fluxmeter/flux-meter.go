package fluxmeter

import (
	"context"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
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
	// Deprecated 1.8.0
	flowSelectorProto := fluxMeterProto.GetFlowSelector()
	if flowSelectorProto != nil {
		selector := &policylangv1.Selector{
			ControlPoint: flowSelectorProto.FlowMatcher.ControlPoint,
			LabelMatcher: flowSelectorProto.FlowMatcher.LabelMatcher,
			Service:      flowSelectorProto.ServiceSelector.Service,
			AgentGroup:   flowSelectorProto.ServiceSelector.AgentGroup,
		}
		fluxMeterProto.Selectors = append(fluxMeterProto.Selectors, selector)
		fluxMeterProto.FlowSelector = nil
	}

	var options []fx.Option

	s := fluxMeterProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	for _, agentGroup := range agentGroups {
		etcdPath := path.Join(paths.FluxMeterConfigPath,
			paths.FluxMeterKey(agentGroup, name))
		configSync := &fluxMeterConfigSync{
			fluxMeterProto: fluxMeterProto,
			policyReadAPI:  policyBaseAPI,
			agentGroupName: agentGroup,
			etcdPath:       etcdPath,
			fluxMeterName:  name,
		}
		options = append(options, fx.Invoke(configSync.doSync))
	}

	return fx.Options(options...), nil
}

// doSync is a method of fluxMeterConfigSync struct that syncs the flux meter configuration to etcd.
// It takes an etcdClient of type etcdclient.Client and a lifecycle of type fx.Lifecycle as input parameters.
// It returns an error in case of any failure during the sync process.
func (configSync *fluxMeterConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	// Get the logger instance from the status registry.
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()

	// Append fx.Hook to the lifecycle.
	lifecycle.Append(fx.Hook{
		// OnStart hook will be called when the application starts.
		OnStart: func(ctx context.Context) error {
			// Create a FluxMeterWrapper using the fluxMeterName and fluxMeterProto provided.
			wrapper := &policysyncv1.FluxMeterWrapper{
				FluxMeterName: configSync.fluxMeterName,
				FluxMeter:     configSync.fluxMeterProto,
			}

			// Marshal the wrapper using protobuf.
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to marshal flux meter config")
				return err
			}

			// Put the marshaled data in etcd using the provided etcdPath and LeaseID.
			// It returns an error in case of any failure.
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				configSync.etcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to put flux meter config")
				return err
			}

			// Return nil to indicate success.
			return nil
		},

		// OnStop hook will be called when the application stops.
		OnStop: func(ctx context.Context) error {
			// Delete the data from etcd using the provided etcdPath.
			// It returns an error in case of any failure.
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), configSync.etcdPath)
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to delete flux meter config")
				return err
			}

			// Return nil to indicate success.
			return nil
		},
	})

	// Return nil to indicate success.
	return nil
}
