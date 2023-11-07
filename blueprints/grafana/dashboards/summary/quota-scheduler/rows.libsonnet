local schedulerRowsFn = import '../scheduler/rows-fn.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  local componentID = component.component_id;
  schedulerRowsFn(datasourceName, policyName, componentID, extraFilters={})
