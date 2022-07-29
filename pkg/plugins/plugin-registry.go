package plugins

import (
	"path"
	"path/filepath"
	"plugin"

	"github.com/spf13/pflag"
	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/filesystem"
	"github.com/FluxNinja/aperture/pkg/info"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/utils"
)

const (
	defaultKey = "plugins"

	// ServicePluginSymbol is symbol that Service level plugins must expose to be loaded.
	// Example usage:
	//
	// func ServicePlugin() ServicePluginIface.
	//
	ServicePluginSymbol = "ServicePlugin"
)

var (
	pluginPrefix    = info.Prefix + "-plugin-"
	pluginGlob      = pluginPrefix + "*.so"
	defaultPath     = path.Join(config.DefaultArtifactsDirectory, "plugins")
	defaultPathFlag = defaultKey + ".plugins_path"
)

// PluginRegistry holds fields used for internal tracking of plugin symbols and disabled symbols or plugins of the service.
type PluginRegistry struct {
	allPlugins      pluginSymbolTrackers
	disabledSymbols []string
	disabledPlugins []string
}

// pluginSymbolTracker tracks plugins for a type.
type pluginSymbolTracker struct {
	plugins PluginTrackers
}

// pluginSymbolTrackers tracks types.
type pluginSymbolTrackers map[string]*pluginSymbolTracker

// PluginTracker tracks single plugin.
type PluginTracker struct {
	FileInfo *filesystem.FileInfo
	Plugin   *plugin.Plugin
	Symbol   plugin.Symbol
}

// PluginTrackers tracks plugin name to plugin.
type PluginTrackers map[string]*PluginTracker

// swagger:operation POST /plugins common-configuration Plugins
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/PluginsConfig"

// PluginsConfig holds configuration for plugins.
// swagger:model
type PluginsConfig struct {
	// Path to plugins directory. This can be set via command line arguments as well.
	PluginsPath string `json:"plugins_path"`
	// Specific plugin types to disable
	DisabledSymbols []string `json:"disabled_symbols"`
	// Specific plugins to disable
	DisabledPlugins []string `json:"disabled_plugins"`
	// Disables all plugins
	DisablePlugins bool `json:"disable_plugins" default:"false"`
}

// Constructor holds fields for constructing a PluginRegistry.
type Constructor struct {
	// Config key
	Key string
	// Path key
	PathKey string
	// Plugin Symbols to look for
	PluginSymbols []string
	// default config
	DefaultConfig PluginsConfig
}

func (constructor Constructor) setFlags(fs *pflag.FlagSet) error {
	fs.String(constructor.PathKey, defaultPath, "path to plugins")
	return nil
}

func (constructor Constructor) provideFlagSetBuilder() config.FlagSetBuilderOut {
	return config.FlagSetBuilderOut{Builder: constructor.setFlags}
}

// ModuleConfig holds configuration for the plugin module.
type ModuleConfig struct {
	PluginSymbols        []string
	OnlyCommandLineFlags bool
}

// Module is a fx module that provides flag set builder and new plugin registry.
func (config ModuleConfig) Module() fx.Option {
	constructor := Constructor{
		Key:           defaultKey,
		PathKey:       defaultPathFlag,
		PluginSymbols: config.PluginSymbols,
	}
	options := fx.Options(
		fx.Provide(constructor.provideFlagSetBuilder),
	)
	if !config.OnlyCommandLineFlags {
		options = fx.Options(
			options,
			fx.Provide(constructor.newPluginRegistry),
		)
	}
	return options
}

