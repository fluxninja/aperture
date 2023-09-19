local barGaugePanel = import '../utils/bar_gauge_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local legendFormat = '{{ instance }} - {{ policy_name }}',

  local wfqSchedulerHeapRequests = barGaugePanel('WFQ Scheduler Heap Requests',
                                                 cfg.dashboard.datasource.name,
                                                 'avg(wfq_requests_total{%(filters)s})',
                                                 stringFilters,
                                                 legendFormat),

  panel: wfqSchedulerHeapRequests.panel,
}
