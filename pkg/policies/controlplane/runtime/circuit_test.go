package runtime_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/policies/controlplane/sim"
)

var _ = Describe("Circuit", func() {
	It("handles loops by using value from previous iteration", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- first_valid:
				in_ports:
					inputs:
						- { signal_name: SUM }
						- { constant_value: 0.0 }
				out_ports:
					output: { signal_name: SUM_OR_ZERO }
			- arithmetic_combinator:
				operator: add
				in_ports:
					lhs: { signal_name: SUM_OR_ZERO }
					rhs: { constant_value: 1.0 }
				out_ports:
					output: { signal_name: SUM }
			`,
			sim.Inputs(nil),
			sim.OutputSignals{"SUM"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Run(3)).To(Equal(
			sim.Outputs{
				"SUM": sim.NewReadings([]float64{1.0, 2.0, 3.0}),
			},
		))
	})
})
