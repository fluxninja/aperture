package concurrencylimiter

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

type concurrencyLimiterSync struct {
	policyReadAPI           iface.Policy
	concurrencyLimiterProto *policylangv1.ConcurrencyLimiter
	etcdClient              *etcdclient.Client
	componentID             string
	configEtcdPaths         []string
}

// NewConfigSyncOptions creates fx options for syncing ConcurrencyLimiter objects with agent groups.
func NewConfigSyncOptions(
	concurrencyLimiterProto *policylangv1.ConcurrencyLimiter,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (fx.Option, error) {
	s := concurrencyLimiterProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.ConcurrencyLimiterConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
	}

	limiterSync := &concurrencyLimiterSync{
		concurrencyLimiterProto: concurrencyLimiterProto,
		policyReadAPI:           policyReadAPI,
		configEtcdPaths:         configEtcdPaths,
		componentID:             componentID.String(),
	}
	return fx.Options(fx.Invoke(limiterSync.setupSync)), nil
}

func (cls *concurrencyLimiterSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	logger := cls.policyReadAPI.GetStatusRegistry().GetLogger()
	cls.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			wrapper := &policysyncv1.ConcurrencyLimiterWrapper{
				ConcurrencyLimiter: cls.concurrencyLimiterProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  cls.policyReadAPI.GetPolicyName(),
					PolicyHash:  cls.policyReadAPI.GetPolicyHash(),
					ComponentId: cls.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal concurrency limiter config")
				return err
			}
			for _, configEtcdPath := range cls.configEtcdPaths {
				etcdClient.Put(configEtcdPath, string(dat))
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			deleteEtcdPath := func(paths []string) {
				for _, path := range paths {
					etcdClient.Delete(path)
				}
			}
			deleteEtcdPath(cls.configEtcdPaths)
			return nil
		},
	})
}
