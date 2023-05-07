package agent

import (
	"crypto/tls"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"k8s.io/client-go/rest"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	agentconfig "github.com/fluxninja/aperture/cmd/aperture-agent/config"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	otelcustom "github.com/fluxninja/aperture/pkg/otelcollector/custom"
)

func provideAgent(
	unmarshaller config.Unmarshaller,
	lis *listener.Listener,
	promClient promapi.Client,
	tlsConfig *tls.Config,
) (*otelconfig.OTelConfig, error) {
	var agentCfg agentconfig.AgentOTelConfig
	if err := unmarshaller.UnmarshalKey("otel", &agentCfg); err != nil {
		return nil, err
	}

	otelCfg := otelconfig.NewOTelConfig()
	otelCfg.SetDebugPort(&agentCfg.CommonOTelConfig)
	otelCfg.AddDebugExtensions(&agentCfg.CommonOTelConfig)

	addLogsPipeline(otelCfg, &agentCfg)
	addTracesPipeline(otelCfg, lis)
	addMetricsPipeline(otelCfg, &agentCfg, tlsConfig, lis, promClient)

	customConfig := map[string]*policylangv1.InfraMeter{}
	if !agentCfg.DisableKubeletScraper {
		if _, ok := customConfig[otelconsts.ReceiverKubeletStats]; !ok {
			customConfig[otelconsts.ReceiverKubeletStats] = otelcustom.InfraMeterForKubeletStats()
		}
	}

	if err := otelcustom.AddCustomMetricsPipelines(otelCfg, customConfig); err != nil {
		return nil, err
	}
	otelconfig.AddAlertsPipeline(otelCfg, agentCfg.CommonOTelConfig, otelconsts.ProcessorAgentResourceLabels)
	return otelCfg, nil
}

func addLogsPipeline(
	config *otelconfig.OTelConfig,
	userConfig *agentconfig.AgentOTelConfig,
) {
	// Common dependencies for pipelines
	config.AddReceiver(otelconsts.ReceiverOTLP, otlpreceiver.Config{})
	// Note: Passing map[string]interface{}{} instead of real config, so that
	// processors' configs' default work.
	config.AddProcessor(otelconsts.ProcessorMetrics, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPrerollup,
		userConfig.BatchPrerollup.Timeout.AsDuration(),
		userConfig.BatchPrerollup.SendBatchSize,
		userConfig.BatchPrerollup.SendBatchMaxSize,
	)
	config.AddProcessor(otelconsts.ProcessorRollup, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPostrollup,
		userConfig.BatchPostrollup.Timeout.AsDuration(),
		userConfig.BatchPostrollup.SendBatchSize,
		userConfig.BatchPostrollup.SendBatchMaxSize,
	)
	config.AddExporter(otelconsts.ExporterLogging, nil)

	processors := []string{
		otelconsts.ProcessorMetrics,
		otelconsts.ProcessorBatchPrerollup,
		otelconsts.ProcessorRollup,
		otelconsts.ProcessorBatchPostrollup,
		otelconsts.ProcessorAgentGroup,
	}

	config.Service.AddPipeline("logs", otelconfig.Pipeline{
		Receivers:  []string{otelconsts.ReceiverOTLP, otelconsts.ConnectorAdapter},
		Processors: processors,
		Exporters:  []string{otelconsts.ExporterLogging},
	})
}

func addTracesPipeline(config *otelconfig.OTelConfig, _ *listener.Listener) {
	config.AddConnector(otelconsts.ConnectorAdapter, map[string]any{})
	config.Service.AddPipeline("traces", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverOTLP},
		Exporters: []string{otelconsts.ConnectorAdapter},
	})
}

func addMetricsPipeline(
	config *otelconfig.OTelConfig,
	agentConfig *agentconfig.AgentOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
	promClient promapi.Client,
) {
	addPrometheusReceiver(config, agentConfig, tlsConfig, lis)
	otelconfig.AddPrometheusRemoteWriteExporter(config, promClient)
	config.Service.AddPipeline("metrics/fast", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverPrometheus},
		Processors: []string{
			otelconsts.ProcessorAgentGroup,
		},
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

func addPrometheusReceiver(
	config *otelconfig.OTelConfig,
	agentConfig *agentconfig.AgentOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
) {
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-self", tlsConfig, lis),
		otelconfig.BuildOTelScrapeConfig("aperture-otel", agentConfig.CommonOTelConfig),
	}

	if !agentConfig.DisableKubernetesScraper {
		_, err := rest.InClusterConfig()
		if err == rest.ErrNotInCluster {
			log.Debug().Msg("K8s environment not detected. Skipping K8s scrape configurations.")
		} else if err != nil {
			log.Warn().Err(err).Msg("Error when discovering k8s environment")
		} else {
			log.Debug().Msg("K8s environment detected. Adding K8s scrape configurations.")
			scrapeConfigs = append(scrapeConfigs, buildKubernetesPodsScrapeConfig())
		}
	}

	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTel. Need to use bare maps instead.
	config.AddReceiver(otelconsts.ReceiverPrometheus, map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval":     "1s",
				"scrape_timeout":      "1s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}

func buildKubernetesPodsScrapeConfig() map[string]any {
	return map[string]any{
		"job_name":     "kubernetes-pods",
		"scheme":       "http",
		"metrics_path": "/metrics",
		"kubernetes_sd_configs": []map[string]any{
			{"role": "pod"},
		},
		"relabel_configs": []map[string]any{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_pod_node_name"},
				"action":        "keep",
				"regex":         info.Hostname,
			},
			// Scrape only pods which have github.com/fluxninja/scrape=true annotation.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_aperture_tech_scrape"},
				"action":        "keep",
				"regex":         "true",
			},
			// Allow rewrite of scheme, path and port where prometheus metrics are served.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
				"action":        "replace",
				"regex":         "(https?)",
				"target_label":  "__scheme__",
			},
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
				"action":        "replace",
				"target_label":  "__metrics_path__",
				"regex":         "(.+)",
			},
			{
				"source_labels": []string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
				"action":        "replace",
				"regex":         `([^:]+)(?::\d+)?;(\d+)`,
				"replacement":   "$$1:$$2",
				"target_label":  "__address__",
			},
		},
		"metric_relabel_configs": []map[string]any{
			// For now, dropping everything. In future, we'll want to filter in some
			// metrics based on policies. See #4632.
			{
				"source_labels": []string{},
				"action":        "drop",
			},
		},
	}
}
