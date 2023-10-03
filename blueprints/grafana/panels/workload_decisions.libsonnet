local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local quotaScheduler = timeSeriesPanel('Workload Decisions',
                                         datasourceName,
                                         'sum by(decision_type) (rate(workload_requests_total{%(filters)s}[$__rate_interval]))',
                                         stringFilters,
                                         'Decisions',
                                         'reqps'),

  panel: quotaScheduler.panel,
}
