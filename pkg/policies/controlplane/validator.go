package controlplane

import (
	"context"
	"errors"

	policiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier/compiler"
)

// ValidateAndCompile checks the validity of a single Policy and compiles it.
func ValidateAndCompile(ctx context.Context, name string, yamlSrc []byte) (CompiledCircuit, bool, string, error) {
	log.Info().Str("name", name).Msg("Validating CM policy yaml")
	if len(yamlSrc) == 0 {
		return nil, false, "empty yaml", nil
	}
	policy := &policiesv1.Policy{}
	err := config.UnmarshalYAML(yamlSrc, policy)
	if err != nil {
		return nil, false, err.Error(), nil
	}
	circuit, err := CompilePolicy(policy)
	if err != nil {
		return nil, false, err.Error(), nil
	}

	if policy.GetResources() != nil {
		for _, c := range policy.GetResources().Classifiers {
			_, err = compiler.CompileRuleset(ctx, name, &wrappersv1.ClassifierWrapper{
				Classifier: c,
			})
			if err != nil {
				if errors.Is(err, compiler.BadExtractor) || errors.Is(err, compiler.BadSelector) ||
					errors.Is(err, compiler.BadRego) || errors.Is(err, compiler.BadLabelName) {
					return nil, false, err.Error(), nil
				} else {
					return nil, false, "", err
				}
			}
		}
	}
	return circuit, true, "", nil
}
