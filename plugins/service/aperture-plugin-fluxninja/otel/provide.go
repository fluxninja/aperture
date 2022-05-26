package otel

import (
	_ "embed"
	"fmt"
	"time"

	grpcclient "aperture.tech/aperture/pkg/net/grpc"
	httpclient "aperture.tech/aperture/pkg/net/http"
	"aperture.tech/aperture/pkg/net/tlsconfig"
	"go.uber.org/fx"

	"aperture.tech/aperture/pkg/config"
	"aperture.tech/aperture/pkg/otel"
	"aperture.tech/aperture/pkg/otelcollector"
	"aperture.tech/aperture/pkg/utils"
	"aperture.tech/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
)

const (
	processorBatchMetricsSlow = "batch/metrics-slow"
	processorRollup           = "rollup"

	exporterFluxninja = "otlp/fluxninja"
)

func ProvideAnnotatedPluginConfig() fx.Option {
	return fx.Option(
		fx.Provide(
			fx.Annotate(
				ProvidePluginOTELConfig,
				fx.ParamTags(config.NameTag("base"), config.NameTag("heartbeats-grpc-client-config"), config.NameTag("heartbeats-http-client-config")),
				fx.ResultTags(config.GroupTag("plugin")),
			),
		),
	)
}

func ProvidePluginOTELConfig(baseConfig *otelcollector.OTELConfig,
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
	unmarshaller config.Unmarshaller,
) (*otelcollector.OTELConfig, error) {
	var pluginConfig pluginconfig.FluxNinjaPluginConfig
	if err := unmarshaller.UnmarshalKey(pluginconfig.PluginConfigKey, &pluginConfig); err != nil {
		return nil, err
	}
	config := otelcollector.NewOTELConfig()
	addFluxninjaExporter(config, &pluginConfig, grpcClientConfig, httpClientConfig)
	if logsPipeline, exists := baseConfig.Service.Pipeline("logs"); exists {
		addPipelineWithFNExporter("logs", config, &pluginConfig, logsPipeline)
	}
	if tracesPipeline, exists := baseConfig.Service.Pipeline("traces"); exists {
		addPipelineWithFNExporter("traces", config, &pluginConfig, tracesPipeline)
	}
	if _, exists := baseConfig.Service.Pipeline("metrics/fast"); exists {
		addMetricsSlowPipeline(config)
	}
	if metricsPipeline, exists := baseConfig.Service.Pipeline("metrics/controller"); exists {
		addPipelineWithFNExporter("metrics/fluxninja", config, &pluginConfig, metricsPipeline)
	}
	return config, nil
}

func addPipelineWithFNExporter(
	name string,
	config *otelcollector.OTELConfig,
	pluginConfig *pluginconfig.FluxNinjaPluginConfig,
	pipeline otelcollector.Pipeline,
) {
	pipeline.Exporters = append(pipeline.Exporters, exporterFluxninja)
	config.Service.AddPipeline(name, pipeline)
}

func addMetricsSlowPipeline(config *otelcollector.OTELConfig) {
	config.AddBatchProcessor(processorBatchMetricsSlow, 10*time.Second, 10000)
	config.Service.AddPipeline("metrics/slow", otelcollector.Pipeline{
		Receivers: []string{otel.ReceiverPrometheus},
		Processors: []string{
			otel.ProcessorEnrichment,
			processorBatchMetricsSlow,
		},
		Exporters: []string{exporterFluxninja},
	})
}

func addFluxninjaExporter(config *otelcollector.OTELConfig,
	pluginConfig *pluginconfig.FluxNinjaPluginConfig,
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
) {
	cfg := map[string]interface{}{
		"endpoint": pluginConfig.FluxNinjaEndpoint,
		"headers": map[string]interface{}{
			"authorization": fmt.Sprintf("Bearer %s", pluginConfig.APIKey),
		},
	}

	var clientTLSConfig tlsconfig.ClientTLSConfig
	tlsCfg := make(map[string]interface{})

	if utils.IsHTTPUrl(pluginConfig.FluxNinjaEndpoint) {
		clientTLSConfig = httpClientConfig.ClientTLSConfig
	} else {
		clientTLSConfig = grpcClientConfig.ClientTLSConfig
		tlsCfg["insecure"] = grpcClientConfig.Insecure
	}

	tlsCfg["insecure_skip_verify"] = clientTLSConfig.InsecureSkipVerify
	tlsCfg["cert_file"] = clientTLSConfig.CertFile
	tlsCfg["key_file"] = clientTLSConfig.KeyFile
	tlsCfg["ca_file"] = clientTLSConfig.CAFile

	cfg["tls"] = tlsCfg

	config.AddExporter(exporterFluxninja, cfg)
}
