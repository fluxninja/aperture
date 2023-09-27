local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(title, dsName, query, strFilters, h=20, w=10) {
  local tablePanel =
    g.panel.table.new(title)
    + g.panel.table.queryOptions.withDatasource(dsName)
    + g.panel.table.queryOptions.withTargets([
      g.query.prometheus.new(dsName, query % { filters: strFilters })
      + g.query.prometheus.withIntervalFactor(1)
      + g.query.prometheus.withRange(true),
    ])
    + g.panel.table.options.sortBy.withDesc(true)
    + g.panel.table.options.sortBy.withDisplayName(true)
    + g.panel.table.gridPos.withH(h)
    + g.panel.table.gridPos.withW(w)
    + g.panel.table.options.withFieldOptions({
      overrides: [
        {
          matcher: g.panel.table.options.fieldMatcher.withId('byName').withOptions('postgresql_table_name'),
          properties: [
            g.panel.table.options.property.withId('displayName').withValue('Table Name'),
          ],
        },
        {
          matcher: g.panel.table.options.fieldMatcher.withId('byName').withOptions('Value'),
          properties: [
            g.panel.table.options.property.withId('displayName').withValue('Size (Bytes)'),
          ],
        },
      ],
    }),
  panel: tablePanel,
}
