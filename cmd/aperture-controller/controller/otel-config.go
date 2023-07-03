package controller

import (
	"crypto/tls"

	promapi "github.com/prometheus/client_golang/api"

	controllerconfig "github.com/fluxninja/aperture/v2/cmd/aperture-controller/config"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
)

func provideController(
	unmarshaller config.Unmarshaller,
	lis *listener.Listener,
	promClient promapi.Client,
	tlsConfig *tls.Config,
) (*otelconfig.Provider, error) {
	var controllerCfg controllerconfig.ControllerOTelConfig
	if err := unmarshaller.UnmarshalKey("otel", &controllerCfg); err != nil {
		return nil, err
	}

	otelCfg := otelconfig.New()
	otelCfg.SetDebugPort(&controllerCfg.CommonOTelConfig)
	otelCfg.AddDebugExtensions(&controllerCfg.CommonOTelConfig)

	addMetricsPipeline(otelCfg, &controllerCfg, tlsConfig, lis, promClient)
	otelCfg.AddExporter(otelconsts.ExporterLogging, nil)
	otelconfig.AddAlertsPipeline(otelCfg, controllerCfg.CommonOTelConfig)

	return otelconfig.NewProvider("service", otelCfg), nil
}

// addMetricsPipeline adds metrics to pipeline for controller OTel collector.
func addMetricsPipeline(
	config *otelconfig.Config,
	controllerConfig *controllerconfig.ControllerOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
	promClient promapi.Client,
) {
	addPrometheusReceiver(config, controllerConfig, tlsConfig, lis)
	otelconfig.AddPrometheusRemoteWriteExporter(config, promClient)
	processors := []string{}
	if !controllerConfig.EnableHighCardinalityPlatformMetrics {
		otelconfig.AddHighCardinalityMetricsFilterProcessor(config)
		// Prepending processor so we drop metrics as soon as possible without any unnecessary operation on them.
		processors = append([]string{otelconsts.ProcessorFilterHighCardinalityMetrics}, processors...)
	}
	config.Service.AddPipeline("metrics/controller-fast", otelconfig.Pipeline{
		Receivers:  []string{otelconsts.ReceiverPrometheus},
		Processors: processors,
		Exporters:  []string{otelconsts.ExporterPrometheusRemoteWrite, otelconsts.ExporterLogging},
	})
}

func addPrometheusReceiver(
	config *otelconfig.Config,
	controllerConfig *controllerconfig.ControllerOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
) {
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-controller-self", tlsConfig, lis),
		otelconfig.BuildOTelScrapeConfig("aperture-controller-otel", controllerConfig.CommonOTelConfig),
	}
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTel. Need to use bare maps instead.
	config.AddReceiver(otelconsts.ReceiverPrometheus, map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval":     "10s",
				"scrape_timeout":      "1s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}
