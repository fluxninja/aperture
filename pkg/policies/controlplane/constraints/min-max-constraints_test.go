package constraints_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/FluxNinja/aperture/pkg/policies/controlplane/constraints"
)

var _ = Describe("MinMaxConstraints", func() {
	var minmax *constraints.MinMaxConstraints

	BeforeEach(func() {
		minmax = constraints.NewMinMaxConstraints()
	})

	It("Constraints max", func() {
		minmax.SetMax(5.0)
		output, constraintType := minmax.Constrain(10.0)
		Expect(output).To(Equal(5.0))
		Expect(constraintType).To(Equal(constraints.MaxConstraint))
	})

	It("Constraints min", func() {
		minmax.SetMin(2.0)
		output, constraintType := minmax.Constrain(1.0)
		Expect(output).To(Equal(2.0))
		Expect(constraintType).To(Equal(constraints.MinConstraint))
	})
})
