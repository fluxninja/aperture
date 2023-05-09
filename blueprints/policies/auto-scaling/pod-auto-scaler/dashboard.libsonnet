local baseDashboardFn = import '../base/dashboard.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config + cfg,
  dashboard: baseDashboardFn(params),
}
