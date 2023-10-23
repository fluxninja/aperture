local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, _, extraFilters, title, withVariables, signalName='') {
  local policyFilters = extraFilters { policy_name: policyName },
  local signalFilters =
    if withVariables then
      policyFilters { signal_name: '${signal_name}', sub_circuit_id: '${sub_circuit_id}' }
    else
      policyFilters { signal_name: signalName },
  local stringFilters = utils.dictToPrometheusFilter(signalFilters),

  local signalAvg = timeSeriesPanel(title,
                                    datasourceName,
                                    'increase(signal_reading_sum{%(filters)s}[$__rate_interval]) / increase(signal_reading_count{%(filters)s}[$__rate_interval])',
                                    stringFilters),

  panel: signalAvg.panel,
}
