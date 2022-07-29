package enrichmentprocessor

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Enrichment Processor - Traces", func() {
	It("Enriches egress traces attributes with data from entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID:         entitycache.EntityID{},
			IPAddress:  "192.0.2.0",
			AgentGroup: "defaultAG",
			Namespace:  "nspc1",
			Services:   []string{"svc1", "svc2"},
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]string{
			"control_point":    "egress",
			"net.host.address": "192.0.2.0:80",
			"net.peer.address": "192.0.2.1:80",
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]string{
			"control_point": "egress",
			"labeled":       "false",
			"agent_group":   "defaultAG",
			"namespace":     "nspc1",
			"services":      "svc1,svc2",
		}))
	})

	It("Does not panic when egress metrics net.host.address attribute empty", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID:        entitycache.EntityID{},
			IPAddress: "192.0.2.0",
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]string{
			"control_point":    "egress",
			"net.host.address": "",
			"net.peer.address": "192.0.2.1:80",
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]string{
			"control_point":    "egress",
			"labeled":          "false",
			"net.host.address": "",
			"net.peer.address": "192.0.2.1:80",
		}))
	})

	It("Enriches ingress traces attributes with data from entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID:         entitycache.EntityID{},
			IPAddress:  "192.0.2.0",
			AgentGroup: "defaultAG",
			Namespace:  "nspc1",
			Services:   []string{"svc1", "svc2"},
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]string{
			"control_point": "ingress",
			"net.host.ip":   "192.0.2.0",
			"net.peer.ip":   "192.0.2.1",
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]string{
			"control_point": "ingress",
			"labeled":       "false",
			"agent_group":   "defaultAG",
			"namespace":     "nspc1",
			"services":      "svc1,svc2",
		}))
	})

	It("Does not enrich when there are no labels in entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]string{
			"control_point":    "egress",
			"net.host.address": "192.0.2.0:80",
			"net.peer.address": "192.0.2.1:80",
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]string{
			"control_point": "egress",
			"labeled":       "false",
		}))
	})

	It("Unpacks fn.flow properly", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]string{
			"control_point": "egress",
			"fn.flow":       `{"foo": "bar", "fizz": "buzz"}`,
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]string{
			"control_point": "egress",
			"fn.flow.foo":   "bar",
			"fn.flow.fizz":  "buzz",
			"labeled":       "true",
		}))
	})

	It("Ignores empty fn.flow", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]string{
			"control_point": "egress",
			"fn.flow":       ``,
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]string{
			"control_point": "egress",
			"labeled":       "false",
		}))
	})

	It("Ignores minus as fn.flow", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]string{
			"control_point": "feature",
			"fn.flow":       `-`,
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]string{
			"control_point": "feature",
			"labeled":       "false",
		}))
	})
})

func tracesFromLabels(labels map[string]string) ptrace.Traces {
	td := ptrace.NewTraces()
	traces := td.ResourceSpans().AppendEmpty().
		ScopeSpans().AppendEmpty().
		Spans()
	spanRecord := traces.AppendEmpty()
	attr := spanRecord.Attributes()
	for k, v := range labels {
		attr.InsertString(k, v)
	}
	return td
}

func assertTracesEqual(act, exp ptrace.Traces) {
	actTraces := act.ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	expTraces := exp.ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	Expect(actTraces.Len()).To(Equal(expTraces.Len()))
	for i := 0; i < expTraces.Len(); i++ {
		assertAttributesEqual(actTraces.At(i).Attributes(),
			expTraces.At(i).Attributes())
	}
}
