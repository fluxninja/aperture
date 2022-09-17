package reading_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/policies/controlplane/reading"
)

var _ = Describe("Reading", func() {
	It("Creates default invalid reading", func() {
		reading := reading.InvalidReading()
		Expect(reading).ToNot(BeNil())
		Expect(reading.Valid()).To(BeFalse())
	})

	It("Creates new reading with invalid value", func() {
		reading := reading.NewReading((math.NaN()))
		Expect(reading.Valid()).To(BeFalse())
		Expect(math.IsNaN(reading.Value())).To(BeTrue())
	})

	It("Creates valid value", func() {
		reading := reading.NewReading(0.2)
		Expect(reading.Valid()).To(BeTrue())
		Expect(reading.Value()).To(Equal(0.2))
	})
})
