package loadscheduler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"path"
	"strconv"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/scheduler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
)

var (
	// FxNameTag is Load Scheduler Watcher's Fx Tag.
	fxNameTag = config.NameTag("load_scheduler_watcher")

	// Array of Label Keys for WFQ and Token Bucket Metrics.
	metricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIDLabel}
)

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
	engineAPI iface.Engine
	registry  status.Registry

	autoTokensFactory *autoTokensFactory
	actuatorFactory   *actuatorFactory

	// WFQ Metrics.
	wfqFlowsGaugeVec    *prometheus.GaugeVec
	wfqRequestsGaugeVec *prometheus.GaugeVec

	// TODO: following will be moved to scheduler.
	incomingTokensCounterVec *prometheus.CounterVec
	acceptedTokensCounterVec *prometheus.CounterVec

	workloadLatencySummaryVec *prometheus.SummaryVec
	workloadCounterVec        *prometheus.CounterVec
}

// setupLoadSchedulerFactory sets up the load scheduler module in the main fx app.
func setupLoadSchedulerFactory(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
	e iface.Engine,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) error {
	agentGroup := ai.GetAgentGroup()

	// Create factories
	actuatorFactory, err := newActuatorFactory(lifecycle, etcdClient, agentGroup, prometheusRegistry)
	if err != nil {
		return err
	}

	autoTokensFactory, err := newAutoTokensFactory(lifecycle, etcdClient, agentGroup)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", "load_scheduler")

	conLimiterFactory := &loadSchedulerFactory{
		engineAPI:         e,
		autoTokensFactory: autoTokensFactory,
		actuatorFactory:   actuatorFactory,
		registry:          reg,
	}

	conLimiterFactory.wfqFlowsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQFlowsMetricName,
			Help: "A gauge that tracks the number of flows in the WFQScheduler",
		},
		metricLabelKeys,
	)
	conLimiterFactory.wfqRequestsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQRequestsMetricName,
			Help: "A gauge that tracks the number of queued requests in the WFQScheduler",
		},
		metricLabelKeys,
	)
	conLimiterFactory.incomingTokensCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.IncomingTokensMetricName,
			Help: "A counter measuring work incoming into Scheduler",
		},
		metricLabelKeys,
	)
	conLimiterFactory.acceptedTokensCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.AcceptedTokensMetricName,
			Help: "A counter measuring work admitted by Scheduler",
		},
		metricLabelKeys,
	)

	conLimiterFactory.workloadLatencySummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.WorkloadLatencyMetricName,
		Help: "Latency summary of workload",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	conLimiterFactory.workloadCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.WorkloadCounterMetricName,
		Help: "Counter of workload requests",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.DecisionTypeLabel,
		metrics.WorkloadIndexLabel,
		metrics.LimiterDroppedLabel,
	})
	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{conLimiterFactory.newLoadSchedulerOptions},
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var merr error

			err := prometheusRegistry.Register(conLimiterFactory.wfqFlowsGaugeVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.wfqRequestsGaugeVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.incomingTokensCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.acceptedTokensCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.workloadLatencySummaryVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(conLimiterFactory.workloadCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			return merr
		},
		OnStop: func(_ context.Context) error {
			var merr error

			if !prometheusRegistry.Unregister(conLimiterFactory.wfqFlowsGaugeVec) {
				err := fmt.Errorf("failed to unregister wfq_flows metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.wfqRequestsGaugeVec) {
				err := fmt.Errorf("failed to unregister wfq_requests metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.incomingTokensCounterVec) {
				err := fmt.Errorf("failed to unregister incoming_tokens_total metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.acceptedTokensCounterVec) {
				err := fmt.Errorf("failed to unregister accepted_tokens_total metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.workloadLatencySummaryVec) {
				err := fmt.Errorf("failed to unregister workload_latency_ms metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(conLimiterFactory.workloadCounterVec) {
				err := fmt.Errorf("failed to unregister workload_counter metric")
				merr = multierr.Append(merr, err)
			}

			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	matchedWorkloads map[int]*policylangv1.LoadScheduler_Scheduler_Workload_Parameters
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[int, multiMatchResult]

// newLoadSchedulerOptions returns fx options for the load scheduler fx app.
func (conLimiterFactory *loadSchedulerFactory) newLoadSchedulerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := conLimiterFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.LoadSchedulerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	loadSchedulerMessage := wrapperMessage.LoadScheduler
	if err != nil || loadSchedulerMessage == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal load scheduler config wrapper")
		return fx.Options(), err
	}

	// Scheduler config
	schedulerMsg := loadSchedulerMessage.Scheduler
	if schedulerMsg == nil {
		err = fmt.Errorf("no scheduler specified")
		reg.SetStatus(status.NewStatus(nil, err))
		return fx.Options(), err
	}
	schedulerParams := schedulerMsg.Parameters
	if schedulerParams == nil {
		err = fmt.Errorf("no scheduler parameters specified")
		reg.SetStatus(status.NewStatus(nil, err))
		return fx.Options(), err
	}
	mm := multimatcher.New[int, multiMatchResult]()
	// Loop through the workloads
	for workloadIndex, workloadProto := range schedulerParams.Workloads {
		labelMatcher, err := selectors.MMExprFromLabelMatcher(workloadProto.GetLabelMatcher())
		if err != nil {
			return fx.Options(), err
		}
		wm := &workloadMatcher{
			workloadIndex: workloadIndex,
			workloadProto: workloadProto,
		}
		err = mm.AddEntry(workloadIndex, labelMatcher, wm.matchCallback)
		if err != nil {
			return fx.Options(), err
		}
	}

	conLimiter := &loadScheduler{
		Component:                    wrapperMessage.GetCommonAttributes(),
		loadSchedulerProto:           loadSchedulerMessage,
		registry:                     reg,
		loadSchedulerFactory:         conLimiterFactory,
		workloadMultiMatcher:         mm,
		defaultWorkloadParametersMsg: schedulerParams.GetDefaultWorkloadParameters(),
		schedulerParameters:          schedulerParams,
	}
	// default workload params is not a required param so it can be nil
	if conLimiter.defaultWorkloadParametersMsg == nil {
		conLimiter.defaultWorkloadParametersMsg = &policylangv1.LoadScheduler_Scheduler_Workload_Parameters{}
		config.SetDefaults(conLimiter.defaultWorkloadParametersMsg)
	}

	return fx.Options(
		fx.Invoke(
			conLimiter.setup,
		),
	), nil
}

type workloadMatcher struct {
	workloadProto *policylangv1.LoadScheduler_Scheduler_Workload
	workloadIndex int
}

func (wm *workloadMatcher) matchCallback(mmr multiMatchResult) multiMatchResult {
	// mmr.matchedWorkloads is nil on first match.
	if mmr.matchedWorkloads == nil {
		mmr.matchedWorkloads = make(map[int]*policylangv1.LoadScheduler_Scheduler_Workload_Parameters)
	}
	mmr.matchedWorkloads[wm.workloadIndex] = wm.workloadProto.GetParameters()
	return mmr
}

// loadScheduler implements load scheduler on the flowcontrol side.
type loadScheduler struct {
	iface.Component
	scheduler                    scheduler.Scheduler
	registry                     status.Registry
	incomingTokensCounter        prometheus.Counter
	acceptedTokensCounter        prometheus.Counter
	loadSchedulerProto           *policylangv1.LoadScheduler
	loadSchedulerFactory         *loadSchedulerFactory
	autoTokens                   *autoTokens
	workloadMultiMatcher         *multiMatcher
	defaultWorkloadParametersMsg *policylangv1.LoadScheduler_Scheduler_Workload_Parameters
	schedulerParameters          *policylangv1.LoadScheduler_Scheduler_Parameters
}

// Make sure LoadScheduler implements the iface.LoadScheduler.
var _ iface.Limiter = &loadScheduler{}

func (conLimiter *loadScheduler) setup(lifecycle fx.Lifecycle) error {
	// Factories
	conLimiterFactory := conLimiter.loadSchedulerFactory
	actuatorFactory := conLimiterFactory.actuatorFactory
	autoTokensFactory := conLimiterFactory.autoTokensFactory
	// Form metric labels
	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = conLimiter.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = conLimiter.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = conLimiter.GetComponentId()
	// Create sub components.
	clock := clockwork.NewRealClock()
	actuator, err := actuatorFactory.newActuator(conLimiter.loadSchedulerProto.GetActuator(),
		conLimiter, conLimiter.registry, clock, lifecycle, metricLabels)
	if err != nil {
		return err
	}
	if conLimiter.schedulerParameters.AutoTokens {
		autoTokens, err := autoTokensFactory.newAutoTokens(
			conLimiter.GetPolicyName(), conLimiter.GetPolicyHash(),
			lifecycle, conLimiter.GetComponentId(), conLimiter.registry)
		if err != nil {
			return err
		}
		conLimiter.autoTokens = autoTokens
	}

	engineAPI := conLimiterFactory.engineAPI
	wfqFlowsGaugeVec := conLimiterFactory.wfqFlowsGaugeVec
	wfqRequestsGaugeVec := conLimiterFactory.wfqRequestsGaugeVec
	incomingTokensCounterVec := conLimiterFactory.incomingTokensCounterVec
	acceptedTokensCounterVec := conLimiterFactory.acceptedTokensCounterVec

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				conLimiter.registry.SetStatus(status.NewStatus(nil, err))
				return err
			}

			wfqFlowsGauge, err := wfqFlowsGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return fmt.Errorf("%w: failed to get wfq flows gauge", err)
			}

			wfqRequestsGauge, err := wfqRequestsGaugeVec.GetMetricWith(metricLabels)
			if err != nil {
				return fmt.Errorf("%w: failed to get wfq requests gauge", err)
			}

			wfqMetrics := &scheduler.WFQMetrics{
				FlowsGauge:        wfqFlowsGauge,
				HeapRequestsGauge: wfqRequestsGauge,
			}

			// setup scheduler
			conLimiter.scheduler = scheduler.NewWFQScheduler(actuator.tokenBucketLoadMultiplier, clock, wfqMetrics)

			conLimiter.incomingTokensCounter, err = incomingTokensCounterVec.GetMetricWith(metricLabels)
			if err != nil {
				return err
			}
			conLimiter.acceptedTokensCounter, err = acceptedTokensCounterVec.GetMetricWith(metricLabels)
			if err != nil {
				return err
			}

			err = engineAPI.RegisterLoadScheduler(conLimiter)
			if err != nil {
				return retErr(err)
			}

			return retErr(nil)
		},
		OnStop: func(context.Context) error {
			var errMulti error

			err := engineAPI.UnregisterLoadScheduler(conLimiter)
			if err != nil {
				errMulti = multierr.Append(errMulti, err)
			}

			// Remove metrics from metric vectors
			deleted := wfqFlowsGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete wfq_flows gauge from its metric vector"))
			}
			deleted = wfqRequestsGaugeVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete wfq_requests gauge from its metric vector"))
			}
			deleted = incomingTokensCounterVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete incoming_tokens_total counter from its metric vector"))
			}
			deleted = acceptedTokensCounterVec.Delete(metricLabels)
			if !deleted {
				errMulti = multierr.Append(errMulti, errors.New("failed to delete accepted_tokens_total counter from its metric vector"))
			}
			deletedCount := conLimiter.loadSchedulerFactory.workloadLatencySummaryVec.DeletePartialMatch(metricLabels)
			if deletedCount == 0 {
				log.Warn().Msg("Could not delete workload_latency_ms summary from its metric vector. No traffic to generate metrics?")
			}
			deletedCount = conLimiter.loadSchedulerFactory.workloadCounterVec.DeletePartialMatch(metricLabels)
			if deletedCount == 0 {
				log.Warn().Msg("Could not delete workload_requests_total counter from its metric vector. No traffic to generate metrics?")
			}

			conLimiter.registry.SetStatus(status.NewStatus(nil, errMulti))
			return errMulti
		},
	})

	return nil
}

