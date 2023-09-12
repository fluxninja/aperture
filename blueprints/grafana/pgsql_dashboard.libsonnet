local panelLibrary = import './panel_library.libsonnet';
local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  local config = defaultConfig + cfg,
  local dashboard = base(config),

  local pgsqlPanels = unwrap(std.toString('pgsql'), {}, config).panel,
  local final = dashboard.baseDashboard + g.dashboard.withPanels(pgsqlPanels),
  dashboard: final,
}
