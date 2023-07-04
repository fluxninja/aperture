package prometheus

import (
	"errors"
	"net/http"
	"strings"

	promapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"go.uber.org/fx"
	"golang.org/x/oauth2"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	commonhttp "github.com/fluxninja/aperture/v2/pkg/net/http"
	promconfig "github.com/fluxninja/aperture/v2/pkg/prometheus/config"
)

var (
	// swagger:operation POST /prometheus common-configuration Prometheus
	// ---
	// x-fn-config-env: true
	// parameters:
	// - in: body
	//   schema:
	//     $ref: "#/definitions/PrometheusConfig"
	// - name: http_client
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/HTTPClientConfig"

	// PrometheusConfigKey is the key used to store the PrometheusConfig in the config.
	prometheusConfigKey = "prometheus"
	// HttpConfigKey is the key used to store the HTTPClientConfig in the config.
	httpConfigKey = strings.Join([]string{prometheusConfigKey, "http_client"}, ".")
)

// Module provides a singleton pointer to prometheusv1.API via FX.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(providePrometheusClient),
		fx.Provide(providePrometheusEnforcer),
		commonhttp.ClientConstructor{Name: "prometheus.http-client", ConfigKey: httpConfigKey}.Annotate(),
	)
}

// ClientIn holds fields, parameters, to provide Prometheus Client.
type ClientIn struct {
	fx.In
	HTTPClient   *http.Client `name:"prometheus.http-client"`
	TokenSource  oauth2.TokenSource
	Unmarshaller config.Unmarshaller
}

func providePrometheusClient(in ClientIn) (prometheusv1.API, promapi.Client, error) {
	var config promconfig.PrometheusConfig
	if err := in.Unmarshaller.UnmarshalKey(prometheusConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("unable to deserialize")
		return nil, nil, err
	}

	if config.Address == "" {
		err := errors.New("prometheus address not specified")
		log.Error().Err(err).Msg("")
		return nil, nil, err
	}

	if in.TokenSource != nil {
		log.Info().Msg("Using Google TokenSource for prometheus API queries")
		oauth2Transport := &oauth2.Transport{
			Source: in.TokenSource,
			Base:   in.HTTPClient.Transport,
		}
		in.HTTPClient.Transport = oauth2Transport
	}

	client, err := promapi.NewClient(promapi.Config{
		Address:      config.Address,
		RoundTripper: in.HTTPClient.Transport,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error creating client")
		return nil, nil, err
	}

	return prometheusv1.NewAPI(client), client, nil
}
