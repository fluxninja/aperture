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

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/labelstatus"
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
	incomingTokensCounterVec *prometheus.CounterVec
	acceptedTokensCounterVec *prometheus.CounterVec
	rejectedTokensCounterVec *prometheus.CounterVec

	requestInQueueDurationSummaryVec *prometheus.SummaryVec

	workloadLatencySummaryVec *prometheus.SummaryVec
	workloadCounterVec        *prometheus.CounterVec

	workloadPreemptedTokensSummaryVec *prometheus.SummaryVec
	workloadDelayedTokensSummaryVec   *prometheus.SummaryVec
	workloadOnTimeCounterVec          *prometheus.CounterVec

	fairnessPreemptedTokensSummaryVec *prometheus.SummaryVec
	fairnessDelayedTokensSummaryVec   *prometheus.SummaryVec
	fairnessOnTimeCounterVec          *prometheus.CounterVec
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
	wsFactory.rejectedTokensCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.RejectedTokensMetricName,
			Help: "A counter measuring work rejected by Scheduler",
		},
		MetricLabelKeys,
	)
	wsFactory.requestInQueueDurationSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.RequestInQueueDurationMetricName,
		Help: "Duration of requests scheduled in Queue",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

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

	wsFactory.workloadPreemptedTokensSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.WorkloadPreemptedTokensMetricName,
		Help: "Number of tokens a request was preempted, measured end-to-end in the scheduler across all workloads.",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	wsFactory.workloadDelayedTokensSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.WorkloadDelayedTokensMetricName,
		Help: "Number of tokens a request was delayed by, measured end-to-end in the scheduler across all workloads.",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	wsFactory.workloadOnTimeCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.WorkloadOnTimeMetricName,
		Help: "Counter of workload requests that were on time, measured end-to-end in the scheduler across all workloads.",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	wsFactory.fairnessPreemptedTokensSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.FairnessPreemptedTokensMetricName,
		Help: "Number of tokens a request was preempted, measured at fairness queues within the same workload.",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	wsFactory.fairnessDelayedTokensSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.FairnessDelayedTokensMetricName,
		Help: "Number of tokens a request was delayed by, measured at fairness queues within the same workload.",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	wsFactory.fairnessOnTimeCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.FairnessOnTimeMetricName,
		Help: "Counter of workload requests that were on time, measured at fairness queues within the same workload.",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	})

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var merr error

			err := prometheusRegistry.Register(wsFactory.incomingTokensCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.acceptedTokensCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.rejectedTokensCounterVec)
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
			err = prometheusRegistry.Register(wsFactory.requestInQueueDurationSummaryVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.workloadPreemptedTokensSummaryVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.workloadDelayedTokensSummaryVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.workloadOnTimeCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.fairnessPreemptedTokensSummaryVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.fairnessDelayedTokensSummaryVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}
			err = prometheusRegistry.Register(wsFactory.fairnessOnTimeCounterVec)
			if err != nil {
				merr = multierr.Append(merr, err)
			}

			return merr
		},
		OnStop: func(_ context.Context) error {
			var merr error

			if !prometheusRegistry.Unregister(wsFactory.incomingTokensCounterVec) {
				err := fmt.Errorf("failed to unregister incoming_tokens_total metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.acceptedTokensCounterVec) {
				err := fmt.Errorf("failed to unregister accepted_tokens_total metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.rejectedTokensCounterVec) {
				err := fmt.Errorf("failed to unregister rejected_tokens_total metric")
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
			if !prometheusRegistry.Unregister(wsFactory.requestInQueueDurationSummaryVec) {
				err := fmt.Errorf("failed to unregister request_in_queue_duration_ms metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.workloadPreemptedTokensSummaryVec) {
				err := fmt.Errorf("failed to unregister workload_preempted_tokens metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.workloadDelayedTokensSummaryVec) {
				err := fmt.Errorf("failed to unregister workload_delayed_tokens metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.workloadOnTimeCounterVec) {
				err := fmt.Errorf("failed to unregister workload_on_time_total metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.fairnessPreemptedTokensSummaryVec) {
				err := fmt.Errorf("failed to unregister fairness_preempted_tokens metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.fairnessDelayedTokensSummaryVec) {
				err := fmt.Errorf("failed to unregister fairness_delayed_tokens metric")
				merr = multierr.Append(merr, err)
			}
			if !prometheusRegistry.Unregister(wsFactory.fairnessOnTimeCounterVec) {
				err := fmt.Errorf("failed to unregister fairness_on_time_total metric")
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
	incomingTokensCounter, err := wsFactory.incomingTokensCounterVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, err
	}

	acceptedTokensCounter, err := wsFactory.acceptedTokensCounterVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, err
	}

	rejectedTokensCounter, err := wsFactory.rejectedTokensCounterVec.GetMetricWith(metricLabels)
	if err != nil {
		return nil, err
	}

	wfqMetrics := &scheduler.WFQMetrics{
		IncomingTokensCounter:          incomingTokensCounter,
		AcceptedTokensCounter:          acceptedTokensCounter,
		RejectedTokensCounter:          rejectedTokensCounter,
		RequestInQueueDurationSummary:  wsFactory.requestInQueueDurationSummaryVec,
		WorkloadPreemptedTokensSummary: wsFactory.workloadPreemptedTokensSummaryVec,
		WorkloadDelayedTokensSummary:   wsFactory.workloadDelayedTokensSummaryVec,
		WorkloadOnTimeCounter:          wsFactory.workloadOnTimeCounterVec,
		FairnessPreemptedTokensSummary: wsFactory.fairnessPreemptedTokensSummaryVec,
		FairnessDelayedTokensSummary:   wsFactory.fairnessDelayedTokensSummaryVec,
		FairnessOnTimeCounter:          wsFactory.fairnessOnTimeCounterVec,
	}

	return &SchedulerMetrics{
		wfqMetrics:   wfqMetrics,
		metricLabels: metricLabels,
		wsFactory:    wsFactory,
	}, nil
}

// Delete removes all metrics from metric vectors.
func (sm *SchedulerMetrics) Delete() error {
	var merr error

	deleted := sm.wsFactory.incomingTokensCounterVec.Delete(sm.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete incoming_tokens_total counter from its metric vector"))
	}
	deleted = sm.wsFactory.acceptedTokensCounterVec.Delete(sm.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete accepted_tokens_total counter from its metric vector"))
	}
	deleted = sm.wsFactory.rejectedTokensCounterVec.Delete(sm.metricLabels)
	if !deleted {
		merr = multierr.Append(merr, errors.New("failed to delete rejected_tokens_total counter from its metric vector"))
	}
	deletedCount := sm.wsFactory.workloadLatencySummaryVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_latency_ms summary from its metric vector. No traffic to generate metrics?")
	}
	deletedCount = sm.wsFactory.workloadCounterVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_requests_total counter from its metric vector. No traffic to generate metrics?")
	}
	deletedCount = sm.wsFactory.requestInQueueDurationSummaryVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete request_in_queue_duration_ms summary from its metric vector. No traffic to generate metrics?")
	}
	deletedCount = sm.wsFactory.workloadPreemptedTokensSummaryVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_preempted_tokens summary from its metric vector.")
	}
	deletedCount = sm.wsFactory.workloadDelayedTokensSummaryVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_delayed_tokens summary from its metric vector.")
	}
	deletedCount = sm.wsFactory.workloadOnTimeCounterVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete workload_on_time_total counter from its metric vector.")
	}
	deletedCount = sm.wsFactory.fairnessPreemptedTokensSummaryVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete fairness_preempted_tokens summary from its metric vector.")
	}
	deletedCount = sm.wsFactory.fairnessDelayedTokensSummaryVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete fairness_delayed_tokens summary from its metric vector.")
	}
	deletedCount = sm.wsFactory.fairnessOnTimeCounterVec.DeletePartialMatch(sm.metricLabels)
	if deletedCount == 0 {
		log.Warn().Msg("Could not delete fairness_on_time_total counter from its metric vector.")
	}
	return merr
}

