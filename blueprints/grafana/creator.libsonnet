local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';
local panelLibrary = import 'panel_library.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyFile, cfg, add_signals=false) {
  local config = defaultConfig + cfg,
  local dashboard = base(config),
  local rateLimiterPanel = panelLibrary.rate_limiter(config),
  local policyJSON = std.parseYaml(policyFile),

  // Flow Control Panels
  local flowControlComponents = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(component, 'flow_control')
    then
      std.objectFields(component.flow_control)
    for component in policyJSON.spec.circuit.components
  ])),

  local flowControlPanels = std.flattenArrays(std.filter(function(x) x != null, [
    if flowControlComponents != null
    then if std.objectHas(panelLibrary, std.toString(component))
    then
      unwrap(std.toString(component), {}, config).panel
    for component in flowControlComponents
  ])),

  // Other first-level Panels
  local otherPanels = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(panelLibrary, std.toString(componentName))
    then
      unwrap(std.toString(componentName), component, config).panel
    for component in policyJSON.spec.circuit.components
    for componentName in std.objectFields(component)
  ])),

  // Add signals panels if requested
  local signals =
    if add_signals
    then
      otherPanels + unwrap(std.toString('signals'), {}, config).panel
    else [],

  local panels = flowControlPanels + otherPanels + signals,
  local final = dashboard.baseDashboard
                + g.dashboard.withPanels(panels),
  dashboard: final,
}
