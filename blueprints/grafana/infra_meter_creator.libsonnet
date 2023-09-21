local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrapInfraMeter = import './utils/unwrap_infra_meter_panel.libsonnet';
local unwrap = import './utils/unwrap_panels.libsonnet';
local panelLibrary = import 'panel_library.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyFile, cfg) {

  local policyName = cfg.policy.policy_name,
  local config = defaultConfig + cfg,
  local dashboard = base(config),
  local policyJSON =
    if std.isObject(policyFile)
    then policyFile
    else if std.startsWith(policyFile, '{')
    then std.parseJson(policyFile)
    else std.parseYaml(policyFile),

  local infraMeters =
    if std.objectHas(policyJSON, 'spec') &&
       std.objectHas(policyJSON.spec, 'resources') &&
       std.objectHas(policyJSON.spec.resources, 'infra_meters')
    then policyJSON.spec.resources.infra_meters
    else {},

  local receiverNames = [
    receiver
    for meter in std.objectFields(infraMeters)
    for receiver in std.objectFields(infraMeters[meter].receivers)
    if std.objectHas(infraMeters[meter].receivers, receiver)
  ],

  local receiverToMeterMap = {
    [receiver]: meter
    for meter in std.objectFields(infraMeters)
    for receiver in std.objectFields(infraMeters[meter].receivers)
  },

  local receiverDashboards = {
    [receiver]: dashboard.baseDashboard + g.dashboard.withPanels(
      std.flattenArrays([
        if std.objectHas(panelLibrary, receiver)
        then unwrapInfraMeter(receiver, policyName, receiverToMeterMap[receiver], cfg.dashboard.datasource, cfg.dashboard.extra_filters).panel
        else [],
      ])
    )
    for receiver in receiverNames
  },

  receiverDashboards: receiverDashboards,
}
