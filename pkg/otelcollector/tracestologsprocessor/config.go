package tracestologsprocessor

import (
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for the Traces to Logs processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	LogsExporter             string                   `mapstructure:"logs_exporter"`
}

var _ config.Processor = (*Config)(nil)

// Validate checks if the exporter configuration is valid.
func (cfg *Config) Validate() error {
	return nil
}
