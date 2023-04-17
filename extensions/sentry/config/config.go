// +kubebuilder:validation:Optional
package config

// SentryConfig holds configuration for Sentry.
// swagger:model
// +kubebuilder:object:generate=true
type SentryConfig struct {
	// If DSN is not set, the client is effectively disabled
	// You can set test project's DSN to send log events.
	// oss-aperture project DSN is set as default.
	Dsn string `json:"dsn" default:"https://6223f112b0ac4344aa67e94d1631eb85@o574197.ingest.sentry.io/6605877"`
	// Environment
	Environment string `json:"environment" default:"production"`
	// Sample rate for sampling traces
	TracesSampleRate float64 `json:"traces_sample_rate" default:"0.2" validate:"gte=0,lte=1"`
	// Sample rate for event submission
	SampleRate float64 `json:"sample_rate" default:"1.0" validate:"gte=0,lte=1"`
	// Debug enables printing of Sentry SDK debug messages
	Debug bool `json:"debug" default:"true"`
	// Configure to generate and attach stack traces to capturing message calls
	AttachStacktrace bool `json:"attach_stack_trace" default:"true"`
	// Sentry crash report disabled
	Disabled bool `json:"disabled" default:"false"`
}
