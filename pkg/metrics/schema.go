package metrics

import "time"

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
	// InfiniteSignalReadingsTotalMetricName - counts infinite signal readings.
	InfiniteSignalReadingsTotalMetricName = "infinite_signal_readings_total"
	// InvalidFluxMeterTotalMetricName - counts invalid flux meters.
	InvalidFluxMeterTotalMetricName = "invalid_flux_meter_total"
	// RateLimiterCounterTotalMetricName - name of the counter describing times rate limiter was triggered.
	RateLimiterCounterTotalMetricName = "rate_limiter_counter_total"
	// ClassifierCounterTotalMetricName - name of the counter describing times classifier was triggered.
	ClassifierCounterTotalMetricName = "classifier_counter_total"
	// SamplerCounterTotalMetricName - name of the counter describing times sampler was triggered.
	SamplerCounterTotalMetricName = "sampler_counter_total"
	// ConcurrencyLimiterCounterTotalMetricName - name of the counter describing times concurrency limiter was triggered.
	ConcurrencyLimiterCounterTotalMetricName = "concurrency_limiter_counter_total"

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
	// DistCachePartitionsCountMetricsName - metric for the current number of non-empty partitions owned by given node.
	DistCachePartitionsCountMetricsName = "distcache_partitions_count"
	// DistCacheBackupPartitionsCountMetricsName - metric for the current number of non-empty backups owned by given node.
	DistCacheBackupPartitionsCountMetricsName = "distcache_backup_partitions_count"
	// DistCacheFragmentMigrationEventsTotalMetricsName - metric for cumulative number of fragment migration (outgoing) events.
	DistCacheFragmentMigrationEventsTotalMetricsName = "distcache_fragment_migration_events_total"
	// DistCacheFragmentReceivedEventsTotalMetricsName - metric for cumulative number of fragment received (incoming) events.
	DistCacheFragmentReceivedEventsTotalMetricsName = "distcache_fragment_received_events_total"
	// DistCachePartitionsLength - metric for current length of partitions on given node.
	DistCachePartitionsLength = "distcache_partitions_length"
	// DistCacheBackupPartitionsLength - metric for current length of backup partitions on given node.
	DistCacheBackupPartitionsLength = "distcache_backup_partitions_length"

	// Persistent cache metrics.

	// PersistentCacheGetMissesMetricName - metric for cumulative number of persistent cache misses.
	PersistentCacheGetMissesMetricName = "persistent_cache_get_misses_total"
	// PersistentCacheGetHitsMetricName - metric for cumulative number of persistent cache hits.
	PersistentCacheGetHitsMetricName = "persistent_cache_get_hits_total"
	// PersistentCacheTypeLabel - denotes if hit or miss sources from in_memory of object_storage layer.
	PersistentCacheTypeLabel = "cache_type"
	// PersistentCacheTypeInMemory - in_memory layer of persistent cache.
	PersistentCacheTypeInMemory = "in_memory"
	// PersistentCacheTypeObjectStorage - object_storage layer of persistent cache.
	PersistentCacheTypeObjectStorage = "object_storage"
	// PersistentCacheOperationDurationMetricName - metric for duration of persistent cache operations.
	PersistentCacheOperationDurationMetricName = "persistent_cache_operation_duration_ms"
	// PersistentCacheOperationLabel - label to define type of operation performed on persistent cache.
	PersistentCacheOperationLabel = "operation"
	// PersistentCacheOperationGet - get operation.
	PersistentCacheOperationGet = "get"
	// PersistentCacheOperationPut - put operation.
	PersistentCacheOperationPut = "put"
	// PersistentCacheOperationDelete - delete operation.
	PersistentCacheOperationDelete = "delete"

	// Cache metrics.

	// CacheLookupHitsTotalMetricName - metric for cumulative number of cache lookup hits.
	CacheLookupHitsTotalMetricName = "cache_lookup_hits_total"
	// CacheLookupMissesTotalMetricName - metric for cumulative number of cache lookup misses.
	CacheLookupMissesTotalMetricName = "cache_lookup_misses_total"
	// CacheOperationResultsTotalMetricName - metric for cumulative number of cache operation results.
	CacheOperationResultsTotalMetricName = "cache_operation_results_total"

	// Workload metrics.

	// WorkloadLatencyMetricName - metric used for grouping latencies per workload.
	WorkloadLatencyMetricName = "workload_latency_ms"
	// WorkloadLatencySumMetricName - metric from workload histogram.
	WorkloadLatencySumMetricName = "workload_latency_ms_sum"
	// WorkloadLatencyCountMetricName - metric from workload histogram.
	WorkloadLatencyCountMetricName = "workload_latency_ms_count"
	// WorkloadCounterMetricName - metric used for counting workload requests.
	WorkloadCounterMetricName = "workload_requests_total"
	// RequestInQueueDurationMetricName - metric used for grouping durations for requests in queue of Scheduler.
	RequestInQueueDurationMetricName = "request_in_queue_duration_ms"
	// WorkloadPreemptedTokensMetricName - metric used for counting tokens preempted per request measured end-to-end in the scheduler across all workloads.
	WorkloadPreemptedTokensMetricName = "workload_preempted_tokens"
	// WorkloadDelayedTokensMetricName - metric used for counting tokens delayed per request measured end-to-end in the scheduler across all workloads.
	WorkloadDelayedTokensMetricName = "workload_delayed_tokens"
	// WorkloadOnTimeMetricName - metric used for counting requests that are on time, neither preempted nor delayed measured end-to-end in the scheduler across all workloads.
	WorkloadOnTimeMetricName = "workload_on_time_total"
	// FairnessPreemptedTokensMetricName - metric used for counting tokens preempted per request measured at fairness queues within the same workload.
	FairnessPreemptedTokensMetricName = "fairness_preempted_tokens"
	// FairnessDelayedTokensMetricName - metric used for counting tokens delayed per request measured at fairness queues within the same workload.
	FairnessDelayedTokensMetricName = "fairness_delayed_tokens"
	// FairnessOnTimeMetricName - metric used for counting requests that are on time, neither preempted nor delayed measured at fairness queues within the same workload.
	FairnessOnTimeMetricName = "fairness_on_time_total"

	// IncomingTokensMetricName - total work measured in tokens of all incoming requests.
	IncomingTokensMetricName = "incoming_tokens_total"
	// AcceptedTokensMetricName - total work measured in tokens of all accepted requests.
	AcceptedTokensMetricName = "accepted_tokens_total"
	// RejectedTokensMetricName - total work measured in tokens of all rejected requests.
	RejectedTokensMetricName = "rejected_tokens_total"

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
	// FlowControlEndsMetricName - counter for flow ends.
	FlowControlEndsMetricName = "flowcontrol_ends_total"

	// Check flow metrics.

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
	// PolicyNameLabel - policy name label.
	PolicyNameLabel = "policy_name"
	// PolicyHashLabel - policy hash label.
	PolicyHashLabel = "policy_hash"
	// InfraMeterNameLabel - Label for name of infra meter in resources.
	InfraMeterNameLabel = "infra_meter_name"
	// ComponentIDLabel - index of component in circuit label.
	ComponentIDLabel = "component_id"
	// DecisionTypeLabel - label for decision type rejected or accepted.
	DecisionTypeLabel = "decision_type"
	// DecisionTypeRejected - decision type rejected.
	DecisionTypeRejected = "DECISION_TYPE_REJECTED"
	// DecisionTypeAccepted - decision type accepted.
	DecisionTypeAccepted = "DECISION_TYPE_ACCEPTED"
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
	// CacheTypeLabel - label for cache type.
	CacheTypeLabel = "type"
	// CacheOperationTypeLabel - label for type of cache operation.
	CacheOperationTypeLabel = "operation_type"
	// CacheOperationTypeLookup - cache lookup operation.
	CacheOperationTypeLookup = "lookup"
	// CacheOperationTypeDelete - cache delete operation.
	CacheOperationTypeDelete = "delete"
	// CacheOperationTypeUpsert - cache upsert operation.
	CacheOperationTypeUpsert = "upsert"
	// CacheOperationStatusLabel - label for cache operation status.
	CacheOperationStatusLabel = "status"
	// CacheOperationStatusSuccess - cache operation success.
	CacheOperationStatusSuccess = "success"
	// CacheOperationStatusError - cache operation error.
	CacheOperationStatusError = "error"
	// ControlPointLabel - label for control point.
	ControlPointLabel = "control_point"

	// DEFAULTS.

	// DefaultWorkloadIndex - when workload is not specified this value is used.
	DefaultWorkloadIndex = "default"
	// DefaultAgentGroup - default agent group.
	DefaultAgentGroup = "default"

	// ScrapeInterval - Scrape interval for metrics.
	ScrapeInterval = time.Second * 10
)
