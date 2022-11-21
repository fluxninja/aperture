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

// Module is a fx module that constructs annotated instance of alerts.Alerter.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				ProvideAlerter,
				fx.ResultTags(AlertsFxTag),
			),
		),
	)
}

// ProvideAlerter creates an alerter.
func ProvideAlerter(unmarshaller config.Unmarshaller) (Alerter, error) {
	var cfg AlerterConfig
	if err := unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize alerter configuration!")
		return nil, err
	}

	return NewSimpleAlerter(cfg.ChannelSize), nil
}
