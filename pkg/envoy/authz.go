package envoy

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/open-policy-agent/opa-envoy-plugin/envoyauth"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/logging"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/classification"
	"github.com/fluxninja/aperture/pkg/entitycache"
	authz_baggage "github.com/fluxninja/aperture/pkg/envoy/baggage"
	"github.com/fluxninja/aperture/pkg/flowcontrol"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
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
	classifier *classification.Classifier,
	entityCache *entitycache.EntityCache,
	fcHandler flowcontrol.HandlerWithValues,
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
	classifier  *classification.Classifier
	propagator  authz_baggage.Propagator
	fcHandler   flowcontrol.HandlerWithValues
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
	log.Trace().Msg("Classifier.Check()")

	createExtAuthzResponse := func(fcResponse *flowcontrolv1.CheckResponse, authzResponse *flowcontrolv1.AuthzResponse, flowLabels classification.FlowLabels) *ext_authz.CheckResponse {
		if fcResponse == nil {
			fcResponse = &flowcontrolv1.CheckResponse{
				DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_UNSPECIFIED,
				DecisionReason: &flowcontrolv1.DecisionReason{
					ErrorReason:  flowcontrolv1.DecisionReason_ERROR_REASON_UNSPECIFIED,
					RejectReason: flowcontrolv1.DecisionReason_REJECT_REASON_UNSPECIFIED,
				},
				LimiterDecisions: nil,
				FluxMeters:       nil,
			}
		}
		marshalledCheckResponse, err := protoMessageAsPbValue(fcResponse)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to marshal check response")
			return nil
		}

		if authzResponse == nil {
			authzResponse = &flowcontrolv1.AuthzResponse{
				Status: flowcontrolv1.AuthzResponse_STATUS_NO_ERROR,
			}
		}
		marshalledAuthzResponse, err := protoMessageAsPbValue(authzResponse)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to marshal authz response")
			return nil
		}

		if flowLabels == nil {
			flowLabels = make(classification.FlowLabels)
		}
		marshalledFlowLabels := flowLabelsAsPbValueForTelemetry(flowLabels)

		log.Trace().Interface("fcResponse", marshalledCheckResponse).Interface("authzResponse", marshalledAuthzResponse).Interface("flowLabels", marshalledFlowLabels).Msg("Created ext_authz.CheckResponse")
		return &ext_authz.CheckResponse{
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					otelcollector.MarshalledLabelsLabel:        marshalledFlowLabels,
					otelcollector.MarshalledCheckResponseLabel: marshalledCheckResponse,
					otelcollector.MarshalledAuthzResponseLabel: marshalledAuthzResponse,
				},
			},
		}
	}

	var direction selectors.TrafficDirection
	headers, _ := metadata.FromIncomingContext(ctx)
	if dirHeader, exists := headers["traffic-direction"]; exists && len(dirHeader) > 0 {
		switch dirHeader[0] {
		case "INBOUND":
			direction = selectors.Ingress
		case "OUTBOUND":
			direction = selectors.Egress
		default:
			log.Warn().Str("traffic-direction", dirHeader[0]).Msg("invalid traffic-direction header")
			authzResponse := &flowcontrolv1.AuthzResponse{
				Status: flowcontrolv1.AuthzResponse_STATUS_INVALID_TRAFFIC_DIRECTION,
			}
			resp := createExtAuthzResponse(nil, authzResponse, nil)
			return resp, errors.New("invalid traffic-direction")
		}
	} else {
		log.Warn().Msg("traffic-direction not set, assuming inbound")
		direction = selectors.Ingress
	}

	var svcs []services.ServiceID
	var err error

	rpcPeer, peerExists := peer.FromContext(ctx)
	if peerExists {
		clientIP := strings.Split(rpcPeer.Addr.String(), ":")[0]
		if h.entityCache != nil {
			entity := h.entityCache.GetByIP(clientIP)
			svcs = entitycache.ServiceIDsFromEntity(entity)
		} else {
			// TODO: should not have a fallback, always expect entity for consistent experience
			log.Warn().Msg("No entity cache, guessing ServiceID based on Host header")
			svcs = []services.ServiceID{guessDstService(req)}
		}
	}

	logger := logging.New().WithFields(map[string]interface{}{"rego": "input"})

	input, err := envoyauth.RequestToInput(req, logger, nil)
	if err != nil {
		authzResponse := &flowcontrolv1.AuthzResponse{
			Status: flowcontrolv1.AuthzResponse_STATUS_CONVERT_TO_MAP_STRUCT,
		}
		resp := createExtAuthzResponse(nil, authzResponse, nil)
		return resp, fmt.Errorf("converting raw input into rego input failed: %v", err)
	}

	inputValue, err := ast.InterfaceToValue(input)
	if err != nil {
		authzResponse := &flowcontrolv1.AuthzResponse{
			Status: flowcontrolv1.AuthzResponse_STATUS_CONVERT_TO_REGO_AST,
		}
		resp := createExtAuthzResponse(nil, authzResponse, nil)
		return resp, fmt.Errorf("converting rego input to value failed: %v", err)
	}

	// Extract previous flow labels from headers
	existingHeaders := authz_baggage.Headers(req.GetAttributes().GetRequest().GetHttp().GetHeaders())
	oldFlowLabels := h.propagator.Extract(existingHeaders)

	labelsForMatching := selectors.NewLabels(
		selectors.LabelSources{
			Flow:    oldFlowLabels.ToPlainMap(),
			Request: req.GetAttributes().GetRequest(),
		},
	)

	newFlowLabels, err := h.classifier.Classify(ctx, svcs, labelsForMatching, direction, inputValue)
	if err != nil {
		authzResponse := &flowcontrolv1.AuthzResponse{
			Status: flowcontrolv1.AuthzResponse_STATUS_CLASSIFY_FLOW_LABELS,
		}
		resp := createExtAuthzResponse(nil, authzResponse, nil)
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
		log.Warn().Err(err).Msg("Failed to inject baggage into headers")
	}

	// Make the freshly created flow labels available to flowcontrol.
	labelsForMatching = labelsForMatching.CombinedWith(selectors.LabelSources{
		Flow: newFlowLabels.ToPlainMap(),
	})
	// Ask flow control service for Ok/Deny
	fcResponse := h.fcHandler.CheckWithValues(ctx, selectors.ControlPoint{Traffic: direction}, svcs, labelsForMatching)

	flowLabels := mergeFlowLabels(oldFlowLabels, newFlowLabels)

	resp := createExtAuthzResponse(fcResponse, nil, flowLabels)

	// Check if fcResponse error is set
	if fcResponse.DecisionType != flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
		resp.Status = &status.Status{
			Code: int32(code.Code_OK),
		}
		resp.HttpResponse = &ext_authz.CheckResponse_OkResponse{
			OkResponse: &ext_authz.OkHttpResponse{
				Headers: newHeaders,
			},
		}
	} else {
		// TODO add rate limiting headers etc.
		resp.Status = &status.Status{
			Code: int32(code.Code_UNAVAILABLE),
		}
		resp.HttpResponse = &ext_authz.CheckResponse_DeniedResponse{
			DeniedResponse: &ext_authz.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_ServiceUnavailable,
				},
			},
		}
	}

	return resp, nil
}

func guessDstService(req *ext_authz.CheckRequest) services.ServiceID {
	host := req.GetAttributes().GetRequest().GetHttp().GetHost()
	host = strings.Split(host, ":")[0]
	return services.ServiceID{
		Service: host,
	}
}

// Functions below transform our classes/proto to structpb.Value required to be sent
// via DynamicMetadata
// "The External Authorization filter supports emitting dynamic metadata as an opaque google.protobuf.Struct."
// from envoy documentation

func flowLabelsAsPbValueForTelemetry(labels classification.FlowLabels) *structpb.Value {
	fields := make(map[string]*structpb.Value, len(labels))
	for k, v := range labels {
		if v.Flags.Hidden {
			continue
		}
		fields[k] = structpb.NewStringValue(v.Value)
	}
	return structpb.NewStructValue(&structpb.Struct{Fields: fields})
}

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
func mergeFlowLabels(first, second classification.FlowLabels) classification.FlowLabels {
	if first == nil {
		return second
	}
	for k, v := range second {
		first[k] = v
	}
	return first
}
