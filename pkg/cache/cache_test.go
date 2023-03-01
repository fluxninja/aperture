package cache_test

import (
	"github.com/fluxninja/aperture/pkg/cache"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	It("Doesn't crash/leak on start/stop", func() {
		c := cache.NewCache[int]()
		c.Start()
		c.Stop()
	})

	It("Handles epoch-based expiration", func() {
		c := cache.NewCache[string]()
		c.Put("a")
		c.Put("a")
		c.Put("b")
		Expect(c.GetAll()).To(ConsistOf("a", "b"))

		c.NewEpochForTest()
		c.Put("b")
		c.Put("c")
		Expect(c.GetAll()).To(ConsistOf("a", "b", "c"))
		Expect(c.Contains("a")).To(BeTrue())
		Expect(c.Contains("b")).To(BeTrue())
		Expect(c.Contains("c")).To(BeTrue())

		c.NewEpochForTest()
		Expect(c.GetAll()).To(ConsistOf("b", "c"))
		Expect(c.Contains("a")).To(BeFalse())

		c.NewEpochForTest()
		Expect(c.GetAll()).To(BeEmpty())
	})
})
