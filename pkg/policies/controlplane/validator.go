package controlplane

import (
	"context"
	"errors"

	policiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/circuitfactory"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier/compiler"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/fluxninja/aperture/pkg/webhooks/policyvalidator"
	"go.uber.org/fx"
)

// FxOut is the output of the controlplane module.
type FxOut struct {
	fx.Out
	Validator policyvalidator.PolicySpecValidator `group:"policy-validators"`
}

// providePolicyValidator provides classification Policy Custom Resource validator
//
// Note: This validator must be registered to be accessible.
func providePolicyValidator() FxOut {
	return FxOut{
		Validator: &PolicySpecValidator{},
	}
}

// PolicySpecValidator Policy implementation of PolicySpecValidator interface.
type PolicySpecValidator struct{}

// ValidateSpec checks the validity of a Policy spec
//
// returns:
// * true, "", nil when Policy is valid
// * false, message, nil when Policy is invalid
// and
// * false, "", err on other errors.
//
// ValidateSpec checks the syntax, validity of extractors, and validity of
// rego modules (by attempting to compile them).
func (v *PolicySpecValidator) ValidateSpec(
	ctx context.Context,
	name string,
	yamlSrc []byte,
) (bool, string, error) {
	_, valid, msg, err := ValidateAndCompile(ctx, name, yamlSrc)
	return valid, msg, err
}

// ValidateAndCompile checks the validity of a single Policy and compiles it.
func ValidateAndCompile(ctx context.Context, name string, yamlSrc []byte) (*circuitfactory.Circuit, bool, string, error) {
	if len(yamlSrc) == 0 {
		return nil, false, "Empty yaml", nil
	}
	policy := &policiesv1.Policy{}

	err := config.UnmarshalYAML(yamlSrc, policy)
	if err != nil {
		return nil, false, err.Error(), nil
	}

	registry := status.NewRegistry(log.GetGlobalLogger())
	circuit, err := CompilePolicy(policy, registry)
	if err != nil {
		return nil, false, err.Error(), err
	}

	if policy.GetResources() != nil {
		for _, c := range policy.GetResources().Classifiers {
			_, err = compiler.CompileRuleset(ctx, name, &policysyncv1.ClassifierWrapper{
				Classifier: c,
				ClassifierAttributes: &policysyncv1.ClassifierAttributes{
					PolicyName:      "dummy",
					PolicyHash:      "dummy",
					ClassifierIndex: 0,
				},
			})
			if err != nil {
				if errors.Is(err, compiler.BadExtractor) || errors.Is(err, compiler.BadSelector) ||
					errors.Is(err, compiler.BadRego) || errors.Is(err, compiler.BadLabelName) {
					return nil, false, err.Error(), nil
				} else {
					return nil, false, err.Error(), err
				}
			}
		}
	}
	return circuit, true, "", nil
}
