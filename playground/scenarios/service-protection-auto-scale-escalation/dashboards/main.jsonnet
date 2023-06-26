local apertureDashboard = import '../../../resources/grafana-dashboard/main.libsonnet';

local aperture = import '../../../../blueprints/main.libsonnet';
local autoScalingPolicyDashboard = aperture.policies.PodAutoScaler.dashboard;

local dashboardMixin = std.parseJson(std.extVar('APERTURE_DASHBOARD'));

local AutoScale =
  autoScalingPolicyDashboard({
    policy+: {
      policy_name: std.extVar('POLICY_NAME'),
    },
  }).dashboard;

local maxId = std.reverse(std.sort(dashboardMixin.panels, keyF=function(panel) '%s' % panel.id))[0].id;
local maxPanelYAxis = std.reverse(std.sort(dashboardMixin.panels, keyF=function(panel) panel.gridPos.y))[0].gridPos.y;

local policyDashBoardMixin =
  dashboardMixin {
    panels+: [
      AutoScale.panels[panel_idx] {
        id: maxId + panel_idx + 1,
        gridPos: { x: 0, y: (maxPanelYAxis + (panel_idx + 1) * 10), w: 24, h: 10 },
      }
      for panel_idx in std.range(0, std.length(AutoScale.panels) - 1)
    ],
  };

apertureDashboard(policyDashBoardMixin).dashboards