func (constructor Constructor) newPluginRegistry(unmarshaller config.Unmarshaller) (*PluginRegistry, fx.Option, error) {
	var pluginOptions fx.Option
	config := constructor.DefaultConfig
	if err := unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize plugin configuration")
		return nil, nil, err
	}

	registry := &PluginRegistry{}
	registry.allPlugins = make(pluginSymbolTrackers, 0)
	registry.disabledSymbols = config.DisabledSymbols
	registry.disabledPlugins = config.DisabledPlugins

	constructor.PluginSymbols = append(constructor.PluginSymbols, ServicePluginSymbol)
	var pluginSymbols []string
	for _, s := range constructor.PluginSymbols {
		if !registry.isSymbolDisabled(s) {
			pluginSymbols = append(pluginSymbols, s)
		}
	}

	if !config.DisablePlugins {
		// Discover Plugins
		var pluginPaths []string
		var err error
		// Make the directory absolute if it isn't already
		if !filepath.IsAbs(config.PluginsPath) {
			config.PluginsPath, err = filepath.Abs(config.PluginsPath)
			if err != nil {
				return nil, nil, err
			}
		}
		pluginPaths, err = filepath.Glob(filepath.Join(config.PluginsPath, pluginGlob))
		if err != nil {
			return nil, nil, err
		}
		log.Debug().Strs("plugins", pluginPaths).Msg("discovered plugins")
		// Find types for these plugins
		for _, pluginPath := range pluginPaths {
			finfo := filesystem.ParseFilePath(pluginPath)
			pluginName := finfo.GetFileName()
			if registry.isPluginDisabled(pluginName) {
				log.Debug().Str("plugin", pluginName).Msg("not loading plugin as it is disabled")
				continue
			}

			plugin, err := plugin.Open(pluginPath)
			if err != nil {
				log.Warn().Err(err).Str("plugin", pluginName).Msg("unable to open plugin")
				continue
			}
			for _, pluginSymbol := range pluginSymbols {
				// lookup symbol
				symbol, err := plugin.Lookup(pluginSymbol)
				if err != nil {
					continue
				}
				log.Debug().Str("symbol", pluginSymbol).Str("plugin", pluginName).Msg("symbol found")

				symbolTracker, ok := registry.allPlugins[pluginSymbol]
				if !ok {
					symbolTracker = &pluginSymbolTracker{}
					symbolTracker.plugins = make(PluginTrackers, 0)
					registry.allPlugins[pluginSymbol] = symbolTracker
				}
				pluginTracker := &PluginTracker{}
				pluginTracker.FileInfo = finfo
				pluginTracker.Plugin = plugin
				pluginTracker.Symbol = symbol
				symbolTracker.plugins[pluginName] = pluginTracker
			}
		}
		// provide fx.Option for service-level plugins
		pluginOptions = registry.GetServicePluginOptions()
	}
	return registry, pluginOptions, nil
}

func (registry *PluginRegistry) isSymbolDisabled(symbolName string) bool {
	return utils.SliceContains(registry.disabledSymbols, symbolName)
}

func (registry *PluginRegistry) isPluginDisabled(pluginName string) bool {
	return utils.SliceContains(registry.disabledPlugins, pluginName)
}

// GetServicePluginOptions returns plugin options for all the service-level plugins via trackers.
func (registry *PluginRegistry) GetServicePluginOptions() fx.Option {
	var pluginOptions fx.Option
	// search symbol
	pluginSymbolTracker, ok := registry.allPlugins[ServicePluginSymbol]
	if !ok {
		return nil
	}
	for _, tracker := range pluginSymbolTracker.plugins {
		symbol := tracker.Symbol
		var options fx.Option
		fxPlugin := symbol.(func() ServicePluginIface)()
		options = fxPlugin.Module()
		if options != nil {
			if pluginOptions == nil {
				pluginOptions = options
			} else {
				pluginOptions = fx.Options(pluginOptions, options)
			}
		}
	}
	return pluginOptions
}

// GetPluginTracker returns the tracker for the plugin.
func (registry *PluginRegistry) GetPluginTracker(symbolName string, pluginName string) (*PluginTracker, bool) {
	var ok bool
	var pluginSymbolTracker *pluginSymbolTracker
	var pluginTracker *PluginTracker
	pluginSymbolTracker, ok = registry.allPlugins[symbolName]
	if !ok {
		return nil, false
	}
	pluginTracker, ok = pluginSymbolTracker.plugins[pluginName]
	return pluginTracker, ok
}
