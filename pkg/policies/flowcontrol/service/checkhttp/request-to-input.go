package checkhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/util"

	flowcontrolhttpv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
)

// RequestToInput - Converts a CheckHTTPRequest to an input map.
//
// The given CheckHTTPRequest should not be modified while returned Input is in
// use (because the Input will be created lazily).
func RequestToInput(req *flowcontrolhttpv1.CheckHTTPRequest) classifier.Input {
	return RequestToInputWithServices(req, nil, nil)
}

var (
	pathStringTerm             = ast.StringTerm("path")
	bodyStringTerm             = ast.StringTerm("body")
	hostStringTerm             = ast.StringTerm("host")
	methodStringTerm           = ast.StringTerm("method")
	schemeStringTerm           = ast.StringTerm("scheme")
	sizeStringTerm             = ast.StringTerm("size")
	protocolStringTerm         = ast.StringTerm("protocol")
	headersStringTerm          = ast.StringTerm("headers")
	socketAddressStringTerm    = ast.StringTerm("socketAddress")
	sourceFQDNsStringTerm      = ast.StringTerm("source_fqdns")
	destinationFQDNsStringTerm = ast.StringTerm("destination_fqdns")
	requestStringTerm          = ast.StringTerm("request")
	sourceStringTerm           = ast.StringTerm("source")
	destinationStringTerm      = ast.StringTerm("destination")
	attributesStringTerm       = ast.StringTerm("attributes")
	parsedPathStringTerm       = ast.StringTerm("parsed_path")
	parsedQueryStringTerm      = ast.StringTerm("parsed_query")
	parsedBodyStringTerm       = ast.StringTerm("parsed_body")
	truncatedBodyStringTerm    = ast.StringTerm("truncated_body")
	addressStringTerm          = ast.StringTerm("address")
	portStringTerm             = ast.StringTerm("port")
	httpStringTerm             = ast.StringTerm("http")
)

// RequestToInputWithServices - Converts a CheckHTTPRequest to an input map
// Additionally sets attributes.source.services and attributes.destination.services with discovered services.
//
// Arguments should not be modified while returned Input is in use (because the
// Input will be created lazily).
func RequestToInputWithServices(
	req *flowcontrolhttpv1.CheckHTTPRequest,
	sourceSvcs,
	destinationSvcs []string,
) classifier.Input {
	return newLazyInput(func() ast.Value {
		return requestToInputWithServices(req, sourceSvcs, destinationSvcs)
	})
}

