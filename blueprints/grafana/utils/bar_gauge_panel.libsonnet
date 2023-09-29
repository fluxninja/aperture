local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, h=10, w=24, legendFormat=null, values=false) {
  local barGaugePanel =
    g.panel.barGauge.new(title)
    + g.panel.barGauge.datasource.withType('prometheus')
    + g.panel.barGauge.datasource.withUid(dsName)
    + g.panel.barGauge.queryOptions.withTargets([
      g.query.prometheus.new(dsName, query % { filters: strFilters })
      + g.query.prometheus.withIntervalFactor(1)
      + g.query.prometheus.withLegendFormat(legendFormat)
      + g.query.prometheus.withFormat('time_series')
      + g.query.prometheus.withInstant(false)
      + g.query.prometheus.withRange(true),
    ])
    + g.panel.barGauge.options.withDisplayMode('gradient')
    + g.panel.barGauge.options.withOrientation('horizontal')
    + g.panel.barGauge.options.reduceOptions.withValues(values)
    + g.panel.barGauge.standardOptions.color.withMode('thresholds')
    + g.panel.barGauge.standardOptions.thresholds.withMode('absolute')
    + g.panel.barGauge.standardOptions.thresholds.withSteps([{ color: 'green', value: null }])
    + g.panel.barGauge.gridPos.withH(h)
    + g.panel.barGauge.gridPos.withW(w),

  panel: barGaugePanel,
}
