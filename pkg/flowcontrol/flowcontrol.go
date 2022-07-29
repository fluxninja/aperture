package flowcontrol

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	flowcontrolv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/FluxNinja/aperture/pkg/entitycache"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/policies/dataplane/iface"
	"github.com/FluxNinja/aperture/pkg/selectors"
	"github.com/FluxNinja/aperture/pkg/services"
)

// Handler implements the flowcontrol.v1 Service
//
// It also accepts a pointer to an EntityCache for services lookup.
type Handler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer
	entityCache *entitycache.EntityCache
	metrics     Metrics
	engine      iface.EngineAPI
}

// NewHandler creates an empty flowcontrol Handler
//
// It also accepts a pointer to an EntityCache for Infra Labels lookup.
func NewHandler(entityCache *entitycache.EntityCache, metrics Metrics, engine iface.EngineAPI) *Handler {
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
		selectors.ControlPoint,
		[]services.ServiceID,
		selectors.Labels,
	) *flowcontrolv1.CheckResponse
}

// CheckWithValues makes decision using collected inferred fields from authz or Handler.
func (h *Handler) CheckWithValues(
	ctx context.Context,
	controlPoint selectors.ControlPoint,
	serviceIDs []services.ServiceID,
	labels selectors.Labels,
) *flowcontrolv1.CheckResponse {
	log.Trace().Msg("FlowControl.CheckWithValues()")

	checkResponse := h.engine.ProcessRequest(controlPoint, serviceIDs, labels)
	h.metrics.CheckResponse(checkResponse.DecisionType, checkResponse.GetReason())
	return checkResponse
}

// Check is the Check method of Flow Control service returns the allow/deny decisions of
// whether to accept the traffic after running the algorithms.
func (h *Handler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	log.Trace().Msg("FlowControl.Check()")

	rpcPeer, peerExists := peer.FromContext(ctx)
	if !peerExists {
		reason := &flowcontrolv1.Reason{
			Reason: &flowcontrolv1.Reason_ErrorReason_{
				ErrorReason: flowcontrolv1.Reason_ERROR_REASON_BAD_CLIENT_IP,
			},
		}
		h.metrics.CheckResponse(flowcontrolv1.DecisionType_DECISION_TYPE_UNSPECIFIED, reason)
		return &flowcontrolv1.CheckResponse{
			DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_UNSPECIFIED,
			Reason:       reason,
		}, errors.New("failed to get client address")
	}

	clientIP := strings.Split(rpcPeer.Addr.String(), ":")[0]
	_ = grpc.SetHeader(ctx, metadata.Pairs("client-ip", clientIP))

	entity := h.entityCache.GetByIP(clientIP)
	if entity == nil {
		log.Warn().Err(errors.New("no entity in cache")).Msg("failed to get services and labels from entities")
		reason := &flowcontrolv1.Reason{
			Reason: &flowcontrolv1.Reason_ErrorReason_{
				ErrorReason: flowcontrolv1.Reason_ERROR_REASON_ENTITY_LOOKUP_FAILED,
			},
		}
		return &flowcontrolv1.CheckResponse{
			DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_UNSPECIFIED,
			Reason:       reason,
		}, nil
	}

	svcs, err := entitycache.ServiceIDsFromEntity(entity)
	if err != nil {
		log.Error().Err(err).Msg("failed to get services and labels from entities")
		reason := &flowcontrolv1.Reason{
			Reason: &flowcontrolv1.Reason_ErrorReason_{
				ErrorReason: flowcontrolv1.Reason_ERROR_REASON_SERVICE_LOOKUP_FAILED,
			},
		}
		return &flowcontrolv1.CheckResponse{
			DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_UNSPECIFIED,
			Reason:       reason,
		}, nil
	}

	// CheckWithValues already pushes result to metrics
	resp := h.CheckWithValues(
		ctx,
		selectors.ControlPoint{Feature: req.Feature},
		svcs,
		selectors.NewLabels(selectors.LabelSources{
			Flow: req.Labels,
		}),
	)
	return resp, nil
}
