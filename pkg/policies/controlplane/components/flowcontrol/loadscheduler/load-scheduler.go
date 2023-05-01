package loadscheduler

import (
	"context"
	"path"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type loadSchedulerConfigSync struct {
	policyBaseAPI      iface.Policy
	loadSchedulerProto *policylangv1.LoadScheduler
	etcdPath           string
	componentID        string
}

// NewLoadSchedulerOptions creates fx options for LoadScheduler and also returns the agent group name associated with it.
func NewLoadSchedulerOptions(
	loadSchedulerProto *policylangv1.LoadScheduler,
	componentStackID string,
	policyReadAPI iface.Policy,
) (fx.Option, []string, error) {
	// Deprecated 1.8.0
	flowSelectorProto := loadSchedulerProto.GetFlowSelector()
	if flowSelectorProto != nil {
		selector := &policylangv1.Selector{
			ControlPoint: flowSelectorProto.FlowMatcher.ControlPoint,
			LabelMatcher: flowSelectorProto.FlowMatcher.LabelMatcher,
			Service:      flowSelectorProto.ServiceSelector.Service,
			AgentGroup:   flowSelectorProto.ServiceSelector.AgentGroup,
		}
		loadSchedulerProto.Selectors = append(loadSchedulerProto.Selectors, selector)
		loadSchedulerProto.FlowSelector = nil
	}

	options := []fx.Option{}

	s := loadSchedulerProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	for _, agentGroup := range agentGroups {
		etcdPath := path.Join(paths.LoadSchedulerConfigPath,
			paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentStackID))
		configSync := &loadSchedulerConfigSync{
			loadSchedulerProto: loadSchedulerProto,
			policyBaseAPI:      policyReadAPI,
			etcdPath:           etcdPath,
			componentID:        componentStackID,
		}
		options = append(options, fx.Invoke(configSync.doSync))
	}

	return fx.Options(options...), agentGroups, nil
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
