package components

const (
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

	// ExporterLogging exports telemetry using Aperture logger.
	ExporterLogging = "aperturelogging"
	// ExporterPrometheusRemoteWrite exports metrics to local prometheus instance.
	ExporterPrometheusRemoteWrite = "prometheusremotewrite"
)
