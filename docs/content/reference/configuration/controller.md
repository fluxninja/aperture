---
title: Aperture Controller Configuration Reference
sidebar_position: 10
sidebar_label: Controller
---

<head>
  <body className="schema-docs" />
</head>

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

### EXTENSION CONFIGURATION

| Key         | Reference                                   |
| ----------- | ------------------------------------------- |
| `fluxninja` | [FluxNinjaExtension](#flux-ninja-extension) |
| `sentry`    | [SentryExtension](#sentry-extension)        |

## Reference

### _client_ {#client}

<dl>

<dt>proxy</dt>
<dd>

([ProxyConfig](#proxy-config))
Env-Var Prefix: `APERTURE_CONTROLLER_CLIENT_PROXY_`

</dd>

</dl>

---

### _etcd_ {#etcd}

<dl>

<dt></dt>
<dd>

([EtcdConfig](#etcd-config))
Env-Var Prefix: `APERTURE_CONTROLLER_ETCD_`

</dd>

</dl>

---

### _fluxninja_ {#flux-ninja-extension}

<dl>

<dt></dt>
<dd>

([FluxNinjaExtensionConfig](#flux-ninja-extension-config))
Env-Var Prefix: `APERTURE_CONTROLLER_FLUXNINJA_`

</dd>

</dl>

---

### _liveness_ {#liveness}

<dl>

<dt>scheduler</dt>
<dd>

([JobGroupConfig](#job-group-config))
Env-Var Prefix: `APERTURE_CONTROLLER_LIVENESS_SCHEDULER_`

</dd>

<dt>service</dt>
<dd>

([JobConfig](#job-config))
Env-Var Prefix: `APERTURE_CONTROLLER_LIVENESS_SERVICE_`

</dd>

</dl>

---

### _log_ {#log}

<dl>

<dt></dt>
<dd>

([LogConfig](#log-config))
Env-Var Prefix: `APERTURE_CONTROLLER_LOG_`

</dd>

</dl>

---

### _metrics_ {#metrics}

<dl>

<dt></dt>
<dd>

([MetricsConfig](#metrics-config))
Env-Var Prefix: `APERTURE_CONTROLLER_METRICS_`

</dd>

</dl>

---

### _otel_ {#o-t-e-l}

<dl>

<dt></dt>
<dd>

([ControllerOTELConfig](#controller-o-t-e-l-config))
Env-Var Prefix: `APERTURE_CONTROLLER_OTEL_`

</dd>

</dl>

---

### _policies_ {#policies-config}

<dl>

<dt>promql_jobs_scheduler</dt>
<dd>

([JobGroupConfig](#job-group-config))
Env-Var Prefix: `APERTURE_CONTROLLER_POLICIES_PROMQL_JOBS_SCHEDULER_`

</dd>

</dl>

---

### _profilers_ {#profilers}

<dl>

<dt></dt>
<dd>

([ProfilersConfig](#profilers-config))
Env-Var Prefix: `APERTURE_CONTROLLER_PROFILERS_`

</dd>

</dl>

---

### _prometheus_ {#prometheus}

<dl>

<dt></dt>
<dd>

([PrometheusConfig](#prometheus-config))
Env-Var Prefix: `APERTURE_CONTROLLER_PROMETHEUS_`

</dd>

<dt>http_client</dt>
<dd>

([HTTPClientConfig](#http-client-config))
Env-Var Prefix: `APERTURE_CONTROLLER_PROMETHEUS_HTTP_CLIENT_`

</dd>

</dl>

---

### _readiness_ {#readiness}

<dl>

<dt>scheduler</dt>
<dd>

([JobGroupConfig](#job-group-config))
Env-Var Prefix: `APERTURE_CONTROLLER_READINESS_SCHEDULER_`

</dd>

<dt>service</dt>
<dd>

([JobConfig](#job-config))
Env-Var Prefix: `APERTURE_CONTROLLER_READINESS_SERVICE_`

</dd>

</dl>

---

### _sentry_ {#sentry-extension}

<dl>

<dt></dt>
<dd>

([SentryConfig](#sentry-config))
Env-Var Prefix: `APERTURE_CONTROLLER_SENTRY_`

</dd>

</dl>

---

### _server_ {#server}

<dl>

<dt>grpc</dt>
<dd>

([GRPCServerConfig](#g-rpc-server-config))
Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_GRPC_`

</dd>

<dt>grpc_gateway</dt>
<dd>

([GRPCGatewayConfig](#g-rpc-gateway-config))
Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_GRPC_GATEWAY_`

</dd>

<dt>http</dt>
<dd>

([HTTPServerConfig](#http-server-config))
Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_HTTP_`

</dd>

<dt>listener</dt>
<dd>

([ListenerConfig](#listener-config))
Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_LISTENER_`

</dd>

<dt>tls</dt>
<dd>

([ServerTLSConfig](#server-tls-config))
Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_TLS_`

</dd>

</dl>

---

### _watchdog_ {#watchdog}

<dl>

<dt>memory</dt>
<dd>

([WatchdogConfig](#watchdog-config))
Env-Var Prefix: `APERTURE_CONTROLLER_WATCHDOG_MEMORY_`

</dd>

</dl>

---

## Objects

---

### AdaptivePolicy {#adaptive-policy}

AdaptivePolicy creates a policy that forces GC when the usage surpasses the configured factor of the available memory. This policy calculates next target as usage+(limit-usage)\*factor.

<dl>
<dt>enabled</dt>
<dd>

(bool) Flag to enable the policy

</dd>
<dt>factor</dt>
<dd>

(float64, minimum: `0`, maximum: `1`, default: `0.5`) Factor sets user-configured limit of available memory

</dd>
</dl>

---

### BackoffConfig {#backoff-config}

BackoffConfig holds configuration for GRPC Client Backoff.

<dl>
<dt>base_delay</dt>
<dd>

(string, default: `"1s"`) Base Delay

</dd>
<dt>jitter</dt>
<dd>

(float64, minimum: `0`, default: `0.2`) Jitter

</dd>
<dt>max_delay</dt>
<dd>

(string, default: `"120s"`) Max Delay

</dd>
<dt>multiplier</dt>
<dd>

(float64, minimum: `0`, default: `1.6`) Backoff multiplier

</dd>
</dl>

---

### BatchAlertsConfig {#batch-alerts-config}

BatchAlertsConfig defines configuration for OTEL batch processor.

<dl>
<dt>send_batch_max_size</dt>
<dd>

(uint32, minimum: `0`) SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
into smaller units.

</dd>
<dt>send_batch_size</dt>
<dd>

(uint32, minimum: `0`) SendBatchSize is the size of a batch which after hit, will trigger it to be sent.

</dd>
<dt>timeout</dt>
<dd>

(string, default: `"1s"`) Timeout sets the time after which a batch will be sent regardless of size.

</dd>
</dl>

---

### ClientConfig {#client-config}

ClientConfig is the client configuration.

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

---

### ClientTLSConfig {#client-tls-config}

ClientTLSConfig is the config for client TLS.

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

---

### ControllerOTELConfig {#controller-o-t-e-l-config}

ControllerOTELConfig is the configuration for Agent's OTEL collector.

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

---

### EtcdConfig {#etcd-config}

EtcdConfig holds configuration for etcd client.

<dl>
<dt>endpoints</dt>
<dd>

([]string, **required**) List of Etcd server endpoints

</dd>
<dt>lease_ttl</dt>
<dd>

(string, default: `"60s"`) Lease time-to-live

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

---

### FluxNinjaExtensionConfig {#flux-ninja-extension-config}

FluxNinjaExtensionConfig is the configuration for FluxNinja ARC integration.

<dl>
<dt>api_key</dt>
<dd>

(string) API Key for this agent. If this key is not set, the extension will not be enabled.

</dd>
<dt>endpoint</dt>
<dd>

(string, format: `empty | hostname_port | url | fqdn`) Address to grpc or http(s) server listening in agent service. To use http protocol, the address must start with http(s)://.

</dd>
<dt>heartbeat_interval</dt>
<dd>

(string, default: `"5s"`) Interval between each heartbeat.

</dd>
<dt>client</dt>
<dd>

([ClientConfig](#client-config))

</dd>
</dl>

---

### GRPCClientConfig {#g-rpc-client-config}

GRPCClientConfig holds configuration for GRPC Client.

<dl>
<dt>insecure</dt>
<dd>

(bool) Disable ClientTLS

</dd>
<dt>min_connection_timeout</dt>
<dd>

(string, default: `"20s"`) Minimum connection timeout

</dd>
<dt>use_proxy</dt>
<dd>

(bool) Use HTTP CONNECT Proxy

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

---

### GRPCGatewayConfig {#g-rpc-gateway-config}

GRPCGatewayConfig holds configuration for grpc-http gateway

<dl>
<dt>grpc_server_address</dt>
<dd>

(string, format: `hostname_port`, default: `"0.0.0.0:1"`) GRPC server address to connect to - By default it points to HTTP server port because FluxNinja stack runs GRPC and HTTP servers on the same port

</dd>
</dl>

---

### GRPCServerConfig {#g-rpc-server-config}

GRPCServerConfig holds configuration for GRPC Server.

<dl>
<dt>connection_timeout</dt>
<dd>

(string, default: `"120s"`) Connection timeout

</dd>
<dt>enable_reflection</dt>
<dd>

(bool) Enable Reflection

</dd>
<dt>latency_buckets_ms</dt>
<dd>

([]float64, default: `[10,25,100,250,1000]`) Buckets specification in latency histogram

</dd>
</dl>

---

### HTTPClientConfig {#http-client-config}

HTTPClientConfig holds configuration for HTTP Client.

<dl>
<dt>disable_compression</dt>
<dd>

(bool) Disable Compression

</dd>
<dt>disable_keep_alives</dt>
<dd>

(bool) Disable HTTP Keep Alives

</dd>
<dt>expect_continue_timeout</dt>
<dd>

(string, default: `"1s"`) Expect Continue Timeout. 0 = no timeout.

</dd>
<dt>idle_connection_timeout</dt>
<dd>

(string, default: `"90s"`) Idle Connection Timeout. 0 = no timeout.

</dd>
<dt>key_log_file</dt>
<dd>

(string) SSL key log file (useful for debugging with wireshark)

</dd>
<dt>max_conns_per_host</dt>
<dd>

(int64, minimum: `0`) Max Connections Per Host. 0 = no limit.

</dd>
<dt>max_idle_connections</dt>
<dd>

(int64, minimum: `0`, default: `100`) Max Idle Connections. 0 = no limit.

</dd>
<dt>max_idle_connections_per_host</dt>
<dd>

(int64, minimum: `0`, default: `5`) Max Idle Connections per host. 0 = no limit.

</dd>
<dt>max_response_header_bytes</dt>
<dd>

(int64, minimum: `0`) Max Response Header Bytes. 0 = no limit.

</dd>
<dt>network_keep_alive</dt>
<dd>

(string, default: `"30s"`) Network level keep-alive duration

</dd>
<dt>network_timeout</dt>
<dd>

(string, default: `"30s"`) Timeout for making network connection

</dd>
<dt>read_buffer_size</dt>
<dd>

(int64, minimum: `0`) Read Buffer Size. 0 = 4KB

</dd>
<dt>response_header_timeout</dt>
<dd>

(string, default: `"0s"`) Response Header Timeout. 0 = no timeout.

</dd>
<dt>tls_handshake_timeout</dt>
<dd>

(string, default: `"10s"`) TLS Handshake Timeout. 0 = no timeout

</dd>
<dt>timeout</dt>
<dd>

(string, default: `"60s"`) HTTP client timeout - Timeouts includes connection time, redirects, reading the response etc. 0 = no timeout.

</dd>
<dt>use_proxy</dt>
<dd>

(bool) Use Proxy

</dd>
<dt>write_buffer_size</dt>
<dd>

(int64, minimum: `0`) Write Buffer Size. 0 = 4KB.

</dd>
<dt>proxy_connect_header</dt>
<dd>

([Header](#header))

</dd>
<dt>tls</dt>
<dd>

([ClientTLSConfig](#client-tls-config))

</dd>
</dl>

---

### HTTPServerConfig {#http-server-config}

HTTPServerConfig holds configuration for HTTP Server.

<dl>
<dt>disable_http_keep_alives</dt>
<dd>

(bool) Disable HTTP Keep Alives

</dd>
<dt>idle_timeout</dt>
<dd>

(string, default: `"30s"`) Idle timeout

</dd>
<dt>latency_buckets_ms</dt>
<dd>

([]float64, default: `[10,25,100,250,1000]`) Buckets specification in latency histogram

</dd>
<dt>max_header_bytes</dt>
<dd>

(int64, minimum: `0`, default: `1048576`) Max header size in bytes

</dd>
<dt>read_header_timeout</dt>
<dd>

(string, default: `"10s"`) Read header timeout

</dd>
<dt>read_timeout</dt>
<dd>

(string, default: `"10s"`) Read timeout

</dd>
<dt>write_timeout</dt>
<dd>

(string, default: `"45s"`) Write timeout

</dd>
</dl>

---

### Header {#header}

A Header represents the key-value pairs in an HTTP header.

The keys should be in canonical form, as returned by
CanonicalHeaderKey.

[Header](#header)

---

### HeapConfig {#heap-config}

HeapConfig holds configuration for heap Watchdog.

<dl>
<dt>limit</dt>
<dd>

(uint64, minimum: `0`) Maximum memory (in bytes) sets limit of process usage. Default = 256MB.

</dd>
<dt>min_gogc</dt>
<dd>

(int64, minimum: `0`, maximum: `100`, default: `25`) Minimum GoGC sets the minimum garbage collection target percentage for heap driven Watchdogs. This setting helps avoid overscheduling.

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

---

### JobConfig {#job-config}

JobConfig is config for Job

<dl>
<dt>execution_period</dt>
<dd>

(string, default: `"10s"`) Time period between job executions. Zero or negative value means that the job will never execute periodically.

</dd>
<dt>execution_timeout</dt>
<dd>

(string, default: `"5s"`) Execution timeout

</dd>
<dt>initially_healthy</dt>
<dd>

(bool) Sets whether the job is initially healthy

</dd>
</dl>

---

### JobGroupConfig {#job-group-config}

JobGroupConfig holds configuration for JobGroup.

<dl>
<dt>blocking_execution</dt>
<dd>

(bool) When true, the scheduler will run jobs synchronously,
waiting for each execution instance of the job to return
before starting the next execution. Running with this
option effectively serializes all job execution.

</dd>
<dt>worker_limit</dt>
<dd>

(int64) Limits how many jobs can be running at the same time. This is
useful when running resource intensive jobs and a precise start time is
not critical. 0 = no limit. If BlockingExecution is set, then WorkerLimit
is ignored.

</dd>
</dl>

---

### ListenerConfig {#listener-config}

ListenerConfig holds configuration for socket listeners.

<dl>
<dt>addr</dt>
<dd>

(string, format: `hostname_port`, default: `":8080"`) Address to bind to in the form of [host%zone]:port

</dd>
<dt>keep_alive</dt>
<dd>

(string, default: `"180s"`) Keep-alive period - 0 = enabled if supported by protocol or OS. If negative then keep-alive is disabled.

</dd>
<dt>network</dt>
<dd>

(string, oneof: `tcp | tcp4 | tcp6`, default: `"tcp"`) TCP networks - "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only)

</dd>
</dl>

---

### LogConfig {#log-config}

LogConfig holds configuration for a logger and log writers.

<dl>
<dt>level</dt>
<dd>

(string, oneof: `debug | DEBUG | info | INFO | warn | WARN | error | ERROR | fatal | FATAL | panic | PANIC | trace | TRACE | disabled | DISABLED`, default: `"info"`) Log level

</dd>
<dt>non_blocking</dt>
<dd>

(bool, default: `true`) Use non-blocking log writer (can lose logs at high throughput)

</dd>
<dt>pretty_console</dt>
<dd>

(bool) Additional log writer: pretty console (stdout) logging (not recommended for prod environments)

</dd>
<dt>writers</dt>
<dd>

([[]LogWriterConfig](#log-writer-config)) Log writers

</dd>
</dl>

---

### LogWriterConfig {#log-writer-config}

LogWriterConfig holds configuration for a log writer.

<dl>
<dt>compress</dt>
<dd>

(bool) Compress

</dd>
<dt>file</dt>
<dd>

(string, default: `"stderr"`) Output file for logs. Keywords allowed - ["stderr", "default"]. "default" maps to `/var/log/fluxninja/<service>.log`

</dd>
<dt>max_age</dt>
<dd>

(int64, minimum: `0`, default: `7`) Max age in days for log files

</dd>
<dt>max_backups</dt>
<dd>

(int64, minimum: `0`, default: `3`) Max log file backups

</dd>
<dt>max_size</dt>
<dd>

(int64, minimum: `0`, default: `50`) Log file max size in MB

</dd>
</dl>

---

### MetricsConfig {#metrics-config}

MetricsConfig holds configuration for service metrics.

<dl>
<dt>enable_go_metrics</dt>
<dd>

(bool) EnableGoCollector controls whether the go collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewGoCollector>

</dd>
<dt>enable_process_collector</dt>
<dd>

(bool) EnableProcessCollector controls whether the process collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewProcessCollector>

</dd>
<dt>pedantic</dt>
<dd>

(bool) Pedantic controls whether a pedantic Registerer is used as the prometheus backend. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewPedanticRegistry>

</dd>
</dl>

---

### PortsConfig {#ports-config}

PortsConfig defines configuration for OTEL debug and extension ports.

<dl>
<dt>debug_port</dt>
<dd>

(uint32, minimum: `0`) Port on which otel collector exposes prometheus metrics on /metrics path.

</dd>
<dt>health_check_port</dt>
<dd>

(uint32, minimum: `0`) Port on which health check extension in exposed.

</dd>
<dt>pprof_port</dt>
<dd>

(uint32, minimum: `0`) Port on which pprof extension in exposed.

</dd>
<dt>zpages_port</dt>
<dd>

(uint32, minimum: `0`) Port on which zpages extension in exposed.

</dd>
</dl>

---

### ProfilersConfig {#profilers-config}

ProfilersConfig holds configuration for profilers.

<dl>
<dt>cpu_profiler</dt>
<dd>

(bool) Flag to enable cpu profiling on process start and save it to a file. HTTP interface will not work if this is enabled as CPU profile will always be running.

</dd>
<dt>profiles_path</dt>
<dd>

(string, default: `"default"`) Path to save performance profiles. "default" path is `/var/log/aperture/<service>/profiles`.

</dd>
<dt>register_http_routes</dt>
<dd>

(bool, default: `true`) Register routes. Profile types profile, symbol and cmdline will be registered at /debug/pprof/{profile,symbol,cmdline}.

</dd>
</dl>

---

### PrometheusConfig {#prometheus-config}

PrometheusConfig holds configuration for Prometheus Server.

<dl>
<dt>address</dt>
<dd>

(string, format: `hostname_port | url | fqdn`, **required**) Address of the prometheus server

</dd>
</dl>

---

### ProxyConfig {#proxy-config}

ProxyConfig holds proxy configuration.

This configuration has preference over environment variables HTTP_PROXY, HTTPS_PROXY or NO_PROXY. See <https://pkg.go.dev/golang.org/x/net/http/httpproxy#Config>

<dl>
<dt>http</dt>
<dd>

(string, format: `empty | url | hostname_port`)

</dd>
<dt>https</dt>
<dd>

(string, format: `empty | url | hostname_port`)

</dd>
<dt>no_proxy</dt>
<dd>

([]string)

</dd>
</dl>

---

### SentryConfig {#sentry-config}

SentryConfig holds configuration for Sentry.

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

(bool) Sentry crash report disabled

</dd>
<dt>dsn</dt>
<dd>

(string, default: `"https://6223f112b0ac4344aa67e94d1631eb85@o574197.ingest.sentry.io/6605877"`) If DSN is not set, the client is effectively disabled
You can set test project's dsn to send log events.
oss-aperture project dsn is set as default.

</dd>
<dt>environment</dt>
<dd>

(string, default: `"production"`) Environment

</dd>
<dt>sample_rate</dt>
<dd>

(float64, default: `1`) Sample rate for event submission i.e. 0.0 to 1.0

</dd>
<dt>traces_sample_rate</dt>
<dd>

(float64, default: `0.2`) Sample rate for sampling traces i.e. 0.0 to 1.0

</dd>
</dl>

---

### ServerTLSConfig {#server-tls-config}

ServerTLSConfig holds configuration for setting up server TLS support.

<dl>
<dt>allowed_cn</dt>
<dd>

(string, format: `empty | fqdn`) Allowed CN

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

(bool) Enabled TLS

</dd>
<dt>key_file</dt>
<dd>

(string) Server Key file path

</dd>
</dl>

---

### WatchdogConfig {#watchdog-config}

WatchdogConfig holds configuration for Watchdog Policy. For each policy, either watermark or adaptive should be configured.

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

---

### WatchdogPolicyType {#watchdog-policy-type}

WatchdogPolicyType holds configuration Watchdog Policy algorithms. If both algorithms are configured then only watermark algorithm is used.

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

---

### WatermarksPolicy {#watermarks-policy}

WatermarksPolicy creates a Watchdog policy that schedules GC at concrete watermarks.

<dl>
<dt>enabled</dt>
<dd>

(bool) Flag to enable the policy

</dd>
<dt>watermarks</dt>
<dd>

([]float64, default: `[0.5,0.75,0.8,0.85,0.9,0.95,0.99]`) Watermarks are increasing limits on which to trigger GC. Watchdog disarms when the last watermark is surpassed. It is recommended to set an extreme watermark for the last element (e.g. 0.99).

</dd>
</dl>

---

<!---
Generated File Ends
-->
