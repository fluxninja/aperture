package config

import (
	"context"

	"go.opentelemetry.io/collector/confmap"
)

const schemeName = "file"

// OTelConfigUnmarshaller can be used as an OTel config map provider.
type OTelConfigUnmarshaller struct {
	config map[string]interface{}
}

// NewOTelConfigUnmarshaler creates a new OTelConfigUnmarshaler instance.
func NewOTelConfigUnmarshaler(config map[string]interface{}) *OTelConfigUnmarshaller {
	return &OTelConfigUnmarshaller{config: config}
}

// Implements MapProvider interface

// Retrieve returns the value to be injected in the configuration and the corresponding watcher.
func (u *OTelConfigUnmarshaller) Retrieve(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
	return confmap.NewRetrieved(u.config)
}

// Shutdown indicates the provider should close.
func (u *OTelConfigUnmarshaller) Shutdown(ctx context.Context) error {
	return nil
}

// Scheme returns the scheme name, location scheme used by Retrieve.
func (u *OTelConfigUnmarshaller) Scheme() string {
	return schemeName
}
