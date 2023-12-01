package preview

import (
	"context"
	"sync"

	"github.com/google/uuid"

	flowpreviewv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/preview/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/consts"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
)

type labelPreviewRequest struct {
	mutex           sync.Mutex
	previewDoneCtx  context.Context
	previewResponse *flowpreviewv1.PreviewFlowLabelsResponse
	previewDone     context.CancelFunc
	previewID       iface.PreviewID
	selectors       []*policylangv1.Selector
	samples         int64
}

// GetPreviewID returns the preview ID for this request.
func (r *labelPreviewRequest) GetPreviewID() iface.PreviewID {
	return r.previewID
}

// GetSelectors returns the selectors for this request.
func (r *labelPreviewRequest) GetSelectors() []*policylangv1.Selector {
	return r.selectors
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
	// generate a unique ID for this request
	previewID := iface.PreviewID{
		RequestID: uuid.New().String(),
	}

	service := consts.AnyService
	if req.Service != "" {
		service = req.Service
	}

	selectors := []*policylangv1.Selector{
		{
			ControlPoint: req.ControlPoint,
			LabelMatcher: req.LabelMatcher,
			AgentGroup:   h.agentGroup,
			Service:      service,
		},
	}

	lr := &labelPreviewRequest{
		previewID:       previewID,
		selectors:       selectors,
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
