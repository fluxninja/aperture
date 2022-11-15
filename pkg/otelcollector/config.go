// +kubebuilder:validation:Optional
package otelcollector

import (
	"context"
	"crypto/tls"
	"fmt"
	"sort"
	"time"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.uber.org/fx"
	"k8s.io/client-go/rest"
)

// OTELConfigUnmarshaller can be used as an OTEL config map provider.
type OTELConfigUnmarshaller struct {
	config map[string]interface{}
}

// NewOTELConfigUnmarshaler creates a new OTELConfigUnmarshaler instance.
func NewOTELConfigUnmarshaler(config map[string]interface{}) *OTELConfigUnmarshaller {
	return &OTELConfigUnmarshaller{config: config}
}

// Implements MapProvider interface

// Retrieve returns the value to be injected in the configuration and the corresponding watcher.
func (u *OTELConfigUnmarshaller) Retrieve(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
	return confmap.NewRetrieved(u.config)
}

// Shutdown indicates the provider should close.
func (u *OTELConfigUnmarshaller) Shutdown(ctx context.Context) error {
	return nil
}

// Scheme returns the scheme name, location scheme used by Retrieve.
func (u *OTELConfigUnmarshaller) Scheme() string {
	return schemeName
}

// OTELConfig represents OTEL Collector configuration.
type OTELConfig struct {
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Receivers  map[string]interface{} `json:"receivers,omitempty"`
	Processors map[string]interface{} `json:"processors,omitempty"`
	Exporters  map[string]interface{} `json:"exporters,omitempty"`
	Service    *OTELService           `json:"service"`
}

// NewOTELConfig creates new empty OTELConfig.
func NewOTELConfig() *OTELConfig {
	return &OTELConfig{
		Extensions: map[string]interface{}{},
		Receivers:  map[string]interface{}{},
		Processors: map[string]interface{}{},
		Exporters:  map[string]interface{}{},
		Service:    NewOTELService(),
	}
}

// AsMap returns map representation of OTELConfig.
func (o *OTELConfig) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"extensions": o.Extensions,
		"receivers":  o.Receivers,
		"processors": o.Processors,
		"exporters":  o.Exporters,
		"service":    o.Service.AsMap(),
	}
}

// AddExtension adds given extension and enables it in service.
func (o *OTELConfig) AddExtension(name string, value interface{}) {
	if value == nil {
		value = map[string]interface{}{}
	}
	o.Extensions[name] = value
	extensions := make([]string, 0, len(o.Extensions))
	for extension := range o.Extensions {
		extensions = append(extensions, extension)
	}
	sort.Strings(extensions)
	o.Service.Extensions = extensions
}

// AddReceiver adds receiver to OTEL config.
func (o *OTELConfig) AddReceiver(name string, value interface{}) {
	o.Receivers[name] = value
}

// AddProcessor adds receiver to OTEL config.
func (o *OTELConfig) AddProcessor(name string, value interface{}) {
	o.Processors[name] = value
}

// AddExporter adds receiver to OTEL config.
func (o *OTELConfig) AddExporter(name string, value interface{}) {
	o.Exporters[name] = value
}

// SetDebugPort configures debug port on which OTEL server /metrics as specified by user.
func (o *OTELConfig) SetDebugPort(userCfg *OtelConfig) {
	portInput := fmt.Sprintf(":%d", userCfg.Ports.DebugPort)
	if val, ok := o.Service.Telemetry["metrics"]; ok {
		if val, ok := val.(map[string]interface{}); ok {
			val["address"] = portInput
			o.Service.Telemetry["metrics"] = val
		}
	} else {
		addressMap := make(map[string]interface{})
		addressMap["address"] = portInput
		o.Service.Telemetry["metrics"] = addressMap
	}
}

// AddDebugExtensions adds common debug extensions and enables them.
func (o *OTELConfig) AddDebugExtensions(userCfg *OtelConfig) {
	o.AddExtension("health_check", map[string]interface{}{
		"endpoint": fmt.Sprintf("localhost:%d", userCfg.Ports.HealthCheckPort),
	})
	o.AddExtension("pprof", map[string]interface{}{
		"endpoint": fmt.Sprintf("localhost:%d", userCfg.Ports.PprofPort),
	})
	o.AddExtension("zpages", map[string]interface{}{
		"endpoint": fmt.Sprintf("localhost:%d", userCfg.Ports.ZpagesPort),
	})
}

