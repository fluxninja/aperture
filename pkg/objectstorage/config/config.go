// +kubebuilder:validation:Optional
package config

// swagger:operation POST /object_storage common-configuration ObjectStorage
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/Config"

// Config for object storage.
// swagger:model
// +kubebuilder:object:generate=true
type Config struct {
	// Enabled denotes if object storage is enabled.
	Enabled bool `json:"enabled" default:"false"`
	// Backend which provides the object storage.
	Backend string `json:"backend" validate:"oneof=gcs" default:"gcs"`
	// Bucket name of the bucket to use. Required if enabled is true.
	Bucket string `json:"bucket"`
	// KeyPrefix to use when writing to bucket. Required if enabled is true.
	KeyPrefix string `json:"key_prefix"`
}
