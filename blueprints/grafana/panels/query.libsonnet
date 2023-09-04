local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),
  local signalName = cfg.component_body.query.promql.out_ports.output.signal_name,
  local query = cfg.component_body.query.promql.query_string,

  local queryPanel = timeSeriesPanel('Query for ' + signalName, cfg.dashboard.datasource.name, query, stringFilters),

  panel: queryPanel.panel,
}
