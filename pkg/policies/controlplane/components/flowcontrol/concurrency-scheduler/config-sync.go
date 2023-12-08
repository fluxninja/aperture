package concurrencyscheduler

import (
	"context"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

type concurrencySchedulerSync struct {
	policyReadAPI             iface.Policy
	concurrencySchedulerProto *policylangv1.ConcurrencyScheduler
	etcdClient                *etcdclient.Client
	componentID               string
	configEtcdPaths           []string
}

// NewConfigSyncOptions creates fx options for syncing LoadScheduler objects with agent groups.
func NewConfigSyncOptions(
	concurrencySchedulerProto *policylangv1.ConcurrencyScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (fx.Option, error) {
	s := concurrencySchedulerProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.ConcurrencySchedulerConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
	}

	css := &concurrencySchedulerSync{
		concurrencySchedulerProto: concurrencySchedulerProto,
		policyReadAPI:             policyReadAPI,
		configEtcdPaths:           configEtcdPaths,
		componentID:               componentID.String(),
	}
	return fx.Options(fx.Invoke(css.setupSync)), nil
}

func (css *concurrencySchedulerSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	logger := css.policyReadAPI.GetStatusRegistry().GetLogger()
	css.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.ConcurrencySchedulerWrapper{
				ConcurrencyScheduler: css.concurrencySchedulerProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  css.policyReadAPI.GetPolicyName(),
					PolicyHash:  css.policyReadAPI.GetPolicyHash(),
					ComponentId: css.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal concurrency scheduler config")
				return err
			}
			for _, configEtcdPath := range css.configEtcdPaths {
				etcdClient.Put(configEtcdPath, string(dat))
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			deleteEtcdPath := func(paths []string) {
				for _, path := range paths {
					etcdClient.Delete(path)
				}
			}
			deleteEtcdPath(css.configEtcdPaths)
			return nil
		},
	})
}
