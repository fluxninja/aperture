package enrichmentprocessor

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/ptrace"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

var _ = Describe("Enrichment Processor - Traces", func() {
	It("Enriches feature trace attributes with data from entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycachev1.Entity{
			Prefix:    "",
			Uid:       "",
			IpAddress: hardCodedIPAddress,
			Name:      hardCodedEntityName,
			Services:  hardCodedServices,
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel:   otelcollector.ControlPointFeature,
			otelcollector.FeatureAddressLabel: "192.0.2.0",
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			otelcollector.LabeledLabel:      "false",
			otelcollector.ServicesLabel:     []string{"svc1", "svc2"},
		}))
	})

	It("Does not panic when egress metrics FeatureAddressLabel attribute empty", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycachev1.Entity{
			Prefix:    "",
			Uid:       "",
			IpAddress: hardCodedIPAddress,
			Name:      hardCodedEntityName,
			Services:  hardCodedServices,
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel:   otelcollector.ControlPointFeature,
			otelcollector.FeatureAddressLabel: "",
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			otelcollector.LabeledLabel:      "false",
		}))
	})

	It("Does not enrich when there is no matching entries in entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycachev1.Entity{
			Prefix:    "",
			Uid:       "",
			IpAddress: "192.0.2.3",
			Name:      hardCodedEntityName,
			Services:  hardCodedServices,
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel:   otelcollector.ControlPointFeature,
			otelcollector.FeatureAddressLabel: "192.0.2.0",
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			otelcollector.LabeledLabel:      "false",
		}))
	})

	It("Unpacks aperture.labels properly", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel:     otelcollector.ControlPointFeature,
			otelcollector.MarshalledLabelsLabel: `{"foo": "bar", "fizz": "buzz"}`,
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			"foo":                           "bar",
			"fizz":                          "buzz",
			otelcollector.LabeledLabel:      "true",
		}))
	})

	It("Ignores empty aperture.labels", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel:     otelcollector.ControlPointFeature,
			otelcollector.MarshalledLabelsLabel: ``,
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			otelcollector.LabeledLabel:      "false",
		}))
	})

	It("Ignores minus as aperture.labels", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		td := tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel:     otelcollector.ControlPointFeature,
			otelcollector.MarshalledLabelsLabel: `-`,
		})
		td, err := processor.ConsumeTraces(context.TODO(), td)
		Expect(err).NotTo(HaveOccurred())

		assertTracesEqual(td, tracesFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			otelcollector.LabeledLabel:      "false",
		}))
	})
})

func tracesFromLabels(labels map[string]interface{}) ptrace.Traces {
	td := ptrace.NewTraces()
	traces := td.ResourceSpans().AppendEmpty().
		ScopeSpans().AppendEmpty().
		Spans()
	spanRecord := traces.AppendEmpty()
	populateAttrsFromLabels(spanRecord.Attributes(), labels)
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
