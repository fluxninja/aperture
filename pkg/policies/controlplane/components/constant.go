package components

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Variable is a dynamically configurable variable signal.
type Variable struct {
	constantValue *policylangv1.ConstantValue
	policyReadAPI iface.Policy
	variableProto *policylangv1.Variable
}

// Name implements runtime.Component.
func (*Variable) Name() string { return "Variable" }

// Type implements runtime.Component.
func (*Variable) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// NewVariable creates a variable component.
func NewVariable(value float64) runtime.Component {
	return &Variable{
		constantValue: &policylangv1.ConstantValue{
			Valid: true,
			Value: value,
		},
	}
}

// NewVariableAndOptions creates a variable components and its fx options.
func NewVariableAndOptions(variableProto *policylangv1.Variable, componentIndex int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	variable := &Variable{
		policyReadAPI: policyReadAPI,
		constantValue: variableProto.DefaultConfig.ConstantValue,
		variableProto: variableProto,
	}

	return variable, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (v *Variable) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	// Always emit the value.
	if !v.constantValue.Valid {
		return runtime.PortToValue{
			"output": []runtime.Reading{runtime.InvalidReading()},
		}, nil
	}

	return runtime.PortToValue{
		"output": []runtime.Reading{runtime.NewReading(v.constantValue.Value)},
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
		v.constantValue = dynamicConfig.ConstantValue
	} else {
		v.constantValue = v.variableProto.GetDefaultConfig().ConstantValue
	}
}
