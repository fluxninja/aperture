package loadscheduler

import (
	"context"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	workloadscheduler "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/workload-scheduler"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/scheduler"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// FxNameTag is Load Scheduler Watcher's Fx Tag.
var fxNameTag = config.NameTag("load_scheduler_watcher")

// loadSchedulerModule returns the fx options for flowcontrol side pieces of load scheduler in the main fx app.
func loadSchedulerModule() fx.Option {
	return fx.Options(
		// Tag the watcher so that other modules can find it.
		fx.Provide(
			fx.Annotate(
				provideWatcher,
				fx.ResultTags(fxNameTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupLoadSchedulerFactory,
				fx.ParamTags(fxNameTag),
			),
		),
	)
}

// provideWatcher provides pointer to load scheduler watcher.
func provideWatcher(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	// Get Agent Group from host info gatherer
	agentGroupName := ai.GetAgentGroup()
	// Scope the sync to the agent group.
	etcdPath := path.Join(paths.LoadSchedulerConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

type loadSchedulerFactory struct {
	engineAPI                          iface.Engine
	registry                           status.Registry
	loadDecisionWatcher                notifiers.Watcher
	tokenBucketLMGaugeVec              *prometheus.GaugeVec
	tokenBucketFillRateGaugeVec        *prometheus.GaugeVec
	tokenBucketBucketCapacityGaugeVec  *prometheus.GaugeVec
	tokenBucketAvailableTokensGaugeVec *prometheus.GaugeVec
	wsFactory                          *workloadscheduler.Factory
	agentGroupName                     string
}

// setupLoadSchedulerFactory sets up the load scheduler module in the main fx app.
func setupLoadSchedulerFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	registry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
	wsFactory *workloadscheduler.Factory,
) error {
	reg := registry.Child("component", "load_scheduler")

	agentGroup := ai.GetAgentGroup()

	// Scope the sync to the agent group.
	etcdDecisionsPath := path.Join(paths.LoadSchedulerDecisionsPath, paths.AgentGroupPrefix(agentGroup))
	loadDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdDecisionsPath)
	if err != nil {
		return err
	}

	lsFactory := &loadSchedulerFactory{
		engineAPI:           e,
		registry:            reg,
		wsFactory:           wsFactory,
		loadDecisionWatcher: loadDecisionWatcher,
		agentGroupName:      ai.GetAgentGroup(),
	}

	// Initialize and register the WFQ and Token Bucket Metric Vectors
	lsFactory.tokenBucketLMGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketLMMetricName,
			Help: "A gauge that tracks the load multiplier",
		},
		workloadscheduler.MetricLabelKeys,
	)
	lsFactory.tokenBucketFillRateGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketFillRateMetricName,
			Help: "A gauge that tracks the fill rate of token bucket in tokens/sec",
		},
		workloadscheduler.MetricLabelKeys,
	)
	lsFactory.tokenBucketBucketCapacityGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketCapacityMetricName,
			Help: "A gauge that tracks the capacity of token bucket",
		},
		workloadscheduler.MetricLabelKeys,
	)
	lsFactory.tokenBucketAvailableTokensGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketAvailableMetricName,
			Help: "A gauge that tracks the number of tokens available in token bucket",
		},
		workloadscheduler.MetricLabelKeys,
	)

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{lsFactory.newLoadSchedulerOptions},
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var merr error

			err := prometheusRegistry.Register(lsFactory.tokenBucketLMGaugeVec)
			if err != nil {
				return err
			}
			err = prometheusRegistry.Register(lsFactory.tokenBucketFillRateGaugeVec)
			if err != nil {
				return err
			}
			err = prometheusRegistry.Register(lsFactory.tokenBucketBucketCapacityGaugeVec)
			if err != nil {
				return err
			}
			err = prometheusRegistry.Register(lsFactory.tokenBucketAvailableTokensGaugeVec)
			if err != nil {
				return err
			}

			err = lsFactory.loadDecisionWatcher.Start()
			if err != nil {
				return err
			}

			return merr
		},
		OnStop: func(_ context.Context) error {
			var merr error

			err := lsFactory.loadDecisionWatcher.Stop()
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			if !prometheusRegistry.Unregister(lsFactory.tokenBucketLMGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketLMMetricName)
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(lsFactory.tokenBucketFillRateGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketFillRateMetricName)
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(lsFactory.tokenBucketBucketCapacityGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketCapacityMetricName)
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(lsFactory.tokenBucketAvailableTokensGaugeVec) {
				err := fmt.Errorf("failed to unregister " + metrics.TokenBucketAvailableMetricName)
				merr = multierr.Append(merr, err)
			}
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// newLoadSchedulerOptions returns fx options for the load scheduler fx app.
func (lsFactory *loadSchedulerFactory) newLoadSchedulerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	registry status.Registry,
) (fx.Option, error) {
	logger := lsFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.LoadSchedulerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	loadSchedulerProto := wrapperMessage.LoadScheduler
	if err != nil || loadSchedulerProto == nil {
		registry.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal load scheduler config wrapper")
		return fx.Options(), err
	}
	loadSchedulerProto.Parameters.Scheduler, err = workloadscheduler.SanitizeSchedulerProto(loadSchedulerProto.Parameters.Scheduler)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to sanitize scheduler proto")
		return fx.Options(), err
	}

	ls := &loadScheduler{
		Component:            wrapperMessage.GetCommonAttributes(),
		proto:                loadSchedulerProto,
		registry:             registry,
		loadSchedulerFactory: lsFactory,
		clock:                clockwork.NewRealClock(),
	}

	return fx.Options(
		fx.Invoke(
			ls.setup,
		),
	), nil
}

