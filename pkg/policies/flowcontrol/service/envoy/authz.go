package envoy

import (
	"bytes"
	"context"
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/checkhttp"
	authz_baggage "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/envoy/baggage"
)

// NewHandler creates new authorization handler for authz api
//
// Authz will use the given classifier to inject flow labels and return them as
// metadata in the response to the Check calls.
func NewHandler(
	classifier *classification.ClassificationEngine,
	serviceGetter servicegetter.ServiceGetter,
	fcHandler check.HandlerWithValues,
) *Handler {
	return &Handler{
		classifier:    classifier,
		serviceGetter: serviceGetter,
		propagator:    authz_baggage.W3Baggage{},
		fcHandler:     fcHandler,
	}
}

// Handler implements envoy.service.auth.v3.Authorization and handles Check call.
type Handler struct {
	serviceGetter servicegetter.ServiceGetter
	classifier    *classification.ClassificationEngine
	propagator    authz_baggage.Propagator
	fcHandler     check.HandlerWithValues
}

var baggageSanitizeRegex *regexp.Regexp = regexp.MustCompile(`[\s\\\/;",]`)

var (
	missingControlPointSampler    = log.NewRatelimitingSampler()
	failedBaggageInjectionSampler = log.NewRatelimitingSampler()
)

// sanitizeBaggageHeaderValue excludes characters that should be url escaped
// Otherwise both baggage.String method and envoy itself will do it.
func sanitizeBaggageHeaderValue(value string) string {
	// All characters allowed except control chars, whitespace, double quote, comma, semicolon, backslash
	// see https://www.w3.org/TR/baggage/#header-content
	cleanValue := baggageSanitizeRegex.ReplaceAll([]byte(value), []byte("-"))
	return string(cleanValue)
}

// Check is the Check method of Authorization service
//
// Check
// * computes flow labels and returns them via DynamicMetadata.
// * makes the allow/deny decision - sends flow labels to flow control's Check function.
func (h *Handler) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// record the start time of the request
	start := time.Now()

	createExtAuthzResponse := func(checkResponse *flowcontrolv1.CheckResponse) *authv3.CheckResponse {
		// We do not care about the particular format we send the CheckResponse,
		// Envoy can treat is as black-box. The only thing we care about is for
		// it to be deserializable by logs processing pipeline.
		// Using protobuf wire format as it is faster to serialize/deserialize
		// than using protojson or roundtripping through structpb.Struct.
		// Additional base64 encoding step is used, as there's no way to push
		// binary data through dynamic metadata and envoy's access log
		// formatter. Overhead of this base64 encoding is small though.
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

		return &authv3.CheckResponse{
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					otelconsts.ApertureCheckResponseLabel: structpb.NewStringValue(checkResponseBase64),
				},
			},
		}
	}

	ctrlPt := ""
	headers, _ := metadata.FromIncomingContext(ctx)
	if ctrlPtHeader, exists := headers["control-point"]; exists && len(ctrlPtHeader) > 0 {
		ctrlPt = ctrlPtHeader[0]
	} else if ctrlPtHeader, exists := req.GetAttributes().GetContextExtensions()["control-point"]; exists && len(ctrlPtHeader) > 0 {
		ctrlPt = ctrlPtHeader
	}

	if ctrlPt == "" {
		// TODO(krdln) metrics
		return nil, grpc.LoggedError(log.Sample(missingControlPointSampler).Warn()).
			Code(codes.InvalidArgument).Msg("missing control-point")
	}

	sourceAddress := req.GetAttributes().GetSource().GetAddress().GetSocketAddress()
	sourceSvcs := h.serviceGetter.ServicesFromSocketAddress(sourceAddress)
	sourceSvcsStr := strings.Join(sourceSvcs, ",")
	destinationAddress := req.GetAttributes().GetDestination().GetAddress().GetSocketAddress()
	destinationSvcs := h.serviceGetter.ServicesFromSocketAddress(destinationAddress)
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

	if req.GetAttributes().GetRequest().GetHttp().GetBody() == "" && len(req.GetAttributes().GetRequest().GetHttp().GetRawBody()) != 0 {
		req.GetAttributes().GetRequest().GetHttp().Body = bytes.NewBuffer(req.GetAttributes().GetRequest().GetHttp().GetRawBody()).String()
	}
	checkHTTPReq := authzRequestToCheckHTTPRequest(req, ctrlPt)
	input := checkhttp.RequestToInputWithServices(checkHTTPReq, sourceSvcs, destinationSvcs)

	// Default flow labels from Authz request
	requestFlowLabels := AuthzRequestToFlowLabels(req.GetAttributes().GetRequest())
	// Extract flow labels from baggage headers
	existingHeaders := authz_baggage.Headers(req.GetAttributes().GetRequest().GetHttp().GetHeaders())
	baggageFlowLabels := h.propagator.Extract(existingHeaders)

	// Merge flow labels from Authz request and baggage headers
	mergedFlowLabels := requestFlowLabels
	// Baggage can overwrite request flow labels
	flowlabel.Merge(mergedFlowLabels, baggageFlowLabels)
	flowlabel.Merge(mergedFlowLabels, sdFlowLabels)

	svcs := h.serviceGetter.ServicesFromContext(ctx)
	classifierMsgs, newFlowLabels := h.classifier.Classify(ctx, svcs, ctrlPt, mergedFlowLabels, input)

	for key, fl := range newFlowLabels {
		cleanValue := sanitizeBaggageHeaderValue(fl.Value)
		fl.Value = cleanValue
		newFlowLabels[key] = fl
	}

	// Add new flow labels to baggage
	newHeaders, err := h.propagator.Inject(newFlowLabels, existingHeaders)
	if err != nil {
		// TODO(krdln) metrics
		log.Sample(failedBaggageInjectionSampler).
			Warn().Err(err).Msg("Failed to inject baggage into headers")
	}

	// Make the freshly created flow labels available to flowcontrol.
	// Newly created flow labels can overwrite existing flow labels.
	flowlabel.Merge(mergedFlowLabels, newFlowLabels)

	// Ask flow control service for Ok/Deny
	checkResponse := h.fcHandler.CheckRequest(ctx,
		iface.RequestContext{
			FlowLabels:   mergedFlowLabels,
			ControlPoint: ctrlPt,
			Services:     svcs,
		},
	)
	checkResponse.ClassifierInfos = classifierMsgs
	// Set telemetry_flow_labels in the CheckResponse
	checkResponse.TelemetryFlowLabels = mergedFlowLabels.TelemetryLabels()
	// add control point type
	checkResponse.TelemetryFlowLabels[otelconsts.ApertureControlPointTypeLabel] = otelconsts.HTTPControlPoint

	resp := createExtAuthzResponse(checkResponse)

	switch checkResponse.DecisionType {
	case flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED:
		resp.Status = &status.Status{
			Code: int32(code.Code_OK),
		}
		resp.HttpResponse = &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: newHeaders,
			},
		}
	case flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED:
		resp.Status = &status.Status{
			Code: int32(code.Code_UNAVAILABLE),
		}
		var statusCode typev3.StatusCode
		switch checkResponse.RejectReason {
		case flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED:
			statusCode = typev3.StatusCode_TooManyRequests
		case flowcontrolv1.CheckResponse_REJECT_REASON_NO_TOKENS:
			statusCode = typev3.StatusCode_ServiceUnavailable
		case flowcontrolv1.CheckResponse_REJECT_REASON_REGULATED:
			statusCode = typev3.StatusCode_Forbidden
		default:
			log.Bug().Stringer("reason", checkResponse.RejectReason).Msg("Unexpected reject reason")
		}
		deniedHTTPResponse := &authv3.DeniedHttpResponse{
			Status: &typev3.HttpStatus{
				Code: statusCode,
			},
		}
		if checkResponse.WaitTime != nil {
			deniedHTTPResponse.Headers = append(
				deniedHTTPResponse.Headers,
				waitTimeToRetryAfter(checkResponse.WaitTime),
			)
			// Clear to avoid redundancy, as we're translating it into header.
			// Logs processor doesn't read it and Envoy doesn't peek into
			// CheckResponse.
			checkResponse.WaitTime = nil
		}
		resp.HttpResponse = &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: deniedHTTPResponse,
		}
	default:
		return nil, grpc.Bug().Stringer("type", checkResponse.DecisionType).
			Msg("unexpected decision type")
	}

	return resp, nil
}

