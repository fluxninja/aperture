package config

import (
	"context"

	"go.opentelemetry.io/collector/confmap"
)

const schemeName = "file"

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
