local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local workloadDecAccepted = timeSeriesPanel('Workload Decisions (accepted)',
                                              datasourceName,
                                              'sum by(workload_index, decision_type) (rate(workload_requests_total{%(filters)s,decision_type="DECISION_TYPE_ACCEPTED"}[$__rate_interval]))',
                                              stringFilters,
                                              'Decisions',
                                              'reqps'),

  panel: workloadDecAccepted.panel,
}
