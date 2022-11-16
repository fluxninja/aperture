package envoy

import (
	"context"
	"encoding/base64"
	"regexp"
	"strings"
	"time"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/open-policy-agent/opa-envoy-plugin/envoyauth"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/logging"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	grpc_status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	authz_baggage "github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/envoy/baggage"
)

// NewHandler creates new authorization handler for authz api
//
// Authz will use the given classifier to inject flow labels and return them as
// metadata in the response to the Check calls
//
// entityCache can be nil. In this case services will be guessed based on Host
// header.  No-entity-cache support is mostly so that authz can be experimented
// with without the need for tagger to run.
func NewHandler(
	classifier *classification.ClassificationEngine,
	entityCache *entitycache.EntityCache,
	fcHandler check.HandlerWithValues,
) *Handler {
	if entityCache == nil {
		log.Warn().Msg("Authz: No entity cache, will guess services based on Host header")
	}
	return &Handler{
		classifier:  classifier,
		entityCache: entityCache,
		propagator:  authz_baggage.W3Baggage{},
		fcHandler:   fcHandler,
	}
}

// Handler implements envoy.service.auth.v3.Authorization and handles Check call.
type Handler struct {
	entityCache *entitycache.EntityCache
	classifier  *classification.ClassificationEngine
	propagator  authz_baggage.Propagator
	fcHandler   check.HandlerWithValues
}

var baggageSanitizeRegex *regexp.Regexp = regexp.MustCompile(`[\s\\\/;",]`)

