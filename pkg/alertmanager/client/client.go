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

	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	commonhttp "github.com/fluxninja/aperture/pkg/net/http"
)

var configKey = "alertmanagers"

// AlertManagerConfig main level config for alertmanager.
// swagger:model
// +kubebuilder:object:generate=true
type AlertManagerConfig struct {
	Clients []AlertManagerClientConfig `json:"clients,omitempty"`
}

// AlertManagerClientConfig config for single alertmanager client.
// swagger:model AlertManagerClientConfig
// +kubebuilder:object:generate=true
type AlertManagerClientConfig struct {
	Name       string                      `json:"name"`
	Address    string                      `json:"address" validate:"hostname_port|url|fqdn"`
	HTTPConfig commonhttp.HTTPClientConfig `json:"http_client"`
	BasePath   string                      `json:"base_path" default:"/"`
}

// ProvideNamedAlertManagerClients provides a list of alertmanager clients from configuration.
func ProvideNamedAlertManagerClients(unmarshaller config.Unmarshaller) []AlertManagerClient {
	clientSlice := []AlertManagerClient{}

	var config AlertManagerConfig
	if err := unmarshaller.UnmarshalKey(configKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize alert managers configuration!")
		return clientSlice
	}

	for _, configItem := range config.Clients {
		httpClient, err := commonhttp.ClientFromConfig(configItem.HTTPConfig)
		if err != nil {
			log.Warn().Msg("Could not create http client from config")
			continue
		}
		amClient := CreateClient(configItem.Name, configItem.Address, configItem.BasePath, httpClient)
		clientSlice = append(clientSlice, amClient)
	}
	return clientSlice
}

// AlertManagerClient provides an interface for alert manager client.
type AlertManagerClient interface {
	SendAlerts(ctx context.Context, alerts []*alerts.Alert) error
	GetName() string
}

// RealAlertManagerClient implements AlertManagerClient interface.
type RealAlertManagerClient struct {
	Name            string
	httpClient      *http.Client
	promAlertClient *promclient.Alertmanager
}

// CreateClient creates a new alertmanager client with provided http client.
func CreateClient(name, address, basePath string, httpClient *http.Client) AlertManagerClient {
	hu, _ := url.Parse(address)
	transport := runtimeclient.NewWithClient(hu.Host, basePath, []string{"http"}, httpClient)
	promClient := promclient.New(transport, strfmt.NewFormats())

	alertMgrClient := &RealAlertManagerClient{
		Name:            name,
		promAlertClient: promClient,
		httpClient:      httpClient,
	}
	return alertMgrClient
}

// SendAlerts sends postable alerts via configured alertmanager http client.
func (ac *RealAlertManagerClient) SendAlerts(ctx context.Context, alerts []*alerts.Alert) error {
	postableAlerts := make([]*prommodels.PostableAlert, len(alerts))
	for i, alert := range alerts {
		postableAlert := alert.PostableAlert()
		postableAlerts[i] = &postableAlert
	}
	postAlertParams := &promalert.PostAlertsParams{
		Context:    ctx,
		HTTPClient: ac.httpClient,
		Alerts:     prommodels.PostableAlerts(postableAlerts),
	}
	_, err := ac.promAlertClient.Alert.PostAlerts(postAlertParams)

	return err
}

// GetName getter func for alert manager client name.
func (ac *RealAlertManagerClient) GetName() string {
	return ac.Name
}
