package plugins

import "go.uber.org/fx"

// FxPluginIface is an interface for all plugins that provide fx.Option that can be loaded into services, policies, rules etc.
type FxPluginIface interface {
	Module() fx.Option
}

// ServicePluginIface is an interface for all service level plugins.
type ServicePluginIface interface {
	FxPluginIface
}
