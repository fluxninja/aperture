package otelcollector

import (
	"context"
	"sort"
	"time"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/processor/batchprocessor"
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

// AddDebugExtensions adds common debug extextions and enables them.
func (o *OTELConfig) AddDebugExtensions() {
	o.AddExtension("health_check", nil)
	o.AddExtension("pprof", map[string]interface{}{
		"endpoint": "localhost:1777",
	})
	o.AddExtension("zpages", map[string]interface{}{
		"endpoint": "localhost:55679",
	})
}

// AddBatchProcessor is a helper function for adding batch processor.
func (o *OTELConfig) AddBatchProcessor(name string, timeout time.Duration, sendBatchSize uint32) {
	o.AddProcessor(name, batchprocessor.Config{
		Timeout:       timeout,
		SendBatchSize: sendBatchSize,
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
