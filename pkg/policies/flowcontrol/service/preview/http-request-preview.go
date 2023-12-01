package preview

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"

	flowpreviewv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/preview/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/consts"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
)

// HTTPRequestsPreviewRequest holds the samples while the preview is being generated.
type HTTPRequestsPreviewRequest struct {
	mutex           sync.Mutex
	selectors       []*policylangv1.Selector
	previewResponse *flowpreviewv1.PreviewHTTPRequestsResponse
	previewDoneCtx  context.Context
	previewDone     context.CancelFunc
	previewID       iface.PreviewID
	samples         int64
}

// GetPreviewID returns the preview ID.
func (r *HTTPRequestsPreviewRequest) GetPreviewID() iface.PreviewID {
	return r.previewID
}

// GetSelectors returns the flow selector.
func (r *HTTPRequestsPreviewRequest) GetSelectors() []*policylangv1.Selector {
	return r.selectors
}

// AddHTTPRequestPreview adds a HTTP request preview to the response.
func (r *HTTPRequestsPreviewRequest) AddHTTPRequestPreview(request map[string]interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.samples > 0 {
		// some fields in request are json.Number and incompatible with structpb
		// so we convert them to float64 by encoding and decoding
		b, err := json.Marshal(request)
		if err != nil {
			log.Errorf("failed to marshal request: %v", err)
			return
		}
		var m map[string]interface{}
		if err = json.Unmarshal(b, &m); err != nil {
			log.Errorf("failed to unmarshal request: %v", err)
			return
		}

		// encode as request as structpb
		structProto, err := structpb.NewStruct(m)
		if err != nil {
			log.Errorf("failed to encode HTTP request as structpb: %v", err)
			return
		}
		r.previewResponse.Samples = append(r.previewResponse.Samples,
			structProto)
		r.samples--
		if r.samples == 0 {
			r.previewDone()
		}
	}
}

// PreviewHTTPRequests implements flowpreview.v1.PreviewHTTPRequests.
func (h *Handler) PreviewHTTPRequests(ctx context.Context, req *flowpreviewv1.PreviewRequest) (*flowpreviewv1.PreviewHTTPRequestsResponse, error) {
	// generate a unique ID for the preview request
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

	// create a new request object
	hr := &HTTPRequestsPreviewRequest{
		previewID:       previewID,
		selectors:       selectors,
		previewResponse: &flowpreviewv1.PreviewHTTPRequestsResponse{},
		samples:         req.Samples,
	}

	// make a context that is canceled when the preview is done
	hr.previewDoneCtx, hr.previewDone = context.WithCancel(ctx)
	defer hr.previewDone()

	// add the request to the classifier
	h.classifier.AddPreview(hr)

	// wait for the preview to complete
	<-hr.previewDoneCtx.Done()

	// remove the request from the classifier
	h.classifier.DropPreview(hr)

	return hr.previewResponse, nil
}
