package alerts_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

var _ = Describe("Alert", func() {
	Context("AsLogs()", func() {
		It("works properly", func() {
			alert := exampleAlert()
			ld := alert.AsLogs()
			Expect(ld.LogRecordCount()).To(Equal(1))
			otelcollector.IterateLogRecords(ld, func(lr plog.LogRecord) error {
				attributes := lr.Attributes()
				expectKeyWithValue(attributes, otelcollector.AlertStartsAtLabel, alert.StartsAt)
				expectKeyWithValue(attributes, otelcollector.AlertGeneratorURLLabel, alert.GeneratorURL)
				for k, v := range alert.Annotations {
					expectKeyWithValue(attributes, otelcollector.AlertAnnotationsLabelPrefix+k, v)
				}
				for k, v := range alert.Labels {
					expectKeyWithValue(attributes, otelcollector.AlertLabelsLabelPrefix+k, v)
				}
				return nil
			})
		})
	})

	Context("AlertsFromLogs()", func() {
		It("works properly", func() {
			ld := exampleLogs()
			alerts := alerts.AlertsFromLogs(ld)
			Expect(alerts).To(HaveLen(2))
			Expect(alerts).To(ContainElement(exampleAlert(0)))
			Expect(alerts).To(ContainElement(exampleAlert(1)))
		})
	})
})

func expectKeyWithValue(attributes pcommon.Map, key, value string) {
	rawValue, exists := attributes.Get(key)
	Expect(exists).To(BeTrue())
	Expect(rawValue.AsString()).To(Equal(value))
}

func exampleAlert(index ...int) *alerts.Alert {
	if len(index) == 0 {
		index = []int{0}
	}
	idx := index[0]
	return &alerts.Alert{
		StartsAt:     fmt.Sprintf("2022-10-26T07:44:14.10%vZ", idx),
		GeneratorURL: fmt.Sprintf("https://www.example.org/alertSource%v", idx),
		Annotations: map[string]string{
			"one": fmt.Sprintf("eleven%v", idx),
			"two": fmt.Sprintf("twelve%v", idx),
		},
		Labels: map[string]string{
			"foo":  fmt.Sprintf("bar%v", idx),
			"fizz": fmt.Sprintf("buzz%v", idx),
		},
	}
}

func exampleLogs() plog.Logs {
	ld := plog.NewLogs()
	logRecords := ld.
		ResourceLogs().AppendEmpty().
		ScopeLogs().AppendEmpty().
		LogRecords()
	for i := 0; i < 2; i++ {
		alert := exampleAlert(i)
		lr := logRecords.AppendEmpty()
		attributes := lr.Attributes()
		attributes.PutStr(otelcollector.AlertStartsAtLabel, alert.StartsAt)
		attributes.PutStr(otelcollector.AlertGeneratorURLLabel, alert.GeneratorURL)
		for k, v := range alert.Annotations {
			attributes.PutStr(otelcollector.AlertAnnotationsLabelPrefix+k, v)
		}
		for k, v := range alert.Labels {
			attributes.PutStr(otelcollector.AlertLabelsLabelPrefix+k, v)
		}
	}
	return ld
}
