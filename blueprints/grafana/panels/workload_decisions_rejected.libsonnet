local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local workloadDecRejected = timeSeriesPanel('Workload Decisions (rejected)',
                                              datasourceName,
                                              'sum by(workload_index, decision_type) (rate(workload_requests_total{%(filters)s,decision_type="DECISION_TYPE_REJECTED"}[$__rate_interval]))',
                                              stringFilters,
                                              'Decisions',
                                              'reqps'),

  panel: workloadDecRejected.panel,
}
