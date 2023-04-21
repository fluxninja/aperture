package checkhttp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	checkhttp_baggage "github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/checkhttp/baggage"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/servicegetter"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/util"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
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

	ctrlPt := req.GetControlPoint()
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

	input := RequestToInput(req)

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
	// checkResponse := h.fcHandler.CheckRequest(ctx, destinationSvcs, ctrlPt, flowLabels)
	checkResponse := h.fcHandler.CheckRequest(ctx,
		iface.RequestContext{
			Services:     destinationSvcs,
			ControlPoint: ctrlPt,
			FlowLabels:   flowLabels,
		},
	)
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
func RequestToInput(req *flowcontrolhttpv1.CheckHTTPRequest) ast.Value {
	return RequestToInputWithServices(req, nil, nil)
}

// RequestToInputWithServices - Converts a CheckHTTPRequest to an input map
// Additionally sets attributes.source.services and attributes.destination.services with discovered services.
func RequestToInputWithServices(req *flowcontrolhttpv1.CheckHTTPRequest, sourceSvcs, destinationSvcs []string) ast.Value {
	request := req.GetRequest()
	path := request.GetPath()
	body := request.GetBody()
	headers := request.GetHeaders()

	http := ast.NewObject()
	http.Insert(ast.StringTerm("path"), ast.StringTerm(path))
	http.Insert(ast.StringTerm("body"), ast.StringTerm(body))
	http.Insert(ast.StringTerm("host"), ast.StringTerm(request.GetHost()))
	http.Insert(ast.StringTerm("method"), ast.StringTerm(request.GetMethod()))
	http.Insert(ast.StringTerm("scheme"), ast.StringTerm(request.GetScheme()))
	http.Insert(ast.StringTerm("size"),
		ast.NumberTerm(json.Number(strconv.FormatInt(request.GetSize(), 10))))
	http.Insert(ast.StringTerm("protocol"), ast.StringTerm(request.GetProtocol()))

	headersObj := ast.NewObject()
	for key, val := range headers {
		headersObj.Insert(ast.StringTerm(key), ast.StringTerm(val))
	}
	http.Insert(ast.StringTerm("headers"), ast.NewTerm(headersObj))

	srcSocketAddress := ast.NewObject()
	srcSocketAddress.Insert(ast.StringTerm("address"), ast.StringTerm(req.GetSource().GetAddress()))
	srcSocketAddress.Insert(ast.StringTerm("port"),
		ast.NumberTerm(json.Number(strconv.FormatUint(uint64(req.GetSource().GetPort()), 10))))

	dstSocketAddress := ast.NewObject()
	dstSocketAddress.Insert(ast.StringTerm("address"), ast.StringTerm(req.GetDestination().GetAddress()))
	dstSocketAddress.Insert(ast.StringTerm("port"),
		ast.NumberTerm(json.Number(strconv.FormatUint(uint64(req.GetDestination().GetPort()), 10))))

	source := ast.NewObject()
	source.Insert(ast.StringTerm("socketAddress"), ast.NewTerm(srcSocketAddress))
	if sourceSvcs != nil {
		srcServicesArray := make([]*ast.Term, 0)
		for _, svc := range sourceSvcs {
			srcServicesArray = append(srcServicesArray, ast.StringTerm(svc))
		}
		source.Insert(ast.StringTerm("services"), ast.NewTerm(ast.NewArray(srcServicesArray...)))
	}

	destination := ast.NewObject()
	destination.Insert(ast.StringTerm("socketAddress"), ast.NewTerm(dstSocketAddress))
	if destinationSvcs != nil {
		dstServicesArray := make([]*ast.Term, 0)
		for _, svc := range destinationSvcs {
			dstServicesArray = append(dstServicesArray, ast.StringTerm(svc))
		}
		destination.Insert(ast.StringTerm("services"), ast.NewTerm(ast.NewArray(dstServicesArray...)))
	}

	requestMap := ast.NewObject()
	requestMap.Insert(ast.StringTerm("http"), ast.NewTerm(http))

	attributes := ast.NewObject()
	attributes.Insert(ast.StringTerm("request"), ast.NewTerm(requestMap))
	attributes.Insert(ast.StringTerm("source"), ast.NewTerm(source))
	attributes.Insert(ast.StringTerm("destination"), ast.NewTerm(destination))

	input := ast.NewObject()
	input.Insert(ast.StringTerm("attributes"), ast.NewTerm(attributes))

	parsedPath, parsedQuery, err := getParsedPathAndQuery(path)
	if err == nil {
		input.Insert(ast.StringTerm("parsed_path"), parsedPath)
		input.Insert(ast.StringTerm("parsed_query"), parsedQuery)
	}

	parsedBody, isBodyTruncated, err := getParsedBody(headers, body)
	if err == nil {
		input.Insert(ast.StringTerm("parsed_body"), parsedBody)
		input.Insert(ast.StringTerm("truncated_body"), ast.BooleanTerm(isBodyTruncated))
	}

	return input
}

