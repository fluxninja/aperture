package otel_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mitchellh/copystructure"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.uber.org/fx"

	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/fluxninja/v1"
	"github.com/fluxninja/aperture/extensions/fluxninja/extconfig"
	"github.com/fluxninja/aperture/extensions/fluxninja/heartbeats"
	"github.com/fluxninja/aperture/extensions/fluxninja/otel"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	"github.com/fluxninja/aperture/pkg/platform"
)

type inStruct struct {
	fx.In
	Actual []*otelconfig.OTELConfig `group:"extension-config"`
}

var _ = DescribeTable("FN Extension OTEL", func(
	baseConfig *otelconfig.OTELConfig,
	expected *otelconfig.OTELConfig,
) {
	cfg := map[string]interface{}{
		"fluxninja": map[string]interface{}{
			"api_key":  "deadbeef",
			"endpoint": "http://localhost:1234",
		},
	}
	marshalledCfg, err := json.Marshal(cfg)
	Expect(err).NotTo(HaveOccurred())
	unmarshaller, err := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(marshalledCfg)
	Expect(err).NotTo(HaveOccurred())

	var in inStruct
	opts := fx.Options(
		grpcclient.ClientConstructor{Name: "heartbeats-grpc-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.grpc"}.Annotate(),
		httpclient.ClientConstructor{Name: "heartbeats-http-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.http"}.Annotate(),
		extconfig.Module(),
		fx.Provide(
			func() config.Unmarshaller {
				return unmarshaller
			},
			func() *heartbeats.Heartbeats {
				return &heartbeats.Heartbeats{
					ControllerInfo: &heartbeatv1.ControllerInfo{
						Id: "controllero",
					},
				}
			},
			fx.Annotate(
				func() *otelconfig.OTELConfig {
					return baseConfig
				},
				fx.ResultTags(config.NameTag("base")),
			),
		),
		otel.Module(),
		fx.Populate(&in),
	)
	app := platform.New(opts)

	originalBaseConfig, err := copystructure.Copy(baseConfig)
	Expect(err).NotTo(HaveOccurred())

	err = app.Err()
	if err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Error().Err(err).Msg("fx.New failed: " + visualize)
	}

	Expect(err).NotTo(HaveOccurred())

	err = app.Start(context.TODO())
	Expect(err).NotTo(HaveOccurred())

	err = app.Stop(context.TODO())
	Expect(err).NotTo(HaveOccurred())

	// Ensure we did not modify the original base config
	Expect(baseConfig).To(BeEquivalentTo(originalBaseConfig))

	Expect(in.Actual).To(HaveLen(1))
	Expect(in.Actual[0].Receivers).To(Equal(expected.Receivers))
	Expect(in.Actual[0].Processors).To(Equal(expected.Processors))
	Expect(in.Actual[0].Exporters).To(Equal(expected.Exporters))
	Expect(in.Actual[0].Service.Pipelines).To(Equal(expected.Service.Pipelines))
},
	Entry(
		"add FN processors and exporters",
		otelconfig.NewOTELConfig(),
		baseExtensionOTELConfig(),
	),
	Entry(
		"add FN exporters to logs pipeline",
		baseOTELConfigWithPipeline("logs", testPipeline()),
		baseExtensionOTELConfigWithPipeline("logs", testPipelineWithFN()),
	),
	Entry(
		"add FN exporters to alerts pipeline",
		baseOTELConfigWithPipeline("logs/alerts", testPipeline()),
		baseExtensionOTELConfigWithPipeline("logs/alerts", testPipelineWithFN()),
	),
	Entry(
		"add FN exporters to user custom metrics pipeline",
		baseOTELConfigWithPipeline("metrics/user-defined-rabbitmq", testPipeline()),
		baseExtensionOTELConfigWithPipeline("metrics/user-defined-rabbitmq", testPipelineWithFN()),
	),
	Entry(
		"add metrics/slow pipeline if metrics/fast pipeline exists",
		baseOTELConfigWithPipeline("metrics/fast", testPipeline()),
		baseExtensionOTELConfigWithMetrics("metrics/slow"),
	),
	Entry(
		"add metrics/controller-slow pipeline if metrics/controller-fast pipeline exists",
		baseOTELConfigWithPipeline("metrics/controller-fast", testPipeline()),
		baseExtensionOTELConfigWithMetrics("metrics/controller-slow"),
	),
)

