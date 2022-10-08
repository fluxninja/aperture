package iface

import (
	"fmt"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
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
	// GetSelector returns the selector.
	GetSelector() *selectorv1.Selector

	// GetClassifierID returns ClassifierID object that should uniquely identify classifier.
	GetClassifierID() ClassifierID

	// GetCounter returns the counter for the classifier.
	GetCounter() prometheus.Counter
}
