local baseAutoScalingDashboardFn = import '../../auto-scaling/pod-auto-scaler/dashboard.libsonnet';
local baseServiceProtectionDashboardFn = import '../../service-protection/average-latency/dashboard.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

function(cfg) {
  local params = config + cfg,
  local policyName = params.policy.policy_name,
  local variantName = params.dashboard.variant_name,

  local baseServiceProtectionDashboard = baseServiceProtectionDashboardFn(params),
  local baseAutoScalingDashboard = baseAutoScalingDashboardFn(params),

  local maxPanelYAxis = std.reverse(std.sort(baseServiceProtectionDashboard.dashboard.panels, keyF=function(panel) panel.gridPos.y))[0].gridPos.y,
  local maxId = std.reverse(std.sort(baseServiceProtectionDashboard.dashboard.panels, keyF=function(panel) '%s' % panel.id))[0].id,

  local extendedDashboard =
    baseServiceProtectionDashboard.dashboard {
      panels+: [
        baseAutoScalingDashboard.dashboard.panels[panel_idx] {
          id: maxId + panel_idx + 1,
          gridPos: { x: 0, y: (maxPanelYAxis + (panel_idx + 1) * 10), w: 24, h: 10 },
        }
        for panel_idx in std.range(0, std.length(baseAutoScalingDashboard.dashboard.panels) - 1)
      ],
    },

  dashboard: extendedDashboard,
}
