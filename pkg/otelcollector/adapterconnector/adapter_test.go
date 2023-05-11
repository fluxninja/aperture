package adapterconnector_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/connector/connectortest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/fluxninja/aperture/v2/pkg/otelcollector/adapterconnector"
)

var _ = Describe("Adapter", func() {
	var (
		f   connector.Factory
		cfg component.Config
	)
	BeforeEach(func() {
		f = adapterconnector.NewFactory()
		Expect(f.Type()).To(BeEquivalentTo("adapter"))
		cfg = f.CreateDefaultConfig()
		Expect(cfg).To(Equal(&struct{}{}))
	})
	It("converts traces to logs", func() {
		ctx := context.Background()
		set := connectortest.NewNopCreateSettings()
		host := componenttest.NewNopHost()

		logsSink := new(consumertest.LogsSink)
		tracesToLogs, err := f.CreateTracesToLogs(ctx, set, cfg, logsSink)
		Expect(err).To(BeNil())
		err = tracesToLogs.Start(ctx, host)
		Expect(err).To(BeNil())

		err = tracesToLogs.ConsumeTraces(ctx, ptrace.NewTraces())
		Expect(err).To(BeNil())

		err = tracesToLogs.Shutdown(ctx)
		Expect(err).To(BeNil())
		Expect(logsSink.AllLogs()).To(HaveLen(1))
	})
})
