package fluxmeter

import (
	"context"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
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
func (configSync *fluxMeterConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
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
				PolicyName:    configSync.policyReadAPI.GetPolicyName(),
			}

			// Marshal the wrapper using protobuf.
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to marshal flux meter config")
				return err
			}

			// Put the marshaled data in etcd using the provided etcdPath and LeaseID.
			etcdClient.Put(configSync.etcdPath, string(dat))

			// Return nil to indicate success.
			return nil
		},

		// OnStop hook will be called when the application stops.
		OnStop: func(ctx context.Context) error {
			// Delete the data from etcd using the provided etcdPath.
			etcdClient.Delete(configSync.etcdPath)

			// Return nil to indicate success.
			return nil
		},
	})
}
