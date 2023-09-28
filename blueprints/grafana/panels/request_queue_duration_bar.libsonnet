local barGaugePanel = import '../utils/bar_gauge_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local legendFormat = '{{workload_index}}',

  local requestsDuration = barGaugePanel('Requests Duration',
                                         cfg.dashboard.datasource.name,
                                         '(sum by (workload_index) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__rate_interval])))/(sum by (workload_index) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__rate_interval])))',
                                         stringFilters,
                                         legendFormat,
                                         values=true),

  panel: requestsDuration.panel,
}
