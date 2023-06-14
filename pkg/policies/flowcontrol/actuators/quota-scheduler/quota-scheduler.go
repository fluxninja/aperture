package quotascheduler

import (
	"context"
	"path"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	workloadscheduler "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/workload-scheduler"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/rate-limiter"
	globaltokenbucket "github.com/fluxninja/aperture/v2/pkg/rate-limiter/global-token-bucket"
	lazysync "github.com/fluxninja/aperture/v2/pkg/rate-limiter/lazy-sync"
	"github.com/fluxninja/aperture/v2/pkg/scheduler"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const quotaSchedulerStatusRoot = "quota_scheduler"

var fxTag = config.NameTag(quotaSchedulerStatusRoot)

func quotaSchedulerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideQuotaSchedulerWatchers,
				fx.ResultTags(fxTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupQuotaSchedulerFactory,
				fx.ParamTags(fxTag),
			),
		),
	)
}

func provideQuotaSchedulerWatchers(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroupName := ai.GetAgentGroup()

	etcdPath := path.Join(paths.QuotaSchedulerConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type quotaSchedulerFactory struct {
	engineAPI        iface.Engine
	registry         status.Registry
	decisionsWatcher notifiers.Watcher
	distCache        *distcache.DistCache
	auditJobGroup    *jobs.JobGroup
	wsFactory        *workloadscheduler.Factory
	agentGroupName   string
}

// main fx app.
func setupQuotaSchedulerFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	distCache *distcache.DistCache,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
	wsFactory *workloadscheduler.Factory,
) error {
	agentGroupName := ai.GetAgentGroup()
	etcdPath := path.Join(paths.QuotaSchedulerDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", quotaSchedulerStatusRoot)

	logger := reg.GetLogger()

	auditJobGroup, err := jobs.NewJobGroup(reg.Child("sync", "audit_jobs"), jobs.JobGroupConfig{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create audit job group")
		return err
	}

	quotaSchedulerFactory := &quotaSchedulerFactory{
		engineAPI:        e,
		distCache:        distCache,
		auditJobGroup:    auditJobGroup,
		decisionsWatcher: decisionsWatcher,
		agentGroupName:   agentGroupName,
		registry:         reg,
		wsFactory:        wsFactory,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{quotaSchedulerFactory.newQuotaSchedulerOptions})
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err = auditJobGroup.Start()
			if err != nil {
				return err
			}
			err = decisionsWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			var err, merr error
			err = decisionsWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = auditJobGroup.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			reg.Detach()
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// per component fx app.
func (qsFactory *quotaSchedulerFactory) newQuotaSchedulerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := qsFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.QuotaSchedulerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.QuotaScheduler == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal quota scheduler config")
		return fx.Options(), err
	}

	qsProto := wrapperMessage.QuotaScheduler
	qsProto.Scheduler, err = workloadscheduler.SanitizeSchedulerProto(qsProto.Scheduler)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to sanitize scheduler proto")
		return fx.Options(), err
	}

	qs := &quotaScheduler{
		Component: wrapperMessage.GetCommonAttributes(),
		proto:     qsProto,
		qsFactory: qsFactory,
		registry:  reg,
		clock:     clockwork.NewRealClock(),
	}
	qs.name = iface.ComponentKey(qs)

	return fx.Options(
		fx.Invoke(
			qs.setup,
		),
	), nil
}

// quotaScheduler implements rate limiter on the data plane side.
type quotaScheduler struct {
	schedulers sync.Map
	iface.Component
	registry         status.Registry
	limiter          ratelimiter.RateLimiter
	clock            clockwork.Clock
	qsFactory        *quotaSchedulerFactory
	inner            *globaltokenbucket.GlobalTokenBucket
	proto            *policylangv1.QuotaScheduler
	schedulerMetrics *workloadscheduler.SchedulerMetrics
	name             string
}

func (qs *quotaScheduler) setup(lifecycle fx.Lifecycle) error {
	logger := qs.registry.GetLogger()
	etcdKey := paths.AgentComponentKey(qs.qsFactory.agentGroupName,
		qs.GetPolicyName(),
		qs.GetComponentId())
	// decision notifier
	decisionUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		qs.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = qs.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = qs.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = qs.GetComponentId()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error

			idleDuration := qs.proto.RateLimiter.GetMaxIdleTime().AsDuration()
			job := jobs.NewBasicJob(qs.name, qs.audit)
			// register job with job group
			err = qs.qsFactory.auditJobGroup.RegisterJob(job, jobs.JobConfig{
				ExecutionPeriod: config.MakeDuration(idleDuration),
			})
			if err != nil {
				return err
			}

			qs.schedulerMetrics, err = qs.qsFactory.wsFactory.NewSchedulerMetrics(metricLabels)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create scheduler metrics")
				return err
			}

			qs.inner, err = globaltokenbucket.NewGlobalTokenBucket(
				qs.qsFactory.distCache,
				qs.name,
				qs.proto.RateLimiter.GetInterval().AsDuration(),
				qs.proto.RateLimiter.GetMaxIdleTime().AsDuration(),
				qs.proto.RateLimiter.GetContinuousFill(),
			)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			qs.limiter = qs.inner

			// check whether lazy limiter is enabled
			if lazySyncConfig := qs.proto.RateLimiter.GetLazySync(); lazySyncConfig != nil {
				if lazySyncConfig.GetEnabled() {
					qs.limiter, err = lazysync.NewLazySyncRateLimiter(qs.limiter,
						qs.proto.RateLimiter.GetInterval().AsDuration(),
						lazySyncConfig.GetNumSync(),
						qs.qsFactory.auditJobGroup)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to create lazy limiter")
						return err
					}
				}
			}

			// add decisions notifier
			err = qs.qsFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}

			// add to data engine
			err = qs.qsFactory.engineAPI.RegisterScheduler(qs)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register scheduler")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error

			// remove from data engine
			err = qs.qsFactory.engineAPI.UnregisterScheduler(qs)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister rate limiter")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = qs.qsFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			err = qs.limiter.Close()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to close limiter")
				merr = multierr.Append(merr, err)
			}
			err = qs.schedulerMetrics.Delete()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete scheduler metrics")
				merr = multierr.Append(merr, err)
			}
			err = qs.qsFactory.auditJobGroup.DeregisterJob(qs.name)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to deregister job")
				merr = multierr.Append(merr, err)
			}

			qs.registry.SetStatus(status.NewStatus(nil, merr))

			return merr
		},
	})
	return nil
}

