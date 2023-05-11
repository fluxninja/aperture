package iface

import policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"

// PreviewID is the ID of a preview.
type PreviewID struct {
	RequestID string
}

// String returns the string representation of the ID.
func (id PreviewID) String() string {
	return id.RequestID
}

// PreviewBase is the base interface for all preview requests.
type PreviewBase interface {
	// GetPreviewID returns the ID of the preview.
	GetPreviewID() PreviewID
	// GetSelectors returns the selectors.
	GetSelectors() []*policylangv1.Selector
}

// LabelPreview interface.
type LabelPreview interface {
	PreviewBase
	// AddLabelPreview adds labels to preview.
	AddLabelPreview(labels map[string]string)
}

// HTTPRequestPreview interface.
type HTTPRequestPreview interface {
	PreviewBase
	// AddHTTPRequestPreview adds labels to preview.
	AddHTTPRequestPreview(request map[string]interface{})
}
