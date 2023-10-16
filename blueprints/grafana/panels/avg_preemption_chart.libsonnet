local barChartPanel = import '../utils/bar_chart_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local preemptionRequest = barChartPanel('Average Preemption Tokens per Request by Workload',
                                          cfg.dashboard.datasource.name,
                                          '( sum by (workload_index) (rate(workload_preempted_tokens_sum{%(filters)s}[$__range])) - sum by (workload_index) (rate(workload_delayed_tokens_sum{%(filters)s}[$__range])) ) / (( sum by (workload_index) (rate(workload_preempted_tokens_count{%(filters)s}[$__range])) + sum by (workload_index) (rate(workload_delayed_tokens_count{%(filters)s}[$__range])) + sum by (workload_index) (rate(workload_on_time_total{%(filters)s}[$__range])) ) != 0)',
                                          stringFilters,
                                          description='Preemption measures the average number of tokens by which a request from a specific workload was preempted, compared to a purely FIFO (first-in-first-out) ordering of requests. A negative preemption value indicates that the workload was delayed rather than preempted.',
                                          instantQuery=true,
                                          range=false),

  panel: preemptionRequest.panel,
}