// loadScheduler implements load scheduler on the flowcontrol side.
type loadScheduler struct {
	// TODO: comment to self: why do we depend on Scheduler to implement Decide and Revert in this case?
	iface.Component
	scheduler            *workloadscheduler.Scheduler
	registry             status.Registry
	proto                *policylangv1.LoadScheduler
	loadSchedulerFactory *loadSchedulerFactory
	clock                clockwork.Clock
	tokenBucket          *scheduler.LoadMultiplierTokenBucket
	schedulerMetrics     *workloadscheduler.SchedulerMetrics
}

// Make sure LoadScheduler implements the iface.LoadScheduler.
var _ iface.Limiter = &loadScheduler{}

func (ls *loadScheduler) setup(lifecycle fx.Lifecycle) error {
	// Factories
	lsFactory := ls.loadSchedulerFactory
	wsFactory := lsFactory.wsFactory
	// Form metric labels
	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = ls.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = ls.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = ls.GetComponentId()

	etcdKey := paths.AgentComponentKey(lsFactory.agentGroupName,
		ls.GetPolicyName(),
		ls.GetComponentId())

	decisionUnmarshaller, protoErr := config.NewProtobufUnmarshaller(nil)
	if protoErr != nil {
		return protoErr
	}

	// decision notifier
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaller,
		ls.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}

	engineAPI := lsFactory.engineAPI

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				ls.registry.SetStatus(status.NewStatus(nil, err))
				return err
			}

			tbLMGauge, err := lsFactory.tokenBucketLMGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(fmt.Errorf("%w: Failed to get token bucket LM gauge", err))
			}

			tokenBucketFillRateGauge, err := lsFactory.tokenBucketFillRateGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(fmt.Errorf("%w: Failed to get token bucket fill rate gauge", err))
			}

			tokenBucketBucketCapacityGauge, err := lsFactory.tokenBucketBucketCapacityGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(fmt.Errorf("%w: Failed to get token bucket bucket capacity gauge", err))
			}

			tokenBucketAvailableTokensGauge, err := lsFactory.tokenBucketAvailableTokensGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return retErr(fmt.Errorf("%w: Failed to get token bucket available tokens gauge", err))
			}

			tbMetrics := &scheduler.TokenBucketMetrics{
				FillRateGauge:        tokenBucketFillRateGauge,
				BucketCapacityGauge:  tokenBucketBucketCapacityGauge,
				AvailableTokensGauge: tokenBucketAvailableTokensGauge,
			}

			// Initialize the token bucket (non continuous tracking mode)
			// TODO: 30s is also used by the controller. Define the constant at a common location.
			ls.tokenBucket = scheduler.NewLoadMultiplierTokenBucket(ls.clock, 30, time.Second, tbLMGauge, tbMetrics)

			ls.schedulerMetrics, err = wsFactory.NewSchedulerMetrics(metricLabels)
			if err != nil {
				return retErr(err)
			}

			// Create a new scheduler
			ls.scheduler, err = wsFactory.NewScheduler(
				ls.clock,
				ls.registry,
				ls.proto.Parameters.Scheduler,
				ls,
				ls.tokenBucket,
				ls.schedulerMetrics,
			)
			if err != nil {
				return retErr(err)
			}

			err = lsFactory.loadDecisionWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				return retErr(err)
			}

			err = engineAPI.RegisterScheduler(ls)
			if err != nil {
				return retErr(err)
			}

			return retErr(nil)
		},
		OnStop: func(context.Context) error {
			var errMulti error

			err := engineAPI.UnregisterScheduler(ls)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			protoErr = lsFactory.loadDecisionWatcher.RemoveKeyNotifier(decisionNotifier)
			if protoErr != nil {
				errMulti = multierr.Append(errMulti, protoErr)
			}

			// delete the metrics
			if ls.schedulerMetrics != nil {
				err = ls.schedulerMetrics.Delete()
				if err != nil {
					errMulti = multierr.Append(errMulti, err)
				}
			}

			deleted := lsFactory.tokenBucketLMGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketLMMetricName+" from its metric vector"))
			}
			deleted = lsFactory.tokenBucketFillRateGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketFillRateMetricName+" gauge from its metric vector"))
			}
			deleted = lsFactory.tokenBucketBucketCapacityGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketCapacityMetricName+" gauge from its metric vector"))
			}
			deleted = lsFactory.tokenBucketAvailableTokensGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete "+metrics.TokenBucketAvailableMetricName+" gauge from its metric vector"))
			}

			ls.registry.SetStatus(status.NewStatus(nil, errMulti))
			return errMulti
		},
	})

	return nil
}

