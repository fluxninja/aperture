package enrichmentprocessor

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Enrichment Processor - Metrics", func() {
	It("Enriches metrics attributes with data from entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID: entitycache.EntityID{
				Prefix: "testPrefix",
				UID:    "1",
			},
			AgentGroup: "fooGroup",
			Namespace:  "fooNS",
			Services:   []string{"fooSvc1", "fooSvc2"},
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		md := metricsFromLabels(map[string]string{
			"preserve":    "this",
			"entity_name": "testPrefix-1",
		})
		md, err := processor.ConsumeMetrics(context.TODO(), md)
		Expect(err).NotTo(HaveOccurred())

		assertMetricsEqual(md, metricsFromLabels(map[string]string{
			"preserve":    "this",
			"agent_group": "fooGroup",
			"namespace":   "fooNS",
			"services":    "fooSvc1,fooSvc2",
		}))
	})

	It("Does not enrich when there are no labels in entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		md := metricsFromLabels(map[string]string{
			"preserve":    "this",
			"entity_name": "bar",
		})
		md, err := processor.ConsumeMetrics(context.TODO(), md)
		Expect(err).NotTo(HaveOccurred())

		assertMetricsEqual(md, metricsFromLabels(map[string]string{
			"preserve": "this",
		}))
	})
})

func metricsFromLabels(labels map[string]string) pmetric.Metrics {
	td := pmetric.NewMetrics()
	metrics := td.ResourceMetrics().AppendEmpty().
		ScopeMetrics().AppendEmpty().Metrics()
	metric := metrics.AppendEmpty()
	metric.SetDataType(pmetric.MetricDataTypeGauge)
	spanRecord := metric.Gauge()
	attr := spanRecord.DataPoints().AppendEmpty().Attributes()
	for k, v := range labels {
		attr.InsertString(k, v)
	}
	return td
}

func assertMetricsEqual(act, exp pmetric.Metrics) {
	actMetrics := act.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	expMetrics := exp.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	Expect(actMetrics.Len()).To(Equal(expMetrics.Len()))
	for i := 0; i < expMetrics.Len(); i++ {
		assertAttributesEqual(actMetrics.At(i).Gauge().DataPoints().At(0).Attributes(),
			expMetrics.At(i).Gauge().DataPoints().At(0).Attributes())
	}
}
