package envoy

import (
	"strconv"
	"strings"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
)

const (
	requestLabelPrefix       = "http."
	requestLabelHeaderPrefix = "http.request.header."
	// Number of always-added request labels (this value is used only for
	// capacity estimation).
	numRequestLabels = 6
)

// AuthzRequestToFlowLabels converts request attributes to new FlowLabels.
func AuthzRequestToFlowLabels(request *authv3.AttributeContext_Request) flowlabel.FlowLabels {
	capacity := numRequestLabels + len(request.GetHttp().GetHeaders())
	flowLabels := make(flowlabel.FlowLabels, capacity)
	if request != nil {
		if http := request.GetHttp(); http != nil {
			flowLabels[requestLabelPrefix+"method"] = flowlabel.FlowLabelValue{
				Value:     http.Method,
				Telemetry: true,
			}
			flowLabels[requestLabelPrefix+"target"] = flowlabel.FlowLabelValue{
				Value:     http.Path,
				Telemetry: true,
			}
			flowLabels[requestLabelPrefix+"host"] = flowlabel.FlowLabelValue{
				Value:     http.Host,
				Telemetry: true,
			}
			flowLabels[requestLabelPrefix+"scheme"] = flowlabel.FlowLabelValue{
				Value:     http.Scheme,
				Telemetry: true,
			}
			flowLabels[requestLabelPrefix+"request_content_length"] = flowlabel.FlowLabelValue{
				Value:     strconv.FormatInt(http.Size, 10),
				Telemetry: false,
			}
			flowLabels[requestLabelPrefix+"flavor"] = flowlabel.FlowLabelValue{
				Value:     canonicalizeOtelHTTPFlavor(http.Protocol),
				Telemetry: true,
			}
		}
		for k, v := range request.GetHttp().GetHeaders() {
			if strings.HasPrefix(k, ":") {
				// Headers starting with `:` are pseudoheaders, so we do not add
				// them.  We do not lose anything, as these values are already
				// available as labels pulled from dedicated fields of
				// Request.Http.
				continue
			}
			flowLabels[requestLabelHeaderPrefix+canonicalizeOtelHeaderKey(k)] = flowlabel.FlowLabelValue{
				Value:     v,
				Telemetry: false,
			}
		}
	}

	return flowLabels
}

// canonicalizeOtelHeaderKey converts header naming convention to Otel's one.
func canonicalizeOtelHeaderKey(key string) string {
	return strings.ReplaceAll(strings.ToLower(key), "-", "_")
}

// canonicalizeOtelHTTPFlavor converts protocol to Otel kind of HTTP protocol.
func canonicalizeOtelHTTPFlavor(protocolName string) string {
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
