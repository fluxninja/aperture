local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanelPanel = import '../../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local dbSize = timeSeriesPanelPanel('DB Size Bytes',
                                      cfg.dashboard.datasource.name,
                                      'rate(postgresql_database_size_bytes{%(filters)s}[1m])',
                                      stringFilters,
                                      'Bytes'),
  panel: dbSize.panel,
}
