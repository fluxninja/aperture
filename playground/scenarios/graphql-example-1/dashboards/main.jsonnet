local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';

local grafana = grafanaOperator.integreatly.v1alpha1.grafana;
local dashboard = grafanaOperator.integreatly.v1alpha1.grafanaDashboard;

local dashboardMixin = std.parseJson(std.extVar('APERTURE_DASHBOARD'));

local dashboards =
  dashboard.new('graphql-static-rate-limiting-dashboard') +
  dashboard.metadata.withNamespace('aperture-controller') +
  dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' }) +
  dashboard.spec.withJson(std.manifestJsonEx(dashboardMixin, indent='  ')) +
  dashboard.spec.withDatasources({
    inputName: 'DS_CONTROLLER-PROMETHEUS',
    datasourceName: 'controller-prometheus',
  });

dashboards
