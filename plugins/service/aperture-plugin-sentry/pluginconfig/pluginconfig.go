package pluginconfig

const (
	// PluginConfigKey is the key for the plugin configuration.
	PluginConfigKey = "sentry_plugin"
)

// SentryPluginConfig is the configuration for Sentry cloud integration plugin.
// swagger:model
type SentryPluginConfig struct {
	// API Key for this agent.
	// Send Agent Key as label to sentry events.
	SentryAgentKey string `json:"sentry_agent_key"`
}
