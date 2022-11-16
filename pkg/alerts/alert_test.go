package alerts_test

import (
	"github.com/go-openapi/strfmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/alertmanager/api/v2/models"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/alerts"
)

var _ = Describe("Alert", func() {
	It("transforms both ways properly", func() {
		alert := exampleAlert()
		ld := alert.AsLogs()
		Expect(ld.LogRecordCount()).To(Equal(1))
		actual := alerts.AlertsFromLogs(ld)
		Expect(actual).To(HaveLen(1))
		Expect(actual[0]).To(Equal(alert))
	})
})

func expectKeyWithValue(attributes pcommon.Map, key, value string) {
	rawValue, exists := attributes.Get(key)
	Expect(exists).To(BeTrue())
	Expect(rawValue.AsString()).To(Equal(value))
}

func exampleAlert() *alerts.Alert {
	startsAt, err := strfmt.ParseDateTime("2022-10-26T07:44:14.101Z")
	Expect(err).NotTo(HaveOccurred())
	return &alerts.Alert{
		models.PostableAlert{
			StartsAt: startsAt,
			Annotations: map[string]string{
				"one": "eleven",
				"two": "twelve",
			},
			Alert: models.Alert{
				GeneratorURL: strfmt.URI("https://www.example.org/alertSource"),
				Labels: map[string]string{
					"alertname": "buzz",
					"severity":  "error",
				},
			},
		},
	}
}
