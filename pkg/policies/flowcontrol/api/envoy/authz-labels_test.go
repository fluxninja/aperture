package envoy_test

import (
	. "github.com/fluxninja/aperture/pkg/policies/flowcontrol/api/envoy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

var _ = Describe("Flow labels", func() {
	It("should contain all the request labels with OTEL-compatible keys", func() {
		req := &ext_authz.AttributeContext_Request{
			Http: &ext_authz.AttributeContext_HttpRequest{
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
