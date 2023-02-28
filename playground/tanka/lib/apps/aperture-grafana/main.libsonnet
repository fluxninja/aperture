local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';
local kubernetesMixin = import 'github.com/kubernetes-monitoring/kubernetes-mixin/mixin.libsonnet';

local aperture = import '../../../../../blueprints/main.libsonnet';
local policyDashboard = aperture.policies.LatencyAIMDConcurrencyLimiting.dashboard;
local rateLimitpolicyDashboard = aperture.policies.StaticRateLimiting.dashboard;
local signalsDashboard = aperture.dashboards.Signals.dashboard;

local grafana = grafanaOperator.integreatly.v1alpha1.grafana;
local dashboard = grafanaOperator.integreatly.v1alpha1.grafanaDashboard;
local dataSource = grafanaOperator.integreatly.v1alpha1.grafanaDataSource;

local dataSources =
  {
    controllerPrometheus:
      dataSource.new('controller-prometheus') +
      dataSource.spec.withName('controller-prometheus') +
      dataSource.spec.withDatasources({
        name: 'controller-prometheus',
        type: 'prometheus',
        access: 'proxy',
        url: 'http://controller-prometheus-server',
        jsonData: {
          timeInterval: '1s',
        },
      }),
    operationsPrometheus:
      dataSource.new('operations-prometheus') +
      dataSource.spec.withName('operations-prometheus') +
      dataSource.spec.withDatasources({
        name: 'operations-prometheus',
        type: 'prometheus',
        access: 'proxy',
        url: 'http://prometheus-operations.monitoring:9090',
      }),
  };

local kubeDashboards =
  (kubernetesMixin {
     _config+: {
       cadvisorSelector: 'job="kubelet", metrics_path="/metrics/cadvisor"',
     },
   }).grafanaDashboards;

local dashboards =
  [
    dashboard.new('k8s-resources') +
    dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' }) +
    dashboard.spec.withJson(std.manifestJsonEx(kubeDashboards['k8s-resources-pod.json'], indent='  ')) +
    dashboard.spec.withDatasources({
      inputName: '${datasource}',
      datasourceName: 'operations-prometheus',
    }),

    dashboard.new('aperture-signals')
    + dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' })
    + dashboard.spec.withJson(std.manifestJsonEx(signalsDashboard({
      policy_name: 'service1-demo-app',
      datasource+: {
        name: 'controller-prometheus',
      },
    }).dashboard, indent='  ', newline='\n'))
    + dashboard.spec.withDatasources({
      inputName: 'datasource',
      datasourceName: 'controller-prometheus',
    }),
  ];

local grafanaMixin =
  grafana.new('aperture-grafana')
  + grafana.spec.withDashboardLabelSelector({ 'fluxninja.com/grafana-instance': 'aperture-grafana' })
  + grafana.spec.config.security.withAdmin_user('fluxninja')
  + grafana.spec.config.security.withAdmin_password('fluxninja')
  + grafana.spec.service.withName('aperture-grafana')
  + {
    spec+: {
      config+: {
        'auth.anonymous': {
          enabled: true,
          org_name: 'Main Org.',
          org_role: 'Admin',
        },
      },
    },
  };

{
  grafana: grafanaMixin,
  dashboards: dashboards,
  datasources: dataSources,
}
