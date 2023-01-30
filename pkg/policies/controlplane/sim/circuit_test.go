package sim_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
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
				runtime.MakeRootSignalID("X"): components.NewConstantSignal(42.0),
			},
			sim.OutputSignals{runtime.MakeRootSignalID("X")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			map[runtime.SignalID]sim.Reading{
				runtime.MakeRootSignalID("X"): sim.NewReading(42.0),
			},
		))
	})

	It("can use array-input", func() {
		circuit, err := sim.NewCircuit(
			nil,
			sim.Inputs{
				runtime.MakeRootSignalID("XX"): sim.NewInput([]float64{42.0, 43.0}),
			},
			sim.OutputSignals{runtime.MakeRootSignalID("XX")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.MakeRootSignalID("XX"): sim.NewReadings([]float64{42.0, 43.0}),
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
				runtime.MakeRootSignalID("INPUT"): sim.NewInput([]float64{4.0, 9.0}),
			},
			sim.OutputSignals{runtime.MakeRootSignalID("SQRT")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				// FIXME: Add some helpers to handle floating-point inaccuracies.
				// (Here we are lucky to have an exactly-representable answer)
				runtime.MakeRootSignalID("SQRT"): sim.NewReadings([]float64{2.0, 3.0}),
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
				runtime.MakeRootSignalID("INPUT"): sim.NewInput([]float64{4.0, 9.0}),
				runtime.MakeRootSignalID("CAP"):   sim.NewConstantInput(2.5),
			},
			sim.OutputSignals{runtime.MakeRootSignalID("SQRT-CAPPED")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.MakeRootSignalID("SQRT-CAPPED"): sim.NewReadings([]float64{2.0, 2.5}),
			},
		))
	})

	It("can handle invalid readings in tests", func() {
		circuit, err := sim.NewCircuit(
			nil,
			sim.Inputs{runtime.MakeRootSignalID("NAN"): sim.NewConstantInput(math.NaN())},
			sim.OutputSignals{runtime.MakeRootSignalID("NAN")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			map[runtime.SignalID]sim.Reading{
				runtime.MakeRootSignalID("NAN"): sim.InvalidReading(),
			},
		))
	})

	// NOTE: This test file test a sanity of circuit simulator itself. Please
	// put tests of actual components in appropriate packages.
})
