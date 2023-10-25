package classifier

import (
	"context"
	"path"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

type classifierConfigSync struct {
	policyReadAPI   iface.Policy
	classifierProto *policylangv1.Classifier
	etcdPath        string
	agentGroupName  string
	classifierIndex int
}

// NewClassifierOptions creates fx options for classifier.
func NewClassifierOptions(
	classifierIndex int,
	classifierProto *policylangv1.Classifier,
	policyBaseAPI iface.Policy,
) (fx.Option, error) {
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

func (configSync *classifierConfigSync) doSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) {
	logger := configSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.ClassifierWrapper{
				Classifier: configSync.classifierProto,
				ClassifierAttributes: &policysyncv1.ClassifierAttributes{
					PolicyName:      configSync.policyReadAPI.GetPolicyName(),
					PolicyHash:      configSync.policyReadAPI.GetPolicyHash(),
					ClassifierIndex: int64(configSync.classifierIndex),
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal classifier")
				return err
			}
			etcdClient.Put(configSync.etcdPath, string(dat))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			etcdClient.Delete(configSync.etcdPath)
			return nil
		},
	})
}
