local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local avgPreemption = timeSeriesPanel('Average preemption per request by Workload',
                                        cfg.dashboard.datasource.name,
                                        '( sum by (workload_index) (rate(workload_preempted_tokens_sum{%(filters)s}[$__rate_interval])) - sum by (workload_index) (rate(workload_delayed_tokens_sum{%(filters)s}[$__rate_interval])) ) / ( sum by (workload_index) (rate(workload_preempted_tokens_count{%(filters)s}[$__rate_interval])) + sum by (workload_index) (rate(workload_delayed_tokens_count{%(filters)s}[$__rate_interval])) + sum by (workload_index) (rate(workload_on_time_total{%(filters)s}[$__rate_interval])) )',
                                        stringFilters,
                                        'Tokens'),

  panel: avgPreemption.panel,
}
