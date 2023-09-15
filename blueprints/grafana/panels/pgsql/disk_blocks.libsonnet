local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local diskBlocks = statPanel('Block Reads 1h',
                               cfg.dashboard.datasource.name,
                               'increase(postgresql_blocks_read_total{%(filters)}[5m])',
                               stringFilters),
  panel: diskBlocks.panel,
}
