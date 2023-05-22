package config

import (
	"context"
	"sync"

	"github.com/fluxninja/aperture/v2/pkg/log"
	"go.opentelemetry.io/collector/confmap"
)

// comply with confmap.Provider interface.
var _ confmap.Provider = (*OTelConfigProvider)(nil)

// OTelConfigProvider can be used as an OTel config map provider.
type OTelConfigProvider struct {
	lock      sync.Mutex  // protects config & watchFunc
	config    *OTelConfig // nil only after Shutdown.
	watchFunc confmap.WatcherFunc
	scheme    string
}

// NewOTelConfigProvider creates a new OTelConfigUnmarshaler instance.
func NewOTelConfigProvider(scheme string, config *OTelConfig) *OTelConfigProvider {
	p := &OTelConfigProvider{
		scheme: scheme,
		config: NewOTelConfig(),
	}
	p.UpdateConfig(config)
	return p
}

// Implements MapProvider interface

// Retrieve returns the value to be injected in the configuration and the corresponding watcher.
func (u *OTelConfigProvider) Retrieve(_ context.Context, _ string, watchFn confmap.WatcherFunc) (*confmap.Retrieved, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.watchFunc = watchFn
	return confmap.NewRetrieved(u.config.AsMap())
}

// Shutdown indicates the provider should close.
func (u *OTelConfigProvider) Shutdown(ctx context.Context) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	// Prevent UpdatesConfig to run after Shutdown.
	u.watchFunc = nil
	u.config = nil
	return nil
}

// Scheme returns the scheme name, location scheme used by Retrieve.
func (u *OTelConfigProvider) Scheme() string {
	return u.scheme
}

// UpdateConfig sets the map to the given map.
func (u *OTelConfigProvider) UpdateConfig(config *OTelConfig) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if config == nil {
		config = NewOTelConfig()
	}
	if u.config == nil {
		log.Warn().Msg("OtelConfigProvider.UpdateConfig called after Shutdown")
		return
	}
	u.config = config
	if u.watchFunc != nil {
		u.watchFunc(&confmap.ChangeEvent{})
	}
}

// GetConfig returns the current config.
func (u *OTelConfigProvider) GetConfig() *OTelConfig {
	u.lock.Lock()
	defer u.lock.Unlock()
	return u.config
}
