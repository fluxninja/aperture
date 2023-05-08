package controller

import (
	"crypto/tls"

	promapi "github.com/prometheus/client_golang/api"

	controllerconfig "github.com/fluxninja/aperture/cmd/aperture-controller/config"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

func provideController(
	unmarshaller config.Unmarshaller,
	lis *listener.Listener,
	promClient promapi.Client,
	tlsConfig *tls.Config,
) (*otelconfig.OTelConfigProvider, error) {
	var controllerCfg controllerconfig.ControllerOTelConfig
	if err := unmarshaller.UnmarshalKey("otel", &controllerCfg); err != nil {
		return nil, err
	}

	otelCfg := otelconfig.NewOTelConfig()
	otelCfg.SetDebugPort(&controllerCfg.CommonOTelConfig)
	otelCfg.AddDebugExtensions(&controllerCfg.CommonOTelConfig)

	addMetricsPipeline(otelCfg, &controllerCfg, tlsConfig, lis, promClient)
	otelCfg.AddExporter(otelconsts.ExporterLogging, nil)
	otelconfig.AddAlertsPipeline(otelCfg, controllerCfg.CommonOTelConfig)
	baseConfigProvider := otelconfig.NewOTelConfigProvider("service", otelCfg)

	return baseConfigProvider, nil
}

// addMetricsPipeline adds metrics to pipeline for controller OTel collector.
func addMetricsPipeline(
	config *otelconfig.OTelConfig,
	controllerConfig *controllerconfig.ControllerOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
	promClient promapi.Client,
) {
	addPrometheusReceiver(config, controllerConfig, tlsConfig, lis)
	otelconfig.AddPrometheusRemoteWriteExporter(config, promClient)
	config.Service.AddPipeline("metrics/controller-fast", otelconfig.Pipeline{
		Receivers:  []string{otelconsts.ReceiverPrometheus},
		Processors: []string{},
		Exporters:  []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

func addPrometheusReceiver(
	config *otelconfig.OTelConfig,
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
				"scrape_interval":     "1s",
				"scrape_timeout":      "1s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}
