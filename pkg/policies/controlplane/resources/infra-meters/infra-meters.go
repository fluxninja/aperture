package inframeters

import (
	"context"
	"encoding/json"
	"path"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
)

type infraMeterConfigSync struct {
	name            string
	policyReadAPI   iface.Policy
	infraMeterProto *policylangv1.InfraMeter
	etcdPath        string
	agentGroupName  string
}

// NewInfraMetersOptions creates fx options InfraMeters.
func NewInfraMetersOptions(
	infraMeters map[string]*policylangv1.InfraMeter,
	policyBaseAPI iface.Policy,
) (fx.Option, error) {
	var options []fx.Option

	for name, infraMeter := range infraMeters {
		agentGroup := infraMeter.GetAgentGroup()
		etcdPath := path.Join(paths.InfraMeterConfigPath,
			paths.InfraMeterKey(agentGroup, policyBaseAPI.GetPolicyName(), name))
		configSync := &infraMeterConfigSync{
			name:            name,
			infraMeterProto: infraMeter,
			policyReadAPI:   policyBaseAPI,
			agentGroupName:  agentGroup,
			etcdPath:        etcdPath,
		}
		options = append(options, fx.Invoke(configSync.doSync))
	}

	return fx.Options(options...), nil
}

func (configSync *infraMeterConfigSync) doSync(scopedKV *etcdclient.SessionScopedKV, lifecycle fx.Lifecycle) {
	// Get the logger instance from the status registry.
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()

	// Append fx.Hook to the lifecycle.
	lifecycle.Append(fx.Hook{
		// OnStart hook will be called when the application starts.
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.InfraMeterWrapper{
				InfraMeter:     configSync.infraMeterProto,
				InfraMeterName: configSync.name,
				PolicyName:     configSync.policyReadAPI.GetPolicyName(),
			}

			// Marshal the infra meter using json marshaler.
			dat, err := json.Marshal(wrapper)
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to marshal infra meter config")
				return err
			}

			// Put the marshaled data in etcd using the provided etcdPath and LeaseID.
			// It returns an error in case of any failure.
			_, err = scopedKV.Put(clientv3.WithRequireLeader(ctx), configSync.etcdPath, string(dat))
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to put infra meter config")
				return err
			}

			// Return nil to indicate success.
			return nil
		},

		// OnStop hook will be called when the application stops.
		OnStop: func(ctx context.Context) error {
			// Delete the data from etcd using the provided etcdPath.
			// It returns an error in case of any failure.
			_, err := scopedKV.Delete(clientv3.WithRequireLeader(ctx), configSync.etcdPath)
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to delete infra meter config")
				return err
			}

			// Return nil to indicate success.
			return nil
		},
	})
}
