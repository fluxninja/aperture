local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters, title) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local targets = [
    g.query.prometheus.new(datasource.name, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="ins"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Insert'),

    g.query.prometheus.new(datasource.name, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="del"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Delete'),

    g.query.prometheus.new(datasource.name, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="upd"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Update'),

    g.query.prometheus.new(datasource.name, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="hot_upd"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Hot Update'),
  ],

  local operations = timeSeriesPanel(title, datasource.name, '', stringFilters, h=8, w=10, targets=targets),
  panel: operations.panel,
}
