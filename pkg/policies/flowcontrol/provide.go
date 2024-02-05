package flowcontrol

import (
	"go.uber.org/fx"

	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/labelstatus"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/fluxmeter"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service"
	servicegetter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service-getter"
)

// Module returns the fx options for dataplane side pieces of policy.
func Module() fx.Option {
	return fx.Options(
		actuators.Module(),
		fluxmeter.Module(),
		classifier.Module(),
		service.Module(),
		servicegetter.Module,
		EngineModule(),
		CacheModule(),
		labelstatus.Module(),
	)
}

// EngineModule returns the fx options for the engine.
func EngineModule() fx.Option {
	return fx.Provide(ProvideEngine)
}

// ProvideEngine provides the engine for the dataplane side of policy.
func ProvideEngine(cache iface.Cache, agentInfo *agentinfo.AgentInfo) iface.Engine {
	engine := NewEngine(agentInfo)
	engine.RegisterCache(cache)
	return engine
}

// CacheModule returns the fx options for the cache.
func CacheModule() fx.Option {
	return fx.Provide(NewCache)
}