// AddBatchProcessor is a helper function for adding batch processor.
func (o *OTELConfig) AddBatchProcessor(
	name string,
	timeout time.Duration,
	sendBatchSize uint32,
	sendBatchMaxSize uint32,
) {
	o.AddProcessor(name, batchprocessor.Config{
		Timeout:          timeout,
		SendBatchSize:    sendBatchSize,
		SendBatchMaxSize: sendBatchMaxSize,
	})
}

// OTELService represents service in OTEL Config.
type OTELService struct {
	Telemetry  map[string]interface{}
	Pipelines  map[string]Pipeline
	Extensions []string
}

// NewOTELService returns new empty OTEL Service.
func NewOTELService() *OTELService {
	return &OTELService{
		Telemetry: map[string]interface{}{
			"logs": map[string]interface{}{
				"level": "INFO",
			},
		},
		Pipelines:  map[string]Pipeline{},
		Extensions: []string{},
	}
}

// AsMap returns map representation of OTELService.
func (o *OTELService) AsMap() map[string]interface{} {
	pipelines := map[string]interface{}{}
	for name, pipe := range o.Pipelines {
		pipelines[name] = pipe.AsMap()
	}
	return map[string]interface{}{
		"telemetry":  o.Telemetry,
		"extensions": o.Extensions,
		"pipelines":  pipelines,
	}
}

// AddPipeline adds pipeline to OTEL Service.
func (o *OTELService) AddPipeline(name string, pipeline Pipeline) {
	o.Pipelines[name] = pipeline
}

// Pipeline gets pipeline with given name from OTEL Service together with `exists` bool.
func (o *OTELService) Pipeline(name string) (Pipeline, bool) {
	pipeline, exists := o.Pipelines[name]
	return pipeline, exists
}

// Pipeline represents OTEL Config pipeline.
type Pipeline struct {
	Receivers  []string
	Processors []string
	Exporters  []string
}

// AsMap returns map representation of Pipeline.
func (p *Pipeline) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"receivers":  p.Receivers,
		"processors": p.Processors,
		"exporters":  p.Exporters,
	}
}

/* Specific to Agent and Controller OTEL collector factories. */

// BaseFxTag is the base name tag for otel components.
var BaseFxTag = config.NameTag("base")

// OtelParams contains parameters for otel collector factories for agent and controller.
type OtelParams struct {
	promClient promapi.Client
	Config     *OTELConfig
	Listener   *listener.Listener
	tlsConfig  *tls.Config
	OtelConfig
}

// swagger:operation POST /otel common-configuration Otel
// ---
// x-fn-config-env: true
// parameters:
// - name: proxy
//   in: body
//   schema:
//     "$ref": "#/definitions/OtelConfig"

// OtelConfig is the configuration for the OTEL collector.
// swagger:model
// +kubebuilder:object:generate=true
type OtelConfig struct {
	// BatchPrerollup configures batch prerollup processor.
	BatchPrerollup BatchPrerollupConfig `json:"batch_prerollup"`
	// BatchPostrollup configures batch postrollup processor.
	BatchPostrollup BatchPostrollupConfig `json:"batch_postrollup"`
	// BatchAlerts configures batch alerts processor.
	BatchAlerts BatchAlertsConfig `json:"batch_alerts"`
	// Ports configures debug, health and extension ports values.
	Ports PortsConfig `json:"ports"`
}

// PortsConfig defines configuration for OTEL debug and extension ports.
// swagger:model
// +kubebuilder:object:generate=true
type PortsConfig struct {
	// Port on which otel collector exposes prometheus metrics on /metrics path.
	DebugPort uint32 `json:"debug_port" validate:"gte=0" default:"8888"`
	// Port on which health check extension in exposed.
	HealthCheckPort uint32 `json:"health_check_port" validate:"gte=0" default:"13133"`
	// Port on which pprof extension in exposed.
	PprofPort uint32 `json:"pprof_port" validate:"gte=0" default:"1777"`
	// Port on which zpages extension in exposed.
	ZpagesPort uint32 `json:"zpages_port" validate:"gte=0" default:"55679"`
}

