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
				runtime.SignalID{SignalName: "INPUT"}: sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, -1.}),
			},
			sim.OutputSignals{runtime.SignalID{SignalName: "NOT"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.SignalID{SignalName: "NOT"}: sim.NewReadings([]float64{nan, 1.0, 0.0, 0.0, 0.0}),
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
				runtime.SignalID{SignalName: "INPUTX"}: sim.NewInput([]float64{nan, 0.0, 1.0, 0.0, 1.0, 1.0, 1.0, -1.}),
				runtime.SignalID{SignalName: "INPUTY"}: sim.NewInput([]float64{nan, nan, nan, 0.0, 0.0, 1.0, 2.0, -2.}),
				runtime.SignalID{SignalName: "INPUTZ"}: sim.NewInput([]float64{nan, 0.0, 1.0, 0.0, 1.0, 1.0, 3.0, 3.0}),
			},
			sim.OutputSignals{runtime.SignalID{SignalName: "AND"}, runtime.SignalID{SignalName: "OR"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.SignalID{SignalName: "AND"}: sim.NewReadings([]float64{nan, 0.0, nan, 0.0, 0.0, 1.0, 1.0, 1.0}),
				runtime.SignalID{SignalName: "OR"}:  sim.NewReadings([]float64{nan, nan, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0}),
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
			sim.OutputSignals{runtime.SignalID{SignalName: "AND"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			map[runtime.SignalID]sim.Reading{
				{SignalName: "AND"}: sim.NewReading(1.0),
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
			sim.OutputSignals{runtime.SignalID{SignalName: "OR"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.Step()).To(Equal(
			map[runtime.SignalID]sim.Reading{
				{SignalName: "OR"}: sim.NewReading(0.0),
			},
		))
	})
})

// so the input/output arrays are visually aligned
var nan = math.NaN()
