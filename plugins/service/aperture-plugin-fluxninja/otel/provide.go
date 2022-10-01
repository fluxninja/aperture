package otel

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
	"github.com/imdario/mergo"
	"github.com/mitchellh/copystructure"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/utils"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/heartbeats"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
)

const (
	receiverPrometheus = "prometheus/fluxninja"

	processorBatchMetricsSlow = "batch/metrics-slow"
	processorRollup           = "rollup"
	processorAttributes       = "attributes/fluxninja"

	exporterFluxninja = "otlp/fluxninja"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideOtelConfig,
				fx.ParamTags(config.NameTag("base"),
					config.NameTag("heartbeats-grpc-client-config"),
					config.NameTag("heartbeats-http-client-config")),
				fx.ResultTags(config.GroupTag("plugin-config")),
			),
		),
	)
}

func provideOtelConfig(baseConfig *otelcollector.OTELConfig,
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
	lifecycle fx.Lifecycle,
	heartbeats *heartbeats.Heartbeats,
	unmarshaller config.Unmarshaller,
) (*otelcollector.OTELConfig, error) {
	var pluginConfig pluginconfig.FluxNinjaPluginConfig
	if err := unmarshaller.UnmarshalKey(pluginconfig.PluginConfigKey, &pluginConfig); err != nil {
		return nil, err
	}

	config := otelcollector.NewOTELConfig()
	config.AddDebugExtensions()
	addFluxninjaExporter(config, &pluginConfig, grpcClientConfig, httpClientConfig)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			controllerID := heartbeats.ControllerInfo.Id

			addAttributesProcessor(config, controllerID)

			if logsPipeline, exists := baseConfig.Service.Pipeline("logs"); exists {
				addFNToPipeline("logs", config, logsPipeline)
			}
			if _, exists := baseConfig.Service.Pipeline("metrics/fast"); exists {
				addMetricsSlowPipeline(baseConfig, config)
			}
			if _, exists := baseConfig.Service.Pipeline("metrics/controller-fast"); exists {
				addMetricsControllerSlowPipeline(baseConfig, config)
			}
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})

	return config, nil
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

func addMetricsSlowPipeline(baseConfig, config *otelcollector.OTELConfig) {
	addFluxninjaPrometheusReceiver(baseConfig, config)
	config.AddBatchProcessor(processorBatchMetricsSlow, 10*time.Second, 10000)
	config.Service.AddPipeline("metrics/slow", otelcollector.Pipeline{
		Receivers: []string{receiverPrometheus},
		Processors: []string{
			otelcollector.ProcessorEnrichment,
			processorBatchMetricsSlow,
			processorAttributes,
		},
		Exporters: []string{exporterFluxninja},
	})
}

func addMetricsControllerSlowPipeline(baseConfig, config *otelcollector.OTELConfig) {
	addFluxninjaPrometheusReceiver(baseConfig, config)
	config.AddBatchProcessor(processorBatchMetricsSlow, 10*time.Second, 10000)
	config.Service.AddPipeline("metrics/controller-slow", otelcollector.Pipeline{
		Receivers: []string{receiverPrometheus},
		Processors: []string{
			processorBatchMetricsSlow,
			processorAttributes,
		},
		Exporters: []string{exporterFluxninja},
	})
}

func addFluxninjaPrometheusReceiver(baseConfig, config *otelcollector.OTELConfig) {
	rawReceiverConfig, _ := baseConfig.Receivers[otelcollector.ReceiverPrometheus].(map[string]any)
	duplicatedReceiverConfig, err := duplicateMap(rawReceiverConfig)
	if err != nil {
		// It should not happen, unless the original config is messed up.
		log.Fatal().Err(err).Msg("failed to duplicate config")
	}
	configPatch := map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval": "10s",
			},
		},
	}
	err = mergo.MergeWithOverwrite(&duplicatedReceiverConfig, configPatch)
	if err != nil {
		// It should not happen, unless the original config is messed up.
		log.Fatal().Err(err).Msg("failed to merge configs")
	}
	config.AddReceiver(receiverPrometheus, duplicatedReceiverConfig)
}

func duplicateMap(in map[string]any) (map[string]any, error) {
	rawDuplicatedMap, err := copystructure.Copy(in)
	if err != nil {
		return nil, err
	}
	duplicatedMap, ok := rawDuplicatedMap.(map[string]any)
	if !ok {
		return nil, errors.New("duplicated object not a map")
	}
	return duplicatedMap, nil
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
