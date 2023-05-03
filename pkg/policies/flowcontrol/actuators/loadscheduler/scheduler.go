package loadscheduler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/scheduler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/status"
)

// Array of Label Keys for WFQ and Token Bucket Metrics.
var metricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIDLabel}

type workloadSchedulerFactory struct {
	registry status.Registry

	// WFQ Metrics.
	wfqFlowsGaugeVec    *prometheus.GaugeVec
	wfqRequestsGaugeVec *prometheus.GaugeVec

	incomingTokensCounterVec *prometheus.CounterVec
	acceptedTokensCounterVec *prometheus.CounterVec

	workloadLatencySummaryVec *prometheus.SummaryVec
	workloadCounterVec        *prometheus.CounterVec
}

// newWorkloadSchedulerFactory sets up the load scheduler module in the main fx app.
func newWorkloadSchedulerFactory(
	lifecycle fx.Lifecycle,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
) (*workloadSchedulerFactory, error) {
	reg := statusRegistry.Child("component", "scheduler")

	wsFactory := &workloadSchedulerFactory{
		registry: reg,
	}

	wsFactory.wfqFlowsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQFlowsMetricName,
			Help: "A gauge that tracks the number of flows in the WFQScheduler",
		},
		metricLabelKeys,
	)
	wsFactory.wfqRequestsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQRequestsMetricName,
			Help: "A gauge that tracks the number of queued requests in the WFQScheduler",
		},
		metricLabelKeys,
	)
	wsFactory.incomingTokensCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.IncomingTokensMetricName,
			Help: "A counter measuring work incoming into Scheduler",
		},
		metricLabelKeys,
	)
	wsFactory.acceptedTokensCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.AcceptedTokensMetricName,
			Help: "A counter measuring work admitted by Scheduler",
		},
		metricLabelKeys,
	)

	wsFactory.workloadLatencySummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.WorkloadLatencyMetricName,
		Help: "Latency summary of workload",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	wsFactory.workloadCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
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

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var merr error

			err := prometheusRegistry.Register(wsFactory.wfqFlowsGaugeVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.wfqRequestsGaugeVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.incomingTokensCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.acceptedTokensCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.workloadLatencySummaryVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.workloadCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			return merr
		},
		OnStop: func(_ context.Context) error {
			var merr error

			if !prometheusRegistry.Unregister(wsFactory.wfqFlowsGaugeVec) {
				err := fmt.Errorf("failed to unregister wfq_flows metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.wfqRequestsGaugeVec) {
				err := fmt.Errorf("failed to unregister wfq_requests metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.incomingTokensCounterVec) {
				err := fmt.Errorf("failed to unregister incoming_tokens_total metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.acceptedTokensCounterVec) {
				err := fmt.Errorf("failed to unregister accepted_tokens_total metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.workloadLatencySummaryVec) {
				err := fmt.Errorf("failed to unregister workload_latency_ms metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.workloadCounterVec) {
				err := fmt.Errorf("failed to unregister workload_counter metric")
				merr = multierr.Append(merr, err)
			}

			return merr
		},
	})

	return wsFactory, nil
}

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	matchedWorkloads map[int]*policylangv1.Scheduler_Workload_Parameters
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[int, multiMatchResult]

// newLoadSchedulerOptions returns fx options for the load scheduler fx app.
func (lsFactory *workloadSchedulerFactory) newLoadSchedulerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := lsFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.LoadSchedulerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	workloadSchedulerProto := wrapperMessage.LoadScheduler
	if err != nil || workloadSchedulerProto == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal load scheduler config wrapper")
		return fx.Options(), err
	}

	// Scheduler config
	schedulerProto := workloadSchedulerProto.Parameters.Scheduler
	if schedulerProto == nil {
		err = fmt.Errorf("no scheduler specified")
		reg.SetStatus(status.NewStatus(nil, err))
		return fx.Options(), err
	}
	mm := multimatcher.New[int, multiMatchResult]()
	// Loop through the workloads
	for workloadIndex, workloadProto := range schedulerProto.Workloads {
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

	ls := &workloadScheduler{
		Component:                wrapperMessage.GetCommonAttributes(),
		workloadSchedulerProto:   workloadSchedulerProto,
		registry:                 reg,
		workloadSchedulerFactory: lsFactory,
		workloadMultiMatcher:     mm,
	}
	// default workload params is not a required param so it can be nil
	if ls.workloadSchedulerProto.Parameters.Scheduler.DefaultWorkloadParameters == nil {
		p := &policylangv1.Scheduler_Workload_Parameters{}
		config.SetDefaults(p)
		ls.workloadSchedulerProto.Parameters.Scheduler.DefaultWorkloadParameters = p
	}

	return fx.Options(
		fx.Invoke(
			ls.setup,
		),
	), nil
}

type workloadMatcher struct {
	workloadProto *policylangv1.Scheduler_Workload
	workloadIndex int
}

func (wm *workloadMatcher) matchCallback(mmr multiMatchResult) multiMatchResult {
	// mmr.matchedWorkloads is nil on first match.
	if mmr.matchedWorkloads == nil {
		mmr.matchedWorkloads = make(map[int]*policylangv1.Scheduler_Workload_Parameters)
	}
	mmr.matchedWorkloads[wm.workloadIndex] = wm.workloadProto.GetParameters()
	return mmr
}

// workloadScheduler implements load scheduler on the flowcontrol side.
type workloadScheduler struct {
	iface.Component
	scheduler                scheduler.Scheduler
	registry                 status.Registry
	incomingTokensCounter    prometheus.Counter
	acceptedTokensCounter    prometheus.Counter
	workloadSchedulerProto   *policylangv1.LoadScheduler
	workloadSchedulerFactory *workloadSchedulerFactory
	autoTokens               *autoTokens
	workloadMultiMatcher     *multiMatcher
}

// Make sure LoadScheduler implements the iface.LoadScheduler.
var _ iface.Limiter = &workloadScheduler{}

func (ls *workloadScheduler) setup(lifecycle fx.Lifecycle) error {
	// Factories
	lsFactory := ls.workloadSchedulerFactory
	actuatorFactory := lsFactory.actuatorFactory
	autoTokensFactory := lsFactory.autoTokensFactory
	// Form metric labels
	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.PolicyNameLabel] = ls.GetPolicyName()
	metricLabels[metrics.PolicyHashLabel] = ls.GetPolicyHash()
	metricLabels[metrics.ComponentIDLabel] = ls.GetComponentId()
	// Create sub components.
	clock := clockwork.NewRealClock()
	actuator, err := actuatorFactory.newActuator(ls.workloadSchedulerProto.GetActuator(),
		ls, ls.registry, clock, lifecycle, metricLabels)
	if err != nil {
		return err
	}
	if ls.workloadSchedulerProto.GetActuator().WorkloadLatencyBasedTokens {
		autoTokens, err := autoTokensFactory.newAutoTokens(
			ls.GetPolicyName(), ls.GetPolicyHash(),
			lifecycle, ls.GetComponentId(), ls.registry)
		if err != nil {
			return err
		}
		ls.autoTokens = autoTokens
	}

	engineAPI := lsFactory.engineAPI
	wfqFlowsGaugeVec := lsFactory.wfqFlowsGaugeVec
	wfqRequestsGaugeVec := lsFactory.wfqRequestsGaugeVec
	incomingTokensCounterVec := lsFactory.incomingTokensCounterVec
	acceptedTokensCounterVec := lsFactory.acceptedTokensCounterVec

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			retErr := func(err error) error {
				ls.registry.SetStatus(status.NewStatus(nil, err))
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
			ls.scheduler = scheduler.NewWFQScheduler(actuator.tokenBucketLoadMultiplier, clock, wfqMetrics)

			ls.incomingTokensCounter, err = incomingTokensCounterVec.GetMetricWith(metricLabels)
			if err != nil {
				return err
			}
			ls.acceptedTokensCounter, err = acceptedTokensCounterVec.GetMetricWith(metricLabels)
			if err != nil {
				return err
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
			deletedCount := ls.workloadSchedulerFactory.workloadLatencySummaryVec.DeletePartialMatch(metricLabels)
			if deletedCount == 0 {
				log.Warn().Msg("Could not delete workload_latency_ms summary from its metric vector. No traffic to generate metrics?")
			}
			deletedCount = ls.workloadSchedulerFactory.workloadCounterVec.DeletePartialMatch(metricLabels)
			if deletedCount == 0 {
				log.Warn().Msg("Could not delete workload_requests_total counter from its metric vector. No traffic to generate metrics?")
			}

			ls.registry.SetStatus(status.NewStatus(nil, errMulti))
			return errMulti
		},
	})

	return nil
}

