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
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/paths"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

type rateLimiterSync struct {
	policyReadAPI     iface.Policy
	rateLimiterProto  *policylangv1.RateLimiter
	decision          *policydecisionsv1.RateLimiterDecision
	configEtcdPath    string
	decisionsEtcdPath string
	decisionWriter    *etcdwriter.Writer
	agentGroupName    string
	componentIndex    int
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
	agentGroupName := selectorProto.GetAgentGroup()
	componentID := paths.DataplaneComponentKey(agentGroupName, policyReadAPI.GetPolicyName(), int64(componentIndex))
	configEtcdPath := path.Join(paths.RateLimiterConfigPath, componentID)
	decisionsEtcdPath := path.Join(paths.RateLimiterDecisionsPath, componentID)

	limiterSync := &rateLimiterSync{
		rateLimiterProto:  rateLimiterProto,
		decision:          &policydecisionsv1.RateLimiterDecision{},
		policyReadAPI:     policyReadAPI,
		configEtcdPath:    configEtcdPath,
		decisionsEtcdPath: decisionsEtcdPath,
		componentIndex:    componentIndex,
		agentGroupName:    agentGroupName,
	}
	return limiterSync, fx.Options(
		fx.Invoke(
			limiterSync.setupSync,
		),
	), nil
}

func (limiterSync *rateLimiterSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &wrappersv1.RateLimiterWrapper{
				RateLimiter:    limiterSync.rateLimiterProto,
				ComponentIndex: int64(limiterSync.componentIndex),
				PolicyName:     limiterSync.policyReadAPI.GetPolicyName(),
				PolicyHash:     limiterSync.policyReadAPI.GetPolicyHash(),
			}
			dat, err := proto.Marshal(wrapper)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal rate limiter config")
				return err
			}
			_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ctx),
				limiterSync.configEtcdPath, string(dat), clientv3.WithLease(etcdClient.LeaseID))
			if err != nil {
				log.Error().Err(err).Msg("failed to put rate limiter config")
				return err
			}
			limiterSync.decisionWriter = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			limiterSync.decisionWriter.Close()
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), limiterSync.configEtcdPath)
			if err != nil {
				log.Error().Err(err).Msg("failed to delete rate limiter config")
				merr = multierr.Append(merr, err)
			}
			_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), limiterSync.decisionsEtcdPath)
			if err != nil {
				log.Error().Err(err).Msg("failed to delete rate limiter decisions")
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
	// Publish only if there's a change
	if limiterSync.decision.GetLimit() != limitValue {
		// Save the decision
		limiterSync.decision.Limit = limitValue
		// Publish decision
		log.Debug().Float64("limit", limitValue).Msg("publishing rate limiter decision")
		wrapper := &wrappersv1.RateLimiterDecisionWrapper{
			RateLimiterDecision: limiterSync.decision,
			ComponentIndex:      int64(limiterSync.componentIndex),
			PolicyName:          limiterSync.policyReadAPI.GetPolicyName(),
			PolicyHash:          limiterSync.policyReadAPI.GetPolicyHash(),
		}
		dat, err := proto.Marshal(wrapper)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal rate limiter decision")
			return err
		}
		if limiterSync.decisionWriter == nil {
			log.Panic().Msg("decision writer is nil")
		}
		limiterSync.decisionWriter.Write(limiterSync.decisionsEtcdPath, dat)
	}
	return nil
}
