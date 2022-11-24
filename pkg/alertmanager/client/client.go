package alertmgrclient

import (
	"context"
	"net/http"
	"net/url"

	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	promclient "github.com/prometheus/alertmanager/api/v2/client"
	promalert "github.com/prometheus/alertmanager/api/v2/client/alert"
	prommodels "github.com/prometheus/alertmanager/api/v2/models"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	commonhttp "github.com/fluxninja/aperture/pkg/net/http"
)

var (
	// AlertMgrClientConfigKey is the key used to store the AlertManagerClientConfig in the config.
	configKey           = "alertmanagers"
	AlertMgrClientFxKey = config.GroupTag("alertmgr_clients")
)

// AlertManagerConfig main level config for alertmanager.
type AlertManagerConfig struct {
	Clients []AlertManagerClientConfig `json:"clients" validate:"required"`
}

// AlertManagerClientConfig config for single alertmanager client.
type AlertManagerClientConfig struct {
	Name             string                      `json:"name" validate:"required"`
	Address          string                      `json:"address" validate:"required,hostname_port|url|fqdn"`
	HttpClientConfig commonhttp.HTTPClientConfig `json:"http_client"`
}

func ProvideNamedAlertManagerClients(unmarshaller config.Unmarshaller) []AlertManagerClient {
	clientSlice := []AlertManagerClient{}

	var config AlertManagerConfig
	if err := unmarshaller.UnmarshalKey(configKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize alert managers configuration!")
		return clientSlice
	}

	for _, configItem := range config.Clients {
		httpClient, err := commonhttp.ClientFromConfig(configItem.HttpClientConfig)
		if err != nil {
			log.Warn().Msg("Could not create http client from config")
			continue
		}
		amClient := CreateClient(configItem.Name, configItem.Address, httpClient)
		clientSlice = append(clientSlice, amClient)
	}
	return clientSlice
}

// AlertManagerClient provides an interface for alert manager client.
type AlertManagerClient interface {
	SendAlert(ctx context.Context, alerts prommodels.PostableAlerts) error
}

// RealAlertManagerClient implements AlertManagerClient interface.
type RealAlertManagerClient struct {
	name            string
	httpClient      *http.Client
	promAlertClient *promclient.Alertmanager
}

// CreateClient creates a new alertmanager client with provided http client.
func CreateClient(name, address string, httpClient *http.Client) AlertManagerClient {
	hu, _ := url.Parse(address)
	transport := runtimeclient.NewWithClient(hu.Host, "/", []string{"http"}, httpClient)
	promClient := promclient.New(transport, strfmt.NewFormats())

	alertMgrClient := &RealAlertManagerClient{
		name:            name,
		promAlertClient: promClient,
		httpClient:      httpClient,
	}
	return alertMgrClient
}

func (ac *RealAlertManagerClient) SendAlert(ctx context.Context, alerts prommodels.PostableAlerts) error {
	postAlertParams := &promalert.PostAlertsParams{
		Context:    ctx,
		HTTPClient: ac.httpClient,
		Alerts:     alerts,
	}
	_, err := ac.promAlertClient.Alert.PostAlerts(postAlertParams)

	return err
}
