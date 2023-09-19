local serviceProtectionDefaults = import '../common-aiad/config-defaults.libsonnet';

serviceProtectionDefaults {
  policy+: {
    /**
    * @param (policy.elasticsearch: elasticsearch) Configuration for Elasticsearch OpenTelemetry receiver. Refer https://docs.fluxninja.com/integrations/metrics/elasticsearch for more information.
    * @schema (elasticsearch.username: string) Username of the Elasticsearch.
    * @schema (elasticsearch.password: string) Password of the Elasticsearch.
    * @schema (elasticsearch.endpoint: string) Endpoint of the Elasticsearch.
    * @schema (elasticsearch.agent_group: string) Name of the Aperture Agent group.
    * @schema (elasticsearch.nodes: []string) Node filters that define which nodes are scraped for node-level and cluster-level metrics.
    * @schema (elasticsearch.indices: []string) Index filters that define which indices are scraped for index-level metrics.
    * @schema (elasticsearch.skip_cluster_metrics: bool) If true, cluster-level metrics will not be scraped.
    * @schema (elasticsearch.collection_interval: string) This receiver collects metrics on an interval.
    * @schema (elasticsearch.initial_delay: string) Defines how long this receiver waits before starting.
    */
    elasticsearch: {
      username: '__REQUIRED_FIELD__',
      password: '__REQUIRED_FIELD__',
      endpoint: '__REQUIRED_FIELD__',
      agent_group: 'default',
    },
  },

  dashboard+: {
    title: 'Aperture Service Protection for Elasticsearch',
    variant_name: 'Elasticsearch Overload Detection',
  },
}
