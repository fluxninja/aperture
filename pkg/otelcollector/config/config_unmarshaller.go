package config

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/confmap"
)

// comply with confmap.Provider interface.
var _ confmap.Provider = (*OTelConfigUnmarshaller)(nil)

const schemeName = "file"

// OTelConfigUnmarshaller can be used as an OTel config map provider.
type OTelConfigUnmarshaller struct {
	lock      sync.Mutex
	config    map[string]interface{}
	watchFunc confmap.WatcherFunc
}

// NewOTelConfigUnmarshaler creates a new OTelConfigUnmarshaler instance.
func NewOTelConfigUnmarshaler(config map[string]interface{}) *OTelConfigUnmarshaller {
	return &OTelConfigUnmarshaller{config: config}
}

// Implements MapProvider interface

// Retrieve returns the value to be injected in the configuration and the corresponding watcher.
func (u *OTelConfigUnmarshaller) Retrieve(_ context.Context, _ string, watchFn confmap.WatcherFunc) (*confmap.Retrieved, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.watchFunc = watchFn
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

// UpdateMap sets the map to the given map.
func (u *OTelConfigUnmarshaller) UpdateMap(config map[string]interface{}) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.config = config
	u.watchFunc(&confmap.ChangeEvent{})
}
