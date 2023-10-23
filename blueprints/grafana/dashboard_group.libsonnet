local infraMetersDashboards = import 'infra_meter_dashboard.libsonnet';
local summaryDashboard = import 'summary_dashboard.libsonnet';

function(policyFile, componentsList, policyName, datasource, extraFilters={}) {
  local summary = summaryDashboard(componentsList, policyName, datasource, extraFilters).dashboard,
  local receivers = infraMetersDashboards(policyFile, policyName, datasource, extraFilters).dashboards,

  dashboards: {
    summary: summary,
  } + receivers,
}