func baseOTELConfigWithPipeline(name string, pipeline otelconfig.Pipeline) *otelconfig.OTELConfig {
	cfg := baseOTELConfig()
	cfg.Service.AddPipeline(name, pipeline)
	return cfg
}

func baseExtensionOTELConfigWithPipeline(name string, pipeline otelconfig.Pipeline) *otelconfig.OTELConfig {
	cfg := baseExtensionOTELConfig()
	cfg.Service.AddPipeline(name, pipeline)
	return cfg
}

func baseExtensionOTELConfigWithMetrics(pipelineName string) *otelconfig.OTELConfig {
	cfg := baseExtensionOTELConfig()
	cfg.AddReceiver("prometheus/fluxninja", map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				// Here is different scrape interval than in the base otel config.
				"scrape_interval": "10s",
			},
			"scrape_configs": []string{"foo", "bar"},
		},
	})
	cfg.AddProcessor("batch/metrics-slow", batchprocessor.Config{
		SendBatchSize:    10000,
		SendBatchMaxSize: 10000,
		Timeout:          5 * time.Second,
	})
	processors := []string{
		"batch/metrics-slow",
		"attributes/fluxninja",
	}
	cfg.Service.AddPipeline(pipelineName, otelconfig.Pipeline{
		Receivers:  []string{"prometheus/fluxninja"},
		Processors: processors,
		Exporters:  []string{"otlp/fluxninja"},
	})
	return cfg
}

func baseOTELConfig() *otelconfig.OTELConfig {
	cfg := otelconfig.NewOTELConfig()
	cfg.AddReceiver("prometheus", map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval": "1s",
			},
			// Put some scrape configs to be sure they are not overwritten.
			"scrape_configs": []string{"foo", "bar"},
		},
	})
	return cfg
}

// baseExtensionOTELConfig as produced by FN Extension
func baseExtensionOTELConfig() *otelconfig.OTELConfig {
	cfg := otelconfig.NewOTELConfig()
	cfg.AddProcessor("attributes/fluxninja", map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    "controller_id",
				"action": "insert",
				"value":  "controllero",
			},
		},
	})
	cfg.AddProcessor("transform/fluxninja", map[string]interface{}{
		"log_statements": []map[string]interface{}{
			{
				"context": "resource",
				"statements": []string{
					`set(attributes["controller_id"], "controllero")`,
				},
			},
		},
	})
	cfg.AddExporter("otlp/fluxninja", map[string]interface{}{
		"endpoint": "http://localhost:1234",
		"headers": map[string]interface{}{
			"authorization": "Bearer deadbeef",
		},
		"sending_queue": map[string]interface{}{
			"num_consumers": 1,
		},
		"tls": map[string]interface{}{
			"key_file":             "",
			"ca_file":              "",
			"cert_file":            "",
			"insecure":             false,
			"insecure_skip_verify": false,
		},
	})
	return cfg
}

func testPipelineWithFN() otelconfig.Pipeline {
	p := testPipeline()
	p.Processors = append(p.Processors, "attributes/fluxninja", "transform/fluxninja")
	p.Exporters = append(p.Exporters, "otlp/fluxninja")
	return p
}

func testPipeline() otelconfig.Pipeline {
	return otelconfig.Pipeline{
		Receivers:  []string{"foo"},
		Processors: []string{"bar"},
		Exporters:  []string{"baz"},
	}
}
