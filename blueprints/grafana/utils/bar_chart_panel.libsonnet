local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, h=10, w=24, legendFormat=null, queryFormat='time_series', instantQuery=false, range=true, labelSpacing=0, axisGridshow=true, axisPlacement='auto', mode='single', sort='sort') {
  local barChartPanel =
    g.panel.barChart.new(title)
    + g.panel.barChart.queryOptions.withDatasource(dsName)
    + g.panel.barChart.queryOptions.withTargets([
      g.query.prometheus.new(dsName, query % { filters: strFilters })
      + g.query.prometheus.withIntervalFactor(1)
      + g.query.prometheus.withLegendFormat(legendFormat)
      + g.query.prometheus.withFormat(queryFormat)
      + g.query.prometheus.withInstant(instantQuery)
      + g.query.prometheus.withRange(range),
    ])
    + g.panel.barChart.options.withOrientation('horizontal')
    + g.panel.barChart.options.withXTickLabelSpacing(labelSpacing)
    + g.panel.barChart.options.withColorByField(legendFormat)
    + g.panel.barChart.options.tooltip.withMode(mode)
    + g.panel.barChart.options.tooltip.withSort(sort)
    + g.panel.barChart.fieldConfig.defaults.custom.withAxisGridShow(axisGridshow)
    + g.panel.barChart.fieldConfig.defaults.custom.withAxisPlacement(axisPlacement)
    + g.panel.barChart.standardOptions.color.withMode('thresholds')
    + g.panel.barChart.standardOptions.thresholds.withMode('absolute')
    + g.panel.barChart.standardOptions.thresholds.withSteps([{ color: 'green', value: null }])
    + g.panel.barChart.gridPos.withH(h)
    + g.panel.barChart.gridPos.withW(w),

  panel: barChartPanel,
}
