package flowcontrol

import (
	"context"

	languagev1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/log"
)

// CMFileValidator Flowcontrol implementation of CMFileValidator interface.
type CMFileValidator struct{}

// CheckCMName checks configmap name is equals to "flowcontrol"
//
// returns:
// * true when config is flowcontrol
// * false when config is not flowcontrol.
func (v *CMFileValidator) CheckCMName(name string) bool {
	if name == "flowcontrol" {
		return true
	}
	log.Trace().Str("name", name).Msg("Not a flowcontrol cm, skipping")
	return false
}

// ValidateFile checks the validity of a single Flowcontrol Policy as yaml file
//
// returns:
// * true, "", nil when config is valid
// * false, message, nil when config is invalid
// and
// * false, "", err on other errors.
//
// ValidateConfig checks the syntax, validity of policies.
func (v *CMFileValidator) ValidateFile(
	ctx context.Context,
	name string,
	yamlSrc []byte,
) (bool, string, error) {
	log.Trace().Str("name", name).Msg("Validating CM flowcontrol yaml")
	if len(yamlSrc) == 0 {
		log.Warn().Msg("Validation: Empty flowcontrol configmap file")
		return false, "", nil
	}
	var policy languagev1.Policy
	err := config.Unmarshal(yamlSrc, &policy)
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", nil
}
