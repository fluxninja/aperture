package tracestologsprocessor

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the Traces to Logs processor.
type Config struct {
	LogsExporter string `mapstructure:"logs_exporter"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid.
func (cfg *Config) Validate() error {
	return nil
}
