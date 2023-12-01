package ratelimiter

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

type rateLimiterSync struct {
	policyReadAPI    iface.Policy
	rateLimiterProto *policylangv1.RateLimiter
	etcdClient       *etcdclient.Client
	componentID      string
	configEtcdPaths  []string
}

// NewConfigSyncOptions creates fx options for syncing LoadScheduler objects with agent groups.
func NewConfigSyncOptions(
	rateLimiterProto *policylangv1.RateLimiter,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (fx.Option, error) {
	s := rateLimiterProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths []string

	for _, agentGroup := range agentGroups {
		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.RateLimiterConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
	}

	limiterSync := &rateLimiterSync{
		rateLimiterProto: rateLimiterProto,
		policyReadAPI:    policyReadAPI,
		configEtcdPaths:  configEtcdPaths,
		componentID:      componentID.String(),
	}
	return fx.Options(fx.Invoke(limiterSync.setupSync)), nil
}

func (rls *rateLimiterSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	logger := rls.policyReadAPI.GetStatusRegistry().GetLogger()
	rls.etcdClient = etcdClient
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.RateLimiterWrapper{
				RateLimiter: rls.rateLimiterProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  rls.policyReadAPI.GetPolicyName(),
					PolicyHash:  rls.policyReadAPI.GetPolicyHash(),
					ComponentId: rls.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal rate limiter config")
				return err
			}
			for _, configEtcdPath := range rls.configEtcdPaths {
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
			deleteEtcdPath(rls.configEtcdPaths)
			return nil
		},
	})
}
