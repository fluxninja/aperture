local promqlDashboardFn = import '../promql/dashboard.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config + cfg,
  local variantName = params.dashboard.variant_name,
  local query = params.policy.promql_query,

  local promqlDashboard = promqlDashboardFn(params),

  dashboard: promqlDashboard.dashboard,
}
