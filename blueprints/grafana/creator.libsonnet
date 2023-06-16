local panelLibrary = import './panels/panel_library.libsonnet';
local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyFile, cfg) {
  local config = defaultConfig + cfg,
  local dashboard = base(config),
  local rateLimiterPanel = panelLibrary.rate_limiter(config),
  local policyJSON = std.parseYaml(policyFile),

  //policyJSON: policyJSON,

  // Flow Control Panels
  local flowControlComponents = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(component, 'flow_control')
    then
      std.objectFields(component.flow_control)
    for component in policyJSON.spec.circuit.components
  ])),
  //flowControlComponents: flowControlComponents,

  local flowControlPanels = std.filter(function(x) x != null, [
    if flowControlComponents != null
    then if std.objectHas(panelLibrary, std.toString(component))
    then
      panelLibrary[std.toString(component)](config).panel
    for component in flowControlComponents
  ]),
  //flowControlPanels: flowControlPanels,

  // Other first-level Panels
  local otherPanels = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(panelLibrary, std.toString(componentName))
    then
      unwrap(std.toString(componentName), component, config).panel
    for component in policyJSON.spec.circuit.components
    for componentName in std.objectFields(component)
  ])),
  otherPanels: otherPanels,

  local panels = flowControlPanels + otherPanels,
  local final = dashboard.baseDashboard
                + g.dashboard.withPanels(panels),
  dashboard: final,
}
