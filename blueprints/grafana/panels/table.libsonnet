local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(title, dsName, query, x=0, h=20, w=10)
  g.panel.table.new(title)
  + g.panel.table.datasource.withType('prometheus')
  + g.panel.table.datasource.withUid(dsName)
  + g.panel.table.queryOptions.withTargets([
    g.query.prometheus.new(dsName, query)
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withRange(true),
  ])
  + g.panel.table.options.sortBy.withDesc(true)
  + g.panel.table.options.sortBy.withDisplayName(true)
  + g.panel.table.gridPos.withX(x)
  + g.panel.table.gridPos.withH(h)
  + g.panel.table.gridPos.withW(w)
