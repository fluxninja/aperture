local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, h=6, w=6, targets=[]) {
  local barChart =
    g.panel.barChart.new(title)
    + g.panel.barChart.queryOptions.withDatasource(dsName)
    + g.panel.barChart.queryOptions.withTargets([
      g.query.prometheus.new(dsName, query % { filters: strFilters })
      + g.query.prometheus.withIntervalFactor(1)
      + g.query.prometheus.withLegendFormat('{{ instance }} - {{ policy_name }}')
      + g.query.prometheus.withRange(true),
    ])
    + g.panel.barChart.standardOptions.color.withMode('thresholds')
    + g.panel.barChart.standardOptions.thresholds.withMode('absolute')
    //+ g.panel.barChart.standardOptions.thresholds.withSteps([{ color: 'green', value: null }])
    + g.panel.barChart.options.withOrientation('vertical')
    + g.panel.barChart.options.withText('auto')
    + g.panel.barChart.gridPos.withH(h)
    + g.panel.barChart.gridPos.withW(w),

  local withMultipleTargets =
    if targets != []
    then
      barChart + g.panel.barChart.withTargets(targets)
    else
      barChart + g.panel.barChart.withTargets([
        g.query.prometheus.new(dsName, query % { filters: strFilters })
        + g.query.prometheus.withIntervalFactor(1),
      ]),

  panel: withMultipleTargets,
}
