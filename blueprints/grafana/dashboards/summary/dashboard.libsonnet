local base = import '../../dashboard/base.libsonnet';
local defaultConfig = import '../../dashboard/default-config.libsonnet';
local panelLibrary = import 'panel-library.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(graphObj, policyName, datasource, extraFilters={})
  local dashboard = base('Policy Summary - %s' % policyName, defaultConfig.refresh_interval);
  local internalComponents = graphObj.internal_components;

  local panels = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(panelLibrary, c.component_name)
    then
      panelLibrary[c.component_name](datasource, policyName, c, extraFilters)
    for c in internalComponents
  ]));

  dashboard.baseDashboard + g.dashboard.withPanels(panels)
