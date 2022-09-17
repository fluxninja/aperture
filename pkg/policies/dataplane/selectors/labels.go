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
// see https://docs.fluxninja.com/docs/development/concepts/flow-control/flow-label#sources
type Labels map[string]string

// LabelSources describes all the sources of labels.
type LabelSources struct {
	// Baggage, classifier and explicit labels
	Flow map[string]string
	// Request labels
	Request *ext_authz.AttributeContext_Request
}

const (
	requestLabelPrefix       = "http."
	requestLabelHeaderPrefix = "http.request.header."
	// Number of always-added request labels (this value is used only for
	// capacity estimation).
	numRequestLabels = 5
)

// ToPlainMap extracts the underlying map of labels.
func (l Labels) ToPlainMap() map[string]string { return map[string]string(l) }

// NewLabels maps all available labels into a flat namespace, as documented in
// https://docs.fluxninja.com/docs/development/concepts/flow-control/flow-label#sources
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
			labels[requestLabelPrefix+"method"] = http.Method
			labels[requestLabelPrefix+"target"] = http.Path
			labels[requestLabelPrefix+"host"] = http.Host
			labels[requestLabelPrefix+"scheme"] = http.Scheme
			labels[requestLabelPrefix+"request_content_length"] = strconv.FormatInt(http.Size, 10)
			labels[requestLabelPrefix+"flavor"] = CanonicalizeOtelHTTPFlavor(http.Protocol)
		}
		for k, v := range newLabels.Request.GetHttp().GetHeaders() {
			if strings.HasPrefix(k, ":") {
				// Headers starting with `:` are pseudoheaders, so we don't add
				// them.  We don't lose anything, as these values are already
				// available as labels pulled from dedicated fields of
				// Request.Http.
				continue
			}
			labels[requestLabelHeaderPrefix+canonicalizeOtelHeaderKey(k)] = v
		}
	}

	return Labels(labels)
}

// canonicalizeOtelHeaderKey converts envoy's header naming convention to Otel's one.
func canonicalizeOtelHeaderKey(key string) string {
	return strings.ReplaceAll(key, "-", "_")
}

// CanonicalizeOtelHTTPFlavor converts envoy's protocol to Otel kind of HTTP protocol.
func CanonicalizeOtelHTTPFlavor(protocolName string) string {
	switch protocolName {
	case "HTTP/1.0":
		return "1.0"
	case "HTTP/1.1":
		return "1.1"
	case "HTTP/2":
		return "2.0"
	case "HTTP/3":
		return "3.0"
	default:
		return protocolName
	}
}
