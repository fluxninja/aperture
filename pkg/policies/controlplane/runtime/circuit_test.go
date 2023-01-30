package runtime_test

import (
	"math"

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
						- { constant_signal: { value: 0.0 } }
				out_ports:
					output: { signal_name: SUM_OR_ZERO }
			- arithmetic_combinator:
				operator: add
				in_ports:
					lhs: { signal_name: SUM_OR_ZERO }
					rhs: { constant_signal: { value: 1.0 } }
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

	It("properly generates all kinds of constant signals", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- variable:
				default_config:
					constant_signal:
						value: 42
				out_ports:
					output: { signal_name: VAR }
			- variable:
				default_config:
					constant_signal:
						special_value: NaN
				out_ports:
					output: { signal_name: VAR2 }
			- variable:
				default_config:
					constant_signal:
						special_value: +Inf
				out_ports:
					output: { signal_name: VAR3 }
			- variable:
				default_config:
					constant_signal:
						special_value: -Inf
				out_ports:
					output: { signal_name: VAR4 }
			`,
			sim.Inputs(nil),
			sim.OutputSignals{"VAR", "VAR2", "VAR3", "VAR4"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			sim.StepOutputs{
				"VAR":  sim.NewReading(42),
				"VAR2": sim.NewReading(math.NaN()),
				"VAR3": sim.NewReading(math.Inf(1)),
				"VAR4": sim.NewReading(math.Inf(-1)),
			},
		))
	})
})
