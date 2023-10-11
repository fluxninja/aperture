local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local rejectedTokens = statPanel('Total Rejected Tokens',
                                   datasourceName,
                                   'sum(increase(incoming_tokens_total{%(filters)s}[$__range]) - increase(accepted_tokens_total{%(filters)s}[$__range]))',
                                   stringFilters,
                                   h=10,
                                   w=8,
                                   panelColor='red',
                                   graphMode='area',
                                   unit='short'),
  panel: rejectedTokens.panel,
}
