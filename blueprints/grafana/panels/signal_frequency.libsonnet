local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg, title, withVariables, signalName='') {
  local policyFilters = cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name },
  local signalFilters =
    if withVariables then
      policyFilters { signal_name: '${signal_name}', sub_circuit_id: '${sub_circuit_id}' }
    else
      policyFilters { signal_name: signalName },
  local stringFilters = utils.dictToPrometheusFilter(signalFilters),

  local targets = [
    g.query.prometheus.new(cfg.dashboard.datasource.name, 'avg(rate(signal_reading_count{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Valid'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(invalid_signal_readings_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Invalid'),
  ],

  local signalFrequency = timeSeriesPanel(title, cfg.dashboard.datasource.name, '', signalFilters, targets=targets),
  panel: signalFrequency.panel,
}
