package otel

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/mitchellh/copystructure"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/extensions/fluxninja/extconfig"
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/heartbeats"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	grpcclient "github.com/fluxninja/aperture/v2/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/v2/pkg/net/http"
	"github.com/fluxninja/aperture/v2/pkg/net/tlsconfig"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

const (
	receiverPrometheus = "prometheus/fluxninja"

	processorBatchMetricsSlow   = "batch/metrics-slow"
	processorRollup             = "rollup"
	processorAttributes         = "attributes/fluxninja"
	processorResourceAttributes = "transform/fluxninja"

	exporterFluxNinja = "otlp/fluxninja"

	scrapeInterval = "10s"
	// BatchTimeout needs to be less than `scrapeInterval` no there is only
	// one data point for given metrics in each batch.
	batchTimeout     = 5 * time.Second
	sendBatchSize    = 10000
	sendBatchSizeMax = 10000
)

// Module provides the OTel configuration for FluxNinja.
func Module() fx.Option {
	return fx.Options(
		fx.Invoke(
			fx.Annotate(
				injectOtelConfig,
				fx.ParamTags(
					config.NameTag("heartbeats-grpc-client-config"),
					config.NameTag("heartbeats-http-client-config"),
				),
			),
		),
	)
}

func injectOtelConfig(
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
	lifecycle fx.Lifecycle,
	heartbeats *heartbeats.Heartbeats,
	extensionConfig *extconfig.FluxNinjaExtensionConfig,
	configProvider *otelconfig.Provider,
) {
	//nolint:staticcheck // SA1019 read APIKey config for backward compatibility
	if extensionConfig.AgentAPIKey == "" && extensionConfig.APIKey == "" {
		return
	}

	lifecycle.Append(fx.StartHook(func() {
		// We must add the config-modifying hook in this start hook, as we need
		// heartbeats.ControllerInfo, that's populated only in heartbeats start.
		configProvider.AddMutatingHook(func(config *otelconfig.Config) {
			addFluxNinjaExporter(config, extensionConfig, grpcClientConfig, httpClientConfig)

			controllerInfoPtr := heartbeats.GetControllerInfoPtr()
			if controllerInfoPtr == nil {
				log.Info().Msg("Heartbeats controller info not set, skipping")
				return
			}

			controllerID := controllerInfoPtr.Id

			addAttributesProcessor(config, controllerID)
			addResourceAttributesProcessor(config, controllerID)

			if logsPipeline, exists := config.Service.Pipeline("logs"); exists {
				addFNToPipeline("logs", config, logsPipeline)
			}
			if alertsPipeline, exists := config.Service.Pipeline("logs/alerts"); exists {
				addFNToPipeline("logs/alerts", config, alertsPipeline)
			}

			disableLocalPipelines := extensionConfig.EnableCloudController || extensionConfig.DisableLocalOTelPipeline

			if _, exists := config.Service.Pipeline("metrics/fast"); exists {
				addMetricsSlowPipeline(config)
				if disableLocalPipelines {
					deleteLocalMetricsPipeline(config, "metrics/fast")
				}
			}
			if _, exists := config.Service.Pipeline("metrics/controller-fast"); exists {
				addMetricsControllerSlowPipeline(config)
				if disableLocalPipelines {
					deleteLocalMetricsPipeline(config, "metrics/controller-fast")
				}
			}
			for pipelineName, customMetricsPipeline := range config.Service.Pipelines {
				if !strings.HasPrefix(pipelineName, "metrics/user-defined-") {
					continue
				}
				if disableLocalPipelines {
					// In case of user defined pipelines, we clean-up list of exporters
					// and then depend on `addFNToPipeline' to add its own exporter.
					customMetricsPipeline.Exporters = []string{}
				}
				addFNToPipeline(pipelineName, config, customMetricsPipeline)
			}
		})
	}))
}

func addAttributesProcessor(config *otelconfig.Config, controllerID string) {
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

func addResourceAttributesProcessor(config *otelconfig.Config, controllerID string) {
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
	config *otelconfig.Config,
	pipeline otelconfig.Pipeline,
) {
	log.Info().Str("pipeline", name).Msg("Adding fluxninja exporter to pipeline")
	// TODO this duplication of `controller_id` insertion should be cleaned up
	// when telemetry logs pipeline is update to follow the same rules as alerts
	// pipeline.
	pipeline.Processors = append(pipeline.Processors, processorAttributes, processorResourceAttributes)
	pipeline.Exporters = append(pipeline.Exporters, exporterFluxNinja)
	config.Service.AddPipeline(name, pipeline)
}

func addMetricsSlowPipeline(config *otelconfig.Config) {
	log.Info().Msg("Adding metrics/slow pipeline")
	addFluxNinjaPrometheusReceiver(config)
	config.AddBatchProcessor(processorBatchMetricsSlow, batchTimeout, sendBatchSize, sendBatchSizeMax)
	config.Service.AddPipeline("metrics/slow", otelconfig.Pipeline{
		Receivers: []string{receiverPrometheus},
		Processors: []string{
			processorBatchMetricsSlow,
			processorAttributes,
			otelconsts.ProcessorAgentGroup,
		},
		Exporters: []string{exporterFluxNinja},
	})
}

func addMetricsControllerSlowPipeline(config *otelconfig.Config) {
	log.Info().Msg("Adding metrics/controller-slow pipeline")
	addFluxNinjaPrometheusReceiver(config)
	config.AddBatchProcessor(processorBatchMetricsSlow, batchTimeout, sendBatchSize, sendBatchSizeMax)
	config.Service.AddPipeline("metrics/controller-slow", otelconfig.Pipeline{
		Receivers: []string{receiverPrometheus},
		Processors: []string{
			processorBatchMetricsSlow,
			processorAttributes,
		},
		Exporters: []string{exporterFluxNinja},
	})
}

func addFluxNinjaPrometheusReceiver(config *otelconfig.Config) {
	rawReceiverConfig, _ := config.Receivers[otelconsts.ReceiverPrometheus].(map[string]any)
	duplicatedReceiverConfig, err := duplicateMap(rawReceiverConfig)
	if err != nil {
		// It should not happen, unless the original config is messed up.
		log.Fatal().Err(err).Msg("failed to duplicate config")
	}
	configPatch := map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval": scrapeInterval,
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

func deleteLocalMetricsPipeline(config *otelconfig.Config, pipeline string) {
	log.Info().Msg("Cleaning up local prometheus pipeline")
	config.Service.DeletePipeline(pipeline)
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

func addFluxNinjaExporter(config *otelconfig.Config,
	extensionConfig *extconfig.FluxNinjaExtensionConfig,
	grpcClientConfig *grpcclient.GRPCClientConfig,
	httpClientConfig *httpclient.HTTPClientConfig,
) {
	apiKey := extensionConfig.APIKey
	if apiKey == "" {
		//nolint:staticcheck // SA1019 read AgentAPIKey config for backward compatibility
		apiKey = extensionConfig.AgentAPIKey
	}
	cfg := map[string]interface{}{
		"endpoint": extensionConfig.Endpoint,
		"headers": map[string]interface{}{
			"authorization": fmt.Sprintf("Bearer %s", apiKey),
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
