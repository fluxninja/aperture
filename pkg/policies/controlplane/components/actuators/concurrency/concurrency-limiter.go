package concurrency

import (
	"context"
	"errors"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type concurrencyLimiterConfigSync struct {
	policyBaseAPI           iface.Policy
	concurrencyLimiterProto *policylangv1.ConcurrencyLimiter
	etcdPath                string
	agentGroupName          string
	componentIndex          int
}

// NewConcurrencyLimiterOptions creates fx options for ConcurrencyLimiter and also returns the agent group name associated with it.
func NewConcurrencyLimiterOptions(
	concurrencyLimiterProto *policylangv1.ConcurrencyLimiter,
	componentStackIndex int,
	policyReadAPI iface.Policy,
) (fx.Option, string, error) {
	// Get Agent Group Name from ConcurrencyLimiter.Scheduler.Selector.AgentGroup
	flowSelectorProto := concurrencyLimiterProto.GetFlowSelector()
	if flowSelectorProto == nil {
		return fx.Options(), "", errors.New("concurrencyLimiter.Selector is nil")
	}
	agentGroupName := flowSelectorProto.ServiceSelector.GetAgentGroup()
	etcdPath := path.Join(paths.ConcurrencyLimiterConfigPath,
		paths.FlowControlComponentKey(agentGroupName, policyReadAPI.GetPolicyName(), int64(componentStackIndex)))
	configSync := &concurrencyLimiterConfigSync{
		concurrencyLimiterProto: concurrencyLimiterProto,
		policyBaseAPI:           policyReadAPI,
		etcdPath:                etcdPath,
		componentIndex:          componentStackIndex,
		agentGroupName:          agentGroupName,
	}

	return fx.Options(
		fx.Invoke(
			configSync.doSync,
		),
	), agentGroupName, nil
}

func (configSync *concurrencyLimiterConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := configSync.policyBaseAPI.GetStatusRegistry().GetLogger()
	// Add/remove file in lifecycle hooks in order to sync with etcd.
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.ConcurrencyLimiterWrapper{
				ConcurrencyLimiter: configSync.concurrencyLimiterProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:     configSync.policyBaseAPI.GetPolicyName(),
					PolicyHash:     configSync.policyBaseAPI.GetPolicyHash(),
					ComponentIndex: int64(configSync.componentIndex),
				},
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
