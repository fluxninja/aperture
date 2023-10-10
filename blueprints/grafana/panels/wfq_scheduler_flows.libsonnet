local barGaugePanel = import '../utils/bar_gauge_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),


  local legendFormat = '{{ instance }} - {{ policy_name }}',

  local wfqSchedulerFlows = barGaugePanel('WFQ Scheduler Flows',
                                          datasourceName,
                                          'avg(wfq_flows_total{%(filters)s})',
                                          stringFilters,
                                          legendFormat),

  panel: wfqSchedulerFlows.panel,
}
