package concurrency

import (
	"context"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policydecisionsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/rs/zerolog"
)

// LoadActuator struct.
type LoadActuator struct {
	policyReadAPI  iface.Policy
	decision       *policydecisionsv1.LoadDecision
	etcdPath       string
	writer         *etcdwriter.Writer
	agentGroupName string
	componentIndex int
}

// NewLoadActuatorAndOptions creates load actuator and its fx options.
func NewLoadActuatorAndOptions(
	_ *policylangv1.LoadActuator,
	componentIndex int,
	policyReadAPI iface.Policy,
	agentGroup string,
) (runtime.Component, fx.Option, error) {
	etcdPath := path.Join(common.LoadDecisionsPath,
		common.DataplaneComponentKey(agentGroup, policyReadAPI.GetPolicyName(), int64(componentIndex)))
	lsa := &LoadActuator{
		policyReadAPI:  policyReadAPI,
		agentGroupName: agentGroup,
		componentIndex: componentIndex,
		etcdPath:       etcdPath,
	}
	lsa.decision = &policydecisionsv1.LoadDecision{}

	return lsa, fx.Options(
		fx.Invoke(lsa.setupWriter),
	), nil
}

func (lsa *LoadActuator) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := lsa.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			lsa.writer = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), lsa.etcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete load decision config")
				return err
			}
			lsa.writer.Close()
			return nil
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (lsa *LoadActuator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	// Get the decision from the port
	lm, ok := inPortReadings["load_multiplier"]
	if ok {
		if len(lm) > 0 {
			lmReading := lm[0]
			var lmValue float64
			if !lmReading.Valid() {
				lmValue = 0
			} else {
				if lmReading.Value() <= 0 {
					lmValue = 0
				} else if lmReading.Value() >= 1 {
					lmValue = 1
				} else {
					lmValue = lmReading.Value()
				}
			}
			return nil, lsa.publishLoadMultiplier(lmValue)
		}
	}
	return nil, nil
}

// DynamicConfigUpdate is a no-op for load actuator.
func (lsa *LoadActuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (lsa *LoadActuator) publishLoadMultiplier(loadMultiplier float64) error {
	logger := lsa.policyReadAPI.GetStatusRegistry().GetLogger()
	// Save load multiplier in decision message
	lsa.decision.LoadMultiplier = loadMultiplier
	// Publish decision
	logger.Sample(zerolog.Often).Debug().Float64("loadMultiplier", loadMultiplier).Msg("Publish load decision")
	wrapper := &wrappersv1.LoadDecisionWrapper{
		LoadDecision: lsa.decision,
		CommonAttributes: &wrappersv1.CommonAttributes{
			PolicyName:     lsa.policyReadAPI.GetPolicyName(),
			PolicyHash:     lsa.policyReadAPI.GetPolicyHash(),
			ComponentIndex: int64(lsa.componentIndex),
		},
	}
	dat, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal policy decision")
		return err
	}
	lsa.writer.Write(lsa.etcdPath, dat)
	return nil
}
