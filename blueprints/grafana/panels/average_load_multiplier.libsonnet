local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local avgLoadMultiplier = statPanel('Average Load Multiplier',
                                      cfg.dashboard.datasource.name,
                                      'avg(token_bucket_lm_ratio{%(filters)s})',
                                      stringFilters),

  panel: avgLoadMultiplier.panel,
}
