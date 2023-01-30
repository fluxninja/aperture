package sim_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/sim"
)

var _ = Describe("Circuit", func() {
	It("can simulate empty circuit", func() {
		circuit, err := sim.NewCircuit(nil, nil, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(BeEmpty())
	})

	It("can use constant output", func() {
		circuit, err := sim.NewCircuit(
			nil,
			sim.Inputs{
				"X": components.NewConstantSignal(42.0),
			},
			sim.OutputSignals{"X"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			sim.StepOutputs{
				"X": sim.NewReading(42.0),
			},
		))
	})

	It("can use array-input", func() {
		circuit, err := sim.NewCircuit(
			nil,
			sim.Inputs{
				"XX": sim.NewInput([]float64{42.0, 43.0}),
			},
			sim.OutputSignals{"XX"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"XX": sim.NewReadings([]float64{42.0, 43.0}),
			},
		))
	})

	It("can simulate a single-component circuit", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- sqrt:
				scale: 1.0 # FIXME, for some reason default doesn't load
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SQRT }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{4.0, 9.0}),
			},
			sim.OutputSignals{"SQRT"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				// FIXME: Add some helpers to handle floating-point inaccuracies.
				// (Here we are lucky to have an exactly-representable answer)
				"SQRT": sim.NewReadings([]float64{2.0, 3.0}),
			},
		))
	})

	It("can simulate a multi-component circuit", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			components:
			- sqrt:
				scale: 1.0 # FIXME, for some reason default doesn't load
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SQRT }
			- min:
				in_ports:
					inputs:
						- { signal_name: SQRT }
						- { signal_name: CAP }
				out_ports:
					output: { signal_name: SQRT-CAPPED }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{4.0, 9.0}),
				"CAP":   sim.NewConstantInput(2.5),
			},
			sim.OutputSignals{"SQRT-CAPPED"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"SQRT-CAPPED": sim.NewReadings([]float64{2.0, 2.5}),
			},
		))
	})

	It("can handle invalid readings in tests", func() {
		circuit, err := sim.NewCircuit(
			nil,
			sim.Inputs{"NAN": sim.NewConstantInput(math.NaN())},
			sim.OutputSignals{"NAN"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			sim.StepOutputs{
				"NAN": sim.InvalidReading(),
			},
		))
	})

	// NOTE: This test file test a sanity of circuit simulator itself. Please
	// put tests of actual components in appropriate packages.
})
