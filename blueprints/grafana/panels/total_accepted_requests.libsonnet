local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local acceptedRequests = statPanel('Total Accepted Requests',
                                     cfg.dashboard.datasource.name,
                                     'sum(increase(workload_requests_total{%(filters)s, decision_type="DECISION_TYPE_ACCEPTED"}[$__range]))',
                                     stringFilters,
                                     h=10,
                                     w=8,
                                     graphMode='area',
                                     unit='short'),
  panel: acceptedRequests.panel,
}
