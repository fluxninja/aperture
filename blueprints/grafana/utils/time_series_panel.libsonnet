local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, axisLabel='', unit='', h=10, w=24) {
  local timeseries =
    g.panel.timeSeries.new(title)
    + g.panel.timeSeries.datasource.withType('prometheus')
    + g.panel.timeSeries.datasource.withUid(dsName)
    + g.panel.timeSeries.withTargets([
      g.query.prometheus.new(dsName, query % { filters: strFilters })
      + g.query.prometheus.withIntervalFactor(1),
    ])
    + g.panel.timeSeries.standardOptions.withUnit(unit)
    + g.panel.timeSeries.fieldConfig.defaults.custom.withAxisLabel(axisLabel)
    + g.panel.timeSeries.gridPos.withH(h)
    + g.panel.timeSeries.gridPos.withW(w),

  panel: timeseries,
}
