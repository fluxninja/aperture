package loadscheduler

import (
	"context"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

type loadSchedulerConfigSync struct {
	policyBaseAPI      iface.Policy
	loadSchedulerProto *policylangv1.LoadScheduler
	etcdPath           string
	componentID        string
}

// NewConfigSyncOptions creates fx options for syncing LoadScheduler objects with agent groups.
func NewConfigSyncOptions(
	loadSchedulerProto *policylangv1.LoadScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (fx.Option, error) {
	options := []fx.Option{}

	s := loadSchedulerProto.Parameters.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	for _, agentGroup := range agentGroups {
		etcdPath := path.Join(
			paths.LoadSchedulerConfigPath,
			paths.AgentComponentKey(
				agentGroup,
				policyReadAPI.GetPolicyName(),
				componentID.String(),
			),
		)
		configSync := &loadSchedulerConfigSync{
			loadSchedulerProto: loadSchedulerProto,
			policyBaseAPI:      policyReadAPI,
			etcdPath:           etcdPath,
			componentID:        componentID.String(),
		}
		options = append(options, fx.Invoke(configSync.doSync))
	}

	return fx.Options(options...), nil
}

func (configSync *loadSchedulerConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := configSync.policyBaseAPI.GetStatusRegistry().GetLogger()
	// Add/remove file in lifecycle hooks in order to sync with etcd.
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.LoadSchedulerWrapper{
				LoadScheduler: configSync.loadSchedulerProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  configSync.policyBaseAPI.GetPolicyName(),
					PolicyHash:  configSync.policyBaseAPI.GetPolicyHash(),
					ComponentId: configSync.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal flux meter config")
				return err
			}
			etcdClient.Put(configSync.etcdPath, string(dat))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			etcdClient.Delete(configSync.etcdPath)
			return nil
		},
	})

	return nil
}
