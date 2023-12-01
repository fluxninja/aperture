package checkhttp_test

import (
	flowcontrolhttpv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/checkhttp"
)

var _ = Describe("Flow labels", func() {
	It("should contain all the request labels with OTel-compatible keys", func() {
		req := &flowcontrolhttpv1.CheckHTTPRequest_HttpRequest{
			Method: "GET",
			Headers: map[string]string{
				"cache-control": "max-age=0",
			},
			Path:     "/foo/bar",
			Host:     "example.com",
			Scheme:   "https",
			Size:     9000,
			Protocol: "HTTP/2",
		}

		flowLabels := CheckHTTPRequestToFlowLabels(req)
		labels := flowLabels.Copy()
		Expect(labels).To(HaveKeyWithValue("http.method", "GET"))
		Expect(labels).To(HaveKeyWithValue("http.target", "/foo/bar"))
		Expect(labels).To(HaveKeyWithValue("http.host", "example.com"))
		Expect(labels).To(HaveKeyWithValue("http.scheme", "https"))
		Expect(labels).To(HaveKeyWithValue("http.request_content_length", "9000"))
		Expect(labels).To(HaveKeyWithValue("http.flavor", "2.0"))
		Expect(labels).To(HaveKeyWithValue("http.request.header.cache_control", "max-age=0"))
	})
})
