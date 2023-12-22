// +kubebuilder:validation:Optional
package config

// swagger:operation POST /object_storage common-configuration ObjectStorage
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ObjectStorageConfig"

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
	// OperationsChannelSize controls size of the channel used for asynchronous puts and deletes.
	OperationsChannelSize int `json:"operations_channel_size" default:"1000"`
}
