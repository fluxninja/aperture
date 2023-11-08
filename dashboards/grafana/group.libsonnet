local inframeterDashboardsFn = import 'dashboards/inframeter/dashboards.libsonnet';
local signalsDashboardFn = import 'dashboards/signals/dashboard.libsonnet';
local summaryDashboardFn = import 'dashboards/summary/dashboard.libsonnet';

function(policyFile, graph, policyName, datasource, extraFilters={}) {
  local graphObj = std.parseJson(graph),
  local summaryDashboard = summaryDashboardFn(graphObj, policyName, datasource, extraFilters),
  local signalsDashboard = signalsDashboardFn(policyName, datasource, extraFilters),
  local inframeterDashboards = inframeterDashboardsFn(policyFile, policyName, datasource, extraFilters),

  dashboards: {
    summary: summaryDashboard,
    signals: signalsDashboard,
  } + inframeterDashboards,
}
