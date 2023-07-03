local barGaugePanel = import '../utils/bar_gauge_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local wfqSchedulerFlows = barGaugePanel('WFQ Scheduler Flows',
                                          cfg.dashboard.datasource.name,
                                          'avg(wfq_flows_total{%(filters)s})',
                                          stringFilters),

  panel: wfqSchedulerFlows.panel,
}
