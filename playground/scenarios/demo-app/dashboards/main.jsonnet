local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';

local aperture = import '../../../../blueprints/main.libsonnet';
local rateLimitpolicyDashboard = aperture.policies.StaticRateLimiting.dashboard;

local grafana = grafanaOperator.integreatly.v1alpha1.grafana;
local dashboard = grafanaOperator.integreatly.v1alpha1.grafanaDashboard;

local dashboardMixin = std.parseJson(std.extVar('APERTURE_DASHBOARD'));

local rateLimitPanel =
  rateLimitpolicyDashboard({
    policy_name: 'service1-demo-app',
  }).dashboard.panels[0];

local policyDashBoardMixin =
  dashboardMixin
  {
    panels+: [rateLimitPanel { id: std.length(dashboardMixin.panels) + 2 }],
  };

local dashboards =
  dashboard.new('service1-demo-app-dashboard') +
  dashboard.metadata.withNamespace('aperture-controller') +
  dashboard.metadata.withLabels({ 'fluxninja.com/grafana-instance': 'aperture-grafana' }) +
  dashboard.spec.withJson(std.manifestJsonEx(policyDashBoardMixin, indent='  ')) +
  dashboard.spec.withDatasources({
    inputName: 'DS_CONTROLLER-PROMETHEUS',
    datasourceName: 'controller-prometheus',
  });

dashboards
