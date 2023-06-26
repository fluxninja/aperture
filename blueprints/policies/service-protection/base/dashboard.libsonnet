local adaptiveLoadSchedulerDashboard = import '../../../../blueprints/dashboards/flow-control/adaptive-load-scheduler/dashboard.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg, params={}) {
  local updatedConfig = config + cfg,
  local protectionDashboard = adaptiveLoadSchedulerDashboard(updatedConfig).dashboard,

  dashboard: protectionDashboard,
}
