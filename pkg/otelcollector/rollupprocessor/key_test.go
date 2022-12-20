package rollupprocessor

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var _ = Describe("Rollup key", func() {
	It("ignores order", func() {
		m1 := pcommon.NewMap()
		m1.PutStr("foo", "fooval")
		m1.PutStr("bar", "barval")
		k1 := key(m1, nil)
		m2 := pcommon.NewMap()
		m2.PutStr("bar", "barval")
		m2.PutStr("foo", "fooval")
		k2 := key(m2, nil)
		Expect(k1).To(Equal(k2))
	})

	It("differs when attrs differ", func() {
		m1 := pcommon.NewMap()
		m1.PutStr("foo", "fooval")
		m1.PutStr("bar", "barval")
		k1 := key(m1, nil)
		m2 := pcommon.NewMap()
		m2.PutStr("foo", "fooval")
		m2.PutStr("bar", "barval2")
		k2 := key(m2, nil)
		Expect(k1).NotTo(Equal(k2))
	})

	It("differs when slice attrs differ", func() {
		m1 := pcommon.NewMap()
		Expect(m1.FromRaw(map[string]any{
			"foo":  "fooval",
			"vals": []any{"x, y"},
		})).To(Succeed())
		k1 := key(m1, nil)
		m2 := pcommon.NewMap()
		Expect(m2.FromRaw(map[string]any{
			"foo":  "fooval",
			"vals": []any{"x, z"},
		})).To(Succeed())
		k2 := key(m2, nil)
		Expect(k1).NotTo(Equal(k2))
	})

	It("ignores ignored attrs", func() {
		ignored := map[string]struct{}{"bar": {}}
		m1 := pcommon.NewMap()
		m1.PutStr("foo", "fooval")
		m1.PutStr("bar", "barval")
		k1 := key(m1, ignored)
		m2 := pcommon.NewMap()
		m2.PutStr("foo", "fooval")
		m2.PutStr("bar", "barval2")
		k2 := key(m2, ignored)
		Expect(k1).To(Equal(k2))
	})
})
