package enrichmentprocessor

import (
	"go.opentelemetry.io/collector/config"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

// Config holds the configuration for the enrichment processor.
type Config struct {
	entityCache              *entitycache.EntityCache
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
}
