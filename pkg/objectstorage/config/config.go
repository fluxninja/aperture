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
	Enabled   bool   `json:"enabled" default:"false"`
	Backend   string `json:"backend" validate:"oneof=gcs" default:"gcs"`
	Bucket    string `json:"bucket" validate:"required"`
	KeyPrefix string `json:"key_prefix" validate:"required"`
}
