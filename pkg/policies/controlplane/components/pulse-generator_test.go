package components_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/sim"
)

var _ = Describe("PulseGenerator", func() {
	It("generates pulse", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 10s
			components:
			- pulse_generator:
				true_for: 30s
				false_for: 20s
				out_ports:
					output: { signal_name: PULSE }
			`,
			nil,
			sim.OutputSignals{"PULSE"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Run(10)).To(Equal(
			sim.Outputs{
				"PULSE": sim.NewReadings([]float64{1, 1, 1, 0, 0, 1, 1, 1, 0, 0}),
			},
		))
	})

	It("rounds up small pulse durations", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- pulse_generator:
				true_for: 0.5s
				false_for: 0.5s
				out_ports:
					output: { signal_name: PULSE }
			`,
			nil,
			sim.OutputSignals{"PULSE"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Run(4)).To(Equal(
			sim.Outputs{
				"PULSE": sim.NewReadings([]float64{1, 0, 1, 0}),
			},
		))
	})

	It("can use defaults", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- pulse_generator:
				out_ports:
					output: { signal_name: PULSE }
			`,
			nil,
			sim.OutputSignals{"PULSE"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Run(11)).To(Equal(
			sim.Outputs{
				"PULSE": sim.NewReadings([]float64{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1}),
			},
		))
	})
})
