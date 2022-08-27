package otel

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/otel"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/utils"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/heartbeats"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
)

const (
	processorBatchMetricsSlow = "batch/metrics-slow"
	processorRollup           = "rollup"
	processorAttributes       = "attributes/insert"

	exporterFluxninja = "otlp/fluxninja"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				Provide,
				fx.ParamTags(config.NameTag("base"), config.NameTag("heartbeats-grpc-client-config"), config.NameTag("heartbeats-http-client-config")),
				fx.ResultTags(config.GroupTag("plugin")),
			),
		),
		fx.Invoke(Invoke),
	)
}

func Provide(baseConfig *otelcollector.OTELConfig,
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
	unmarshaller config.Unmarshaller,
) (*otelcollector.OTELConfig, error) {
	var pluginConfig pluginconfig.FluxNinjaPluginConfig
	if err := unmarshaller.UnmarshalKey(pluginconfig.PluginConfigKey, &pluginConfig); err != nil {
		return nil, err
	}
	config := otelcollector.NewOTELConfig()
	config.AddDebugExtensions()
	addFluxninjaExporter(config, &pluginConfig, grpcClientConfig, httpClientConfig)

	return config, nil
}

type ConstructorIn struct {
	fx.In
	Lifecycle     fx.Lifecycle
	BaseConfig    *otelcollector.OTELConfig   `name:"base"`
	PluginConfigs []*otelcollector.OTELConfig `group:"plugin"`
	Heartbeats    *heartbeats.Heartbeats
}

func Invoke(in ConstructorIn) {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			controllerID := in.Heartbeats.GetControllerInfo().Id

			for _, config := range in.PluginConfigs {
				addAttributesProcessor(config, controllerID)

				if logsPipeline, exists := in.BaseConfig.Service.Pipeline("logs"); exists {
					addFNToPipeline("logs", config, logsPipeline)
				}
				if tracesPipeline, exists := in.BaseConfig.Service.Pipeline("traces"); exists {
					addFNToPipeline("traces", config, tracesPipeline)
				}
				if _, exists := in.BaseConfig.Service.Pipeline("metrics/fast"); exists {
					addMetricsSlowPipeline(config)
				}
				if _, exists := in.BaseConfig.Service.Pipeline("metrics/controller-fast"); exists {
					addMetricsControllerSlowPipeline(config)
				}
			}
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})
}

func addAttributesProcessor(config *otelcollector.OTELConfig, controllerID string) {
	config.AddProcessor(processorAttributes, map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    "controller_id",
				"action": "insert",
				"value":  controllerID,
			},
		},
	})
}

func addFNToPipeline(
	name string,
	config *otelcollector.OTELConfig,
	pipeline otelcollector.Pipeline,
) {
	pipeline.Processors = append(pipeline.Processors, processorAttributes)
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
			processorAttributes,
		},
		Exporters: []string{exporterFluxninja},
	})
}

func addMetricsControllerSlowPipeline(config *otelcollector.OTELConfig) {
	config.AddBatchProcessor(processorBatchMetricsSlow, 10*time.Second, 10000)
	config.Service.AddPipeline("metrics/controller-slow", otelcollector.Pipeline{
		Receivers: []string{otel.ReceiverPrometheus},
		Processors: []string{
			processorBatchMetricsSlow,
			processorAttributes,
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
