local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrapInfraMeter = import './utils/unwrap_infra_meter_panel.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';
local panelLibrary = import 'panel_library.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyFile, cfg) {
  local config = defaultConfig + cfg,
  local dashboard = base(config),
  local policyName = cfg.policy.policy_name,
  local receiverDashboard = base(config { dashboard+: { title: 'Receiver Dashboard - %s' % policyName } }),

  local policyJSON =
    if std.isObject(policyFile)
    then policyFile
    else if std.startsWith(policyFile, '{')
    then std.parseJson(policyFile)
    else std.parseYaml(policyFile),

  local componentsJSON =
    if std.objectHas(policyJSON, 'spec')
    then
      policyJSON.spec.circuit.components
    else policyJSON.node.component.components,

  // Flow Control Panels
  local flowControlComponents = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(component, 'flow_control')
    then
      std.objectFields(component.flow_control)
    for component in componentsJSON
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
    for component in componentsJSON
    for componentName in std.objectFields(component)
  ])),

  local infraMeters =
    if std.objectHas(policyJSON, 'spec') &&
       std.objectHas(policyJSON.spec, 'resources') &&
       std.objectHas(policyJSON.spec.resources, 'infra_meters') &&
       std.length(std.objectFields(policyJSON.spec.resources.infra_meters)) > 0
    then policyJSON.spec.resources.infra_meters
    else {},

  local receiverDashboards = {
    ['receiver' + '-' + policyName + '-' + infraMeter + '-' + receiver + '.json']:
      receiverDashboard.baseDashboard + g.dashboard.withPanels(
        unwrapInfraMeter(receiver, policyName, infraMeter, cfg.dashboard.datasource, cfg.dashboard.extra_filters).panel
      )
    for infraMeter in std.objectFields(infraMeters)
    if std.objectHas(infraMeters[infraMeter], 'receivers')
    for receiver in std.objectFields(infraMeters[infraMeter].receivers)
  },

  local panels = flowControlPanels + otherPanels,
  local final = dashboard.baseDashboard + g.dashboard.withPanels(panels),

  dashboard: final,
  receiverDashboards: receiverDashboards,
}
