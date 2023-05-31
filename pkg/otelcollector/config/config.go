// +kubebuilder:validation:Optional
package otelconfig

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/mitchellh/copystructure"
)

// Config represents OTel Collector configuration.
type Config struct {
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Receivers  map[string]interface{} `json:"receivers,omitempty"`
	Processors map[string]interface{} `json:"processors,omitempty"`
	Exporters  map[string]interface{} `json:"exporters,omitempty"`
	Connectors map[string]interface{} `json:"connectors,omitempty"`
	Service    *Service               `json:"service"`
}

// New creates new empty Config.
func New() *Config {
	return &Config{
		Extensions: map[string]interface{}{},
		Receivers:  map[string]interface{}{},
		Processors: map[string]interface{}{},
		Exporters:  map[string]interface{}{},
		Connectors: map[string]interface{}{},
		Service:    NewService(),
	}
}

// AsMap returns map representation of Config.
func (o *Config) AsMap() map[string]interface{} {
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
func (o *Config) AddExtension(name string, value interface{}) {
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
func (o *Config) AddReceiver(name string, value interface{}) {
	o.Receivers[name] = value
}

// AddProcessor adds processor to OTel config.
func (o *Config) AddProcessor(name string, value interface{}) {
	o.Processors[name] = value
}

// AddExporter adds exporter to OTel config.
func (o *Config) AddExporter(name string, value interface{}) {
	o.Exporters[name] = value
}

// AddConnector adds connector to OTel config.
func (o *Config) AddConnector(name string, value interface{}) {
	o.Connectors[name] = value
}

// SetDebugPort configures debug port on which OTel server /metrics as specified by user.
func (o *Config) SetDebugPort(userCfg *CommonOTelConfig) {
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
func (o *Config) AddDebugExtensions(userCfg *CommonOTelConfig) {
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
func (o *Config) AddBatchProcessor(
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

// Copy returns a deep copy of the config.
//
// This should error only in pathological cases.
func (o *Config) Copy() (*Config, error) {
	copyInterface, err := copystructure.Copy(o)
	if err != nil {
		return nil, err
	}

	copy, ok := copyInterface.(*Config)
	if !ok {
		return nil, errors.New("copy has wrong type")
	}
	return copy, nil
}

// MustCopy returns a deep copy of the config or panics.
func (o *Config) MustCopy() *Config {
	copy, err := o.Copy()
	if err != nil {
		log.Panic().Err(err).Msg("cannot copy config")
	}
	return copy
}

// Service represents service in OTel Config.
type Service struct {
	Telemetry  map[string]interface{}
	Pipelines  map[string]Pipeline
	Extensions []string
}

// NewService returns new empty OTel Service.
func NewService() *Service {
	return &Service{
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
func (o *Service) AsMap() map[string]interface{} {
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
func (o *Service) AddPipeline(name string, pipeline Pipeline) {
	o.Pipelines[name] = pipeline
}

// Pipeline gets pipeline with given name from OTel Service together with `exists` bool.
func (o *Service) Pipeline(name string) (Pipeline, bool) {
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
