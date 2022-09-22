local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';
local kubernetesMixin = import 'github.com/kubernetes-monitoring/kubernetes-mixin/mixin.libsonnet';

local policyDashboard = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/dashboards/latency-gradient.libsonnet';

local grafana = grafanaOperator.integreatly.v1alpha1.grafana;
local dashboard = grafanaOperator.integreatly.v1alpha1.grafanaDashboard;
local dataSource = grafanaOperator.integreatly.v1alpha1.grafanaDataSource;

local dataSources =
  {
    cloudPrometheus:
      dataSource.new('controller-prometheus') +
      dataSource.spec.withName('controller-prometheus') +
      dataSource.spec.withDatasources({
        name: 'controller-prometheus',
        type: 'prometheus',
        access: 'proxy',
        url: 'http://controller-prometheus-server',
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
    dashboard.new('example-dashboard') +
    dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' }) +
    dashboard.spec.withJson(std.manifestJsonEx(policyDashboard({
      policyName: 'service1-demoapp',
    }).dashboard, indent='  ')) +
    dashboard.spec.withDatasources({
      inputName: 'DS_CONTROLLER-PROMETHEUS',
      datasourceName: 'controller-prometheus',
    }),

    dashboard.new('k8s-resources') +
    dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' }) +
    dashboard.spec.withJson(std.manifestJsonEx(kubeDashboards['k8s-resources-pod.json'], indent='  ')) +
    dashboard.spec.withDatasources({
      inputName: '${datasource}',
      datasourceName: 'operations-prometheus',
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
