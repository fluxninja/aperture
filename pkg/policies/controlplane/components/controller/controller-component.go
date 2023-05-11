package controller

import (
	"math"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// ControllerComponent provides a runtime.Component for Controllers. It can be initialized with a Controller implementation. It exposes a ControllerStateReadAPI for Controllers.
type ControllerComponent struct {
	controller Controller
	// Signal's last reading
	signal runtime.Reading
	// Setpoint's last reading
	setpoint runtime.Reading
	// Control variable's last reading
	controlVariable runtime.Reading
	// Controller output's last reading
	output           runtime.Reading
	policyReadAPI    iface.Policy
	defaultConfig    *policylangv1.GradientController_DynamicConfig
	dynamicConfigKey string
	componentName    string
	shortDescription string
	componentID      string
	manualMode       bool
}

// NewControllerComponent creates a new ControllerComponent.
func NewControllerComponent(
	controller Controller,
	componentName string,
	shortDescription string,
	componentID string,
	policyReadAPI iface.Policy,
	dynamicConfigKey string,
	defaultConfig *policylangv1.GradientController_DynamicConfig,
) *ControllerComponent {
	cc := &ControllerComponent{
		signal:           runtime.InvalidReading(),
		setpoint:         runtime.InvalidReading(),
		controlVariable:  runtime.InvalidReading(),
		output:           runtime.InvalidReading(),
		controller:       controller,
		componentName:    componentName,
		shortDescription: shortDescription,
		componentID:      componentID,
		policyReadAPI:    policyReadAPI,
		dynamicConfigKey: dynamicConfigKey,
		defaultConfig:    defaultConfig,
	}
	cc.setConfig(defaultConfig)
	return cc
}

// Name implements runtime.Component.
func (cc *ControllerComponent) Name() string { return cc.componentName }

// Type implements runtime.Component.
func (cc *ControllerComponent) Type() runtime.ComponentType {
	return runtime.ComponentTypeSignalProcessor
}

// ShortDescription implements runtime.Component.
func (cc *ControllerComponent) ShortDescription() string {
	return cc.shortDescription
}

// IsActuator implements runtime.Component.
func (*ControllerComponent) IsActuator() bool { return false }

// Execute implements runtime.Component.Execute.
func (cc *ControllerComponent) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (outPortReadings runtime.PortToReading, err error) {
	retErr := func(err error) (runtime.PortToReading, error) {
		return runtime.PortToReading{
			"output": []runtime.Reading{runtime.InvalidReading()},
		}, err
	}

	signal := inPortReadings.ReadSingleReadingPort("signal")
	setpoint := inPortReadings.ReadSingleReadingPort("setpoint")
	optimize := inPortReadings.ReadSingleReadingPort("optimize")
	max := inPortReadings.ReadSingleReadingPort("max")
	min := inPortReadings.ReadSingleReadingPort("min")
	controlVariable := inPortReadings.ReadSingleReadingPort("control_variable")
	output := runtime.InvalidReading()

	prevSetpoint := cc.setpoint

	// Save readings for the current tick so that Controller may access them via ControllerStateReadAPI
	cc.signal = signal
	cc.setpoint = setpoint
	cc.controlVariable = controlVariable

	if signal.Valid() && setpoint.Valid() {
		// ComputeOutput
		computedOutput, err := cc.controller.ComputeOutput(signal, setpoint, controlVariable, cc, tickInfo)
		if err != nil {
			return retErr(err)
		}
		output = computedOutput
	}

	// Check if the setpoint has changed
	if setpoint.Valid() && setpoint.Value() != prevSetpoint.Value() {
		// Try to maintain output
		err := cc.controller.MaintainOutput(prevSetpoint, setpoint, cc, tickInfo)
		if err != nil {
			return retErr(err)
		}
	}

	// Optimize
	if output.Valid() && optimize.Valid() {
		targetOutput := runtime.NewReading(output.Value() + optimize.Value())
		// Wind output
		windedOutput, err := cc.controller.WindOutput(output, targetOutput, cc, tickInfo)
		output = windedOutput
		if err != nil {
			return retErr(err)
		}
	}

	if output.Valid() {
		outputReading := output
		// Constrain output
		if max.Valid() {
			outputReading = runtime.NewReading(math.Min(output.Value(), max.Value()))
		}
		if min.Valid() {
			outputReading = runtime.NewReading(math.Max(outputReading.Value(), min.Value()))
		}
		if outputReading.Value() != output.Value() {
			// Wind output
			windedOutput, err := cc.controller.WindOutput(output, outputReading, cc, tickInfo)
			output = windedOutput
			if err != nil {
				return retErr(err)
			}
		}
	}

	// Set output to control variable in-case of Manual mode
	if cc.manualMode {
		// wind the controller output to the control variable
		windedOutput, err := cc.controller.WindOutput(output, controlVariable, cc, tickInfo)
		if err != nil {
			return retErr(err)
		}
		output = windedOutput
	}

	// Save readings for the next tick so that Controller may access them via ControllerStateReadAPI
	cc.output = output

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate handles setting of controller.ControllerMode.
func (cc *ControllerComponent) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := cc.policyReadAPI.GetStatusRegistry().GetLogger()
	dynamicConfig := &policylangv1.GradientController_DynamicConfig{}
	if unmarshaller.IsSet(cc.dynamicConfigKey) {
		err := unmarshaller.UnmarshalKey(cc.dynamicConfigKey, dynamicConfig)
		if err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal dynamic config")
			return
		}
		cc.setConfig(dynamicConfig)
	} else {
		cc.setConfig(cc.defaultConfig)
	}
}

func (cc *ControllerComponent) setConfig(config *policylangv1.GradientController_DynamicConfig) {
	if config != nil {
		cc.manualMode = config.ManualMode
	} else {
		cc.manualMode = false
	}
}

// GetSignal returns the signal's last reading.
func (cc *ControllerComponent) GetSignal() runtime.Reading {
	return cc.signal
}

// GetSetpoint returns the setpoint's last reading.
func (cc *ControllerComponent) GetSetpoint() runtime.Reading {
	return cc.setpoint
}

// GetControlVariable returns the control variable's last reading.
func (cc *ControllerComponent) GetControlVariable() runtime.Reading {
	return cc.controlVariable
}

// GetControllerOutput returns the controller output's last reading.
func (cc *ControllerComponent) GetControllerOutput() runtime.Reading {
	return cc.output
}
