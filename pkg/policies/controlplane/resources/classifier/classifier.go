package classifier

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

type classifierConfigSync struct {
	policyReadAPI   iface.Policy
	classifierProto *policylangv1.Classifier
	etcdPath        string
	agentGroupName  string
	classifierIndex int64
}

// NewClassifierOptions creates fx options for classifier.
func NewClassifierOptions(
	index int64,
	classifierProto *policylangv1.Classifier,
	policyBaseAPI iface.Policy,
) (fx.Option, error) {
	// Get Agent Group Name for Classifier.Selector.AgentGroup
	flowSelectorProto := classifierProto.GetFlowSelector()
	if flowSelectorProto == nil {
		return nil, errors.New("Classifier.Selector is nil")
	}
	agentGroup := flowSelectorProto.ServiceSelector.GetAgentGroup()

	etcdPath := path.Join(paths.ClassifiersPath,
		paths.ClassifierKey(agentGroup, policyBaseAPI.GetPolicyName(), index))
	configSync := &classifierConfigSync{
		classifierProto: classifierProto,
		policyReadAPI:   policyBaseAPI,
		agentGroupName:  agentGroup,
		etcdPath:        etcdPath,
		classifierIndex: index,
	}

	return fx.Options(
		fx.Invoke(
			configSync.doSync,
		),
	), nil
}

func (configSync *classifierConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.ClassifierWrapper{
				Classifier: configSync.classifierProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:     configSync.policyReadAPI.GetPolicyName(),
					PolicyHash:     configSync.policyReadAPI.GetPolicyHash(),
					ComponentIndex: configSync.classifierIndex,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal classifier")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				configSync.etcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				logger.Error().Err(err).Msg("Failed to put classifier")
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), configSync.etcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete classifier")
				return err
			}
			return nil
		},
	})

	return nil
}