func requestToInputWithServices(
	req *flowcontrolhttpv1.CheckHTTPRequest,
	sourceSvcs,
	destinationSvcs []string,
) ast.Value {
	request := req.GetRequest()
	path := request.GetPath()
	body := request.GetBody()
	headers := request.GetHeaders()

	headersKV := make([][2]*ast.Term, 0, len(headers))
	for key, val := range headers {
		headersKV = append(headersKV, [2]*(ast.Term){ast.StringTerm(key), ast.StringTerm(val)})
	}
	headersObj := ast.NewObject(headersKV...)
	http := ast.NewObject(
		[2]*ast.Term{pathStringTerm, ast.StringTerm(path)},
		[2]*ast.Term{bodyStringTerm, ast.StringTerm(body)},
		[2]*ast.Term{hostStringTerm, ast.StringTerm(request.GetHost())},
		[2]*ast.Term{methodStringTerm, ast.StringTerm(request.GetMethod())},
		[2]*ast.Term{schemeStringTerm, ast.StringTerm(request.GetScheme())},
		[2]*ast.Term{sizeStringTerm, ast.NumberTerm(json.Number(strconv.FormatInt(request.GetSize(), 10)))},
		[2]*ast.Term{protocolStringTerm, ast.StringTerm(request.GetProtocol())},
		[2]*ast.Term{headersStringTerm, ast.NewTerm(headersObj)},
	)

	srcSocketAddress := ast.NewObject(
		[2]*ast.Term{addressStringTerm, ast.StringTerm(req.GetSource().GetAddress())},
		[2]*ast.Term{portStringTerm, ast.NumberTerm(json.Number(strconv.FormatUint(uint64(req.GetSource().GetPort()), 10)))},
	)

	dstSocketAddress := ast.NewObject(
		[2]*ast.Term{addressStringTerm, ast.StringTerm(req.GetDestination().GetAddress())},
		[2]*ast.Term{portStringTerm, ast.NumberTerm(json.Number(strconv.FormatUint(uint64(req.GetDestination().GetPort()), 10)))},
	)

	// sourceKV array is used to call ast.NewObject only once (to avoid reallocations in Insert).
	sourceKV := make([][2]*ast.Term, 0, 2)
	sourceKV = append(sourceKV, [2]*ast.Term{socketAddressStringTerm, ast.NewTerm(srcSocketAddress)})
	if sourceSvcs != nil {
		srcServicesArray := make([]*ast.Term, 0, len(sourceSvcs))
		for _, svc := range sourceSvcs {
			srcServicesArray = append(srcServicesArray, ast.StringTerm(svc))
		}
		sourceKV = append(sourceKV, [2]*ast.Term{sourceFQDNsStringTerm, ast.NewTerm(ast.NewArray(srcServicesArray...))})
	}
	source := ast.NewObject(sourceKV...)

	// see comment on sourceKV
	destinationKV := make([][2]*ast.Term, 0, 2)
	destinationKV = append(destinationKV, [2]*ast.Term{socketAddressStringTerm, ast.NewTerm(dstSocketAddress)})
	if destinationSvcs != nil {
		dstServicesArray := make([]*ast.Term, 0, len(destinationSvcs))
		for _, svc := range destinationSvcs {
			dstServicesArray = append(dstServicesArray, ast.StringTerm(svc))
		}
		destinationKV = append(destinationKV, [2]*ast.Term{destinationFQDNsStringTerm, ast.NewTerm(ast.NewArray(dstServicesArray...))})
	}
	destination := ast.NewObject(destinationKV...)

	requestMap := ast.NewObject(
		[2]*ast.Term{httpStringTerm, ast.NewTerm(http)},
	)

	attributes := ast.NewObject(
		[2]*ast.Term{requestStringTerm, ast.NewTerm(requestMap)},
		[2]*ast.Term{sourceStringTerm, ast.NewTerm(source)},
		[2]*ast.Term{destinationStringTerm, ast.NewTerm(destination)},
	)

	// see comment on sourceKV
	inputKV := make([][2]*ast.Term, 0, 5)
	inputKV = append(inputKV, [2]*ast.Term{attributesStringTerm, ast.NewTerm(attributes)})

	parsedPath, parsedQuery, err := getParsedPathAndQuery(path)
	if err == nil {
		inputKV = append(inputKV, [2]*ast.Term{parsedPathStringTerm, parsedPath})
		inputKV = append(inputKV, [2]*ast.Term{parsedQueryStringTerm, parsedQuery})
	}

	parsedBody, isBodyTruncated, err := getParsedBody(headers, body)
	if err == nil {
		inputKV = append(inputKV, [2]*ast.Term{parsedBodyStringTerm, parsedBody})
		inputKV = append(inputKV, [2]*ast.Term{truncatedBodyStringTerm, ast.BooleanTerm(isBodyTruncated)})
	}
	input := ast.NewObject(inputKV...)

	return input
}

func getParsedPathAndQuery(path string) (*ast.Term, *ast.Term, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return ast.NullTerm(), ast.NullTerm(), err
	}

	parsedPath := strings.Split(strings.TrimLeft(parsedURL.Path, "/"), "/")
	parsedPathSlice := make([]*ast.Term, 0, len(parsedPath))
	for _, v := range parsedPath {
		parsedPathSlice = append(parsedPathSlice, ast.StringTerm(v))
	}

	parsedQueryInterface := ast.NewObject()
	for paramKey, paramValues := range parsedURL.Query() {
		queryValues := make([]*ast.Term, 0, len(paramValues))
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

type lazyInput struct {
	mkValue   func() ast.Value
	valueOnce sync.Once
	value     ast.Value
	ifaceOnce sync.Once
	iface     interface{}
}

func newLazyInput(mkValueFunc func() ast.Value) classifier.Input {
	return &lazyInput{mkValue: mkValueFunc}
}

// Value implements classifier.Input.
func (i *lazyInput) Value() ast.Value {
	i.valueOnce.Do(func() {
		i.value = i.mkValue()
	})
	return i.value
}

// Interface implements classifier.Input.
func (i *lazyInput) Interface() interface{} {
	i.ifaceOnce.Do(func() {
		var err error
		i.iface, err = ast.ValueToInterface(i.Value(), emptyResolver{})
		if err != nil {
			// This should never happen as we're using only "simple" types in our
			// Value (objects, arrays, strings, bools, ints), and resolver we
			// use is infallible anyway. Log just in case.
			log.Bug().Msgf("failed to convert value to interface: %v", err)
			i.iface = map[string]interface{}{"error": "failed to convert value"}
		}
	})
	return i.iface
}

type emptyResolver struct{}

// Resolve implements ast.ValueResolver interface by returning empty object for
// any unknown ref.
//
// Note: We assume this will never be called.
func (emptyResolver) Resolve(ref ast.Ref) (interface{}, error) {
	return make(map[string]interface{}), nil
}