func (sm *SchedulerMetrics) appendWorkloadLabel(workloadLabel string) prometheus.Labels {
	baseMetricsLabels := sm.metricLabels
	metricsLabels := make(prometheus.Labels, len(baseMetricsLabels)+1)
	metricsLabels[metrics.WorkloadIndexLabel] = workloadLabel
	for k, v := range baseMetricsLabels {
		metricsLabels[k] = v
	}

	return metricsLabels
}

// Scheduler implements load scheduler on the flowcontrol side.
type Scheduler struct {
	component              iface.Component
	scheduler              scheduler.Scheduler
	registry               status.Registry
	proto                  *policylangv1.Scheduler
	defaultWorkload        *workload
	workloadMultiMatcher   *multiMatcher
	tokensByWorkloadIndex  map[string]float64
	metrics                *SchedulerMetrics
	mutex                  sync.RWMutex
	tokensLabelKeyStatus   *labelstatus.LabelStatus
	priorityLabelKeyStatus *labelstatus.LabelStatus
	workloadLabelKeyStatus *labelstatus.LabelStatus
	fairnessLabelKeyStatus *labelstatus.LabelStatus
}

// NewScheduler returns fx options for the load scheduler fx app.
func (wsFactory *Factory) NewScheduler(
	clk clockwork.Clock,
	registry status.Registry,
	proto *policylangv1.Scheduler,
	component iface.Component,
	tokenManger scheduler.TokenManager,
	schedulerMetrics *SchedulerMetrics,
	tokensLabelKeyStatus *labelstatus.LabelStatus,
	priorityLabelKeyStatus *labelstatus.LabelStatus,
	workloadLabelKeyStatus *labelstatus.LabelStatus,
	fairnessLabelKeyStatus *labelstatus.LabelStatus,
) (*Scheduler, error) {
	initPreemptMetrics := func(workloadLabel string) error {
		if schedulerMetrics == nil {
			return nil
		}
		workloadLabels := schedulerMetrics.appendWorkloadLabel(workloadLabel)
		var err error
		_, err = schedulerMetrics.wsFactory.workloadPreemptedTokensSummaryVec.GetMetricWith(workloadLabels)
		if err != nil {
			return fmt.Errorf("%w: failed to get workload_preempted_tokens summary", err)
		}
		_, err = schedulerMetrics.wsFactory.workloadDelayedTokensSummaryVec.GetMetricWith(workloadLabels)
		if err != nil {
			return fmt.Errorf("%w: failed to get workload_delayed_tokens summary", err)
		}
		_, err = schedulerMetrics.wsFactory.workloadOnTimeCounterVec.GetMetricWith(workloadLabels)
		if err != nil {
			return fmt.Errorf("%w: failed to get workload_on_time_total counter", err)
		}
		_, err = schedulerMetrics.wsFactory.fairnessPreemptedTokensSummaryVec.GetMetricWith(workloadLabels)
		if err != nil {
			return fmt.Errorf("%w: failed to get fairness_preempted_tokens summary", err)
		}
		_, err = schedulerMetrics.wsFactory.fairnessDelayedTokensSummaryVec.GetMetricWith(workloadLabels)
		if err != nil {
			return fmt.Errorf("%w: failed to get fairness_delayed_tokens summary", err)
		}
		_, err = schedulerMetrics.wsFactory.fairnessOnTimeCounterVec.GetMetricWith(workloadLabels)
		if err != nil {
			return fmt.Errorf("%w: failed to get fairness_on_time_total counter", err)
		}
		return nil
	}

	mm := multimatcher.New[int, multiMatchResult]()
	for workloadIndex, workloadProto := range proto.Workloads {
		labelMatcher, err := selectors.MMExprFromLabelMatcher(workloadProto.GetLabelMatcher())
		if err != nil {
			return nil, err
		}
		wm := &workloadMatcher{
			workloadIndex: workloadIndex,
			workload: &workload{
				proto:    workloadProto,
				priority: workloadProto.Parameters.Priority,
			},
		}
		err = mm.AddEntry(workloadIndex, labelMatcher, wm.matchCallback)
		if err != nil {
			return nil, err
		}
		err = initPreemptMetrics(getWorkloadLabel(workloadIndex, workloadProto))
		if err != nil {
			return nil, err
		}
	}
	// default workload
	err := initPreemptMetrics(metrics.DefaultWorkloadIndex)
	if err != nil {
		return nil, err
	}

	ws := &Scheduler{
		proto: proto,
		defaultWorkload: &workload{
			priority: proto.DefaultWorkloadParameters.Priority,
			proto: &policylangv1.Scheduler_Workload{
				Parameters: proto.DefaultWorkloadParameters,
				Name:       metrics.DefaultWorkloadIndex,
			},
		},
		registry:               registry,
		workloadMultiMatcher:   mm,
		component:              component,
		metrics:                schedulerMetrics,
		tokensLabelKeyStatus:   tokensLabelKeyStatus,
		priorityLabelKeyStatus: priorityLabelKeyStatus,
		workloadLabelKeyStatus: workloadLabelKeyStatus,
		fairnessLabelKeyStatus: fairnessLabelKeyStatus,
	}

	var wfqMetrics *scheduler.WFQMetrics
	if schedulerMetrics != nil {
		wfqMetrics = schedulerMetrics.wfqMetrics
	}

	// setup scheduler
	ws.scheduler = scheduler.NewWFQScheduler(clk, tokenManger, wfqMetrics, schedulerMetrics.metricLabels)

	return ws, nil
}

