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
func (o *OTELConfig) SetDebugPort(userCfg *CommonOTELConfig) {
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
func (o *OTELConfig) AddDebugExtensions(userCfg *CommonOTELConfig) {
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
	// Note: Not passing batchprocessor.Config struct to avoid depending on batchprocessor.
	o.AddProcessor(name, map[string]interface{}{
		"timeout":             timeout,
		"send_batch_size":     sendBatchSize,
		"send_batch_max_size": sendBatchMaxSize,
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

// FxIn consumes parameters via Fx.
type FxIn struct {
	fx.In
	Unmarshaller    config.Unmarshaller
	Listener        *listener.Listener
	PromClient      promapi.Client
	TLSConfig       *tls.Config
	ServerTLSConfig tlsconfig.ServerTLSConfig
}
