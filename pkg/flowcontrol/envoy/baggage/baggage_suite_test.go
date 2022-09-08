package baggage_test

import (
	"testing"

	envoy_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/pkg/flowcontrol/envoy/baggage"
	class "github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier/compiler"
	"github.com/fluxninja/aperture/pkg/utils"
)

func TestBaggage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Baggage Suite")
}

func fl(s string) class.FlowLabelValue {
	return class.FlowLabelValue{
		Value: s,
		Flags: compiler.LabelFlags{Propagate: true},
	}
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Prefixed propagator", func() {
	var propagator baggage.Propagator = baggage.Prefixed{Prefix: "myprefix-"}

	It("extracts flow labels from baggage", func() {
		Expect(propagator.Extract(map[string]string{
			"myprefix-foo": "bar",
			"myprefix-baz": "quux",
			"content-type": "application/json",
		})).To(Equal(class.FlowLabels{
			"foo": fl("bar"),
			"baz": fl("quux"),
		}))
	})

	It("handles urlencoded values", func() {
		Expect(propagator.Extract(map[string]string{
			"myprefix-foo": "%20",
		})).To(Equal(class.FlowLabels{
			"foo": fl(" "),
		}))
	})

	It("creates injecting instructions for envoy", func() {
		newHeaders, err := propagator.Inject(class.FlowLabels{
			"foo": fl("bar"),
			"baz": fl("quux"),
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(newHeaders).To(HaveLen(2))
		Expect(newHeaders).To(ContainElement(&envoy_core.HeaderValueOption{
			Header: &envoy_core.HeaderValue{Key: "myprefix-foo", Value: "bar"},
			Append: wrapperspb.Bool(false),
		}))
		Expect(newHeaders).To(ContainElement(&envoy_core.HeaderValueOption{
			Header: &envoy_core.HeaderValue{Key: "myprefix-baz", Value: "quux"},
			Append: wrapperspb.Bool(false),
		}))
	})
})

var _ = Describe("W3 Baggage propagator", func() {
	var propagator baggage.Propagator = baggage.W3Baggage{}

	Context("when there's no baggage header", func() {
		It("reads no flow labels", func() {
			Expect(propagator.Extract(map[string]string{
				"content-type": "application/json",
			})).To(Equal(class.FlowLabels{}))
		})

		It("creates injecting instructions for envoy", func() {
			newHeaders, err := propagator.Inject(class.FlowLabels{
				"foo": fl("bar"),
			}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(newHeaders).To(Equal([]*envoy_core.HeaderValueOption{{
				Header: &envoy_core.HeaderValue{Key: "baggage", Value: "foo=bar"},
				Append: wrapperspb.Bool(false),
			}}))
		})
	})

	Context("when baggage header exists", func() {
		It("extracts flow labels from baggage", func() {
			Expect(propagator.Extract(map[string]string{
				"baggage":      "foo=bar,baz=quux",
				"content-type": "application/json",
			})).To(Equal(class.FlowLabels{
				"foo": fl("bar"),
				"baz": fl("quux"),
			}))
		})

		It("creates injecting instructions for envoy", func() {
			newHeaders, err := propagator.Inject(class.FlowLabels{
				"hello": fl("world"),
			}, baggage.Headers(map[string]string{
				"baggage":      "foo=bar,baz=quux",
				"content-type": "application/json",
			}))
			Expect(err).NotTo(HaveOccurred())
			Expect(newHeaders).To(Equal([]*envoy_core.HeaderValueOption{{
				Header: &envoy_core.HeaderValue{Key: "baggage", Value: "hello=world"},
				Append: wrapperspb.Bool(true),
			}}))
		})
	})

	Context("when flow label has 'hidden' flag", func() {
		It("extracts it correctly", func() {
			Expect(propagator.Extract(map[string]string{
				"baggage":      "foo=bar;hidden",
				"content-type": "application/json",
			})).To(Equal(class.FlowLabels{
				"foo": class.FlowLabelValue{
					Value: "bar",
					Flags: compiler.LabelFlags{Hidden: true, Propagate: true},
				},
			}))
		})

		It("injects it correctly", func() {
			newHeaders, err := propagator.Inject(class.FlowLabels{
				"foo": class.FlowLabelValue{
					Value: "bar",
					Flags: compiler.LabelFlags{Hidden: true, Propagate: true},
				},
			}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(newHeaders).To(Equal([]*envoy_core.HeaderValueOption{{
				Header: &envoy_core.HeaderValue{Key: "baggage", Value: "foo=bar;hidden"},
				Append: wrapperspb.Bool(false),
			}}))
		})
	})

	It("ignores member properties", func() {
		Expect(propagator.Extract(map[string]string{
			"baggage": "foo=bar;props",
		})).To(Equal(class.FlowLabels{
			"foo": fl("bar"),
		}))
	})

	It("handles urlencoded values", func() {
		Expect(propagator.Extract(map[string]string{
			// TODO make sure this is correct FLUX-1290
			"baggage": "foo=%2520",
		})).To(Equal(class.FlowLabels{
			"foo": fl(" "),
		}))
	})
})
