package runtime

import (
	"math"

	"github.com/fluxninja/aperture/pkg/mapstruct"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
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
	p.Ins[portName] = signals
}

// AddOutPort adds an output port to the PortMapping.
func (p *PortMapping) AddOutPort(portName string, signals []Signal) {
	p.Outs[portName] = signals
}

// ExistsInPort returns true if the port exists in the PortMapping.
func (p *PortMapping) ExistsInPort(portName string) bool {
	_, ok := p.Ins[portName]
	return ok
}

// ExistsOutPort returns true if the port exists in the PortMapping.
func (p *PortMapping) ExistsOutPort(portName string) bool {
	_, ok := p.Outs[portName]
	return ok
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

// PortToSignals is a map from port name to a list of ports.
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

// PortsFromComponentConfig extracts Ports from component's config.
func PortsFromComponentConfig(componentConfig mapstruct.Object, circuitID string) (PortMapping, error) {
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
			signals[i].CircuitID = circuitID
		}
	}
	for _, signals := range ports.Outs {
		for i := range signals {
			signals[i].CircuitID = circuitID
		}
	}
	return ports, err
}

// SignalType enum.
type SignalType int

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
	// Note: pointers are used to detect fields being not set.
	SignalName     string `mapstructure:"signal_name"`
	CircuitID      string
	ConstantSignal ConstantSignal `mapstructure:"constant_signal"`
	Looped         bool
}

// SignalType returns the Signal type of the port.
func (s *Signal) SignalType() SignalType {
	if s.SignalName != "" {
		return SignalTypeNamed
	}
	return SignalTypeConstant
}

// GetConstantSignalValue returns the value of the constant signal.
func (s *Signal) GetConstantSignalValue() float64 {
	constantSignal := s.ConstantSignal
	value := 0.0
	specialValue := constantSignal.SpecialValue
	if specialValue != "" {
		switch specialValue {
		case "NaN":
			value = math.NaN()
		case "+Inf":
			value = math.Inf(1)
		case "-Inf":
			value = math.Inf(-1)
		}
	} else {
		value = constantSignal.Value
	}

	return value
}
