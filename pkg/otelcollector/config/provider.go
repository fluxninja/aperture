package otelconfig

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/confmap"
)

// comply with confmap.Provider interface.
var _ confmap.Provider = (*Provider)(nil)

// Provider is an OTel config map provider.
//
// It allows updating the config and registering hooks.
type Provider struct {
	watchFunc confmap.WatcherFunc
	scheme    string
	hooks     []func(*Config)
	lock      sync.Mutex // protects watchFunc & hooks
}

// NewProvider creates a new OTelConfigProvider.
func NewProvider(scheme string) *Provider {
	p := &Provider{
		scheme: scheme,
	}
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

	p.watchFunc = watchFn

	config := p.getConfig()

	return confmap.NewRetrieved(config.AsMap())
}

// GetConfig returns the current config.
func (p *Provider) GetConfig() *Config {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.getConfig()
}

func (p *Provider) getConfig() *Config {
	config := New()
	for _, hook := range p.hooks {
		hook(config)
	}
	return config
}

// Shutdown implements confmap.Provider.
func (p *Provider) Shutdown(ctx context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	// Prevent UpdateConfig to run after Shutdown.
	p.watchFunc = nil
	return nil
}

// Scheme implements confmap.Provider.
func (p *Provider) Scheme() string { return p.scheme }

// UpdateConfig triggers Collector update asynchronously.
func (p *Provider) UpdateConfig() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.notifyWatchFunc()
}

func (p *Provider) notifyWatchFunc() {
	// Set post-hooks config and trigger update, assuming p.lock locked.
	if p.watchFunc != nil {
		p.watchFunc(&confmap.ChangeEvent{})
		// Prevent calling watchFunc another time before the next Retrieve
		// (perhaps it's even illegal). We won't miss any update though, as
		// explained below.
		p.watchFunc = nil
	}
	// If watchFunc is nil, then:
	// * either we have not any Retrieve yet, so there's no reason to notify, or
	// * we already notified of the change, but didn't got Retrieve yet.
}

// AddMutatingHook adds a hook to be run before applying config.
func (p *Provider) AddMutatingHook(hook func(*Config)) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.hooks = append(p.hooks, hook)

	p.notifyWatchFunc()
}
