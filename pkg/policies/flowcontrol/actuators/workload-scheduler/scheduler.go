package workloadscheduler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	multimatcher "github.com/fluxninja/aperture/v2/pkg/multi-matcher"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/v2/pkg/scheduler"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// Module provides the fx options for the workload scheduler.
func Module() fx.Option {
	return fx.Provide(newFactory)
}

// MetricLabelKeys is an array of Label Keys for WFQ and Token Bucket Metrics.
var MetricLabelKeys = []string{metrics.PolicyNameLabel, metrics.PolicyHashLabel, metrics.ComponentIDLabel}

// Factory is a factory for creating load schedulers.
type Factory struct {
	registry status.Registry

	// WFQ Metrics.
	wfqFlowsGaugeVec    *prometheus.GaugeVec
	wfqRequestsGaugeVec *prometheus.GaugeVec

	incomingTokensCounterVec *prometheus.CounterVec
	acceptedTokensCounterVec *prometheus.CounterVec

	workloadLatencySummaryVec *prometheus.SummaryVec
	workloadCounterVec        *prometheus.CounterVec
}

// newFactory sets up the load scheduler module in the main fx app.
func newFactory(
	lifecycle fx.Lifecycle,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
) (*Factory, error) {
	reg := statusRegistry.Child("component", "scheduler")

	wsFactory := &Factory{
		registry: reg,
	}

	wsFactory.wfqFlowsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQFlowsMetricName,
			Help: "A gauge that tracks the number of flows in the WFQScheduler",
		},
		MetricLabelKeys,
	)
	wsFactory.wfqRequestsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metrics.WFQRequestsMetricName,
			Help: "A gauge that tracks the number of queued requests in the WFQScheduler",
		},
		MetricLabelKeys,
	)
	wsFactory.incomingTokensCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.IncomingTokensMetricName,
			Help: "A counter measuring work incoming into Scheduler",
		},
		MetricLabelKeys,
	)
	wsFactory.acceptedTokensCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.AcceptedTokensMetricName,
			Help: "A counter measuring work admitted by Scheduler",
		},
		MetricLabelKeys,
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

// GetLatencyObserver returns a latency observer for a given workload.
func (wsFactory *Factory) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	latencySummary, err := wsFactory.workloadLatencySummaryVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return nil
	}

	return latencySummary
}

// GetRequestCounter returns a request counter for a given workload.
func (wsFactory *Factory) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := wsFactory.workloadCounterVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting counter")
		return nil
	}

	return counter
}

// SchedulerMetrics is a struct that holds all metrics for Scheduler.
type SchedulerMetrics struct {
	wfqMetrics   *scheduler.WFQMetrics
	metricLabels prometheus.Labels
	wsFactory    *Factory
}

// NewSchedulerMetrics creates a new SchedulerMetrics instance.
func (wsFactory *Factory) NewSchedulerMetrics(metricLabels prometheus.Labels) (*SchedulerMetrics, error) {
	wfqFlowsGauge, err := wsFactory.wfqFlowsGaugeVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get wfq flows gauge", err)
	}

	wfqRequestsGauge, err := wsFactory.wfqRequestsGaugeVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get wfq requests gauge", err)
	}

	incomingTokensCounter, err := wsFactory.incomingTokensCounterVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, err
	}

	acceptedTokensCounter, err := wsFactory.acceptedTokensCounterVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, err
	}

	wfqMetrics := &scheduler.WFQMetrics{
		FlowsGauge:            wfqFlowsGauge,
		HeapRequestsGauge:     wfqRequestsGauge,
		IncomingTokensCounter: incomingTokensCounter,
		AcceptedTokensCounter: acceptedTokensCounter,
	}

	return &SchedulerMetrics{
		wfqMetrics:   wfqMetrics,
		metricLabels: metricLabels,
		wsFactory:    wsFactory,
	}, nil
}

// Delete removes all metrics from metric vectors.
func (metrics *SchedulerMetrics) Delete() error {
	var merr error

	// Remove metrics from metric vectors
	deleted := metrics.wsFactory.wfqFlowsGaugeVec.Delete(metrics.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete wfq_flows gauge from its metric vector"))
	}
	deleted = metrics.wsFactory.wfqRequestsGaugeVec.Delete(metrics.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete wfq_requests gauge from its metric vector"))
	}
	deleted = metrics.wsFactory.incomingTokensCounterVec.Delete(metrics.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete incoming_tokens_total counter from its metric vector"))
	}
	deleted = metrics.wsFactory.acceptedTokensCounterVec.Delete(metrics.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete accepted_tokens_total counter from its metric vector"))
	}
	deletedCount := metrics.wsFactory.workloadLatencySummaryVec.DeletePartialMatch(metrics.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_latency_ms summary from its metric vector. No traffic to generate metrics?")
	}
	deletedCount = metrics.wsFactory.workloadCounterVec.DeletePartialMatch(metrics.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_requests_total counter from its metric vector. No traffic to generate metrics?")
	}
	return merr
}