func getParsedPathAndQuery(path string) (*ast.Term, *ast.Term, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return ast.NullTerm(), ast.NullTerm(), err
	}

	parsedPath := strings.Split(strings.TrimLeft(parsedURL.Path, "/"), "/")
	parsedPathSlice := make([]*ast.Term, 0)
	for _, v := range parsedPath {
		parsedPathSlice = append(parsedPathSlice, ast.StringTerm(v))
	}

	parsedQueryInterface := ast.NewObject()
	for paramKey, paramValues := range parsedURL.Query() {
		queryValues := make([]*ast.Term, 0)
		for _, v := range paramValues {
			queryValues = append(queryValues, ast.StringTerm(v))
		}
		parsedQueryInterface.Insert(ast.StringTerm(paramKey), ast.NewTerm(ast.NewArray(queryValues...)))
	}

	return ast.NewTerm(ast.NewArray(parsedPathSlice...)), ast.NewTerm(parsedQueryInterface), nil
}

func getParsedBody(headers map[string]string, body string) (*ast.Term, bool, error) {
	data := ast.NewObject()

	if val, ok := headers["content-type"]; ok {
		if strings.Contains(val, "application/json") {
			if body == "" {
				return ast.NullTerm(), false, nil
			}

			if headerVal, ok := headers["content-length"]; ok {
				truncated, err := checkIfHTTPBodyTruncated(headerVal, int64(len(body)))
				if err != nil {
					return ast.NullTerm(), false, err
				}
				if truncated {
					return ast.NullTerm(), true, nil
				}
			}

			astValue, err := ast.ValueFromReader(bytes.NewReader([]byte(body)))
			if err != nil {
				return ast.NullTerm(), false, err
			}
			return ast.NewTerm(astValue), false, nil
		} else if strings.Contains(val, "application/x-www-form-urlencoded") {
			var payload string
			switch {
			case body != "":
				payload = body
			default:
				return ast.NullTerm(), false, nil
			}

			if headerVal, ok := headers["content-length"]; ok {
				truncated, err := checkIfHTTPBodyTruncated(headerVal, int64(len(payload)))
				if err != nil {
					return ast.NullTerm(), false, err
				}
				if truncated {
					return ast.NullTerm(), true, nil
				}
			}

			parsed, err := url.ParseQuery(payload)
			if err != nil {
				return ast.NullTerm(), false, err
			}
			for key, valArray := range parsed {
				helperArr := make([]*ast.Term, 0)
				for _, val := range valArray {
					helperArr = append(helperArr, ast.StringTerm(val))
				}
				data.Insert(ast.StringTerm(key), ast.NewTerm(ast.NewArray(helperArr...)))
			}
		} else if strings.Contains(val, "multipart/form-data") {
			var payload string
			switch {
			case body != "":
				payload = body
			default:
				return ast.NullTerm(), false, nil
			}

			if headerVal, ok := headers["content-length"]; ok {
				truncated, err := checkIfHTTPBodyTruncated(headerVal, int64(len(payload)))
				if err != nil {
					return ast.NullTerm(), false, err
				}
				if truncated {
					return ast.NullTerm(), true, nil
				}
			}

			_, params, err := mime.ParseMediaType(headers["content-type"])
			if err != nil {
				return ast.NullTerm(), false, err
			}

			boundary, ok := params["boundary"]
			if !ok {
				return ast.NullTerm(), false, nil
			}

			values := ast.NewObject()

			mr := multipart.NewReader(strings.NewReader(payload), boundary)
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					return ast.NullTerm(), false, err
				}

				name := p.FormName()
				if name == "" {
					continue
				}

				value, err := io.ReadAll(p)
				if err != nil {
					return ast.NullTerm(), false, err
				}

				switch {
				case strings.Contains(p.Header.Get("Content-Type"), "application/json"):
					var jsonValue interface{}
					if err := util.UnmarshalJSON(value, &jsonValue); err != nil {
						return ast.NullTerm(), false, err
					}
					jsonData, err := ast.InterfaceToValue(jsonValue)
					if err != nil {
						return ast.NullTerm(), false, err
					}
					values.Insert(ast.StringTerm(name),
						ast.NewTerm(ast.NewArray(ast.NewTerm(jsonData))))
				default:
					values.Insert(ast.StringTerm(name),
						ast.NewTerm(ast.NewArray((ast.StringTerm(string(value))))))
				}
			}

			data = values
		} else {
			log.Debug().Msgf("rego content-type: %s parsing not supported", val)
		}
	} else {
		log.Debug().Msg("rego no content-type header supplied, performing no body parsing")
	}

	return ast.NewTerm(data), false, nil
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
