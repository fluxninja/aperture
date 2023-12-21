local statPanel = import '../../../panels/stat.libsonnet';
local timeSeriesPanel = import '../../../panels/time-series.libsonnet';
local portUtils = import '../../../utils/port.libsonnet';
local prometheusUtils = import '../../../utils/prometheus.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  local componentID = std.get(component.component, 'parent_component_id', default=component.component_id);
  local stringFilters = prometheusUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID });


  local row1 = [
    statPanel(
      'Total Requests',
      datasourceName,
      'sum(increase(rate_limiter_counter_total{%(filters)s}[$__range]))' % { filters: stringFilters },
      x=0,
      h=10,
      w=8,
      panelColor='blue',
      graphMode='area',
      unit='short'
    ),
    statPanel(
      'Total Accepted Requests',
      datasourceName,
      'sum(increase(rate_limiter_counter_total{%(filters)s, decision_type="DECISION_TYPE_ACCEPTED"}[$__range]))' % { filters: stringFilters },
      x=8,
      h=10,
      w=8,
      graphMode='area',
      unit='short'
    ),
    statPanel(
      'Total Rejected Requests',
      datasourceName,
      'sum(increase(rate_limiter_counter_total{%(filters)s, decision_type="DECISION_TYPE_REJECTED"}[$__range]))' % { filters: stringFilters },
      x=16,
      h=10,
      w=8,
      panelColor='red',
      graphMode='area',
      noValue='No rejected requests',
      unit='short'
    ),
  ];

  local targets =
    [
      g.query.prometheus.new(datasourceName, 'sum by(decision_type) (rate(rate_limiter_counter_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
      + g.query.prometheus.withIntervalFactor(1),
    ] +
    if 'limit_by_label_key' in component.component && component.component.limit_by_label_key != ''
    then
      portUtils.targetsForInPort(datasourceName, component, 'fill_amount', policyName, extraFilters)
    else
      [];

  local row2 = [
    timeSeriesPanel(
      'Aperture Rate Limiter',
      datasourceName,
      axisLabel='Decisions',
      unit='reqps',
      targets=targets,
    ),
  ];

  [row1, row2]
