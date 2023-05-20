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

// NewFactory sets up the load scheduler module in the main fx app.
func NewFactory(
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

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	matchedWorkloads map[int]*policylangv1.Scheduler_Workload_Parameters
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
		mmr.matchedWorkloads = make(map[int]*policylangv1.Scheduler_Workload_Parameters)
	}
	mmr.matchedWorkloads[wm.workloadIndex] = wm.workloadProto.GetParameters()
	return mmr
}

// Scheduler implements load scheduler on the flowcontrol side.
type Scheduler struct {
	mutex                 sync.RWMutex
	component             iface.Component
	scheduler             scheduler.Scheduler
	registry              status.Registry
	incomingTokensCounter prometheus.Counter
	acceptedTokensCounter prometheus.Counter
	proto                 *policylangv1.Scheduler
	wsFactory             *Factory
	workloadMultiMatcher  *multiMatcher
	tokensByWorkloadIndex map[string]uint64
	metricLabels          prometheus.Labels
}

// NewScheduler returns fx options for the load scheduler fx app.
func (wsFactory *Factory) NewScheduler(
	registry status.Registry,
	proto *policylangv1.Scheduler,
	component iface.Component,
	tokenManger scheduler.TokenManager,
	clock clockwork.Clock,
	metricLabels prometheus.Labels,
) (*Scheduler, error) {
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
		wsFactory:            wsFactory,
		workloadMultiMatcher: mm,
		metricLabels:         metricLabels,
		component:            component,
	}

	wfqFlowsGauge, err := wsFactory.wfqFlowsGaugeVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get wfq flows gauge", err)
	}

	wfqRequestsGauge, err := wsFactory.wfqRequestsGaugeVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get wfq requests gauge", err)
	}

	wfqMetrics := &scheduler.WFQMetrics{
		FlowsGauge:        wfqFlowsGauge,
		HeapRequestsGauge: wfqRequestsGauge,
	}

	// setup scheduler
	ws.scheduler = scheduler.NewWFQScheduler(tokenManger, clock, wfqMetrics)

	ws.incomingTokensCounter, err = wsFactory.incomingTokensCounterVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, err
	}
	ws.acceptedTokensCounter, err = wsFactory.acceptedTokensCounterVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// Close closes the scheduler.
func (ws *Scheduler) Close() error {
	var merr error

	// Remove metrics from metric vectors
	deleted := ws.wsFactory.wfqFlowsGaugeVec.Delete(ws.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete wfq_flows gauge from its metric vector"))
	}
	deleted = ws.wsFactory.wfqRequestsGaugeVec.Delete(ws.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete wfq_requests gauge from its metric vector"))
	}
	deleted = ws.wsFactory.incomingTokensCounterVec.Delete(ws.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete incoming_tokens_total counter from its metric vector"))
	}
	deleted = ws.wsFactory.acceptedTokensCounterVec.Delete(ws.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete accepted_tokens_total counter from its metric vector"))
	}
	deletedCount := ws.wsFactory.workloadLatencySummaryVec.DeletePartialMatch(ws.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_latency_ms summary from its metric vector. No traffic to generate metrics?")
	}
	deletedCount = ws.wsFactory.workloadCounterVec.DeletePartialMatch(ws.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_requests_total counter from its metric vector. No traffic to generate metrics?")
	}

	ws.registry.SetStatus(status.NewStatus(nil, merr))
	return merr
}

// Decide processes a single flow by load scheduler in a blocking manner.
//
// Context is used to ensure that requests are not scheduled for longer than its deadline allows.
func (ws *Scheduler) Decide(ctx context.Context,
	labels map[string]string,
) *flowcontrolv1.LimiterDecision {
	var matchedWorkloadProto *policylangv1.Scheduler_Workload_Parameters
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
		matchedWorkloadProto = mmr.matchedWorkloads[smallestWorkloadIndex]
		matchedWorkloadIndex = strconv.Itoa(smallestWorkloadIndex)
	} else {
		// no match, return default workload
		matchedWorkloadProto = ws.proto.DefaultWorkloadParameters
		matchedWorkloadIndex = metrics.DefaultWorkloadIndex
	}

	fairnessLabel := "workload:" + matchedWorkloadIndex

	if val, ok := labels[matchedWorkloadProto.FairnessKey]; ok {
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

	if matchedWorkloadProto.Tokens != 0 {
		tokens = matchedWorkloadProto.Tokens
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

	req := scheduler.Request{
		FairnessLabel: fairnessLabel,
		Priority:      uint8(matchedWorkloadProto.Priority),
		Tokens:        tokens,
	}

	accepted := ws.scheduler.Schedule(reqCtx, req)

	tokensConsumed := uint64(0)
	if accepted {
		tokensConsumed = req.Tokens
	}

	// update load scheduler metrics and decisionType
	ws.incomingTokensCounter.Add(float64(req.Tokens) / 1000)

	if accepted {
		ws.acceptedTokensCounter.Add(float64(req.Tokens) / 1000)
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:  ws.component.GetPolicyName(),
		PolicyHash:  ws.component.GetPolicyHash(),
		ComponentId: ws.component.GetComponentId(),
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
func (ws *Scheduler) Revert(labels map[string]string, decision *flowcontrolv1.LimiterDecision) {
	if lsDecision, ok := decision.GetDetails().(*flowcontrolv1.LimiterDecision_LoadSchedulerInfo_); ok {
		tokens := lsDecision.LoadSchedulerInfo.TokensConsumed
		if tokens > 0 {
			ws.scheduler.Revert(tokens)
		}
	}
}

// GetLatencyObserver returns histogram for specific workload.
func (ws *Scheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	latencySummary, err := ws.wsFactory.workloadLatencySummaryVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return nil
	}

	return latencySummary
}

// GetRequestCounter returns request counter for specific workload.
func (ws *Scheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	counter, err := ws.wsFactory.workloadCounterVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting counter")
		return nil
	}

	return counter
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
