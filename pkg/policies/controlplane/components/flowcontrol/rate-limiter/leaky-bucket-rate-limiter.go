package ratelimiter

import (
	"context"
	"path"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

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

type leakyBucketRateLimiterSync struct {
	policyReadAPI     iface.Policy
	rateLimiterProto  *policylangv1.LeakyBucketRateLimiter
	decision          *policysyncv1.LeakyBucketRateLimiterDecision
	decisionWriter    *etcdwriter.Writer
	componentID       string
	configEtcdPaths   []string
	decisionEtcdPaths []string
}

// Name implements runtime.Component.
func (*leakyBucketRateLimiterSync) Name() string { return "LeakyBucketRateLimiter" }

// Type implements runtime.Component.
func (*leakyBucketRateLimiterSync) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (limiterSync *leakyBucketRateLimiterSync) ShortDescription() string {
	return iface.GetSelectorsShortDescription(limiterSync.rateLimiterProto.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*leakyBucketRateLimiterSync) IsActuator() bool { return true }

// NewLeakyBucketRateLimiterAndOptions creates fx options for RateLimiter and also returns agent group name associated with it.
func NewLeakyBucketRateLimiterAndOptions(
	rateLimiterProto *policylangv1.LeakyBucketRateLimiter,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := rateLimiterProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths, decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {

		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.LeakyBucketRateLimiterConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
		decisionEtcdPath := path.Join(paths.LeakyBucketRateLimiterDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	limiterSync := &leakyBucketRateLimiterSync{
		rateLimiterProto: rateLimiterProto,
		decision: &policysyncv1.LeakyBucketRateLimiterDecision{
			BucketCapacity: -1,
		},
		policyReadAPI:     policyReadAPI,
		configEtcdPaths:   configEtcdPaths,
		decisionEtcdPaths: decisionEtcdPaths,
		componentID:       componentID.String(),
	}
	return limiterSync, fx.Options(
		fx.Invoke(
			limiterSync.setupSync,
		),
	), nil
}

func (limiterSync *leakyBucketRateLimiterSync) setupSync(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.LeakyBucketRateLimiterWrapper{
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
func (limiterSync *leakyBucketRateLimiterSync) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	bucketCapacity := inPortReadings.ReadSingleReadingPort("bucket_capacity")
	leakAmount := inPortReadings.ReadSingleReadingPort("leak_amount")
	leakIntervalMs := inPortReadings.ReadSingleReadingPort("leak_interval_ms")

	decision := &policysyncv1.LeakyBucketRateLimiterDecision{
		BucketCapacity: -1,
		LeakAmount:     0,
		LeakInterval:   &durationpb.Duration{Seconds: 0},
	}
	if !bucketCapacity.Valid() || !leakAmount.Valid() || !leakIntervalMs.Valid() {
		return nil, limiterSync.publishDecision(decision)
	}

	leakInterval := time.Duration(leakIntervalMs.Value()) * time.Millisecond

	decision.BucketCapacity = bucketCapacity.Value()
	decision.LeakAmount = leakAmount.Value()
	decision.LeakInterval = &durationpb.Duration{Seconds: int64(leakInterval.Seconds())}

	return nil, limiterSync.publishDecision(decision)
}

func (limiterSync *leakyBucketRateLimiterSync) publishDecision(decision *policysyncv1.LeakyBucketRateLimiterDecision) error {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if !proto.Equal(limiterSync.decision, decision) {
		limiterSync.decision = decision
		// Publish decision
		logger.Debug().Msg("publishing rate limiter decision")
		wrapper := &policysyncv1.LeakyBucketRateLimiterDecisionWrapper{
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
	}
	return nil
}

// DynamicConfigUpdate is a no-op for rate limiter.
func (limiterSync *leakyBucketRateLimiterSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
