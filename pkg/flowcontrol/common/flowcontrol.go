package common

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/selectors"
)

// Handler implements the flowcontrol.v1 Service
//
// It also accepts a pointer to an EntityCache for services lookup.
type Handler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer
	entityCache *entitycache.EntityCache
	metrics     Metrics
	engine      iface.Engine
}

// NewHandler creates an empty flowcontrol Handler
//
// It also accepts a pointer to an EntityCache for Infra Labels lookup.
func NewHandler(entityCache *entitycache.EntityCache, metrics Metrics, engine iface.Engine) *Handler {
	return &Handler{
		entityCache: entityCache,
		metrics:     metrics,
		engine:      engine,
	}
}

// HandlerWithValues implements the flowcontrol.v1 service using collected inferred values.
type HandlerWithValues interface {
	CheckWithValues(
		context.Context,
		[]string,
		selectors.ControlPoint,
		map[string]string,
	) *flowcontrolv1.CheckResponse
}

// CheckWithValues makes decision using collected inferred fields from authz or Handler.
func (h *Handler) CheckWithValues(
	ctx context.Context,
	serviceIDs []string,
	controlPoint selectors.ControlPoint,
	labels map[string]string,
) *flowcontrolv1.CheckResponse {
	checkResponse := h.engine.ProcessRequest(ctx, controlPoint, serviceIDs, labels)
	h.metrics.CheckResponse(checkResponse.DecisionType, checkResponse.GetRejectReason(), checkResponse.GetError())
	return checkResponse
}

// Check is the Check method of Flow Control service returns the allow/deny decisions of
// whether to accept the traffic after running the algorithms.
func (h *Handler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	log.Trace().Msg("FlowControl.Check()")
	// record the start time of the request
	start := time.Now()

	var serviceIDs []string

	rpcPeer, peerExists := peer.FromContext(ctx)
	if peerExists {
		clientIP := strings.Split(rpcPeer.Addr.String(), ":")[0]
		entity, err := h.entityCache.GetByIP(clientIP)
		if err == nil {
			serviceIDs = entity.Services
		}
	}

	// CheckWithValues already pushes result to metrics
	resp := h.CheckWithValues(
		ctx,
		serviceIDs,
		selectors.NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_FEATURE, req.Feature),
		req.Labels,
	)
	end := time.Now()
	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)
	return resp, nil
}
