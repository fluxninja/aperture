local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local workloadDecRejected = timeSeriesPanel('Workload Decisions (rejected)',
                                              cfg.dashboard.datasource.name,
                                              'sum by(workload_index, decision_type) (rate(workload_requests_total{%(filters)s,decision_type="DECISION_TYPE_REJECTED"}[$__rate_interval]))',
                                              stringFilters,
                                              'Decisions',
                                              'reqps'),

  panel: workloadDecRejected.panel,
}
