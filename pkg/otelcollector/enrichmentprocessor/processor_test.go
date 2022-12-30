package enrichmentprocessor

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Enrichment Processor", func() {
	It("Creates default config", func() {
		entityCache := entitycache.NewEntityCache()
		expected := &Config{
			entityCache: entityCache,
		}
		actual := createDefaultConfig(entityCache)()

		Expect(actual).To(Equal(expected))
	})
})
