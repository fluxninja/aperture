local promqlDefaults = import '../promql/config.libsonnet';

promqlDefaults {
  policy+: {
    /**
    * @param (policy.promql_query: string) PromQL query to detect ElasticSearch overload.
    */
    promql_query: 'elasticsearch_node_thread_pool_tasks_queued{thread_pool_name="search"}',
  },

  dashboard+: {
    title: 'Aperture Service Protection for Elasticsearch',
    variant_name: 'Elasticsearch Overload Detection',
  },
}
