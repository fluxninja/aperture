package classifier

import (
	"context"
	"path"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
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
	classifierIndex int64,
	classifierProto *policylangv1.Classifier,
	policyBaseAPI iface.Policy,
) (fx.Option, error) {
	// Deprecated 1.8.0
	flowSelectorProto := classifierProto.GetFlowSelector()
	if flowSelectorProto != nil {
		selector := &policylangv1.Selector{
			ControlPoint: flowSelectorProto.FlowMatcher.ControlPoint,
			LabelMatcher: flowSelectorProto.FlowMatcher.LabelMatcher,
			Service:      flowSelectorProto.ServiceSelector.Service,
			AgentGroup:   flowSelectorProto.ServiceSelector.AgentGroup,
		}
		classifierProto.Selectors = append(classifierProto.Selectors, selector)
		classifierProto.FlowSelector = nil
	}

	var options []fx.Option

	s := classifierProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	for _, agentGroup := range agentGroups {
		etcdPath := path.Join(paths.ClassifiersPath,
			paths.ClassifierKey(agentGroup, policyBaseAPI.GetPolicyName(), classifierIndex))
		configSync := &classifierConfigSync{
			classifierProto: classifierProto,
			policyReadAPI:   policyBaseAPI,
			agentGroupName:  agentGroup,
			etcdPath:        etcdPath,
			classifierIndex: classifierIndex,
		}
		options = append(options, fx.Invoke(configSync.doSync))
	}

	return fx.Options(options...), nil
}

func (configSync *classifierConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.ClassifierWrapper{
				Classifier: configSync.classifierProto,
				ClassifierAttributes: &policysyncv1.ClassifierAttributes{
					PolicyName:      configSync.policyReadAPI.GetPolicyName(),
					PolicyHash:      configSync.policyReadAPI.GetPolicyHash(),
					ClassifierIndex: configSync.classifierIndex,
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
