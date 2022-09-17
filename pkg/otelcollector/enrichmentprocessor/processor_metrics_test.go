package enrichmentprocessor

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pmetric"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Enrichment Processor - Metrics", func() {
	It("Enriches metrics attributes with data from entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycachev1.Entity{
			Prefix:    "testPrefix",
			Uid:       "1",
			IpAddress: "",
			Name:      "someName",
			Services:  []string{"fooSvc1", "fooSvc2"},
		})

		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		md := metricsFromLabels(map[string]interface{}{
			"preserve":    "this",
			"entity_name": "someName",
		})
		md, err := processor.ConsumeMetrics(context.TODO(), md)
		Expect(err).NotTo(HaveOccurred())

		assertMetricsEqual(md, metricsFromLabels(map[string]interface{}{
			"preserve": "this",
			"services": []string{"fooSvc1", "fooSvc2"},
		}))
	})

	It("Does not enrich when there are no labels in entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		md := metricsFromLabels(map[string]interface{}{
			"preserve":    "this",
			"entity_name": "bar",
		})
		md, err := processor.ConsumeMetrics(context.TODO(), md)
		Expect(err).NotTo(HaveOccurred())

		assertMetricsEqual(md, metricsFromLabels(map[string]interface{}{
			"preserve": "this",
		}))
	})
})

func metricsFromLabels(labels map[string]interface{}) pmetric.Metrics {
	td := pmetric.NewMetrics()
	metrics := td.ResourceMetrics().AppendEmpty().
		ScopeMetrics().AppendEmpty().Metrics()
	metric := metrics.AppendEmpty()
	metric.SetDataType(pmetric.MetricDataTypeGauge)
	spanRecord := metric.Gauge()
	populateAttrsFromLabels(spanRecord.DataPoints().AppendEmpty().Attributes(), labels)
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
