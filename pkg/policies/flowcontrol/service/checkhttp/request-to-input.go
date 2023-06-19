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

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/util"

	flowcontrolhttpv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// RequestToInput - Converts a CheckHTTPRequest to an input map.
func RequestToInput(req *flowcontrolhttpv1.CheckHTTPRequest) ast.Value {
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
)

// RequestToInputWithServices - Converts a CheckHTTPRequest to an input map
// Additionally sets attributes.source.services and attributes.destination.services with discovered services.
func RequestToInputWithServices(req *flowcontrolhttpv1.CheckHTTPRequest, sourceSvcs, destinationSvcs []string) ast.Value {
	request := req.GetRequest()
	path := request.GetPath()
	body := request.GetBody()
	headers := request.GetHeaders()

	http := ast.NewObject()
	http.Insert(pathStringTerm, ast.StringTerm(path))
	http.Insert(bodyStringTerm, ast.StringTerm(body))
	http.Insert(hostStringTerm, ast.StringTerm(request.GetHost()))
	http.Insert(methodStringTerm, ast.StringTerm(request.GetMethod()))
	http.Insert(schemeStringTerm, ast.StringTerm(request.GetScheme()))
	http.Insert(sizeStringTerm, ast.NumberTerm(json.Number(strconv.FormatInt(request.GetSize(), 10))))
	http.Insert(protocolStringTerm, ast.StringTerm(request.GetProtocol()))

	headersObj := ast.NewObject()
	for key, val := range headers {
		headersObj.Insert(ast.StringTerm(key), ast.StringTerm(val))
	}
	http.Insert(headersStringTerm, ast.NewTerm(headersObj))

	srcSocketAddress := ast.NewObject()
	srcSocketAddress.Insert(addressStringTerm, ast.StringTerm(req.GetSource().GetAddress()))
	srcSocketAddress.Insert(portStringTerm, ast.NumberTerm(json.Number(strconv.FormatUint(uint64(req.GetSource().GetPort()), 10))))

	dstSocketAddress := ast.NewObject()
	dstSocketAddress.Insert(addressStringTerm, ast.StringTerm(req.GetDestination().GetAddress()))
	dstSocketAddress.Insert(portStringTerm, ast.NumberTerm(json.Number(strconv.FormatUint(uint64(req.GetDestination().GetPort()), 10))))

	source := ast.NewObject()
	source.Insert(socketAddressStringTerm, ast.NewTerm(srcSocketAddress))
	if sourceSvcs != nil {
		srcServicesArray := make([]*ast.Term, 0)
		for _, svc := range sourceSvcs {
			srcServicesArray = append(srcServicesArray, ast.StringTerm(svc))
		}
		source.Insert(sourceFQDNsStringTerm, ast.NewTerm(ast.NewArray(srcServicesArray...)))
	}

	destination := ast.NewObject()
	destination.Insert(socketAddressStringTerm, ast.NewTerm(dstSocketAddress))
	if destinationSvcs != nil {
		dstServicesArray := make([]*ast.Term, 0)
		for _, svc := range destinationSvcs {
			dstServicesArray = append(dstServicesArray, ast.StringTerm(svc))
		}
		destination.Insert(destinationFQDNsStringTerm, ast.NewTerm(ast.NewArray(dstServicesArray...)))
	}

	requestMap := ast.NewObject()
	requestMap.Insert(ast.StringTerm("http"), ast.NewTerm(http))

	attributes := ast.NewObject()
	attributes.Insert(requestStringTerm, ast.NewTerm(requestMap))
	attributes.Insert(sourceStringTerm, ast.NewTerm(source))
	attributes.Insert(destinationStringTerm, ast.NewTerm(destination))

	input := ast.NewObject()
	input.Insert(attributesStringTerm, ast.NewTerm(attributes))

	parsedPath, parsedQuery, err := getParsedPathAndQuery(path)
	if err == nil {
		input.Insert(parsedPathStringTerm, parsedPath)
		input.Insert(parsedQueryStringTerm, parsedQuery)
	}

	parsedBody, isBodyTruncated, err := getParsedBody(headers, body)
	if err == nil {
		input.Insert(parsedBodyStringTerm, parsedBody)
		input.Insert(truncatedBodyStringTerm, ast.BooleanTerm(isBodyTruncated))
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
