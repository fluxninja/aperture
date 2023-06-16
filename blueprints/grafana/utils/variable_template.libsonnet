local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(name, query, dsName, label) {
  local variable =
    g.var.query.new(name, query)
    + g.var.query.queryTypes.withLabelValues(
      'up',
      'instance',
    )
    + g.var.query.withDatasource('prometheus', dsName)
    + g.var.query.generalOptions.withLabel(label)
    + g.var.query.selectionOptions.withIncludeAll(false)
    + g.var.query.selectionOptions.withMulti(false)
    + g.var.query.withSort(),

  variable: variable,
}
