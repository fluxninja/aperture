package quotascheduler

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

type quotaSchedulerSync struct {
	policyReadAPI       iface.Policy
	quotaSchedulerProto *policylangv1.QuotaScheduler
	etcdClient          *etcdclient.Client
	componentID         string
	configEtcdPaths     []string
}

// NewConfigSyncOptions creates fx options for syncing QuotaScheduler objects with agent groups.
func NewConfigSyncOptions(
	quotaSchedulerProto *policylangv1.QuotaScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (fx.Option, error) {
	s := quotaSchedulerProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.QuotaSchedulerConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
	}

	qss := &quotaSchedulerSync{
		quotaSchedulerProto: quotaSchedulerProto,
		policyReadAPI:       policyReadAPI,
		configEtcdPaths:     configEtcdPaths,
		componentID:         componentID.String(),
	}
	return fx.Options(fx.Invoke(qss.setupSync)), nil
}

func (qss *quotaSchedulerSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	logger := qss.policyReadAPI.GetStatusRegistry().GetLogger()
	qss.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.QuotaSchedulerWrapper{
				QuotaScheduler: qss.quotaSchedulerProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  qss.policyReadAPI.GetPolicyName(),
					PolicyHash:  qss.policyReadAPI.GetPolicyHash(),
					ComponentId: qss.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal quota scheduler config")
				return err
			}
			for _, configEtcdPath := range qss.configEtcdPaths {
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
			deleteEtcdPath(qss.configEtcdPaths)
			return nil
		},
	})
}
