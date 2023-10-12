local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local rejectedRequests = statPanel('Total Rejected Requests',
                                     datasourceName,
                                     'sum(increase(workload_requests_total{%(filters)s, decision_type="DECISION_TYPE_REJECTED"}[$__range]))',
                                     stringFilters,
                                     h=10,
                                     w=8,
                                     panelColor='red',
                                     graphMode='area',
                                     noValue='No rejected requests',
                                     unit='short'),

  panel: rejectedRequests.panel,
}
