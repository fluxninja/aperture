package alertmanager

import (
	"go.uber.org/fx"

	amclient "github.com/fluxninja/aperture/pkg/alertmanager/client"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
)

const (
	configKey = "alertmanagers"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideNamedAlertManagerClients,
			fx.Annotate(
				ProvideAlertManager,
				fx.ParamTags(amclient.AlertMgrClientFxKey),
			)),
	)
}

type AlertManagerConfig struct {
	Name             string                       `json:"name" validate:"required"`
	Address          string                       `json:"address" validate:"required,hostname_port|url|fqdn"`
	HttpClientConfig *httpclient.HTTPClientConfig `json:"http_client"`
}

func ProvideNamedAlertManagerClients(unmarshaller config.Unmarshaller) fx.Option {
	var configList []*AlertManagerConfig
	if err := unmarshaller.UnmarshalKey(configKey, &configList); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize alert managers configuration!")
		return fx.Options()
	}
	log.Bug().Msgf("DARIA LOG ALERT MANAGER CONFIG LIST: %+v", configList)

	var optionList []fx.Option
	for _, configItem := range configList {
		log.Warn().Msgf("DARIA LOG ALERT CLIENT NAME: %+v", configItem.Name)
		options := amclient.ProvideNamedClient(configItem.Name, configItem.HttpClientConfig)
		optionList = append(optionList, options)
	}

	return fx.Options(
		optionList...,
	)
}

type AlertManager struct {
	Clients []amclient.AlertManagerClient
}

func ProvideAlertManager(clientSlice ...amclient.AlertManagerClient) *AlertManager {

	log.Bug().Msgf("DARIA LOG GOT CLIENT SLICE: %+v", clientSlice)

	return &AlertManager{
		Clients: clientSlice,
	}
}
