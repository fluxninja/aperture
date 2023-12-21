package consts

import "github.com/fluxninja/aperture/v2/pkg/config"

const (
	/* Common labels available on all flow control integrations. */

	// ApertureSourceLabel is the label used to identify the flow control integration.
	ApertureSourceLabel = "aperture.source"
	// ApertureSourceSDK const for SDK source.
	ApertureSourceSDK = "sdk"
	// ApertureSourceEnvoy const for Envoy source.
	ApertureSourceEnvoy = "envoy"
	// ApertureSourceLua const for Lua source.
	ApertureSourceLua = "lua"
	// FeatureControlPoint const for feature control point.
	FeatureControlPoint = "feature"
	// HTTPControlPoint for envoy control point.
	HTTPControlPoint = "http"

	// ApertureCheckResponseLabel contains JSON encoded check response struct.
	ApertureCheckResponseLabel = "aperture.check_response"

	// ResponseReceivedLabel designates whether a response was received.
	ResponseReceivedLabel = "response_received"
	// ResponseReceivedTrue const for true value.
	ResponseReceivedTrue = "true"
	// ResponseReceivedFalse const for false value.
	ResponseReceivedFalse = "false"

	/* Derived label that is applied based on content of labels. */

	// ApertureServicesLabel describes services to which metrics refer.
	ApertureServicesLabel = "aperture.services"
	// ApertureControlPointLabel describes control point to which metrics refer.
	ApertureControlPointLabel = "aperture.control_point"
	// ApertureControlPointTypeLabel describes the type of control point.
	ApertureControlPointTypeLabel = "aperture.control_point_type"
	// WorkloadDurationLabel describes duration of the workload in milliseconds.
	WorkloadDurationLabel = "workload_duration_ms"
	// FlowDurationLabel describes duration of the flow in milliseconds.
	FlowDurationLabel = "flow_duration_ms"
	// ApertureProcessingDurationLabel describes Aperture's processing duration in milliseconds.
	ApertureProcessingDurationLabel = "aperture_processing_duration_ms"
	// ApertureSourceServiceLabel describes the source service FQDNs of the flow.
	ApertureSourceServiceLabel = "aperture.source_fqdns"
	// ApertureDestinationServiceLabel describes the destination service FQDNs of the flow.
	ApertureDestinationServiceLabel = "aperture.destination_fqdns"

	/* The following are derived labels that are applied based on contents of check response. */

	// ApertureDecisionTypeLabel describes the decision type taken by policy.
	ApertureDecisionTypeLabel = "aperture.decision_type"
	// ApertureRejectReasonLabel describes the reject reason of the decision taken by policy.
	ApertureRejectReasonLabel = "aperture.reject_reason"
	// ApertureRateLimitersLabel describes rate limiters matched to the traffic.
	ApertureRateLimitersLabel = "aperture.rate_limiters"
	// ApertureDroppingRateLimitersLabel describes rate limiters dropping the traffic.
	ApertureDroppingRateLimitersLabel = "aperture.dropping_rate_limiters"
	// ApertureLoadSchedulersLabel describes load schedulers matched to the traffic.
	ApertureLoadSchedulersLabel = "aperture.load_schedulers"
	// ApertureQuotaSchedulersLabel describes quota schedulers matched to the traffic.
	ApertureQuotaSchedulersLabel = "aperture.quota_schedulers"
	// ApertureDroppingLoadSchedulersLabel describes load schedulers dropping the traffic.
	ApertureDroppingLoadSchedulersLabel = "aperture.dropping_load_schedulers"
	// ApertureDroppingQuotaSchedulersLabel describes quota schedulers dropping the traffic.
	ApertureDroppingQuotaSchedulersLabel = "aperture.dropping_quota_schedulers"
	// ApertureSamplersLabel describes samplers matched to the traffic.
	ApertureSamplersLabel = "aperture.samplers"
	// ApertureDroppingSamplersLabel describes samplers dropping the traffic.
	ApertureDroppingSamplersLabel = "aperture.dropping_samplers"
	// ApertureWorkloadsLabel describes workloads matched to the traffic.
	ApertureWorkloadsLabel = "aperture.workloads"
	// ApertureDroppingWorkloadsLabel describes workloads dropping the traffic.
	ApertureDroppingWorkloadsLabel = "aperture.dropping_workloads"
	// ApertureFluxMetersLabel describes flux meters matched to the traffic.
	ApertureFluxMetersLabel = "aperture.flux_meters"
	// ApertureFlowLabelKeysLabel describes keys of flow labels matched to the traffic.
	ApertureFlowLabelKeysLabel = "aperture.flow_label_keys"
	// ApertureClassifiersLabel describes classifiers matched to the traffic.
	ApertureClassifiersLabel = "aperture.classifiers"
	// ApertureClassifierErrorsLabel describes encountered classifier errors for specified policy.
	ApertureClassifierErrorsLabel = "aperture.classifier_errors"
	// ApertureFlowStatusLabel label to denote OK or Error across all protocols.
	ApertureFlowStatusLabel = "aperture.flow.status"
	// ApertureFlowStatusOK const for OK status.
	ApertureFlowStatusOK = "OK"
	// ApertureFlowStatusError const for error status.
	ApertureFlowStatusError = "Error"
	// ApertureResultCacheLookupStatusLabel describes status of the cache lookup.
	ApertureResultCacheLookupStatusLabel = "aperture.cache_lookup_status"
	// ApertureResultCacheOperationStatusLabel describes status of the cache operation.
	ApertureResultCacheOperationStatusLabel = "aperture.cache_operation_status"
	// ApertureGlobalCacheLookupStatusesLabel describes statuses of the global cache lookups.
	ApertureGlobalCacheLookupStatusesLabel = "aperture.global_cache_lookup_statuses"
	// ApertureGlobalCacheOperationStatusesLabel describes statuses of the global cache operations.
	ApertureGlobalCacheOperationStatusesLabel = "aperture.global_cache_operation_statuses"

	/* HTTP Specific labels. */

	// HTTPStatusCodeLabel describes HTTP status code of the response.
	HTTPStatusCodeLabel = "http.status_code"
	// HTTPRequestContentLength describes length of the HTTP request content in bytes.
	HTTPRequestContentLength = "http.request_content_length"
	// HTTPResponseContentLength describes length of the HTTP response content in bytes.
	HTTPResponseContentLength = "http.response_content_length"

	/* Labels specific to Envoy. */

	// EnvoyAuthzDurationLabel describes duration of the Authz call in milliseconds.
	EnvoyAuthzDurationLabel = "authz_duration"

	// CheckHTTPDurationLabel describes duration of the CheckHTTP call in milliseconds.
	CheckHTTPDurationLabel = "checkhttp_duration"

	// ResponseDurationLabel from access logs.
	ResponseDurationLabel = "RESPONSE_DURATION"
	// BytesReceivedLabel from access logs.
	BytesReceivedLabel = "BYTES_RECEIVED"
	// BytesSentLabel from access logs.
	BytesSentLabel = "BYTES_SENT"

	// EnvoyMissingAttributeValue is a special attribute value, which can
	// happen when (eg. Envoy's) logger tries to send attribute value, but it
	// is not available. Eg. In case authz could not reach agent, so we know
	// nothing about flowcontrol policies.  Note that this is different case
	// from "just empty", eg. "", "[]" or "{}".
	EnvoyMissingAttributeValue = "-"

	// LuaMissingAttributeValue is a special attribute value, which can
	// happen when (eg. Lua's) logger tries to send attribute value, but it
	// is not available.
	LuaMissingAttributeValue = "-"

	/* SDK specific labels. */

	// ApertureFlowStartTimestampLabel is the start timestamp of the flow in nano seconds.
	// Deprecated: v3.0.0. Use `aperture.flow_start_timestamp_ms` instead.
	ApertureFlowStartTimestampLabel = "aperture.flow_start_timestamp"
	// ApertureFlowStartTimestampLabelMs is the start timestamp of the flow in milli seconds.
	ApertureFlowStartTimestampLabelMs = "aperture.flow_start_timestamp_ms"
	// ApertureFlowEndTimestampLabel is the end timestamp of the flow in nano seconds.
	// Deprecated: v3.0.0. Use `aperture.flow_end_timestamp_ms` instead.
	ApertureFlowEndTimestampLabel = "aperture.flow_end_timestamp"
	// ApertureFlowEndTimestampLabelMs is the end timestamp of the flow in milli seconds.
	ApertureFlowEndTimestampLabelMs = "aperture.flow_end_timestamp_ms"
	// ApertureWorkloadStartTimestampLabel is the start timestamp of the workload in nano seconds.
	// Deprecated: v3.0.0. Use `aperture.workload_start_timestamp_ms` instead.
	ApertureWorkloadStartTimestampLabel = "aperture.workload_start_timestamp"
	// ApertureWorkloadStartTimestampLabelMs is the start timestamp of the workload in milli seconds.
	ApertureWorkloadStartTimestampLabelMs = "aperture.workload_start_timestamp_ms"

	/* Aperture specific enrichment labels. */

	// AgentGroupLabel describes agent group to which metrics refer.
	AgentGroupLabel = "agent_group"
	// InstanceLabel describes agent group to which metrics refer.
	InstanceLabel = "instance"

	/* Specific to Agent and Controller OTel collector factories. */

	// ReceiverOTLP collects logs from libraries and SDKs.
	ReceiverOTLP = "otlp"
	// ReceiverPrometheus collects metrics from environment and services.
	ReceiverPrometheus = "prometheus"
	// ReceiverAlerts collects alerts from alerter.
	ReceiverAlerts = "alerts"
	// ReceiverKubeletStats collects metrics from kubelet.
	ReceiverKubeletStats = "kubeletstats"

	// ProcessorMetrics generates metrics based on logs and exposes them
	// on application prometheus metrics endpoint.
	ProcessorMetrics = "metrics"
	// ProcessorBatchPrerollup batches incoming data before rolling up. This is
	// required, as rollup processor can only roll up data inside a single batch.
	ProcessorBatchPrerollup = "batch/prerollup"
	// ProcessorBatchPostrollup batches data after rolling up, as roll up process
	// shrinks number of data points significantly.
	ProcessorBatchPostrollup = "batch/postrollup"
	// ProcessorBatchAlerts batches alerts before passing them to exporters.
	// This reduces number of calls to the Alertmanager.
	ProcessorBatchAlerts = "batch/alerts"
	// ProcessorRollup rolls up data to decrease cardinality.
	ProcessorRollup = "rollup"
	// ProcessorAgentGroup adds `agent_group` attribute.
	ProcessorAgentGroup = "attributes/agent_group"
	// ProcessorInfraMeter adds `service.name` resource attribute.
	ProcessorInfraMeter = "resource/infra_meter"
	// ProcessorAgentResourceLabels adds `instance` and `agent_group` resource attributes.
	ProcessorAgentResourceLabels = "transform/agent_resource_labels"
	// ProcessorAlertsNamespace adds host info as `namespace` attribute.
	ProcessorAlertsNamespace = "attributes/alerts"
	// ProcessorFilterKubeletStats filters in only metrics of interest.
	ProcessorFilterKubeletStats = "filter/kubeletstats"
	// ProcessorFilterHighCardinalityMetrics filters out high cardinality Aperture platform metrics.
	ProcessorFilterHighCardinalityMetrics = "filter/high_cardinality_metrics"
	// ProcessorK8sAttributes enriches metrics with k8s metadata.
	ProcessorK8sAttributes = "k8sattributes"
	// ProcessorK8sAttributesSelectors is the key name of selectors field in k8sattributes processor.
	ProcessorK8sAttributesSelectors = "selectors"

	// ExporterLogging exports telemetry using Aperture logger.
	ExporterLogging = "logging"
	// ExporterPrometheusRemoteWrite exports metrics to local prometheus instance.
	ExporterPrometheusRemoteWrite = "prometheusremotewrite"
	// ExporterAlerts exports alerts via alertmanager clients.
	ExporterAlerts = "alerts"

	// ConnectorAdapter adapts OTEL signals between pipelines.
	ConnectorAdapter = "adapter"

	/* Specific to alerts pipeline. */

	// AlertGeneratorURLLabel describes.
	AlertGeneratorURLLabel = "generator_url"
	// AlertNameLabel describes name of the alert.
	AlertNameLabel = "alertname"
	// AlertSeverityLabel also known as log level. Human readable string.
	AlertSeverityLabel = "severity"
	// AlertChannelsLabel is a comma-separated list of channels to which alert is assigned.
	AlertChannelsLabel = "alert_channels"
	// IsAlertLabel helps to differentiate normal logs from alert logs.
	IsAlertLabel = "is_alert"
	// AlertNamespaceLabel which required by Alertmanager. Set to host name.
	AlertNamespaceLabel = "namespace"
)

// FX tags used to pass OTel Collector factories.
var (
	BaseFxTag               = config.NameTag("base")
	ReceiverFactoriesFxTag  = config.GroupTag("otel-collector-receiver-factories")
	ProcessorFactoriesFxTag = config.GroupTag("otel-collector-processor-factories")
)
