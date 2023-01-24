package awsgateway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/servicegetter"
	"github.com/fluxninja/aperture/pkg/utils"
)

const (
	requestLabelPrefix       = "http."
	requestLabelHeaderPrefix = "http.request.header."
)

// Handler implements the flowcontrol.v1 Service
//
// It also accepts a pointer to an EntityCache for services lookup.
type Handler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer
	serviceGetter servicegetter.ServiceGetter
	classifier    *classification.ClassificationEngine
	fcHandler     check.HandlerWithValues
}

// NewHandler creates new authorization handler for authz api
//
// Authz will use the given classifier to inject flow labels and return them as
// metadata in the response to the Check calls.
func NewHandler(
	classifier *classification.ClassificationEngine,
	serviceGetter servicegetter.ServiceGetter,
	fcHandler check.HandlerWithValues,
) *Handler {
	return &Handler{
		classifier:    classifier,
		serviceGetter: serviceGetter,
		fcHandler:     fcHandler,
	}
}

// AWSGatewayCheck .
func (h *Handler) AWSGatewayCheck(ctx context.Context, req *flowcontrolv1.AWSGatewayCheckRequest) (*httpbody.HttpBody, error) {
	// record the start time of the request
	start := time.Now()

	event := &events.APIGatewayV2CustomAuthorizerV2Request{}
	err := json.Unmarshal([]byte(req.Payload), event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal payload in request")
		return nil, err
	}

	svcs := []string{event.RouteArn}

	input := inputFromPayload(event)

	requestFlowLabels := flowLabelsFromPayload(event)

	classifierInfos, newLabels := h.classifier.Classify(ctx, svcs, "ingress", requestFlowLabels.ToPlainMap(), input)

	checkResponse := h.fcHandler.CheckWithValues(ctx, svcs, "ingress", newLabels.ToPlainMap())

	checkResponse.ClassifierInfos = classifierInfos

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

func inputFromPayload(event *events.APIGatewayV2CustomAuthorizerV2Request) map[string]interface{} {
	input := map[string]interface{}{}

	input["httpMethod"] = event.RequestContext.HTTP.Method
	input["path"] = event.RequestContext.HTTP.Path
	input["queryStringParameters"] = event.QueryStringParameters
	input["headers"] = event.Headers

	return input
}

func flowLabelsFromPayload(event *events.APIGatewayV2CustomAuthorizerV2Request) flowlabel.FlowLabels {
	flowLabels := make(flowlabel.FlowLabels)

	flowLabels[requestLabelPrefix+"method"] = flowlabel.FlowLabelValue{
		Value:     event.RequestContext.HTTP.Method,
		Telemetry: true,
	}
	flowLabels[requestLabelPrefix+"target"] = flowlabel.FlowLabelValue{
		Value:     event.RequestContext.HTTP.Path,
		Telemetry: true,
	}
	flowLabels[requestLabelPrefix+"host"] = flowlabel.FlowLabelValue{
		Value:     event.RequestContext.HTTP.SourceIP,
		Telemetry: true,
	}
	flowLabels[requestLabelPrefix+"flavor"] = flowlabel.FlowLabelValue{
		Value:     utils.CanonicalizeOtelHTTPFlavor(event.RequestContext.HTTP.Protocol),
		Telemetry: true,
	}

	headers := event.Headers

	if headers != nil {
		flowLabels[requestLabelPrefix+"scheme"] = flowlabel.FlowLabelValue{
			Value:     event.Headers["x-forwarded-proto"],
			Telemetry: true,
		}
		delete(headers, "x-forwarded-proto")
		flowLabels[requestLabelPrefix+"request_content_length"] = flowlabel.FlowLabelValue{
			Value:     event.Headers["content-length"],
			Telemetry: false,
		}
		delete(headers, "content-length")

		for k, v := range headers {
			if strings.HasPrefix(k, ":") {
				// Headers starting with `:` are pseudoheaders, so we don't add
				// them.  We don't lose anything, as these values are already
				// available as labels pulled from dedicated fields of
				// Request.Http.
				continue
			}
			flowLabels[requestLabelHeaderPrefix+utils.CanonicalizeOtelHeaderKey(k)] = flowlabel.FlowLabelValue{
				Value:     v,
				Telemetry: false,
			}
		}
	}

	return flowLabels
}
