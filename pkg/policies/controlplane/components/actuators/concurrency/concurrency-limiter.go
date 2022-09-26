package concurrency

import (
	"context"
	"errors"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type concurrencyLimiterConfigSync struct {
	policyBaseAPI           iface.PolicyBase
	concurrencyLimiterProto *policylangv1.ConcurrencyLimiter
	etcdPath                string
	agentGroupName          string
	componentIndex          int
}

// NewConcurrencyLimiterOptions creates fx options for ConcurrencyLimiter and also returns the agent group name associated with it.
func NewConcurrencyLimiterOptions(
	concurrencyLimiterProto *policylangv1.ConcurrencyLimiter,
	componentStackIndex int,
	policyBaseAPI iface.PolicyBase,
) (fx.Option, string, error) {
	// Get Agent Group Name from ConcurrencyLimiter.Scheduler.Selector.AgentGroup
	schedulerProto := concurrencyLimiterProto.GetScheduler()
	if schedulerProto == nil {
		return fx.Options(), "", errors.New("concurrencyLimiter.Scheduler is nil")
	}
	selectorProto := schedulerProto.GetSelector()
	if selectorProto == nil {
		return fx.Options(), "", errors.New("concurrencyLimiter.Scheduler.Selector is nil")
	}
	agentGroupName := selectorProto.ServiceSelector.GetAgentGroup()
	etcdPath := path.Join(common.ConcurrencyLimiterConfigPath, common.DataplaneComponentKey(agentGroupName, policyBaseAPI.GetPolicyName(), int64(componentStackIndex)))
	configSync := &concurrencyLimiterConfigSync{
		concurrencyLimiterProto: concurrencyLimiterProto,
		policyBaseAPI:           policyBaseAPI,
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
	// Add/remove file in lifecycle hooks in order to sync with etcd.
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &wrappersv1.ConcurrencyLimiterWrapper{
				ConcurrencyLimiter: configSync.concurrencyLimiterProto,
				ComponentIndex:     int64(configSync.componentIndex),
				PolicyName:         configSync.policyBaseAPI.GetPolicyName(),
				PolicyHash:         configSync.policyBaseAPI.GetPolicyHash(),
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal flux meter config")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				configSync.etcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				log.Error().Err(err).Msg("Failed to put flux meter config")
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), configSync.etcdPath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete flux meter config")
				return err
			}
			return nil
		},
	})

	return nil
}
