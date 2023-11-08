local portUtils = import '../../../utils/port.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  [portUtils.panelsForInPort('Desired Replicas', datasourceName, component, 'replicas', policyName, extraFilters)]
