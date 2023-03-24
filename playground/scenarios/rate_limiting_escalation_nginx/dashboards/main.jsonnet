local apertureDashboard = import '../../../resources/grafana-dashboard/main.libsonnet';

local aperture = import '../../../../blueprints/main.libsonnet';
local rateLimitpolicyDashboard = aperture.policies.StaticRateLimiting.dashboard;

local dashboardMixin = std.parseJson(std.extVar('APERTURE_DASHBOARD'));

local rateLimitPanel =
  rateLimitpolicyDashboard({
    policy_name: std.extVar('POLICY_NAME'),
  }).dashboard.panels[0];

local policyDashBoardMixin =
  dashboardMixin
  {
    panels+: [rateLimitPanel { id: std.length(dashboardMixin.panels) + 2 }],
  };

apertureDashboard(policyDashBoardMixin).dashboards
