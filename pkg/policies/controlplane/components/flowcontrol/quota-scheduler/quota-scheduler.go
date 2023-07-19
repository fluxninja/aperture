package quotascheduler

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
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime/tristate"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

type quotaSchedulerSync struct {
	policyReadAPI       iface.Policy
	quotaSchedulerProto *policylangv1.QuotaScheduler
	decision            *policysyncv1.RateLimiterDecision
	decisionWriter      *etcdwriter.Writer
	componentID         string
	configEtcdPaths     []string
	decisionEtcdPaths   []string
}

// Name implements runtime.Component.
func (*quotaSchedulerSync) Name() string { return "QuotaScheduler" }

// Type implements runtime.Component.
func (*quotaSchedulerSync) Type() runtime.ComponentType { return runtime.ComponentTypeSink }

// ShortDescription implements runtime.Component.
func (limiterSync *quotaSchedulerSync) ShortDescription() string {
	return iface.GetSelectorsShortDescription(limiterSync.quotaSchedulerProto.GetSelectors())
}

// IsActuator implements runtime.Component.
func (*quotaSchedulerSync) IsActuator() bool { return true }

// NewQuotaSchedulerAndOptions creates fx options for QuotaScheduler.
func NewQuotaSchedulerAndOptions(
	quotaSchedulerProto *policylangv1.QuotaScheduler,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	s := quotaSchedulerProto.GetSelectors()

	agentGroups := selectors.UniqueAgentGroups(s)

	var configEtcdPaths, decisionEtcdPaths []string

	for _, agentGroup := range agentGroups {

		etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), componentID.String())
		configEtcdPath := path.Join(paths.QuotaSchedulerConfigPath, etcdKey)
		configEtcdPaths = append(configEtcdPaths, configEtcdPath)
		decisionEtcdPath := path.Join(paths.QuotaSchedulerDecisionsPath, etcdKey)
		decisionEtcdPaths = append(decisionEtcdPaths, decisionEtcdPath)
	}

	limiterSync := &quotaSchedulerSync{
		quotaSchedulerProto: quotaSchedulerProto,
		decision:            &policysyncv1.RateLimiterDecision{},
		policyReadAPI:       policyReadAPI,
		configEtcdPaths:     configEtcdPaths,
		decisionEtcdPaths:   decisionEtcdPaths,
		componentID:         componentID.String(),
	}
	return limiterSync, fx.Options(
		fx.Invoke(
			limiterSync.setupSync,
		),
	), nil
}

func (limiterSync *quotaSchedulerSync) setupSync(scopedKV *etcdclient.SessionScopedKV, lifecycle fx.Lifecycle) error {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper := &policysyncv1.QuotaSchedulerWrapper{
				QuotaScheduler: limiterSync.quotaSchedulerProto,
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
				_, err = scopedKV.Put(clientv3.WithRequireLeader(ctx), configEtcdPath, string(dat))
				if err != nil {
					logger.Error().Err(err).Msg("failed to put rate limiter config")
					merr = multierr.Append(merr, err)
				}
			}
			limiterSync.decisionWriter = etcdwriter.NewWriter(&scopedKV.KVWrapper)
			return merr
		},
		OnStop: func(ctx context.Context) error {
			limiterSync.decisionWriter.Close()
			deleteEtcdPath := func(paths []string) error {
				var merr error
				for _, path := range paths {
					_, err := scopedKV.Delete(clientv3.WithRequireLeader(ctx), path)
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
func (limiterSync *quotaSchedulerSync) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	bucketCapacity := inPortReadings.ReadSingleReadingPort("bucket_capacity")
	fillAmount := inPortReadings.ReadSingleReadingPort("fill_amount")
	ptBool := tristate.FromReading(inPortReadings.ReadSingleReadingPort("pass_through"))

	decision := &policysyncv1.RateLimiterDecision{
		PassThrough: true,
	}
	if !bucketCapacity.Valid() || !fillAmount.Valid() {
		return nil, limiterSync.publishDecision(decision)
	}

	decision.BucketCapacity = bucketCapacity.Value()
	decision.FillAmount = fillAmount.Value()
	decision.PassThrough = ptBool.IsTrue()

	return nil, limiterSync.publishDecision(decision)
}

func (limiterSync *quotaSchedulerSync) publishDecision(decision *policysyncv1.RateLimiterDecision) error {
	logger := limiterSync.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if !proto.Equal(limiterSync.decision, decision) {
		limiterSync.decision = decision
		// Publish decision
		logger.Debug().Msg("publishing rate limiter decision")
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
	}
	return nil
}

// DynamicConfigUpdate is a no-op for rate limiter.
func (limiterSync *quotaSchedulerSync) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
