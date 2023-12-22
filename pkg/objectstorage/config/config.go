// +kubebuilder:validation:Optional
package config

import "time"

// swagger:operation POST /object_storage common-configuration ObjectStorage
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ObjectStorageConfig"

// BackoffConfig for object storage backoff configuration
// swagger:model
// +kubebuilder:object:generate=true
type BackoffConfig struct {
	Initial    time.Duration `json:"initial" default:"500ms"`
	Multiplier float64       `json:"multiplier" default:"1.5"`
	Maximum    time.Duration `json:"maximum" default:"2s"`
}

// RetryPolicy for object storage retry configuration
// swagger:model
// +kubebuilder:object:generate=true
type RetryPolicy struct {
	Timeout time.Duration `json:"timeout" default:"10s"`
	Backoff BackoffConfig `json:"backoff,omitempty"`
}

// ObjectStorageConfig configures object storage structure.
// swagger:model
// +kubebuilder:object:generate=true
type ObjectStorageConfig struct {
	// Enabled denotes if object storage is enabled.
	Enabled bool `json:"enabled" default:"false"`
	// Backend which provides the object storage.
	Backend string `json:"backend" validate:"oneof=gcs" default:"gcs"`
	// Bucket name of the bucket to use. Required if enabled is true.
	Bucket string `json:"bucket"`
	// KeyPrefix to use when writing to bucket. Required if enabled is true.
	KeyPrefix string `json:"key_prefix"`
	// RetryPolicy to configure retries and timeouts for object storage requests.
	RetryPolicy RetryPolicy `json:"retry_policy,omitempty"`
	// OperationsChannelSize controls size of the channel used for asynchronous puts and deletes.
	OperationsChannelSize int `json:"operations_channel_size" default:"1000"`
}
