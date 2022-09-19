package otel_test

import (
	"context"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.uber.org/fx"

	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/heartbeats"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/otel"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
)

type inStruct struct {
	fx.In
	Actual []*otelcollector.OTELConfig `group:"plugin-config"`
}

var _ = DescribeTable("FN Plugin OTEL", func(
	baseConfig *otelcollector.OTELConfig,
	expected *otelcollector.OTELConfig,
) {
	cfg := map[string]interface{}{
		"fluxninja_plugin": map[string]interface{}{
			"api_key":            "deadbeef",
			"fluxninja_endpoint": "http://localhost:1234",
		},
	}
	marshalledCfg, err := json.Marshal(cfg)
	Expect(err).NotTo(HaveOccurred())
	unmarshaller, err := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(marshalledCfg)
	Expect(err).NotTo(HaveOccurred())

	var in inStruct
	opts := fx.Options(
		grpcclient.ClientConstructor{Name: "heartbeats-grpc-client", ConfigKey: pluginconfig.PluginConfigKey + ".client.grpc"}.Annotate(),
		httpclient.ClientConstructor{Name: "heartbeats-http-client", ConfigKey: pluginconfig.PluginConfigKey + ".client.http"}.Annotate(),
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
				func() *otelcollector.OTELConfig {
					return baseConfig
				},
				fx.ResultTags(config.NameTag("base")),
			),
		),
		otel.Module(),
		fx.Populate(&in),
	)
	app := platform.New(opts)

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

	Expect(in.Actual).To(HaveLen(1))
	Expect(in.Actual[0].Receivers).To(Equal(expected.Receivers))
	Expect(in.Actual[0].Processors).To(Equal(expected.Processors))
	Expect(in.Actual[0].Exporters).To(Equal(expected.Exporters))
	Expect(in.Actual[0].Service.Pipelines).To(Equal(expected.Service.Pipelines))
},
	Entry(
		"add FN processors and exporters",
		otelcollector.NewOTELConfig(),
		basePluginOTELConfig(),
	),
	Entry(
		"add FN exporters to logs pipeline",
		baseOTELConfigWithPipeline("logs", testPipeline()),
		basePluginOTELConfigWithPipeline("logs", testPipelineWithFN()),
	),
	Entry(
		"add FN exporters to traces pipeline",
		baseOTELConfigWithPipeline("traces", testPipeline()),
		basePluginOTELConfigWithPipeline("traces", testPipelineWithFN()),
	),
	Entry(
		"add metrics/slow pipeline if metrics/fast pipeline exists",
		baseOTELConfigWithPipeline("metrics/fast", testPipeline()),
		basePluginOTELConfigWithMetrics("metrics/slow"),
	),
	Entry(
		"add metrics/controller-slow pipeline if metrics/controller-fast pipeline exists",
		baseOTELConfigWithPipeline("metrics/controller-fast", testPipeline()),
		basePluginOTELConfigWithMetrics("metrics/controller-slow"),
	),
)

func baseOTELConfigWithPipeline(name string, pipeline otelcollector.Pipeline) *otelcollector.OTELConfig {
	cfg := baseOTELConfig()
	cfg.Service.AddPipeline(name, pipeline)
	return cfg
}

func basePluginOTELConfigWithPipeline(name string, pipeline otelcollector.Pipeline) *otelcollector.OTELConfig {
	cfg := basePluginOTELConfig()
	cfg.Service.AddPipeline(name, pipeline)
	return cfg
}

func basePluginOTELConfigWithMetrics(pipelineName string) *otelcollector.OTELConfig {
	cfg := basePluginOTELConfig()
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
		SendBatchSize: 10000,
		Timeout:       10 * time.Second,
	})
	processors := []string{
		"batch/metrics-slow",
		"attributes/fluxninja",
	}
	if pipelineName == "metrics/slow" {
		processors = append([]string{"enrichment"}, processors...)
	}
	cfg.Service.AddPipeline(pipelineName, otelcollector.Pipeline{
		Receivers:  []string{"prometheus/fluxninja"},
		Processors: processors,
		Exporters:  []string{"otlp/fluxninja"},
	})
	return cfg
}

func baseOTELConfig() *otelcollector.OTELConfig {
	cfg := otelcollector.NewOTELConfig()
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

// basePluginOTELConfig as produced by FN plugin
func basePluginOTELConfig() *otelcollector.OTELConfig {
	cfg := otelcollector.NewOTELConfig()
	cfg.AddProcessor("attributes/fluxninja", map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    "controller_id",
				"action": "insert",
				"value":  "controllero",
			},
		},
	})
	cfg.AddExporter("otlp/fluxninja", map[string]interface{}{
		"endpoint": "http://localhost:1234",
		"headers": map[string]interface{}{
			"authorization": "Bearer deadbeef",
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

func testPipelineWithFN() otelcollector.Pipeline {
	p := testPipeline()
	p.Processors = append(p.Processors, "attributes/fluxninja")
	p.Exporters = append(p.Exporters, "otlp/fluxninja")
	return p
}

func testPipeline() otelcollector.Pipeline {
	return otelcollector.Pipeline{
		Receivers:  []string{"foo"},
		Processors: []string{"bar"},
		Exporters:  []string{"baz"},
	}
}
