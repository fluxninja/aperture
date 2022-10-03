package rate

import (
	"context"
	"errors"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
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
)

type rateLimiterSync struct {
	policyReadAPI         iface.Policy
	rateLimiterProto      *policylangv1.RateLimiter
	decision              *policydecisionsv1.RateLimiterDecision
	configEtcdPath        string
	decisionsEtcdPath     string
	dynamicConfigEtcdPath string
	decisionWriter        *etcdwriter.Writer
	dynamicConfigWriter   *etcdwriter.Writer
	agentGroupName        string
	componentIndex        int
}

// NewRateLimiterAndOptions creates fx options for RateLimiter and also returns agent group name associated with it.
func NewRateLimiterAndOptions(
	rateLimiterProto *policylangv1.RateLimiter,
	componentIndex int,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	// Get the agent group name.
	selectorProto := rateLimiterProto.GetSelector()
	if selectorProto == nil {
		return nil, fx.Options(), errors.New("selector is nil")
	}
	agentGroupName := selectorProto.ServiceSelector.GetAgentGroup()
	componentID := common.DataplaneComponentKey(agentGroupName, policyReadAPI.GetPolicyName(), int64(componentIndex))
	configEtcdPath := path.Join(common.RateLimiterConfigPath, componentID)
	decisionsEtcdPath := path.Join(common.RateLimiterDecisionsPath, componentID)
	dynamicConfigEtcdPath := path.Join(common.RateLimiterDynamicConfigPath, componentID)

	limiterSync := &rateLimiterSync{
		rateLimiterProto:      rateLimiterProto,
		decision:              &policydecisionsv1.RateLimiterDecision{},
		policyReadAPI:         policyReadAPI,
		configEtcdPath:        configEtcdPath,
		decisionsEtcdPath:     decisionsEtcdPath,
		dynamicConfigEtcdPath: dynamicConfigEtcdPath,
		componentIndex:        componentIndex,
		agentGroupName:        agentGroupName,
	}
	return limiterSync, fx.Options(
		fx.Invoke(
			limiterSync.setupSync,
		),
	), nil
}

func (limiterSync *rateLimiterSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &wrappersv1.RateLimiterWrapper{
				RateLimiter: limiterSync.rateLimiterProto,
				CommonAttributes: &wrappersv1.CommonAttributes{
					PolicyName:     limiterSync.policyReadAPI.GetPolicyName(),
					PolicyHash:     limiterSync.policyReadAPI.GetPolicyHash(),
					ComponentIndex: int64(limiterSync.componentIndex),
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal rate limiter config")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				limiterSync.configEtcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				logger.Error().Err(err).Msg("failed to put rate limiter config")
				return err
			}
			limiterSync.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			limiterSync.dynamicConfigWriter = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			limiterSync.dynamicConfigWriter.Close()
			limiterSync.decisionWriter.Close()
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), limiterSync.configEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete rate limiter config")
				merr = multierr.Append(merr, err)
			}
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), limiterSync.decisionsEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete rate limiter decisions")
				merr = multierr.Append(merr, err)
			}
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), limiterSync.dynamicConfigEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete rate limiter dynamic config")
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (limiterSync *rateLimiterSync) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	limit, ok := inPortReadings["limit"]
	if !ok {
		return nil, nil
	}

	if len(limit) == 0 {
		return nil, nil
	}

	limitReading := limit[0]
	var limitValue float64
	if !limitReading.Valid() {
		limitValue = -1.0 // no limit is applied
	} else {
		limitValue = limitReading.Value()
	}
	return nil, limiterSync.publishLimit(limitValue)
}

func (limiterSync *rateLimiterSync) publishLimit(limitValue float64) error {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if limiterSync.decision.GetLimit() != limitValue {
		// Save the decision
		limiterSync.decision.Limit = limitValue
		// Publish decision
		logger.Debug().Float64("limit", limitValue).Msg("publishing rate limiter decision")
		wrapper := &wrappersv1.RateLimiterDecisionWrapper{
			RateLimiterDecision: limiterSync.decision,
			CommonAttributes: &wrappersv1.CommonAttributes{
				PolicyName:     limiterSync.policyReadAPI.GetPolicyName(),
				PolicyHash:     limiterSync.policyReadAPI.GetPolicyHash(),
				ComponentIndex: int64(limiterSync.componentIndex),
			},
		}
		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("failed to marshal rate limiter decision")
			return err
		}
		if limiterSync.decisionWriter == nil {
			logger.Panic().Msg("decision writer is nil")
		}
		limiterSync.decisionWriter.Write(limiterSync.decisionsEtcdPath, dat)
	}
	return nil
}

// DynamicConfigUpdate handles overrides.
func (limiterSync *rateLimiterSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	dynamicConfig := &policylangv1.RateLimiter_DynamicConfig{}
	key := limiterSync.rateLimiterProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal dynamic config")
			return
		}
		wrapper := &wrappersv1.RateLimiterDynamicConfigWrapper{
			RateLimiterDynamicConfig: dynamicConfig,
			CommonAttributes: &wrappersv1.CommonAttributes{
				PolicyName:     limiterSync.policyReadAPI.GetPolicyName(),
				PolicyHash:     limiterSync.policyReadAPI.GetPolicyHash(),
				ComponentIndex: int64(limiterSync.componentIndex),
			},
		}
		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("failed to marshal dynamic config")
			return
		}
		if limiterSync.dynamicConfigWriter == nil {
			logger.Panic().Msg("dynamic config writer is nil")
		}
		limiterSync.dynamicConfigWriter.Write(limiterSync.dynamicConfigEtcdPath, dat)
		logger.Info().Msg("rate limiter dynamic config updated")
	}
}
