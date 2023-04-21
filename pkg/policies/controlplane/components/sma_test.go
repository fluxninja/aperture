package components_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/policies/controlplane/sim"
)

var _ = Describe("SMA component", func() {
	It("computes simple moving average with sma_window:2s, valid_during_warmup:false", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- sma:
				parameters:
					sma_window: 2s
					valid_during_warmup: false
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SMA }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0}),
			},
			sim.OutputSignals{"SMA"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"SMA": sim.NewReadings([]float64{nan, nan, 0.5, 1.5, 2.5, 3.5, 4.5, 5.5}),
			},
		))
	})

	It("computes simple moving average with sma_window:2s, valid_during_warmup:true", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- sma:
				parameters:
					sma_window: 2s
					valid_during_warmup: true
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SMA }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0}),
			},
			sim.OutputSignals{"SMA"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"SMA": sim.NewReadings([]float64{nan, 0.0, 0.5, 1.5, 2.5, 3.5, 4.5, 5.5}),
			},
		))
	})

	It("computes simple moving average with sma_window:5s, valid_during_warmup:false", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- sma:
				parameters:
					sma_window: 5s
					valid_during_warmup: false
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SMA }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0}),
			},
			sim.OutputSignals{"SMA"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"SMA": sim.NewReadings([]float64{nan, nan, nan, nan, nan, 2.0, 3.0, 4.0}),
			},
		))
	})

	It("computes simple moving average with sma_window:5s, valid_during_warmup:true", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 1s
			components:
			- sma:
				parameters:
					sma_window: 5s
					valid_during_warmup: true
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SMA }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0}),
			},
			sim.OutputSignals{"SMA"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"SMA": sim.NewReadings([]float64{nan, 0.0, 0.5, 1, 1.5, 2.0, 3.0, 4.0}),
			},
		))
	})

	It("computes simple moving average with evaluation_interval:0.5s, sma_window:5s, valid_during_warmup:true", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 0.5s
			components:
			- sma:
				parameters:
					sma_window: 5s
					valid_during_warmup: true
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SMA }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0}),
			},
			sim.OutputSignals{"SMA"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"SMA": sim.NewReadings([]float64{nan, 0.0, 0.5, 1, 1.5, 2.0, 2.5, 3.0}),
			},
		))
	})

	It("computes simple moving average with evaluation_interval:2s, sma_window:1s, valid_during_warmup:true", func() {
		circuit, err := sim.NewCircuitFromYaml(
			`
			evaluation_interval: 2s
			components:
			- sma:
				parameters:
					sma_window: 1s
					valid_during_warmup: true
				in_ports:
					input: { signal_name: INPUT }
				out_ports:
					output: { signal_name: SMA }
			`,
			sim.Inputs{
				"INPUT": sim.NewInput([]float64{nan, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0}),
			},
			sim.OutputSignals{"SMA"},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(circuit.RunDrainInputs()).To(Equal(
			sim.Outputs{
				"SMA": sim.NewReadings([]float64{nan, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}),
			},
		))
	})
})
