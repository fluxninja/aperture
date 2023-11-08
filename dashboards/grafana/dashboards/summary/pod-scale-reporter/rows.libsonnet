local timeSeriesPanel = import '../../../panels/time-series.libsonnet';
local portUtils = import '../../../utils/port.libsonnet';

function(datasourceName, policyName, component, extraFilters={})

  local targets =
    portUtils.targetsForOutPort(datasourceName, component, 'actual_replicas', policyName, extraFilters) +
    portUtils.targetsForOutPort(datasourceName, component, 'configured_replicas', policyName, extraFilters);

  [[timeSeriesPanel(
    'Actual and configured replicas',
    datasourceName,
    targets=targets,
  )]]