// Decide processes a single flow by load scheduler in a blocking manner.
// Context is used to ensure that requests are not scheduled for longer than its deadline allows.
func (s *Scheduler) Decide(ctx context.Context, labels labels.Labels) (*flowcontrolv1.LimiterDecision, string) {
	var matchedWorkloadParametersProto *policylangv1.Scheduler_Workload_Parameters
	var invPriority float64
	var priority float64
	var matchedWorkloadLabel string
	// match labels against ws.workloadMultiMatcher
	mmr := s.workloadMultiMatcher.Match(labels)
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
		priority = matchedWorkload.priority
		invPriority = 1 / matchedWorkload.priority
		matchedWorkloadParametersProto = matchedWorkload.proto.GetParameters()
		matchedWorkloadLabel = getWorkloadLabel(smallestWorkloadIndex, matchedWorkload.proto)
	} else {
		// no match, return default workload
		priority = s.defaultWorkload.priority
		invPriority = 1 / s.defaultWorkload.priority
		matchedWorkloadParametersProto = s.defaultWorkload.proto.Parameters
		matchedWorkloadLabel = s.defaultWorkload.proto.Name
	}

	tokens := float64(1)
	// Precedence order:
	// 1. Label tokens
	// 2. Estimated Tokens
	// 3. Workload tokens
	if matchedWorkloadParametersProto.GetTokens() != 0 {
		tokens = matchedWorkloadParametersProto.GetTokens()
	}

	if s.proto.WorkloadLabelKey != "" {
		val, ok := labels.Get(s.proto.WorkloadLabelKey)
		if ok {
			matchedWorkloadLabel = val
		} else {
			s.workloadLabelKeyStatus.SetMissing()
		}
	}

	if estimatedTokens, ok := s.GetEstimatedTokens(matchedWorkloadLabel); ok {
		tokens = estimatedTokens
	}

	if s.proto.TokensLabelKey != "" {
		val, ok := labels.Get(s.proto.TokensLabelKey)
		if ok {
			if parsedTokens, err := strconv.ParseFloat(val, 64); err == nil {
				tokens = parsedTokens
			} else {
				s.tokensLabelKeyStatus.SetMissing()
			}
		}
	}

	if s.proto.PriorityLabelKey != "" {
		val, ok := labels.Get(s.proto.PriorityLabelKey)
		if ok {
			if parsedPriority, err := strconv.ParseFloat(val, 64); err == nil {
				if parsedPriority > 0 {
					priority = parsedPriority
					invPriority = 1 / parsedPriority
				}
			}
		} else {
			s.priorityLabelKeyStatus.SetMissing()
		}
	}

	reqCtx := ctx

	var matchedWorkloadTimeout time.Duration
	hasWorkloadTimeout := false
	if matchedWorkloadParametersProto.QueueTimeout != nil {
		matchedWorkloadTimeout = matchedWorkloadParametersProto.QueueTimeout.AsDuration()
		hasWorkloadTimeout = true
	}

	clientDeadline, hasClientDeadline := ctx.Deadline()
	if hasClientDeadline {
		// The clientDeadline is calculated based on client's timeout, passed
		// as grpc-timeout. Our goal is for the response to be received by the
		// client before its deadline passes (otherwise we risk fail-open on
		// timeout). To allow some headroom for transmitting the response to
		// the client, we set an "internal" deadline to a bit before client's
		// deadline, subtracting the configured margin.
		clientTimeout := time.Until(clientDeadline)
		timeout := clientTimeout - s.proto.DecisionDeadlineMargin.AsDuration()
		if timeout < 0 {
			// we will still schedule the request and it will get
			// dropped if it doesn't get the tokens immediately.
			timeout = 0
		}

		// find the minimum of matchedWorkloadTimeout and client's timeout
		if hasWorkloadTimeout && matchedWorkloadTimeout < timeout {
			timeout = matchedWorkloadTimeout
		}
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		reqCtx = timeoutCtx
	} else if hasWorkloadTimeout {
		// If there is no client deadline but there is a workload timeout, we create a new context with the workload timeout.
		timeoutCtx, cancel := context.WithTimeout(ctx, matchedWorkloadTimeout)
		defer cancel()
		reqCtx = timeoutCtx
	}

	var fairnessLabel string
	if s.proto.FairnessLabelKey != "" {
		val, ok := labels.Get(s.proto.FairnessLabelKey)
		if ok {
			fairnessLabel = val
		} else {
			s.fairnessLabelKeyStatus.SetMissing()
		}
	}

	req := scheduler.NewRequest(matchedWorkloadLabel, fairnessLabel, tokens, invPriority)

	accepted, remaining, current, reqID := s.scheduler.Schedule(reqCtx, req)

	tokensConsumed := float64(0)
	if accepted {
		tokensConsumed = req.Tokens
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:               s.component.GetPolicyName(),
		PolicyHash:               s.component.GetPolicyHash(),
		ComponentId:              s.component.GetComponentId(),
		Dropped:                  !accepted,
		DeniedResponseStatusCode: s.proto.GetDeniedResponseStatusCode(),
		Details: &flowcontrolv1.LimiterDecision_LoadSchedulerInfo{
			LoadSchedulerInfo: &flowcontrolv1.LimiterDecision_SchedulerInfo{
				WorkloadIndex: matchedWorkloadLabel,
				TokensInfo: &flowcontrolv1.LimiterDecision_TokensInfo{
					Consumed:  tokensConsumed,
					Remaining: remaining,
					Current:   current,
				},
				Priority: priority,
			},
		},
	}, reqID
}

