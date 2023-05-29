package otelconfig

import (
	"context"
	"errors"
	"sync"

	"github.com/fluxninja/aperture/v2/pkg/log"
	"go.opentelemetry.io/collector/confmap"
)

// comply with confmap.Provider interface.
var _ confmap.Provider = (*Provider)(nil)

// Provider is an OTel config map provider.
//
// It allows updating the config and registering hooks.
type Provider struct {
	lock      sync.Mutex // protects config, watchFunc & hooks
	config    *Config    // nil only after Shutdown.
	watchFunc confmap.WatcherFunc
	hooks     []func(*Config)
	scheme    string
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
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.config == nil {
		log.Bug().Msg("Retrieve after Shutdown")
		return nil, errors.New("already shut down")
	}

	p.watchFunc = watchFn
	return confmap.NewRetrieved(p.config.AsMap())
}

// Shutdown implements confmap.Provider.
func (p *Provider) Shutdown(ctx context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	// Prevent UpdateConfig to run after Shutdown.
	p.watchFunc = nil
	p.config = nil
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

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, hook := range p.hooks {
		hook(config)
	}

	p.setConfig(config)
}

// Set post-hooks config and trigger update, assuming p.lock locked.
func (p *Provider) setConfig(config *Config) {
	if p.config == nil {
		log.Warn().Msg("OtelConfigProvider: tried to update config after Shutdown")
		return
	}

	p.config = config
	if p.watchFunc != nil {
		p.watchFunc(&confmap.ChangeEvent{})
	}
}

// AddMutatingHook adds a hook to be run before applying config.
//
// The hook should treat the given config as temporary.
// The hook will also be executed immediately, to ensure that current config
// was passed through all the added hooks.
func (p *Provider) AddMutatingHook(hook func(*Config)) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.config == nil {
		log.Warn().Msg("OtelConfigProvider.AddHook: already shut down")
		return
	}

	p.hooks = append(p.hooks, hook)

	// Now the config provided to collector is outdated, run the newly-added hook.
	hook(p.config)
	p.setConfig(p.config)
}

// MustGetConfig returns a snapshot of the current config.
func (p *Provider) MustGetConfig() *Config {
	p.lock.Lock()
	defer p.lock.Unlock()
	// Copying to avoid concurrent modification by hooks.
	return p.config.MustCopy()
}
