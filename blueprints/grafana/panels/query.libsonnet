local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),
  local signalName = component.component.out_ports.output.signal_name,
  local query = component.component.query_string,

  local queryPanel = timeSeriesPanel('Query for ' + signalName, datasourceName, query, stringFilters),

  panel: queryPanel.panel,
}
