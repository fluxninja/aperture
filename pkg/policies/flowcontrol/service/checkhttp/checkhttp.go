package checkhttp

import (
	"context"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	flowcontrolhttpv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpc"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	checkhttp_baggage "github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/checkhttp/baggage"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/servicegetter"
	"github.com/open-policy-agent/opa/logging"
	"github.com/open-policy-agent/opa/util"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var baggageSanitizeRegex *regexp.Regexp = regexp.MustCompile(`[\s\\\/;",]`)

var (
	missingControlPointSampler    = log.NewRatelimitingSampler()
	failedBaggageInjectionSampler = log.NewRatelimitingSampler()
)

// NewHandler creates new handler for flowcontrol CheckHTTP
//
// It will use the given classifier to inject flow labels and return them as
// metadata in the response to the Check calls.
func NewHandler(
	classifier *classification.ClassificationEngine,
	serviceGetter servicegetter.ServiceGetter,
	fcHandler check.HandlerWithValues,
) *Handler {
	return &Handler{
		classifier:    classifier,
		serviceGetter: serviceGetter,
		propagator:    checkhttp_baggage.W3Baggage{},
		fcHandler:     fcHandler,
	}
}

// sanitizeBaggageHeaderValue excludes characters that should be url escaped
// Otherwise both baggage.String method and envoy itself will do it.
func sanitizeBaggageHeaderValue(value string) string {
	// All characters allowed except control chars, whitespace, double quote, comma, semicolon, backslash
	// see https://www.w3.org/TR/baggage/#header-content
	cleanValue := baggageSanitizeRegex.ReplaceAll([]byte(value), []byte("-"))
	return string(cleanValue)
}

// Handler implements aperture.flowcontrol.v1.FlowControlServiceHTTP and handles Check call.
type Handler struct {
	flowcontrolhttpv1.UnimplementedFlowControlServiceHTTPServer
	serviceGetter servicegetter.ServiceGetter
	classifier    *classification.ClassificationEngine
	propagator    checkhttp_baggage.Propagator
	fcHandler     check.HandlerWithValues
}

// CheckHTTP is the Check method of Flow Control service returns the allow/deny decisions of
// whether to accept the traffic after running the algorithms.
func (h *Handler) CheckHTTP(ctx context.Context, req *flowcontrolhttpv1.CheckHTTPRequest) (*flowcontrolhttpv1.CheckHTTPResponse, error) {
	// record the start time of the request
	start := time.Now()
	createResponse := func(checkResponse *flowcontrolv1.CheckResponse) *flowcontrolhttpv1.CheckHTTPResponse {
		marshalledCheckResponse, err := proto.Marshal(checkResponse)
		if err != nil {
			log.Bug().Err(err).Msg("bug: Failed to marshal check response")
			return nil
		}
		checkResponseBase64 := base64.StdEncoding.EncodeToString(marshalledCheckResponse)

		// record the end time of the request
		end := time.Now()
		checkResponse.Start = timestamppb.New(start)
		checkResponse.End = timestamppb.New(end)

		return &flowcontrolhttpv1.CheckHTTPResponse{
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					otelconsts.ApertureCheckResponseLabel: structpb.NewStringValue(checkResponseBase64),
				},
			},
		}
	}

	ctrlPt := ""
	headers, _ := metadata.FromIncomingContext(ctx)
	if ctrlPtHeader, exists := headers["control-point"]; exists && len(ctrlPtHeader) > 0 {
		ctrlPt = ctrlPtHeader[0]
	}

	if ctrlPt == "" {
		return nil, grpc.LoggedError(log.Sample(missingControlPointSampler).Warn()).
			Code(codes.InvalidArgument).Msg("missing control-point")
	}

	sourceAddress := req.GetSource().GetAddress()
	sourceSvcs := h.serviceGetter.ServicesFromAddress(sourceAddress)
	if sourceSvcs == nil {
		sourceSvcs = []string{"UNKNOWN"}
	}
	sourceSvcsStr := strings.Join(sourceSvcs, ",")
	destinationAddress := req.GetDestination().GetAddress()
	destinationSvcs := h.serviceGetter.ServicesFromAddress(destinationAddress)
	if destinationSvcs == nil {
		sourceSvcs = []string{"UNKNOWN"}
	}
	destinationSvcsStr := strings.Join(destinationSvcs, ",")

	// make flowlabels from source and destination services
	sdFlowLabels := make(flowlabel.FlowLabels, 2)
	sdFlowLabels[otelconsts.ApertureSourceServiceLabel] = flowlabel.FlowLabelValue{
		Value:     sourceSvcsStr,
		Telemetry: true,
	}
	sdFlowLabels[otelconsts.ApertureDestinationServiceLabel] = flowlabel.FlowLabelValue{
		Value:     destinationSvcsStr,
		Telemetry: true,
	}

	logger := logging.New().WithFields(map[string]interface{}{"rego": "input"})
	input := RequestToInput(req, logger)

	// Default flow labels from request
	requestFlowLabels := CheckHTTPRequestToFlowLabels(req.GetRequest())
	existingHeaders := checkhttp_baggage.Headers(req.GetRequest().GetHeaders())
	baggageFlowLabels := h.propagator.Extract(existingHeaders)

	// Merge flow labels from request and baggage headers
	mergedFlowLabels := requestFlowLabels
	// Baggage can overwrite request flow labels
	flowlabel.Merge(mergedFlowLabels, baggageFlowLabels)
	flowlabel.Merge(mergedFlowLabels, sdFlowLabels)

	classifierMsgs, newFlowLabels := h.classifier.Classify(ctx, destinationSvcs, ctrlPt, mergedFlowLabels.ToPlainMap(), input)

	for key, fl := range newFlowLabels {
		cleanValue := sanitizeBaggageHeaderValue(fl.Value)
		fl.Value = cleanValue
		newFlowLabels[key] = fl
	}

	// Add new flow labels to baggage
	newHeaders, err := h.propagator.Inject(newFlowLabels, existingHeaders)
	if err != nil {
		log.Sample(failedBaggageInjectionSampler).
			Warn().Err(err).Msg("Failed to inject baggage into headers")
	}

	// Make the freshly created flow labels available to flowcontrol.
	// Newly created flow labels can overwrite existing flow labels.
	flowlabel.Merge(mergedFlowLabels, newFlowLabels)
	flowLabels := mergedFlowLabels.ToPlainMap()

	// Ask flow control service for Ok/Deny
	checkResponse := h.fcHandler.CheckWithValues(ctx, destinationSvcs, ctrlPt, flowLabels)
	checkResponse.ClassifierInfos = classifierMsgs
	// Set telemetry_flow_labels in the CheckResponse
	checkResponse.TelemetryFlowLabels = flowLabels
	// add control point type
	checkResponse.TelemetryFlowLabels[otelconsts.ApertureControlPointTypeLabel] = otelconsts.HTTPControlPoint
	checkResponse.TelemetryFlowLabels[otelconsts.ApertureSourceServiceLabel] = strings.Join(sourceSvcs, ",")
	checkResponse.TelemetryFlowLabels[otelconsts.ApertureDestinationServiceLabel] = strings.Join(destinationSvcs, ",")

	resp := createResponse(checkResponse)

	switch checkResponse.DecisionType {
	case flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED:
		resp.Status = &status.Status{
			Code: int32(code.Code_OK),
		}
		resp.HttpResponse = &flowcontrolhttpv1.CheckHTTPResponse_OkResponse{
			OkResponse: &flowcontrolhttpv1.OkHttpResponse{
				Headers: newHeaders,
			},
		}
	case flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED:
		resp.Status = &status.Status{
			Code: int32(code.Code_UNAVAILABLE),
		}
		switch checkResponse.RejectReason {
		case flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED:
			resp.HttpResponse = &flowcontrolhttpv1.CheckHTTPResponse_DeniedResponse{
				DeniedResponse: &flowcontrolhttpv1.DeniedHttpResponse{
					Status: http.StatusTooManyRequests,
				},
			}
		case flowcontrolv1.CheckResponse_REJECT_REASON_CONCURRENCY_LIMITED:
			resp.HttpResponse = &flowcontrolhttpv1.CheckHTTPResponse_DeniedResponse{
				DeniedResponse: &flowcontrolhttpv1.DeniedHttpResponse{
					Status: http.StatusServiceUnavailable,
				},
			}
		default:
			log.Bug().Stringer("reason", checkResponse.RejectReason).Msg("Unexpected reject reason")
		}
	default:
		return nil, grpc.Bug().Stringer("type", checkResponse.DecisionType).
			Msg("unexpected decision type")
	}

	return resp, nil
}

