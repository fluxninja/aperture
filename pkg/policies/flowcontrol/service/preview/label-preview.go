package preview

import (
	"context"
	"fmt"
	"sync"

	flowpreviewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/google/uuid"
)

type labelPreviewRequest struct {
	mutex           sync.Mutex
	flowSelector    *policylangv1.FlowSelector
	previewResponse *flowpreviewv1.PreviewFlowLabelsResponse
	previewDoneCtx  context.Context
	previewDone     context.CancelFunc
	previewID       iface.PreviewID
	samples         int64
}

// GetPreviewID returns the preview ID for this request.
func (r *labelPreviewRequest) GetPreviewID() iface.PreviewID {
	return r.previewID
}

// GetFlowSelector returns the flow selector for this request.
func (r *labelPreviewRequest) GetFlowSelector() *policylangv1.FlowSelector {
	return r.flowSelector
}

// AddLabelPreview adds a label preview to the response.
func (r *labelPreviewRequest) AddLabelPreview(labels map[string]string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.samples > 0 {
		r.previewResponse.Samples = append(r.previewResponse.Samples,
			&flowpreviewv1.PreviewFlowLabelsResponse_FlowLabels{
				Labels: labels,
			})
		r.samples--
		if r.samples == 0 {
			r.previewDone()
		}
	}
}

// PreviewFlowLabels implements flowpreview.v1.PreviewFlowLabels.
func (h *Handler) PreviewFlowLabels(ctx context.Context, req *flowpreviewv1.PreviewRequest) (*flowpreviewv1.PreviewFlowLabelsResponse, error) {
	if req.Samples < 1 {
		return nil, fmt.Errorf("invalid number of samples: %d", req.Samples)
	}

	// generate a unique ID for this request
	previewID := iface.PreviewID{
		RequestID: uuid.New().String(),
	}

	flowSelector := &policylangv1.FlowSelector{
		ServiceSelector: &policylangv1.ServiceSelector{
			AgentGroup: h.agentGroup,
			Service:    req.Service,
		},
		FlowMatcher: &policylangv1.FlowMatcher{
			ControlPoint: req.ControlPoint,
		},
	}

	lr := &labelPreviewRequest{
		previewID:       previewID,
		flowSelector:    flowSelector,
		previewResponse: &flowpreviewv1.PreviewFlowLabelsResponse{},
		samples:         req.Samples,
	}

	// make a child context that will be canceled when the preview is done
	lr.previewDoneCtx, lr.previewDone = context.WithCancel(ctx)
	defer lr.previewDone()

	// register the label preview request
	err := h.engine.RegisterLabelPreview(lr)
	if err != nil {
		return nil, err
	}

	// wait for the preview to be done
	<-lr.previewDoneCtx.Done()

	// unregister the label preview request
	err = h.engine.UnregisterLabelPreview(lr)
	if err != nil {
		log.Errorf("failed to unregister label preview request: %v", err)
	}

	log.Info().Msgf("Generated preview. Samples: %d", len(lr.previewResponse.Samples))

	// return the preview response
	return lr.previewResponse, nil
}
