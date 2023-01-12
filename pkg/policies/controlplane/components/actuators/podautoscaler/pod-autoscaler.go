package podautoscaler

import (
	"context"
	"errors"
	"path"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

// Module returns the fx options for pod autoscaler in the main app.
func Module() fx.Option {
	return fx.Options(
		scaleReporterModule(),
	)
}

// PodAutoscaler struct.
type podAutoscalerConfigSync struct {
	policyReadAPI      iface.Policy
	podAutoscalerProto *policylangv1.PodAutoscaler
	etcdPath           string
	componentIndex     int
}

// NewPodAutoscalerOptions creates fx options for PodAutoscaler and also returns the agent group name associated with it.
func NewPodAutoscalerOptions(
	podAutoscalerProto *policylangv1.PodAutoscaler,
	componentStackIndex int,
	policyReadAPI iface.Policy,
) (fx.Option, string, error) {
	// Get Agent Group Name from PodAutoscaler.KubernetesObjectSelector.AgentGroup
	k8sObjectSelectorProto := podAutoscalerProto.GetKubernetesObjectSelector()
	if k8sObjectSelectorProto == nil {
		return fx.Options(), "", errors.New("podAutoscaler.Selector is nil")
	}
	agentGroup := k8sObjectSelectorProto.GetAgentGroup()
	etcdPath := path.Join(paths.PodAutoscalerConfigPath,
		paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), int64(componentStackIndex)))
	configSync := &podAutoscalerConfigSync{
		podAutoscalerProto: podAutoscalerProto,
		policyReadAPI:      policyReadAPI,
		etcdPath:           etcdPath,
		componentIndex:     componentStackIndex,
	}

	return fx.Options(
		fx.Invoke(
			configSync.doSync,
		),
	), agentGroup, nil
}

func (configSync *podAutoscalerConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()
	// Add/remove file in lifecycle hooks in order to sync with etcd.
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.PodAutoscalerWrapper{
				PodAutoscaler: configSync.podAutoscalerProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:     configSync.policyReadAPI.GetPolicyName(),
					PolicyHash:     configSync.policyReadAPI.GetPolicyHash(),
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
