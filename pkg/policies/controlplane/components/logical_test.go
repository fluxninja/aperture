package components_test

import (
	"math"

	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/sim"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
				runtime.MakeRootSignalID("INPUT"): sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, -1.}),
			},
			sim.OutputSignals{runtime.MakeRootSignalID("NOT")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.MakeRootSignalID("NOT"): sim.NewReadings([]float64{nan, 1.0, 0.0, 0.0, 0.0}),
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
				runtime.MakeRootSignalID("INPUTX"): sim.NewInput([]float64{nan, 0.0, 1.0, 0.0, 1.0, 1.0, 1.0, -1.}),
				runtime.MakeRootSignalID("INPUTY"): sim.NewInput([]float64{nan, nan, nan, 0.0, 0.0, 1.0, 2.0, -2.}),
				runtime.MakeRootSignalID("INPUTZ"): sim.NewInput([]float64{nan, 0.0, 1.0, 0.0, 1.0, 1.0, 3.0, 3.0}),
			},
			sim.OutputSignals{runtime.MakeRootSignalID("AND"), runtime.MakeRootSignalID("OR")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.MakeRootSignalID("AND"): sim.NewReadings([]float64{nan, 0.0, nan, 0.0, 0.0, 1.0, 1.0, 1.0}),
				runtime.MakeRootSignalID("OR"):  sim.NewReadings([]float64{nan, nan, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0}),
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
			sim.OutputSignals{runtime.MakeRootSignalID("AND")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			map[runtime.SignalID]sim.Reading{
				runtime.MakeRootSignalID("AND"): sim.NewReading(1.0),
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
			sim.OutputSignals{runtime.MakeRootSignalID("OR")},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			map[runtime.SignalID]sim.Reading{
				runtime.MakeRootSignalID("OR"): sim.NewReading(0.0),
			},
		))
	})
})

// so the input/output arrays are visually aligned
var nan = math.NaN()
