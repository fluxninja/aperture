local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local ready = statPanel('Ready Messages',
                          datasource.name,
                          'sum(rabbitmq_message_current{%(filters)s, state="ready"})' % { filters: stringFilters },
                          stringFilters,
                          instantQuery=true,
                          range=false),

  local unacknowledged = statPanel('Unacknowledged Messages',
                                   datasource.name,
                                   'sum(rabbitmq_message_current{%(filters)s, state="unacknowledged"})' % { filters: stringFilters },
                                   stringFilters,
                                   instantQuery=true,
                                   panelColor='red',
                                   range=false),

  local incoming = statPanel('Rate of Incoming Messages',
                             datasource.name,
                             'sum(rate(rabbitmq_message_published_total{%(filters)s}[$__range]))' % { filters: stringFilters },
                             stringFilters,
                             instantQuery=true,
                             range=false),

  local outgoing = statPanel('Rate of Outgoing Messages',
                             datasource.name,
                             'sum(rate(rabbitmq_message_delivered_total{%(filters)s}[$__range])) + sum(rate(rabbitmq_message_acknowledged_total{%(filters)s}[$__range]))' % { filters: stringFilters },
                             stringFilters,
                             instantQuery=true,
                             range=false),

  local consumers = statPanel('Consumers',
                              datasource.name,
                              'sum(rabbitmq_consumer_count{%(filters)s})' % { filters: stringFilters },
                              stringFilters,
                              instantQuery=true,
                              panelColor='blue',
                              range=false),

  local queues = statPanel('Queues',
                           datasource.name,
                           'count without(rabbitmq_queue_name) (rabbitmq_message_published_total{%(filters)s})' % { filters: stringFilters },
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
