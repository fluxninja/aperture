local portUtils = import '../../../utils/port.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  [portUtils.panelsForOutPort(
    component.component_description,
    datasourceName,
    component,
    'output',
    policyName,
    extraFilters,
    description='Signal derived from periodic execution of query: ' + component.component.query_string,
  )]
