package envoy

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/open-policy-agent/opa-envoy-plugin/envoyauth"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/logging"
	"github.com/rs/zerolog"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/api/base"
	authz_baggage "github.com/fluxninja/aperture/pkg/policies/flowcontrol/api/envoy/baggage"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
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
	fcHandler base.HandlerWithValues,
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
	fcHandler   base.HandlerWithValues
}

var baggageSanitizeRegex *regexp.Regexp = regexp.MustCompile(`[\s\\\/;",]`)

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
		marshalledCheckResponse, err := protoMessageAsPbValue(checkResponse)
		if err != nil {
			log.Sample(zerolog.Sometimes).Error().Err(err).Msg("Failed to marshal check response")
			return nil
		}

		// record the end time of the request
		end := time.Now()
		checkResponse.Start = timestamppb.New(start)
		checkResponse.End = timestamppb.New(end)

		return &ext_authz.CheckResponse{
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					otelcollector.ApertureCheckResponseLabel: marshalledCheckResponse,
				},
			},
		}
	}

	ctrlPt := selectors.NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_UNKNOWN, "")
	headers, _ := metadata.FromIncomingContext(ctx)
	if dirHeader, exists := headers["traffic-direction"]; exists && len(dirHeader) > 0 {
		switch dirHeader[0] {
		case "INBOUND":
			ctrlPt = selectors.NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_INGRESS, "")
		case "OUTBOUND":
			ctrlPt = selectors.NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_EGRESS, "")
		default:
			log.Sample(zerolog.Sometimes).Error().Str("traffic-direction", dirHeader[0]).Msg("invalid traffic-direction header")
			checkResponse := &flowcontrolv1.CheckResponse{
				Error:            flowcontrolv1.CheckResponse_ERROR_INVALID_TRAFFIC_DIRECTION,
				ControlPointInfo: ctrlPt.ToControlPointInfoProto(),
			}
			resp := createExtAuthzResponse(checkResponse)
			return resp, errors.New("invalid traffic-direction")
		}
	} else {
		log.Sample(zerolog.Sometimes).Error().Msg("traffic-direction not set")
		checkResponse := &flowcontrolv1.CheckResponse{
			Error: flowcontrolv1.CheckResponse_ERROR_MISSING_TRAFFIC_DIRECTION,
		}
		resp := createExtAuthzResponse(checkResponse)
		return resp, errors.New("invalid traffic-direction")
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
		checkResponse := &flowcontrolv1.CheckResponse{
			Error:    flowcontrolv1.CheckResponse_ERROR_CONVERT_TO_MAP_STRUCT,
			Services: svcs,
		}
		resp := createExtAuthzResponse(checkResponse)
		return resp, fmt.Errorf("converting raw input into rego input failed: %v", err)
	}

	inputValue, err := ast.InterfaceToValue(input)
	if err != nil {
		checkResponse := &flowcontrolv1.CheckResponse{
			Error:    flowcontrolv1.CheckResponse_ERROR_CONVERT_TO_REGO_AST,
			Services: svcs,
		}
		resp := createExtAuthzResponse(checkResponse)
		return resp, fmt.Errorf("converting rego input to value failed: %v", err)
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

	classifierMsgs, newFlowLabels, err := h.classifier.Classify(ctx, svcs, ctrlPt, mergedFlowLabels.ToPlainMap(), inputValue)
	if err != nil {
		checkResponse := &flowcontrolv1.CheckResponse{
			Error:           flowcontrolv1.CheckResponse_ERROR_CLASSIFY,
			Services:        svcs,
			ClassifierInfos: classifierMsgs,
		}
		resp := createExtAuthzResponse(checkResponse)
		return resp, fmt.Errorf("failed to classify: %v", err)
	}

	for key, fl := range newFlowLabels {
		cleanValue := sanitizeBaggageHeaderValue(fl.Value)
		fl.Value = cleanValue
		newFlowLabels[key] = fl
	}

	// Add new flow labels to baggage
	newHeaders, err := h.propagator.Inject(newFlowLabels, existingHeaders)
	if err != nil {
		log.Sample(zerolog.Sometimes).Warn().Err(err).Msg("Failed to inject baggage into headers")
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
			log.Sample(zerolog.Sometimes).Error().Msg("Unexpected reject reason: " + checkResponse.RejectReason.String())
		}
	}

	return resp, nil
}

// Functions below transform our classes/proto to structpb.Value required to be sent
// via DynamicMetadata
// "The External Authorization filter supports emitting dynamic metadata as an opaque google.protobuf.Struct."
// from envoy documentation

func protoMessageAsPbValue(message protoreflect.ProtoMessage) (*structpb.Value, error) {
	mBytes, err := protojson.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal proto message into JSON")
		return nil, err

	}
	mStruct := new(structpb.Struct)
	err = protojson.Unmarshal(mBytes, mStruct)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal JSON bytes into structpb.Struct")
		return nil, err
	}

	return structpb.NewStructValue(mStruct), nil
}

// merges two flow labels maps.
//
// If key exists in both, the value from second one will be taken.
// Nil maps should be handled fine.
/*func mergeFlowLabels(first, second classification.FlowLabels) classification.FlowLabels {
	if first == nil {
		return second
	}
	for k, v := range second {
		first[k] = v
	}
	return first
}*/
