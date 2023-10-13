local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local totalRequests = statPanel('Total Requests',
                                  cfg.dashboard.datasource.name,
                                  'sum(increase(workload_requests_total{%(filters)s}[$__range]))',
                                  stringFilters,
                                  h=10,
                                  w=8,
                                  panelColor='blue',
                                  graphMode='area',
                                  unit='short'),

  panel: totalRequests.panel,
}
