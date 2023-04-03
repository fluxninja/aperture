local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;

function(cfg) {
  local params = config.common + config.dashboard + cfg,

  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja',
      editable=true,
      schemaVersion=18,
      refresh=params.refresh_interval,
      time_from=params.time_from,
      time_to=params.time_to
    ),

  dashboard: dashboardDef,
}
