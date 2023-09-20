local commonConfig = import '../common-aiad/config-defaults.libsonnet';

commonConfig {
  policy+: {
    /**
    * @param (policy.postgresql: postgresql) Configuration for PostgreSQL OpenTelemetry receiver. Refer https://docs.fluxninja.com/integrations/metrics/postgresql for more information.
    * @schema (postgresql.username: string) Username of the PostgreSQL.
    * @schema (postgresql.password: string) Password of the PostgreSQL.
    * @schema (postgresql.endpoint: string) Endpoint of the PostgreSQL.
    * @schema (postgresql.transport: string) The transport protocol being used to connect to postgresql. Available options are tcp and unix.
    * @schema (postgresql.databases: []string) The list of databases for which the receiver will attempt to collect statistics.
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
  },

  dashboard+: {
    title: 'Aperture Service Protection for PostgreSQL',
    variant_name: 'PostgreSQL Overload Detection',
  },
}
