package checkhttp

import (
	"context"
	"encoding/base64"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	flowcontrolhttpv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpc"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	flowlabel "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/label"
	classification "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	servicegetter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service-getter"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/check"
	checkhttp_baggage "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/checkhttp/baggage"
)

var baggageSanitizeRegex *regexp.Regexp = regexp.MustCompile(`[\s\\\/;",]`)

var (
	missingControlPointSampler    = log.NewRatelimitingSampler()
	failedBaggageInjectionSampler = log.NewRatelimitingSampler()
)

// NewHandler creates new handler for flowcontrol CheckHTTP
//
// It will use the given classifier to inject flow labels and return them as
// metadata in the response to the Check calls.
func NewHandler(
	classifier *classification.ClassificationEngine,
	serviceGetter servicegetter.ServiceGetter,
	fcHandler check.HandlerWithValues,
) *Handler {
	return &Handler{
		classifier:    classifier,
		serviceGetter: serviceGetter,
		propagator:    checkhttp_baggage.W3Baggage{},
		fcHandler:     fcHandler,
	}
}

// sanitizeBaggageHeaderValue excludes characters that should be url escaped
// Otherwise both baggage.String method and envoy itself will do it.
func sanitizeBaggageHeaderValue(value string) string {
	// All characters allowed except control chars, whitespace, double quote, comma, semicolon, backslash
	// see https://www.w3.org/TR/baggage/#header-content
	cleanValue := baggageSanitizeRegex.ReplaceAll([]byte(value), []byte("-"))
	return string(cleanValue)
}

// Handler implements aperture.flowcontrol.v1.FlowControlServiceHTTP and handles Check call.
type Handler struct {
	flowcontrolhttpv1.UnimplementedFlowControlServiceHTTPServer
	serviceGetter servicegetter.ServiceGetter
	classifier    *classification.ClassificationEngine
	propagator    checkhttp_baggage.Propagator
	fcHandler     check.HandlerWithValues
}

