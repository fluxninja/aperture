package runtime

import "github.com/mitchellh/mapstructure"

// Ports is description of a component's ports mapping.
//
// This struct is meant to be deserializable from map-struct serialized
// representation of _any_ of the components. Eg. EMA component defines:
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
// And such EMA component could be serialized and deserialized into Ports as:
//
// ```go
//
//	Ports {
//	  InPorts: map[string]InPort {
//	    "input": []InPort {{ ... }},
//	  },
//	}
//
// ```
//
// Note how "input" is a concrete field in EMA definition, but a dynamic map
// key in Ports.
type Ports struct {
	// Note: Not using policylangv1.InPort and OutPort directly to avoid
	// runtime depending on proto.
	InPorts  map[string][]Port `mapstructure:"in_ports"`
	OutPorts map[string][]Port `mapstructure:"out_ports"`
}

// PortsFromMapStruct extracts Ports from component serialized previously to
// MapStruct via encodeMapStruct.
//
// Note: This relies on every proto structure providing Marshal/UnmarshalJSON
// (via protojson with protojson.MarshalOptions.UseProtoNames).
func PortsFromMapStruct(componentMapStruct map[string]any) (Ports, error) {
	var ports Ports

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true, // So that singular ports will transparently be converted to lists.
		Result:           &ports,
	})
	if err != nil {
		return Ports{}, err
	}

	err = decoder.Decode(componentMapStruct)
	return ports, err
}

// Port describes an input or output port of a component
//
// Only one field should be set.
type Port struct {
	// Note: pointers are used to detect fields being not set.
	SignalName    *string  `mapstructure:"signal_name"`
	ConstantValue *float64 `mapstructure:"constant_value"`
}
