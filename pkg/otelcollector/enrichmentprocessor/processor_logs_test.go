package enrichmentprocessor

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Enrichment Processor - Logs", func() {
	It("Enriches egress logs attributes with data from entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID:         entitycache.EntityID{},
			IPAddress:  "192.0.2.0",
			AgentGroup: "defaultAG",
			Services:   []string{"svc1", "svc2"},
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point":    "egress",
			"net.host.address": "192.0.2.0:80",
			"net.peer.address": "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point": "egress",
			"labeled":       "false",
			"agent_group":   "defaultAG",
			"services":      "svc1,svc2",
		}))
	})

	It("Does not panic when egress logs net.host.address attribute empty", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID:        entitycache.EntityID{},
			IPAddress: "192.0.2.0",
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point":    "egress",
			"net.host.address": "-",
			"net.peer.address": "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point":    "egress",
			"labeled":          "false",
			"net.host.address": "-",
			"net.peer.address": "192.0.2.1:80",
		}))
	})

	It("Does not panic when egress logs net.host.address attribute is garbage", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID:        entitycache.EntityID{},
			IPAddress: "192.0.2.0",
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point":    "egress",
			"net.host.address": "this is :: definitely not an IP ::::",
			"net.peer.address": "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point": "egress",
			"labeled":       "false",
		}))
	})

	It("Enriches ingress logs attributes with data from entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		entityCache.Put(&entitycache.Entity{
			ID:         entitycache.EntityID{},
			IPAddress:  "192.0.2.0",
			AgentGroup: "defaultAG",
			Services:   []string{"svc1", "svc2"},
		})
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point": "ingress",
			"net.host.ip":   "192.0.2.0",
			"net.peer.ip":   "192.0.2.1",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point": "ingress",
			"labeled":       "false",
			"agent_group":   "defaultAG",
			"services":      "svc1,svc2",
		}))
	})

	It("Does not enrich when there are no labels in entity cache", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point":    "egress",
			"net.host.address": "192.0.2.0:80",
			"net.peer.address": "192.0.2.1:80",
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point": "egress",
			"labeled":       "false",
		}))
	})

	It("Unpacks aperture.labels properly", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point":   "egress",
			"aperture.labels": `{"foo": "bar", "fizz": "buzz"}`,
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point": "egress",
			"foo":           "bar",
			"fizz":          "buzz",
			"labeled":       "true",
		}))
	})

	It("Ignores empty aperture.labels", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point":   "egress",
			"aperture.labels": ``,
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point": "egress",
			"labeled":       "false",
		}))
	})

	It("Ignores minus as aperture.labels", func() {
		entityCache := entitycache.NewEntityCache()
		processor := newProcessor(entityCache)
		Expect(processor).NotTo(BeNil())

		ld := logsFromLabels(map[string]string{
			"control_point":   "feature",
			"aperture.labels": `-`,
		})
		ld, err := processor.ConsumeLogs(context.TODO(), ld)
		Expect(err).NotTo(HaveOccurred())

		assertLogsEqual(ld, logsFromLabels(map[string]string{
			"control_point": "feature",
			"labeled":       "false",
		}))
	})
})

func logsFromLabels(labels map[string]string) plog.Logs {
	ld := plog.NewLogs()
	logs := ld.ResourceLogs().AppendEmpty().
		ScopeLogs().AppendEmpty().
		LogRecords()
	logRecord := logs.AppendEmpty()
	attr := logRecord.Attributes()
	for k, v := range labels {
		attr.InsertString(k, v)
	}
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
