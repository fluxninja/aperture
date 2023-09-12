local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local totalCommits = timeSeriesPanel('Commits Per Minute',
                                       cfg.dashboard.datasource.name,
                                       'rate(postgresql_commits_total{%(filters)s}[1m])*60',
                                       stringFilters,
                                       'Commits'),

  panel: totalCommits.panel,
}
