package policyvalidator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
)

// limit the concurrent validation requests (some validations incurs rego
// compilation, which might be heavy).
const maxValidationConcurrency = 3

// PolicySpecValidator is an interface for Policy Custom Resource validation.
type PolicySpecValidator interface {
	ValidateSpec(
		ctx context.Context,
		name string,
		yamlSrc []byte,
	) (bool, string, error)
}

// PolicyValidator validates the Policy Custom Resource.
type PolicyValidator struct {
	tokens           chan concurrencyToken
	policyValidators []PolicySpecValidator
}

// NewPolicyValidator creates a new instance of PolicyValidator.
func NewPolicyValidator(validators []PolicySpecValidator) *PolicyValidator {
	tokens := make(chan concurrencyToken, maxValidationConcurrency)
	for i := 0; i < maxValidationConcurrency; i++ {
		tokens <- concurrencyToken{}
	}
	return &PolicyValidator{
		tokens:           tokens,
		policyValidators: validators,
	}
}

type concurrencyToken struct{}

// ValidateObject checks the validity of a object as a k8s object.
func (v *PolicyValidator) ValidateObject(
	ctx context.Context,
	req *admissionv1.AdmissionRequest,
) (ok bool, msg string, err error) {
	select {
	case <-ctx.Done():
		return false, "", errors.New("context expired before concurrency token was ready")
	case tok := <-v.tokens:
		defer func() { v.tokens <- tok }()
	}

	policyKind := api.GroupVersion.WithKind("Policy")
	expectedKind := metav1.GroupVersionKind{
		Group:   policyKind.Group,
		Version: policyKind.Version,
		Kind:    policyKind.Kind,
	}
	if req.Kind != expectedKind {
		return false, "object is not a Policy", nil
	}

	if req.Namespace != os.Getenv("APERTURE_CONTROLLER_NAMESPACE") {
		return false, "Policy should be created in the same namespace as Aperture Controller", nil
	}

	var policy policyv1alpha1.Policy
	err = json.Unmarshal(req.Object.Raw, &policy)
	if err != nil {
		return false, "Not a valid Policy Object", err
	}

	for _, validator := range v.policyValidators {
		ok, msg, err := validator.ValidateSpec(ctx, policy.GetName(), policy.Spec.Raw)
		if err != nil {
			return false, "Validator error", err
		}

		if !ok {
			return false, fmt.Sprintf("%s: %s", policy.GetName(), msg), nil
		}
	}

	return true, fmt.Sprintf("%s: Valid Policy", policy.GetName()), nil
}