// RequestToInput - Converts a CheckHTTPRequest to an input map.
func RequestToInput(req *flowcontrolhttpv1.CheckHTTPRequest, logger logging.Logger) map[string]interface{} {
	var err error
	input := map[string]interface{}{}

	path := req.GetRequest().GetPath()
	body := req.GetRequest().GetBody()
	headers := req.GetRequest().GetHeaders()

	http := map[string]interface{}{}
	http["path"] = path
	http["body"] = body
	http["headers"] = headers
	http["host"] = req.GetRequest().GetHost()
	http["method"] = req.GetRequest().GetMethod()
	http["scheme"] = req.GetRequest().GetScheme()
	http["size"] = req.GetRequest().GetSize()
	http["protocol"] = req.GetRequest().GetProtocol()

	sourceSocketAddress := map[string]interface{}{}
	sourceSocketAddress["address"] = req.GetSource().GetAddress()
	sourceSocketAddress["port"] = req.GetSource().GetPort()

	destinationSocketAddress := map[string]interface{}{}
	destinationSocketAddress["address"] = req.GetDestination().GetAddress()
	destinationSocketAddress["port"] = req.GetDestination().GetPort()

	source := map[string]interface{}{}
	source["socketAddress"] = sourceSocketAddress

	destination := map[string]interface{}{}
	destination["socketAddress"] = destinationSocketAddress

	request := map[string]interface{}{}
	request["http"] = http

	attributes := map[string]interface{}{}
	attributes["request"] = request
	attributes["source"] = source
	attributes["destination"] = destination

	input["attributes"] = attributes

	parsedPath, parsedQuery, err := getParsedPathAndQuery(path)
	if err == nil {
		input["parsed_path"] = parsedPath
		input["parsed_query"] = parsedQuery
	}

	parsedBody, isBodyTruncated, err := getParsedBody(logger, headers, body)
	if err == nil {
		input["parsed_body"] = parsedBody
		input["truncated_body"] = isBodyTruncated
	}

	return input
}

