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

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/scheduler"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
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
	wsFactory                          *SchedulerFactory
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
) error {
	reg := registry.Child("component", "load_scheduler")

	agentGroup := ai.GetAgentGroup()

	// Scope the sync to the agent group.
	etcdDecisionsPath := path.Join(paths.LoadSchedulerDecisionsPath,
		paths.AgentGroupPrefix(agentGroup))
	loadDecisionWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdDecisionsPath)
	if err != nil {
		return err
	}

	wsFactory, err := NewSchedulerFactory(
		lifecycle,
		reg,
		prometheusRegistry,
	)
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
		metricLabelKeys,
	)
	lsFactory.tokenBucketFillRateGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketFillRateMetricName,
			Help: "A gauge that tracks the fill rate of token bucket in tokens/sec",
		},
		metricLabelKeys,
	)
	lsFactory.tokenBucketBucketCapacityGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketCapacityMetricName,
			Help: "A gauge that tracks the capacity of token bucket",
		},
		metricLabelKeys,
	)
	lsFactory.tokenBucketAvailableTokensGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.TokenBucketAvailableMetricName,
			Help: "A gauge that tracks the number of tokens available in token bucket",
		},
		metricLabelKeys,
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

	// Scheduler config
	schedulerProto := loadSchedulerProto.Parameters.Scheduler
	if schedulerProto == nil {
		err = fmt.Errorf("no scheduler specified")
		registry.SetStatus(status.NewStatus(nil, err))
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
	*Scheduler
	iface.Component
	registry                  status.Registry
	proto                     *policylangv1.LoadScheduler
	loadSchedulerFactory      *loadSchedulerFactory
	clock                     clockwork.Clock
	tokenBucketLoadMultiplier *scheduler.TokenBucketLoadMultiplier
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

			tokenBucketLMGauge, err := lsFactory.tokenBucketLMGaugeVec.GetMetricWith(metricLabels)
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

			tokenBucketMetrics := &scheduler.TokenBucketLoadMultiplierMetrics{
				LMGauge: tokenBucketLMGauge,
				TokenBucketMetrics: &scheduler.TokenBucketMetrics{
					FillRateGauge:        tokenBucketFillRateGauge,
					BucketCapacityGauge:  tokenBucketBucketCapacityGauge,
					AvailableTokensGauge: tokenBucketAvailableTokensGauge,
				},
			}

			// Initialize the token bucket (non continuous tracking mode)
			ls.tokenBucketLoadMultiplier = scheduler.NewTokenBucketLoadMultiplier(ls.clock.Now(), 10, time.Second, tokenBucketMetrics)
			// Initialize with PassThrough mode
			ls.tokenBucketLoadMultiplier.SetPassThrough(true)

			// Create a new scheduler
			ls.Scheduler, err = wsFactory.NewScheduler(
				ls.registry,
				ls.proto.Parameters.Scheduler,
				ls,
				ls.tokenBucketLoadMultiplier,
				ls.clock,
				metricLabels,
			)
			if err != nil {
				return retErr(err)
			}

			err = lsFactory.loadDecisionWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				return retErr(err)
			}

			err = engineAPI.RegisterLoadScheduler(ls)
			if err != nil {
				return retErr(err)
			}

			return retErr(nil)
		},
		OnStop: func(context.Context) error {
			var errMulti error

			err := engineAPI.UnregisterLoadScheduler(ls)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			protoErr = lsFactory.loadDecisionWatcher.RemoveKeyNotifier(decisionNotifier)
			if protoErr != nil {
				errMulti = multierr.Append(errMulti, protoErr)
			}

			// Stop the scheduler
			err = ls.Close()
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
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
		ls.tokenBucketLoadMultiplier.SetPassThrough(true)
		return
	}

	var wrapperMessage policysyncv1.LoadDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	loadDecision := wrapperMessage.LoadDecision
	if err != nil || loadDecision == nil {
		statusMsg := "Failed to unmarshal config wrapper"
		logger.Warn().Err(err).Msg(statusMsg)
		ls.registry.SetStatus(status.NewStatus(nil, err))
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		statusMsg := "Failed to get common attributes from LoadDecisionWrapper"
		logger.Error().Err(err).Msg(statusMsg)
		ls.registry.SetStatus(status.NewStatus(nil, err))
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

	if loadDecision.PassThrough {
		logger.Autosample().Debug().Msg("Setting pass through mode")
		ls.tokenBucketLoadMultiplier.SetPassThrough(true)
	} else {
		logger.Autosample().Debug().Float64("loadMultiplier", loadDecision.LoadMultiplier).Msg("Setting load multiplier")
		ls.tokenBucketLoadMultiplier.SetLoadMultiplier(ls.clock.Now(), loadDecision.LoadMultiplier)
		ls.tokenBucketLoadMultiplier.SetPassThrough(false)
	}

	ls.SetEstimatedTokens(loadDecision.TokensByWorkloadIndex)
}
