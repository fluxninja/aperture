local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, axisLabel='', unit='', description='', h=10, w=24, targets=[]) {
  local timeseries =
    g.panel.timeSeries.new(title)
    + g.panel.timeSeries.panelOptions.withDescription(description)
    + g.panel.timeSeries.datasource.withType('prometheus')
    + g.panel.timeSeries.datasource.withUid(dsName)
    + g.panel.timeSeries.standardOptions.withUnit(unit)
    + g.panel.timeSeries.fieldConfig.defaults.custom.withAxisLabel(axisLabel)
    + g.panel.timeSeries.fieldConfig.defaults.custom.withFillOpacity(10)
    + g.panel.timeSeries.gridPos.withH(h)
    + g.panel.timeSeries.gridPos.withW(w)
    + g.panel.timeSeries.queryOptions.withInterval('10s'),

  local withMultipleTargets =
    if targets != []
    then
      timeseries + g.panel.timeSeries.withTargets(targets)
    else
      timeseries + g.panel.timeSeries.withTargets([
        g.query.prometheus.new(dsName, query % { filters: strFilters })
        + g.query.prometheus.withIntervalFactor(1),
      ]),

  panel: withMultipleTargets,
}
