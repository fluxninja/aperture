package runtime

import (
	"errors"
	"fmt"
	"math"

	"github.com/mitchellh/mapstructure"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/mapstruct"
)

// PortMapping is description of a component's ports mapping.
//
// This struct is meant to be decodable from Mapstruct representation of _any_
// of the components's config. Eg. EMA component defines:
//
// ```proto
//
//	message Ins {
//	  InPort input = 1;
//	}
//
// Ins in_ports = 1;
// ...
// ```
//
// And such EMA component's config could be encoded and then decoded into
// PortMapping as:
//
// ```go
//
//	PortMapping {
//	  InPorts: map[string]InPort {
//	    "input": []InPort {{ ... }},
//	  },
//	}
//
// ```
//
// Note how "input" is a concrete field in EMA definition, but a dynamic map
// key in PortMapping.
type PortMapping struct {
	// Note: Not using policylangv1.InPort and OutPort directly to avoid
	// runtime depending on proto.
	Ins  PortToSignals `mapstructure:"in_ports"`
	Outs PortToSignals `mapstructure:"out_ports"`
}

// AddInPort adds an input port to the PortMapping.
func (p *PortMapping) AddInPort(portName string, signals []Signal) {
	if p.Ins == nil {
		p.Ins = make(PortToSignals)
	}
	p.Ins[portName] = signals
}

// AddOutPort adds an output port to the PortMapping.
func (p *PortMapping) AddOutPort(portName string, signals []Signal) {
	if p.Outs == nil {
		p.Outs = make(PortToSignals)
	}
	p.Outs[portName] = signals
}

// GetInPort returns true if the port exists in the PortMapping.
func (p *PortMapping) GetInPort(portName string) ([]Signal, bool) {
	if p.Ins == nil {
		return nil, false
	}
	signals, ok := p.Ins[portName]
	return signals, ok
}

// GetOutPort returns true if the port exists in the PortMapping.
func (p *PortMapping) GetOutPort(portName string) ([]Signal, bool) {
	if p.Outs == nil {
		return nil, false
	}
	signals, ok := p.Outs[portName]
	return signals, ok
}

// NewPortMapping creates a new PortMapping.
func NewPortMapping() PortMapping {
	return PortMapping{
		Ins:  make(PortToSignals),
		Outs: make(PortToSignals),
	}
}

// Merge merges two PortMappings.
func (p *PortMapping) Merge(other PortMapping) error {
	err := p.Ins.merge(other.Ins)
	if err != nil {
		return err
	}
	err = p.Outs.merge(other.Outs)
	return err
}

// PortToSignals is a map from port name to a list of signals.
type PortToSignals map[string][]Signal

func (p PortToSignals) merge(other PortToSignals) error {
	for portName, signals := range other {
		if _, ok := p[portName]; !ok {
			p[portName] = signals
		} else {
			return errors.New("duplicate port definition")
		}
	}
	return nil
}

// ConstantSignal is a mirror struct to same proto message.
type ConstantSignal struct {
	SpecialValue string  `mapstructure:"special_value"`
	Value        float64 `mapstructure:"value"`
}

// Description returns a description of the constant signal.
func (constantSignal *ConstantSignal) Description() string {
	specialValue := constantSignal.SpecialValue
	value := constantSignal.Value
	var description string
	if specialValue != "" {
		description = specialValue
	} else {
		description = fmt.Sprintf("%0.2f", value)
	}

	return description
}

// Float returns the float value of the constant signal.
func (constantSignal *ConstantSignal) Float() float64 {
	specialValue := constantSignal.SpecialValue
	value := constantSignal.Value
	if specialValue == "NaN" {
		return math.NaN()
	}
	if specialValue == "+Inf" {
		return math.Inf(1)
	}
	if specialValue == "-Inf" {
		return math.Inf(-1)
	}
	return value
}

// ConstantSignalFromProto creates a ConstantSignal from a proto message.
func ConstantSignalFromProto(constantSignalProto *policylangv1.ConstantSignal) *ConstantSignal {
	return &ConstantSignal{
		SpecialValue: constantSignalProto.GetSpecialValue(),
		Value:        constantSignalProto.GetValue(),
	}
}

// PortsFromComponentConfig extracts Ports from component's config.
func PortsFromComponentConfig(componentConfig mapstruct.Object, subCircuitID string) (PortMapping, error) {
	var ports PortMapping

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true, // So that singular ports will transparently be converted to lists.
		Result:           &ports,
	})
	if err != nil {
		return PortMapping{}, err
	}

	err = decoder.Decode(componentConfig)
	// Add circuitID to all signals.
	for _, signals := range ports.Ins {
		for i := range signals {
			signals[i].SubCircuitID = subCircuitID
		}
	}
	for _, signals := range ports.Outs {
		for i := range signals {
			signals[i].SubCircuitID = subCircuitID
		}
	}
	return ports, err
}

// SignalType enum.
type SignalType int

// MakeRootSignalID creates SignalID with "root" SubCircuitID.
func MakeRootSignalID(signalName string) SignalID {
	return SignalID{
		SubCircuitID: "root",
		SignalName:   signalName,
	}
}

// SignalID is a unique identifier for a signal.
type SignalID struct {
	SubCircuitID string
	SignalName   string
}

const (
	// SignalTypeNamed is a named signal.
	SignalTypeNamed = iota
	// SignalTypeConstant is a constant signal.
	SignalTypeConstant
)

// Signal describes an input or output port of a component
//
// Only one field should be set.
type Signal struct {
	SubCircuitID   string
	SignalName     string         `mapstructure:"signal_name"`
	ConstantSignal ConstantSignal `mapstructure:"constant_signal"`
	Looped         bool
}

// SignalID returns the Signal ID.
func (s *Signal) SignalID() SignalID {
	return SignalID{
		SubCircuitID: s.SubCircuitID,
		SignalName:   s.SignalName,
	}
}

// SignalType returns the Signal type of the port.
func (s *Signal) SignalType() SignalType {
	if s.SignalName != "" {
		return SignalTypeNamed
	}
	return SignalTypeConstant
}

// ConstantSignalValue returns the value of the constant signal.
func (s *Signal) ConstantSignalValue() float64 {
	constantSignal := s.ConstantSignal
	value := constantSignal.Float()

	return value
}