// GetSelectors returns the selectors for the rate limiter.
func (qs *quotaScheduler) GetSelectors() []*policylangv1.Selector {
	return qs.proto.GetSelectors()
}

func (qs *quotaScheduler) getLabelKey(labels map[string]string) (string, bool) {
	labelKey := qs.proto.RateLimiter.GetLabelKey()
	var label string
	if labelKey == "" {
		label = "default"
	} else {
		labelValue, found := labels[labelKey]
		if !found {
			return "", false
		}
		label = labelKey + ":" + labelValue
	}
	return label, true
}

// Decide runs the limiter.
func (qs *quotaScheduler) Decide(ctx context.Context, labels map[string]string) iface.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED
	dropped := false
	label := ""
	schedulerInfo := &flowcontrolv1.LimiterDecision_SchedulerInfo{
		WorkloadIndex: metrics.DefaultWorkloadIndex,
	}

	returnDecision := func() iface.LimiterDecision {
		return iface.LimiterDecision{
			LimiterDecision: &flowcontrolv1.LimiterDecision{
				PolicyName:  qs.GetPolicyName(),
				PolicyHash:  qs.GetPolicyHash(),
				ComponentId: qs.GetComponentId(),
				Dropped:     dropped,
				Reason:      reason,
				Details: &flowcontrolv1.LimiterDecision_QuotaSchedulerInfo_{
					QuotaSchedulerInfo: &flowcontrolv1.LimiterDecision_QuotaSchedulerInfo{
						Label:         label,
						SchedulerInfo: schedulerInfo,
					},
				},
			},
		}
	}

	if qs.limiter.GetPassThrough() {
		return returnDecision()
	}

	var found bool
	label, found = qs.getLabelKey(labels)

	if !found {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
		return returnDecision()
	}

	// lookup the scheduler
	existing, found := qs.schedulers.Load(label)
	if !found {
		tokenBucket := scheduler.NewGlobalTokenBucket(label, qs.limiter)
		s, err := qs.qsFactory.wsFactory.NewScheduler(
			qs.clock,
			qs.registry,
			qs.proto.Scheduler,
			qs,
			tokenBucket,
			qs.schedulerMetrics,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create scheduler")
			return returnDecision()
		}
		// add to the map using LoadOrStore
		existing, _ = qs.schedulers.LoadOrStore(label, s)
	}

	s := existing.(*workloadscheduler.Scheduler)

	schedulerDecision := s.Decide(ctx, labels)
	schedulerInfo = schedulerDecision.GetLoadSchedulerInfo()
	dropped = schedulerDecision.GetDropped()

	return returnDecision()
}

// Revert returns the tokens to the limiter.
func (qs *quotaScheduler) Revert(ctx context.Context, labels map[string]string, decision *flowcontrolv1.LimiterDecision) {
	// return to the underlying rate limiter
	if qsDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_QuotaSchedulerInfo_); ok {
		tokens := qsDecision.QuotaSchedulerInfo.SchedulerInfo.TokensConsumed
		if tokens > 0 {
			label, found := qs.getLabelKey(labels)
			if found {
				qs.limiter.TakeIfAvailable(ctx, label, float64(-tokens))
			}
		}
	}
}

func (qs *quotaScheduler) audit(ctx context.Context) (proto.Message, error) {
	now := time.Now()
	// range through the map and sync the counters
	qs.schedulers.Range(func(label, value interface{}) bool {
		s := value.(*workloadscheduler.Scheduler)

		lastAccess, size := s.Info()

		// if this counter has not synced in a while, then remove it from the map
		if now.After(lastAccess.Add(qs.proto.RateLimiter.MaxIdleTime.AsDuration())) &&
			size == 0 {
			qs.schedulers.Delete(label)
			return true
		}
		return true
	})
	return nil, nil
}

func (qs *quotaScheduler) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := qs.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		qs.limiter.SetPassThrough(true)
		return
	}

	var wrapperMessage policysyncv1.RateLimiterDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.RateLimiterDecision == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != qs.GetPolicyHash() {
		return
	}
	limitDecision := wrapperMessage.RateLimiterDecision
	qs.inner.SetBucketCapacity(limitDecision.BucketCapacity)
	qs.inner.SetFillAmount(limitDecision.FillAmount)
	qs.inner.SetPassThrough(limitDecision.PassThrough)
}

// GetLimiterID returns the limiter ID.
func (qs *quotaScheduler) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  qs.GetPolicyName(),
		PolicyHash:  qs.GetPolicyHash(),
		ComponentID: qs.GetComponentId(),
	}
}

// GetLatencyObserver returns the latency observer.
func (qs *quotaScheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	return qs.qsFactory.wsFactory.GetLatencyObserver(labels)
}

// GetRequestCounter returns counter for tracking number of times rateLimiter was triggered.
func (qs *quotaScheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	return qs.qsFactory.wsFactory.GetRequestCounter(labels)
}
