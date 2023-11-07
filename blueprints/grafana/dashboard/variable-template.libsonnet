local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';
local var = g.dashboard.variable;

function(name, query, dsName, label) {
  local variable =
    var.query.new(name, query)
    + var.query.withDatasource('prometheus', dsName)
    + var.query.generalOptions.withLabel(label)
    + var.query.selectionOptions.withIncludeAll(false)
    + var.query.selectionOptions.withMulti(false)
    + var.query.withSort(1),

  variable: variable,
}
