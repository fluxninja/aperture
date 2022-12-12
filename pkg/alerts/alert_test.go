package alerts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/alerts"
)

var _ = Describe("Alert", func() {
	It("transforms both ways properly", func() {
		alert := alerts.NewAlert(
			alerts.WithAnnotation("one", "eleven"),
			alerts.WithAnnotation("two", "twelve"),
			alerts.WithName("buzz"),
			alerts.WithSeverity(alerts.SeverityCrit),
		)
		finalAlert := alert
		finalAlert.SetLabel("is_alert", "true")
		ld := alert.AsLogs()
		Expect(ld.LogRecordCount()).To(Equal(1))
		actual := alerts.AlertsFromLogs(ld)
		Expect(actual).To(HaveLen(1))
		Expect(actual[0]).To(Equal(finalAlert))
	})
})

func expectKeyWithValue(attributes pcommon.Map, key, value string) {
	rawValue, exists := attributes.Get(key)
	Expect(exists).To(BeTrue())
	Expect(rawValue.AsString()).To(Equal(value))
}
