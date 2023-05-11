package policyvalidator

import (
	"github.com/fluxninja/aperture/v2/pkg/webhooks"
	"go.uber.org/fx"
)

// Module provides fx module for Policy Custom Resource validator.
func Module() fx.Option {
	return fx.Options(
		fx.Invoke(registerPolicyValidator),
	)
}

// FxIn is a struct that contains all dependencies for Policy Custom Resource validator.
type FxIn struct {
	fx.In
	Webhooks   *webhooks.K8sRegistry
	Validators []PolicySpecValidator `group:"policy-validators"`
}

// registerPolicyValidator registers Policy Custom Resource validator as k8s webhook.
func registerPolicyValidator(in FxIn) {
	validators := make([]PolicySpecValidator, 0, len(in.Validators))
	for _, validator := range in.Validators {
		if validator != nil {
			validators = append(validators, validator)
		}
	}

	if len(validators) == 0 {
		return
	}

	// The path is not configurable – if one does not want default path, one
	// could just write their own Register function
	in.Webhooks.RegisterValidator("/validate/policy", NewPolicyValidator(validators))
}
