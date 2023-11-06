local schedulerPanelsFn = import '../scheduler/panels-fn.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  local componentID = component.component_id;
  schedulerPanelsFn(datasourceName, policyName, componentID, extraFilters={})
