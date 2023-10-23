local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local avgPreemption = timeSeriesPanel('Average preemption per request by Workload',
                                        datasourceName,
                                        '( sum by (workload_index) (rate(workload_preempted_tokens_sum{%(filters)s}[$__rate_interval])) - sum by (workload_index) (rate(workload_delayed_tokens_sum{%(filters)s}[$__rate_interval])) ) / (( sum by (workload_index) (rate(workload_preempted_tokens_count{%(filters)s}[$__rate_interval])) + sum by (workload_index) (rate(workload_delayed_tokens_count{%(filters)s}[$__rate_interval])) + sum by (workload_index) (rate(workload_on_time_total{%(filters)s}[$__rate_interval])) ) != 0)',
                                        stringFilters,
                                        axisLabel='Tokens',
                                        description='Preemption measures the average number of tokens by which a request from a specific workload was preempted, compared to a purely FIFO (first-in-first-out) ordering of requests. A negative preemption value indicates that the workload was delayed rather than preempted.'),

  panel: avgPreemption.panel,
}
