package checkhttp

import (
	"strconv"
	"strings"

	flowcontrolhttpv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	flowlabel "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/label"
)

const (
	requestLabelPrefix       = "http."
	requestLabelHeaderPrefix = "http.request.header."
	// Number of always-added request labels (this value is used only for
	// capacity estimation).
	numRequestLabels = 6
)

// CheckHTTPRequestToFlowLabels converts request attributes to new FlowLabels.
// It takes a flowcontrolhttpv1.CheckHTTPRequest_HttpRequest object as input and returns a flowlabel.FlowLabels object.
// The function adds several labels to the flowlabel.FlowLabels object based on the attributes of the flowcontrolhttpv1.CheckHTTPRequest_HttpRequest object.
// If the request is nil, the function returns an empty flowlabel.FlowLabels object.
func CheckHTTPRequestToFlowLabels(request *flowcontrolhttpv1.CheckHTTPRequest_HttpRequest) flowlabel.FlowLabels {
	if request == nil {
		return flowlabel.FlowLabels{}
	}

	capacity := numRequestLabels + len(request.GetHeaders())
	flowLabels := make(flowlabel.FlowLabels, capacity)
	flowLabels[requestLabelPrefix+"method"] = flowlabel.FlowLabelValue{
		Value:     request.GetMethod(),
		Telemetry: true,
	}
	flowLabels[requestLabelPrefix+"target"] = flowlabel.FlowLabelValue{
		Value:     request.GetPath(),
		Telemetry: true,
	}
	flowLabels[requestLabelPrefix+"host"] = flowlabel.FlowLabelValue{
		Value:     request.GetHost(),
		Telemetry: true,
	}
	flowLabels[requestLabelPrefix+"scheme"] = flowlabel.FlowLabelValue{
		Value:     request.GetScheme(),
		Telemetry: true,
	}
	flowLabels[requestLabelPrefix+"request_content_length"] = flowlabel.FlowLabelValue{
		Value:     strconv.FormatInt(request.GetSize(), 10),
		Telemetry: false,
	}
	flowLabels[requestLabelPrefix+"flavor"] = flowlabel.FlowLabelValue{
		Value:     canonicalizeOtelHTTPFlavor(request.GetProtocol()),
		Telemetry: true,
	}
	for k, v := range request.GetHeaders() {
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
