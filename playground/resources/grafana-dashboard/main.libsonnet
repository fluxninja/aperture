local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';
local dashboard = grafanaOperator.integreatly.v1alpha1.grafanaDashboard;
local policyName = std.extVar('POLICY_NAME');

function(dashboardMixin) {
  dashboards: [
    dashboard.new('%s-dashboard' % policyName)
    + dashboard.metadata.withNamespace('aperture-controller')
    + dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' })
    + dashboard.spec.withJson(std.manifestJsonEx(dashboardMixin, indent='  '))
    + dashboard.spec.withDatasources({
      inputName: 'DS_CONTROLLER-PROMETHEUS',
      datasourceName: 'controller-prometheus',
    }),
  ],
}
