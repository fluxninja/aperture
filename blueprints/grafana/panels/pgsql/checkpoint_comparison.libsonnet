local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters, title) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local targets = [
    g.query.prometheus.new(datasource.name, 'sum(rate(postgresql_bgwriter_checkpoint_count_total{%(filters)s,infra_meter_name="%(infra_meter)s",type="requested"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Requested'),

    g.query.prometheus.new(datasource.name, 'sum(rate(postgresql_bgwriter_checkpoint_count_total{%(filters)s,infra_meter_name="%(infra_meter)s",type="scheduled"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Scheduled'),
  ],

  local checkpointComparison = timeSeriesPanel(title, datasource.name, '', stringFilters, targets=targets),
  panel: checkpointComparison.panel,
}
