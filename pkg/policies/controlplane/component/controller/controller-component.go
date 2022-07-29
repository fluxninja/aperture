package controller

import (
	"github.com/fluxninja/aperture/pkg/policies/apis/policyapi"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/constraints"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/reading"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// ControllerComponent provides a runtime.Component for Controllers. It can be initialized with a Controller implementation. It exposes a ControllerStateReadAPI for Controllers.
type ControllerComponent struct {
	controller Controller
	// Signal's last reading
	signal reading.Reading
	// Setpoint's last reading
	setpoint reading.Reading
	// Control variable's last reading
	controlVariable reading.Reading
	// Controller output's last reading
	output         reading.Reading
	componentIndex int
	policyReadAPI  policyapi.PolicyReadAPI
}

// NewControllerComponent creates a new ControllerComponent.
func NewControllerComponent(controller Controller, componentIndex int, policyReadAPI policyapi.PolicyReadAPI) *ControllerComponent {
	return &ControllerComponent{
		signal:          reading.NewInvalid(),
		setpoint:        reading.NewInvalid(),
		controlVariable: reading.NewInvalid(),
		output:          reading.NewInvalid(),
		controller:      controller,
		componentIndex:  componentIndex,
		policyReadAPI:   policyReadAPI,
	}
}

// Execute implements runtime.Component.Execute.
func (cc *ControllerComponent) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (outPortReadings runtime.PortToValue, err error) {
	retErr := func(err error) (runtime.PortToValue, error) {
		return runtime.PortToValue{
			"output": []reading.Reading{reading.NewInvalid()},
		}, err
	}

	signal := inPortReadings.ReadSingleValuePort("signal")
	setpoint := inPortReadings.ReadSingleValuePort("setpoint")
	optimize := inPortReadings.ReadSingleValuePort("optimize")
	max := inPortReadings.ReadSingleValuePort("max")
	min := inPortReadings.ReadSingleValuePort("min")
	controlVariable := inPortReadings.ReadSingleValuePort("control_variable")
	output := reading.NewInvalid()

	prevSetpoint := cc.setpoint

	// Save readings for the current tick so that Controller may access them via ControllerStateReadAPI
	cc.signal = signal
	cc.setpoint = setpoint
	cc.controlVariable = controlVariable

	if signal.Valid && setpoint.Valid {
		// ComputeOutput
		computedOutput, err := cc.controller.ComputeOutput(signal, setpoint, controlVariable, cc, tickInfo)
		if err != nil {
			return retErr(err)
		}
		output = computedOutput
	}

	// Check if the setpoint has changed
	if setpoint.Valid && setpoint.Value != prevSetpoint.Value {
		// Try to maintain output
		err := cc.controller.MaintainOutput(prevSetpoint, setpoint, cc, tickInfo)
		if err != nil {
			return retErr(err)
		}
	}

	// Optimize
	if output.Valid && optimize.Valid {
		targetOutput := reading.New(output.Value + optimize.Value)
		// Wind output
		windedOutput, err := cc.controller.WindOutput(output, targetOutput, cc, tickInfo)
		output = windedOutput
		if err != nil {
			return retErr(err)
		}
	}

	// Constraints
	minMaxConstraints := constraints.NewMinMaxConstraints()
	if max.Valid {
		// minxMaxConstraints' Max, Min are set as math.MaxFloat64, -math.MaxFloat64 initially; no error.
		err := minMaxConstraints.SetMax(max.Value)
		if err != nil {
			return retErr(err)
		}
	}
	if min.Valid {
		err := minMaxConstraints.SetMin(min.Value)
		if err != nil {
			// To make sure min is less than max; otherwise, emits invalid signal.
			return retErr(err)
		}
	}

	if output.Valid {
		// Constrain output
		outputConstrained, _ := minMaxConstraints.Constrain(output.Value)
		outputReading := reading.New(outputConstrained)
		if outputReading.Value != output.Value {
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
		"output": []reading.Reading{output},
	}, nil
}

// GetSignal returns the signal's last reading.
func (cc *ControllerComponent) GetSignal() reading.Reading {
	return cc.signal
}

// GetSetpoint returns the setpoint's last reading.
func (cc *ControllerComponent) GetSetpoint() reading.Reading {
	return cc.setpoint
}

// GetControlVariable returns the control variable's last reading.
func (cc *ControllerComponent) GetControlVariable() reading.Reading {
	return cc.controlVariable
}

// GetControllerOutput returns the controller output's last reading.
func (cc *ControllerComponent) GetControllerOutput() reading.Reading {
	return cc.output
}
