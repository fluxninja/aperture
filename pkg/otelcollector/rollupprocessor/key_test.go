package rollupprocessor

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var _ = Describe("Rollup key", func() {
	It("ignores order", func() {
		m1 := pcommon.NewMap()
		m1.PutStr("foo", "fooval")
		m1.PutStr("bar", "barval")
		k1 := key(m1, nil)
		m2 := pcommon.NewMap()
		m2.PutStr("bar", "barval")
		m2.PutStr("foo", "fooval")
		k2 := key(m2, nil)
		Expect(k1).To(Equal(k2))
	})

	It("differs when attrs differ", func() {
		m1 := pcommon.NewMap()
		m1.PutStr("foo", "fooval")
		m1.PutStr("bar", "barval")
		k1 := key(m1, nil)
		m2 := pcommon.NewMap()
		m2.PutStr("foo", "fooval")
		m2.PutStr("bar", "barval2")
		k2 := key(m2, nil)
		Expect(k1).NotTo(Equal(k2))
	})

	It("differs when slice attrs differ", func() {
		m1 := pcommon.NewMap()
		Expect(m1.FromRaw(map[string]any{
			"foo":  "fooval",
			"vals": []any{"x, y"},
		})).To(Succeed())
		k1 := key(m1, nil)
		m2 := pcommon.NewMap()
		Expect(m2.FromRaw(map[string]any{
			"foo":  "fooval",
			"vals": []any{"x, z"},
		})).To(Succeed())
		k2 := key(m2, nil)
		Expect(k1).NotTo(Equal(k2))
	})

	It("ignores ignored attrs", func() {
		ignored := map[string]struct{}{"bar": {}}
		m1 := pcommon.NewMap()
		m1.PutStr("foo", "fooval")
		m1.PutStr("bar", "barval")
		k1 := key(m1, ignored)
		m2 := pcommon.NewMap()
		m2.PutStr("foo", "fooval")
		m2.PutStr("bar", "barval2")
		k2 := key(m2, ignored)
		Expect(k1).To(Equal(k2))
	})
})

// Run with: go test -bench=. -run='^$' -benchmem ./...
func BenchmarkKey(b *testing.B) {
	var attributes [8]pcommon.Map
	// Try to achieve different key orderings for _less_ variance in tests and
	// more strain on branch predictor.
	for i := 0; i < 8; i++ {
		attributes[i] = exampleAttributes()
	}
	for i := 0; i < b.N; i++ {
		_ = key(attributes[i%8], nil)
	}
}

func exampleAttributes() pcommon.Map {
	var attrsRaw map[string]any
	if err := json.Unmarshal([]byte(exampleAttributesJSON), &attrsRaw); err != nil {
		panic(err)
	}
	attrs := pcommon.NewMap()
	if err := attrs.FromRaw(attrsRaw); err != nil {
		panic(err)
	}
	return attrs
}

// Note: these are attributes after running processors (just as key receives it).
// Scraped from default playground scenario on 2023-06-21.
const exampleAttributesJSON = `
{
  "aperture.classifier_errors": [],
  "aperture.classifiers": [
    "policy_name:service-protection,classifier_index:0"
  ],
  "aperture.control_point": "awesomeFeature",
  "aperture.control_point_type": "http",
  "aperture.decision_type": "DECISION_TYPE_ACCEPTED",
  "aperture.destination_fqdns": "service1-demo-app.demoapp.svc.cluster.local",
  "aperture.dropping_load_schedulers": [],
  "aperture.dropping_rate_limiters": [],
  "aperture.dropping_workloads": [],
  "aperture.flow.status": "OK",
  "aperture.flow_label_keys": [
    "aperture.destination_fqdns",
    "aperture.source_fqdns",
    "http.flavor",
    "http.host",
    "http.method",
    "http.request.header.content_length",
    "http.request.header.content_type",
    "http.request.header.cookie",
    "http.request.header.host",
    "http.request.header.user_agent",
    "http.request.header.user_id",
    "http.request.header.user_type",
    "http.request_content_length",
    "http.scheme",
    "http.target",
    "user_type"
  ],
  "aperture.flux_meters": [],
  "aperture.load_schedulers": [
    "policy_name:service-protection,component_id:root.0.10,policy_hash:DVWnkjji2OWF8Vuxtm+XG8fLzA4PMkX2BWJXiE3wpMs="
  ],
  "aperture.rate_limiters": [],
  "aperture.reject_reason": "REJECT_REASON_NONE",
  "aperture.services": [
    "service1-demo-app.demoapp.svc.cluster.local"
  ],
  "aperture.source": "sdk",
  "aperture.source_fqdns": "UNKNOWN",
  "aperture.workloads": [
    "policy_name:service-protection,component_id:root.0.10,workload_index:0,policy_hash:DVWnkjji2OWF8Vuxtm+XG8fLzA4PMkX2BWJXiE3wpMs="
  ],
  "aperture_processing_duration_ms": 0,
  "flow_duration_ms": 70,
  "http.flavor": "1.1",
  "http.host": "10.244.2.99",
  "http.method": "POST",
  "http.request.header.content_length": "225",
  "http.request.header.content_type": "application/json",
  "http.request.header.cookie": "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
  "http.request.header.host": "service1-demo-app.demoapp.svc.cluster.local",
  "http.request.header.user_agent": "k6/0.45.0 (https://k6.io/)",
  "http.request.header.user_id": "132",
  "http.request.header.user_type": "guest",
  "http.request_content_length": "225",
  "http.scheme": "http",
  "http.target": "/request",
  "response_received": "true",
  "user_type": "guest",
  "workload_duration_ms": 60
}
`
