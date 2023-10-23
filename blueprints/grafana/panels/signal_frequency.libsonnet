local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(datasourceName, policyName, _, extraFilters, title, withVariables, signalName='') {
  local policyFilters = extraFilters { policy_name: policyName },
  local signalFilters =
    if withVariables then
      policyFilters { signal_name: '${signal_name}', sub_circuit_id: '${sub_circuit_id}' }
    else
      policyFilters { signal_name: signalName },
  local stringFilters = utils.dictToPrometheusFilter(signalFilters),

  local targets = [
    g.query.prometheus.new(datasourceName, 'avg(rate(signal_reading_count{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Valid'),

    g.query.prometheus.new(datasourceName, 'sum(rate(invalid_signal_readings_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Invalid'),
  ],

  local signalFrequency = timeSeriesPanel(title, datasourceName, '', stringFilters, targets=targets),
  panel: signalFrequency.panel,
}
