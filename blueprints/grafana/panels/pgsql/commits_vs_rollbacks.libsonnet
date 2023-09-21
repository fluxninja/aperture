local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters, title) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local targets = [
    g.query.prometheus.new(datasource.name, 'rate(postgresql_commits_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval])' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Commits'),

    g.query.prometheus.new(datasource.name, 'rate(postgresql_rollbacks_total{%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval])' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Rollbacks'),
  ],

  local commitVsRollback = timeSeriesPanel(title, datasource.name, '', stringFilters, h=8, w=10, targets=targets),
  panel: commitVsRollback.panel,
}
