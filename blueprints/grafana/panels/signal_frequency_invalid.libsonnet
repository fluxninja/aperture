local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg, title, withVariables, signalName='') {
  local policyFilters = cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name },
  local signalFilters =
    if withVariables then
      policyFilters { signal_name: '${signal_name}', sub_circuit_id: '${sub_circuit_id}' }
    else
      policyFilters { signal_name: signalName },
  local stringFilters = utils.dictToPrometheusFilter(signalFilters),

  local signalFrequency = timeSeriesPanel(title,
                                          cfg.dashboard.datasource.name,
                                          'sum(rate(invalid_signal_readings_total{%(filters)s}[$__rate_interval]))',
                                          signalFilters),

  panel: signalFrequency.panel,
}
