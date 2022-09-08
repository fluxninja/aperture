package selectors

import (
	"fmt"
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
	requestLabelPrefix       = "http"
	requestLabelHeaderPrefix = "http.request.header"
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
			var contentLength int64
			if http.Size < 0 {
				contentLength = http.Size
			} else {
				// TODO: compute if unknown
			}
			labels[fmt.Sprintf("%s.method", requestLabelPrefix)] = http.Method
			labels[fmt.Sprintf("%s.target", requestLabelPrefix)] = http.Path
			labels[fmt.Sprintf("%s.scheme", requestLabelPrefix)] = http.Scheme
			labels[fmt.Sprintf("%s.request_content_length", requestLabelPrefix)] = strconv.FormatInt(contentLength, 10)
			labels[fmt.Sprintf("%s.flavor", requestLabelPrefix)] = convertProtocol(http.Protocol)
		}
		for k, v := range newLabels.Request.GetHttp().GetHeaders() {
			if strings.HasPrefix(k, ":") {
				// Headers starting with `:` are pseudoheaders, so we don't add
				// them.  We don't lose anything, as these values are already
				// available as labels pulled from dedicated fields of
				// Request.Http.
				continue
			}
			labels[fmt.Sprintf("%s.%s", requestLabelHeaderPrefix, k)] = v
		}
	}

	return Labels(labels)
}

func convertHeaderKey(key string) string {
	return strings.ReplaceAll(key, "-", ".")
}

// convertProtocol converts envoy's protocol to OTEL kind of HTTP protocol.
func convertProtocol(protocolName string) string {
	protocols := map[string]string{
		"HTTP/1.0": "1.0",
		"HTTP/1.1": "1.1",
		"HTTP/2.0": "2.0",
		"HTTP/3.0": "3.0",
	}
	flavor, ok := protocols[protocolName]
	if !ok {
		return protocolName
	}
	return flavor
}
