local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local workloadLatency = timeSeriesPanel('Request in Queue Duration',
                                          cfg.dashboard.datasource.name,
                                          '(sum by (component_id) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__rate_interval])))/(sum by (component_id) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__rate_interval])))',
                                          stringFilters,
                                          'Wait Time',
                                          'ms'),

  panel: workloadLatency.panel,
}
