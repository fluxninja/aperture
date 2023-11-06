local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';
local var = g.dashboard.variable;

function(title, refresh_interval) {
  local datasourceVar =
    var.datasource.new('datasource', 'prometheus')
    + var.datasource.generalOptions.withLabel('Data Source')
    + var.datasource.selectionOptions.withMulti(false)
    + var.datasource.selectionOptions.withIncludeAll(false),

  local dashboardDef =
    g.dashboard.new(title)
    + g.dashboard.time.withFrom('now-15m')
    + g.dashboard.withTimezone('browser')
    + g.dashboard.withRefresh(refresh_interval)
    + g.dashboard.withVariables([datasourceVar]),

  baseDashboard: dashboardDef,
}
