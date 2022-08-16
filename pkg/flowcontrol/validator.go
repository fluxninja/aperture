package flowcontrol

import (
	"context"

	languagev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

// CMFileValidator policies implementation of CMFileValidator interface.
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

// ValidateFile checks the validity of a single policies Policy as yaml file
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
	log.Trace().Str("name", name).Msg("Validating CM policies yaml")
	if len(yamlSrc) == 0 {
		log.Warn().Msg("Validation: Empty policies configmap file")
		return false, "", nil
	}
	var policy languagev1.Policy
	err := config.Unmarshal(yamlSrc, &policy)
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", nil
}
