local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local currentIndex = statPanel('Index cache hit rate %',
                                 cfg.dashboard.datasource.name,
                                 'sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="postgresql",source="idx_hit"}) / (sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="postgresql",source="idx_hit"}) + sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="postgresql",source="idx_read"}))',
                                 stringFilters),
  panel: currentIndex.panel,
}
