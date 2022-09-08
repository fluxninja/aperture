package enrichmentprocessor

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/config"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Enrichment Processor", func() {
	It("Creates default config", func() {
		entityCache := entitycache.NewEntityCache()
		expected := &Config{
			ProcessorSettings: config.NewProcessorSettings(config.NewComponentID("enrichment")),
			entityCache:       entityCache,
		}
		actual := createDefaultConfig(entityCache)()

		Expect(actual).To(Equal(expected))
	})
})