// GetSelectors returns selectors.
func (ls *loadScheduler) GetSelectors() []*policylangv1.Selector {
	return ls.proto.Parameters.GetSelectors()
}

// GetLimiterID returns the limiter ID.
func (ls *loadScheduler) GetLimiterID() iface.LimiterID {
	// TODO: move this to limiter base.
	return iface.LimiterID{
		PolicyName:  ls.GetPolicyName(),
		PolicyHash:  ls.GetPolicyHash(),
		ComponentID: ls.GetComponentId(),
	}
}

func (ls *loadScheduler) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := ls.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision was removed, set pass through mode")
		ls.tokenBucket.SetPassThrough(true)
		return
	}

	var wrapperMessage policysyncv1.LoadDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil {
		statusMsg := "Failed to unmarshal config wrapper"
		logger.Warn().Err(err).Msg(statusMsg)
		ls.registry.SetStatus(status.NewStatus(nil, err))
		return
	}

	loadDecision := wrapperMessage.GetLoadDecision()
	if loadDecision == nil {
		statusMsg := "load decision is nil"
		logger.Error().Msg(statusMsg)
		ls.registry.SetStatus(status.NewStatus(nil, fmt.Errorf("failed to get load decision from LoadDecisionWrapper: %s", statusMsg)))
		return
	}

	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		statusMsg := "common attributes is nil"
		logger.Error().Msg(statusMsg)
		ls.registry.SetStatus(status.NewStatus(nil, fmt.Errorf("failed to get common attributes from LoadDecisionWrapper: %s", statusMsg)))
		return
	}

	// check if this decision is for the same policy id as what we have
	if commonAttributes.PolicyHash != ls.GetPolicyHash() {
		err = errors.New("policy id mismatch")
		statusMsg := fmt.Sprintf("Expected policy hash: %s, Got: %s", ls.GetPolicyHash(), commonAttributes.PolicyHash)
		logger.Warn().Err(err).Msg(statusMsg)
		ls.registry.SetStatus(status.NewStatus(nil, err))
		return
	}

	logger.Autosample().Debug().Bool("passThrough", loadDecision.PassThrough).Float64("loadMultiplier", loadDecision.LoadMultiplier).Msg("Setting load multiplier")
	ls.tokenBucket.SetLoadDecisionValues(loadDecision)
	ls.tokenBucket.SetPassThrough(loadDecision.PassThrough)
	ls.scheduler.SetEstimatedTokens(loadDecision.TokensByWorkloadIndex)
}

// Decide implements iface.Limiter.
func (ls *loadScheduler) Decide(ctx context.Context, labels labels.Labels) *flowcontrolv1.LimiterDecision {
	limiterDecision, _ := ls.scheduler.Decide(ctx, labels)
	return limiterDecision
}

// Revert implements iface.Limiter.
func (ls *loadScheduler) Revert(ctx context.Context, labels labels.Labels, decision *flowcontrolv1.LimiterDecision) {
	if lsDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_LoadSchedulerInfo); ok {
		tokens := lsDecision.LoadSchedulerInfo.TokensInfo.Consumed
		if tokens > 0 {
			ls.tokenBucket.Return(ctx, tokens, "")
		}
	}
}

// GetRampMode is required by iface.Limiter.
func (ls *loadScheduler) GetRampMode() bool {
	return false
}

// GetRequestCounter is required by iface.Limiter.
func (ls *loadScheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	return ls.scheduler.GetRequestCounter(labels)
}

// GetLatencyObserver is required by iface.Scheduler.
func (ls *loadScheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	return ls.scheduler.GetLatencyObserver(labels)
}
