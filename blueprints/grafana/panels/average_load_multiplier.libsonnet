local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(datasourceName, policyName, component, extra_filters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extra_filters { policy_name: policyName, component_id: component.component_id }),

  local avgLoadMultiplier = statPanel('Average Load Multiplier',
                                      datasourceName,
                                      'avg(token_bucket_lm_ratio{%(filters)s})',
                                      stringFilters),

  panel: avgLoadMultiplier.panel,
}
