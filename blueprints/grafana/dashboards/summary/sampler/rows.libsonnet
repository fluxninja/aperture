local timeSeriesPanel = import '../../../panels/time-series.libsonnet';
local portUtils = import '../../../utils/port.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  local componentID = component.component_id;
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID });

  local row1 = [
    timeSeriesPanel(
      'Throughput - Accept/Reject',
      datasourceName,
      query='sum(rate(sampler_counter_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters },
    ),
  ];

  local row2 = portUtils.panelsForInPort('Accept Percentage', datasourceName, component, 'accept_percentage', policyName, extraFilters);

  [row1, row2]
