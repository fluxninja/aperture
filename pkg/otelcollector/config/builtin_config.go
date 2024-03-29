package otelconfig

import (
	"crypto/tls"
	"fmt"

	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/net/listener"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	promapi "github.com/prometheus/client_golang/api"
)

// AddAlertsPipeline adds reusable alerts pipeline.
func AddAlertsPipeline(config *Config, cfg CommonOTelConfig, extraProcessors ...string) {
	config.AddReceiver(otelconsts.ReceiverAlerts, map[string]any{})
	config.AddProcessor(otelconsts.ProcessorAlertsNamespace, map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    otelconsts.AlertNamespaceLabel,
				"action": "insert",
				"value":  info.Hostname,
			},
		},
	})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchAlerts,
		cfg.BatchAlerts.Timeout.AsDuration(),
		cfg.BatchAlerts.SendBatchSize,
		cfg.BatchAlerts.SendBatchMaxSize,
	)
	config.AddExporter(otelconsts.ExporterAlerts, nil)

	processors := []string{
		otelconsts.ProcessorBatchAlerts,
		otelconsts.ProcessorAlertsNamespace,
	}
	processors = append(processors, extraProcessors...)
	config.Service.AddPipeline("logs/alerts", Pipeline{
		Receivers:  []string{otelconsts.ReceiverAlerts},
		Processors: processors,
		Exporters:  []string{otelconsts.ExporterLogging, otelconsts.ExporterAlerts},
	})
}

// AddPrometheusRemoteWriteExporter adds Prometheus remote write exporter which
// writes to controller Prometheus instance.
func AddPrometheusRemoteWriteExporter(config *Config, promClient promapi.Client) {
	endpoint := promClient.URL("api/v1/write", nil)
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTel. Need to use bare maps instead.
	config.AddExporter(otelconsts.ExporterPrometheusRemoteWrite, map[string]any{
		"endpoint": endpoint.String(),
		"resource_to_telemetry_conversion": map[string]any{
			"enabled": true,
		},
	})
}

// BuildApertureSelfScrapeConfig is a helper to create prometheus configuration
// which scrapes localhost.
func BuildApertureSelfScrapeConfig(
	name string,
	tlsConfig *tls.Config,
	lis *listener.Listener,
) map[string]any {
	scheme := "http"
	if tlsConfig != nil {
		scheme = "https"
	}
	return map[string]any{
		"job_name": name,
		"scheme":   scheme,
		"tls_config": map[string]any{
			"insecure_skip_verify": true,
		},
		"metrics_path": "/metrics",
		"static_configs": []map[string]any{
			{
				"targets": []string{lis.GetAddr()},
				"labels": map[string]any{
					metrics.InstanceLabel:    info.Hostname,
					metrics.ProcessUUIDLabel: info.UUID,
				},
			},
		},
	}
}

// BuildOTelScrapeConfig is a helper to create prometheus sonfiguration which
// scrapes OTel instance running on localhost.
func BuildOTelScrapeConfig(name string, cfg CommonOTelConfig) map[string]any {
	otelDebugTarget := fmt.Sprintf(":%d", cfg.Ports.DebugPort)
	return map[string]any{
		"job_name": name,
		"scheme":   "http",
		"tls_config": map[string]any{
			"insecure_skip_verify": true,
		},
		"metrics_path": "/metrics",
		"static_configs": []map[string]any{
			{
				"targets": []string{otelDebugTarget},
				"labels": map[string]any{
					metrics.InstanceLabel:    info.Hostname,
					metrics.ProcessUUIDLabel: info.UUID,
				},
			},
		},
	}
}

// AddHighCardinalityMetricsFilterProcessor adds filter processor which filters
// out high cardinality Aperture platform metrics.
func AddHighCardinalityMetricsFilterProcessor(config *Config) {
	config.AddProcessor(otelconsts.ProcessorFilterHighCardinalityMetrics, map[string]any{
		"metrics": map[string]interface{}{
			"exclude": map[string]interface{}{
				"match_type": "regexp",
				"regexp": map[string]interface{}{
					"cacheenabled":       true,
					"cachemaxnumentries": 1000,
				},
				"metric_names": []string{
					"grpc_server_handled.*",
					"grpc_server_handling.*",
					"grpc_server_msg.*",
					"grpc_server_started.*",
				},
			},
		},
	})
}
