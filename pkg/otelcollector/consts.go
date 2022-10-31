package otelcollector

const (
	/* Common labels available on all flow control integrations. */

	// ApertureSourceLabel is the label used to identify the flow control integration.
	ApertureSourceLabel = "aperture.source"
	// ApertureSourceSDK const for SDK source.
	ApertureSourceSDK = "sdk"
	// ApertureSourceEnvoy const for Envoy source.
	ApertureSourceEnvoy = "envoy"

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
	// WorkloadDurationLabel describes duration of the workload in milliseconds.
	WorkloadDurationLabel = "workload_duration_ms"
	// FlowDurationLabel describes duration of the flow in milliseconds.
	FlowDurationLabel = "flow_duration_ms"
	// ApertureProcessingDurationLabel describes Aperture's processing duration in milliseconds.
	ApertureProcessingDurationLabel = "aperture_processing_duration_ms"

	/* The following are derived labels that are applied based on contents of check response. */

	// ApertureDecisionTypeLabel describes the decision type taken by policy.
	ApertureDecisionTypeLabel = "aperture.decision_type"
	// ApertureErrorLabel describes the error reason of the decision taken by policy.
	ApertureErrorLabel = "aperture.error"
	// ApertureRejectReasonLabel describes the reject reason of the decision taken by policy.
	ApertureRejectReasonLabel = "aperture.reject_reason"
	// ApertureRateLimitersLabel describes rate limiters matched to the traffic.
	ApertureRateLimitersLabel = "aperture.rate_limiters"
	// ApertureDroppingRateLimitersLabel describes rate limiters dropping the traffic.
	ApertureDroppingRateLimitersLabel = "aperture.dropping_rate_limiters"
	// ApertureConcurrencyLimitersLabel describes rate limiters matched to the traffic.
	ApertureConcurrencyLimitersLabel = "aperture.concurrency_limiters"
	// ApertureDroppingConcurrencyLimitersLabel describes rate limiters dropping the traffic.
	ApertureDroppingConcurrencyLimitersLabel = "aperture.dropping_concurrency_limiters"
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
	// ApertureResponseStatusLabel label to denote OK or Error across all protocols.
	ApertureResponseStatusLabel = "aperture.response_status"
	// ApertureResponseStatusOK OK response across all protocols.
	ApertureResponseStatusOK = ApertureFeatureStatusOK
	// ApertureResponseStatusError Error response across all protocols.
	ApertureResponseStatusError = ApertureFeatureStatusError

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
	// EnvoyResponseDurationLabel from envoy access logs.
	EnvoyResponseDurationLabel = "RESPONSE_DURATION"
	// EnvoyBytesReceivedLabel from envoy access logs.
	EnvoyBytesReceivedLabel = "BYTES_RECEIVED"
	// EnvoyBytesSentLabel from envoy access logs.
	EnvoyBytesSentLabel = "BYTES_SENT"

	// EnvoyMissingAttributeValue is a special attribute value, which can
	// happen when (eg. Envoy's) logger tries to send attribute value, but it
	// is not available. Eg. In case authz couldn't reach agent, so we know
	// nothing about flowcontrol policies.  Note that this is different case
	// from "just empty", eg. "", "[]" or "{}".
	EnvoyMissingAttributeValue = "-"

	/* SDK specific labels. */

	// ApertureFeatureStatusLabel describes the status of the feature.
	ApertureFeatureStatusLabel = "aperture.feature.status"
	// ApertureFeatureStatusOK const for OK status.
	ApertureFeatureStatusOK = "OK"
	// ApertureFeatureStatusError const for error status.
	ApertureFeatureStatusError = "Error"
	// ApertureFlowStartTimestampLabel is the start timestamp of the flow.
	ApertureFlowStartTimestampLabel = "aperture.flow_start_timestamp"
	// ApertureFlowEndTimestampLabel is the end timestamp of the flow.
	ApertureFlowEndTimestampLabel = "aperture.flow_end_timestamp"
	// ApertureWorkloadStartTimestampLabel is the start timestamp of the workload.
	ApertureWorkloadStartTimestampLabel = "aperture.workload_start_timestamp"

	/* Specific to infra metrics pipeline. */

	// EntityNameLabel describes entity name e.g. pod name.
	EntityNameLabel = "entity_name"

	/* Aperture specific enrichment labels. */

	// AgentGroupLabel describes agent group to which metrics refer.
	AgentGroupLabel = "agent_group"

	/* Specific to Agent and Controller OTEL collector factories. */

	// ReceiverOTLP collects logs from libraries and SDKs.
	ReceiverOTLP = "otlp"
	// ReceiverPrometheus collects metrics from environment and services.
	ReceiverPrometheus = "prometheus"

	// ProcessorEnrichment enriches metrics with discovery data.
	ProcessorEnrichment = "enrichment"
	// ProcessorMetrics generates metrics based on logs and exposes them
	// on application prometheus metrics endpoint.
	ProcessorMetrics = "metrics"
	// ProcessorBatchPrerollup batches incoming data before rolling up. This is
	// required, as rollup processor can only roll up data inside a single batch.
	ProcessorBatchPrerollup = "batch/prerollup"
	// ProcessorBatchPostrollup batches data after rolling up, as roll up process
	// shrinks number of data points significantly.
	ProcessorBatchPostrollup = "batch/postrollup"
	// ProcessorRollup rolls up data to decrease cardinality.
	ProcessorRollup = "rollup"
	// ProcessorAgentGroup adds `agent_group` attribute.
	ProcessorAgentGroup = "attributes/agent_group"
	// ProcessorTracesToLogs converts received tracess to logs and passes them to configured
	// log exporter.
	ProcessorTracesToLogs = "tracestologs"

	// ExporterLogging exports telemetry using Aperture logger.
	ExporterLogging = "logging"
	// ExporterPrometheusRemoteWrite exports metrics to local prometheus instance.
	ExporterPrometheusRemoteWrite = "prometheusremotewrite"
	// ExporterOTLPLoopback exports OTLP data to local OTLP receiver. To be used only
	// with ProcessorSpanToLog.
	ExporterOTLPLoopback = "otlp/loopback"
)
