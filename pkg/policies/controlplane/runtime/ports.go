package runtime

import (
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

// NewPortMapping creates a new PortMapping.
func NewPortMapping() PortMapping {
	return PortMapping{
		Ins:  make(PortToSignals),
		Outs: make(PortToSignals),
	}
}

// PortsFromComponentConfig extracts Ports from component's config.
func PortsFromComponentConfig(componentConfig mapstruct.Object) (PortMapping, error) {
	var ports PortMapping

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true, // So that singular ports will transparently be converted to lists.
		Result:           &ports,
	})
	if err != nil {
		return PortMapping{}, err
	}

	err = decoder.Decode(componentConfig)
	return ports, err
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
	SignalName    string `mapstructure:"signal_name"`
	CircuitID     string
	ConstantValue float64 `mapstructure:"constant_value"`
	Looped        bool
}

// SignalType returns the Signal type of the port.
func (p *Signal) SignalType() SignalType {
	if p.SignalName != "" {
		return SignalTypeNamed
	}
	return SignalTypeConstant
}

// MakeNamedSignal creates a new named Signal.
func MakeNamedSignal(name string, looped bool) Signal {
	return Signal{
		SignalName: name,
		Looped:     looped,
	}
}

// MakeConstantSignal creates a new constant Signal.
func MakeConstantSignal(value float64) Signal {
	return Signal{
		ConstantValue: value,
	}
}
