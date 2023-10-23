local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID }),
  local signalName = component.component.out_ports.output.signal_name,
  local query = component.component.query_string,

  local queryPanel = timeSeriesPanel('Query for ' + signalName, datasourceName, query, stringFilters),

  panel: queryPanel.panel,
}
