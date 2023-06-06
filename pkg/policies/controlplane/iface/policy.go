package iface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const (
	// PoliciesRoot - path in config and status registry for policies results.
	PoliciesRoot = "policies"
)

// FxOptionsFuncTag allows sub-modules to provide their options to per policy apps independently.
var FxOptionsFuncTag = config.GroupTag("policy-fx-funcs")

// PolicyBase is for read only access to base policy info.
type PolicyBase interface {
	GetPolicyName() string
	GetPolicyHash() string
}

// Policy is for read only access to full policy state.
type Policy interface {
	PolicyBase
	GetEvaluationInterval() time.Duration
	GetStatusRegistry() status.Registry
}

// GetSelectorsShortDescription returns a short description of the selectors.
func GetSelectorsShortDescription(selectors []*policylangv1.Selector) string {
	return fmt.Sprintf("%d selectors", len(selectors))
}

// PolicyMessage is used for passing policies to be wrapped.
type PolicyMessage struct {
	*policylangv1.Policy
	PolicyMetadata *policylangv1.PolicyMetadata `json:"policy_metadata"`
}

// MarshalJSON implements the json.Marshaler interface.
func (m *PolicyMessage) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	pb, err := m.Policy.MarshalJSON()
	if err != nil {
		return nil, err
	}
	buf.Write(pb[:len(pb)-1])
	if m.PolicyMetadata != nil {
		buf.WriteString(`,"policy_metadata":`)
		pmb, err := m.PolicyMetadata.MarshalJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(pmb)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *PolicyMessage) UnmarshalJSON(b []byte) error {
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	err := unmarshaler.Unmarshal(b, m.Policy)
	if err != nil {
		return err
	}
	tmp := struct {
		PolicyMetadata *policylangv1.PolicyMetadata `json:"policy_metadata"`
	}{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	m.PolicyMetadata = tmp.PolicyMetadata
	return nil
}
