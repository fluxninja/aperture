local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local bgwriterAllocated = statPanel('Total Buffer Allocated',
                                      cfg.dashboard.datasource.name,
                                      'postgresql_bgwriter_buffers_allocated_total{%(filters)s}',
                                      stringFilters),
  panel: bgwriterAllocated.panel,
}
