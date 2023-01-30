package controller

import (
	"crypto/tls"

	promapi "github.com/prometheus/client_golang/api"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
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

// OTELFxIn consumes parameters via Fx.
type OTELFxIn struct {
	fx.In
	Unmarshaller    config.Unmarshaller
	Listener        *listener.Listener
	PromClient      promapi.Client
	TLSConfig       *tls.Config
	ServerTLSConfig tlsconfig.ServerTLSConfig
}

func provideController(in OTELFxIn) (*otelconfig.OTELConfig, error) {
	var controllerCfg ControllerOTELConfig
	if err := in.Unmarshaller.UnmarshalKey("otel", &controllerCfg); err != nil {
		return nil, err
	}

	otelCfg := otelconfig.NewOTELConfig()
	otelCfg.SetDebugPort(&controllerCfg.CommonOTELConfig)
	otelCfg.AddDebugExtensions(&controllerCfg.CommonOTELConfig)

	addMetricsPipeline(otelCfg, &controllerCfg, in.TLSConfig, in.Listener, in.PromClient)
	otelCfg.AddExporter(otelconsts.ExporterLogging, nil)
	otelconfig.AddAlertsPipeline(otelCfg, controllerCfg.CommonOTELConfig)
	return otelCfg, nil
}

// addMetricsPipeline adds metrics to pipeline for controller OTEL collector.
func addMetricsPipeline(
	config *otelconfig.OTELConfig,
	controllerConfig *ControllerOTELConfig,
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
	config *otelconfig.OTELConfig,
	controllerConfig *ControllerOTELConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
) {
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-controller-self", tlsConfig, lis),
		otelconfig.BuildOTELScrapeConfig("aperture-controller-otel", controllerConfig.CommonOTELConfig),
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
