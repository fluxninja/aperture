package alerts

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

const configKey = "alerter"

// AlerterConfig for alerter.
// swagger:model
// +kubebuilder:object:generate=true
type AlerterConfig struct {
	// ChannelSize size of the alerts channel in the alerter. Alerts should be
	// consument from it quickly, so no big sizes are needed.
	ChannelSize int `json:"channel_size" validate:"gt=0" default:"100"`
}

// FxIn describes parameters passed to alerter constructor.
type FxIn struct {
	fx.In
	Unmarshaller config.Unmarshaller
}

// ProvideAlerter creates an alerter.
func ProvideAlerter(in FxIn) (Alerter, error) {
	var cfg AlerterConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize alerter configuration!")
		return nil, err
	}

	return NewSimpleAlerter(cfg.ChannelSize), nil
}
