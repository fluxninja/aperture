package check

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/servicegetter"
)

// Handler implements the flowcontrol.v1 Service
//
// It also accepts a pointer to an EntityCache for services lookup.
type Handler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer
	serviceGetter servicegetter.ServiceGetter
	metrics       Metrics
	engine        iface.Engine
}

// NewHandler creates a flowcontrol Handler.
func NewHandler(
	serviceGetter servicegetter.ServiceGetter,
	metrics Metrics,
	engine iface.Engine,
) *Handler {
	return &Handler{
		serviceGetter: serviceGetter,
		metrics:       metrics,
		engine:        engine,
	}
}

// HandlerWithValues implements the flowcontrol.v1 service using collected inferred values.
type HandlerWithValues interface {
	CheckWithValues(
		context.Context,
		[]string,
		string,
		map[string]string,
	) *flowcontrolv1.CheckResponse
}

// CheckWithValues makes decision using collected inferred fields from authz or Handler.
func (h *Handler) CheckWithValues(
	ctx context.Context,
	serviceIDs []string,
	controlPoint string,
	labels map[string]string,
) *flowcontrolv1.CheckResponse {
	checkResponse := h.engine.ProcessRequest(ctx, controlPoint, serviceIDs, labels)
	h.metrics.CheckResponse(checkResponse.DecisionType, checkResponse.GetRejectReason())
	return checkResponse
}

// Check is the Check method of Flow Control service returns the allow/deny decisions of
// whether to accept the traffic after running the algorithms.
func (h *Handler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	// record the start time of the request
	start := time.Now()

	// CheckWithValues already pushes result to metrics
	resp := h.CheckWithValues(
		ctx,
		h.serviceGetter.ServicesFromContext(ctx),
		req.ControlPoint,
		req.Labels,
	)
	end := time.Now()
	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)
	return resp, nil
}

// GatewayCheck .
func (h *Handler) GatewayCheck(ctx context.Context, req *flowcontrolv1.GatewayCheckRequest) (*httpbody.HttpBody, error) {
	log.Info().Str("payload", req.Payload).Msg("Received FlowControl Gateway Check request")

	event := &events.APIGatewayV2CustomAuthorizerV2Request{}
	err := json.Unmarshal([]byte(req.Payload), event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal payload in request")
		return nil, err
	}

	svcs := []string{event.RequestContext.DomainName}

	// record the start time of the request
	start := time.Now()

	// TODO: extract fields from the request payload into input.
	input := map[string]string{}
	checkResponse := h.CheckWithValues(ctx, svcs, "gateway", input)

	end := time.Now()
	checkResponse.Start = timestamppb.New(start)
	checkResponse.End = timestamppb.New(end)

	marshalledCheckResponse, err := proto.Marshal(checkResponse)
	if err != nil {
		log.Bug().Err(err).Msg("bug: Failed to marshal check response")
		return nil, err
	}
	checkResponseBase64 := base64.StdEncoding.EncodeToString(marshalledCheckResponse)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(checkResponseBase64),
	}, nil
}
