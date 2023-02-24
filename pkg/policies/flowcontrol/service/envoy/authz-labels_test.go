package envoy_test

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/envoy"
)

var _ = Describe("Flow labels", func() {
	It("should contain all the request labels with OTEL-compatible keys", func() {
		req := &authv3.AttributeContext_Request{
			Http: &authv3.AttributeContext_HttpRequest{
				Id:     "123",
				Method: "GET",
				Headers: map[string]string{
					"cache-control": "max-age=0",
				},
				Path:     "/foo/bar",
				Host:     "example.com",
				Scheme:   "https",
				Size:     9000,
				Protocol: "HTTP/2",
			},
		}

		flowLabels := AuthzRequestToFlowLabels(req)
		labels := flowLabels.ToPlainMap()
		Expect(labels).To(HaveKeyWithValue("http.method", "GET"))
		Expect(labels).To(HaveKeyWithValue("http.target", "/foo/bar"))
		Expect(labels).To(HaveKeyWithValue("http.host", "example.com"))
		Expect(labels).To(HaveKeyWithValue("http.scheme", "https"))
		Expect(labels).To(HaveKeyWithValue("http.request_content_length", "9000"))
		Expect(labels).To(HaveKeyWithValue("http.flavor", "2.0"))
		Expect(labels).To(HaveKeyWithValue("http.request.header.cache_control", "max-age=0"))
	})
})
