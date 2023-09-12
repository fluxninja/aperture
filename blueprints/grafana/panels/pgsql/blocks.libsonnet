local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local currentBlocks = timeSeriesPanel('Blocks Per Min',
                                        cfg.dashboard.datasource.name,
                                        'rate(postgresql_blocks_read_total{policy_name="postgres-connections"}[1m])*60',
                                        stringFilters,
                                        'Blocks'),
  panel: currentBlocks.panel,
}
