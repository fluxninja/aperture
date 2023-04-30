package iface

import (
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate mockgen -source=classifier.go -destination=../../mocks/mock_classifier.go -package=mocks

// ClassifierID is the ID of the Classifier.
type ClassifierID struct {
	PolicyName      string
	PolicyHash      string
	ClassifierIndex int64
}

// String function returns the ClassifierID as a string.
func (cID ClassifierID) String() string {
	return fmt.Sprintf("policy_name-%s-policy_hash-%s-%d", cID.PolicyName, cID.PolicyHash, cID.ClassifierIndex)
}

// Classifier interface.
type Classifier interface {
	// GetSelectors returns the selectors.
	GetSelectors() []*policylangv1.Selector

	// GetClassifierID returns ClassifierID object that should uniquely identify classifier.
	GetClassifierID() ClassifierID

	// GetRequestCounter returns the counter for the classifier.
	GetRequestCounter() prometheus.Counter
}
