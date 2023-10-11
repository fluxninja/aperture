local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local ready = statPanel('Ready Messages',
                          datasource,
                          'sum(rabbitmq_message_current{%(filters)s, infra_meter_name="%(infra_meter)s", state="ready"})' % { filters: stringFilters, infra_meter: infraMeterName },
                          stringFilters,
                          instantQuery=true,
                          range=false),

  local unacknowledged = statPanel('Unacknowledged Messages',
                                   datasource,
                                   'sum(rabbitmq_message_current{%(filters)s, infra_meter_name="%(infra_meter)s", state="unacknowledged"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                   stringFilters,
                                   instantQuery=true,
                                   panelColor='red',
                                   range=false),

  local incoming = statPanel('Incoming Messages',
                             datasource,
                             'sum(rate(rabbitmq_message_published_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__range]))' % { filters: stringFilters, infra_meter: infraMeterName },
                             stringFilters,
                             instantQuery=true,
                             range=false),

  local outgoing = statPanel('Outgoing Messages',
                             datasource,
                             'sum(rate(rabbitmq_message_delivered_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__range])) + sum(rate(rabbitmq_message_acknowledged_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__range]))' % { filters: stringFilters, infra_meter: infraMeterName },
                             stringFilters,
                             instantQuery=true,
                             range=false),

  local consumers = statPanel('Consumers',
                              datasource,
                              'sum(rabbitmq_consumer_count{%(filters)s, infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                              stringFilters,
                              instantQuery=true,
                              panelColor='blue',
                              range=false),

  local queues = statPanel('Queues',
                           datasource,
                           'count without(rabbitmq_queue_name) (rabbitmq_message_published_total{%(filters)s, infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                           stringFilters,
                           instantQuery=true,
                           panelColor='blue',
                           range=false),

  ready: ready.panel,
  unacknowledged: unacknowledged.panel,
  incoming: incoming.panel,
  outgoing: outgoing.panel,
  consumers: consumers.panel,
  queues: queues.panel,
}
