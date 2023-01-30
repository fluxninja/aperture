package controller

import (
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

// swagger:operation POST /otel controller-configuration OTEL
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ControllerOTELConfig"

// ControllerOTELConfig is the configuration for Agent's OTEL collector.
// swagger:model
// +kubebuilder:object:generate=true
type ControllerOTELConfig struct {
	otelconfig.CommonOTELConfig `json:",inline"`
}

func provideController(cfg *otelconfig.OTELParams) *otelconfig.OTELConfig {
	addMetricsPipeline(cfg)
	cfg.Config.AddExporter(otelconsts.ExporterLogging, nil)
	otelconfig.AddAlertsPipeline(cfg)
	return cfg.Config
}

// addMetricsPipeline adds metrics to pipeline for controller OTEL collector.
func addMetricsPipeline(cfg *otelconfig.OTELParams) {
	config := cfg.Config
	addPrometheusReceiver(config, cfg)
	otelconfig.AddPrometheusRemoteWriteExporter(config, cfg.PromClient)
	config.Service.AddPipeline("metrics/controller-fast", otelconfig.Pipeline{
		Receivers:  []string{otelconsts.ReceiverPrometheus},
		Processors: []string{},
		Exporters:  []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

func addPrometheusReceiver(config *otelconfig.OTELConfig, cfg *otelconfig.OTELParams) {
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-controller-self", cfg),
		otelconfig.BuildOTELScrapeConfig("aperture-controller-otel", cfg),
	}
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
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
