local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local totalRollbacks = timeSeriesPanel('Rollbacks Per Minute',
                                         cfg.dashboard.datasource.name,
                                         'rate(postgresql_rollbacks_total{%(filters)s}[1m])',
                                         stringFilters,
                                         'Rollbacks'),
  panel: totalRollbacks.panel,
}