// GetLatencyObserver returns histogram for specific workload.
func (s *Scheduler) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	return s.metrics.wsFactory.GetLatencyObserver(labels)
}

// GetRequestCounter returns request counter for specific workload.
func (s *Scheduler) GetRequestCounter(labels map[string]string) prometheus.Counter {
	return s.metrics.wsFactory.GetRequestCounter(labels)
}

// GetEstimatedTokens returns estimated tokens for specific workload.
func (s *Scheduler) GetEstimatedTokens(workloadIndex string) (float64, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.tokensByWorkloadIndex[workloadIndex]
	return val, ok
}

// SetEstimatedTokens sets estimated tokens for specific workload.
func (s *Scheduler) SetEstimatedTokens(tokensByWorkloadIndex map[string]float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tokensByWorkloadIndex = tokensByWorkloadIndex
}

// Info returns information about the scheduler.
func (s *Scheduler) Info() (time.Time, int) {
	return s.scheduler.Info()
}

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	matchedWorkloads map[int]*workload
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[int, multiMatchResult]

type workload struct {
	proto    *policylangv1.Scheduler_Workload
	priority float64
}

type workloadMatcher struct {
	workload      *workload
	workloadIndex int
}

func (wm *workloadMatcher) matchCallback(mmr multiMatchResult) multiMatchResult {
	// mmr.matchedWorkloads is nil on first match.
	if mmr.matchedWorkloads == nil {
		mmr.matchedWorkloads = make(map[int]*workload)
	}

	mmr.matchedWorkloads[wm.workloadIndex] = wm.workload
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

func getWorkloadLabel(workloadIndex int, workloadProto *policylangv1.Scheduler_Workload) string {
	var workloadLabel string
	if workloadProto.GetName() != "" {
		workloadLabel = workloadProto.GetName()
	} else {
		workloadLabel = strconv.Itoa(workloadIndex)
	}
	return workloadLabel
}
