package controller

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/constraints"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
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
	output         runtime.Reading
	policyReadAPI  iface.Policy
	componentIndex int
}

// NewControllerComponent creates a new ControllerComponent.
func NewControllerComponent(controller Controller, componentIndex int, policyReadAPI iface.Policy) *ControllerComponent {
	return &ControllerComponent{
		signal:          runtime.InvalidReading(),
		setpoint:        runtime.InvalidReading(),
		controlVariable: runtime.InvalidReading(),
		output:          runtime.InvalidReading(),
		controller:      controller,
		componentIndex:  componentIndex,
		policyReadAPI:   policyReadAPI,
	}
}

// Execute implements runtime.Component.Execute.
func (cc *ControllerComponent) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (outPortReadings runtime.PortToValue, err error) {
	retErr := func(err error) (runtime.PortToValue, error) {
		return runtime.PortToValue{
			"output": []runtime.Reading{runtime.InvalidReading()},
		}, err
	}

	signal := inPortReadings.ReadSingleValuePort("signal")
	setpoint := inPortReadings.ReadSingleValuePort("setpoint")
	optimize := inPortReadings.ReadSingleValuePort("optimize")
	max := inPortReadings.ReadSingleValuePort("max")
	min := inPortReadings.ReadSingleValuePort("min")
	controlVariable := inPortReadings.ReadSingleValuePort("control_variable")
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

	// Constraints
	minMaxConstraints := constraints.NewMinMaxConstraints()
	if max.Valid() {
		// minxMaxConstraints' Max, Min are set as math.MaxFloat64, -math.MaxFloat64 initially; no error.
		err := minMaxConstraints.SetMax(max.Value())
		if err != nil {
			return retErr(err)
		}
	}
	if min.Valid() {
		err := minMaxConstraints.SetMin(min.Value())
		if err != nil {
			// To make sure min is less than max; otherwise, emits invalid signal.
			return retErr(err)
		}
	}

	if output.Valid() {
		// Constrain output
		outputConstrained, _ := minMaxConstraints.Constrain(output.Value())
		outputReading := runtime.NewReading(outputConstrained)
		if outputReading.Value() != output.Value() {
			// Wind output
			windedOutput, err := cc.controller.WindOutput(output, outputReading, cc, tickInfo)
			output = windedOutput
			if err != nil {
				return retErr(err)
			}
		}
	}

	// Save readings for the next tick so that Controller may access them via ControllerStateReadAPI
	cc.output = output

	return runtime.PortToValue{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for ControllerComponent.
func (cc *ControllerComponent) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
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