// Scheduler implements load scheduler on the flowcontrol side.
type Scheduler struct {
	component             iface.Component
	scheduler             scheduler.Scheduler
	registry              status.Registry
	proto                 *policylangv1.Scheduler
	workloadMultiMatcher  *multiMatcher
	tokensByWorkloadIndex map[string]uint64
	metrics               *SchedulerMetrics
	mutex                 sync.RWMutex
}

// NewScheduler returns fx options for the load scheduler fx app.
func (wsFactory *Factory) NewScheduler(
	clk clockwork.Clock,
	registry status.Registry,
	proto *policylangv1.Scheduler,
	component iface.Component,
	tokenManger scheduler.TokenManager,
	metrics *SchedulerMetrics,
) (*Scheduler, error) {
	mm := multimatcher.New[int, multiMatchResult]()
	// Loop through the workloads
	for workloadIndex, workloadProto := range proto.Workloads {
		labelMatcher, err := selectors.MMExprFromLabelMatcher(workloadProto.GetLabelMatcher())
		if err != nil {
			return nil, err
		}

		wm := &workloadMatcher{
			workloadIndex: workloadIndex,
			workloadProto: workloadProto,
		}
		err = mm.AddEntry(workloadIndex, labelMatcher, wm.matchCallback)
		if err != nil {
			return nil, err
		}
	}

	ws := &Scheduler{
		proto:                proto,
		registry:             registry,
		workloadMultiMatcher: mm,
		component:            component,
		metrics:              metrics,
	}

	var wfqMetrics *scheduler.WFQMetrics
	if metrics != nil {
		wfqMetrics = metrics.wfqMetrics
	}

	// setup scheduler
	ws.scheduler = scheduler.NewWFQScheduler(clk, tokenManger, wfqMetrics)

	return ws, nil
}

// Decide processes a single flow by load scheduler in a blocking manner.
//
// Context is used to ensure that requests are not scheduled for longer than its deadline allows.
func (ws *Scheduler) Decide(ctx context.Context, labels map[string]string) *flowcontrolv1.LimiterDecision {
	var matchedWorkloadParametersProto *policylangv1.Scheduler_Workload_Parameters
	var matchedWorkloadIndex string
	// match labels against ws.workloadMultiMatcher
	mmr := ws.workloadMultiMatcher.Match(multimatcher.Labels(labels))
	// if at least one match, return workload with lowest index
	if len(mmr.matchedWorkloads) > 0 {
		// select the smallest workloadIndex
		smallestWorkloadIndex := math.MaxInt32
		for workloadIndex := range mmr.matchedWorkloads {
			if workloadIndex < smallestWorkloadIndex {
				smallestWorkloadIndex = workloadIndex
			}
		}
		matchedWorkload := mmr.matchedWorkloads[smallestWorkloadIndex]
		matchedWorkloadParametersProto = matchedWorkload.GetParameters()
		if matchedWorkload.GetName() != "" {
			matchedWorkloadIndex = matchedWorkload.GetName()
		} else {
			matchedWorkloadIndex = strconv.Itoa(smallestWorkloadIndex)
		}
	} else {
		// no match, return default workload
		matchedWorkloadParametersProto = ws.proto.DefaultWorkloadParameters
		matchedWorkloadIndex = metrics.DefaultWorkloadIndex
	}

	fairnessLabel := "workload:" + matchedWorkloadIndex

	if val, ok := labels[matchedWorkloadParametersProto.FairnessKey]; ok {
		fairnessLabel = fairnessLabel + "," + val
	}

	tokens := uint64(1)
	// Precedence order (lowest to highest):
	// 1. Estimated Tokens
	// 2. Workload tokens
	// 3. Label tokens
	if tokensEstimated, ok := ws.GetEstimatedTokens(matchedWorkloadIndex); ok {
		tokens = tokensEstimated
	}

	if matchedWorkloadParametersProto.Tokens != 0 {
		tokens = matchedWorkloadParametersProto.Tokens
	}

	if ws.proto.TokensLabelKey != "" {
		if val, ok := labels[ws.proto.TokensLabelKey]; ok {
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
		timeout := clientTimeout - ws.proto.DecisionDeadlineMargin.AsDuration()
		if timeout < 0 {
			// we will still schedule the request and it will get
			// dropped if it doesn't get the tokens immediately.
			timeout = 0
		}
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		reqCtx = timeoutCtx
	}

	req := scheduler.NewRequest(fairnessLabel, tokens, uint8(matchedWorkloadParametersProto.Priority))

	accepted := ws.scheduler.Schedule(reqCtx, req)

	tokensConsumed := uint64(0)
	if accepted {
		tokensConsumed = req.Tokens
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:  ws.component.GetPolicyName(),
		PolicyHash:  ws.component.GetPolicyHash(),
		ComponentId: ws.component.GetComponentId(),
		Dropped:     !accepted,
		Details: &flowcontrolv1.LimiterDecision_LoadSchedulerInfo{
			LoadSchedulerInfo: &flowcontrolv1.LimiterDecision_SchedulerInfo{
				WorkloadIndex:  matchedWorkloadIndex,
				TokensConsumed: tokensConsumed,
			},
		},
	}
}

// Revert reverts the decision made by the limiter.
func (ws *Scheduler) Revert(ctx context.Context, labels map[string]string, decision *flowcontrolv1.LimiterDecision) {
	if lsDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_LoadSchedulerInfo); ok {
		tokens := lsDecision.LoadSchedulerInfo.TokensConsumed
		if tokens > 0 {
			ws.scheduler.Revert(ctx, tokens)
		}
	}
}

