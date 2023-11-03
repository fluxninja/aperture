package otel_test

import (
	"context"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"

	heartbeatv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/fluxninja/v1"
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/extconfig"
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/heartbeats"
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/otel"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	grpcclient "github.com/fluxninja/aperture/v2/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/v2/pkg/net/http"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	"github.com/fluxninja/aperture/v2/pkg/platform"
)

var _ = DescribeTable("FN Extension OTel", func(
	baseConfig *otelconfig.Config,
	expected *otelconfig.Config,
) {
	cfg := map[string]interface{}{
		"fluxninja": map[string]interface{}{
			"agent_api_key": "deadbeef",
			"endpoint":      "http://localhost:1234",
		},
	}
	marshalledCfg, err := json.Marshal(cfg)
	Expect(err).NotTo(HaveOccurred())
	unmarshaller, err := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(marshalledCfg)
	Expect(err).NotTo(HaveOccurred())

	configProvider := otelconfig.NewProvider("base")
	configProvider.AddMutatingHook(func(otelCfg *otelconfig.Config) {
		*otelCfg = *baseConfig
	})

	heartbeats := &heartbeats.Heartbeats{}
	heartbeats.SetControllerInfoPtr(
		&heartbeatv1.ControllerInfo{
			Id: "controllero",
		})
	opts := fx.Options(
		grpcclient.ClientConstructor{Name: "heartbeats-grpc-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.grpc"}.Annotate(),
		httpclient.ClientConstructor{Name: "heartbeats-http-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.http"}.Annotate(),
		extconfig.Module(),
		fx.Provide(
			func() config.Unmarshaller {
				return unmarshaller
			},
		),
		fx.Supply(
			heartbeats,
		),
		fx.Supply(configProvider),
		otel.Module(),
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

	config := configProvider.GetConfig()
	Expect(config.Receivers).To(Equal(expected.Receivers))
	Expect(config.Processors).To(Equal(expected.Processors))
	Expect(config.Exporters).To(Equal(expected.Exporters))
	Expect(config.Service.Pipelines).To(Equal(expected.Service.Pipelines))

	err = app.Stop(context.TODO())
	Expect(err).NotTo(HaveOccurred())
},
	Entry(
		"add FN processors and exporters",
		otelconfig.New(),
		configBuilder{cfg: otelconfig.New()}.withExtensionConfig().cfg,
	),
	Entry(
		"add FN exporters to logs pipeline",
		baseOTelConfig().withPipeline("logs", testPipeline()).cfg,
		baseOTelConfig().withPipeline("logs", testPipelineWithFN()).withExtensionConfig().cfg,
	),
	Entry(
		"add FN exporters to alerts pipeline",
		baseOTelConfig().withPipeline("logs/alerts", testPipeline()).cfg,
		baseOTelConfig().withPipeline("logs/alerts", testPipelineWithFN()).withExtensionConfig().cfg,
	),
	Entry(
		"add FN exporters to user custom metrics pipeline",
		baseOTelConfig().withPipeline("metrics/user-defined-rabbitmq", testPipeline()).cfg,
		baseOTelConfig().withPipeline("metrics/user-defined-rabbitmq", testPipelineWithFN()).withExtensionConfig().cfg,
	),
	Entry(
		"add metrics/slow pipeline if metrics/fast pipeline exists",
		baseOTelConfig().withPipeline("metrics/fast", testPipeline()).cfg,
		baseOTelConfig().withPipeline("metrics/fast", testPipeline()).withMetrics("metrics/slow").withExtensionConfig().cfg,
	),
	Entry(
		"add metrics/controller-slow pipeline if metrics/controller-fast pipeline exists",
		baseOTelConfig().withPipeline("metrics/controller-fast", testPipeline()).cfg,
		baseOTelConfig().withPipeline("metrics/controller-fast", testPipeline()).withMetrics("metrics/controller-slow").withExtensionConfig().cfg,
	),
)

// configBuilder is wrapper around otelconfig.Config with builder-style methods.
type configBuilder struct{ cfg *otelconfig.Config }

func (b configBuilder) withPipeline(name string, pipeline otelconfig.Pipeline) configBuilder {
	b.cfg.Service.AddPipeline(name, pipeline)
	return b
}

func (b configBuilder) withMetrics(pipelineName string) configBuilder {
	b.cfg.AddReceiver("prometheus/fluxninja", map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				// Here is different scrape interval than in the base otel config.
				"scrape_interval": "10s",
			},
			"scrape_configs": []string{"foo", "bar"},
		},
	})
	b.cfg.AddProcessor("batch/metrics-slow", map[string]any{
		"send_batch_size":     uint32(10000),
		"send_batch_max_size": uint32(10000),
		"timeout":             5 * time.Second,
	})
	processors := []string{
		"batch/metrics-slow",
		"attributes/fluxninja",
	}
	b.cfg.Service.AddPipeline(pipelineName, otelconfig.Pipeline{
		Receivers:  []string{"prometheus/fluxninja"},
		Processors: processors,
		Exporters:  []string{"otlp/fluxninja"},
	})
	return b
}

func baseOTelConfig() configBuilder {
	cfg := otelconfig.New()
	cfg.AddReceiver("prometheus", map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval": "10s",
			},
			// Put some scrape configs to be sure they are not overwritten.
			"scrape_configs": []string{"foo", "bar"},
		},
	})
	return configBuilder{cfg: cfg}
}

// appends basic processors / exporters unconditionally added by fluxninja extension
func (b configBuilder) withExtensionConfig() configBuilder {
	b.cfg.AddProcessor("attributes/fluxninja", map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    "controller_id",
				"action": "insert",
				"value":  "controllero",
			},
		},
	})
	b.cfg.AddProcessor("transform/fluxninja", map[string]interface{}{
		"log_statements": []map[string]interface{}{
			{
				"context": "resource",
				"statements": []string{
					`set(attributes["controller_id"], "controllero")`,
				},
			},
		},
	})
	b.cfg.AddExporter("otlp/fluxninja", map[string]interface{}{
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
	return b
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
