local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local incomingTokens = statPanel('Total Incoming Tokens',
                                   datasourceName,
                                   'sum(increase(incoming_tokens_total{%(filters)s}[$__range]))',
                                   stringFilters,
                                   h=10,
                                   w=8,
                                   panelColor='blue',
                                   graphMode='area',
                                   unit='short'),

  panel: incomingTokens.panel,
}
