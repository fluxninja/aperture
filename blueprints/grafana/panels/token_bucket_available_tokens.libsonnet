local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local tokenBucketAvailableTokens = statPanel('Token Bucket Available Tokens',
                                               cfg.dashboard.datasource.name,
                                               'avg(token_bucket_available_tokens_total{%(filters)s})',
                                               stringFilters),

  panel: tokenBucketAvailableTokens.panel,
}
