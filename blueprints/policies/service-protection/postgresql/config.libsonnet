local promqlDefaults = import '../promql/config.libsonnet';


promqlDefaults {
  policy+: {
    /**
    * @param (policy.promql_query: string) PromQL query to detect PostgreSQL overload.
    */
    promql_query: '(sum(postgresql_backends) / sum(postgresql_connection_max)) * 100',

    /**
    * @param (policy.postgresql: postgresql) Configuration for PostgreSQL OpenTelemetry receiver. Refer https://docs.fluxninja.com/integrations/metrics/postgresql for more information.
    * @schema (postgresql.username: string) Username of the PostgreSQL.
    * @schema (postgresql.password: string) Password of the PostgreSQL.
    * @schema (postgresql.endpoint: string) Endpoint of the PostgreSQL.
    * @schema (postgresql.transport: string) The transport protocol being used to connect to postgresql. Available options are tcp and unix.
    * @schema (postgresql.database: []string) The list of databases for which the receiver will attempt to collect statistics.
    * @schema (postgresql.collection_interval: string) This receiver collects metrics on an interval.
    * @schema (postgresql.initial_delay: string) Defines how long this receiver waits before starting.
    * @schema (postgresql.agent_group: string) Name of the Aperture Agent group.
    * @schema (postgresql.tls.insecure: bool) Whether to enable client transport security for the postgresql connection.
    * @schema (postgresql.tls.insecure_skip_verify: bool) Whether to validate server name and certificate if client transport security is enabled.
    * @schema (postgresql.tls.cert_file: string) A cerficate used for client authentication, if necessary.
    * @schema (postgresql.tls.key_file: string) An SSL key used for client authentication, if necessary.
    * @schema (postgresql.tls.ca_file: string) A set of certificate authorities used to validate the database server SSL certificate.
    */
    postgresql: {
      username: '__REQUIRED_FIELD__',
      password: '__REQUIRED_FIELD__',
      endpoint: '__REQUIRED_FIELD__',
      agent_group: 'default',
    },

    /**
    * @param (policy.service_protection_core.cpu_overload_confirmation.query_string: string) The Prometheus query to be run to get the PostgreSQL CPU utilization. Must return a scalar or a vector with a single element.
    * @param (policy.service_protection_core.cpu_overload_confirmation.threshold: float64) Threshold value for CPU utilizatio if it has to be used as overload confirmation.
    * @param (policy.service_protection_core.cpu_overload_confirmation.operator: string) The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`.
    */
    service_protection_core+: {
      cpu_overload_confirmation+: {
        query_string: 'avg(k8s_pod_cpu_utilization_ratio{k8s_statefulset_name="__REQUIRED_FIELD__"})',
        threshold: null,
        operator: 'gte',
      },
    },
  },

  dashboard+: {
    title: 'Aperture Service Protection for PostgreSQL',
    variant_name: 'PostgreSQL Overload Detection',
  },
}
