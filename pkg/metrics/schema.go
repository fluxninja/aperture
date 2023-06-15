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
	// InvalidSignalReadingsTotalMetricName - counts invalid signal readings.
	InvalidSignalReadingsTotalMetricName = "invalid_signal_readings_total"
	// InvalidFluxMeterTotalMetricName - counts invalid flux meters.
	InvalidFluxMeterTotalMetricName = "invalid_flux_meter_total"
	// RateLimiterCounterTotalMetricName - name of the counter describing times rate limiter was triggered.
	RateLimiterCounterTotalMetricName = "rate_limiter_counter_total"
	// ClassifierCounterTotalMetricName - name of the counter describing times classifier was triggered.
	ClassifierCounterTotalMetricName = "classifier_counter_total"
	// SamplerCounterTotalMetricName - name of the counter describing times sampler was triggered.
	SamplerCounterTotalMetricName = "sampler_counter_total"

	// DistCache metrics scraped from Olric DMaps statistics.

	// DistCacheEntriesTotalMetricName - metric for the total number of entries (including replicas) stored during the life of this instance.
	DistCacheEntriesTotalMetricName = "distcache_entries_total"
	// DistCacheDeleteHitsMetricName - metric for number of deletion requests resulting in an item being removed.
	DistCacheDeleteHitsMetricName = "distcache_delete_hits"
	// DistCacheDeleteMissesMetricName - metric for number of deletion requests for missing keys.
	DistCacheDeleteMissesMetricName = "distcache_delete_misses"
	// DistCacheGetMissesMetricName - metric for number of entries that have been requested and not found.
	DistCacheGetMissesMetricName = "distcache_get_misses"
	// DistCacheGetHitsMetricName - metric for number of entries that have been requested and found present.
	DistCacheGetHitsMetricName = "distcache_get_hits"
	// DistCacheEvictedTotalMetricName - metric for number of entries removed from cache to free memory for new entries.
	DistCacheEvictedTotalMetricName = "distcache_evicted_total"

	// Workload metrics.

	// WorkloadLatencyMetricName - metric used for grouping latencies per workload.
	WorkloadLatencyMetricName = "workload_latency_ms"
	// WorkloadLatencySumMetricName - metric from workload histogram.
	WorkloadLatencySumMetricName = "workload_latency_ms_sum"
	// WorkloadLatencyCountMetricName - metric from workload histogram.
	WorkloadLatencyCountMetricName = "workload_latency_ms_count"
	// WorkloadCounterMetricName - metric used for counting workload requests.
	WorkloadCounterMetricName = "workload_requests_total"

	// IncomingTokensMetricName - total work measured in tokens (auto-tokens are seconds of estimated duration) of all incoming requests.
	IncomingTokensMetricName = "incoming_tokens_total"
	// AcceptedTokensMetricName - total work measured in tokens (auto-tokens are seconds of estimated duration) of all accepted requests.
	AcceptedTokensMetricName = "accepted_tokens_total"

	// WFQFlowsMetricName - weighted fair queuing number of flows gauge.
	WFQFlowsMetricName = "wfq_flows_total"
	// WFQRequestsMetricName - weighted fair queuing number of requests gauge.
	WFQRequestsMetricName = "wfq_requests_total"
	// TokenBucketLMMetricName - a gauge that tracks the load multiplier.
	TokenBucketLMMetricName = "token_bucket_lm_ratio"
	// TokenBucketFillRateMetricName - a gauge that tracks the fill rate of token bucket.
	TokenBucketFillRateMetricName = "token_bucket_fill_rate"
	// TokenBucketCapacityMetricName - a gauge that tracks the capacity of token bucket.
	TokenBucketCapacityMetricName = "token_bucket_capacity_total"
	// TokenBucketAvailableMetricName - a gauge that tracks the number of tokens available in token bucket.
	TokenBucketAvailableMetricName = "token_bucket_available_tokens_total"

	// ServiceLookupsMetricName - counter for IP to services lookups.
	ServiceLookupsMetricName = "service_lookups_total"

	// FlowControlRequestsMetricName - counter for Check requests for flowcontrol.
	FlowControlRequestsMetricName = "flowcontrol_requests_total"
	// FlowControlDecisionsMetricName - counter for Check requests per decision type.
	FlowControlDecisionsMetricName = "flowcontrol_decisions_total"
	// FlowControlRejectReasonsMetricName - metric for reject reason on FCS Check requests.
	FlowControlRejectReasonsMetricName = "flowcontrol_reject_reasons_total"

	// OTel metrics.

	// RollupMetricName - logs rollup histogram.
	RollupMetricName = "rollup"

	// PROMETHEUS LABELS.

	// InstanceLabel used to identify the host name on which an Aperture process is running.
	InstanceLabel = "instance"
	// ProcessUUIDLabel used to uniquely identify an Aperture process.
	ProcessUUIDLabel = "process_uuid"
	// DistCacheMemberIDLabel - label specifying unique identifier of the node in the olric cluster.
	DistCacheMemberIDLabel = "distcache_member_id"
	// DistCacheMemberNameLabel - label specifying name of the node in the olric cluster.
	DistCacheMemberNameLabel = "distcache_member_name"
	// PolicyNameLabel - label used in prometheus.
	PolicyNameLabel = "policy_name"
	// PolicyHashLabel - label used in prometheus.
	PolicyHashLabel = "policy_hash"
	// ComponentIDLabel - index of component in circuit label.
	ComponentIDLabel = "component_id"
	// DecisionTypeLabel - label for decision type rejected or accepted.
	DecisionTypeLabel = "decision_type"
	// WorkloadIndexLabel - label for choosing correct workload.
	WorkloadIndexLabel = "workload_index"
	// LimiterDroppedLabel - label to indicate that the particular limiter has dropped the request.
	LimiterDroppedLabel = "limiter_dropped"
	// SamplerDroppedLabel - label to indicate that the particular sampler has dropped the request.
	SamplerDroppedLabel = "sampler_dropped"
	// SignalNameLabel - label for saving circuit signal metrics.
	SignalNameLabel = "signal_name"
	// SubCircuitIDLabel - label for saving circuit id in signal metrics.
	SubCircuitIDLabel = "sub_circuit_id"
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
	// ServiceLookupsStatusLabel - status for ServiceLookupsMetricName.
	ServiceLookupsStatusLabel = "status"
	// ServiceLookupsStatusOK - service lookup status OK.
	ServiceLookupsStatusOK = FlowStatusOK
	// ServiceLookupsStatusError - service lookup status Error.
	ServiceLookupsStatusError = FlowStatusError
	// FlowStatusLabel - flow status.
	FlowStatusLabel = "flow_status"
	// FlowStatusOK - flow status OK.
	FlowStatusOK = "OK"
	// FlowStatusError - flow status Error.
	FlowStatusError = "Error"
	// FlowControlCheckDecisionTypeLabel - label for decision type dropped or accepted.
	FlowControlCheckDecisionTypeLabel = "decision_type"
	// FlowControlCheckErrorReasonLabel - label for error reason on FCS Check request.
	FlowControlCheckErrorReasonLabel = "error_reason"
	// FlowControlCheckRejectReasonLabel - label for reject reason on FCS Check request.
	FlowControlCheckRejectReasonLabel = "reject_reason"
	// AgentGroupLabel - label for agent group.
	AgentGroupLabel = "agent_group"

	// DEFAULTS.

	// DefaultWorkloadIndex - when workload is not specified this value is used.
	DefaultWorkloadIndex = "default"
	// DefaultAgentGroup - default agent group.
	DefaultAgentGroup = "default"

	// K8S METRICS.

	// K8sPodCount - number of pods in the cluster.
	K8sPodCount = "k8s_pod_count"
	// K8sNamespaceName - namespace of a resource.
	K8sNamespaceName = "k8s_namespace_name"
	// K8sNodeName - name of a node.
	K8sNodeName = "k8s_node_name"
	// K8sReplicasetName - name of a replicaset.
	K8sReplicasetName = "k8s_replicaset_name"
	// K8sDaemonsetName - name of a daemonset.
	K8sDaemonsetName = "k8s_daemonset_name"
	// K8sStatefulsetName - name of a statefulset.
	K8sStatefulsetName = "k8s_statefulset_name"
	// K8sDeploymentName - name of a deployment.
	K8sDeploymentName = "k8s_deployment_name"
)
