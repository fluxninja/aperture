# Controller

Aperture Controller

## Content negotiation

### URI Schemes

- http

### Consumes

- application/json

### Produces

- application/json

## All endpoints

### common_configuration

| Method | URI                             | Name                                | Summary |
| ------ | ------------------------------- | ----------------------------------- | ------- |
| POST   | /aperture-controller/client     | [client](#client)                   |         |
| POST   | /aperture-controller/etcd       | [etcd](#etcd)                       |         |
| POST   | /aperture-controller/liveness   | [liveness](#liveness)               |         |
| POST   | /aperture-controller/log        | [log](#log)                         |         |
| POST   | /aperture-controller/metrics    | [metrics](#metrics)                 |         |
| POST   | /aperture-controller/policies   | [policies config](#policies-config) |         |
| POST   | /aperture-controller/profilers  | [profilers](#profilers)             |         |
| POST   | /aperture-controller/prometheus | [prometheus](#prometheus)           |         |
| POST   | /aperture-controller/readiness  | [readiness](#readiness)             |         |
| POST   | /aperture-controller/server     | [server](#server)                   |         |
| POST   | /aperture-controller/watchdog   | [watchdog](#watchdog)               |         |

### controller_configuration

| Method | URI                       | Name            | Summary |
| ------ | ------------------------- | --------------- | ------- |
| POST   | /aperture-controller/otel | [o tel](#o-tel) |         |

### extension_configuration

| Method | URI                            | Name                                          | Summary |
| ------ | ------------------------------ | --------------------------------------------- | ------- |
| POST   | /aperture-controller/fluxninja | [flux ninja extension](#flux-ninja-extension) |         |
| POST   | /aperture-controller/sentry    | [sentry extension](#sentry-extension)         |         |

## Paths

### <span id="client"></span> client (_Client_)

```
POST /aperture-controller/client
```

#### Parameters

| Name  | Source | Type                         | Go type              | Separator | Required | Default | Description |
| ----- | ------ | ---------------------------- | -------------------- | --------- | :------: | ------- | ----------- |
| proxy | `body` | [ProxyConfig](#proxy-config) | `models.ProxyConfig` |           |          |         |             |

#### All responses

| Code                       | Status | Description | Has headers | Schema                           |
| -------------------------- | ------ | ----------- | :---------: | -------------------------------- |
| [default](#client-default) |        |             |             | [schema](#client-default-schema) |

#### Responses

##### <span id="client-default"></span> Default Response

###### <span id="client-default-schema"></span> Schema

empty schema

### <span id="etcd"></span> etcd (_Etcd_)

```
POST /aperture-controller/etcd
```

#### Parameters

| Name | Source | Type                       | Go type             | Separator | Required | Default | Description |
| ---- | ------ | -------------------------- | ------------------- | --------- | :------: | ------- | ----------- |
|      | `body` | [EtcdConfig](#etcd-config) | `models.EtcdConfig` |           |          |         |             |

#### All responses

| Code                     | Status | Description | Has headers | Schema                         |
| ------------------------ | ------ | ----------- | :---------: | ------------------------------ |
| [default](#etcd-default) |        |             |             | [schema](#etcd-default-schema) |

#### Responses

##### <span id="etcd-default"></span> Default Response

###### <span id="etcd-default-schema"></span> Schema

empty schema

### <span id="flux-ninja-extension"></span> flux ninja extension (_FluxNinjaExtension_)

```
POST /aperture-controller/fluxninja
```

#### Parameters

| Name | Source | Type                                                     | Go type                           | Separator | Required | Default | Description |
| ---- | ------ | -------------------------------------------------------- | --------------------------------- | --------- | :------: | ------- | ----------- |
|      | `body` | [FluxNinjaExtensionConfig](#flux-ninja-extension-config) | `models.FluxNinjaExtensionConfig` |           |          |         |             |

#### All responses

| Code                                     | Status | Description | Has headers | Schema                                         |
| ---------------------------------------- | ------ | ----------- | :---------: | ---------------------------------------------- |
| [default](#flux-ninja-extension-default) |        |             |             | [schema](#flux-ninja-extension-default-schema) |

#### Responses

##### <span id="flux-ninja-extension-default"></span> Default Response

###### <span id="flux-ninja-extension-default-schema"></span> Schema

empty schema

### <span id="liveness"></span> liveness (_Liveness_)

```
POST /aperture-controller/liveness
```

#### Parameters

| Name      | Source | Type                                | Go type                 | Separator | Required | Default | Description |
| --------- | ------ | ----------------------------------- | ----------------------- | --------- | :------: | ------- | ----------- |
| scheduler | `body` | [JobGroupConfig](#job-group-config) | `models.JobGroupConfig` |           |          |         |             |
| service   | `body` | [JobConfig](#job-config)            | `models.JobConfig`      |           |          |         |             |

#### All responses

| Code                         | Status | Description | Has headers | Schema                             |
| ---------------------------- | ------ | ----------- | :---------: | ---------------------------------- |
| [default](#liveness-default) |        |             |             | [schema](#liveness-default-schema) |

#### Responses

##### <span id="liveness-default"></span> Default Response

###### <span id="liveness-default-schema"></span> Schema

empty schema

### <span id="log"></span> log (_Log_)

```
POST /aperture-controller/log
```

#### Parameters

| Name | Source | Type                     | Go type            | Separator | Required | Default | Description |
| ---- | ------ | ------------------------ | ------------------ | --------- | :------: | ------- | ----------- |
|      | `body` | [LogConfig](#log-config) | `models.LogConfig` |           |          |         |             |

#### All responses

| Code                    | Status | Description | Has headers | Schema                        |
| ----------------------- | ------ | ----------- | :---------: | ----------------------------- |
| [default](#log-default) |        |             |             | [schema](#log-default-schema) |

#### Responses

##### <span id="log-default"></span> Default Response

###### <span id="log-default-schema"></span> Schema

empty schema

### <span id="metrics"></span> metrics (_Metrics_)

```
POST /aperture-controller/metrics
```

#### Parameters

| Name | Source | Type                             | Go type                | Separator | Required | Default | Description |
| ---- | ------ | -------------------------------- | ---------------------- | --------- | :------: | ------- | ----------- |
|      | `body` | [MetricsConfig](#metrics-config) | `models.MetricsConfig` |           |          |         |             |

#### All responses

| Code                        | Status | Description | Has headers | Schema                            |
| --------------------------- | ------ | ----------- | :---------: | --------------------------------- |
| [default](#metrics-default) |        |             |             | [schema](#metrics-default-schema) |

#### Responses

##### <span id="metrics-default"></span> Default Response

###### <span id="metrics-default-schema"></span> Schema

empty schema

### <span id="o-tel"></span> o tel (_OTel_)

```
POST /aperture-controller/otel
```

#### Parameters

| Name | Source | Type                                             | Go type                       | Separator | Required | Default | Description |
| ---- | ------ | ------------------------------------------------ | ----------------------------- | --------- | :------: | ------- | ----------- |
|      | `body` | [ControllerOTelConfig](#controller-o-tel-config) | `models.ControllerOTelConfig` |           |          |         |             |

#### All responses

| Code                      | Status | Description | Has headers | Schema                          |
| ------------------------- | ------ | ----------- | :---------: | ------------------------------- |
| [default](#o-tel-default) |        |             |             | [schema](#o-tel-default-schema) |

#### Responses

##### <span id="o-tel-default"></span> Default Response

###### <span id="o-tel-default-schema"></span> Schema

empty schema

### <span id="policies-config"></span> policies config (_PoliciesConfig_)

```
POST /aperture-controller/policies
```

#### Parameters

| Name                  | Source | Type                                   | Go type                  | Separator | Required | Default | Description |
| --------------------- | ------ | -------------------------------------- | ------------------------ | --------- | :------: | ------- | ----------- |
| cr_watcher            | `body` | [CRWatcherConfig](#c-r-watcher-config) | `models.CRWatcherConfig` |           |          |         |             |
| promql_jobs_scheduler | `body` | [JobGroupConfig](#job-group-config)    | `models.JobGroupConfig`  |           |          |         |             |

#### All responses

| Code                                | Status | Description | Has headers | Schema                                    |
| ----------------------------------- | ------ | ----------- | :---------: | ----------------------------------------- |
| [default](#policies-config-default) |        |             |             | [schema](#policies-config-default-schema) |

#### Responses

##### <span id="policies-config-default"></span> Default Response

###### <span id="policies-config-default-schema"></span> Schema

empty schema

### <span id="profilers"></span> profilers (_Profilers_)

```
POST /aperture-controller/profilers
```

#### Parameters

| Name | Source | Type                                 | Go type                  | Separator | Required | Default | Description |
| ---- | ------ | ------------------------------------ | ------------------------ | --------- | :------: | ------- | ----------- |
|      | `body` | [ProfilersConfig](#profilers-config) | `models.ProfilersConfig` |           |          |         |             |

#### All responses

| Code                          | Status | Description | Has headers | Schema                              |
| ----------------------------- | ------ | ----------- | :---------: | ----------------------------------- |
| [default](#profilers-default) |        |             |             | [schema](#profilers-default-schema) |

#### Responses

##### <span id="profilers-default"></span> Default Response

###### <span id="profilers-default-schema"></span> Schema

empty schema

### <span id="prometheus"></span> prometheus (_Prometheus_)

```
POST /aperture-controller/prometheus
```

#### Parameters

| Name        | Source | Type                                    | Go type                   | Separator | Required | Default | Description |
| ----------- | ------ | --------------------------------------- | ------------------------- | --------- | :------: | ------- | ----------- |
|             | `body` | [PrometheusConfig](#prometheus-config)  | `models.PrometheusConfig` |           |          |         |             |
| http_client | `body` | [HTTPClientConfig](#http-client-config) | `models.HTTPClientConfig` |           |          |         |             |

#### All responses

| Code                           | Status | Description | Has headers | Schema                               |
| ------------------------------ | ------ | ----------- | :---------: | ------------------------------------ |
| [default](#prometheus-default) |        |             |             | [schema](#prometheus-default-schema) |

#### Responses

##### <span id="prometheus-default"></span> Default Response

###### <span id="prometheus-default-schema"></span> Schema

empty schema

### <span id="readiness"></span> readiness (_Readiness_)

```
POST /aperture-controller/readiness
```

#### Parameters

| Name      | Source | Type                                | Go type                 | Separator | Required | Default | Description |
| --------- | ------ | ----------------------------------- | ----------------------- | --------- | :------: | ------- | ----------- |
| scheduler | `body` | [JobGroupConfig](#job-group-config) | `models.JobGroupConfig` |           |          |         |             |
| service   | `body` | [JobConfig](#job-config)            | `models.JobConfig`      |           |          |         |             |

#### All responses

| Code                          | Status | Description | Has headers | Schema                              |
| ----------------------------- | ------ | ----------- | :---------: | ----------------------------------- |
| [default](#readiness-default) |        |             |             | [schema](#readiness-default-schema) |

#### Responses

##### <span id="readiness-default"></span> Default Response

###### <span id="readiness-default-schema"></span> Schema

empty schema

### <span id="sentry-extension"></span> sentry extension (_SentryExtension_)

```
POST /aperture-controller/sentry
```

#### Parameters

| Name | Source | Type                           | Go type               | Separator | Required | Default | Description |
| ---- | ------ | ------------------------------ | --------------------- | --------- | :------: | ------- | ----------- |
|      | `body` | [SentryConfig](#sentry-config) | `models.SentryConfig` |           |          |         |             |

#### All responses

| Code                                 | Status | Description | Has headers | Schema                                     |
| ------------------------------------ | ------ | ----------- | :---------: | ------------------------------------------ |
| [default](#sentry-extension-default) |        |             |             | [schema](#sentry-extension-default-schema) |

#### Responses

##### <span id="sentry-extension-default"></span> Default Response

###### <span id="sentry-extension-default-schema"></span> Schema

empty schema

### <span id="server"></span> server (_Server_)

```
POST /aperture-controller/server
```

#### Parameters

| Name         | Source | Type                                       | Go type                    | Separator | Required | Default | Description |
| ------------ | ------ | ------------------------------------------ | -------------------------- | --------- | :------: | ------- | ----------- |
| grpc         | `body` | [GRPCServerConfig](#g-rpc-server-config)   | `models.GRPCServerConfig`  |           |          |         |             |
| grpc_gateway | `body` | [GRPCGatewayConfig](#g-rpc-gateway-config) | `models.GRPCGatewayConfig` |           |          |         |             |
| http         | `body` | [HTTPServerConfig](#http-server-config)    | `models.HTTPServerConfig`  |           |          |         |             |
| listener     | `body` | [ListenerConfig](#listener-config)         | `models.ListenerConfig`    |           |          |         |             |
| tls          | `body` | [ServerTLSConfig](#server-tls-config)      | `models.ServerTLSConfig`   |           |          |         |             |

#### All responses

| Code                       | Status | Description | Has headers | Schema                           |
| -------------------------- | ------ | ----------- | :---------: | -------------------------------- |
| [default](#server-default) |        |             |             | [schema](#server-default-schema) |

#### Responses

##### <span id="server-default"></span> Default Response

###### <span id="server-default-schema"></span> Schema

empty schema

### <span id="watchdog"></span> watchdog (_Watchdog_)

```
POST /aperture-controller/watchdog
```

#### Parameters

| Name   | Source | Type                               | Go type                 | Separator | Required | Default | Description |
| ------ | ------ | ---------------------------------- | ----------------------- | --------- | :------: | ------- | ----------- |
| memory | `body` | [WatchdogConfig](#watchdog-config) | `models.WatchdogConfig` |           |          |         |             |

#### All responses

| Code                         | Status | Description | Has headers | Schema                             |
| ---------------------------- | ------ | ----------- | :---------: | ---------------------------------- |
| [default](#watchdog-default) |        |             |             | [schema](#watchdog-default-schema) |

#### Responses

##### <span id="watchdog-default"></span> Default Response

###### <span id="watchdog-default-schema"></span> Schema

empty schema

## Models

### <span id="adaptive-policy"></span> AdaptivePolicy

**Properties**

| Name    | Type                      | Go type   | Required | Default | Description                                           | Example |
| ------- | ------------------------- | --------- | :------: | ------- | ----------------------------------------------------- | ------- |
| Enabled | boolean                   | `bool`    |          |         | Flag to enable the policy                             |         |
| Factor  | double (formatted number) | `float64` |          |         | Factor sets user-configured limit of available memory |         |

### <span id="backoff-config"></span> BackoffConfig

**Properties**

| Name       | Type                      | Go type   | Required | Default | Description        | Example |
| ---------- | ------------------------- | --------- | :------: | ------- | ------------------ | ------- |
| BaseDelay  | string (formatted string) | `string`  |          |         | Base Delay         |         |
| Jitter     | double (formatted number) | `float64` |          |         | Jitter             |         |
| MaxDelay   | string (formatted string) | `string`  |          |         | Max Delay          |         |
| Multiplier | double (formatted number) | `float64` |          |         | Backoff multiplier |         |

### <span id="batch-alerts-config"></span> BatchAlertsConfig

**Properties**

| Name                | Type                       | Go type  | Required | Default | Description                                                                         | Example |
| ------------------- | -------------------------- | -------- | :------: | ------- | ----------------------------------------------------------------------------------- | ------- |
| SendBatchMaxSize    | uint32 (formatted integer) | `uint32` |          |         | SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split |
| into smaller units. |                            |
| SendBatchSize       | uint32 (formatted integer) | `uint32` |          |         | SendBatchSize is the size of a batch which after hit, will trigger it to be sent.   |         |
| Timeout             | string (formatted string)  | `string` |          |         | Timeout sets the time after which a batch will be sent regardless of size.          |         |

### <span id="c-r-watcher-config"></span> CRWatcherConfig

**Properties**

| Name    | Type    | Go type | Required | Default | Description                                                  | Example |
| ------- | ------- | ------- | :------: | ------- | ------------------------------------------------------------ | ------- |
| Enabled | boolean | `bool`  |          |         | Enabled indicates whether the Kubernetes watcher is enabled. |         |

### <span id="client-config"></span> ClientConfig

**Properties**

| Name | Type                                     | Go type            | Required | Default | Description | Example |
| ---- | ---------------------------------------- | ------------------ | :------: | ------- | ----------- | ------- |
| grpc | [GRPCClientConfig](#g-rpc-client-config) | `GRPCClientConfig` |          |         |             |         |
| http | [HTTPClientConfig](#http-client-config)  | `HTTPClientConfig` |          |         |             |         |

### <span id="client-tls-config"></span> ClientTLSConfig

**Properties**

| Name               | Type    | Go type  | Required | Default | Description | Example |
| ------------------ | ------- | -------- | :------: | ------- | ----------- | ------- |
| CAFile             | string  | `string` |          |         |             |         |
| CertFile           | string  | `string` |          |         |             |         |
| InsecureSkipVerify | boolean | `bool`   |          |         |             |         |
| KeyFile            | string  | `string` |          |         |             |         |
| KeyLogWriter       | string  | `string` |          |         |             |         |

### <span id="controller-o-tel-config"></span> ControllerOTelConfig

**Properties**

| Name         | Type                                      | Go type             | Required | Default | Description | Example |
| ------------ | ----------------------------------------- | ------------------- | :------: | ------- | ----------- | ------- |
| batch_alerts | [BatchAlertsConfig](#batch-alerts-config) | `BatchAlertsConfig` |          |         |             |         |
| ports        | [PortsConfig](#ports-config)              | `PortsConfig`       |          |         |             |         |

### <span id="etcd-config"></span> EtcdConfig

**Properties**

| Name      | Type                                  | Go type           | Required | Default | Description                   | Example |
| --------- | ------------------------------------- | ----------------- | :------: | ------- | ----------------------------- | ------- |
| Endpoints | []string                              | `[]string`        |          |         | List of etcd server endpoints |         |
| LeaseTTL  | string (formatted string)             | `string`          |          |         | Lease time-to-live            |         |
| Password  | string                                | `string`          |          |         |                               |         |
| Username  | string                                | `string`          |          |         | Authentication                |         |
| tls       | [ClientTLSConfig](#client-tls-config) | `ClientTLSConfig` |          |         |                               |         |

### <span id="flux-ninja-extension-config"></span> FluxNinjaExtensionConfig

**Properties**

| Name              | Type                           | Go type        | Required | Default | Description                                                                                                                   | Example |
| ----------------- | ------------------------------ | -------------- | :------: | ------- | ----------------------------------------------------------------------------------------------------------------------------- | ------- |
| AgentAPIKey       | string                         | `string`       |          |         | API Key for this agent. If this key is not set, the extension won't be enabled.                                               |         |
| Endpoint          | string                         | `string`       |          |         | Address to gRPC or HTTP(s) server listening in agent service. To use HTTP protocol, the address must start with `http(s)://`. |         |
| HeartbeatInterval | string (formatted string)      | `string`       |          |         | Interval between each heartbeat.                                                                                              |         |
| InstallationMode  | string                         | `string`       |          |         | Installation mode describes on which underlying platform the Agent or the Controller is being run.                            |         |
| client            | [ClientConfig](#client-config) | `ClientConfig` |          |         |                                                                                                                               |         |

### <span id="g-rpc-client-config"></span> GRPCClientConfig

**Properties**

| Name                 | Type                                  | Go type           | Required | Default | Description                | Example |
| -------------------- | ------------------------------------- | ----------------- | :------: | ------- | -------------------------- | ------- |
| Insecure             | boolean                               | `bool`            |          |         | Disable ClientTLS          |         |
| MinConnectionTimeout | string (formatted string)             | `string`          |          |         | Minimum connection timeout |         |
| UseProxy             | boolean                               | `bool`            |          |         | Use HTTP CONNECT Proxy     |         |
| backoff              | [BackoffConfig](#backoff-config)      | `BackoffConfig`   |          |         |                            |         |
| tls                  | [ClientTLSConfig](#client-tls-config) | `ClientTLSConfig` |          |         |                            |         |

### <span id="g-rpc-gateway-config"></span> GRPCGatewayConfig

> GRPCGatewayConfig holds configuration for gRPC to HTTP gateway

**Properties**

| Name     | Type   | Go type  | Required | Default | Description                                                                                                                                      | Example |
| -------- | ------ | -------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| GRPCAddr | string | `string` |          |         | gRPC server address to connect to - By default it points to HTTP server port because FluxNinja stack runs gRPC and HTTP servers on the same port |         |

### <span id="g-rpc-server-config"></span> GRPCServerConfig

**Properties**

| Name              | Type                        | Go type     | Required | Default | Description                                | Example |
| ----------------- | --------------------------- | ----------- | :------: | ------- | ------------------------------------------ | ------- |
| ConnectionTimeout | string (formatted string)   | `string`    |          |         | Connection timeout                         |         |
| EnableReflection  | boolean                     | `bool`      |          |         | Enable Reflection                          |         |
| LatencyBucketsMS  | []double (formatted number) | `[]float64` |          |         | Buckets specification in latency histogram |         |

### <span id="http-client-config"></span> HTTPClientConfig

**Properties**

| Name                   | Type                                  | Go type           | Required | Default | Description                                                                                                        | Example |
| ---------------------- | ------------------------------------- | ----------------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------ | ------- |
| DisableCompression     | boolean                               | `bool`            |          |         | Disable Compression                                                                                                |         |
| DisableKeepAlives      | boolean                               | `bool`            |          |         | Disable HTTP Keepalive                                                                                             |         |
| ExpectContinueTimeout  | string (formatted string)             | `string`          |          |         | Expect Continue Timeout. 0 = no timeout.                                                                           |         |
| IdleConnTimeout        | string (formatted string)             | `string`          |          |         | Idle Connection Timeout. 0 = no timeout.                                                                           |         |
| KeyLogWriter           | string                                | `string`          |          |         | SSL/TLS key log file (useful for debugging)                                                                        |         |
| MaxConnsPerHost        | int64 (formatted integer)             | `int64`           |          |         | Max Connections Per Host. 0 = no limit.                                                                            |         |
| MaxIdleConns           | int64 (formatted integer)             | `int64`           |          |         | Max Idle Connections. 0 = no limit.                                                                                |         |
| MaxIdleConnsPerHost    | int64 (formatted integer)             | `int64`           |          |         | Max Idle Connections per host. 0 = no limit.                                                                       |         |
| MaxResponseHeaderBytes | int64 (formatted integer)             | `int64`           |          |         | Max Response Header Bytes. 0 = no limit.                                                                           |         |
| NetworkKeepAlive       | string (formatted string)             | `string`          |          |         | Network level keep-alive duration                                                                                  |         |
| NetworkTimeout         | string (formatted string)             | `string`          |          |         | Timeout for making network connection                                                                              |         |
| ReadBufferSize         | int64 (formatted integer)             | `int64`           |          |         | Read Buffer Size. 0 = 4 KB                                                                                         |         |
| ResponseHeaderTimeout  | string (formatted string)             | `string`          |          |         | Response Header Timeout. 0 = no timeout.                                                                           |         |
| TLSHandshakeTimeout    | string (formatted string)             | `string`          |          |         | TLS Handshake Timeout. 0 = no timeout                                                                              |         |
| Timeout                | string (formatted string)             | `string`          |          |         | HTTP client timeout - Timeouts include connection time, redirects, reading the response and so on. 0 = no timeout. |         |
| UseProxy               | boolean                               | `bool`            |          |         | Use Proxy                                                                                                          |         |
| WriteBufferSize        | int64 (formatted integer)             | `int64`           |          |         | Write Buffer Size. 0 = 4 KB.                                                                                       |         |
| proxy_connect_header   | [Header](#header)                     | `Header`          |          |         |                                                                                                                    |         |
| tls                    | [ClientTLSConfig](#client-tls-config) | `ClientTLSConfig` |          |         |                                                                                                                    |         |

### <span id="http-server-config"></span> HTTPServerConfig

**Properties**

| Name                  | Type                        | Go type     | Required | Default | Description                                | Example |
| --------------------- | --------------------------- | ----------- | :------: | ------- | ------------------------------------------ | ------- |
| DisableHTTPKeepAlives | boolean                     | `bool`      |          |         | Disable HTTP Keepalive                     |         |
| IdleTimeout           | string (formatted string)   | `string`    |          |         | Idle timeout                               |         |
| LatencyBucketsMS      | []double (formatted number) | `[]float64` |          |         | Buckets specification in latency histogram |         |
| MaxHeaderBytes        | int64 (formatted integer)   | `int64`     |          |         | Max header size in bytes                   |         |
| ReadHeaderTimeout     | string (formatted string)   | `string`    |          |         | Read header timeout                        |         |
| ReadTimeout           | string (formatted string)   | `string`    |          |         | Read timeout                               |         |
| WriteTimeout          | string (formatted string)   | `string`    |          |         | Write timeout                              |         |

### <span id="header"></span> Header

> The keys should be in canonical form, as returned by CanonicalHeaderKey.

[Header](#header)

### <span id="heap-config"></span> HeapConfig

**Properties**

| Name              | Type                                   | Go type            | Required | Default | Description                                                                                                                             | Example |
| ----------------- | -------------------------------------- | ------------------ | :------: | ------- | --------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| Limit             | uint64 (formatted integer)             | `uint64`           |          |         | Maximum memory (in bytes) sets limit of process usage. Default = 256MB.                                                                 |         |
| MinGoGC           | int64 (formatted integer)              | `int64`            |          |         | Minimum GoGC sets the minimum garbage collection target percentage for heap driven Watchdogs. This setting helps avoid over scheduling. |         |
| adaptive_policy   | [AdaptivePolicy](#adaptive-policy)     | `AdaptivePolicy`   |          |         |                                                                                                                                         |         |
| watermarks_policy | [WatermarksPolicy](#watermarks-policy) | `WatermarksPolicy` |          |         |                                                                                                                                         |         |

### <span id="job-config"></span> JobConfig

> JobConfig is configuration for a periodic job

**Properties**

| Name             | Type                      | Go type  | Required | Default | Description                                                                                             | Example |
| ---------------- | ------------------------- | -------- | :------: | ------- | ------------------------------------------------------------------------------------------------------- | ------- |
| ExecutionPeriod  | string (formatted string) | `string` |          |         | Time between job executions. Zero or negative value means that the job will never execute periodically. |         |
| ExecutionTimeout | string (formatted string) | `string` |          |         | Execution timeout                                                                                       |         |
| InitiallyHealthy | boolean                   | `bool`   |          |         | Sets whether the job is initially healthy                                                               |         |

### <span id="job-group-config"></span> JobGroupConfig

**Properties**

| Name              | Type    | Go type | Required | Default | Description                                           | Example |
| ----------------- | ------- | ------- | :------: | ------- | ----------------------------------------------------- | ------- |
| BlockingExecution | boolean | `bool`  |          |         | When true, the scheduler will run jobs synchronously, |

waiting for each execution instance of the job to return before starting the
next execution. Running with this option effectively serializes all job
execution. | | | WorkerLimit | int64 (formatted integer)| `int64` | | | Limits
how many jobs can be running at the same time. This is useful when running
resource intensive jobs and a precise start time is not critical. 0 = no limit.
If BlockingExecution is set, then WorkerLimit is ignored. | |

### <span id="listener-config"></span> ListenerConfig

**Properties**

| Name      | Type                      | Go type  | Required | Default | Description                                                                                                             | Example |
| --------- | ------------------------- | -------- | :------: | ------- | ----------------------------------------------------------------------------------------------------------------------- | ------- |
| Addr      | string                    | `string` |          |         | Address to bind to in the form of `[host%zone]:port`                                                                    |         |
| KeepAlive | string (formatted string) | `string` |          |         | Keep-alive period - 0 = enabled if supported by protocol or operating system. If negative, then keep-alive is disabled. |         |
| Network   | string                    | `string` |          |         | TCP networks - `tcp`, `tcp4` (IPv4-only), `tcp6` (IPv6-only)                                                            |         |

### <span id="log-config"></span> LogConfig

**Properties**

| Name          | Type                                    | Go type              | Required | Default | Description                                                                                      | Example |
| ------------- | --------------------------------------- | -------------------- | :------: | ------- | ------------------------------------------------------------------------------------------------ | ------- |
| LogLevel      | string                                  | `string`             |          |         | Log level                                                                                        |         |
| NonBlocking   | boolean                                 | `bool`               |          |         | Use non-blocking log writer (can lose logs at high throughput)                                   |         |
| PrettyConsole | boolean                                 | `bool`               |          |         | Additional log writer: pretty console (`stdout`) logging (not recommended for prod environments) |         |
| Writers       | [][LogWriterConfig](#log-writer-config) | `[]*LogWriterConfig` |          |         | Log writers                                                                                      |         |

### <span id="log-writer-config"></span> LogWriterConfig

**Properties**

| Name       | Type                      | Go type  | Required | Default | Description                                                                                                          | Example |
| ---------- | ------------------------- | -------- | :------: | ------- | -------------------------------------------------------------------------------------------------------------------- | ------- |
| Compress   | boolean                   | `bool`   |          |         | Compress                                                                                                             |         |
| File       | string                    | `string` |          |         | Output file for logs. Keywords allowed - [`stderr`, `default`]. `default` maps to `/var/log/fluxninja/<service>.log` |         |
| MaxAge     | int64 (formatted integer) | `int64`  |          |         | Max age in days for log files                                                                                        |         |
| MaxBackups | int64 (formatted integer) | `int64`  |          |         | Max log file backups                                                                                                 |         |
| MaxSize    | int64 (formatted integer) | `int64`  |          |         | Log file max size in MB                                                                                              |         |

### <span id="metrics-config"></span> MetricsConfig

**Properties**

| Name                   | Type    | Go type | Required | Default | Description                                                                                                                                                                        | Example |
| ---------------------- | ------- | ------- | :------: | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| EnableGoCollector      | boolean | `bool`  |          |         | EnableGoCollector controls whether the go collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewGoCollector>                |         |
| EnableProcessCollector | boolean | `bool`  |          |         | EnableProcessCollector controls whether the process collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewProcessCollector> |         |
| Pedantic               | boolean | `bool`  |          |         | Pedantic controls whether a pedantic registry is used. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewPedanticRegistry>                                  |         |

### <span id="ports-config"></span> PortsConfig

**Properties**

| Name            | Type                       | Go type  | Required | Default | Description                                                               | Example |
| --------------- | -------------------------- | -------- | :------: | ------- | ------------------------------------------------------------------------- | ------- |
| DebugPort       | uint32 (formatted integer) | `uint32` |          |         | Port on which OTel collector exposes Prometheus metrics on /metrics path. |         |
| HealthCheckPort | uint32 (formatted integer) | `uint32` |          |         | Port on which health check extension in exposed.                          |         |
| PprofPort       | uint32 (formatted integer) | `uint32` |          |         | Port on which `pprof` extension in exposed.                               |         |
| ZpagesPort      | uint32 (formatted integer) | `uint32` |          |         | Port on which `zpages` extension in exposed.                              |         |

### <span id="profilers-config"></span> ProfilersConfig

**Properties**

| Name               | Type    | Go type  | Required | Default | Description                                                                                                                                                                 | Example |
| ------------------ | ------- | -------- | :------: | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| CPUProfile         | boolean | `bool`   |          |         | Flag to enable CPU profiling on process start and save it to a file. The Browser (HTTP) interface won't work if this is enabled, as the CPU profile will always be running. |         |
| ProfilesPath       | string  | `string` |          |         | Path to save performance profiles. "default" path is `/var/log/aperture/<service>/profiles`.                                                                                |         |
| RegisterHTTPRoutes | boolean | `bool`   |          |         | Register routes. Profile types `profile`, `symbol` and `cmdline` will be registered at `/debug/pprof/{profile,symbol,cmdline}`.                                             |         |

### <span id="prometheus-config"></span> PrometheusConfig

**Properties**

| Name    | Type   | Go type  | Required | Default | Description                      | Example |
| ------- | ------ | -------- | :------: | ------- | -------------------------------- | ------- |
| Address | string | `string` |          |         | Address of the Prometheus server |         |

### <span id="proxy-config"></span> ProxyConfig

> This configuration has preference over environment variables HTTP_PROXY,
> HTTPS_PROXY or NO_PROXY. See
> <https://pkg.go.dev/golang.org/x/net/http/httpproxy#Config>

**Properties**

| Name       | Type     | Go type    | Required | Default | Description | Example |
| ---------- | -------- | ---------- | :------: | ------- | ----------- | ------- |
| HTTPProxy  | string   | `string`   |          |         |             |         |
| HTTPSProxy | string   | `string`   |          |         |             |         |
| NoProxy    | []string | `[]string` |          |         |             |         |

### <span id="sentry-config"></span> SentryConfig

**Properties**

| Name             | Type    | Go type  | Required | Default | Description                                                              | Example |
| ---------------- | ------- | -------- | :------: | ------- | ------------------------------------------------------------------------ | ------- |
| AttachStacktrace | boolean | `bool`   |          |         | Configure to generate and attach stack traces to capturing message calls |         |
| Debug            | boolean | `bool`   |          |         | Debug enables printing of Sentry SDK debug messages                      |         |
| Disabled         | boolean | `bool`   |          |         | Sentry crash report disabled                                             |         |
| Dsn              | string  | `string` |          |         | If DSN is not set, the client is effectively disabled                    |

You can set test project's DSN to send log events. oss-aperture project DSN is
set as default. | | | Environment | string| `string` | | | Environment | | |
SampleRate | double (formatted number)| `float64` | | | Sample rate for event
submission | | | TracesSampleRate | double (formatted number)| `float64` | | |
Sample rate for sampling traces | |

### <span id="server-tls-config"></span> ServerTLSConfig

**Properties**

| Name         | Type    | Go type  | Required | Default | Description           | Example |
| ------------ | ------- | -------- | :------: | ------- | --------------------- | ------- |
| AllowedCN    | string  | `string` |          |         | Allowed CN            |         |
| CertFile     | string  | `string` |          |         | Server Cert file path |         |
| ClientCAFile | string  | `string` |          |         | Client CA file path   |         |
| Enabled      | boolean | `bool`   |          |         | Enabled TLS           |         |
| KeyFile      | string  | `string` |          |         | Server Key file path  |         |

### <span id="watchdog-config"></span> WatchdogConfig

**Properties**

| Name   | Type                                        | Go type              | Required | Default | Description | Example |
| ------ | ------------------------------------------- | -------------------- | :------: | ------- | ----------- | ------- |
| cgroup | [WatchdogPolicyType](#watchdog-policy-type) | `WatchdogPolicyType` |          |         |             |         |
| heap   | [HeapConfig](#heap-config)                  | `HeapConfig`         |          |         |             |         |
| job    | [JobConfig](#job-config)                    | `JobConfig`          |          |         |             |         |
| system | [WatchdogPolicyType](#watchdog-policy-type) | `WatchdogPolicyType` |          |         |             |         |

### <span id="watchdog-policy-type"></span> WatchdogPolicyType

**Properties**

| Name              | Type                                   | Go type            | Required | Default | Description | Example |
| ----------------- | -------------------------------------- | ------------------ | :------: | ------- | ----------- | ------- |
| adaptive_policy   | [AdaptivePolicy](#adaptive-policy)     | `AdaptivePolicy`   |          |         |             |         |
| watermarks_policy | [WatermarksPolicy](#watermarks-policy) | `WatermarksPolicy` |          |         |             |         |

### <span id="watermarks-policy"></span> WatermarksPolicy

**Properties**

| Name       | Type                        | Go type     | Required | Default | Description                                                                                                                                                                                            | Example |
| ---------- | --------------------------- | ----------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| Enabled    | boolean                     | `bool`      |          |         | Flag to enable the policy                                                                                                                                                                              |         |
| Watermarks | []double (formatted number) | `[]float64` |          |         | Watermarks are increasing limits on which to trigger GC. Watchdog disarms when the last watermark is surpassed. It's recommended to set an extreme watermark for the last element (for example, 0.99). |         |
