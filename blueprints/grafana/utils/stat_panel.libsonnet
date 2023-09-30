local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, h=6, w=6, instantQuery=false, range=true, panelColor='green', graphMode='none') {
  local statPanel =
    g.panel.stat.new(title)
    + g.panel.stat.datasource.withType('prometheus')
    + g.panel.stat.datasource.withUid(dsName)
    + g.panel.stat.queryOptions.withTargets([
      g.query.prometheus.new(dsName, query % { filters: strFilters })
      + g.query.prometheus.withIntervalFactor(1)
      + g.query.prometheus.withLegendFormat('{{ instance }} - {{ policy_name }}')
      + g.query.prometheus.withInstant(instantQuery)
      + g.query.prometheus.withRange(range),
    ])
    + g.panel.stat.standardOptions.color.withMode('thresholds')
    + g.panel.stat.standardOptions.thresholds.withMode('absolute')
    + g.panel.stat.standardOptions.thresholds.withSteps([{ color: panelColor, value: null }])
    + g.panel.stat.options.withColorMode('value')
    + g.panel.stat.options.withGraphMode(graphMode)
    + g.panel.stat.options.withJustifyMode('center')
    + g.panel.stat.options.withOrientation('horizontal')
    + g.panel.stat.options.withTextMode('auto')
    + g.panel.stat.options.reduceOptions.withCalcs(['lastNotNull'])
    + g.panel.stat.options.reduceOptions.withFields('')
    + g.panel.stat.options.reduceOptions.withValues(false)
    + g.panel.stat.gridPos.withH(h)
    + g.panel.stat.gridPos.withW(w),

  panel: statPanel,
}
