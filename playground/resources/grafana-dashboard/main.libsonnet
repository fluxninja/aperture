local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';
local dashboard = grafanaOperator.integreatly.v1alpha1.grafanaDashboard;
local dashboardName = std.extVar('DASHBOARD_NAME');

[
  dashboard.new('%s-dashboard' % dashboardName)
  + dashboard.metadata.withNamespace('aperture-controller')
  + dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' })
  + dashboard.spec.withJson(std.manifestJsonEx(std.parseJson(std.extVar('APERTURE_DASHBOARD')), indent='  '))
  + dashboard.spec.withDatasources({
    inputName: 'DS_CONTROLLER-PROMETHEUS',
    datasourceName: 'controller-prometheus',
  }),
]
