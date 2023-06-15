package components_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/sim"
)

var _ = Describe("Not component", func() {
	It("performs logical negation", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- inverter:
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: NOT }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, -1.}),
			},
			sim.OutputSignals{"NOT"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"NOT": sim.NewReadings([]float64{nan, 1.0, 0.0, 0.0, 0.0}),
			},
		))
	})
})

var _ = Describe("And and Or component", func() {
	It("performs logical operations", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- and:
				in_ports:
					inputs:
					- { signal_name: INPUTX }
					- { signal_name: INPUTY }
					- { signal_name: INPUTZ }
				out_ports:
					output: { signal_name: AND }
			- or:
				in_ports:
					inputs:
					- { signal_name: INPUTX }
					- { signal_name: INPUTY }
					- { signal_name: INPUTZ }
				out_ports:
					output: { signal_name: OR }
			`,
			sim.Inputs{
				"INPUTX": sim.NewInput([]float64{nan, 0.0, 1.0, 0.0, 1.0, 1.0, 1.0, -1.}),
				"INPUTY": sim.NewInput([]float64{nan, nan, nan, 0.0, 0.0, 1.0, 2.0, -2.}),
				"INPUTZ": sim.NewInput([]float64{nan, 0.0, 1.0, 0.0, 1.0, 1.0, 3.0, 3.0}),
			},
			sim.OutputSignals{"AND", "OR"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"AND": sim.NewReadings([]float64{nan, 0.0, nan, 0.0, 0.0, 1.0, 1.0, 1.0}),
				"OR":  sim.NewReadings([]float64{nan, nan, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0}),
			},
		))
	})

	Specify("And returns true on no inputs", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- and:
				out_ports:
					output: { signal_name: AND }
			`,
			nil,
			sim.OutputSignals{"AND"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			sim.StepOutputs{
				"AND": sim.NewReading(1.0),
			},
		))
	})

	Specify("Or returns false on no inputs", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- or:
				out_ports:
					output: { signal_name: OR }
			`,
			nil,
			sim.OutputSignals{"OR"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			sim.StepOutputs{
				"OR": sim.NewReading(0.0),
			},
		))
	})
})

// so the input/output arrays are visually aligned
var nan = math.NaN()
