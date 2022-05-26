package validation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"aperture.tech/aperture/pkg/log"
)

// limit the concurrent validation requests (some validations incurs rego
// compilation, which might be heavy).
const maxValidationConcurrency = 3

// CMFileValidator is an interface for configmap validation.
type CMFileValidator interface {
	CheckCMName(name string) bool
	ValidateFile(
		ctx context.Context,
		name string,
		yamlSrc []byte,
	) (bool, string, error)
}

// CMValidator validates the flowcontrol configmap.
type CMValidator struct {
	tokens         chan concurrencyToken
	fileValidators []CMFileValidator
}

// NewCMValidator creates a NewCMValidator.
func NewCMValidator() *CMValidator {
	tokens := make(chan concurrencyToken, maxValidationConcurrency)
	for i := 0; i < maxValidationConcurrency; i++ {
		tokens <- concurrencyToken{}
	}
	return &CMValidator{
		tokens:         tokens,
		fileValidators: make([]CMFileValidator, 0),
	}
}

type concurrencyToken struct{}

// RegisterCMFileValidator adds a configmap file validator to be handled on validator
//
// This function should be only called before Start phase.
func (v *CMValidator) RegisterCMFileValidator(validator CMFileValidator) {
	v.fileValidators = append(v.fileValidators, validator)
}

// ValidateObject checks the validity of a object as a k8s object.
func (v *CMValidator) ValidateObject(
	ctx context.Context,
	req *admissionv1.AdmissionRequest,
) (ok bool, msg string, err error) {
	log.Trace().Msg("ValidateObject start")

	select {
	case <-ctx.Done():
		return false, "", errors.New("context expired before concurrency token was ready")
	case tok := <-v.tokens:
		defer func() { v.tokens <- tok }()
	}

	cmKind := corev1.SchemeGroupVersion.WithKind("ConfigMap")
	expectedKind := metav1.GroupVersionKind{
		Group:   cmKind.Group,
		Version: cmKind.Version,
		Kind:    cmKind.Kind,
	}
	if req.Kind != expectedKind {
		return false, "object is not a configmap", nil
	}

	var cm corev1.ConfigMap
	err = json.Unmarshal(req.Object.Raw, &cm)
	if err != nil {
		return
	}

	return v.ValidateConfigMap(ctx, cm)
}

// ValidateConfigMap checks if configmap is valid
//
// returns:
// * true, "", nil when config is valid
// * false, message, nil when config is invalid
// and
// * false, "", err on other errors.
//
func (v *CMValidator) ValidateConfigMap(ctx context.Context, cm corev1.ConfigMap) (bool, string, error) {
	files := make(map[string][]byte, len(cm.Data)+len(cm.BinaryData))

	for filename, contents := range cm.Data {
		files[filename] = []byte(contents)
	}

	for filename, contents := range cm.BinaryData {
		if _, exists := files[filename]; exists {
			return false, fmt.Sprintf("duplicate file %q", filename), nil
		}
		files[filename] = contents
	}

	for _, validator := range v.fileValidators {
		if !validator.CheckCMName(cm.Name) {
			continue
		}
		for filename, contents := range files {
			if !strings.HasSuffix(filename, ".yaml") {
				continue
			}

			ok, msg, err := validator.ValidateFile(ctx, strings.TrimSuffix(filename, ".yaml"), contents)
			if err != nil {
				return false, "", err
			}

			if !ok {
				return ok, fmt.Sprintf("%s: %s", filename, msg), err
			}
		}
	}

	return true, "", nil
}
