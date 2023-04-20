package check

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
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
	CheckRequest(
		context.Context,
		iface.RequestContext,
	) *flowcontrolv1.CheckResponse
}

// CheckRequest makes decision using collected inferred fields from authz or Handler.
func (h *Handler) CheckRequest(
	ctx context.Context,
	requestContext iface.RequestContext,
) *flowcontrolv1.CheckResponse {
	checkResponse := h.engine.ProcessRequest(ctx, requestContext)
	h.metrics.CheckResponse(checkResponse.DecisionType, checkResponse.GetRejectReason())
	return checkResponse
}

// Check is the Check method of Flow Control service returns the allow/deny decisions of
// whether to accept the traffic after running the algorithms.
func (h *Handler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	// record the start time of the request
	start := time.Now()

	// handle empty labels
	labels := req.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	services := h.serviceGetter.ServicesFromContext(ctx)

	// CheckRequest already pushes result to metrics
	resp := h.CheckRequest(
		ctx,
		iface.RequestContext{
			FlowLabels:   labels,
			ControlPoint: req.ControlPoint,
			Services:     services,
		},
	)
	end := time.Now()
	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)
	resp.TelemetryFlowLabels = labels
	// add control point type
	resp.TelemetryFlowLabels[otelconsts.ApertureControlPointTypeLabel] = otelconsts.FeatureControlPoint
	return resp, nil
}
