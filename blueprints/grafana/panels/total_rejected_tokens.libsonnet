local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local rejectedTokens = statPanel('Total Rejected Tokens',
                                   cfg.dashboard.datasource.name,
                                   'sum(increase(rejected_tokens_total{%(filters)s}[$__range]))',
                                   stringFilters,
                                   h=10,
                                   w=8,
                                   panelColor='red',
                                   graphMode='area',
                                   unit='short'),
  panel: rejectedTokens.panel,
}
