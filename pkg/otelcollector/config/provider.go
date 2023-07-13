package otelconfig

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/collector/confmap"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

// comply with confmap.Provider interface.
var _ confmap.Provider = (*Provider)(nil)

// Provider is an OTel config map provider.
//
// It allows updating the config and registering hooks.
type Provider struct {
	configLock    sync.RWMutex // protects config, watchFunc & hooks
	watchFuncLock sync.Mutex   // protects watchFunc
	config        *Config      // nil only after Shutdown.
	watchFunc     confmap.WatcherFunc
	hooks         []func(*Config)
	scheme        string
}

// NewProvider creates a new OTelConfigProvider.
func NewProvider(scheme string, config *Config) *Provider {
	p := &Provider{
		scheme: scheme,
		config: New(),
	}
	p.UpdateConfig(config)
	return p
}

// Retrieve implements confmap.Provider.
func (p *Provider) Retrieve(
	_ context.Context,
	_ string,
	watchFn confmap.WatcherFunc,
) (*confmap.Retrieved, error) {
	c := p.getConfig()
	if c == nil {
		log.Bug().Msg("Retrieve after Shutdown")
		return nil, errors.New("already shut down")
	}

	p.setWatchFunc(watchFn)
	return confmap.NewRetrieved(p.config.AsMap())
}

// Shutdown implements confmap.Provider.
func (p *Provider) Shutdown(ctx context.Context) error {
	// Prevent UpdateConfig to run after Shutdown.
	p.setWatchFunc(nil)
	p.setConfig(nil)
	return nil
}

// Scheme implements confmap.Provider.
func (p *Provider) Scheme() string { return p.scheme }

// UpdateConfig sets the new config, replacing the old one.
// Before new config is set, hooks are allowed to modify the config.
// Collector update is triggered asynchronously.
//
// Note: Caller should not use the passed config object after calling this function.
func (p *Provider) UpdateConfig(config *Config) {
	if config == nil {
		// Maintain the p.config not-nil invariant.
		config = New()
	}

	for _, hook := range p.hooks {
		hook(config)
	}

	p.setConfigIfNotNil(config)
	wf := p.getWatchFunc()
	if wf != nil {
		wf(&confmap.ChangeEvent{})
	}
}

// AddMutatingHook adds a hook to be run before applying config.
//
// The hook should treat the given config as temporary.
// The hook will also be executed immediately, to ensure that current config
// was passed through all the added hooks.
//
// WARNING: This is supposed to be called only during initialization.
func (p *Provider) AddMutatingHook(hook func(*Config)) {
	if p.config == nil {
		log.Warn().Msg("OtelConfigProvider.AddHook: already shut down")
		return
	}

	p.hooks = append(p.hooks, hook)

	p.UpdateConfig(p.config)
}

func (p *Provider) getConfig() *Config {
	p.configLock.RLock()
	defer p.configLock.RUnlock()
	return p.config
}

func (p *Provider) setConfig(config *Config) {
	p.configLock.Lock()
	defer p.configLock.Unlock()
	p.config = config
}

func (p *Provider) setConfigIfNotNil(config *Config) {
	p.configLock.Lock()
	defer p.configLock.Unlock()
	if p.config == nil {
		log.Warn().Msg("OtelConfigProvider: tried to update config after Shutdown")
		return
	}
	p.config = config
}

func (p *Provider) getWatchFunc() confmap.WatcherFunc {
	p.watchFuncLock.Lock()
	defer p.watchFuncLock.Unlock()
	return p.watchFunc
}

func (p *Provider) setWatchFunc(watchFunc confmap.WatcherFunc) {
	p.watchFuncLock.Lock()
	defer p.watchFuncLock.Unlock()
	p.watchFunc = watchFunc
}
