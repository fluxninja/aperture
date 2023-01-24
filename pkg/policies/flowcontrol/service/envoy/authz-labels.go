package envoy

import (
	"strconv"
	"strings"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	"github.com/fluxninja/aperture/pkg/utils"
)

const (
	requestLabelPrefix       = "http."
	requestLabelHeaderPrefix = "http.request.header."
	// Number of always-added request labels (this value is used only for
	// capacity estimation).
	numRequestLabels = 6
)

// AuthzRequestToFlowLabels converts request attributes to new FlowLabels.
func AuthzRequestToFlowLabels(request *ext_authz.AttributeContext_Request) flowlabel.FlowLabels {
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
				Value:     utils.CanonicalizeOtelHTTPFlavor(http.Protocol),
				Telemetry: true,
			}
		}
		for k, v := range request.GetHttp().GetHeaders() {
			if strings.HasPrefix(k, ":") {
				// Headers starting with `:` are pseudoheaders, so we don't add
				// them.  We don't lose anything, as these values are already
				// available as labels pulled from dedicated fields of
				// Request.Http.
				continue
			}
			flowLabels[requestLabelHeaderPrefix+utils.CanonicalizeOtelHeaderKey(k)] = flowlabel.FlowLabelValue{
				Value:     v,
				Telemetry: false,
			}
		}
	}

	return flowLabels
}
