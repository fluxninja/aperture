package controller

import (
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// Controller is the interface for the controller.
type Controller interface {
	// ComputeOutput: Compute the output given the current and previous signal readings. Refer to last values of other readings via ControlLoopReadAPI.
	ComputeOutput(signal, setpoint, controlVariable runtime.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) (runtime.Reading, error)
	// WindOutput: Wind the output given the previous and target output readings. Refer to last values of other readings via ControlLoopReadAPI.
	WindOutput(currentOutput, targetOutput runtime.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) (runtime.Reading, error)
	// MaintainOutput: Try to maintain the output on setpoint change. Refer to last values of other readings via ControlLoopReadAPI.
	MaintainOutput(prevSetpoint, currentSetpoint runtime.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) error
}
