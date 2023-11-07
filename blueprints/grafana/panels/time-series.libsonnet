local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(title, dsName, axisLabel='', unit='', x=0, h=10, w=24, query='', targets=[], description='')
  local timeseries =
    g.panel.timeSeries.new(title)
    + g.panel.timeSeries.panelOptions.withDescription(description)
    + g.panel.timeSeries.datasource.withType('prometheus')
    + g.panel.timeSeries.datasource.withUid(dsName)
    + g.panel.timeSeries.standardOptions.withUnit(unit)
    + g.panel.timeSeries.fieldConfig.defaults.custom.withAxisLabel(axisLabel)
    + g.panel.timeSeries.fieldConfig.defaults.custom.withFillOpacity(10)
    + g.panel.timeSeries.gridPos.withX(x)
    + g.panel.timeSeries.gridPos.withH(h)
    + g.panel.timeSeries.gridPos.withW(w)
    + g.panel.timeSeries.queryOptions.withInterval('10s');

  if targets != []
  then
    timeseries + g.panel.timeSeries.queryOptions.withTargets(targets)
  else
    timeseries + g.panel.timeSeries.queryOptions.withTargets([
      g.query.prometheus.new(dsName, query)
      + g.query.prometheus.withIntervalFactor(1),
    ])