func waitTimeToRetryAfter(waitTime *durationpb.Duration) *corev3.HeaderValueOption {
	seconds := waitTime.Seconds
	// Retry-after header doesn't have full second resolution, we need to round
	// up to avoid clients retrying too early.
	if waitTime.Nanos != 0 {
		seconds += 1
	}
	return &corev3.HeaderValueOption{
		Header: &corev3.HeaderValue{
			Key:   "retry-after",
			Value: strconv.FormatInt(seconds, 10),
		},
	}
}

func authzRequestToCheckHTTPRequest(
	req *authv3.CheckRequest,
	controlPoint string,
) *flowcontrolhttpv1.CheckHTTPRequest {
	checkHTTPReq := &flowcontrolhttpv1.CheckHTTPRequest{
		ControlPoint: controlPoint,
	}

	if http := req.GetAttributes().GetRequest().GetHttp(); http != nil {
		httpRequest := &flowcontrolhttpv1.CheckHTTPRequest_HttpRequest{
			Method:   http.GetMethod(),
			Headers:  http.GetHeaders(),
			Path:     http.GetPath(),
			Host:     http.GetHost(),
			Scheme:   http.GetScheme(),
			Size:     http.GetSize(),
			Protocol: http.GetProtocol(),
			Body:     http.GetBody(),
		}
		checkHTTPReq.Request = httpRequest
	}

	src := req.GetAttributes().GetSource()
	if src != nil {
		srcSocketAddr := src.GetAddress().GetSocketAddress()
		checkHTTPReq.Source = &flowcontrolhttpv1.SocketAddress{
			Address:  srcSocketAddr.GetAddress(),
			Port:     srcSocketAddr.GetPortValue(),
			Protocol: flowcontrolhttpv1.SocketAddress_Protocol(srcSocketAddr.GetProtocol()),
		}
	}

	dst := req.GetAttributes().GetDestination()
	if dst != nil {
		dstSocketAddr := dst.GetAddress().GetSocketAddress()
		checkHTTPReq.Destination = &flowcontrolhttpv1.SocketAddress{
			Address:  dstSocketAddr.GetAddress(),
			Port:     dstSocketAddr.GetPortValue(),
			Protocol: flowcontrolhttpv1.SocketAddress_Protocol(dstSocketAddr.GetProtocol()),
		}
	}

	return checkHTTPReq
}
