local base = import './utils/base_dashboard.libsonnet';
local defaultConfig = import './utils/default_config.libsonnet';
local unwrapInfraMeter = import './utils/unwrap_infra_meter_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyFile, policyName, datasource, extraFilters={}) {
  local receiverDashboard = base('Receiver Dashboard - %s' % policyName, defaultConfig.dashboard.refresh_interval),
  local policyJSON = std.parseYaml(policyFile),

  local infraMeters =
    if std.objectHas(policyJSON, 'spec') &&
       std.objectHas(policyJSON.spec, 'resources') &&
       std.objectHas(policyJSON.spec.resources, 'infra_meters') &&
       std.length(std.objectFields(policyJSON.spec.resources.infra_meters)) > 0
    then policyJSON.spec.resources.infra_meters
    else {},

  local receiverDashboards = {
    ['receiver' + '-' + infraMeter + '-' + receiver]:
      receiverDashboard.baseDashboard + g.dashboard.withPanels(
        unwrapInfraMeter(receiver, policyName, infraMeter, datasource, extraFilters).panel
      )
    for infraMeter in std.objectFields(infraMeters)
    if std.objectHas(infraMeters[infraMeter], 'receivers')
    for receiver in std.objectFields(infraMeters[infraMeter].receivers)
  },

  dashboards: receiverDashboards,
}
