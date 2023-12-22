local schedulerRowsFn = import '../scheduler/rows-fn.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  local componentID = std.get(component.component, 'parent_component_id', default=component.component_id);
  schedulerRowsFn(datasourceName, policyName, componentID, extraFilters={})
