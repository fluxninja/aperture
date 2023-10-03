local barGaugePanel = import '../utils/bar_gauge_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local preemptionRequest = barGaugePanel('Average preemption per request by Workload',
                                          cfg.dashboard.datasource.name,
                                          '( sum by (workload_index) (rate(workload_preempted_tokens_sum[$__rate_interval])) - sum by (workload_index) (rate(workload_delayed_tokens_sum[$__rate_interval])) ) / ( sum by (workload_index) (rate(workload_preempted_tokens_count[$__rate_interval])) + sum by (workload_index) (rate(workload_delayed_tokens_count[$__rate_interval])) )',
                                          stringFilters,
                                          unit='short',
                                          instantQuery=true,
                                          range=false),

  panel: preemptionRequest.panel,
}
