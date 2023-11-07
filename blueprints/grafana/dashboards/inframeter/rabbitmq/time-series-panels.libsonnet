local timeSeriesPanel = import '../../../panels/time-series.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, infra_meter_name: infraMeterName }),

  local readyQuery = g.query.prometheus.new(datasource, 'sum(rabbitmq_message_current{%(filters)s, state="ready"})' % { filters: stringFilters })
                     + g.query.prometheus.withIntervalFactor(1)
                     + g.query.prometheus.withLegendFormat('Ready'),
  local ready = timeSeriesPanel('Messages Ready For Consumers', datasource, axisLabel='Messages', unit='short', targets=[readyQuery]),

  local acknowledgedQuery = g.query.prometheus.new(datasource, 'sum(rate(rabbitmq_message_acknowledged_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
                            + g.query.prometheus.withIntervalFactor(1)
                            + g.query.prometheus.withLegendFormat('Acknowledged'),
  local acknowledged = timeSeriesPanel('Messages Acknowledged By Consumers Per Second', datasource, axisLabel='Messages/Second', unit='mps', targets=[acknowledgedQuery]),

  local unacknowledgedQuery = g.query.prometheus.new(datasource, 'sum(rabbitmq_message_current{%(filters)s, state="unacknowledged"})' % { filters: stringFilters })
                              + g.query.prometheus.withIntervalFactor(1)
                              + g.query.prometheus.withLegendFormat('Unacknowledged'),
  local unacknowledged = timeSeriesPanel('Messages Unacknowledged By Consumers', datasource, axisLabel='Messages', unit='short', targets=[unacknowledgedQuery]),

  local readyQueryPerVhost = g.query.prometheus.new(datasource, 'sum by(rabbitmq_vhost_name) (rabbitmq_message_current{%(filters)s, state="ready"})' % { filters: stringFilters })
                             + g.query.prometheus.withIntervalFactor(1)
                             + g.query.prometheus.withLegendFormat('{{rabbitmq_vhost_name}}'),
  local readyPerVhost = timeSeriesPanel('Messages Ready For Consumers Per Second by Vhost', datasource, axisLabel='Messages', unit='short', targets=[readyQueryPerVhost]),

  local acknowledgedQueryPerVhost = g.query.prometheus.new(datasource, 'sum by(rabbitmq_vhost_name) (rate(rabbitmq_message_acknowledged_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
                                    + g.query.prometheus.withIntervalFactor(1)
                                    + g.query.prometheus.withLegendFormat('{{rabbitmq_vhost_name}}'),
  local acknowledgedPerVhost = timeSeriesPanel('Messages Acknowledged By Consumers Per Second by Vhost', datasource, axisLabel='Messages/Second', unit='mps', targets=[acknowledgedQueryPerVhost]),

  local unacknowledgedQueryPerVhost = g.query.prometheus.new(datasource, 'sum by(rabbitmq_vhost_name) (rabbitmq_message_current{%(filters)s, state="unacknowledged"})' % { filters: stringFilters })
                                      + g.query.prometheus.withIntervalFactor(1)
                                      + g.query.prometheus.withLegendFormat('{{rabbitmq_vhost_name}}'),
  local unacknowledgedPerVhost = timeSeriesPanel('Messages Unacknowledged By Consumers Per Vhost', datasource, axisLabel='Messages', unit='short', targets=[unacknowledgedQueryPerVhost]),

  local queuesGrowthQuery = g.query.prometheus.new(datasource, 'sum by (rabbitmq_queue_name) (rate(rabbitmq_message_published_total{%(filters)s}[$__rate_interval])) - sum by (rabbitmq_queue_name) (rate(rabbitmq_message_acknowledged_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
                            + g.query.prometheus.withIntervalFactor(1)
                            + g.query.prometheus.withLegendFormat('{{rabbitmq_queue_name}}'),
  local queuesGrowth = timeSeriesPanel('Queue Growth Per Second', datasource, axisLabel='Messages/Second', unit='mps', targets=[queuesGrowthQuery]),

  local publishedQuery = g.query.prometheus.new(datasource, 'sum(rate(rabbitmq_message_published_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
                         + g.query.prometheus.withIntervalFactor(1)
                         + g.query.prometheus.withLegendFormat('Published'),
  local published = timeSeriesPanel('Messages Published Per Second', datasource, axisLabel='Messages/Second', unit='mps', targets=[publishedQuery]),

  local deliveredQuery = g.query.prometheus.new(datasource, 'sum(rate(rabbitmq_message_delivered_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
                         + g.query.prometheus.withIntervalFactor(1)
                         + g.query.prometheus.withLegendFormat('Delivered'),
  local delivered = timeSeriesPanel('Messages Delivered Per Second', datasource, axisLabel='Messages/Second', unit='mps', targets=[deliveredQuery]),

  ready: ready,
  acknowledged: acknowledged,
  unacknowledged: unacknowledged,
  readyPerVhost: readyPerVhost,
  acknowledgedPerVhost: acknowledgedPerVhost,
  unacknowledgedPerVhost: unacknowledgedPerVhost,
  queuesGrowth: queuesGrowth,
  published: published,
  delivered: delivered,
}
