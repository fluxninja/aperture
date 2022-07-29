package authz

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
	"google.golang.org/protobuf/types/known/structpb"

	flowcontrolv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	authz_baggage "github.com/FluxNinja/aperture/pkg/authz/baggage"
	"github.com/FluxNinja/aperture/pkg/classification"
	"github.com/FluxNinja/aperture/pkg/entitycache"
	"github.com/FluxNinja/aperture/pkg/flowcontrol"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/selectors"
	"github.com/FluxNinja/aperture/pkg/services"
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
		classifier: classifier,
		entities:   entityCache,
		propagator: authz_baggage.W3Baggage{},
		fcHandler:  fcHandler,
	}
}

// Handler implements envoy.service.auth.v3.Authorization and handles Check call.
type Handler struct {
	entities   *entitycache.EntityCache
	classifier *classification.Classifier
	propagator authz_baggage.Propagator
	fcHandler  flowcontrol.HandlerWithValues
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

	var direction selectors.TrafficDirection
	headers, _ := metadata.FromIncomingContext(ctx)
	if dirHeader, exists := headers["traffic-direction"]; exists && len(dirHeader) > 0 {
		switch dirHeader[0] {
		case "INBOUND":
			direction = selectors.Ingress
		case "OUTBOUND":
			direction = selectors.Egress
		default:
			log.Warn().Str("traffic-direction", dirHeader[0]).Msg(
				"invalid traffic-direction header",
			)
			return nil, errors.New("invalid traffic-direction")
		}
	} else {
		log.Warn().Msg("traffic-direction not set, assuming inbound")
		direction = selectors.Ingress
	}

	rpcPeer, peerExists := peer.FromContext(ctx)
	if !peerExists {
		log.Warn().Msg("Cannot get client address")
		return nil, errors.New("failed to get client address")
	}
	clientIP := strings.Split(rpcPeer.Addr.String(), ":")[0]

	var svcs []services.ServiceID
	var err error
	if h.entities != nil {
		entity := h.entities.GetByIP(clientIP)
		if entity == nil {
			log.Warn().Str("IP", clientIP).Msg("No entity in cache")
			return nil, fmt.Errorf("missing entity for %s", clientIP)
		}
		svcs, err = entitycache.ServiceIDsFromEntity(entity)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot get services from entity")
			return nil, fmt.Errorf("failed to get services from entity: %v", err)
		}
	} else {
		// TODO: should not have a fallback, always expect entity for consistent experience
		log.Warn().Msg("No entity cache, guessing service id based on Host header")
		svcs = []services.ServiceID{guessDstService(req)}
	}

	logger := logging.New().WithFields(map[string]interface{}{"rego": "input"})
	input, err := envoyauth.RequestToInput(req, logger, nil)
	if err != nil {
		return nil, fmt.Errorf("converting raw input into rego input failed: %v", err)
	}

	inputValue, err := ast.InterfaceToValue(input)
	if err != nil {
		return nil, fmt.Errorf("converting rego input to value failed: %v", err)
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
		return nil, fmt.Errorf("failed to classify: %v", err)
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
	fcResponse := h.fcHandler.CheckWithValues(
		ctx,
		selectors.ControlPoint{Traffic: direction},
		svcs,
		labelsForMatching,
	)
	if err != nil {
		log.Warn().Err(err).Msg("Flow control Check had an error")
	}

	flowLabels := mergeFlowLabels(oldFlowLabels, newFlowLabels)
	resp := ext_authz.CheckResponse{
		// Put all non-hidden flow labels and policy details as dynamic metadata
		DynamicMetadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"fn.flow":              flowLabelsAsPbValueForTelemetry(flowLabels),
				"fn.limiter_decisions": limiterDecisionsAsPbValueForTelemetry(fcResponse.LimiterDecisions),
				"fn.fluxmeters":        fluxmeterIDsAsPbValueForTelemetry(fcResponse.FluxMeterIds),
			},
		},
	}

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

	return &resp, nil
}

func guessDstService(req *ext_authz.CheckRequest) services.ServiceID {
	host := req.GetAttributes().GetRequest().GetHttp().GetHost()
	host = strings.Split(host, ":")[0]
	hostParts := strings.Split(host, ".")
	ns := "default"
	if len(hostParts) >= 2 {
		ns = hostParts[1]
	}
	return services.ServiceID{
		AgentGroup: "default",
		Namespace:  ns,
		Service:    hostParts[0],
	}
}

// Functions below transform our classes/proto to structpb.Value required to be sent
// via DynamicMetadata
// "The External Authorization filter supports emitting dynamic metadata as an opaque google.protobuf.Struct."
// from envoy documentation

// TODO (hasit): The following marshaling functions
// 1. should be 1 and generic
// 2. should be moved to a common package along with the corresponding unmarshal method

func flowLabelsAsPbValueForTelemetry(m classification.FlowLabels) *structpb.Value {
	fields := make(map[string]*structpb.Value, len(m))
	for k, fl := range m {
		if fl.Flags.Hidden {
			continue
		}
		fields[k] = structpb.NewStringValue(fl.Value)
	}
	return structpb.NewStructValue(&structpb.Struct{Fields: fields})
}

func fluxmeterIDsAsPbValueForTelemetry(fluxmeters []string) *structpb.Value {
	values := make([]*structpb.Value, 0, len(fluxmeters))
	for _, fluxmeter := range fluxmeters {
		values = append(values, structpb.NewStringValue(fluxmeter))
	}
	return structpb.NewListValue(&structpb.ListValue{Values: values})
}

func limiterDecisionsAsPbValueForTelemetry(limiterDecisions []*flowcontrolv1.LimiterDecision) *structpb.Value {
	policies := make([]*structpb.Value, len(limiterDecisions))
	for i, ld := range limiterDecisions {
		var structVal *structpb.Value

		switch decision := ld.Decision.(type) {
		case *flowcontrolv1.LimiterDecision_RateLimiterDecision_:
			structVal = structpb.NewStructValue(&structpb.Struct{
				Fields: map[string]*structpb.Value{
					"policy_name":     structpb.NewStringValue(decision.RateLimiterDecision.PolicyName),
					"component_index": structpb.NewNumberValue(float64(decision.RateLimiterDecision.ComponentIndex)),
					"dropped":         structpb.NewBoolValue(ld.Dropped),
					"reason":          structpb.NewStringValue(ld.Reason.Enum().String()),
				},
			})
		case *flowcontrolv1.LimiterDecision_ConcurrencyLimiterDecision:
			structVal = structpb.NewStructValue(&structpb.Struct{
				Fields: map[string]*structpb.Value{
					"policy_name":     structpb.NewStringValue(decision.ConcurrencyLimiterDecision.PolicyName),
					"component_index": structpb.NewNumberValue(float64(decision.ConcurrencyLimiterDecision.ComponentIndex)),
					"workload":        structpb.NewStringValue(decision.ConcurrencyLimiterDecision.Workload),
					"dropped":         structpb.NewBoolValue(ld.Dropped),
					"reason":          structpb.NewStringValue(ld.Reason.Enum().String()),
				},
			})
		}

		policies[i] = structVal
	}
	return structpb.NewListValue(&structpb.ListValue{Values: policies})
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