// GetSelectors returns selectors.
func (conLimiter *loadScheduler) GetSelectors() []*policylangv1.Selector {
	return conLimiter.loadSchedulerProto.GetSelectors()
}

// Decide processes a single flow by load scheduler in a blocking manner.
//
// Context is used to ensure that requests are not scheduled for longer than its deadline allows.
func (conLimiter *loadScheduler) Decide(ctx context.Context,
	labels map[string]string,
) *flowcontrolv1.LimiterDecision {
	var matchedWorkloadProto *policylangv1.LoadScheduler_Scheduler_Workload_Parameters
	var matchedWorkloadIndex string
	// match labels against conLimiter.workloadMultiMatcher
	mmr := conLimiter.workloadMultiMatcher.Match(multimatcher.Labels(labels))
	// if at least one match, return workload with lowest index
	if len(mmr.matchedWorkloads) > 0 {
		// select the smallest workloadIndex
		smallestWorkloadIndex := math.MaxInt32
		for workloadIndex := range mmr.matchedWorkloads {
			if workloadIndex < smallestWorkloadIndex {
				smallestWorkloadIndex = workloadIndex
			}
		}
		matchedWorkloadProto = mmr.matchedWorkloads[smallestWorkloadIndex]
		matchedWorkloadIndex = strconv.Itoa(smallestWorkloadIndex)
	} else {
		// no match, return default workload
		matchedWorkloadProto = conLimiter.defaultWorkloadParametersMsg
		matchedWorkloadIndex = metrics.DefaultWorkloadIndex
	}

	fairnessLabel := "workload:" + matchedWorkloadIndex

	if val, ok := labels[matchedWorkloadProto.FairnessKey]; ok {
		fairnessLabel = fairnessLabel + "," + val
	}

	tokens := uint64(1)
	// Precedence order (lowest to highest):
	// 1. AutoTokens
	// 2. Workload tokens
	// 3. Label tokens
	if conLimiter.schedulerParameters.AutoTokens {
		if tokensAuto, ok := conLimiter.autoTokens.GetTokensForWorkload(matchedWorkloadIndex); ok {
			tokens = tokensAuto
		}
	}

	if matchedWorkloadProto.Tokens != 0 {
		tokens = matchedWorkloadProto.Tokens
	}

	if conLimiter.schedulerParameters.TokensLabelKey != "" {
		if val, ok := labels[conLimiter.schedulerParameters.TokensLabelKey]; ok {
			if parsedTokens, err := strconv.ParseUint(val, 10, 64); err == nil {
				tokens = parsedTokens
			}
		}
	}

	reqCtx := ctx

	if clientDeadline, hasDeadline := ctx.Deadline(); hasDeadline {
		// The clientDeadline is calculated based on client's timeout, passed
		// as grpc-timeout. Our goal is for the response to be received by the
		// client before its deadline passes (otherwise we risk fail-open on
		// timeout). To allow some headroom for transmitting the response to
		// the client, we set an "internal" deadline to a bit before client's
		// deadline, subtracting the configured margin.
		clientTimeout := time.Until(clientDeadline)
		timeout := clientTimeout - conLimiter.schedulerParameters.DecisionDeadlineMargin.AsDuration()
		if timeout < 0 {
			// we will still schedule the request and it will get
			// dropped if it doesn't get the tokens immediately.
			timeout = 0
		}
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		reqCtx = timeoutCtx
	}

	req := scheduler.Request{
		FairnessLabel: fairnessLabel,
		Priority:      uint8(matchedWorkloadProto.Priority),
		Tokens:        tokens,
	}

	accepted := conLimiter.scheduler.Schedule(reqCtx, req)

	tokensConsumed := uint64(0)
	if accepted {
		tokensConsumed = req.Tokens
	}

	// update load scheduler metrics and decisionType
	conLimiter.incomingTokensCounter.Add(float64(req.Tokens) / 1000)

	if accepted {
		conLimiter.acceptedTokensCounter.Add(float64(req.Tokens) / 1000)
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:  conLimiter.GetPolicyName(),
		PolicyHash:  conLimiter.GetPolicyHash(),
		ComponentId: conLimiter.GetComponentId(),
		Dropped:     !accepted,
		Details: &flowcontrolv1.LimiterDecision_LoadSchedulerInfo_{
			LoadSchedulerInfo: &flowcontrolv1.LimiterDecision_LoadSchedulerInfo{
				WorkloadIndex:  matchedWorkloadIndex,
				TokensConsumed: tokensConsumed,
			},
		},
	}
}

// Revert reverts the decision made by the limiter.
func (conLimiter *loadScheduler) Revert(labels map[string]string, decision *flowcontrolv1.LimiterDecision) {
	if conLimiterDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_LoadSchedulerInfo_); ok {
		tokens := conLimiterDecision.LoadSchedulerInfo.TokensConsumed
		if tokens > 0 {
			conLimiter.scheduler.Revert(tokens)
		}
	}
}

// GetLimiterID returns the limiter ID.
func (conLimiter *loadScheduler) GetLimiterID() iface.LimiterID {
	// TODO: move this to limiter base.
	return iface.LimiterID{
		PolicyName:  conLimiter.GetPolicyName(),
		PolicyHash:  conLimiter.GetPolicyHash(),
		ComponentID: conLimiter.GetComponentId(),
	}
}

// GetLatencyObserver returns histogram for specific workload.
func (conLimiter *loadScheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	latencySummary, err := conLimiter.loadSchedulerFactory.workloadLatencySummaryVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return nil
	}

	return latencySummary
}

// GetRequestCounter returns request counter for specific workload.
func (conLimiter *loadScheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := conLimiter.loadSchedulerFactory.workloadCounterVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting counter")
		return nil
	}

	return counter
}
