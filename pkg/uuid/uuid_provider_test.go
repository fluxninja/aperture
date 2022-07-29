package uuid_test

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/uuid"
)

var _ = Describe(" Provider", func() {
	var (
		provider uuid.Provider
		valid    *regexp.Regexp
	)

	BeforeEach(func() {
		rawRegex := "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
		valid = regexp.MustCompile(rawRegex)
	})

	Context("DefaultProvider", func() {
		BeforeEach(func() {
			provider = uuid.NewDefaultProvider()
		})

		It("Returns proper ", func() {
			id := provider.New()
			Expect(valid.MatchString(id)).To(BeTrue())
		})
	})

	Context("TestProvider", func() {
		var test string

		BeforeEach(func() {
			test = "deadbeef-dead-4ead-beef-deaddeafbeef"
			provider = uuid.NewTestProvider(test)
		})

		It("Returns proper ", func() {
			id := provider.New()
			Expect(valid.MatchString(id)).To(BeTrue())
		})

		It("Returns preset ", func() {
			id := provider.New()
			Expect(id).To(Equal(test))
		})
	})
})
