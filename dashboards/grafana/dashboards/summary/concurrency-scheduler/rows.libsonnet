local schedulerRowsFn = import '../scheduler/rows-fn.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  local componentID = component.component.parent_component_id;
  local scheduler = if 'scheduler' in component.component then component.component.scheduler else {};
  schedulerRowsFn(datasourceName, policyName, componentID, scheduler, extraFilters={})