var (
	missingTrafficDirectionSampler = log.NewRatelimitingSampler()
	invalidTrafficDirectionSampler = log.NewRatelimitingSampler()
	failedReqToInputSampler        = log.NewRatelimitingSampler()
	failedBaggageInjectionSampler  = log.NewRatelimitingSampler()
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
func (h *Handler) Check(ctx context.Context, req *ext_authz.CheckRequest) (*ext_authz.CheckResponse, error) {
	// record the start time of the request
	start := time.Now()

	createExtAuthzResponse := func(checkResponse *flowcontrolv1.CheckResponse) *ext_authz.CheckResponse {
		// We don't care about the particular format we send the CheckResponse,
		// Envoy can treat is as black-box. The only thing we care about is for
		// it to be deserializable by logs processing pipeline.
		// Using protobuf wire format as it's faster to serialize/deserialize
		// than using protojson or roundtripping through structpb.Struct.
		// Additional base64 encoding step is used, as there's no way to push
		// binary data through dynamic metadata and envoy's access log
		// formatter. Overhead of this base64 encoding is small though.
		marshalledCheckResponse, err := proto.Marshal(checkResponse)
		if err != nil {
			log.Bug().Err(err).Msg("bug: Failed to marshal check response")
			return nil
		}
		checkResponseBase64 := base64.StdEncoding.EncodeToString(marshalledCheckResponse)

		// record the end time of the request
		end := time.Now()
		checkResponse.Start = timestamppb.New(start)
		checkResponse.End = timestamppb.New(end)

		return &ext_authz.CheckResponse{
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					otelcollector.ApertureCheckResponseLabel: structpb.NewStringValue(checkResponseBase64),
				},
			},
		}
	}

	var ctrlPt selectors.ControlPoint
	trafficDirectionHeader := ""
	headers, _ := metadata.FromIncomingContext(ctx)
	if dirHeader, exists := headers["traffic-direction"]; exists && len(dirHeader) > 0 {
		trafficDirectionHeader = dirHeader[0]
	} else if dirHeader, exists := req.GetAttributes().GetContextExtensions()["traffic-direction"]; exists && len(dirHeader) > 0 {
		trafficDirectionHeader = dirHeader
	}

	if trafficDirectionHeader != "" {
		switch trafficDirectionHeader {
		case "INBOUND":
			ctrlPt = selectors.NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_INGRESS, "")
		case "OUTBOUND":
			ctrlPt = selectors.NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_EGRESS, "")
		default:
			// TODO(krdln) metrics
			log.Sample(invalidTrafficDirectionSampler).
				Warn().Str("traffic-direction", trafficDirectionHeader).Msg("invalid traffic-direction")
			return nil, grpc_status.Error(grpc_codes.InvalidArgument, "invalid traffic-direction")
		}
	} else {
		// TODO(krdln) metrics
		log.Sample(missingTrafficDirectionSampler).
			Warn().Msg("missing traffic-direction")
		return nil, grpc_status.Error(grpc_codes.InvalidArgument, "missing traffic-direction")
	}

	var svcs []string
	rpcPeer, peerExists := peer.FromContext(ctx)
	if peerExists {
		clientIP := strings.Split(rpcPeer.Addr.String(), ":")[0]
		if h.entityCache != nil {
			entity, err := h.entityCache.GetByIP(clientIP)
			if err == nil {
				svcs = entity.Services
			}
		}
	}

	logger := logging.New().WithFields(map[string]interface{}{"rego": "input"})

	input, err := envoyauth.RequestToInput(req, logger, nil)
	if err != nil {
		// TODO(krdln) This conversion should be made infallible instead.
		// https://github.com/fluxninja/aperture/issues/903
		// TODO(krdln) metrics
		log.Sample(failedReqToInputSampler).
			Warn().Err(err).Msg("converting raw input into rego input failed")
		return nil, grpc_status.Error(grpc_codes.InvalidArgument, "converting raw input into rego input failed")
	}

	inputValue, err := ast.InterfaceToValue(input)
	if err != nil {
		// RequestToInput should never produce anything that's not convertible
		// to ast.Value, so in theory it shouldn't happen.
		// TODO(krdln) metrics
		log.Bug().Err(err).Msg("bug: converting rego input to value failed")
		return nil, grpc_status.Error(grpc_codes.Internal, "converting rego input to value failed")
	}

	// Default flow labels from Authz request
	requestFlowLabels := AuthzRequestToFlowLabels(req.GetAttributes().GetRequest())
	// Extract flow labels from baggage headers
	existingHeaders := authz_baggage.Headers(req.GetAttributes().GetRequest().GetHttp().GetHeaders())
	baggageFlowLabels := h.propagator.Extract(existingHeaders)

	// Merge flow labels from Authz request and baggage headers
	mergedFlowLabels := requestFlowLabels
	// Baggage can overwrite request flow labels
	flowlabel.Merge(mergedFlowLabels, baggageFlowLabels)

	classifierMsgs, newFlowLabels := h.classifier.Classify(ctx, svcs, ctrlPt, mergedFlowLabels.ToPlainMap(), inputValue)

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
	flowLabels := mergedFlowLabels.ToPlainMap()

	// Ask flow control service for Ok/Deny
	checkResponse := h.fcHandler.CheckWithValues(ctx, svcs, ctrlPt, flowLabels)
	checkResponse.ClassifierInfos = classifierMsgs
	// Set telemetry_flow_labels in the CheckResponse
	checkResponse.TelemetryFlowLabels = flowLabels

	resp := createExtAuthzResponse(checkResponse)

	// Check if fcResponse error is set
	if checkResponse.DecisionType != flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED {
		resp.Status = &status.Status{
			Code: int32(code.Code_OK),
		}
		resp.HttpResponse = &ext_authz.CheckResponse_OkResponse{
			OkResponse: &ext_authz.OkHttpResponse{
				Headers: newHeaders,
			},
		}
	} else {
		resp.Status = &status.Status{
			Code: int32(code.Code_UNAVAILABLE),
		}
		if checkResponse.RejectReason == flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED {
			resp.HttpResponse = &ext_authz.CheckResponse_DeniedResponse{
				DeniedResponse: &ext_authz.DeniedHttpResponse{
					Status: &envoy_type.HttpStatus{
						Code: envoy_type.StatusCode_TooManyRequests,
					},
				},
			}
		} else if checkResponse.RejectReason == flowcontrolv1.CheckResponse_REJECT_REASON_CONCURRENCY_LIMITED {
			resp.HttpResponse = &ext_authz.CheckResponse_DeniedResponse{
				DeniedResponse: &ext_authz.DeniedHttpResponse{
					Status: &envoy_type.HttpStatus{
						Code: envoy_type.StatusCode_ServiceUnavailable,
					},
				},
			}
		} else {
			log.Bug().Stringer("reason", checkResponse.RejectReason).Msg("Unexpected reject reason")
		}
	}

	return resp, nil
}