// GetLatencyObserver returns histogram for specific workload.
func (ws *Scheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	return ws.metrics.wsFactory.GetLatencyObserver(labels)
}

// GetRequestCounter returns request counter for specific workload.
func (ws *Scheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	return ws.metrics.wsFactory.GetRequestCounter(labels)
}

// GetEstimatedTokens returns estimated tokens for specific workload.
func (ws *Scheduler) GetEstimatedTokens(workloadIndex string) (uint64, bool) {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	val, ok := ws.tokensByWorkloadIndex[workloadIndex]
	return val, ok
}

// SetEstimatedTokens sets estimated tokens for specific workload.
func (ws *Scheduler) SetEstimatedTokens(tokensByWorkloadIndex map[string]uint64) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	ws.tokensByWorkloadIndex = tokensByWorkloadIndex
}

// Info returns information about the scheduler.
func (ws *Scheduler) Info() (time.Time, int) {
	return ws.scheduler.Info()
}

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	matchedWorkloads map[int]*policylangv1.Scheduler_Workload
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[int, multiMatchResult]

type workloadMatcher struct {
	workloadProto *policylangv1.Scheduler_Workload
	workloadIndex int
}

func (wm *workloadMatcher) matchCallback(mmr multiMatchResult) multiMatchResult {
	// mmr.matchedWorkloads is nil on first match.
	if mmr.matchedWorkloads == nil {
		mmr.matchedWorkloads = make(map[int]*policylangv1.Scheduler_Workload)
	}
	mmr.matchedWorkloads[wm.workloadIndex] = wm.workloadProto
	return mmr
}

// SanitizeSchedulerProto sanitizes the scheduler proto.
func SanitizeSchedulerProto(proto *policylangv1.Scheduler) (*policylangv1.Scheduler, error) {
	if proto == nil {
		p := &policylangv1.Scheduler{}
		config.SetDefaults(p)
		proto = p
	}

	// default workload params is not a required param so it can be nil
	if proto.DefaultWorkloadParameters == nil {
		p := &policylangv1.Scheduler_Workload_Parameters{}
		config.SetDefaults(p)
		proto.DefaultWorkloadParameters = p
	}

	workloadNames := make(map[string]bool)
	workloadNames[metrics.DefaultWorkloadIndex] = true

	// Loop through the workloads
	for workloadIndex, workloadProto := range proto.Workloads {
		workloadIndexStr := strconv.Itoa(workloadIndex)
		workloadNames[workloadIndexStr] = true
		if workloadProto.GetName() != "" {
			if workloadNames[workloadProto.GetName()] {
				return nil, fmt.Errorf("duplicate workload name %s at %d", workloadProto.Name, workloadIndex)
			}
			workloadNames[workloadProto.Name] = true
		}

		if workloadProto.GetParameters() == nil {
			p := &policylangv1.Scheduler_Workload_Parameters{}
			config.SetDefaults(p)
			workloadProto.Parameters = p
		}
	}

	return proto, nil
}
