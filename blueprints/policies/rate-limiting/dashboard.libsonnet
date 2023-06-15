local rateLimiterDashboard = import '../../../blueprints/dashboards/flow-control/rate-limiter/dashboard.libsonnet';
local utils = import '../policy-utils.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config + cfg,

  local dashboardDef = rateLimiterDashboard(params).dashboard,
  dashboard: dashboardDef,
}
