package prometheus

import (
	plp "github.com/prometheus-community/prom-label-proxy/injectproxy"
	promlabels "github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	promconfig "github.com/fluxninja/aperture/v2/pkg/prometheus/config"
)

// EnforcerIn holds fields, parameters, to provide Prometheus Label Injector .
type EnforcerIn struct {
	fx.In
	Unmarshaller config.Unmarshaller
}

// PrometheusEnforcer is responsible for enforcing extra set of labels that
// should be present in every PromQL query executed.
type PrometheusEnforcer struct {
	labels []*promlabels.Matcher
}

// EnforceLabels transforms given query, making sure that all the required
// labels are in place.
func (e *PrometheusEnforcer) EnforceLabels(query string) (string, error) {
	if len(e.labels) == 0 {
		return query, nil
	}

	expr, err := parser.ParseExpr(query)
	if err != nil {
		return "", err
	}

	enforcer := plp.NewEnforcer(false, e.labels...)

	if err := enforcer.EnforceNode(expr); err != nil {
		return "", err
	}

	return expr.String(), nil
}

func providePrometheusEnforcer(in EnforcerIn) (*PrometheusEnforcer, error) {
	var config promconfig.PrometheusConfig
	if err := in.Unmarshaller.UnmarshalKey(prometheusConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("unable to deserialize")
		return nil, err
	}

	labels := []*promlabels.Matcher{}
	for _, label := range config.Labels {
		labels = append(labels, &promlabels.Matcher{
			Name:  label.Name,
			Type:  promlabels.MatchEqual,
			Value: label.Value,
		})
	}

	return &PrometheusEnforcer{labels: labels}, nil
}
