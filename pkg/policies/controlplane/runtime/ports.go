package runtime

import (
	"github.com/fluxninja/aperture/pkg/mapstruct"
	"github.com/mitchellh/mapstructure"
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
	Ins  map[string][]Port `mapstructure:"in_ports"`
	Outs map[string][]Port `mapstructure:"out_ports"`
}

// ConstantSignal is a mirror struct to same proto message.
type ConstantSignal struct {
	SpecialValue *string  `mapstructure:"special_value"`
	Value        *float64 `mapstructure:"value"`
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

// Port describes an input or output port of a component
//
// Only one field should be set.
type Port struct {
	// Note: pointers are used to detect fields being not set.
	SignalName     *string         `mapstructure:"signal_name"`
	ConstantSignal *ConstantSignal `mapstructure:"constant_signal"`
}
