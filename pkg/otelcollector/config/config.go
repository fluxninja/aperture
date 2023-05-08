// +kubebuilder:validation:Optional
package config

import (
	"crypto/tls"
	"fmt"
	"sort"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
)

// OTelConfig represents OTel Collector configuration.
type OTelConfig struct {
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Receivers  map[string]interface{} `json:"receivers,omitempty"`
	Processors map[string]interface{} `json:"processors,omitempty"`
	Exporters  map[string]interface{} `json:"exporters,omitempty"`
	Connectors map[string]interface{} `json:"connectors,omitempty"`
	Service    *OTelService           `json:"service"`
}

// NewOTelConfig creates new empty OTelConfig.
func NewOTelConfig() *OTelConfig {
	return &OTelConfig{
		Extensions: map[string]interface{}{},
		Receivers:  map[string]interface{}{},
		Processors: map[string]interface{}{},
		Exporters:  map[string]interface{}{},
		Connectors: map[string]interface{}{},
		Service:    NewOTelService(),
	}
}

// AsMap returns map representation of OTelConfig.
func (o *OTelConfig) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"extensions": o.Extensions,
		"receivers":  o.Receivers,
		"processors": o.Processors,
		"exporters":  o.Exporters,
		"connectors": o.Connectors,
		"service":    o.Service.AsMap(),
	}
}

// AddExtension adds given extension and enables it in service.
func (o *OTelConfig) AddExtension(name string, value interface{}) {
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

// AddReceiver adds receiver to OTel config.
func (o *OTelConfig) AddReceiver(name string, value interface{}) {
	o.Receivers[name] = value
}

// AddProcessor adds processor to OTel config.
func (o *OTelConfig) AddProcessor(name string, value interface{}) {
	o.Processors[name] = value
}

// AddExporter adds exporter to OTel config.
func (o *OTelConfig) AddExporter(name string, value interface{}) {
	o.Exporters[name] = value
}

// AddConnector adds connector to OTel config.
func (o *OTelConfig) AddConnector(name string, value interface{}) {
	o.Connectors[name] = value
}

// SetDebugPort configures debug port on which OTel server /metrics as specified by user.
func (o *OTelConfig) SetDebugPort(userCfg *CommonOTelConfig) {
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
func (o *OTelConfig) AddDebugExtensions(userCfg *CommonOTelConfig) {
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
func (o *OTelConfig) AddBatchProcessor(
	name string,
	timeout time.Duration,
	sendBatchSize uint32,
	sendBatchMaxSize uint32,
) {
	// Note: Not passing batchprocessor.Config struct to avoid depending on batchprocessor.
	o.AddProcessor(name, map[string]interface{}{
		"timeout":             timeout,
		"send_batch_size":     sendBatchSize,
		"send_batch_max_size": sendBatchMaxSize,
	})
}

// OTelService represents service in OTel Config.
type OTelService struct {
	Telemetry  map[string]interface{}
	Pipelines  map[string]Pipeline
	Extensions []string
}

// NewOTelService returns new empty OTel Service.
func NewOTelService() *OTelService {
	return &OTelService{
		Telemetry: map[string]interface{}{
			"logs": map[string]interface{}{
				"level": "INFO",
			},
		},
		Pipelines:  map[string]Pipeline{},
		Extensions: []string{},
	}
}

// AsMap returns map representation of OTelService.
func (o *OTelService) AsMap() map[string]interface{} {
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

// AddPipeline adds pipeline to OTel Service.
func (o *OTelService) AddPipeline(name string, pipeline Pipeline) {
	o.Pipelines[name] = pipeline
}

// Pipeline gets pipeline with given name from OTel Service together with `exists` bool.
func (o *OTelService) Pipeline(name string) (Pipeline, bool) {
	pipeline, exists := o.Pipelines[name]
	return pipeline, exists
}

// Pipeline represents OTel Config pipeline.
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

/* Specific to Agent and Controller OTel collector factories. */

// BaseFxTag is the base name tag for otel components.
var (
	BaseFxTag               = config.NameTag("base")
	TelemetryCollectorFxTag = config.NameTag("telemetry-collector")
)

// FxIn consumes parameters via Fx.
type FxIn struct {
	fx.In
	Unmarshaller    config.Unmarshaller
	Listener        *listener.Listener
	PromClient      promapi.Client
	TLSConfig       *tls.Config
	ServerTLSConfig tlsconfig.ServerTLSConfig
}
