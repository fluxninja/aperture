package validation

import (
	"github.com/fluxninja/aperture/pkg/webhooks"
	"go.uber.org/fx"
)

// Module provides fx module for configmap validator.
func Module() fx.Option {
	return fx.Options(
		fx.Invoke(registerCMValidator),
	)
}

// FxIn is a struct that contains all dependencies for configmap validator.
type FxIn struct {
	fx.In
	Webhooks   *webhooks.K8sRegistry
	Validators []CMFileValidator `group:"cm-file-validators"`
}

// registerCMValidator registers configmap validator as k8s webhook.
func registerCMValidator(in FxIn) {
	// The path is not configurable – if one doesn't want default path, one
	// could just write their own Register function
	in.Webhooks.RegisterValidator("/validate/configmap", NewCMValidator(in.Validators))
}
