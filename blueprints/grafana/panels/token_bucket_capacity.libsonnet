local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local tokenBucketCapacity = statPanel('Token Bucket Capacity',
                                        cfg.dashboard.datasource.name,
                                        'avg(token_bucket_capacity_total{%(filters)s})',
                                        stringFilters),

  panel: tokenBucketCapacity.panel,
}
