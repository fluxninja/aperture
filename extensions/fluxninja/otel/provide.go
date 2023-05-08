package otel

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/mitchellh/copystructure"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/extensions/fluxninja/extconfig"
	"github.com/fluxninja/aperture/extensions/fluxninja/heartbeats"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/utils"
)

const (
	receiverPrometheus = "prometheus/fluxninja"

	processorBatchMetricsSlow   = "batch/metrics-slow"
	processorRollup             = "rollup"
	processorAttributes         = "attributes/fluxninja"
	processorResourceAttributes = "transform/fluxninja"

	exporterFluxNinja = "otlp/fluxninja"
)

// Module provides the OTel configuration for FluxNinja.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideOtelConfig,
				fx.ParamTags(config.NameTag("base"),
					config.NameTag("heartbeats-grpc-client-config"),
					config.NameTag("heartbeats-http-client-config")),
				fx.ResultTags(config.GroupTag("extension-config")),
			),
		),
	)
}

func provideOtelConfig(
	baseConfigProvider *otelconfig.OTelConfigProvider,
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
	lifecycle fx.Lifecycle,
	heartbeats *heartbeats.Heartbeats,
	extensionConfig *extconfig.FluxNinjaExtensionConfig,
) (*otelconfig.OTelConfigProvider, error) {
	if extensionConfig.APIKey == "" {
		return nil, nil
	}

	baseConfig := baseConfigProvider.GetConfig()

	config := otelconfig.NewOTelConfig()
	addFluxNinjaExporter(config, extensionConfig, grpcClientConfig, httpClientConfig)
	// This is to prevent overwriting extensions with empty config.
	config.Extensions = baseConfig.Extensions

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			controllerID := heartbeats.ControllerInfo.Id

			addAttributesProcessor(config, controllerID)
			addResourceAttributesProcessor(config, controllerID)

			if logsPipeline, exists := baseConfig.Service.Pipeline("logs"); exists {
				addFNToPipeline("logs", config, logsPipeline)
			}
			if alertsPipeline, exists := baseConfig.Service.Pipeline("logs/alerts"); exists {
				addFNToPipeline("logs/alerts", config, alertsPipeline)
			}
			if _, exists := baseConfig.Service.Pipeline("metrics/fast"); exists {
				addMetricsSlowPipeline(baseConfig, config)
			}
			if _, exists := baseConfig.Service.Pipeline("metrics/controller-fast"); exists {
				addMetricsControllerSlowPipeline(baseConfig, config)
			}
			for pipelineName, customMetricsPipeline := range baseConfig.Service.Pipelines {
				if !strings.HasPrefix(pipelineName, "metrics/user-defined-") {
					continue
				}
				addFNToPipeline(pipelineName, config, customMetricsPipeline)
			}
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})

	fnConfigProvider := otelconfig.NewOTelConfigProvider("fluxninja", config)
	return fnConfigProvider, nil
}

func addAttributesProcessor(config *otelconfig.OTelConfig, controllerID string) {
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

func addResourceAttributesProcessor(config *otelconfig.OTelConfig, controllerID string) {
	config.AddProcessor(processorResourceAttributes, map[string]interface{}{
		"log_statements": []map[string]interface{}{
			{
				"context": "resource",
				"statements": []string{
					fmt.Sprintf(`set(attributes["%v"], "%v")`, "controller_id", controllerID),
				},
			},
		},
	})
}

func addFNToPipeline(
	name string,
	config *otelconfig.OTelConfig,
	pipeline otelconfig.Pipeline,
) {
	// TODO this duplication of `controller_id` insertion should be cleaned up
	// when telemetry logs pipeline is update to follow the same rules as alerts
	// pipeline.
	pipeline.Processors = append(pipeline.Processors, processorAttributes, processorResourceAttributes)
	pipeline.Exporters = append(pipeline.Exporters, exporterFluxNinja)
	config.Service.AddPipeline(name, pipeline)
}

func addMetricsSlowPipeline(baseConfig, config *otelconfig.OTelConfig) {
	addFluxNinjaPrometheusReceiver(baseConfig, config)
	config.AddBatchProcessor(processorBatchMetricsSlow, 5*time.Second, 10000, 10000)
	config.Service.AddPipeline("metrics/slow", otelconfig.Pipeline{
		Receivers: []string{receiverPrometheus},
		Processors: []string{
			processorBatchMetricsSlow,
			processorAttributes,
		},
		Exporters: []string{exporterFluxNinja},
	})
}

func addMetricsControllerSlowPipeline(baseConfig, config *otelconfig.OTelConfig) {
	addFluxNinjaPrometheusReceiver(baseConfig, config)
	config.AddBatchProcessor(processorBatchMetricsSlow, 5*time.Second, 10000, 10000)
	config.Service.AddPipeline("metrics/controller-slow", otelconfig.Pipeline{
		Receivers: []string{receiverPrometheus},
		Processors: []string{
			processorBatchMetricsSlow,
			processorAttributes,
		},
		Exporters: []string{exporterFluxNinja},
	})
}

func addFluxNinjaPrometheusReceiver(baseConfig, config *otelconfig.OTelConfig) {
	rawReceiverConfig, _ := baseConfig.Receivers[otelconsts.ReceiverPrometheus].(map[string]any)
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

func addFluxNinjaExporter(config *otelconfig.OTelConfig,
	extensionConfig *extconfig.FluxNinjaExtensionConfig,
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
) {
	cfg := map[string]interface{}{
		"endpoint": extensionConfig.Endpoint,
		"headers": map[string]interface{}{
			"authorization": fmt.Sprintf("Bearer %s", extensionConfig.APIKey),
		},
		"sending_queue": map[string]interface{}{
			// Needed to avoid sending metrics out of order, which leads to metrics
			// being dropped.
			"num_consumers": 1,
		},
	}

	var clientTLSConfig tlsconfig.ClientTLSConfig
	tlsCfg := make(map[string]interface{})

	if utils.IsHTTPUrl(extensionConfig.Endpoint) {
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

	config.AddExporter(exporterFluxNinja, cfg)
}
