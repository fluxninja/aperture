local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';

local aperture = import '../../../blueprints/main.libsonnet';
local signalsDashboard = aperture.dashboards.Signals.dashboard;

local dashboard = grafanaOperator.integreatly.v1alpha1.grafanaDashboard;
local policyName = std.extVar('POLICY_NAME');

function(dashboardMixin) {
  dashboards:
    [
      dashboard.new('aperture-signals-%s' % policyName)
      + dashboard.metadata.withNamespace('aperture-controller')
      + dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' })
      + dashboard.spec.withJson(std.manifestJsonEx(signalsDashboard({
        policy+: {
          policy_name: policyName,
        },
        dashboard+: {
          datasource+: {
            name: 'controller-prometheus',
          },
        },
      }).dashboard, indent='  ', newline='\n'))
      + dashboard.spec.withDatasources({
        inputName: 'datasource',
        datasourceName: 'controller-prometheus',
      }),

      dashboard.new('%s-dashboard' % policyName) +
      dashboard.metadata.withNamespace('aperture-controller') +
      dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' }) +
      dashboard.spec.withJson(std.manifestJsonEx(dashboardMixin, indent='  ')) +
      dashboard.spec.withDatasources({
        inputName: 'DS_CONTROLLER-PROMETHEUS',
        datasourceName: 'controller-prometheus',
      }),
    ],
}
