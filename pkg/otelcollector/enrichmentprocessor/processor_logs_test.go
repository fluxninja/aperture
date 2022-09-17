package enrichmentprocessor

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/plog"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

var _ = Describe("Enrichment Processor - Logs", func() {
	It("Enriches egress logs attributes with data from entity cache", func() {
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

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.HostAddressLabel:  "192.0.2.0:80",
			otelcollector.PeerAddressLabel:  "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.LabeledLabel:      "false",
			"services":                      []string{"svc1", "svc2"},
		}))
	})

	It("Does not panic when egress logs net.host.address attribute empty", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycachev1.Entity{
			Prefix:    "",
			Uid:       "",
			IpAddress: hardCodedIPAddress,
			Name:      hardCodedEntityName,
			Services:  nil,
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.HostAddressLabel:  "-",
			otelcollector.PeerAddressLabel:  "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.LabeledLabel:      "false",
		}))
	})

	It("Does not panic when egress logs net.host.address attribute is garbage", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycachev1.Entity{
			Prefix:    "",
			Uid:       "",
			IpAddress: hardCodedIPAddress,
			Name:      hardCodedEntityName,
			Services:  nil,
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.HostAddressLabel:  "this is :: definitely not an IP ::::",
			otelcollector.PeerAddressLabel:  "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.LabeledLabel:      "false",
		}))
	})

	It("Enriches ingress logs attributes with data from entity cache", func() {
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

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointIngress,
			"net.host.ip":                   "192.0.2.0",
			"net.peer.ip":                   "192.0.2.1",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointIngress,
			otelcollector.LabeledLabel:      "false",
			"services":                      []string{"svc1", "svc2"},
		}))
	})

	It("Does not enrich when there are no labels in entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.HostAddressLabel:  "192.0.2.0:80",
			otelcollector.PeerAddressLabel:  "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.LabeledLabel:      "false",
		}))
	})

	It("Unpacks aperture.labels properly", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			"aperture.labels":               `{"foo": "bar", "fizz": "buzz"}`,
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			"foo":                           "bar",
			"fizz":                          "buzz",
			otelcollector.LabeledLabel:      "true",
		}))
	})

	It("Ignores empty aperture.labels", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			"aperture.labels":               ``,
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointEgress,
			otelcollector.LabeledLabel:      "false",
		}))
	})

	It("Ignores minus as aperture.labels", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			"aperture.labels":               `-`,
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]interface{}{
			otelcollector.ControlPointLabel: otelcollector.ControlPointFeature,
			otelcollector.LabeledLabel:      "false",
		}))
	})
})

func logsFromLabels(labels map[string]interface{}) plog.Logs {
	ld := plog.NewLogs()
	logs := ld.ResourceLogs().AppendEmpty().
		ScopeLogs().AppendEmpty().
		LogRecords()
	logRecord := logs.AppendEmpty()
	populateAttrsFromLabels(logRecord.Attributes(), labels)
	return ld
}

func assertLogsEqual(act, exp plog.Logs) {
	actLogs := act.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
	expLogs := exp.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
	Expect(actLogs.Len()).To(Equal(expLogs.Len()))
	for i := 0; i < expLogs.Len(); i++ {
		assertAttributesEqual(actLogs.At(i).Attributes(),
			expLogs.At(i).Attributes())
	}
}
