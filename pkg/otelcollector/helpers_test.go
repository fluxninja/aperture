package otelcollector_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/v2/pkg/otelcollector"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

var _ = Describe("EnforceIncludeList", func() {
	It("removes entries outside include list", func() {
		includeList := utils.SliceToSet([]string{"foo"})
		attributes := pcommon.NewMap()
		attributes.PutStr("foo", "x")
		attributes.PutStr("bar", "x")
		otelcollector.EnforceIncludeList(attributes, includeList)
		Expect(attributes.AsRaw()).To(Equal(map[string]interface{}{
			"foo": "x",
		}))
	})
})

var _ = Describe("EnforceExcludeList", func() {
	It("removes entries from exclude list", func() {
		excludeList := utils.SliceToSet([]string{"foo"})
		attributes := pcommon.NewMap()
		attributes.PutStr("foo", "x")
		attributes.PutStr("bar", "x")
		otelcollector.EnforceExcludeList(attributes, excludeList)
		Expect(attributes.AsRaw()).To(Equal(map[string]interface{}{
			"bar": "x",
		}))
	})
})
