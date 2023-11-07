local statPanel = import '../../../panels/stat.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, infra_meter_name: infraMeterName }),

  local ready = statPanel('Ready Messages',
                          datasource,
                          'sum(rabbitmq_message_current{%(filters)s", state="ready"})' % { filters: stringFilters },
                          instantQuery=true,
                          range=false),

  local unacknowledged = statPanel('Unacknowledged Messages',
                                   datasource,
                                   'sum(rabbitmq_message_current{%(filters)s", state="unacknowledged"})' % { filters: stringFilters },
                                   instantQuery=true,
                                   panelColor='red',
                                   range=false),

  local incoming = statPanel('Messages Published per Second',
                             datasource,
                             'sum(rate(rabbitmq_message_published_total{%(filters)s"}[$__range]))' % { filters: stringFilters },
                             instantQuery=true,
                             range=false),

  local outgoing = statPanel('Messages Delivered per Second',
                             datasource,
                             'sum(rate(rabbitmq_message_delivered_total{%(filters)s"}[$__range]))' % { filters: stringFilters },
                             instantQuery=true,
                             range=false),

  local consumers = statPanel('Consumers',
                              datasource,
                              'sum(rabbitmq_consumer_count{%(filters)s"})' % { filters: stringFilters },
                              instantQuery=true,
                              panelColor='blue',
                              range=false),

  local queues = statPanel('Messages in Queue',
                           datasource,
                           'sum(rabbitmq_message_current_total{%(filters)s"})' % { filters: stringFilters },
                           instantQuery=true,
                           panelColor='blue',
                           range=false),

  ready: ready,
  unacknowledged: unacknowledged,
  incoming: incoming,
  outgoing: outgoing,
  consumers: consumers,
  queues: queues,
}