// CheckHTTP is the Check method of Flow Control service returns the allow/deny decisions of
// whether to accept the traffic after running the algorithms.
func (h *Handler) CheckHTTP(ctx context.Context, req *flowcontrolhttpv1.CheckHTTPRequest) (*flowcontrolhttpv1.CheckHTTPResponse, error) {
	// Put inner fields back into pool.
	// Note: Not pooling the whole CheckHTTPRequest, as we don't control the object creation.
	defer req.ResetVT()
	// record the start time of the request
	start := time.Now()
	createResponse := func(checkResponse *flowcontrolv1.CheckResponse) *flowcontrolhttpv1.CheckHTTPResponse {
		marshalledCheckResponse, err := checkResponse.MarshalVT()
		if err != nil {
			log.Bug().Err(err).Msg("bug: Failed to marshal check response")
			return nil
		}
		checkResponseBase64 := base64.StdEncoding.EncodeToString(marshalledCheckResponse)

		// record the end time of the request
		end := time.Now()
		checkResponse.Start = timestamppb.New(start)
		checkResponse.End = timestamppb.New(end)

		return &flowcontrolhttpv1.CheckHTTPResponse{
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					otelconsts.ApertureCheckResponseLabel: structpb.NewStringValue(checkResponseBase64),
				},
			},
		}
	}

	ctrlPt := req.GetControlPoint()
	if ctrlPt == "" {
		return nil, grpc.LoggedError(log.Sample(missingControlPointSampler).Warn()).
			Code(codes.InvalidArgument).Msg("missing control-point")
	}

	sourceAddress := req.GetSource().GetAddress()
	sourceSvcs := h.serviceGetter.ServicesFromAddress(sourceAddress)
	if sourceSvcs == nil {
		sourceSvcs = []string{"UNKNOWN"}
	}
	sourceSvcsStr := strings.Join(sourceSvcs, ",")
	destinationAddress := req.GetDestination().GetAddress()
	destinationSvcs := h.serviceGetter.ServicesFromAddress(destinationAddress)
	if destinationSvcs == nil {
		sourceSvcs = []string{"UNKNOWN"}
	}
	destinationSvcsStr := strings.Join(destinationSvcs, ",")

	// make flowlabels from source and destination services
	sdFlowLabels := make(flowlabel.FlowLabels, 2)
	sdFlowLabels[otelconsts.ApertureSourceServiceLabel] = flowlabel.FlowLabelValue{
		Value:     sourceSvcsStr,
		Telemetry: true,
	}
	sdFlowLabels[otelconsts.ApertureDestinationServiceLabel] = flowlabel.FlowLabelValue{
		Value:     destinationSvcsStr,
		Telemetry: true,
	}

	input := RequestToInputWithServices(req, sourceSvcs, destinationSvcs)

	// Default flow labels from request
	requestFlowLabels := CheckHTTPRequestToFlowLabels(req.GetRequest())
	existingHeaders := checkhttp_baggage.Headers(req.GetRequest().GetHeaders())
	baggageFlowLabels := h.propagator.Extract(existingHeaders)

	// Merge flow labels from request and baggage headers
	mergedFlowLabels := requestFlowLabels
	// Baggage can overwrite request flow labels
	flowlabel.Merge(mergedFlowLabels, baggageFlowLabels)
	flowlabel.Merge(mergedFlowLabels, sdFlowLabels)

	classifierMsgs, newFlowLabels := h.classifier.Classify(ctx, destinationSvcs, ctrlPt, mergedFlowLabels, input)

	for key, fl := range newFlowLabels {
		cleanValue := sanitizeBaggageHeaderValue(fl.Value)
		fl.Value = cleanValue
		newFlowLabels[key] = fl
	}

	// Add new flow labels to baggage
	newHeaders, err := h.propagator.Inject(newFlowLabels, existingHeaders)
	if err != nil {
		log.Sample(failedBaggageInjectionSampler).
			Warn().Err(err).Msg("Failed to inject baggage into headers")
	}

	// Make the freshly created flow labels available to flowcontrol.
	// Newly created flow labels can overwrite existing flow labels.
	flowlabel.Merge(mergedFlowLabels, newFlowLabels)

	// Ask flow control service for Ok/Deny
	// checkResponse := h.fcHandler.CheckRequest(ctx, destinationSvcs, ctrlPt, flowLabels)
	checkResponse := h.fcHandler.CheckRequest(ctx,
		iface.RequestContext{
			Services:     destinationSvcs,
			ControlPoint: ctrlPt,
			FlowLabels:   mergedFlowLabels,
		},
	)
	checkResponse.ClassifierInfos = classifierMsgs
	// Set telemetry_flow_labels in the CheckResponse
	checkResponse.TelemetryFlowLabels = mergedFlowLabels.TelemetryLabels()
	// add control point type
	checkResponse.TelemetryFlowLabels[otelconsts.ApertureControlPointTypeLabel] = otelconsts.HTTPControlPoint
	checkResponse.TelemetryFlowLabels[otelconsts.ApertureSourceServiceLabel] = strings.Join(sourceSvcs, ",")
	checkResponse.TelemetryFlowLabels[otelconsts.ApertureDestinationServiceLabel] = strings.Join(destinationSvcs, ",")

	resp := createResponse(checkResponse)

	switch checkResponse.DecisionType {
	case flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED:
		resp.Status = &status.Status{
			Code: int32(code.Code_OK),
		}
		resp.HttpResponse = &flowcontrolhttpv1.CheckHTTPResponse_OkResponse{
			OkResponse: &flowcontrolhttpv1.OkHttpResponse{
				Headers: newHeaders,
			},
		}
	case flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED:
		resp.Status = &status.Status{
			Code: int32(code.Code_UNAVAILABLE),
		}

		var statusCode int32
		switch checkResponse.RejectReason {
		case flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED:
			statusCode = http.StatusTooManyRequests
		case flowcontrolv1.CheckResponse_REJECT_REASON_NO_TOKENS:
			statusCode = http.StatusServiceUnavailable
		case flowcontrolv1.CheckResponse_REJECT_REASON_REGULATED:
			statusCode = http.StatusForbidden
		default:
			log.Bug().Stringer("reason", checkResponse.RejectReason).Msg("Unexpected reject reason")
		}

		deniedHTTPResponse := &flowcontrolhttpv1.DeniedHttpResponse{
			Status: statusCode,
		}
		if checkResponse.WaitTime != nil {
			deniedHTTPResponse.Headers["retry-after"] = waitTimeToRetryAfter(checkResponse.WaitTime)
			// Clear to avoid redundancy, as we're translating it into header.
			// Logs processor doesn't read it and clients aren't supposed to
			// peek into CheckResponse.
			checkResponse.WaitTime = nil
		}
		resp.HttpResponse = &flowcontrolhttpv1.CheckHTTPResponse_DeniedResponse{
			DeniedResponse: deniedHTTPResponse,
		}
	default:
		return nil, grpc.Bug().Stringer("type", checkResponse.DecisionType).
			Msg("unexpected decision type")
	}

	return resp, nil
}

func waitTimeToRetryAfter(waitTime *durationpb.Duration) string {
	seconds := waitTime.Seconds
	// Retry-after header doesn't have full second resolution, we need to round
	// up to avoid clients retrying too early.
	if waitTime.Nanos != 0 {
		seconds += 1
	}
	return strconv.FormatInt(seconds, 10)
}
