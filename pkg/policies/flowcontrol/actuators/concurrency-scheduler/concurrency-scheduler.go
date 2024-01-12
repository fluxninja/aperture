package concurrencyscheduler

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

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	concurrencylimiter "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/concurrency-limiter"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	workloadscheduler "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/workload-scheduler"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/scheduler"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const concurrencySchedulerStatusRoot = "concurrency_scheduler"

var fxTag = config.NameTag(concurrencySchedulerStatusRoot)

func concurrencySchedulerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideConcurrencySchedulerWatchers,
				fx.ResultTags(fxTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupConcurrencySchedulerFactory,
				fx.ParamTags(fxTag),
			),
		),
	)
}

func provideConcurrencySchedulerWatchers(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroupName := ai.GetAgentGroup()

	etcdPath := path.Join(paths.ConcurrencySchedulerConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type concurrencySchedulerFactory struct {
	engineAPI        iface.Engine
	registry         status.Registry
	decisionsWatcher notifiers.Watcher
	distCache        *distcache.DistCache
	auditJobGroup    *jobs.JobGroup
	wsFactory        *workloadscheduler.Factory
	agentGroupName   string
}

// main fx app.
func setupConcurrencySchedulerFactory(
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
	etcdPath := path.Join(paths.ConcurrencySchedulerDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", concurrencySchedulerStatusRoot)

	logger := reg.GetLogger()

	auditJobGroup, err := jobs.NewJobGroup(reg.Child("sync", "audit_jobs"), jobs.JobGroupConfig{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create audit job group")
		return err
	}

	concurrencySchedulerFactory := &concurrencySchedulerFactory{
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
		[]notifiers.FxOptionsFunc{concurrencySchedulerFactory.newConcurrencySchedulerOptions})
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
func (csFactory *concurrencySchedulerFactory) newConcurrencySchedulerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := csFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.ConcurrencySchedulerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.ConcurrencyScheduler == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal concurrency scheduler config")
		return fx.Options(), err
	}

	csProto := wrapperMessage.ConcurrencyScheduler
	csProto.Scheduler, err = workloadscheduler.SanitizeSchedulerProto(csProto.Scheduler)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to sanitize scheduler proto")
		return fx.Options(), err
	}

	cs := &concurrencyScheduler{
		Component: wrapperMessage.GetCommonAttributes(),
		proto:     csProto,
		csFactory: csFactory,
		registry:  reg,
		clock:     clockwork.NewRealClock(),
	}
	cs.name = iface.ComponentKey(cs)

	return fx.Options(
		fx.Invoke(
			cs.setup,
		),
	), nil
}

// concurrencyScheduler implements rate limiter on the data plane side.
type concurrencyScheduler struct {
	schedulers sync.Map
	iface.Component
	registry         status.Registry
	clock            clockwork.Clock
	csFactory        *concurrencySchedulerFactory
	limiter          concurrencylimiter.ConcurrencyLimiter
	proto            *policylangv1.ConcurrencyScheduler
	schedulerMetrics *workloadscheduler.SchedulerMetrics
	name             string
}

func (cs *concurrencyScheduler) setup(lifecycle fx.Lifecycle) error {
	logger := cs.registry.GetLogger()
	etcdKey := paths.AgentComponentKey(cs.csFactory.agentGroupName,
		cs.GetPolicyName(),
		cs.GetComponentId())
	// decision notifier
	decisionUnmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		cs.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = cs.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = cs.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = cs.GetComponentId()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error

			idleDuration := cs.proto.ConcurrencyLimiter.GetMaxIdleTime().AsDuration()
			job := jobs.NewBasicJob(cs.name, cs.audit)
			// register job with job group
			err = cs.csFactory.auditJobGroup.RegisterJob(job, jobs.JobConfig{
				ExecutionPeriod: config.MakeDuration(idleDuration),
			})
			if err != nil {
				return err
			}

			cs.schedulerMetrics, err = cs.csFactory.wsFactory.NewSchedulerMetrics(metricLabels)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create scheduler metrics")
				return err
			}

			inner, err := concurrencylimiter.NewGlobalTokenCounter(
				cs.csFactory.distCache,
				cs.name,
				cs.proto.ConcurrencyLimiter.GetMaxIdleTime().AsDuration(),
				cs.proto.ConcurrencyLimiter.GetMaxInflightDuration().AsDuration(),
			)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create limiter")
				return err
			}
			cs.limiter = inner

			// add decisions notifier
			err = cs.csFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}

			// add to data engine
			err = cs.csFactory.engineAPI.RegisterConcurrencyScheduler(cs)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to register scheduler")
				return err
			}

			return nil
		},
		OnStop: func(context.Context) error {
			var merr, err error

			// remove from data engine
			err = cs.csFactory.engineAPI.UnregisterConcurrencyScheduler(cs)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to unregister rate limiter")
				merr = multierr.Append(merr, err)
			}
			// remove decisions notifier
			err = cs.csFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierr.Append(merr, err)
			}
			err = cs.limiter.Close()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to close limiter")
				merr = multierr.Append(merr, err)
			}
			err = cs.schedulerMetrics.Delete()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete scheduler metrics")
				merr = multierr.Append(merr, err)
			}
			err = cs.csFactory.auditJobGroup.DeregisterJob(cs.name)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to deregister job")
				merr = multierr.Append(merr, err)
			}

			cs.registry.SetStatus(status.NewStatus(nil, merr))

			return merr
		},
	})
	return nil
}

// GetSelectors returns the selectors for the rate limiter.
func (cs *concurrencyScheduler) GetSelectors() []*policylangv1.Selector {
	return cs.proto.GetSelectors()
}