func getParsedPathAndQuery(path string) ([]interface{}, map[string]interface{}, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}

	parsedPath := strings.Split(strings.TrimLeft(parsedURL.Path, "/"), "/")
	parsedPathInterface := make([]interface{}, len(parsedPath))
	for i, v := range parsedPath {
		parsedPathInterface[i] = v
	}

	parsedQueryInterface := make(map[string]interface{})
	for paramKey, paramValues := range parsedURL.Query() {
		queryValues := make([]interface{}, len(paramValues))
		for i, v := range paramValues {
			queryValues[i] = v
		}
		parsedQueryInterface[paramKey] = queryValues
	}

	return parsedPathInterface, parsedQueryInterface, nil
}

func getParsedBody(logger logging.Logger, headers map[string]string, body string) (interface{}, bool, error) {
	var data interface{}

	if val, ok := headers["content-type"]; ok {
		if strings.Contains(val, "application/json") {

			if body == "" {
				return nil, false, nil
			}

			if headerVal, ok := headers["content-length"]; ok {
				truncated, err := checkIfHTTPBodyTruncated(headerVal, int64(len(body)))
				if err != nil {
					return nil, false, err
				}
				if truncated {
					return nil, true, nil
				}
			}

			err := util.UnmarshalJSON([]byte(body), &data)
			if err != nil {
				return nil, false, err
			}
		} else if strings.Contains(val, "application/x-www-form-urlencoded") {
			var payload string
			switch {
			case body != "":
				payload = body
			default:
				return nil, false, nil
			}

			if headerVal, ok := headers["content-length"]; ok {
				truncated, err := checkIfHTTPBodyTruncated(headerVal, int64(len(payload)))
				if err != nil {
					return nil, false, err
				}
				if truncated {
					return nil, true, nil
				}
			}

			parsed, err := url.ParseQuery(payload)
			if err != nil {
				return nil, false, err
			}

			data = map[string][]string(parsed)
		} else if strings.Contains(val, "multipart/form-data") {
			var payload string
			switch {
			case body != "":
				payload = body
			default:
				return nil, false, nil
			}

			if headerVal, ok := headers["content-length"]; ok {
				truncated, err := checkIfHTTPBodyTruncated(headerVal, int64(len(payload)))
				if err != nil {
					return nil, false, err
				}
				if truncated {
					return nil, true, nil
				}
			}

			_, params, err := mime.ParseMediaType(headers["content-type"])
			if err != nil {
				return nil, false, err
			}

			boundary, ok := params["boundary"]
			if !ok {
				return nil, false, nil
			}

			values := map[string][]interface{}{}

			mr := multipart.NewReader(strings.NewReader(payload), boundary)
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, false, err
				}

				name := p.FormName()
				if name == "" {
					continue
				}

				value, err := io.ReadAll(p)
				if err != nil {
					return nil, false, err
				}

				switch {
				case strings.Contains(p.Header.Get("Content-Type"), "application/json"):
					var jsonValue interface{}
					if err := util.UnmarshalJSON(value, &jsonValue); err != nil {
						return nil, false, err
					}
					values[name] = append(values[name], jsonValue)
				default:
					values[name] = append(values[name], string(value))
				}
			}

			data = values
		} else {
			logger.Debug("content-type: %s parsing not supported", val)
		}
	} else {
		logger.Debug("no content-type header supplied, performing no body parsing")
	}

	return data, false, nil
}

func checkIfHTTPBodyTruncated(contentLength string, bodyLength int64) (bool, error) {
	cl, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return false, err
	}
	if cl != -1 && cl > bodyLength {
		return true, nil
	}
	return false, nil
}
