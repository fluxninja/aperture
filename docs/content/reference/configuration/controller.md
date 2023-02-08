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

<dt></dt>
<dd>

Env-Var Prefix: `APERTURE_CONTROLLER_SERVER_`
Type: [ListenerConfig](#listener-config)

</dd>

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

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `enabled`            |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `factor`            |
| Type          | _float64_           |
| Default Value | `, default: `0.50`` |
| Description   | Lorem Ipsum         |

### BackoffConfig {#backoff-config}

BackoffConfig holds configuration for GRPC Client Backoff.

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `base_delay`      |
| Type          | _string_          |
| Default Value | `, default: `1s`` |
| Description   | Lorem Ipsum       |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `jitter`           |
| Type          | _float64_          |
| Default Value | `, default: `0.2`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `max_delay`         |
| Type          | _string_            |
| Default Value | `, default: `120s`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `multiplier`       |
| Type          | _float64_          |
| Default Value | `, default: `1.6`` |
| Description   | Lorem Ipsum        |

### BatchAlertsConfig {#batch-alerts-config}

BatchAlertsConfig defines configuration for OTEL batch processor.

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `send_batch_max_size` |
| Type          | _uint32_              |
| Default Value | `, default: `100``    |
| Description   | Lorem Ipsum           |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `send_batch_size`  |
| Type          | _uint32_           |
| Default Value | `, default: `100`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `timeout`         |
| Type          | _string_          |
| Default Value | `, default: `1s`` |
| Description   | Lorem Ipsum       |

### ClientConfig {#client-config}

ClientConfig is the client configuration.

| <!-- -->      | <!-- -->                                   |
| ------------- | ------------------------------------------ |
| Property      | `grpc`                                     |
| Type          | _[GRPCClientConfig](#g-rpc-client-config)_ |
| Default Value | ``                                         |
| Description   | Lorem Ipsum                                |

| <!-- -->      | <!-- -->                                  |
| ------------- | ----------------------------------------- |
| Property      | `http`                                    |
| Type          | _[HTTPClientConfig](#http-client-config)_ |
| Default Value | ``                                        |
| Description   | Lorem Ipsum                               |

### ClientTLSConfig {#client-tls-config}

ClientTLSConfig is the config for client TLS.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `ca_file`   |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `cert_file` |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `insecure_skip_verify` |
| Type          | _bool_                 |
| Default Value | ``                     |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `key_file`  |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->       |
| ------------- | -------------- |
| Property      | `key_log_file` |
| Type          | _string_       |
| Default Value | ``             |
| Description   | Lorem Ipsum    |

### ControllerOTELConfig {#controller-o-t-e-l-config}

ControllerOTELConfig is the configuration for Agent's OTEL collector.

| <!-- -->      | <!-- -->                                    |
| ------------- | ------------------------------------------- |
| Property      | `batch_alerts`                              |
| Type          | _[BatchAlertsConfig](#batch-alerts-config)_ |
| Default Value | ``                                          |
| Description   | Lorem Ipsum                                 |

| <!-- -->      | <!-- -->                       |
| ------------- | ------------------------------ |
| Property      | `ports`                        |
| Type          | _[PortsConfig](#ports-config)_ |
| Default Value | ``                             |
| Description   | Lorem Ipsum                    |

### EtcdConfig {#etcd-config}

EtcdConfig holds configuration for etcd client.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `endpoints` |
| Type          | _[]string_  |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `lease_ttl`        |
| Type          | _string_           |
| Default Value | `, default: `60s`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `password`  |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `username`  |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `tls`                                   |
| Type          | _[ClientTLSConfig](#client-tls-config)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

### FluxNinjaPluginConfig {#flux-ninja-plugin-config}

FluxNinjaPluginConfig is the configuration for FluxNinja ARC integration plugin.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `api_key`   |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `fluxninja_endpoint` |
| Type          | _string_             |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `heartbeat_interval` |
| Type          | _string_             |
| Default Value | `, default: `5s``    |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                         |
| ------------- | -------------------------------- |
| Property      | `client`                         |
| Type          | _[ClientConfig](#client-config)_ |
| Default Value | ``                               |
| Description   | Lorem Ipsum                      |

### GRPCClientConfig {#g-rpc-client-config}

GRPCClientConfig holds configuration for GRPC Client.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `insecure`           |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                 |
| ------------- | ------------------------ |
| Property      | `min_connection_timeout` |
| Type          | _string_                 |
| Default Value | `, default: `20s``       |
| Description   | Lorem Ipsum              |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `use_proxy`          |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                           |
| ------------- | ---------------------------------- |
| Property      | `backoff`                          |
| Type          | _[BackoffConfig](#backoff-config)_ |
| Default Value | ``                                 |
| Description   | Lorem Ipsum                        |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `tls`                                   |
| Type          | _[ClientTLSConfig](#client-tls-config)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

### GRPCGatewayConfig {#g-rpc-gateway-config}

GRPCGatewayConfig holds configuration for grpc-http gateway

| <!-- -->      | <!-- -->                 |
| ------------- | ------------------------ |
| Property      | `grpc_server_address`    |
| Type          | _string_                 |
| Default Value | `, default: `0.0.0.0:1`` |
| Description   | Lorem Ipsum              |

### GRPCServerConfig {#g-rpc-server-config}

GRPCServerConfig holds configuration for GRPC Server.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `connection_timeout` |
| Type          | _string_             |
| Default Value | `, default: `120s``  |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `enable_reflection`  |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `latency_buckets_ms`                          |
| Type          | _[]float64_                                   |
| Default Value | `, default: `[10.0,25.0,100.0,250.0,1000.0]`` |
| Description   | Lorem Ipsum                                   |

### HTTPClientConfig {#http-client-config}

HTTPClientConfig holds configuration for HTTP Client.

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `disable_compression` |
| Type          | _bool_                |
| Default Value | `, default: `false``  |
| Description   | Lorem Ipsum           |

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `disable_keep_alives` |
| Type          | _bool_                |
| Default Value | `, default: `false``  |
| Description   | Lorem Ipsum           |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `expect_continue_timeout` |
| Type          | _string_                  |
| Default Value | `, default: `1s``         |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `idle_connection_timeout` |
| Type          | _string_                  |
| Default Value | `, default: `90s``        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->       |
| ------------- | -------------- |
| Property      | `key_log_file` |
| Type          | _string_       |
| Default Value | ``             |
| Description   | Lorem Ipsum    |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `max_conns_per_host` |
| Type          | _int64_              |
| Default Value | `, default: `0``     |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `max_idle_connections` |
| Type          | _int64_                |
| Default Value | `, default: `100``     |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->                        |
| ------------- | ------------------------------- |
| Property      | `max_idle_connections_per_host` |
| Type          | _int64_                         |
| Default Value | `, default: `5``                |
| Description   | Lorem Ipsum                     |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `max_response_header_bytes` |
| Type          | _int64_                     |
| Default Value | `, default: `0``            |
| Description   | Lorem Ipsum                 |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `network_keep_alive` |
| Type          | _string_             |
| Default Value | `, default: `30s``   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `network_timeout`  |
| Type          | _string_           |
| Default Value | `, default: `30s`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `read_buffer_size` |
| Type          | _int64_            |
| Default Value | `, default: `0``   |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `response_header_timeout` |
| Type          | _string_                  |
| Default Value | `, default: `0s``         |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                |
| ------------- | ----------------------- |
| Property      | `tls_handshake_timeout` |
| Type          | _string_                |
| Default Value | `, default: `10s``      |
| Description   | Lorem Ipsum             |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `timeout`          |
| Type          | _string_           |
| Default Value | `, default: `60s`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `use_proxy`          |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `write_buffer_size` |
| Type          | _int64_             |
| Default Value | `, default: `0``    |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `proxy_connect_header` |
| Type          | _[Header](#header)_    |
| Default Value | ``                     |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `tls`                                   |
| Type          | _[ClientTLSConfig](#client-tls-config)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

### HTTPServerConfig {#http-server-config}

HTTPServerConfig holds configuration for HTTP Server.

| <!-- -->      | <!-- -->                   |
| ------------- | -------------------------- |
| Property      | `disable_http_keep_alives` |
| Type          | _bool_                     |
| Default Value | `, default: `false``       |
| Description   | Lorem Ipsum                |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `idle_timeout`     |
| Type          | _string_           |
| Default Value | `, default: `30s`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `latency_buckets_ms`                          |
| Type          | _[]float64_                                   |
| Default Value | `, default: `[10.0,25.0,100.0,250.0,1000.0]`` |
| Description   | Lorem Ipsum                                   |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `max_header_bytes`     |
| Type          | _int64_                |
| Default Value | `, default: `1048576`` |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `read_header_timeout` |
| Type          | _string_              |
| Default Value | `, default: `10s``    |
| Description   | Lorem Ipsum           |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `read_timeout`     |
| Type          | _string_           |
| Default Value | `, default: `10s`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `write_timeout`    |
| Type          | _string_           |
| Default Value | `, default: `45s`` |
| Description   | Lorem Ipsum        |

### Header {#header}

A Header represents the key-value pairs in an HTTP header.

The keys should be in canonical form, as returned by
CanonicalHeaderKey.

[Header](#header)

### HeapConfig {#heap-config}

HeapConfig holds configuration for heap Watchdog.

| <!-- -->      | <!-- -->                 |
| ------------- | ------------------------ |
| Property      | `limit`                  |
| Type          | _uint64_                 |
| Default Value | `, default: `268435456`` |
| Description   | Lorem Ipsum              |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `min_gogc`        |
| Type          | _int64_           |
| Default Value | `, default: `25`` |
| Description   | Lorem Ipsum       |

| <!-- -->      | <!-- -->                             |
| ------------- | ------------------------------------ |
| Property      | `adaptive_policy`                    |
| Type          | _[AdaptivePolicy](#adaptive-policy)_ |
| Default Value | ``                                   |
| Description   | Lorem Ipsum                          |

| <!-- -->      | <!-- -->                                 |
| ------------- | ---------------------------------------- |
| Property      | `watermarks_policy`                      |
| Type          | _[WatermarksPolicy](#watermarks-policy)_ |
| Default Value | ``                                       |
| Description   | Lorem Ipsum                              |

### JobConfig {#job-config}

JobConfig is config for Job

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `execution_period` |
| Type          | _string_           |
| Default Value | `, default: `10s`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `execution_timeout` |
| Type          | _string_            |
| Default Value | `, default: `5s``   |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `initial_delay`   |
| Type          | _string_          |
| Default Value | `, default: `0s`` |
| Description   | Lorem Ipsum       |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `initially_healthy`  |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

### JobGroupConfig {#job-group-config}

JobGroupConfig holds configuration for JobGroup.

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `max_concurrent_jobs` |
| Type          | _int64_               |
| Default Value | `, default: `0``      |
| Description   | Lorem Ipsum           |

### ListenerConfig {#listener-config}

ListenerConfig holds configuration for socket listeners.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `addr`               |
| Type          | _string_             |
| Default Value | `, default: `:8080`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `keep_alive`        |
| Type          | _string_            |
| Default Value | `, default: `180s`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `network`          |
| Type          | _string_           |
| Default Value | `, default: `tcp`` |
| Description   | Lorem Ipsum        |

### LogConfig {#log-config}

LogConfig holds configuration for a logger and log writers.

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `level`             |
| Type          | _string_            |
| Default Value | `, default: `info`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `non_blocking`      |
| Type          | _bool_              |
| Default Value | `, default: `true`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `pretty_console`     |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                  |
| ------------- | ----------------------------------------- |
| Property      | `writers`                                 |
| Type          | _[[]LogWriterConfig](#log-writer-config)_ |
| Default Value | ``                                        |
| Description   | Lorem Ipsum                               |

### LogWriterConfig {#log-writer-config}

LogWriterConfig holds configuration for a log writer.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `compress`           |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `file`                |
| Type          | _string_              |
| Default Value | `, default: `stderr`` |
| Description   | Lorem Ipsum           |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `max_age`        |
| Type          | _int64_          |
| Default Value | `, default: `7`` |
| Description   | Lorem Ipsum      |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `max_backups`    |
| Type          | _int64_          |
| Default Value | `, default: `3`` |
| Description   | Lorem Ipsum      |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `max_size`        |
| Type          | _int64_           |
| Default Value | `, default: `50`` |
| Description   | Lorem Ipsum       |

### MetricsConfig {#metrics-config}

MetricsConfig holds configuration for service metrics.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `enable_go_metrics`  |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                   |
| ------------- | -------------------------- |
| Property      | `enable_process_collector` |
| Type          | _bool_                     |
| Default Value | `, default: `false``       |
| Description   | Lorem Ipsum                |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `pedantic`           |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

### PluginsConfig {#plugins-config}

PluginsConfig holds configuration for plugins.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `disable_plugins`    |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `disabled_plugins` |
| Type          | _[]string_         |
| Default Value | ``                 |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `disabled_symbols` |
| Type          | _[]string_         |
| Default Value | ``                 |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `plugins_path`         |
| Type          | _string_               |
| Default Value | `, default: `default`` |
| Description   | Lorem Ipsum            |

### PortsConfig {#ports-config}

PortsConfig defines configuration for OTEL debug and extension ports.

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `debug_port`        |
| Type          | _uint32_            |
| Default Value | `, default: `8888`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `health_check_port`  |
| Type          | _uint32_             |
| Default Value | `, default: `13133`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `pprof_port`        |
| Type          | _uint32_            |
| Default Value | `, default: `1777`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `zpages_port`        |
| Type          | _uint32_             |
| Default Value | `, default: `55679`` |
| Description   | Lorem Ipsum          |

### ProfilersConfig {#profilers-config}

ProfilersConfig holds configuration for profilers.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `cpu_profiler`       |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `profiles_path`        |
| Type          | _string_               |
| Default Value | `, default: `default`` |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `register_http_routes` |
| Type          | _bool_                 |
| Default Value | `, default: `true``    |
| Description   | Lorem Ipsum            |

### PrometheusConfig {#prometheus-config}

PrometheusConfig holds configuration for Prometheus Server.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `address`   |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### ProxyConfig {#proxy-config}

ProxyConfig holds proxy configuration.

This configuration has preference over environment variables HTTP_PROXY, HTTPS_PROXY or NO_PROXY. See <https://pkg.go.dev/golang.org/x/net/http/httpproxy#Config>

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `http`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `https`     |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `no_proxy`  |
| Type          | _[]string_  |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### SentryConfig {#sentry-config}

SentryConfig holds configuration for Sentry.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `attach_stack_trace` |
| Type          | _bool_               |
| Default Value | `, default: `true``  |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `debug`             |
| Type          | _bool_              |
| Default Value | `, default: `true`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `disabled`           |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                                                                 |
| ------------- | ---------------------------------------------------------------------------------------- |
| Property      | `dsn`                                                                                    |
| Type          | _string_                                                                                 |
| Default Value | `, default: `https://6223f112b0ac4344aa67e94d1631eb85@o574197.ingest.sentry.io/6605877`` |
| Description   | Lorem Ipsum                                                                              |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `environment`             |
| Type          | _string_                  |
| Default Value | `, default: `production`` |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `sample_rate`      |
| Type          | _float64_          |
| Default Value | `, default: `1.0`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `traces_sample_rate` |
| Type          | _float64_            |
| Default Value | `, default: `0.2``   |
| Description   | Lorem Ipsum          |

### ServerTLSConfig {#server-tls-config}

ServerTLSConfig holds configuration for setting up server TLS support.

| <!-- -->      | <!-- -->     |
| ------------- | ------------ |
| Property      | `allowed_cn` |
| Type          | _string_     |
| Default Value | ``           |
| Description   | Lorem Ipsum  |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `cert_file` |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `client_ca_file` |
| Type          | _string_         |
| Default Value | ``               |
| Description   | Lorem Ipsum      |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `enabled`            |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `key_file`  |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### WatchdogConfig {#watchdog-config}

WatchdogConfig holds configuration for Watchdog Policy. For each policy, either watermark or adaptive should be configured.

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `cgroup`                                      |
| Type          | _[WatchdogPolicyType](#watchdog-policy-type)_ |
| Default Value | ``                                            |
| Description   | Lorem Ipsum                                   |

| <!-- -->      | <!-- -->                     |
| ------------- | ---------------------------- |
| Property      | `heap`                       |
| Type          | _[HeapConfig](#heap-config)_ |
| Default Value | ``                           |
| Description   | Lorem Ipsum                  |

| <!-- -->      | <!-- -->                   |
| ------------- | -------------------------- |
| Property      | `job`                      |
| Type          | _[JobConfig](#job-config)_ |
| Default Value | ``                         |
| Description   | Lorem Ipsum                |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `system`                                      |
| Type          | _[WatchdogPolicyType](#watchdog-policy-type)_ |
| Default Value | ``                                            |
| Description   | Lorem Ipsum                                   |

### WatchdogPolicyType {#watchdog-policy-type}

WatchdogPolicyType holds configuration Watchdog Policy algorithms. If both algorithms are configured then only watermark algorithm is used.

| <!-- -->      | <!-- -->                             |
| ------------- | ------------------------------------ |
| Property      | `adaptive_policy`                    |
| Type          | _[AdaptivePolicy](#adaptive-policy)_ |
| Default Value | ``                                   |
| Description   | Lorem Ipsum                          |

| <!-- -->      | <!-- -->                                 |
| ------------- | ---------------------------------------- |
| Property      | `watermarks_policy`                      |
| Type          | _[WatermarksPolicy](#watermarks-policy)_ |
| Default Value | ``                                       |
| Description   | Lorem Ipsum                              |

### WatermarksPolicy {#watermarks-policy}

WatermarksPolicy creates a Watchdog policy that schedules GC at concrete watermarks.

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `enabled`            |
| Type          | _bool_               |
| Default Value | `, default: `false`` |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                            |
| ------------- | --------------------------------------------------- |
| Property      | `watermarks`                                        |
| Type          | _[]float64_                                         |
| Default Value | `, default: `[0.50,0.75,0.80,0.85,0.90,0.95,0.99]`` |
| Description   | Lorem Ipsum                                         |

<!---
Generated File Ends
-->
