package validation

import (
	"github.com/fluxninja/aperture/pkg/webhooks"
	"go.uber.org/fx"
)

// Module provides fx module for configmap validator.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideCMValidator),
		fx.Invoke(registerCMValidator),
	)
}

// provideCMValidator provides config map validator
//
// Note: This validator must be registered to be accessible.
func provideCMValidator() *CMValidator {
	return NewCMValidator()
}

// registerCMValidator registers configmap validator as k8s webhook.
func registerCMValidator(validator *CMValidator, webhooks *webhooks.K8sRegistry) {
	// The path is not configurable – if one doesn't want default path, one
	// could just write their own Register function
	webhooks.RegisterValidator("/validate/configmap", validator)
}
