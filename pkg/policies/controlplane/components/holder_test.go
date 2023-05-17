package components_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/sim"
)

var _ = Describe("Holder", func() {
	It("holds value", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- holder:
				in_ports:
					input: { signal_name: INPUT }
				hold_for: 4s
				out_ports:
					output: { signal_name: HOLDER }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{1, nan, 3, nan, nan}),
			},
			sim.OutputSignals{"HOLDER"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"HOLDER": sim.NewReadings([]float64{1, 1, 1, 1, nan}),
			},
		))
	})

	It("immediately holds next valid signal", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- holder:
				in_ports:
					input: { signal_name: INPUT }
				hold_for: 4s
				out_ports:
					output: { signal_name: HOLDER }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{1, 2, 3, 4, 6, 7, 8}),
			},
			sim.OutputSignals{"HOLDER"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"HOLDER": sim.NewReadings([]float64{1, 1, 1, 1, 6, 6, 6}),
			},
		))
	})
})
