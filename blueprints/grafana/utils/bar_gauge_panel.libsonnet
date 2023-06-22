local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, h=6, w=12) {
  local barGaugePanel =
    g.panel.barGauge.new(title)
    + g.panel.barGauge.datasource.withType('prometheus')
    + g.panel.barGauge.datasource.withUid(dsName)
    + g.panel.barGauge.queryOptions.withTargets([
      g.query.prometheus.new(dsName, query % { filters: strFilters })
      + g.query.prometheus.withIntervalFactor(1)
      + g.query.prometheus.withLegendFormat('{{ instance }} - {{ policy_name }}')
      + g.query.prometheus.withFormat('time_series')
      + g.query.prometheus.withInstant(false)
      + g.query.prometheus.withRange(true),
    ])
    + g.panel.barGauge.options.withDisplayMode('gradient')
    + g.panel.barGauge.options.withOrientation('horizontal')
    + g.panel.barGauge.standardOptions.color.withMode('thresholds')
    + g.panel.barGauge.standardOptions.tresholds.withMode('absolute')
    + g.panel.barGauge.standardOptions.tresholds.withSteps([{ color: 'green', value: null }])
    + g.panel.barGauge.gridPos.withH(h)
    + g.panel.barGauge.gridPos.withW(w),

  panel: barGaugePanel,
}
