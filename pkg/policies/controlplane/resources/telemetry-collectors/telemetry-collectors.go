package telemetrycollectors

import (
	"context"
	"encoding/json"
	"path"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
)

type tcConfigSync struct {
	policyReadAPI  iface.Policy
	tcProto        *policylangv1.TelemetryCollector
	etcdPath       string
	agentGroupName string
}

// NewTelemetryCollectorsOptions creates fx options InfraMeters.
func NewTelemetryCollectorsOptions(
	tcProtos []*policylangv1.TelemetryCollector,
	policyBaseAPI iface.Policy,
) (fx.Option, error) {
	var options []fx.Option

	for i, tcProto := range tcProtos {
		agentGroup := tcProto.GetAgentGroup()
		etcdPath := path.Join(paths.TelemetryCollectorConfigPath,
			paths.TelemetryCollectorKey(agentGroup, policyBaseAPI.GetPolicyName(), i))
		configSync := &tcConfigSync{
			tcProto:        tcProto,
			policyReadAPI:  policyBaseAPI,
			agentGroupName: agentGroup,
			etcdPath:       etcdPath,
		}
		options = append(options, fx.Invoke(configSync.doSync))
	}

	return fx.Options(options...), nil
}

func (configSync *tcConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	// Get the logger instance from the status registry.
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()

	// Append fx.Hook to the lifecycle.
	lifecycle.Append(fx.Hook{
		// OnStart hook will be called when the application starts.
		OnStart: func(ctx context.Context) error {
			// Marshal the telemetry collector using json marshaler.
			dat, err := json.Marshal(configSync.tcProto)
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to marshal telemetry collector config")
				return err
			}

			// Put the marshaled data in etcd using the provided etcdPath and LeaseID.
			// It returns an error in case of any failure.
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				configSync.etcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				// Log the error and return it in case of any failure.
				logger.Error().Err(err).Msg("Failed to put telemetry collector config")
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
				logger.Error().Err(err).Msg("Failed to delete telemetry collector config")
				return err
			}

			// Return nil to indicate success.
			return nil
		},
	})

	// Return nil to indicate success.
	return nil
}
