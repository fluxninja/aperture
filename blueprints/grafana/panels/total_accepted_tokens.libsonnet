local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local acceptedTokens = statPanel('Total Accepted Tokens',
                                   datasourceName,
                                   'sum(increase(accepted_tokens_total{%(filters)s}[$__range]))',
                                   stringFilters,
                                   h=10,
                                   w=8,
                                   graphMode='area',
                                   unit='short'),

  panel: acceptedTokens.panel,
}