// GetSelectors returns selectors.
func (ls *workloadScheduler) GetSelectors() []*policylangv1.Selector {
	return ls.workloadSchedulerProto.GetSelectors()
}

// Decide processes a single flow by load scheduler in a blocking manner.
//
// Context is used to ensure that requests are not scheduled for longer than its deadline allows.
func (ls *workloadScheduler) Decide(ctx context.Context,
	labels map[string]string,
) *flowcontrolv1.LimiterDecision {
	var matchedWorkloadProto *policylangv1.Scheduler_Workload_Parameters
	var matchedWorkloadIndex string
	// match labels against ls.workloadMultiMatcher
	mmr := ls.workloadMultiMatcher.Match(multimatcher.Labels(labels))
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
		matchedWorkloadProto = ls.defaultWorkloadParametersMsg
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
	if ls.workloadSchedulerProto.GetActuator().WorkloadLatencyBasedTokens {
		if tokensAuto, ok := ls.autoTokens.GetTokensForWorkload(matchedWorkloadIndex); ok {
			tokens = tokensAuto
		}
	}

	if matchedWorkloadProto.Tokens != 0 {
		tokens = matchedWorkloadProto.Tokens
	}

	if ls.schedulerParameters.TokensLabelKey != "" {
		if val, ok := labels[ls.schedulerParameters.TokensLabelKey]; ok {
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
		timeout := clientTimeout - ls.schedulerParameters.DecisionDeadlineMargin.AsDuration()
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

	accepted := ls.scheduler.Schedule(reqCtx, req)

	tokensConsumed := uint64(0)
	if accepted {
		tokensConsumed = req.Tokens
	}

	// update load scheduler metrics and decisionType
	ls.incomingTokensCounter.Add(float64(req.Tokens) / 1000)

	if accepted {
		ls.acceptedTokensCounter.Add(float64(req.Tokens) / 1000)
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:  ls.GetPolicyName(),
		PolicyHash:  ls.GetPolicyHash(),
		ComponentId: ls.GetComponentId(),
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
func (ls *workloadScheduler) Revert(labels map[string]string, decision *flowcontrolv1.LimiterDecision) {
	if lsDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_LoadSchedulerInfo_); ok {
		tokens := lsDecision.LoadSchedulerInfo.TokensConsumed
		if tokens > 0 {
			ls.scheduler.Revert(tokens)
		}
	}
}

// GetLimiterID returns the limiter ID.
func (ls *workloadScheduler) GetLimiterID() iface.LimiterID {
	// TODO: move this to limiter base.
	return iface.LimiterID{
		PolicyName:  ls.GetPolicyName(),
		PolicyHash:  ls.GetPolicyHash(),
		ComponentID: ls.GetComponentId(),
	}
}

// GetLatencyObserver returns histogram for specific workload.
func (ls *workloadScheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	latencySummary, err := ls.workloadSchedulerFactory.workloadLatencySummaryVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return nil
	}

	return latencySummary
}

// GetRequestCounter returns request counter for specific workload.
func (ls *workloadScheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := ls.workloadSchedulerFactory.workloadCounterVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting counter")
		return nil
	}

	return counter
}
