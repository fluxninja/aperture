package reading_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"aperture.tech/aperture/pkg/policies/controlplane/reading"
)

var _ = Describe("Reading", func() {

	It("Creates default invalid reading", func() {
		reading := reading.NewInvalid()
		Expect(reading).ToNot(BeNil())
		Expect(reading.Valid).To(BeFalse())
	})

	It("Sets invalid to true if NaN is passed", func() {
		reading := reading.New(math.NaN())
		Expect(reading.Valid).To(BeFalse())
		Expect(math.IsNaN(reading.Value)).To(BeTrue())
	})

	It("Creates valid value", func() {
		reading := reading.New(0.2)
		Expect(reading.Valid).To(BeTrue())
		Expect(reading.Value).To(Equal(0.2))
	})
})
