package runtime

import (
	"fmt"

	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// ConfiguredComponent consists of a Component and its PortMapping.
type ConfiguredComponent struct {
	Component
	// Which signals this component wants to have connected on its ports.
	PortMapping PortMapping
	// Mapstruct representation of proto config that was used to create this
	// component.  This Config is used only for observability purposes.
	//
	// Note: PortMapping is also part of Config.
	Config utils.MapStruct
	// ComponentID is the unique ID of this component within the circuit.
	ComponentID ComponentID
}

// NewConfiguredComponent creates a new ConfiguredComponent.
func NewConfiguredComponent(
	component Component,
	config any,
	componentID ComponentID,
	doParsePortMapping bool,
) (*ConfiguredComponent, error) {
	subCircuitID, ok := componentID.ParentID()
	if !ok {
		return nil, fmt.Errorf("component %s is not in a circuit", componentID.String())
	}

	mapStruct, err := utils.ToMapStruct(config)
	if err != nil {
		return nil, err
	}

	ports := NewPortMapping()
	if doParsePortMapping {
		ports, err = PortsFromComponentConfig(mapStruct, subCircuitID.String())
		if err != nil {
			return nil, err
		}
	}

	return &ConfiguredComponent{
		Component:   component,
		PortMapping: ports,
		Config:      mapStruct,
		ComponentID: componentID,
	}, nil
}
