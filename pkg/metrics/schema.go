package metrics

const (
	// METRIC NAMES.

	// HTTP metrics.

	// HTTPRequestMetricName - metric from http server.
	HTTPRequestMetricName = "http_requests_total"
	// HTTPRequestLatencyMetricName - metric from http server.
	HTTPRequestLatencyMetricName = "http_requests_latency_ms"
	// HTTPErrorMetricName - metric from http server.
	HTTPErrorMetricName = "http_errors_total"

	// Circuit metrics.

	// SignalReadingMetricName - used in circuit metrics.
	SignalReadingMetricName = "signal_reading"
	// FluxMeterMetricName name of fluxmeter metrics.
	FluxMeterMetricName = "flux_meter"
	// RateLimiterCounterMetricName - name of the counter describing times rate limiter was triggered.
	RateLimiterCounterMetricName = "rate_limiter_counter"

	// Workload metrics.

	// WorkloadLatencyMetricName - metric used for grouping latencies per workload.
	WorkloadLatencyMetricName = "workload_latency_ms"
	// WorkloadLatencySumMetricName - metric from workload histogram.
	WorkloadLatencySumMetricName = "workload_latency_ms_sum"
	// WorkloadLatencyCountMetricName - metric from workload histogram.
	WorkloadLatencyCountMetricName = "workload_latency_ms_count"

	// AcceptedConcurrencyMetricName - metric for measuring latencies of accepted requests.
	AcceptedConcurrencyMetricName = "accepted_concurrency_ms"
	// IncomingConcurrencyMetricName - metric for measuring latencies of all incoming requests.
	IncomingConcurrencyMetricName = "incoming_concurrency_ms"

	// WFQFlowsMetricName - weighted fair queuing number of flows gauge.
	WFQFlowsMetricName = "wfq_flows_total"
	// WFQRequestsMetricName - weighted fair queuing number of requests gauge.
	WFQRequestsMetricName = "wfq_requests_total"
	// TokenBucketMetricName - a gauge that tracks the load shed factor.
	TokenBucketMetricName = "token_bucket_lsf_ratio"
	// TokenBucketFillRateMetricName - a gauge that tracks the fill rate of token bucket.
	TokenBucketFillRateMetricName = "token_bucket_fill_rate"
	// TokenBucketCapacityMetricName - a gauge that tracks the capacity of token bucket.
	TokenBucketCapacityMetricName = "token_bucket_capacity_total"
	// TokenBucketAvailableMetricName - a gauge that tracks the number of tokens available in token bucket.
	TokenBucketAvailableMetricName = "token_bucket_available_tokens_total"

	// FlowControlRequestsMetricName - counter for Check requests for flowcontrol.
	FlowControlRequestsMetricName = "flowcontrol_requests_total"
	// FlowControlDecisionsMetricName - counter for Check requests per decision type.
	FlowControlDecisionsMetricName = "flowcontrol_decisions_total"
	// FlowControlErrorReasonsMetricName - metric for error reason on FCS Check requests.
	FlowControlErrorReasonsMetricName = "flowcontrol_error_reasons_total"
	// FlowControlRejectReasonsMetricName - metric for reject reason on FCS Check requests.
	FlowControlRejectReasonsMetricName = "flowcontrol_reject_reasons_total"

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
	// ClassifierIndexLabel - prometheus label specifying clasiffier index.
	ClassifierIndexLabel = "classifier_index"
	// StatusCodeLabel - http status code.
	StatusCodeLabel = "http_status_code"
	// MethodLabel - label from http method.
	MethodLabel = "http_method"
	// HandlerName - name of the http handler. Defaults to 'default'.
	HandlerName = "handler_name"
	// FeatureStatusLabel - feature status.
	FeatureStatusLabel = "feature_status"
	// FeatureStatusOK - feature status OK.
	FeatureStatusOK = "OK"
	// FeatureStatusError - feature status Error.
	FeatureStatusError = "Error"
	// ResponseStatusLabel - response status. A common label to denote OK or Error across all protocols.
	ResponseStatusLabel = "response_status"
	// ResponseStatusOK - response status OK.
	ResponseStatusOK = FeatureStatusOK
	// ResponseStatusError - response status Error.
	ResponseStatusError = FeatureStatusError
	// FlowControlCheckDecisionTypeLabel - label for decision type dropped or accepted.
	FlowControlCheckDecisionTypeLabel = "decision_type"
	// FlowControlCheckErrorReasonLabel - label for error reason on FCS Check request.
	FlowControlCheckErrorReasonLabel = "error_reason"
	// FlowControlCheckRejectReasonLabel - label for reject reason on FCS Check request.
	FlowControlCheckRejectReasonLabel = "reject_reason"

	// DEFAULTS.

	// DefaultWorkloadIndex - when workload is not specified this value is used.
	DefaultWorkloadIndex = "default"
	// DefaultAgentGroup - default agent group.
	DefaultAgentGroup = "default"
)
