local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(datasourceName, policyName, component, extra_filters) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extra_filters { policy_name: policyName, component_id: componentID }),

  local tokenBucketFillrate = statPanel('Token Bucket FillRate',
                                        datasourceName,
                                        'avg(token_bucket_fill_rate{%(filters)s})',
                                        stringFilters),

  panel: tokenBucketFillrate.panel,
}
