local promqlDefaults = import '../promql/config.libsonnet';

/**
* @param (policy.policy_name: string required) Name of the policy.
* @param (policy.promql_query: string) PromQL query to detect PostgreSQL overload.
* @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
* @param (policy.resources: aperture.spec.v1.Resources) Additional resources.
* @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
* @param (policy.service_protection_core.overload_confirmations: []overload_confirmation) List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.
* @schema (overload_confirmation.query_string: string required) The Prometheus query to be run. Must return a scalar or a vector with a single element.
* @schema (overload_confirmation.threshold: float64) The threshold for the overload confirmation criteria.
* @schema (overload_confirmation.operator: string) The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`
* @param (policy.service_protection_core.adaptive_load_scheduler: aperture.spec.v1.AdaptiveLoadSchedulerParameters required) Parameters for Adaptive Load Scheduler.
* @param (policy.service_protection_core.dry_run: bool) Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.
*/

promqlDefaults {
  policy+: {
    promql_query: '(sum(postgresql_backends) / sum(postgresql_connection_max)) * 100',
    /**
    * @param (policy.setpoint: float64 required) Setpoint.
    */
    setpoint: '__REQUIRED_FIELD__',

    /**
    * @param (policy.postgresql: postgresql required) Configuration for PostgreSQL OpenTelemetry receiver. Refer https://docs.fluxninja.com/integrations/metrics/postgresql for more information.
    * @schema (postgresql.username: string required) Username of the PostgreSQL.
    * @schema (postgresql.password: string required) Password of the PostgreSQL.
    * @schema (postgresql.endpoint: string required) Endpoint of the PostgreSQL.
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

  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) Time from of dashboard.
  * @param (dashboard.time_to: string) Time to of dashboard.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard: {
    refresh_interval: '15s',
    time_from: 'now-15m',
    time_to: 'now',
    extra_filters: {},
    title: 'Aperture Service Protection for PostgreSQL',
    /**
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
    variant_name: 'PostgreSQL Overload Detection',
  },
}
