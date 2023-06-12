local loadRampDashboard = import '../../../../blueprints/dashboards/flow-control/load-ramp/dashboard.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config + cfg,
  local dashboardDef = loadRampDashboard(params).dashboard,
  dashboard: dashboardDef,
}
