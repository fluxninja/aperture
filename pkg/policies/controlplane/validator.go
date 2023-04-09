package controlplane

import (
	"context"
	"errors"
	"fmt"
	"sync"

	policiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/alerts"
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
		Validator: &PolicySpecValidator{
			fluxmeterNames: make(map[string]bool),
		},
	}
}

// PolicySpecValidator Policy implementation of PolicySpecValidator interface.
type PolicySpecValidator struct {
	mu             sync.Mutex
	fluxmeterNames map[string]bool
}

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
	_, policy, err := ValidateAndCompile(ctx, name, yamlSrc)
	if err != nil {
		// there is no need to handle validator errors. just return validation result.
		return false, err.Error(), nil
	}

	fluxmeters := policy.GetResources().GetFlowControl().GetFluxMeters()
	// Deprecated: v1.5.0
	if fluxmeters == nil {
		fluxmeters = policy.GetResources().GetFluxMeters()
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	newNames := make(map[string]bool)
	for fluxMeterName := range fluxmeters {
		newNames[fluxMeterName] = true
		if v.fluxmeterNames[fluxMeterName] {
			return false, fmt.Sprintf("fluxmeter name \"%s\" already used in other policy", fluxMeterName), nil
		}
	}
	for key, val := range newNames {
		v.fluxmeterNames[key] = val
	}

	return true, "", nil
}

// ValidateAndCompile checks the validity of a single Policy and compiles it.
func ValidateAndCompile(ctx context.Context, name string, yamlSrc []byte) (*circuitfactory.Circuit, *policiesv1.Policy, error) {
	if len(yamlSrc) == 0 {
		return nil, nil, errors.New("empty policy")
	}
	policy := &policiesv1.Policy{}

	err := config.UnmarshalYAML(yamlSrc, policy)
	if err != nil {
		return nil, nil, err
	}

	alerter := alerts.NewSimpleAlerter(100)
	registry := status.NewRegistry(log.GetGlobalLogger(), alerter)
	circuit, err := CompilePolicy(policy, registry)
	if err != nil {
		return nil, nil, err
	}

	if policy.GetResources() != nil {
		classifiers := policy.GetResources().GetFlowControl().GetClassifiers()
		// Deprecated: v1.5.0
		if classifiers == nil {
			classifiers = policy.GetResources().GetClassifiers()
		}
		for _, c := range classifiers {
			_, err = compiler.CompileRuleset(ctx, name, &policysyncv1.ClassifierWrapper{
				Classifier: c,
				ClassifierAttributes: &policysyncv1.ClassifierAttributes{
					PolicyName:      "dummy",
					PolicyHash:      "dummy",
					ClassifierIndex: 0,
				},
			})
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return circuit, policy, nil
}
