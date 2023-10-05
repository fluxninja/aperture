local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';
local panelLibrary = import 'panel_library.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(componentsList, policyName, datasource, extraFilters={}) {
  local dashboard = base('Aperture Dashboard - %s' % policyName, defaultConfig.dashboard.refresh_interval),
  local components = std.parseJson(componentsList),

  local panels = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(panelLibrary, c.component_name)
    then
      unwrap(datasource, policyName, c, extraFilters).panel
    for c in components.internal_components
  ])),

  local final = dashboard.baseDashboard + g.dashboard.withPanels(panels),
  dashboard: final,
}
