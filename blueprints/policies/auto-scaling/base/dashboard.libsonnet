local autoScaleDashboard = import '../../../../blueprints/dashboards/auto-scale/dashboard.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg) {
  local params = config + cfg,
  local dashboardDef = autoScaleDashboard(params).dashboard,

  dashboard: dashboardDef,
}
