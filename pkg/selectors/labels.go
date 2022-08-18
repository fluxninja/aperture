package selectors

import (
	"strconv"
	"strings"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

// Labels is a flattened map of labels expected to be in canonical form to be
// used in LabelMatchers in Selectors.
//
// Canonical form implies appropriate prefixes for different types of labels,
// see docs on policy/language/v1.Selector.
type Labels map[string]string

// LabelSources describes all the sources of labels (flow and request).
type LabelSources struct {
	Flow    map[string]string
	Request *ext_authz.AttributeContext_Request
}

const (
	requestLabelPrefix       = "request_"
	requestLabelHeaderPrefix = "request_header_"
	numRequestLabels         = 7 // used only for capacity estimation
)

// ToPlainMap extracts the underlying map of labels.
func (l Labels) ToPlainMap() map[string]string { return map[string]string(l) }

// NewLabels maps all available labels into a flat namespace, as documented in
// policy/language/v1.Selector.
func NewLabels(ls LabelSources) Labels {
	return Labels(nil).CombinedWith(ls)
}

// CombinedWith combines existing labels with new labels from LabelSources.
//
// The original labels map is unchanged.
// In case of duplicates, values from `newLabels` will be used.
func (l Labels) CombinedWith(newLabels LabelSources) Labels {
	capacity := len(l) + len(newLabels.Flow)
	if newLabels.Request != nil {
		capacity += numRequestLabels + len(newLabels.Request.GetHttp().GetHeaders())
	}

	labels := make(map[string]string, capacity)

	for k, v := range l {
		labels[k] = v
	}

	for k, v := range newLabels.Flow {
		labels[k] = v
	}

	if newLabels.Request != nil {
		if http := newLabels.Request.GetHttp(); http != nil {
			labels[requestLabelPrefix+"id"] = http.Id
			labels[requestLabelPrefix+"method"] = http.Method
			labels[requestLabelPrefix+"path"] = http.Path
			labels[requestLabelPrefix+"host"] = http.Host
			labels[requestLabelPrefix+"scheme"] = http.Scheme
			labels[requestLabelPrefix+"size"] = strconv.FormatInt(http.Size, 10)
			labels[requestLabelPrefix+"protocol"] = http.Protocol
		}
		for k, v := range newLabels.Request.GetHttp().GetHeaders() {
			if strings.HasPrefix(k, ":") {
				// Headers starting with `:` are pseudoheaders, so we don't add
				// them.  We don't lose anything, as these values are already
				// available as labels pulled from dedicated fields of
				// Request.Http.
				continue
			}
			labels[requestLabelHeaderPrefix+k] = v
		}
	}

	return Labels(labels)
}
