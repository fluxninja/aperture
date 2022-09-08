package classifier

import (
	"context"
	"errors"
	"path"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

type classifierConfigSync struct {
	policyBaseAPI   iface.PolicyBase
	classifierProto *policylangv1.Classifier
	etcdPath        string
	agentGroupName  string
	classifierIndex int64
}

// NewClassifierOptions creates fx options for classifier.
func NewClassifierOptions(
	index int64,
	classifierProto *policylangv1.Classifier,
	policyBaseAPI iface.PolicyBase,
) (fx.Option, error) {
	// Get Agent Group Name for Classifier.Selector.AgentGroup
	selectorProto := classifierProto.GetSelector()
	if selectorProto == nil {
		return nil, errors.New("Classifier.Selector is nil")
	}
	agentGroup := selectorProto.GetAgentGroup()

	etcdPath := path.Join(paths.ClassifiersConfigPath, paths.ClassifierKey(agentGroup, policyBaseAPI.GetPolicyName(), index))
	configSync := &classifierConfigSync{
		classifierProto: classifierProto,
		policyBaseAPI:   policyBaseAPI,
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
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &wrappersv1.ClassifierWrapper{
				PolicyName:      configSync.policyBaseAPI.GetPolicyName(),
				PolicyHash:      configSync.policyBaseAPI.GetPolicyHash(),
				ClassifierIndex: configSync.classifierIndex,
				Classifier:      configSync.classifierProto,
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal classifier")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				configSync.etcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				log.Error().Err(err).Msg("Failed to put classifier")
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), configSync.etcdPath)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete classifier")
				return err
			}
			return nil
		},
	})

	return nil
}
