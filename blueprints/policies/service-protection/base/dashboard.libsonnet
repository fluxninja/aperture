local baseAutoScalingDashboardFn = import '../../../../blueprints/dashboards/auto-scale/dashboard.libsonnet';
local adaptiveLoadSchedulerDashboard = import '../../../../blueprints/dashboards/flow-control/adaptive-load-scheduler/dashboard.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg) {
  local params = config + cfg,
  local protectionDashboard = adaptiveLoadSchedulerDashboard(params).dashboard,
  local autoScalingParams = {
    policy+: params.policy.auto_scaling {
      policy_name: params.policy.policy_name,
    },

    dashboard: params.dashboard,
  },

  local baseAutoScalingDashboard = baseAutoScalingDashboardFn(autoScalingParams).dashboard,

  local maxPanelYAxis = std.reverse(std.sort(protectionDashboard.panels, keyF=function(panel) panel.gridPos.y))[0].gridPos.y,
  local maxId = std.reverse(std.sort(protectionDashboard.panels, keyF=function(panel) '%s' % panel.id))[0].id,


  local protectionAndEscalationDashboard =
    protectionDashboard {
      panels+: [
        baseAutoScalingDashboard.dashboard.panels[panel_idx] {
          id: maxId + panel_idx + 1,
          gridPos: { x: 0, y: (maxPanelYAxis + (panel_idx + 1) * 10), w: 24, h: 10 },
        }
        for panel_idx in std.range(0, std.length(baseAutoScalingDashboard.dashboard.panels) - 1)
      ],
    },

  dashboard: if std.objectHas(params.policy, 'auto_scaling') then protectionAndEscalationDashboard else protectionDashboard,
}
