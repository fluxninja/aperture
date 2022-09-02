package iface

import selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"

// ClassifierID is the ID of the Classifier.
type ClassifierID struct {
	PolicyName string
	PolicyHash string
}

// String function returns the ClassifierID as a string.
func (cID ClassifierID) String() string {
	return "policy_name-" + cID.PolicyName + "-policy_hash-" + cID.PolicyHash
}

// Classifier interface.
type Classifier interface {
	// GetSelector returns the selector.
	GetSelector() *selectorv1.Selector
	// GetClassifierID returns ClassifierID object that should uniquely identify classifier.
	GetClassifierID() ClassifierID
}
