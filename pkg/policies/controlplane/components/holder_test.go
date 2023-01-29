package components_test

import (
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/sim"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
				runtime.SignalID{SignalName: "INPUT"}: sim.NewInput([]float64{1, nan, 3, nan, nan}),
			},
			sim.OutputSignals{runtime.SignalID{SignalName: "HOLDER"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.SignalID{SignalName: "HOLDER"}: sim.NewReadings([]float64{1, 1, 1, 1, nan}),
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
				runtime.SignalID{SignalName: "INPUT"}: sim.NewInput([]float64{1, 2, 3, 4, 6, 7, 8}),
			},
			sim.OutputSignals{runtime.SignalID{SignalName: "HOLDER"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				runtime.SignalID{SignalName: "HOLDER"}: sim.NewReadings([]float64{1, 1, 1, 1, 6, 6, 6}),
			},
		))
	})
})
