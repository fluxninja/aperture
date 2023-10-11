local barGaugePanel = import '../utils/bar_gauge_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local legendFormat = '{{ instance }} - {{ policy_name }}',

  local wfqSchedulerHeapRequests = barGaugePanel('WFQ Scheduler Heap Requests',
                                                 datasourceName,
                                                 'avg(wfq_requests_total{%(filters)s})',
                                                 stringFilters,
                                                 legendFormat),

  panel: wfqSchedulerHeapRequests.panel,
}
