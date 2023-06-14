local base = import './base_dashboard.libsonnet';
local panelLibrary = import './panels/panel_library.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyFile, cfg) {
  local config = defaultConfig + cfg,
  local dashboard = base(config),
  local rateLimiterPanel = panelLibrary.rate_limiter(config),
  local policyJSON = std.parseYaml(policyFile),

  local flowControlComponents = std.flattenArrays([
    if std.objectHas(component, 'flow_control')
    then
      std.objectFields(component.flow_control)
    for component in policyJSON.spec.circuit.components
  ]),
  local flowControlPanels = [
    panelLibrary[std.toString(component)](config).panel
    for component in flowControlComponents
  ],

  local final = dashboard.baseDashboard
                + g.dashboard.withPanels(flowControlPanels),

  dashboard: final,
}