func (cs *concurrencyScheduler) getLabelKey(labels labels.Labels) (string, bool) {
	labelKey := cs.proto.ConcurrencyLimiter.GetLimitByLabelKey()
	var label string
	if labelKey == "" {
		label = "default"
	} else {
		labelValue, found := labels.Get(labelKey)
		if !found {
			return "", false
		}
		label = labelKey + ":" + labelValue
	}
	return label, true
}

// Decide runs the limiter.
func (cs *concurrencyScheduler) Decide(ctx context.Context, labels labels.Labels) *flowcontrolv1.LimiterDecision {
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED
	dropped := false
	label := ""
	reqID := ""
	schedulerInfo := &flowcontrolv1.LimiterDecision_SchedulerInfo{
		WorkloadIndex: metrics.DefaultWorkloadIndex,
	}

	returnDecision := func() *flowcontrolv1.LimiterDecision {
		return &flowcontrolv1.LimiterDecision{
			PolicyName:               cs.GetPolicyName(),
			PolicyHash:               cs.GetPolicyHash(),
			ComponentId:              cs.GetComponentId(),
			Dropped:                  dropped,
			DeniedResponseStatusCode: cs.proto.GetScheduler().GetDeniedResponseStatusCode(),
			Reason:                   reason,
			Details: &flowcontrolv1.LimiterDecision_ConcurrencySchedulerInfo_{
				ConcurrencySchedulerInfo: &flowcontrolv1.LimiterDecision_ConcurrencySchedulerInfo{
					Label:         label,
					WorkloadIndex: schedulerInfo.WorkloadIndex,
					TokensInfo:    schedulerInfo.TokensInfo,
					Priority:      schedulerInfo.Priority,
					RequestId:     reqID,
				},
			},
		}
	}

	if cs.limiter.GetPassThrough() {
		return returnDecision()
	}

	var found bool
	label, found = cs.getLabelKey(labels)

	if !found {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
		return returnDecision()
	}

	// lookup the scheduler
	existing, found := cs.schedulers.Load(label)
	if !found {
		tokenCounter := scheduler.NewGlobalTokenCounter(label, cs.limiter)
		s, err := cs.csFactory.wsFactory.NewScheduler(
			cs.clock,
			cs.registry,
			cs.proto.Scheduler,
			cs,
			tokenCounter,
			cs.schedulerMetrics,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create scheduler")
			return returnDecision()
		}
		// add to the map using LoadOrStore
		existing, _ = cs.schedulers.LoadOrStore(label, s)
	}

	s := existing.(*workloadscheduler.Scheduler)

	schedulerDecision, reqID := s.Decide(ctx, labels)
	schedulerInfo = schedulerDecision.GetLoadSchedulerInfo()
	dropped = schedulerDecision.GetDropped()

	return returnDecision()
}

// Revert returns the tokens to the limiter.
func (cs *concurrencyScheduler) Revert(ctx context.Context, labels labels.Labels, decision *flowcontrolv1.LimiterDecision) {
	// return to the underlying rate limiter
	if csDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_ConcurrencySchedulerInfo_); ok {
		csInfo := csDecision.ConcurrencySchedulerInfo
		tokens := csInfo.TokensInfo.Consumed
		if tokens > 0 {
			_, err := cs.limiter.Return(ctx, csInfo.Label, tokens, csInfo.RequestId)
			if err != nil {
				log.Autosample().Error().Err(err).Msg("Failed to return tokens")
			}
		}
	}
}

// Return returns the tokens to the limiter.
func (cs *concurrencyScheduler) Return(ctx context.Context, label string, tokens float64, requestID string) (bool, error) {
	// return to the underlying rate limiter
	return cs.limiter.Return(ctx, label, tokens, requestID)
}

func (cs *concurrencyScheduler) audit(ctx context.Context) (proto.Message, error) {
	now := time.Now()
	// range through the map and sync the counters
	cs.schedulers.Range(func(label, value interface{}) bool {
		s := value.(*workloadscheduler.Scheduler)

		lastAccess, size := s.Info()

		// if this counter has not synced in a while, then remove it from the map
		if now.After(lastAccess.Add(cs.proto.ConcurrencyLimiter.MaxIdleTime.AsDuration())) &&
			size == 0 {
			cs.schedulers.Delete(label)
			return true
		}
		return true
	})
	return nil, nil
}

func (cs *concurrencyScheduler) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := cs.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		cs.limiter.SetPassThrough(true)
		return
	}

	var wrapperMessage policysyncv1.ConcurrencyLimiterDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.ConcurrencyLimiterDecision == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != cs.GetPolicyHash() {
		return
	}
	limitDecision := wrapperMessage.ConcurrencyLimiterDecision
	cs.limiter.SetCapacity(limitDecision.MaxConcurrency)
	cs.limiter.SetPassThrough(limitDecision.PassThrough)
}

// GetLimiterID returns the limiter ID.
func (cs *concurrencyScheduler) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  cs.GetPolicyName(),
		PolicyHash:  cs.GetPolicyHash(),
		ComponentID: cs.GetComponentId(),
	}
}

// GetLatencyObserver returns the latency observer.
func (cs *concurrencyScheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	return cs.csFactory.wsFactory.GetLatencyObserver(labels)
}

// GetRequestCounter returns counter for tracking number of times rateLimiter was triggered.
func (cs *concurrencyScheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	return cs.csFactory.wsFactory.GetRequestCounter(labels)
}

// GetRampMode is always false for quotaSchedulers.
func (cs *concurrencyScheduler) GetRampMode() bool {
	return false
}
