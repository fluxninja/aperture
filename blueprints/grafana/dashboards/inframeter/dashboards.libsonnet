local base = import '../../dashboard/base.libsonnet';
local defaultConfig = import '../../dashboard/default-config.libsonnet';
local panelLibrary = import './panel-library.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyFile, policyName, datasource, extraFilters={})
  local receiverDashboard = base('Receiver Dashboard - %s' % policyName, defaultConfig.refresh_interval);
  local policyJSON = std.parseYaml(policyFile);

  local infraMeters =
    if std.objectHas(policyJSON, 'spec') &&
       std.objectHas(policyJSON.spec, 'resources') &&
       std.objectHas(policyJSON.spec.resources, 'infra_meters') &&
       std.length(std.objectFields(policyJSON.spec.resources.infra_meters)) > 0
    then policyJSON.spec.resources.infra_meters
    else {};


  local generatePanels(receiverKey, policyName, infraMeter, datasource, extraFilters) =
    if std.objectHas(panelLibrary, receiverKey)
    then panelLibrary[receiverKey](policyName, infraMeter, datasource, extraFilters)
    else null;

  local generateDashboardKey(receiver, infraMeter, panels) =
    if panels != null
    then 'receiver' + '-' + infraMeter + '-' + receiver
    else null;

  local dashboardEntries = [
    {
      local panels = generatePanels(std.split(receiver, '/')[0], policyName, infraMeter, datasource, extraFilters),
      key: generateDashboardKey(receiver, infraMeter, panels),
      dashboard: receiverDashboard.baseDashboard + g.dashboard.withPanels(panels),
    }
    for infraMeter in std.objectFields(infraMeters)
    if std.objectHas(infraMeters[infraMeter], 'receivers')
    for receiver in std.objectFields(infraMeters[infraMeter].receivers)
  ];

  { [entry.key]: entry.dashboard for entry in dashboardEntries if entry.key != null }
