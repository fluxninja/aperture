package components

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Variable is a dynamically configurable variable signal.
type Variable struct {
	constantSignal *policylangv1.ConstantSignal
	policyReadAPI  iface.Policy
	variableProto  *policylangv1.Variable
}

// Name implements runtime.Component.
func (*Variable) Name() string { return "Variable" }

// Type implements runtime.Component.
func (*Variable) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// NewConstantSignal creates a variable component with a value that's always valid.
func NewConstantSignal(value float64) runtime.Component {
	return &Variable{
		constantSignal: &policylangv1.ConstantSignal{
			Const: &policylangv1.ConstantSignal_Value{Value: value},
		},
	}
}

// NewVariableAndOptions creates a variable components and its fx options.
func NewVariableAndOptions(variableProto *policylangv1.Variable, _ string, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	variable := &Variable{
		policyReadAPI:  policyReadAPI,
		constantSignal: variableProto.DefaultConfig.ConstantSignal,
		variableProto:  variableProto,
	}

	return variable, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (v *Variable) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	// Always emit the value.
	if specialValue := v.constantSignal.GetSpecialValue(); specialValue != "" {
		output := runtime.InvalidReading()
		switch specialValue {
		case "NaN":
			output = runtime.NewReading(math.NaN())
		case "+Inf":
			output = runtime.NewReading(math.Inf(1))
		case "-Inf":
			output = runtime.NewReading(math.Inf(-1))
		}
		return runtime.PortToReading{
			"output": []runtime.Reading{output},
		}, nil
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{runtime.NewReading(v.constantSignal.GetValue())},
	}, nil
}

// DynamicConfigUpdate finds the dynamic config and syncs the constant value.
func (v *Variable) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := v.policyReadAPI.GetStatusRegistry().GetLogger()
	key := v.variableProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		dynamicConfig := &policylangv1.Variable_DynamicConfig{}
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal dynamic config")
			return
		}
		v.constantSignal = dynamicConfig.ConstantSignal
	} else {
		v.constantSignal = v.variableProto.GetDefaultConfig().ConstantSignal
	}
}
