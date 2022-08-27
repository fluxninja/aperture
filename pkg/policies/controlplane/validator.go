package controlplane

import (
	"context"
	"errors"

	policiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier"
	"github.com/fluxninja/aperture/pkg/webhooks/validation"
)

// ProvideCMFileValidator provides classification config map file validator
//
// Note: This validator must be registered to be accessible.
func ProvideCMFileValidator() *CMFileValidator {
	return &CMFileValidator{}
}

// RegisterCMFileValidator registers classification configmap validator as configmap file validator.
func RegisterCMFileValidator(validator *CMFileValidator, configMapValidator *validation.CMValidator) {
	// The path is not configurable – if one doesn't want default path, one
	// could just write their own Register function
	configMapValidator.RegisterCMFileValidator(validator)
}

// CMFileValidator Policy implementation of CMFileValidator interface.
type CMFileValidator struct{}

// CheckCMName checks configmap name is equals to "policies"
//
// returns:
// * true when config is policies
// * false when config is not policies.
func (v *CMFileValidator) CheckCMName(name string) bool {
	if name == "policies" {
		return true
	}
	log.Trace().Str("name", name).Msg("Not a policies cm, skipping")
	return false
}

// ValidateFile checks the validity of a single Policy as yaml file
//
// returns:
// * true, "", nil when config is valid
// * false, message, nil when config is invalid
// and
// * false, "", err on other errors.
//
// ValidateConfig checks the syntax, validity of extractors, and validity of
// rego modules (by attempting to compile them).
func (v *CMFileValidator) ValidateFile(
	ctx context.Context,
	name string,
	yamlSrc []byte,
) (bool, string, error) {
	log.Info().Str("name", name).Msg("Validating CM policy yaml")
	if len(yamlSrc) == 0 {
		return false, "empty yaml", nil
	}
	var policy policiesv1.Policy
	err := config.Unmarshal(yamlSrc, &policy)
	if err != nil {
		return false, err.Error(), nil
	}
	_, err = CompilePolicy(&policy)
	if err != nil {
		return false, err.Error(), nil
	}

	if policy.GetResources() != nil {
		for _, c := range policy.GetResources().Classifiers {
			_, err = classifier.CompileRuleset(ctx, name, c)
			if err != nil {
				if errors.Is(err, classifier.BadExtractor) || errors.Is(err, classifier.BadSelector) ||
					errors.Is(err, classifier.BadRego) || errors.Is(err, classifier.BadLabelName) {
					return false, err.Error(), nil
				} else {
					return false, "", err
				}
			}
		}
	}
	return true, "", nil
}
