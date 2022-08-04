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

| Method | URI                                 | Name                                | Summary |
| ------ | ----------------------------------- | ----------------------------------- | ------- |
| POST   | /aperture-controller/agent_info     | [agent info](#agent-info)           |         |
| POST   | /aperture-controller/client         | [client](#client)                   |         |
| POST   | /aperture-controller/controller     | [controller](#controller)           |         |
| POST   | /aperture-controller/etcd           | [etcd](#etcd)                       |         |
| POST   | /aperture-controller/liveness       | [liveness](#liveness)               |         |
| POST   | /aperture-controller/log            | [log](#log)                         |         |
| POST   | /aperture-controller/metrics        | [metrics](#metrics)                 |         |
| POST   | /aperture-controller/peer_discovery | [peer discovery](#peer-discovery)   |         |
| POST   | /aperture-controller/plugins        | [plugins](#plugins)                 |         |
| POST   | /aperture-controller/policies       | [policies config](#policies-config) |         |
| POST   | /aperture-controller/profilers      | [profilers](#profilers)             |         |
| POST   | /aperture-controller/prometheus     | [prometheus](#prometheus)           |         |
| POST   | /aperture-controller/readiness      | [readiness](#readiness)             |         |
| POST   | /aperture-controller/server         | [server](#server)                   |         |
| POST   | /aperture-controller/watchdog       | [watchdog](#watchdog)               |         |

### plugin_configuration

| Method | URI                                   | Name                                    | Summary |
| ------ | ------------------------------------- | --------------------------------------- | ------- |
| POST   | /aperture-controller/fluxninja_plugin | [flux ninja plugin](#flux-ninja-plugin) |         |
| POST   | /aperture-controller/sentry_plugin    | [sentry plugin](#sentry-plugin)         |         |

## Paths

### <span id="agent-info"></span> agent info (_AgentInfo_)

```
POST /aperture-controller/agent_info
```

#### Parameters

| Name | Source | Type                                  | Go type                  | Separator | Required | Default | Description |
| ---- | ------ | ------------------------------------- | ------------------------ | --------- | :------: | ------- | ----------- |
|      | `body` | [AgentInfoConfig](#agent-info-config) | `models.AgentInfoConfig` |           |          |         |             |

#### All responses

| Code                           | Status | Description | Has headers | Schema                               |
| ------------------------------ | ------ | ----------- | :---------: | ------------------------------------ |
| [default](#agent-info-default) |        |             |             | [schema](#agent-info-default-schema) |

#### Responses

##### <span id="agent-info-default"></span> Default Response

###### <span id="agent-info-default-schema"></span> Schema

empty schema

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

### <span id="controller"></span> controller (_Controller_)

```
POST /aperture-controller/controller
```

#### Parameters

| Name             | Source  | Type   | Go type  | Separator | Required | Default | Description                               |
| ---------------- | ------- | ------ | -------- | --------- | :------: | ------- | ----------------------------------------- |
| classifiers_path | `query` | string | `string` |           |          |         | Directory containing classification rules |
| policies_path    | `query` | string | `string` |           |          |         | Directory containing policies rules       |

#### All responses

| Code                           | Status | Description | Has headers | Schema                               |
| ------------------------------ | ------ | ----------- | :---------: | ------------------------------------ |
| [default](#controller-default) |        |             |             | [schema](#controller-default-schema) |

#### Responses

##### <span id="controller-default"></span> Default Response

###### <span id="controller-default-schema"></span> Schema

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

### <span id="flux-ninja-plugin"></span> flux ninja plugin (_FluxNinjaPlugin_)

```
POST /aperture-controller/fluxninja_plugin
```

#### Parameters

| Name        | Source | Type                                               | Go type                        | Separator | Required | Default | Description |
| ----------- | ------ | -------------------------------------------------- | ------------------------------ | --------- | :------: | ------- | ----------- |
|             | `body` | [FluxNinjaPluginConfig](#flux-ninja-plugin-config) | `models.FluxNinjaPluginConfig` |           |          |         |             |
| client_grpc | `body` | [GRPCClientConfig](#g-rpc-client-config)           | `models.GRPCClientConfig`      |           |          |         |             |
| client_http | `body` | [HTTPClientConfig](#http-client-config)            | `models.HTTPClientConfig`      |           |          |         |             |

#### All responses

| Code                                  | Status | Description | Has headers | Schema                                      |
| ------------------------------------- | ------ | ----------- | :---------: | ------------------------------------------- |
| [default](#flux-ninja-plugin-default) |        |             |             | [schema](#flux-ninja-plugin-default-schema) |

#### Responses

##### <span id="flux-ninja-plugin-default"></span> Default Response

###### <span id="flux-ninja-plugin-default-schema"></span> Schema

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

### <span id="peer-discovery"></span> peer discovery (_PeerDiscovery_)

```
POST /aperture-controller/peer_discovery
```

#### Parameters

| Name | Source | Type                                          | Go type                      | Separator | Required | Default | Description |
| ---- | ------ | --------------------------------------------- | ---------------------------- | --------- | :------: | ------- | ----------- |
|      | `body` | [PeerDiscoveryConfig](#peer-discovery-config) | `models.PeerDiscoveryConfig` |           |          |         |             |

#### All responses

| Code                               | Status | Description | Has headers | Schema                                   |
| ---------------------------------- | ------ | ----------- | :---------: | ---------------------------------------- |
| [default](#peer-discovery-default) |        |             |             | [schema](#peer-discovery-default-schema) |

#### Responses

##### <span id="peer-discovery-default"></span> Default Response

###### <span id="peer-discovery-default-schema"></span> Schema

empty schema

### <span id="plugins"></span> plugins (_Plugins_)

```
POST /aperture-controller/plugins
```

#### Parameters

| Name | Source | Type                             | Go type                | Separator | Required | Default | Description |
| ---- | ------ | -------------------------------- | ---------------------- | --------- | :------: | ------- | ----------- |
|      | `body` | [PluginsConfig](#plugins-config) | `models.PluginsConfig` |           |          |         |             |

#### All responses

| Code                        | Status | Description | Has headers | Schema                            |
| --------------------------- | ------ | ----------- | :---------: | --------------------------------- |
| [default](#plugins-default) |        |             |             | [schema](#plugins-default-schema) |

#### Responses

##### <span id="plugins-default"></span> Default Response

###### <span id="plugins-default-schema"></span> Schema

empty schema

### <span id="policies-config"></span> policies config (_PoliciesConfig_)

```
POST /aperture-controller/policies
```

#### Parameters

| Name                  | Source | Type                                | Go type                 | Separator | Required | Default | Description |
| --------------------- | ------ | ----------------------------------- | ----------------------- | --------- | :------: | ------- | ----------- |
| promql_jobs_scheduler | `body` | [JobGroupConfig](#job-group-config) | `models.JobGroupConfig` |           |          |         |             |

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

### <span id="sentry-plugin"></span> sentry plugin (_SentryPlugin_)

```
POST /aperture-controller/sentry_plugin
```

#### Parameters

| Name   | Source | Type                           | Go type               | Separator | Required | Default | Description |
| ------ | ------ | ------------------------------ | --------------------- | --------- | :------: | ------- | ----------- |
| sentry | `body` | [SentryConfig](#sentry-config) | `models.SentryConfig` |           |          |         |             |

#### All responses

| Code                              | Status | Description | Has headers | Schema                                  |
| --------------------------------- | ------ | ----------- | :---------: | --------------------------------------- |
| [default](#sentry-plugin-default) |        |             |             | [schema](#sentry-plugin-default-schema) |

#### Responses

##### <span id="sentry-plugin-default"></span> Default Response

###### <span id="sentry-plugin-default-schema"></span> Schema

empty schema

### <span id="server"></span> server (_Server_)

```
POST /aperture-controller/server
```

#### Parameters

| Name         | Source | Type                                       | Go type                    | Separator | Required | Default | Description |
| ------------ | ------ | ------------------------------------------ | -------------------------- | --------- | :------: | ------- | ----------- |
|              | `body` | [ListenerConfig](#listener-config)         | `models.ListenerConfig`    |           |          |         |             |
| grpc         | `body` | [GRPCServerConfig](#g-rpc-server-config)   | `models.GRPCServerConfig`  |           |          |         |             |
| grpc_gateway | `body` | [GRPCGatewayConfig](#g-rpc-gateway-config) | `models.GRPCGatewayConfig` |           |          |         |             |
| http         | `body` | [HTTPServerConfig](#http-server-config)    | `models.HTTPServerConfig`  |           |          |         |             |
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

| Name      | Source | Type                                | Go type                 | Separator | Required | Default | Description |
| --------- | ------ | ----------------------------------- | ----------------------- | --------- | :------: | ------- | ----------- |
| memory    | `body` | [WatchdogConfig](#watchdog-config)  | `models.WatchdogConfig` |           |          |         |             |
| scheduler | `body` | [JobGroupConfig](#job-group-config) | `models.JobGroupConfig` |           |          |         |             |

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

### <span id="agent-info-config"></span> AgentInfoConfig

**Properties**

| Name       | Type   | Go type  | Required | Default | Description                                                                                                             | Example |
| ---------- | ------ | -------- | :------: | ------- | ----------------------------------------------------------------------------------------------------------------------- | ------- |
| AgentGroup | string | `string` |          |         | All agents within an agent_group receive the same data-plane configuration (e.g. schedulers, FluxMeters, rate limiter). |         |

### <span id="backoff-config"></span> BackoffConfig

**Properties**

| Name       | Type                      | Go type   | Required | Default | Description        | Example |
| ---------- | ------------------------- | --------- | :------: | ------- | ------------------ | ------- |
| BaseDelay  | string (formatted string) | `string`  |          |         | Base Delay         |         |
| Jitter     | double (formatted number) | `float64` |          |         | Jitter             |         |
| MaxDelay   | string (formatted string) | `string`  |          |         | Max Delay          |         |
| Multiplier | double (formatted number) | `float64` |          |         | Backoff multiplier |         |

### <span id="client-tls-config"></span> ClientTLSConfig

**Properties**

| Name               | Type    | Go type  | Required | Default | Description | Example |
| ------------------ | ------- | -------- | :------: | ------- | ----------- | ------- |
| CAFile             | string  | `string` |          |         |             |         |
| CertFile           | string  | `string` |          |         |             |         |
| InsecureSkipVerify | boolean | `bool`   |          |         |             |         |
| KeyFile            | string  | `string` |          |         |             |         |
| KeyLogWriter       | string  | `string` |          |         |             |         |

### <span id="etcd-config"></span> EtcdConfig

**Properties**

| Name      | Type                      | Go type    | Required | Default | Description                   | Example |
| --------- | ------------------------- | ---------- | :------: | ------- | ----------------------------- | ------- |
| Endpoints | []string                  | `[]string` |          |         | List of Etcd server endpoints |         |
| LeaseTTL  | string (formatted string) | `string`   |          |         | Lease time-to-live            |         |

### <span id="flux-ninja-plugin-config"></span> FluxNinjaPluginConfig

**Properties**

| Name              | Type                      | Go type  | Required | Default | Description                                                                                                                 | Example |
| ----------------- | ------------------------- | -------- | :------: | ------- | --------------------------------------------------------------------------------------------------------------------------- | ------- |
| APIKey            | string                    | `string` |          |         | API Key for this agent.                                                                                                     |         |
| FluxNinjaEndpoint | string                    | `string` |          |         | Address to grpc or http(s) server listening in agent service. To use http protocol, the address must start with http(s)://. |         |
| HeartbeatInterval | string (formatted string) | `string` |          |         | Interval between each heartbeat.                                                                                            |         |

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

> GRPCGatewayConfig holds configuration for grpc-http gateway

**Properties**

| Name     | Type   | Go type  | Required | Default | Description                                                                                                                                      | Example |
| -------- | ------ | -------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| GRPCAddr | string | `string` |          |         | GRPC server address to connect to - By default it points to HTTP server port because FluxNinja stack runs GRPC and HTTP servers on the same port |         |

### <span id="g-rpc-server-config"></span> GRPCServerConfig

**Properties**

| Name              | Type                      | Go type  | Required | Default | Description        | Example |
| ----------------- | ------------------------- | -------- | :------: | ------- | ------------------ | ------- |
| ConnectionTimeout | string (formatted string) | `string` |          |         | Connection timeout |         |
| EnableReflection  | boolean                   | `bool`   |          |         | Enable Reflection  |         |

### <span id="http-client-config"></span> HTTPClientConfig

**Properties**

| Name                   | Type                                  | Go type           | Required | Default | Description                                                                                                   | Example |
| ---------------------- | ------------------------------------- | ----------------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------- | ------- |
| DisableCompression     | boolean                               | `bool`            |          |         | Disable Compression                                                                                           |         |
| DisableKeepAlives      | boolean                               | `bool`            |          |         | Disable HTTP Keep Alives                                                                                      |         |
| ExpectContinueTimeout  | string (formatted string)             | `string`          |          |         | Expect Continue Timeout. 0 = no timeout.                                                                      |         |
| IdleConnTimeout        | string (formatted string)             | `string`          |          |         | Idle Connection Timeout. 0 = no timeout.                                                                      |         |
| KeyLogWriter           | string                                | `string`          |          |         | SSL key log file (useful for debugging with wireshark)                                                        |         |
| MaxConnsPerHost        | int64 (formatted integer)             | `int64`           |          |         | Max Connections Per Host. 0 = no limit.                                                                       |         |
| MaxIdleConns           | int64 (formatted integer)             | `int64`           |          |         | Max Idle Connections. 0 = no limit.                                                                           |         |
| MaxIdleConnsPerHost    | int64 (formatted integer)             | `int64`           |          |         | Max Idle Connections per host. 0 = no limit.                                                                  |         |
| MaxResponseHeaderBytes | int64 (formatted integer)             | `int64`           |          |         | Max Response Header Bytes. 0 = no limit.                                                                      |         |
| NetworkKeepAlive       | string (formatted string)             | `string`          |          |         | Network level keep-alive duration                                                                             |         |
| NetworkTimeout         | string (formatted string)             | `string`          |          |         | Timeout for making network connection                                                                         |         |
| ReadBufferSize         | int64 (formatted integer)             | `int64`           |          |         | Read Buffer Size. 0 = 4KB                                                                                     |         |
| ResponseHeaderTimeout  | string (formatted string)             | `string`          |          |         | Response Header Timeout. 0 = no timeout.                                                                      |         |
| TLSHandshakeTimeout    | string (formatted string)             | `string`          |          |         | TLS Handshake Timeout. 0 = no timeout                                                                         |         |
| Timeout                | string (formatted string)             | `string`          |          |         | HTTP client timeout - Timeouts includes connection time, redirects, reading the response etc. 0 = no timeout. |         |
| UseProxy               | boolean                               | `bool`            |          |         | Use Proxy                                                                                                     |         |
| WriteBufferSize        | int64 (formatted integer)             | `int64`           |          |         | Write Buffer Size. 0 = 4KB.                                                                                   |         |
| proxy_connect_header   | [Header](#header)                     | `Header`          |          |         |                                                                                                               |         |
| tls                    | [ClientTLSConfig](#client-tls-config) | `ClientTLSConfig` |          |         |                                                                                                               |         |

### <span id="http-server-config"></span> HTTPServerConfig

**Properties**

| Name                  | Type                      | Go type   | Required | Default | Description                                | Example |
| --------------------- | ------------------------- | --------- | :------: | ------- | ------------------------------------------ | ------- |
| DisableHTTPKeepAlives | boolean                   | `bool`    |          |         | Disable HTTP Keep Alives                   |         |
| IdleTimeout           | string (formatted string) | `string`  |          |         | Idle timeout                               |         |
| LatencyBucketCount    | int64 (formatted integer) | `int64`   |          |         | The number of buckets in latency histogram |         |
| LatencyBucketStartMS  | double (formatted number) | `float64` |          |         | The lowest bucket in latency histogram     |         |
| LatencyBucketWidthMS  | double (formatted number) | `float64` |          |         | The bucket width in latency histogram      |         |
| MaxHeaderBytes        | int64 (formatted integer) | `int64`   |          |         | Max header size in bytes                   |         |
| ReadHeaderTimeout     | string (formatted string) | `string`  |          |         | Read header timeout                        |         |
| ReadTimeout           | string (formatted string) | `string`  |          |         | Read timeout                               |         |
| WriteTimeout          | string (formatted string) | `string`  |          |         | Write timeout                              |         |

### <span id="header"></span> Header

> The keys should be in canonical form, as returned by
> CanonicalHeaderKey.

[Header](#header)

### <span id="heap-config"></span> HeapConfig

**Properties**

| Name              | Type                                   | Go type            | Required | Default | Description                                                                                                                            | Example |
| ----------------- | -------------------------------------- | ------------------ | :------: | ------- | -------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| Limit             | uint64 (formatted integer)             | `uint64`           |          |         | Maximum memory (in bytes) sets limit of process usage. Default = 256MB.                                                                |         |
| MinGoGC           | int64 (formatted integer)              | `int64`            |          |         | Minimum GoGC sets the minimum garbage collection target percentage for heap driven Watchdogs. This setting helps avoid overscheduling. |         |
| adaptive_policy   | [AdaptivePolicy](#adaptive-policy)     | `AdaptivePolicy`   |          |         |                                                                                                                                        |         |
| watermarks_policy | [WatermarksPolicy](#watermarks-policy) | `WatermarksPolicy` |          |         |                                                                                                                                        |         |

### <span id="job-config"></span> JobConfig

> JobConfig is config for Job

**Properties**

| Name             | Type                      | Go type  | Required | Default | Description                                                                                                                         | Example |
| ---------------- | ------------------------- | -------- | :------: | ------- | ----------------------------------------------------------------------------------------------------------------------------------- | ------- |
| ExecutionPeriod  | string (formatted string) | `string` |          |         | Time period between job executions. Zero or negative value means that the job will never execute periodically.                      |         |
| ExecutionTimeout | string (formatted string) | `string` |          |         | Execution timeout                                                                                                                   |         |
| InitialDelay     | string (formatted string) | `string` |          |         | Initial delay to start the job. Zero value will schedule the job immediately. Negative value will wait for next scheduled interval. |         |
| InitiallyHealthy | boolean                   | `bool`   |          |         | Sets whether the job is initially healthy                                                                                           |         |

### <span id="job-group-config"></span> JobGroupConfig

**Properties**

| Name              | Type                      | Go type | Required | Default | Description                                                                                                                                                       | Example |
| ----------------- | ------------------------- | ------- | :------: | ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| MaxConcurrentJobs | int64 (formatted integer) | `int64` |          |         | Limits how many jobs can be running at the same time. This is useful when running resource intensive jobs and a precise start time is not critical. 0 = no limit. |         |

### <span id="listener-config"></span> ListenerConfig

**Properties**

| Name      | Type                      | Go type  | Required | Default | Description                                                                                                | Example |
| --------- | ------------------------- | -------- | :------: | ------- | ---------------------------------------------------------------------------------------------------------- | ------- |
| Addr      | string                    | `string` |          |         | Address to bind to in the form of [host%zone]:port                                                         |         |
| KeepAlive | string (formatted string) | `string` |          |         | Keep-alive period - 0 = enabled if supported by protocol or OS. If negative then keep-alives are disabled. |         |
| Network   | string                    | `string` |          |         | TCP networks - "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only)                                               |         |

### <span id="log-config"></span> LogConfig

**Properties**

| Name          | Type                                    | Go type              | Required | Default | Description                                                                                                                    | Example |
| ------------- | --------------------------------------- | -------------------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------------------ | ------- |
| Compress      | boolean                                 | `bool`               |          |         | Compress                                                                                                                       |         |
| File          | string                                  | `string`             |          |         | Output file for logs. Keywords allowed - ["stderr", "stderr", "default"]. "default" maps to `/var/log/fluxninja/<service>.log` |         |
| LogLevel      | string                                  | `string`             |          |         | Log level                                                                                                                      |         |
| MaxAge        | int64 (formatted integer)               | `int64`              |          |         | Max age in days for log files                                                                                                  |         |
| MaxBackups    | int64 (formatted integer)               | `int64`              |          |         | Max log file backups                                                                                                           |         |
| MaxSize       | int64 (formatted integer)               | `int64`              |          |         | Log file max size in MB                                                                                                        |         |
| NonBlocking   | boolean                                 | `bool`               |          |         | Use non-blocking log writer (can lose logs at high throughput)                                                                 |         |
| PrettyConsole | boolean                                 | `bool`               |          |         | Additional log writer: pretty console (stdout) logging (not recommended for prod environments)                                 |         |
| Writers       | [][logwriterconfig](#log-writer-config) | `[]*LogWriterConfig` |          |         | Additional log writers                                                                                                         |         |

### <span id="log-writer-config"></span> LogWriterConfig

**Properties**

| Name       | Type                      | Go type  | Required | Default | Description                                                                                                                    | Example |
| ---------- | ------------------------- | -------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------------------ | ------- |
| Compress   | boolean                   | `bool`   |          |         | Compress                                                                                                                       |         |
| File       | string                    | `string` |          |         | Output file for logs. Keywords allowed - ["stderr", "stderr", "default"]. "default" maps to `/var/log/fluxninja/<service>.log` |         |
| MaxAge     | int64 (formatted integer) | `int64`  |          |         | Max age in days for log files                                                                                                  |         |
| MaxBackups | int64 (formatted integer) | `int64`  |          |         | Max log file backups                                                                                                           |         |
| MaxSize    | int64 (formatted integer) | `int64`  |          |         | Log file max size in MB                                                                                                        |         |

### <span id="metrics-config"></span> MetricsConfig

**Properties**

| Name                   | Type    | Go type | Required | Default | Description                                                                                                                                                                        | Example |
| ---------------------- | ------- | ------- | :------: | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| EnableGoCollector      | boolean | `bool`  |          |         | EnableGoCollector controls whether the go collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewGoCollector>                |         |
| EnableProcessCollector | boolean | `bool`  |          |         | EnableProcessCollector controls whether the process collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewProcessCollector> |         |
| Pedantic               | boolean | `bool`  |          |         | Pedantic controls whether a pedantic Registerer is used as the prometheus backend. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewPedanticRegistry>      |         |

### <span id="peer-discovery-config"></span> PeerDiscoveryConfig

**Properties**

| Name              | Type   | Go type  | Required | Default | Description                                                                                                                                          | Example |
| ----------------- | ------ | -------- | :------: | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| AdvertisementAddr | string | `string` |          |         | Network address of aperture server to advertise to peers - this address should be reachable from other agents. Used for nat traversal when provided. |         |

### <span id="plugins-config"></span> PluginsConfig

**Properties**

| Name            | Type     | Go type    | Required | Default | Description                                                                    | Example |
| --------------- | -------- | ---------- | :------: | ------- | ------------------------------------------------------------------------------ | ------- |
| DisablePlugins  | boolean  | `bool`     |          |         | Disables all plugins                                                           |         |
| DisabledPlugins | []string | `[]string` |          |         | Specific plugins to disable                                                    |         |
| DisabledSymbols | []string | `[]string` |          |         | Specific plugin types to disable                                               |         |
| PluginsPath     | string   | `string`   |          |         | Path to plugins directory. This can be set via command line arguments as well. |         |

### <span id="profilers-config"></span> ProfilersConfig

**Properties**

| Name         | Type    | Go type  | Required | Default | Description                                                                            | Example |
| ------------ | ------- | -------- | :------: | ------- | -------------------------------------------------------------------------------------- | ------- |
| CPUProfile   | boolean | `bool`   |          |         | Flag to enable cpu profiling                                                           |         |
| ProfilesPath | string  | `string` |          |         | Path to save performance profiles. This can be set via command line arguments as well. |         |

### <span id="prometheus-config"></span> PrometheusConfig

**Properties**

| Name    | Type   | Go type  | Required | Default | Description                      | Example |
| ------- | ------ | -------- | :------: | ------- | -------------------------------- | ------- |
| Address | string | `string` |          |         | Address of the prometheus server |         |

### <span id="proxy-config"></span> ProxyConfig

> This configuration has preference over environment variables HTTP_PROXY, HTTPS_PROXY or NO_PROXY. See <https://pkg.go.dev/golang.org/x/net/http/httpproxy#Config>

**Properties**

| Name       | Type     | Go type    | Required | Default | Description | Example |
| ---------- | -------- | ---------- | :------: | ------- | ----------- | ------- |
| HTTPProxy  | string   | `string`   |          |         |             |         |
| HTTPSProxy | string   | `string`   |          |         |             |         |
| NoProxy    | []string | `[]string` |          |         |             |         |

### <span id="sentry-config"></span> SentryConfig

**Properties**

| Name             | Type    | Go type  | Required | Default | Description                                                             | Example |
| ---------------- | ------- | -------- | :------: | ------- | ----------------------------------------------------------------------- | ------- |
| AttachStacktrace | boolean | `bool`   |          |         | Configure to generate and attach stacktraces to capturing message calls |         |
| Debug            | boolean | `bool`   |          |         | Debug enables printing of Sentry SDK debug messages                     |         |
| Disabled         | boolean | `bool`   |          |         | Sentry crash report disabled                                            |         |
| Dsn              | string  | `string` |          |         | If DSN is not set, the client is effectively disabled                   |

You can set test project's dsn to send log events.
i.e. oss-aperture: <https://6223f112b0ac4344aa67e94d1631eb85@o574197.ingest.sentry.io/6605877> | |
| Environment | string| `string` | | | Environment | |
| SampleRate | double (formatted number)| `float64` | | | Sample rate for event submission i.e. 0.0 to 1.0 | |
| TracesSampleRate | double (formatted number)| `float64` | | | Sample rate for sampling traces i.e. 0.0 to 1.0 | |

### <span id="server-tls-config"></span> ServerTLSConfig

**Properties**

| Name       | Type    | Go type  | Required | Default | Description                                                              | Example |
| ---------- | ------- | -------- | :------: | ------- | ------------------------------------------------------------------------ | ------- |
| AllowedCN  | string  | `string` |          |         | Allowed CN                                                               |         |
| CertsPath  | string  | `string` |          |         | Path to credentials. This can be set via command line arguments as well. |         |
| ClientCA   | string  | `string` |          |         | Client CA file                                                           |         |
| Enable     | boolean | `bool`   |          |         | Enable TLS                                                               |         |
| ServerCert | string  | `string` |          |         | Server Cert file                                                         |         |
| ServerKey  | string  | `string` |          |         | Server Key file                                                          |         |

### <span id="watchdog-config"></span> WatchdogConfig

**Properties**

| Name   | Type                                        | Go type              | Required | Default | Description | Example |
| ------ | ------------------------------------------- | -------------------- | :------: | ------- | ----------- | ------- |
| cgroup | [WatchdogPolicyType](#watchdog-policy-type) | `WatchdogPolicyType` |          |         |             |         |
| heap   | [HeapConfig](#heap-config)                  | `HeapConfig`         |          |         |             |         |
| system | [WatchdogPolicyType](#watchdog-policy-type) | `WatchdogPolicyType` |          |         |             |         |

### <span id="watchdog-policy-type"></span> WatchdogPolicyType

**Properties**

| Name              | Type                                   | Go type            | Required | Default | Description | Example |
| ----------------- | -------------------------------------- | ------------------ | :------: | ------- | ----------- | ------- |
| adaptive_policy   | [AdaptivePolicy](#adaptive-policy)     | `AdaptivePolicy`   |          |         |             |         |
| watermarks_policy | [WatermarksPolicy](#watermarks-policy) | `WatermarksPolicy` |          |         |             |         |

### <span id="watermarks-policy"></span> WatermarksPolicy

**Properties**

| Name       | Type                        | Go type     | Required | Default | Description                                                                                                                                                                                     | Example |
| ---------- | --------------------------- | ----------- | :------: | ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| Enabled    | boolean                     | `bool`      |          |         | Flag to enable the policy                                                                                                                                                                       |         |
| Watermarks | []double (formatted number) | `[]float64` |          |         | Watermarks are increasing limits on which to trigger GC. Watchdog disarms when the last watermark is surpassed. It is recommended to set an extreme watermark for the last element (e.g. 0.99). |         |
