local apertureGrafanaApp = import 'apps/aperture-grafana/main.libsonnet';
local grafanaOperatorApp = import 'apps/grafana-operator/main.libsonnet';

local grafanaOperatorMixin =
  grafanaOperatorApp {
    environment+:: {
      namespace: 'default',
      name: 'aperture-grafana-operator',
    },
    values+:: {
      grafana+: {
        enabled: false,
      },
    },
  };

{
  grafanaOperator: grafanaOperatorMixin,
  apertureGrafana: apertureGrafanaApp,
}
