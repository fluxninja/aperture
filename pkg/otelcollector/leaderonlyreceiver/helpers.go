package leaderonlyreceiver

import (
	"go.opentelemetry.io/collector/component"
)

// WrapConfig returns configuration of leader-only receiver wrapping given receiver.
func WrapConfig(id component.ID, cfg any) (component.ID, any) {
	return component.NewIDWithName(type_, id.String()), map[string]any{
		"type":   id.Type(),
		"config": cfg,
	}
}

// WrapConfigIf returns configuration of leader-only receiver wrapping given
// receiver if condition is true. Otherwise, returns original id and config.
func WrapConfigIf(condition bool, id component.ID, cfg any) (component.ID, any) {
	if condition {
		return WrapConfig(id, cfg)
	} else {
		return id, cfg
	}
}
