package enrichmentprocessor

import (
	"github.com/fluxninja/aperture/pkg/entitycache"
)

// Config holds the configuration for the enrichment processor.
type Config struct {
	entityCache *entitycache.EntityCache
}
