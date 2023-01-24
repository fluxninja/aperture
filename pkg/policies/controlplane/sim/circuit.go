package sim

import (
	"errors"
	"fmt"
	"io"
	"time"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/circuitfactory"
	rt "github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/status"
)

// Introducing some newtypes so that tests themselves are more readable.

// Inputs map signal names to components that will emit test input.
// Input components are required to emit an "output" signal, such as Input or
// components.Variable. These can be created with NewInput and NewConstantInput.
type Inputs map[string]rt.Component

// OutputSignals is a list of signal names that comprise test output.
type OutputSignals []string

// Outputs map signal names to captured output readings.
type Outputs map[string][]Reading

// Circuit is a simulated circuit intended to be used in tests.
type Circuit struct {
	meta    *simPolicyMeta
	circuit *rt.Circuit
	inputs  Inputs
	outputs map[string]*output
	tickNo  int
	time    time.Time // virtual time of next tick
}

// NewCircuitFromYaml creates a new simulated Circuit based on yaml circuit description.
//
// Differences from real circuit:
// * Fx options for components will be ignored.
func NewCircuitFromYaml(
	circuitYaml string,
	inputs Inputs,
	outputSignals OutputSignals,
) (*Circuit, error) {
	var err error
	circuitYaml, err = SanitizeYaml(circuitYaml)
	if err != nil {
		return nil, err
	}

	var circuitProto policylangv1.Circuit

	if err = config.UnmarshalYAML([]byte(circuitYaml), &circuitProto); err != nil {
		return nil, err
	}

	policyMeta := newSimPolicyMeta(circuitProto.EvaluationInterval.AsDuration())

	components, _, _, err := circuitfactory.CreateComponents(circuitProto.Components, policyMeta)
	if err != nil {
		return nil, err
	}

	return newCircuit(components, inputs, outputSignals, policyMeta)
}

// NewCircuit creates a new simulated Circuit.
func NewCircuit(
	components []rt.ConfiguredComponent,
	inputs Inputs,
	outputSignals OutputSignals,
) (*Circuit, error) {
	return newCircuit(components, inputs, outputSignals, newSimPolicyMeta(500*time.Millisecond))
}

func newCircuit(
	components []rt.ConfiguredComponent,
	inputs Inputs,
	outputSignals OutputSignals,
	policyMeta *simPolicyMeta,
) (*Circuit, error) {
	for inputSignal, input := range inputs {
		components = append(components, ConfigureInputComponent(input, inputSignal))
	}

	outputs := make(map[string]*output, len(outputSignals))

	for _, outputSignal := range outputSignals {
		output := &output{}
		components = append(components, ConfigureOutputComponent(outputSignal, output))
		outputs[outputSignal] = output
	}

	compiledCircuit, err := rt.Compile(components, policyMeta.Registry.GetLogger())
	if err != nil {
		return nil, err
	}

	runtimeCircuit, _ := rt.NewCircuitAndOptions(compiledCircuit, policyMeta)
	if runtimeCircuit == nil {
		return nil, errors.New("cannot create circuit")
	}

	return &Circuit{
		meta:    policyMeta,
		circuit: runtimeCircuit,
		inputs:  inputs,
		outputs: outputs,
		time:    time.Now(),
		tickNo:  0,
	}, nil
}

// Step runs one tick of circuit execution and returns values of output signals.
func (s *Circuit) Step() map[string]Reading {
	s.execStep()

	outputs := make(map[string]Reading, len(s.outputs))
	for outputSignal, output := range s.outputs {
		readings := output.TakeReadings()
		if len(readings) != 1 {
			panic("unexpected output readings len")
		}
		outputs[outputSignal] = ReadingFromRt(readings[0])
	}
	return outputs
}

// Run runs given number of tick of circuit execution and returns values of output signals.
func (s *Circuit) Run(steps int) Outputs {
	for iStep := 0; iStep < steps; iStep++ {
		s.execStep()
	}

	outputs := make(map[string][]Reading, len(s.outputs))
	for outputSignal, output := range s.outputs {
		outputs[outputSignal] = ReadingsFromRt(output.TakeReadings())
	}
	return outputs
}

// RunDrainInputs runs the circuit for as long as inputs are defined.
//
// Returns values of output signals.
// There must be at least one input of type Input defined and all of them must
// have same lengths.
func (s *Circuit) RunDrainInputs() Outputs {
	return s.Run(s.inputLen())
}

func (s *Circuit) inputLen() int {
	steps := -1
	for _, input := range s.inputs {
		if input, ok := input.(*Input); ok {
			if steps != -1 && steps != len(*input) {
				panic("input len mismatch")
			}
			steps = len(*input)
		}
	}
	if steps == -1 {
		panic("no Inputs")
	}
	return steps
}

func (s *Circuit) execStep() {
	err := s.circuit.Execute(
		rt.NewTickInfo(
			s.time,
			s.time.Add(s.meta.EvaluationInterval),
			s.tickNo,
			s.meta.EvaluationInterval,
		),
	)
	s.tickNo += 1
	s.time = s.time.Add(s.meta.EvaluationInterval)
	if err != nil {
		// Note: Assuming circuit execution will usually be infallible and thus
		// all methods on CircuitSim are infallible. If a need will occur to
		// test a execution when component is expected to fail, consider
		// introducing new fallible methods, like TryStep, TryRun, etc.
		panic(fmt.Errorf("circuit.Execute errored: %w", err))
	}
}

// simPolicyMeta implements Policy interface for usage in simulated circuits.
type simPolicyMeta struct {
	EvaluationInterval time.Duration
	Registry           status.Registry
}

func newSimPolicyMeta(evaluationInternal time.Duration) *simPolicyMeta {
	logger := log.NewLogger(io.Discard, "panic")
	registry := status.NewRegistry(logger)
	return &simPolicyMeta{
		Registry:           registry,
		EvaluationInterval: evaluationInternal,
	}
}

// GetPolicyName implements PolicyReadAPI.
func (p *simPolicyMeta) GetPolicyName() string { return "test-policy" }

// GetPolicyHash implements PolicyReadAPI.
func (p *simPolicyMeta) GetPolicyHash() string { return "test-policy-hash" }

// GetEvaluationInterval implements PolicyReadAPI.
func (p *simPolicyMeta) GetEvaluationInterval() time.Duration { return p.EvaluationInterval }

// GetStatusRegistry implements PolicyReadAPI.
func (p *simPolicyMeta) GetStatusRegistry() status.Registry { return p.Registry }
