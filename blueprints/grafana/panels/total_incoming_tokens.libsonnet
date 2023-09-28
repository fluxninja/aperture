local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local incomingTokens = statPanel('Total Incoming Tokens',
                                   cfg.dashboard.datasource.name,
                                   'sum(increase(incoming_tokens_total{%(filters)s}[$__range]))',
                                   stringFilters,
                                   h=10,
                                   w=8,
                                   panelColor='blue',
                                   graphMode='area'),

  panel: incomingTokens.panel,
}
