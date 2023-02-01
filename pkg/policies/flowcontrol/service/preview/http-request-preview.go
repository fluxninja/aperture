package preview

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	flowpreviewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
)

// HTTPRequestsPreviewRequest holds the samples while the preview is being generated.
type HTTPRequestsPreviewRequest struct {
	mutex           sync.Mutex
	flowSelector    *policylangv1.FlowSelector
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

// GetFlowSelector returns the flow selector.
func (r *HTTPRequestsPreviewRequest) GetFlowSelector() *policylangv1.FlowSelector {
	return r.flowSelector
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
	config.SetDefaults(req)
	if err := config.ValidateStruct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// generate a unique ID for the preview request
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
			LabelMatcher: req.LabelMatcher,
		},
	}

	// create a new request object
	hr := &HTTPRequestsPreviewRequest{
		previewID:       previewID,
		flowSelector:    flowSelector,
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
