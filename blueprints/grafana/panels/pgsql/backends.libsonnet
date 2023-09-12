local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local backendConn = timeSeriesPanel('Connection Per Minute',
                                      cfg.dashboard.datasource.name,
                                      'rate(postgresql_backends{%(filters)s}[1m])*60',
                                      stringFilters,
                                      'Connections'),

  panel: backendConn.panel,
}
