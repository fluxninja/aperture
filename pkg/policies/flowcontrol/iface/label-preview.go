package iface

import policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"

// LabelPreviewID is the ID of a preview.
type LabelPreviewID struct {
	RequestID string
}

// String returns the string representation of the ID.
func (id LabelPreviewID) String() string {
	return id.RequestID
}

// LabelPreview interface.
type LabelPreview interface {
	// GetLabelPreviewID returns the ID of the preview.
	GetLabelPreviewID() LabelPreviewID
	// GetFlowSelector returns the flow selector.
	GetFlowSelector() *policylangv1.FlowSelector
	// AddLabelPreview adds labels to preview.
	AddLabelPreview(labels map[string]string)
}
