package metrics

const (
	// METRIC NAMES.

	// SignalReadingMetricName - used in circuit metrics.
	SignalReadingMetricName = "signal_reading"
	// FluxMeterMetricName name of fluxmeter metrics.
	FluxMeterMetricName = "flux_meter"
	// WorkloadLatencyMetricName - metric used for grouping latencies per workload.
	WorkloadLatencyMetricName = "workload_latency_ms"
	// RequestCounterMetricName - metric from http server.
	RequestCounterMetricName = "http_server_request_counter"
	// ErrorCountMetricName - metric from http server.
	ErrorCountMetricName = "http_server_error_counter"
	// LatencyHistogramMetricName - metric from http server.
	LatencyHistogramMetricName = "http_server_request_latency_seconds"
	// AcceptedConcurrencyMetricName - metric for measuring latencies of accepted requests.
	AcceptedConcurrencyMetricName = "accepted_concurrency"
	// IncomingConcurrencyMetricName - metric for measuring latencies of all incoming requests.
	IncomingConcurrencyMetricName = "incoming_concurrency"
	// WorkloadLatencySumMetricName - metric from workload histogram.
	WorkloadLatencySumMetricName = "workload_latency_ms_sum"
	// WorkloadLatencyCountMetricName - metric from workload histogram.
	WorkloadLatencyCountMetricName = "workload_latency_ms_count"
	// WFQFlowsMetricName - weighted fair queuing number of flows gauge.
	WFQFlowsMetricName = "wfq_flows"
	// WFQRequestsMetricName - weighted fair queuing number of requests gauge.
	WFQRequestsMetricName = "wfq_requests"
	// FlowControlCheckRequestsMetricName - counter for Check requests for flowcontrol.
	FlowControlCheckRequestsMetricName = "flowcontrol_check_requests_total"
	// FlowControlCheckDecisionsMetricName - counter for Check requests per decision type.
	FlowControlCheckDecisionsMetricName = "flowcontrol_check_decisions_total"
	// FlowControlCheckErrorReasonMetricName - metric for error reason on FCS Check requests.
	FlowControlCheckErrorReasonMetricName = "flowcontrol_check_error_reason_total"
	// FlowControlCheckRejectReasonMetricName - metric for reject reason on FCS Check requests.
	FlowControlCheckRejectReasonMetricName = "flowcontrol_check_reject_reason_total"
	// TokenBucketMetricName - a gauge that tracks the load shed factor.
	TokenBucketMetricName = "token_bucket_lsf"
	// TokenBucketFillRateMetricName - a gauge that tracks the fill rate of token bucket.
	TokenBucketFillRateMetricName = "token_bucket_bucket_fill_rate"
	// TokenBucketCapacityMetricName - a gauge that tracks the capacity of token bucket.
	TokenBucketCapacityMetricName = "token_bucket_bucket_capacity"
	// TokenBucketAvailableMetricName - a gauge that tracks the number of tokens available in token bucket.
	TokenBucketAvailableMetricName = "token_bucket_available_tokens"
	// GroupJobRegisteredMetricName - current number of group job registered.
	GroupJobRegisteredMetricName = "group_job_registered_number"
	// GroupJobScheduledMetricName - current number of group job scheduled.
	GroupJobScheduledMetricName = "group_job_scheduled_number"
	// GroupJobCompletedMetricName - total number of group job completed.
	GroupJobCompletedMetricName = "group_job_completed_total"
	// GroupJobLatencyMetricName - the latency of the group jobs.
	GroupJobLatencyMetricName = "group_job_latency_seconds"

	// PROMETHEUS LABELS.

	// PolicyNameLabel - label used in prometheus.
	PolicyNameLabel = "policy_name"
	// PolicyHashLabel - label used in prometheus.
	PolicyHashLabel = "policy_hash"
	// ComponentIndexLabel - index of component in circuit label.
	ComponentIndexLabel = "component_index"
	// DecisionTypeLabel - label for decision type dropped or accepted.
	DecisionTypeLabel = "decision_type"
	// WorkloadIndexLabel - label for choosing correct workload.
	WorkloadIndexLabel = "workload_index"
	// SignalNameLabel - label for saving circuit signal metrics.
	SignalNameLabel = "signal_name"
	// FluxMeterNameLabel - specifying flux meter's name.
	FluxMeterNameLabel = "flux_meter_name"
	// StatusCodeLabel - http/library return code.
	StatusCodeLabel = "status_code"
	// MethodLabel - label from http method.
	MethodLabel = "method"
	// ResponseStatusCodeLabel - label from response status code.
	ResponseStatusCodeLabel = "response_status_code"
	// FlowControlCheckDecisionTypeLabel - label for decision type dropped or accepted.
	FlowControlCheckDecisionTypeLabel = "flowcontrol_check_decision_type"
	// FlowControlCheckErrorReasonLabel - label for error reason on FCS Check request.
	FlowControlCheckErrorReasonLabel = "flowcontrol_check_error_reason"
	// FlowControlCheckRejectReasonLabel - label for reject reason on FCS Check request.
	FlowControlCheckRejectReasonLabel = "flowcontrol_check_reject_reason"

	// DEFAULTS.

	// DefaultWorkloadIndex - when workload is not specified this value is used.
	DefaultWorkloadIndex = "default"
	// DefaultAgentGroup - default agent group.
	DefaultAgentGroup = "default"
)
