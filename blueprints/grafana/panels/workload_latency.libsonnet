local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local workloadLatency = timeSeriesPanel('Workload Latency',
                                          cfg.dashboard.datasource.name,
                                          '(sum by (workload_index) (increase(workload_latency_ms_sum{%(filters)s}[$__rate_interval])))/(sum by (workload_index) (increase(workload_latency_ms_count{%(filters)s}[$__rate_interval])))',
                                          stringFilters,
                                          'Latency',
                                          'ms'),

  panel: workloadLatency.panel,
}
