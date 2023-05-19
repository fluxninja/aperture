package rate

import (
	"context"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/v2/pkg/etcd/writer"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

type rateLimiterSync struct {
	policyReadAPI     iface.Policy
	rateLimiterProto  *policylangv1.RateLimiter
	decisionWriter    *etcdwriter.Writer
	componentID       string
	configEtcdPaths   []string
	decisionEtcdPaths []string
	customLimits      []*policylangv1.RateLimiter_CustomLimit
}

// Name implements runtime.Component.
func (*rateLimiterSync) Name() string { return "RateLimiter" }

// Type implements runtime.Component.
func (*rateLimiterSync) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (limiterSync *rateLimiterSync) ShortDescription() string {
	return iface.GetSelectorsShortDescription(limiterSync.rateLimiterProto.GetParameters().GetSelectors())
}

// IsActuator implements runtime.Component.
func (*rateLimiterSync) IsActuator() bool { return true }

// NewRateLimiterAndOptions creates fx options for RateLimiter and also returns agent group name associated with it.
func NewRateLimiterAndOptions(
	rateLimiterProto *policylangv1.RateLimiter,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := rateLimiterProto.GetParameters().GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths, decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {

		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.RateLimiterConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
		decisionEtcdPath := path.Join(paths.RateLimiterDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	limiterSync := &rateLimiterSync{
		rateLimiterProto:  rateLimiterProto,
		policyReadAPI:     policyReadAPI,
		configEtcdPaths:   configEtcdPaths,
		decisionEtcdPaths: decisionEtcdPaths,
		componentID:       componentID.String(),
		customLimits:      rateLimiterProto.GetCustomLimits(),
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
			wrapper := &policysyncv1.RateLimiterWrapper{
				RateLimiter: limiterSync.rateLimiterProto,
				CommonAttributes: &policysyncv1.CommonAttributes{
					PolicyName:  limiterSync.policyReadAPI.GetPolicyName(),
					PolicyHash:  limiterSync.policyReadAPI.GetPolicyHash(),
					ComponentId: limiterSync.componentID,
				},
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal rate limiter config")
				return err
			}
			var merr error
			for _, configEtcdPath := range limiterSync.configEtcdPaths {
				_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
					configEtcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
				if err != nil {
					logger.Error().Err(err).Msg("failed to put rate limiter config")
					merr = multierr.Append(merr, err)
				}
			}
			limiterSync.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			return merr
		},
		OnStop: func(ctx context.Context) error {
			limiterSync.decisionWriter.Close()
			deleteEtcdPath := func(paths []string) error {
				var merr error
				for _, path := range paths {
					_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), path)
					if err != nil {
						logger.Error().Err(err).Msgf("failed to delete etcd path %s", path)
						merr = multierr.Append(merr, err)
					}
				}
				return merr
			}

			merr := deleteEtcdPath(limiterSync.configEtcdPaths)
			merr = multierr.Append(merr, deleteEtcdPath(limiterSync.decisionEtcdPaths))
			return merr
		},
	})
	return nil
}

// Execute implements runtime.Component.Execute.
func (limiterSync *rateLimiterSync) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
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
	decision := &policysyncv1.RateLimiterDecision{
		Limit:        limitValue,
		CustomLimits: limiterSync.customLimits,
	}
	// Publish decision
	logger.Debug().Float64("limit", limitValue).Msg("publishing rate limiter decision")
	wrapper := &policysyncv1.RateLimiterDecisionWrapper{
		RateLimiterDecision: decision,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:  limiterSync.policyReadAPI.GetPolicyName(),
			PolicyHash:  limiterSync.policyReadAPI.GetPolicyHash(),
			ComponentId: limiterSync.componentID,
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
	for _, decisionEtcdPath := range limiterSync.decisionEtcdPaths {
		limiterSync.decisionWriter.Write(decisionEtcdPath, dat)
	}
	return nil
}

// DynamicConfigUpdate handles overrides.
func (limiterSync *rateLimiterSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	key := limiterSync.rateLimiterProto.GetCustomLimitsConfigKey()
	customLimits := []*policylangv1.RateLimiter_CustomLimit{}
	if !unmarshaller.IsSet(key) {
		limiterSync.customLimits = limiterSync.rateLimiterProto.GetCustomLimits()
		return
	}
	err := unmarshaller.UnmarshalKey(key, &customLimits)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to unmarshal key %s", key)
		limiterSync.customLimits = limiterSync.rateLimiterProto.GetCustomLimits()
		return
	}
	limiterSync.customLimits = customLimits
}
