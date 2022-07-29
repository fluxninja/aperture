package classification

import (
	"context"
	"errors"

	classificationv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/classification/v1"
	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/log"
)

// CMFileValidator Classification implementation of CMFileValidator interface.
type CMFileValidator struct{}

// CheckCMName checks configmap name is equals to "classification"
//
// returns:
// * true when config is classification
// * false when config is not classification.
func (v *CMFileValidator) CheckCMName(name string) bool {
	if name == "classification" {
		return true
	}
	log.Trace().Str("name", name).Msg("Not a classification cm, skipping")
	return false
}

// ValidateFile checks the validity of a single Classification Ruleset as yaml file
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
	log.Trace().Str("name", name).Msg("Validating CM classification yaml")
	var classifier classificationv1.Classifier
	err := config.Unmarshal(yamlSrc, &classifier)
	if err != nil {
		return false, err.Error(), nil
	}
	_, err = compileRuleset(ctx, name, &classifier)
	if err != nil {
		if errors.Is(err, BadExtractor) || errors.Is(err, BadSelector) ||
			errors.Is(err, BadRego) || errors.Is(err, BadLabelName) {
			return false, err.Error(), nil
		} else {
			return false, "", err
		}
	}
	return true, "", nil
}
