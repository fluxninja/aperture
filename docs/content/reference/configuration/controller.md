---
title: Aperture Controller Configuration Reference
sidebar_position: 10
sidebar_label: Controller
---

:::info
See also [Aperture Controller installation](/get-started/installation/controller/controller.md).
:::

List of all config parameters for Aperture Controller.

<!---
Generated File Starts
-->

## Table of contents

### COMMON CONFIGURATION

| Key          | Reference                          |
| ------------ | ---------------------------------- |
| `client`     | [Client](#client)                  |
| `etcd`       | [Etcd](#etcd)                      |
| `liveness`   | [Liveness](#liveness)              |
| `log`        | [Log](#log)                        |
| `metrics`    | [Metrics](#metrics)                |
| `plugins`    | [Plugins](#plugins)                |
| `policies`   | [PoliciesConfig](#policies-config) |
| `profilers`  | [Profilers](#profilers)            |
| `prometheus` | [Prometheus](#prometheus)          |
| `readiness`  | [Readiness](#readiness)            |
| `server`     | [Server](#server)                  |
| `watchdog`   | [Watchdog](#watchdog)              |

### CONTROLLER CONFIGURATION

| Key    | Reference        |
| ------ | ---------------- |
| `otel` | [OTEL](#o-t-e-l) |

### PLUGIN CONFIGURATION

| Key                | Reference                             |
| ------------------ | ------------------------------------- |
| `fluxninja_plugin` | [FluxNinjaPlugin](#flux-ninja-plugin) |
| `sentry_plugin`    | [SentryPlugin](#sentry-plugin)        |

## Reference

### _Client_ {#client}

Key: `client`

Env-Var Prefix: `APERTURE_CONTROLLER_CLIENT_`

#### Members

<dl>

<dt>proxy</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_CLIENT_PROXY_`
Type: [ProxyConfig](#proxy-config)

</dd>

</dl>

### _Etcd_ {#etcd}

Key: `etcd`

Env-Var Prefix: `APERTURE_CONTROLLER_ETCD_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_ETCD_`
Type: [EtcdConfig](#etcd-config)

</dd>

</dl>

### _FluxNinjaPlugin_ {#flux-ninja-plugin}

Key: `fluxninja_plugin`

Env-Var Prefix: `APERTURE_CONTROLLER_FLUXNINJA_PLUGIN_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_FLUXNINJA_PLUGIN_`
Type: [FluxNinjaPluginConfig](#flux-ninja-plugin-config)

</dd>

</dl>

### _Liveness_ {#liveness}

Key: `liveness`

Env-Var Prefix: `APERTURE_CONTROLLER_LIVENESS_`

#### Members

<dl>

<dt>scheduler</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_LIVENESS_SCHEDULER_`
Type: [JobGroupConfig](#job-group-config)

</dd>

<dt>service</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_LIVENESS_SERVICE_`
Type: [JobConfig](#job-config)

</dd>

</dl>

### _Log_ {#log}

Key: `log`

Env-Var Prefix: `APERTURE_CONTROLLER_LOG_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_LOG_`
Type: [LogConfig](#log-config)

</dd>

</dl>

### _Metrics_ {#metrics}

Key: `metrics`

Env-Var Prefix: `APERTURE_CONTROLLER_METRICS_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_METRICS_`
Type: [MetricsConfig](#metrics-config)

</dd>

</dl>

### _OTEL_ {#o-t-e-l}

Key: `otel`

Env-Var Prefix: `APERTURE_CONTROLLER_OTEL_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_OTEL_`
Type: [ControllerOTELConfig](#controller-o-t-e-l-config)

</dd>

</dl>

### _Plugins_ {#plugins}

Key: `plugins`

Env-Var Prefix: `APERTURE_CONTROLLER_PLUGINS_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_PLUGINS_`
Type: [PluginsConfig](#plugins-config)

</dd>

</dl>

### _PoliciesConfig_ {#policies-config}

Key: `policies`

Env-Var Prefix: `APERTURE_CONTROLLER_POLICIES_`

#### Members

<dl>

<dt>promql_jobs_scheduler</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_POLICIES_PROMQL_JOBS_SCHEDULER_`
Type: [JobGroupConfig](#job-group-config)

</dd>

</dl>

### _Profilers_ {#profilers}

Key: `profilers`

Env-Var Prefix: `APERTURE_CONTROLLER_PROFILERS_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_PROFILERS_`
Type: [ProfilersConfig](#profilers-config)

</dd>

</dl>

### _Prometheus_ {#prometheus}

Key: `prometheus`

Env-Var Prefix: `APERTURE_CONTROLLER_PROMETHEUS_`

#### Members

<dl>

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_PROMETHEUS_`
Type: [PrometheusConfig](#prometheus-config)

</dd>

<dt>http_client</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_PROMETHEUS_HTTP_CLIENT_`
Type: [HTTPClientConfig](#http-client-config)

</dd>

</dl>

### _Readiness_ {#readiness}

Key: `readiness`

Env-Var Prefix: `APERTURE_CONTROLLER_READINESS_`

#### Members

<dl>

<dt>scheduler</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_READINESS_SCHEDULER_`
Type: [JobGroupConfig](#job-group-config)

</dd>

<dt>service</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_READINESS_SERVICE_`
Type: [JobConfig](#job-config)

</dd>

</dl>

### _SentryPlugin_ {#sentry-plugin}

Key: `sentry_plugin`

Env-Var Prefix: `APERTURE_CONTROLLER_SENTRY_PLUGIN_`

#### Members

<dl>

<dt>sentry</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_SENTRY_PLUGIN_SENTRY_`
Type: [SentryConfig](#sentry-config)

</dd>

</dl>

### _Server_ {#server}

Key: `server`

Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_`

#### Members

<dl>

<dt>grpc</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_GRPC_`
Type: [GRPCServerConfig](#g-rpc-server-config)

</dd>

<dt>grpc_gateway</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_GRPC_GATEWAY_`
Type: [GRPCGatewayConfig](#g-rpc-gateway-config)

</dd>

<dt>http</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_HTTP_`
Type: [HTTPServerConfig](#http-server-config)

</dd>

<dt>listener</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_LISTENER_`
Type: [ListenerConfig](#listener-config)

</dd>

<dt>tls</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_TLS_`
Type: [ServerTLSConfig](#server-tls-config)

</dd>

</dl>

### _Watchdog_ {#watchdog}

Key: `watchdog`

Env-Var Prefix: `APERTURE_CONTROLLER_WATCHDOG_`

#### Members

<dl>

<dt>memory</dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_WATCHDOG_MEMORY_`
Type: [WatchdogConfig](#watchdog-config)

</dd>

</dl>

## Objects

### AdaptivePolicy {#adaptive-policy}

AdaptivePolicy creates a policy that forces GC when the usage surpasses the configured factor of the available memory. This policy calculates next target as usage+(limit-usage)\*factor.

#### Properties

<dl>
<dt>enabled</dt>
<dd>

(bool, default: `false`) Flag to enable the policy

</dd>
<dt>factor</dt>
<dd>

(float64, `gte=0,lte=1`, default: `0.50`) Factor sets user-configured limit of available memory

</dd>
</dl>

### BackoffConfig {#backoff-config}

BackoffConfig holds configuration for GRPC Client Backoff.

#### Properties

<dl>
<dt>base_delay</dt>
<dd>

(string, `gte=0`, default: `1s`) Base Delay

</dd>
<dt>jitter</dt>
<dd>

(float64, `gte=0`, default: `0.2`) Jitter

</dd>
<dt>max_delay</dt>
<dd>

(string, `gte=0`, default: `120s`) Max Delay

</dd>
<dt>multiplier</dt>
<dd>

(float64, `gte=0`, default: `1.6`) Backoff multiplier

</dd>
</dl>

### BatchAlertsConfig {#batch-alerts-config}

BatchAlertsConfig defines configuration for OTEL batch processor.

#### Properties

<dl>
<dt>send_batch_max_size</dt>
<dd>

(uint32, `gte=0`, default: `100`) SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
into smaller units.

</dd>
<dt>send_batch_size</dt>
<dd>

(uint32, `gt=0`, default: `100`) SendBatchSize is the size of a batch which after hit, will trigger it to be sent.

</dd>
<dt>timeout</dt>
<dd>

(string, `gt=0`, default: `1s`) Timeout sets the time after which a batch will be sent regardless of size.

</dd>
</dl>

### ClientConfig {#client-config}

ClientConfig is the client configuration.

#### Properties

<dl>
<dt>grpc</dt>
<dd>

([GRPCClientConfig](#g-rpc-client-config))

</dd>
<dt>http</dt>
<dd>

([HTTPClientConfig](#http-client-config))

</dd>
</dl>

### ClientTLSConfig {#client-tls-config}

ClientTLSConfig is the config for client TLS.

#### Properties

<dl>
<dt>ca_file</dt>
<dd>

(string)

</dd>
<dt>cert_file</dt>
<dd>

(string)

</dd>
<dt>insecure_skip_verify</dt>
<dd>

(bool)

</dd>
<dt>key_file</dt>
<dd>

(string)

</dd>
<dt>key_log_file</dt>
<dd>

(string)

</dd>
</dl>

### ControllerOTELConfig {#controller-o-t-e-l-config}

ControllerOTELConfig is the configuration for Agent's OTEL collector.

#### Properties

<dl>
<dt>batch_alerts</dt>
<dd>

([BatchAlertsConfig](#batch-alerts-config))

</dd>
<dt>ports</dt>
<dd>

([PortsConfig](#ports-config))

</dd>
</dl>

### EtcdConfig {#etcd-config}

EtcdConfig holds configuration for etcd client.

#### Properties

<dl>
<dt>endpoints</dt>
<dd>

([]string, **required**, `required,gt=0,dive,hostname_port|url|fqdn`) List of Etcd server endpoints

</dd>
<dt>lease_ttl</dt>
<dd>

(string, `gte=1s`, default: `60s`) Lease time-to-live

</dd>
<dt>password</dt>
<dd>

(string)

</dd>
<dt>username</dt>
<dd>

(string) Authentication

</dd>
<dt>tls</dt>
<dd>

([ClientTLSConfig](#client-tls-config))

</dd>
</dl>

### FluxNinjaPluginConfig {#flux-ninja-plugin-config}

FluxNinjaPluginConfig is the configuration for FluxNinja ARC integration plugin.

#### Properties

<dl>
<dt>api_key</dt>
<dd>

(string) API Key for this agent.

</dd>
<dt>fluxninja_endpoint</dt>
<dd>

(string, `omitempty,hostname_port|url|fqdn`) Address to grpc or http(s) server listening in agent service. To use http protocol, the address must start with http(s)://.

</dd>
<dt>heartbeat_interval</dt>
<dd>

(string, `gte=0s`, default: `5s`) Interval between each heartbeat.

</dd>
<dt>client</dt>
<dd>

([ClientConfig](#client-config))

</dd>
</dl>

### GRPCClientConfig {#g-rpc-client-config}

GRPCClientConfig holds configuration for GRPC Client.

#### Properties

<dl>
<dt>insecure</dt>
<dd>

(bool, default: `false`) Disable ClientTLS

</dd>
<dt>min_connection_timeout</dt>
<dd>

(string, `gte=0`, default: `20s`) Minimum connection timeout

</dd>
<dt>use_proxy</dt>
<dd>

(bool, default: `false`) Use HTTP CONNECT Proxy

</dd>
<dt>backoff</dt>
<dd>

([BackoffConfig](#backoff-config))

</dd>
<dt>tls</dt>
<dd>

([ClientTLSConfig](#client-tls-config))

</dd>
</dl>

### GRPCGatewayConfig {#g-rpc-gateway-config}

GRPCGatewayConfig holds configuration for grpc-http gateway

#### Properties

<dl>
<dt>grpc_server_address</dt>
<dd>

(string, `hostname_port`, default: `0.0.0.0:1`) GRPC server address to connect to - By default it points to HTTP server port because FluxNinja stack runs GRPC and HTTP servers on the same port

</dd>
</dl>

### GRPCServerConfig {#g-rpc-server-config}

GRPCServerConfig holds configuration for GRPC Server.

#### Properties

<dl>
<dt>connection_timeout</dt>
<dd>

(string, `gte=0s`, default: `120s`) Connection timeout

</dd>
<dt>enable_reflection</dt>
<dd>

(bool, default: `false`) Enable Reflection

</dd>
<dt>latency_buckets_ms</dt>
<dd>

([]float64, `gte=0`, default: `[10.0,25.0,100.0,250.0,1000.0]`) Buckets specification in latency histogram

</dd>
</dl>

### HTTPClientConfig {#http-client-config}

HTTPClientConfig holds configuration for HTTP Client.

#### Properties

<dl>
<dt>disable_compression</dt>
<dd>

(bool, default: `false`) Disable Compression

</dd>
<dt>disable_keep_alives</dt>
<dd>

(bool, default: `false`) Disable HTTP Keep Alives

</dd>
<dt>expect_continue_timeout</dt>
<dd>

(string, `gte=0s`, default: `1s`) Expect Continue Timeout. 0 = no timeout.

</dd>
<dt>idle_connection_timeout</dt>
<dd>

(string, `gte=0s`, default: `90s`) Idle Connection Timeout. 0 = no timeout.

</dd>
<dt>key_log_file</dt>
<dd>

(string) SSL key log file (useful for debugging with wireshark)

</dd>
<dt>max_conns_per_host</dt>
<dd>

(int64, `gte=0`, default: `0`) Max Connections Per Host. 0 = no limit.

</dd>
<dt>max_idle_connections</dt>
<dd>

(int64, `gte=0`, default: `100`) Max Idle Connections. 0 = no limit.

</dd>
<dt>max_idle_connections_per_host</dt>
<dd>

(int64, `gte=0`, default: `5`) Max Idle Connections per host. 0 = no limit.

</dd>
<dt>max_response_header_bytes</dt>
<dd>

(int64, `gte=0`, default: `0`) Max Response Header Bytes. 0 = no limit.

</dd>
<dt>network_keep_alive</dt>
<dd>

(string, `gte=0s`, default: `30s`) Network level keep-alive duration

</dd>
<dt>network_timeout</dt>
<dd>

(string, `gte=0s`, default: `30s`) Timeout for making network connection

</dd>
<dt>read_buffer_size</dt>
<dd>

(int64, `gte=0`, default: `0`) Read Buffer Size. 0 = 4KB

</dd>
<dt>response_header_timeout</dt>
<dd>

(string, `gte=0s`, default: `0s`) Response Header Timeout. 0 = no timeout.

</dd>
<dt>tls_handshake_timeout</dt>
<dd>

(string, `gte=0s`, default: `10s`) TLS Handshake Timeout. 0 = no timeout

</dd>
<dt>timeout</dt>
<dd>

(string, `gte=0s`, default: `60s`) HTTP client timeout - Timeouts includes connection time, redirects, reading the response etc. 0 = no timeout.

</dd>
<dt>use_proxy</dt>
<dd>

(bool, default: `false`) Use Proxy

</dd>
<dt>write_buffer_size</dt>
<dd>

(int64, `gte=0`, default: `0`) Write Buffer Size. 0 = 4KB.

</dd>
<dt>proxy_connect_header</dt>
<dd>

([Header](#header), `omitempty`)

</dd>
<dt>tls</dt>
<dd>

([ClientTLSConfig](#client-tls-config))

</dd>
</dl>

### HTTPServerConfig {#http-server-config}

HTTPServerConfig holds configuration for HTTP Server.

#### Properties

<dl>
<dt>disable_http_keep_alives</dt>
<dd>

(bool, default: `false`) Disable HTTP Keep Alives

</dd>
<dt>idle_timeout</dt>
<dd>

(string, `gte=0s`, default: `30s`) Idle timeout

</dd>
<dt>latency_buckets_ms</dt>
<dd>

([]float64, `gte=0`, default: `[10.0,25.0,100.0,250.0,1000.0]`) Buckets specification in latency histogram

</dd>
<dt>max_header_bytes</dt>
<dd>

(int64, `gte=0`, default: `1048576`) Max header size in bytes

</dd>
<dt>read_header_timeout</dt>
<dd>

(string, `gte=0s`, default: `10s`) Read header timeout

</dd>
<dt>read_timeout</dt>
<dd>

(string, `gte=0s`, default: `10s`) Read timeout

</dd>
<dt>write_timeout</dt>
<dd>

(string, `gte=0s`, default: `45s`) Write timeout

</dd>
</dl>

### Header {#header}

A Header represents the key-value pairs in an HTTP header.

The keys should be in canonical form, as returned by
CanonicalHeaderKey.

[Header](#header)

### HeapConfig {#heap-config}

HeapConfig holds configuration for heap Watchdog.

#### Properties

<dl>
<dt>limit</dt>
<dd>

(uint64, `gt=0`, default: `268435456`) Maximum memory (in bytes) sets limit of process usage. Default = 256MB.

</dd>
<dt>min_gogc</dt>
<dd>

(int64, `gt=0,lte=100`, default: `25`) Minimum GoGC sets the minimum garbage collection target percentage for heap driven Watchdogs. This setting helps avoid overscheduling.

</dd>
<dt>adaptive_policy</dt>
<dd>

([AdaptivePolicy](#adaptive-policy))

</dd>
<dt>watermarks_policy</dt>
<dd>

([WatermarksPolicy](#watermarks-policy))

</dd>
</dl>

### JobConfig {#job-config}

JobConfig is config for Job

#### Properties

<dl>
<dt>execution_period</dt>
<dd>

(string, default: `10s`) Time period between job executions. Zero or negative value means that the job will never execute periodically.

</dd>
<dt>execution_timeout</dt>
<dd>

(string, `gte=0s`, default: `5s`) Execution timeout

</dd>
<dt>initially_healthy</dt>
<dd>

(bool, default: `false`) Sets whether the job is initially healthy

</dd>
</dl>

### JobGroupConfig {#job-group-config}

JobGroupConfig holds configuration for JobGroup.

#### Properties

<dl>
<dt>blocking_execution</dt>
<dd>

(bool, default: `false`) When true, the scheduler will run jobs synchronously,
waiting for each execution instance of the job to return
before starting the next execution. Running with this
option effectively serializes all job execution.

</dd>
<dt>worker_limit</dt>
<dd>

(int64, default: `0`) Limits how many jobs can be running at the same time. This is
useful when running resource intensive jobs and a precise start time is
not critical. 0 = no limit. If BlockingExecution is set, then WorkerLimit
is ignored.

</dd>
</dl>

### ListenerConfig {#listener-config}

ListenerConfig holds configuration for socket listeners.

#### Properties

<dl>
<dt>addr</dt>
<dd>

(string, `hostname_port`, default: `:8080`) Address to bind to in the form of [host%zone]:port

</dd>
<dt>keep_alive</dt>
<dd>

(string, `gte=0s`, default: `180s`) Keep-alive period - 0 = enabled if supported by protocol or OS. If negative then keep-alive is disabled.

</dd>
<dt>network</dt>
<dd>

(string, `oneof=tcp tcp4 tcp6`, default: `tcp`) TCP networks - "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only)

</dd>
</dl>

### LogConfig {#log-config}

LogConfig holds configuration for a logger and log writers.

#### Properties

<dl>
<dt>level</dt>
<dd>

(string, `oneof=debug DEBUG info INFO warn WARN error ERROR fatal FATAL panic PANIC trace TRACE disabled DISABLED`, default: `info`) Log level

</dd>
<dt>non_blocking</dt>
<dd>

(bool, default: `true`) Use non-blocking log writer (can lose logs at high throughput)

</dd>
<dt>pretty_console</dt>
<dd>

(bool, default: `false`) Additional log writer: pretty console (stdout) logging (not recommended for prod environments)

</dd>
<dt>writers</dt>
<dd>

([[]LogWriterConfig](#log-writer-config), `omitempty,dive,omitempty`) Log writers

</dd>
</dl>

### LogWriterConfig {#log-writer-config}

LogWriterConfig holds configuration for a log writer.

#### Properties

<dl>
<dt>compress</dt>
<dd>

(bool, default: `false`) Compress

</dd>
<dt>file</dt>
<dd>

(string, default: `stderr`) Output file for logs. Keywords allowed - ["stderr", "default"]. "default" maps to `/var/log/fluxninja/<service>.log`

</dd>
<dt>max_age</dt>
<dd>

(int64, `gte=0`, default: `7`) Max age in days for log files

</dd>
<dt>max_backups</dt>
<dd>

(int64, `gte=0`, default: `3`) Max log file backups

</dd>
<dt>max_size</dt>
<dd>

(int64, `gte=0`, default: `50`) Log file max size in MB

</dd>
</dl>

### MetricsConfig {#metrics-config}

MetricsConfig holds configuration for service metrics.

#### Properties

<dl>
<dt>enable_go_metrics</dt>
<dd>

(bool, default: `false`) EnableGoCollector controls whether the go collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewGoCollector>

</dd>
<dt>enable_process_collector</dt>
<dd>

(bool, default: `false`) EnableProcessCollector controls whether the process collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewProcessCollector>

</dd>
<dt>pedantic</dt>
<dd>

(bool, default: `false`) Pedantic controls whether a pedantic Registerer is used as the prometheus backend. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewPedanticRegistry>

</dd>
</dl>

### PluginsConfig {#plugins-config}

PluginsConfig holds configuration for plugins.

#### Properties

<dl>
<dt>disable_plugins</dt>
<dd>

(bool, default: `false`) Disables all plugins

</dd>
<dt>disabled_plugins</dt>
<dd>

([]string, `omitempty`) Specific plugins to disable

</dd>
<dt>disabled_symbols</dt>
<dd>

([]string, `omitempty`) Specific plugin types to disable

</dd>
<dt>plugins_path</dt>
<dd>

(string, default: `default`) Path to plugins directory. "default" points to `/var/lib/aperture/<service>/plugins`.

</dd>
</dl>

### PortsConfig {#ports-config}

PortsConfig defines configuration for OTEL debug and extension ports.

#### Properties

<dl>
<dt>debug_port</dt>
<dd>

(uint32, `gte=0`, default: `8888`) Port on which otel collector exposes prometheus metrics on /metrics path.

</dd>
<dt>health_check_port</dt>
<dd>

(uint32, `gte=0`, default: `13133`) Port on which health check extension in exposed.

</dd>
<dt>pprof_port</dt>
<dd>

(uint32, `gte=0`, default: `1777`) Port on which pprof extension in exposed.

</dd>
<dt>zpages_port</dt>
<dd>

(uint32, `gte=0`, default: `55679`) Port on which zpages extension in exposed.

</dd>
</dl>

### ProfilersConfig {#profilers-config}

ProfilersConfig holds configuration for profilers.

#### Properties

<dl>
<dt>cpu_profiler</dt>
<dd>

(bool, default: `false`) Flag to enable cpu profiling on process start and save it to a file. HTTP interface will not work if this is enabled as CPU profile will always be running.

</dd>
<dt>profiles_path</dt>
<dd>

(string, default: `default`) Path to save performance profiles. "default" path is `/var/log/aperture/<service>/profiles`.

</dd>
<dt>register_http_routes</dt>
<dd>

(bool, default: `true`) Register routes. Profile types profile, symbol and cmdline will be registered at /debug/pprof/{profile,symbol,cmdline}.

</dd>
</dl>

### PrometheusConfig {#prometheus-config}

PrometheusConfig holds configuration for Prometheus Server.

#### Properties

<dl>
<dt>address</dt>
<dd>

(string, **required**, `required,hostname_port|url|fqdn`) Address of the prometheus server

</dd>
</dl>

### ProxyConfig {#proxy-config}

ProxyConfig holds proxy configuration.

This configuration has preference over environment variables HTTP_PROXY, HTTPS_PROXY or NO_PROXY. See <https://pkg.go.dev/golang.org/x/net/http/httpproxy#Config>

#### Properties

<dl>
<dt>http</dt>
<dd>

(string, `omitempty,url|hostname_port`)

</dd>
<dt>https</dt>
<dd>

(string, `omitempty,url|hostname_port`)

</dd>
<dt>no_proxy</dt>
<dd>

([]string, `omitempty,dive,ip|cidr|fqdn|hostname_port`)

</dd>
</dl>

### SentryConfig {#sentry-config}

SentryConfig holds configuration for Sentry.

#### Properties

<dl>
<dt>attach_stack_trace</dt>
<dd>

(bool, default: `true`) Configure to generate and attach stacktraces to capturing message calls

</dd>
<dt>debug</dt>
<dd>

(bool, default: `true`) Debug enables printing of Sentry SDK debug messages

</dd>
<dt>disabled</dt>
<dd>

(bool, default: `false`) Sentry crash report disabled

</dd>
<dt>dsn</dt>
<dd>

(string, default: `https://6223f112b0ac4344aa67e94d1631eb85@o574197.ingest.sentry.io/6605877`) If DSN is not set, the client is effectively disabled
You can set test project's dsn to send log events.
oss-aperture project dsn is set as default.

</dd>
<dt>environment</dt>
<dd>

(string, default: `production`) Environment

</dd>
<dt>sample_rate</dt>
<dd>

(float64, default: `1.0`) Sample rate for event submission i.e. 0.0 to 1.0

</dd>
<dt>traces_sample_rate</dt>
<dd>

(float64, default: `0.2`) Sample rate for sampling traces i.e. 0.0 to 1.0

</dd>
</dl>

### ServerTLSConfig {#server-tls-config}

ServerTLSConfig holds configuration for setting up server TLS support.

#### Properties

<dl>
<dt>allowed_cn</dt>
<dd>

(string, `omitempty,fqdn`) Allowed CN

</dd>
<dt>cert_file</dt>
<dd>

(string) Server Cert file path

</dd>
<dt>client_ca_file</dt>
<dd>

(string) Client CA file path

</dd>
<dt>enabled</dt>
<dd>

(bool, default: `false`) Enabled TLS

</dd>
<dt>key_file</dt>
<dd>

(string) Server Key file path

</dd>
</dl>

### WatchdogConfig {#watchdog-config}

WatchdogConfig holds configuration for Watchdog Policy. For each policy, either watermark or adaptive should be configured.

#### Properties

<dl>
<dt>cgroup</dt>
<dd>

([WatchdogPolicyType](#watchdog-policy-type))

</dd>
<dt>heap</dt>
<dd>

([HeapConfig](#heap-config))

</dd>
<dt>job</dt>
<dd>

([JobConfig](#job-config))

</dd>
<dt>system</dt>
<dd>

([WatchdogPolicyType](#watchdog-policy-type))

</dd>
</dl>

### WatchdogPolicyType {#watchdog-policy-type}

WatchdogPolicyType holds configuration Watchdog Policy algorithms. If both algorithms are configured then only watermark algorithm is used.

#### Properties

<dl>
<dt>adaptive_policy</dt>
<dd>

([AdaptivePolicy](#adaptive-policy))

</dd>
<dt>watermarks_policy</dt>
<dd>

([WatermarksPolicy](#watermarks-policy))

</dd>
</dl>

### WatermarksPolicy {#watermarks-policy}

WatermarksPolicy creates a Watchdog policy that schedules GC at concrete watermarks.

#### Properties

<dl>
<dt>enabled</dt>
<dd>

(bool, default: `false`) Flag to enable the policy

</dd>
<dt>watermarks</dt>
<dd>

([]float64, `omitempty,dive,gte=0,lte=1`, default: `[0.50,0.75,0.80,0.85,0.90,0.95,0.99]`) Watermarks are increasing limits on which to trigger GC. Watchdog disarms when the last watermark is surpassed. It is recommended to set an extreme watermark for the last element (e.g. 0.99).

</dd>
</dl>

<!---
Generated File Ends
-->
