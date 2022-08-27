package validation

import "github.com/fluxninja/aperture/pkg/webhooks"

// ProvideCMValidator provides config map validator
//
// Note: This validator must be registered to be accessible.
func ProvideCMValidator() *CMValidator {
	return NewCMValidator()
}

// RegisterCMValidator registers configmap validator as k8s webhook.
func RegisterCMValidator(validator *CMValidator, webhooks *webhooks.K8sRegistry) {
	// The path is not configurable – if one doesn't want default path, one
	// could just write their own Register function
	webhooks.RegisterValidator("/validate/configmap", validator)
}
