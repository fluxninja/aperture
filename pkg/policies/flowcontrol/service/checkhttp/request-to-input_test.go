package checkhttp_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	flowcontrolhttpv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/checkhttp"
	"github.com/open-policy-agent/opa/ast"
)

var _ = Describe("RequestToInput", func() {

	var defaultRequest *flowcontrolhttpv1.CheckHTTPRequest
	var commonInput map[string]interface{}

	BeforeEach(func() {
		defaultRequest = &flowcontrolhttpv1.CheckHTTPRequest{
			ControlPoint: "ingress",
			Source: &flowcontrolhttpv1.SocketAddress{
				Address: "test",
				Port:    1234,
			},
			Destination: &flowcontrolhttpv1.SocketAddress{
				Address: "test2",
				Port:    1234,
			},
			Request: &flowcontrolhttpv1.CheckHTTPRequest_HttpRequest{
				Method: "GET",
				Path:   "/test/test2?param1=val",
				Headers: map[string]string{
					"content-type": "application/json",
				},
				Host:     "localhost",
				Body:     "{ \"myjson\": \"value\" }",
				Protocol: "HTTP/1.1",
				Scheme:   "http",
				Size:     21,
			},
		}

		commonInput = map[string]interface{}{
			"attributes": map[string]interface{}{
				"destination": map[string]interface{}{
					"socketAddress": map[string]interface{}{
						"address": "test2",
						"port":    json.Number("1234"),
					},
				},
				"source": map[string]interface{}{
					"socketAddress": map[string]interface{}{
						"address": "test",
						"port":    json.Number("1234"),
					},
				},
				"request": map[string]interface{}{
					"http": map[string]interface{}{
						"body": "{ \"myjson\": \"value\" }",
						"headers": map[string]interface{}{
							"content-type": "application/json",
						},
						"host":     "localhost",
						"method":   "GET",
						"path":     "/test/test2?param1=val",
						"protocol": "HTTP/1.1",
						"scheme":   "http",
						"size":     json.Number("21"),
					},
				},
			},
			"parsed_body": map[string]interface{}{
				"myjson": "value",
			},
			"parsed_query": map[string]interface{}{
				"param1": []interface{}{"val"},
			},
			"parsed_path": []interface{}{
				"test", "test2",
			},
			"truncated_body": false,
		}
	})

	It("Can process empty HTTP request", func() {
		value := checkhttp.RequestToInput(&flowcontrolhttpv1.CheckHTTPRequest{})
		iface, err := ast.ValueToInterface(value, valueResolver{})
		Expect(iface).To(Equal(emptyInput()))
		Expect(err).NotTo(HaveOccurred())
	})

	It("Can process valid HTTP request", func() {
		value := checkhttp.RequestToInput(defaultRequest)
		iface, err := ast.ValueToInterface(value, valueResolver{})
		Expect(iface).To(Equal(commonInput))
		Expect(err).NotTo(HaveOccurred())
	})

	It("Can process valid request with truncated body", func() {
		req := defaultRequest
		req.Request.Headers["content-length"] = "64"
		value := checkhttp.RequestToInput(req)
		iface, err := ast.ValueToInterface(value, valueResolver{})
		Expect(iface).To(Equal(truncatedInput(commonInput)))
		Expect(err).NotTo(HaveOccurred())
	})

	It("Can process valid HTTP request with content type url encoded", func() {
		req := defaultRequest
		req.Request.Headers["content-type"] = "application/x-www-form-urlencoded"
		req.Request.Headers["content-length"] = "12"
		req.Request.Body = "myjson=value"
		value := checkhttp.RequestToInput(req)
		iface, err := ast.ValueToInterface(value, valueResolver{})
		Expect(iface).To(Equal(urlEncoded(commonInput)))
		Expect(err).NotTo(HaveOccurred())
	})

	It("Can process valid HTTP request with truncated url encoded body", func() {
		req := defaultRequest
		req.Request.Headers["content-type"] = "application/x-www-form-urlencoded"
		req.Request.Headers["content-length"] = "64"
		req.Request.Body = "myjson=value"
		value := checkhttp.RequestToInput(req)
		iface, err := ast.ValueToInterface(value, valueResolver{})
		Expect(iface).To(Equal(urlEncodedTruncated(commonInput)))
		Expect(err).NotTo(HaveOccurred())
	})
})

func emptyInput() map[string]interface{} {
	input := map[string]interface{}{
		"attributes": map[string]interface{}{
			"destination": map[string]interface{}{
				"socketAddress": map[string]interface{}{
					"address": "",
					"port":    json.Number("0"),
				},
			},
			"source": map[string]interface{}{
				"socketAddress": map[string]interface{}{
					"address": "",
					"port":    json.Number("0"),
				},
			},
			"request": map[string]interface{}{
				"http": map[string]interface{}{
					"body":     "",
					"headers":  map[string]interface{}{},
					"host":     "",
					"method":   "",
					"path":     "",
					"protocol": "",
					"scheme":   "",
					"size":     json.Number("0"),
				},
			},
		},
		"parsed_body":    map[string]interface{}{},
		"parsed_query":   map[string]interface{}{},
		"parsed_path":    []interface{}{""},
		"truncated_body": false,
	}

	return input
}

func truncatedInput(commonInput map[string]interface{}) map[string]interface{} {
	input := commonInput
	input["attributes"].(map[string]interface{})["request"].(map[string]interface{})["http"].(map[string]interface{})["headers"] = map[string]interface{}{
		"content-type":   "application/json",
		"content-length": "64",
	}
	input["truncated_body"] = true
	input["parsed_body"] = nil

	return input
}

func urlEncoded(commonInput map[string]interface{}) map[string]interface{} {
	input := commonInput
	httpMap := input["attributes"].(map[string]interface{})["request"].(map[string]interface{})["http"].(map[string]interface{})
	httpMap["headers"] = map[string]interface{}{
		"content-type":   "application/x-www-form-urlencoded",
		"content-length": "12",
	}
	httpMap["body"] = "myjson=value"
	input["parsed_body"] = map[string]interface{}{
		"myjson": []interface{}{"value"},
	}

	return input
}

func urlEncodedTruncated(commonInput map[string]interface{}) map[string]interface{} {
	input := commonInput
	httpMap := input["attributes"].(map[string]interface{})["request"].(map[string]interface{})["http"].(map[string]interface{})
	httpMap["headers"] = map[string]interface{}{
		"content-type":   "application/x-www-form-urlencoded",
		"content-length": "64",
	}
	httpMap["body"] = "myjson=value"
	input["truncated_body"] = true
	input["parsed_body"] = nil

	return input
}

type valueResolver struct{}

// Resolve implements ast.ValueResolver interface.
func (valueResolver) Resolve(ref ast.Ref) (interface{}, error) {
	return make(map[string]interface{}), nil
}
