package agentfunctions

import (
	"context"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/preview"
	previewconfig "github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/preview/config"
	"github.com/fluxninja/aperture/pkg/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PreviewHandler is a handler for preview-family of functions.
//
// Simply forwards methods to grpc handler, which provides actual implementation.
type PreviewHandler struct {
	handler   *preview.Handler
	isEnabled bool
}

// ProvidePreviewHandler provides PreviewHandler.
func ProvidePreviewHandler(
	handler *preview.Handler,
	unmarshaller config.Unmarshaller,
) (*PreviewHandler, error) {
	var previewConfig previewconfig.FlowPreviewConfig
	if err := unmarshaller.UnmarshalKey(previewconfig.Key, &previewConfig); err != nil {
		return nil, err
	}

	return &PreviewHandler{
		handler:   handler,
		isEnabled: previewConfig.Enabled,
	}, nil
}

// PreviewFlowLabels previews flow labels on given control point.
func (h *PreviewHandler) PreviewFlowLabels(
	ctx context.Context,
	req *cmdv1.PreviewFlowLabelsRequest,
) (*previewv1.PreviewFlowLabelsResponse, error) {
	if !h.isEnabled {
		return nil, status.Error(codes.FailedPrecondition, "preview disabled")
	}
	// GRPC handlers assume non-nillness of argument.
	if req.Request == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	return h.handler.PreviewFlowLabels(ctx, req.Request)
}

// PreviewHTTPRequests previews HTTP requests on given control point.
func (h *PreviewHandler) PreviewHTTPRequests(
	ctx context.Context,
	req *cmdv1.PreviewHTTPRequestsRequest,
) (*previewv1.PreviewHTTPRequestsResponse, error) {
	if !h.isEnabled {
		return nil, status.Error(codes.FailedPrecondition, "preview disabled")
	}
	// GRPC handlers assume non-nillness of argument.
	if req.Request == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	return h.handler.PreviewHTTPRequests(ctx, req.Request)
}

// RegisterPreviewHandler registers PreviewHandler in handler registry.
func RegisterPreviewHandler(handler *PreviewHandler, registry *rpc.HandlerRegistry) error {
	// Note: Registering also when handler is disabled, so that we can send
	// more specific error code than Unimplemented.
	if err := rpc.RegisterFunction(registry, handler.PreviewFlowLabels); err != nil {
		return err
	}
	if err := rpc.RegisterFunction(registry, handler.PreviewHTTPRequests); err != nil {
		return err
	}
	return nil
}
