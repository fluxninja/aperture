local timeSeriesPanel = import '../../panels/time-series.libsonnet';
local promUtils = import '../../utils/prometheus.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(datasourceName, policyName, extraFilters={})
  local signalFilters = extraFilters { policy_name: policyName, signal_name: '${signal_name}', sub_circuit_id: '${sub_circuit_id}' };
  local stringFilters = promUtils.dictToPrometheusFilter(signalFilters);

  local row1 = timeSeriesPanel(
    'Signal Average - ${signal_name} (${sub_circuit_id})',
    datasourceName,
    query='increase(signal_reading_sum{%(filters)s}[$__rate_interval]) / increase(signal_reading_count{%(filters)s}[$__rate_interval])' % { filters: stringFilters },
  );

  local targets = [
    g.query.prometheus.new(datasourceName, 'avg(rate(signal_reading_count{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Valid'),

    g.query.prometheus.new(datasourceName, 'sum(rate(invalid_signal_readings_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Invalid'),
  ];

  local row2 = timeSeriesPanel(
    'Signal Validity (Frequency) - ${signal_name} (${sub_circuit_id})',
    datasourceName,
    targets=targets
  );

  [
    row1,
    row2,
  ]
