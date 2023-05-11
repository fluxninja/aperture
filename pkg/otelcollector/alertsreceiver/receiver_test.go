package alertsreceiver

import (
	"context"
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/v2/pkg/alerts"
)

var _ = Describe("Alerts receiver", func() {
	var (
		alerter  alerts.Alerter
		tc       *testConsumer
		receiver *alertsReceiver
	)

	BeforeEach(func() {
		alerter = alerts.NewSimpleAlerter(1)
		cfg := &Config{
			alerter: alerter,
		}
		var err error
		receiver, err = newReceiver(cfg)
		Expect(err).NotTo(HaveOccurred())
		tc = newTestConsumer()
		receiver.logsConsumer = tc

		go func() {
			err := receiver.Start(nil, nil)
			Expect(err).NotTo(HaveOccurred())
		}()
	})

	AfterEach(func() {
		err := receiver.Shutdown(nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("consumes single alert properly", func() {
		alert := alerts.NewAlert(alerts.WithName("foo"))
		Eventually(func() error {
			alerter.AddAlert(alert)
			return nil
		}).Should(Succeed())

		Eventually(func() int {
			return len(tc.ReceivedLogs())
		}).Should(Equal(1))

		Expect(tc.ReceivedLogs()[0]).To(Equal(alert.AsLogs()))
	})

	It("consumes multiple alerts properly", func() {
		alertsObj := []*alerts.Alert{}
		for i := 0; i < 10; i++ {
			alertsObj = append(alertsObj, alerts.NewAlert(
				alerts.WithName(fmt.Sprintf("foo%v", i)),
			))
		}

		for i := 0; i < 10; i++ {
			Eventually(func() error {
				alerter.AddAlert(alertsObj[i])
				return nil
			}).Should(Succeed())
		}

		Eventually(func() int {
			return len(tc.ReceivedLogs())
		}).Should(Equal(10))

		for i := 0; i < 10; i++ {
			Expect(tc.ReceivedLogs()).To(ContainElement(alertsObj[i].AsLogs()))
		}
	})
})

func newTestConsumer() *testConsumer {
	return &testConsumer{receivedLogs: []plog.Logs{}}
}

type testConsumer struct {
	sync.Mutex
	receivedLogs []plog.Logs
}

func (t *testConsumer) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	t.Lock()
	defer t.Unlock()
	t.receivedLogs = append(t.receivedLogs, ld)
	return nil
}

func (t *testConsumer) ReceivedLogs() []plog.Logs {
	t.Lock()
	defer t.Unlock()
	return t.receivedLogs
}

func (t *testConsumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: false,
	}
}