// BatchPrerollupConfig defines configuration for OTEL batch processor.
// swagger:model
// +kubebuilder:object:generate=true
type BatchPrerollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"10000"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"10000"`
}

// BatchPostrollupConfig defines configuration for OTEL batch processor.
// swagger:model
// +kubebuilder:object:generate=true
type BatchPostrollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"100"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"100"`
}

// BatchAlertsConfig defines configuration for OTEL batch processor.
// swagger:model
// +kubebuilder:object:generate=true
type BatchAlertsConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"100"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"100"`
}

// FxIn consumes parameters via Fx.
type FxIn struct {
	fx.In
	Unmarshaller    config.Unmarshaller
	Listener        *listener.Listener
	PromClient      promapi.Client
	TLSConfig       *tls.Config
	ServerTLSConfig tlsconfig.ServerTLSConfig
}

// NewOtelConfig returns OTEL parameters for OTEL collectors.
func NewOtelConfig(in FxIn) (*OtelParams, error) {
	config := NewOTELConfig()

	var userCfg OtelConfig
	if err := in.Unmarshaller.UnmarshalKey("otel", &userCfg); err != nil {
		return nil, err
	}

	config.SetDebugPort(&userCfg)
	config.AddDebugExtensions(&userCfg)

	cfg := &OtelParams{
		OtelConfig: userCfg,
		Listener:   in.Listener,
		promClient: in.PromClient,
		tlsConfig:  in.TLSConfig,
		Config:     config,
	}
	return cfg, nil
}

// NewDefaultOtelConfig creates OtelConfig with all the default values set.
func NewDefaultOtelConfig() *OtelConfig {
	return &OtelConfig{
		Ports: PortsConfig{
			DebugPort:       8888,
			HealthCheckPort: 13133,
			PprofPort:       1777,
			ZpagesPort:      55679,
		},
	}
}

// AddMetricsPipeline adds metrics to pipeline for agent OTEL collector.
func AddMetricsPipeline(cfg *OtelParams) {
	config := cfg.Config
	addPrometheusReceiver(cfg)
	config.AddProcessor(ProcessorEnrichment, nil)
	addPrometheusRemoteWriteExporter(config, cfg.promClient)
	config.Service.AddPipeline("metrics/fast", Pipeline{
		Receivers: []string{ReceiverPrometheus},
		Processors: []string{
			ProcessorEnrichment,
			ProcessorAgentGroup,
		},
		Exporters: []string{ExporterPrometheusRemoteWrite},
	})
}

// AddControllerMetricsPipeline adds metrics to pipeline for controller OTEL collector.
func AddControllerMetricsPipeline(cfg *OtelParams) {
	config := cfg.Config
	addControllerPrometheusReceiver(config, cfg)
	addPrometheusRemoteWriteExporter(config, cfg.promClient)
	config.Service.AddPipeline("metrics/controller-fast", Pipeline{
		Receivers:  []string{ReceiverPrometheus},
		Processors: []string{},
		Exporters:  []string{ExporterPrometheusRemoteWrite},
	})
}

// AddAlertsPipeline adds reusable alerts pipeline.
func AddAlertsPipeline(cfg *OtelParams, extraProcessors ...string) {
	config := cfg.Config
	config.AddReceiver(ReceiverAlerts, map[string]any{})
	config.AddBatchProcessor(
		ProcessorBatchAlerts,
		cfg.BatchAlerts.Timeout.AsDuration(),
		cfg.BatchAlerts.SendBatchSize,
		cfg.BatchAlerts.SendBatchMaxSize,
	)
	processors := []string{ProcessorBatchAlerts}
	processors = append(processors, extraProcessors...)
	config.Service.AddPipeline("logs/alerts", Pipeline{
		Receivers:  []string{ReceiverAlerts},
		Processors: processors,
		Exporters:  []string{ExporterLogging},
	})
}

func addPrometheusReceiver(cfg *OtelParams) {
	config := cfg.Config
	scrapeConfigs := []map[string]any{
		buildApertureSelfScrapeConfig("aperture-self", cfg),
	}

	_, err := rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		log.Debug().Msg("K8s environment not detected. Skipping K8s scrape configurations.")
	} else if err != nil {
		log.Warn().Err(err).Msg("Error when discovering k8s environment")
	} else {
		log.Debug().Msg("K8s environment detected. Adding K8s scrape configurations.")
		scrapeConfigs = append(scrapeConfigs, buildKubernetesNodesScrapeConfig(cfg), buildKubernetesPodsScrapeConfig(cfg))
	}

	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddReceiver(ReceiverPrometheus, map[string]any{
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

func addControllerPrometheusReceiver(config *OTELConfig, cfg *OtelParams) {
	scrapeConfigs := []map[string]any{
		buildApertureSelfScrapeConfig("aperture-controller-self", cfg),
	}
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddReceiver(ReceiverPrometheus, map[string]any{
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

func addPrometheusRemoteWriteExporter(config *OTELConfig, promClient promapi.Client) {
	endpoint := promClient.URL("api/v1/write", nil)
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddExporter(ExporterPrometheusRemoteWrite, map[string]any{
		"endpoint": endpoint.String(),
	})
}

func buildApertureSelfScrapeConfig(name string, cfg *OtelParams) map[string]any {
	scheme := "http"
	if cfg.tlsConfig != nil {
		scheme = "https"
	}
	return map[string]any{
		"job_name": name,
		"scheme":   scheme,
		"tls_config": map[string]any{
			"insecure_skip_verify": true,
		},
		"metrics_path": "/metrics",
		"static_configs": []map[string]any{
			{
				"targets": []string{cfg.Listener.GetAddr()},
				"labels": map[string]any{
					metrics.InstanceLabel:    info.Hostname,
					metrics.ProcessUUIDLabel: info.UUID,
				},
			},
		},
	}
}

func buildKubernetesNodesScrapeConfig(cfg *OtelParams) map[string]any {
	return map[string]any{
		"job_name":     "kubernetes-nodes",
		"scheme":       "https",
		"metrics_path": "/metrics/cadvisor",
		"authorization": map[string]any{
			"credentials_file": "/var/run/secrets/kubernetes.io/serviceaccount/token",
		},
		"tls_config": map[string]any{
			"insecure_skip_verify": true,
		},
		"kubernetes_sd_configs": []map[string]any{
			{"role": "node"},
		},
		"relabel_configs": []map[string]any{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_node_name"},
				"action":        "keep",
				"regex":         info.Hostname,
			},
		},
		"metric_relabel_configs": []map[string]any{
			{
				"source_labels": []string{"__name__"},
				"action":        "keep",
				"regex":         "container_memory_working_set_bytes|container_spec_memory_limit_bytes|container_spec_cpu_(?:quota|period)|container_cpu_usage_seconds_total",
			},
			{
				"source_labels": []string{"pod"},
				"action":        "replace",
				"target_label":  "entity_name",
			},
		},
	}
}

func buildKubernetesPodsScrapeConfig(cfg *OtelParams) map[string]any {
	return map[string]any{
		"job_name":     "kubernetes-pods",
		"scheme":       "http",
		"metrics_path": "/metrics",
		"kubernetes_sd_configs": []map[string]any{
			{"role": "pod"},
		},
		"relabel_configs": []map[string]any{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_pod_node_name"},
				"action":        "keep",
				"regex":         info.Hostname,
			},
			// Scrape only pods which have github.com/fluxninja/scrape=true annotation.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_aperture_tech_scrape"},
				"action":        "keep",
				"regex":         "true",
			},
			// Allow rewrite of scheme, path and port where prometheus metrics are served.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
				"action":        "replace",
				"regex":         "(https?)",
				"target_label":  "__scheme__",
			},
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
				"action":        "replace",
				"target_label":  "__metrics_path__",
				"regex":         "(.+)",
			},
			{
				"source_labels": []string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
				"action":        "replace",
				"regex":         `([^:]+)(?::\d+)?;(\d+)`,
				"replacement":   "$$1:$$2",
				"target_label":  "__address__",
			},
		},
		"metric_relabel_configs": []map[string]any{
			// For now, dropping everything. In future, we'll want to filter in some
			// metrics based on policies. See #4632.
			{
				"source_labels": []string{},
				"action":        "drop",
			},
		},
	}
}
