---
title: Aperture Agent Configuration Reference
sidebar_position: 2
sidebar_label: Agent
---

<!-- markdownlint-disable -->
<!-- vale off -->

<head>
  <body className="schema-docs" />
</head>

<!-- vale on -->

:::info

See also
[Aperture Agent installation](/get-started/installation/agent/agent.md).

:::

List of all configuration parameters for Aperture Agent.

<!---
Generated File Starts
-->

## Table of contents

### AGENT CONFIGURATION

<!-- vale off -->

| Key    | Reference      |
| ------ | -------------- |
| `otel` | [OTel](#o-tel) |

### COMMON CONFIGURATION

<!-- vale off -->

| Key                 | Reference                              |
| ------------------- | -------------------------------------- |
| `agent_info`        | [AgentInfo](#agent-info)               |
| `auto_scale`        | [AutoScaleConfig](#auto-scale-config)  |
| `client`            | [Client](#client)                      |
| `dist_cache`        | [DistCache](#dist-cache)               |
| `etcd`              | [Etcd](#etcd)                          |
| `flow_control`      | [FlowControl](#flow-control)           |
| `agent_functions`   | [Functions](#functions)                |
| `kubernetes_client` | [KubernetesClient](#kubernetes-client) |
| `liveness`          | [Liveness](#liveness)                  |
| `log`               | [Log](#log)                            |
| `metrics`           | [Metrics](#metrics)                    |
| `peer_discovery`    | [PeerDiscovery](#peer-discovery)       |
| `profilers`         | [Profilers](#profilers)                |
| `prometheus`        | [Prometheus](#prometheus)              |
| `readiness`         | [Readiness](#readiness)                |
| `server`            | [Server](#server)                      |
| `service_discovery` | [ServiceDiscovery](#service-discovery) |
| `watchdog`          | [Watchdog](#watchdog)                  |

### EXTENSION CONFIGURATION

<!-- vale off -->

| Key         | Reference                                   |
| ----------- | ------------------------------------------- |
| `fluxninja` | [FluxNinjaExtension](#flux-ninja-extension) |
| `sentry`    | [SentryExtension](#sentry-extension)        |

<!-- vale on -->

## Reference

<!-- vale off -->

### _agent_info_ {#agent-info}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([AgentInfoConfig](#agent-info-config)) Environment variable prefix:
`APERTURE_AGENT_AGENT_INFO_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _auto_scale_ {#auto-scale-config}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>kubernetes</dt>
<dd>

([AutoScaleKubernetesConfig](#auto-scale-kubernetes-config)) Environment
variable prefix: `APERTURE_AGENT_AUTO_SCALE_KUBERNETES_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _client_ {#client}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>proxy</dt>
<dd>

([ProxyConfig](#proxy-config)) Environment variable prefix:
`APERTURE_AGENT_CLIENT_PROXY_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _dist_cache_ {#dist-cache}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([DistCacheConfig](#dist-cache-config)) Environment variable prefix:
`APERTURE_AGENT_DIST_CACHE_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _etcd_ {#etcd}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([EtcdConfig](#etcd-config)) Environment variable prefix: `APERTURE_AGENT_ETCD_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _flow_control_ {#flow-control}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>preview_service</dt>
<dd>

([FlowPreviewConfig](#flow-preview-config)) Environment variable prefix:
`APERTURE_AGENT_FLOW_CONTROL_PREVIEW_SERVICE_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _fluxninja_ {#flux-ninja-extension}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([FluxNinjaExtensionConfig](#flux-ninja-extension-config)) Environment variable
prefix: `APERTURE_AGENT_FLUXNINJA_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _agent_functions_ {#functions}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([AgentFunctionsConfig](#agent-functions-config)) Environment variable prefix:
`APERTURE_AGENT_AGENT_FUNCTIONS_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _kubernetes_client_ {#kubernetes-client}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>http_client</dt>
<dd>

([HTTPClientConfig](#http-client-config)) Environment variable prefix:
`APERTURE_AGENT_KUBERNETES_CLIENT_HTTP_CLIENT_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _liveness_ {#liveness}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>scheduler</dt>
<dd>

([JobGroupConfig](#job-group-config)) Environment variable prefix:
`APERTURE_AGENT_LIVENESS_SCHEDULER_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>service</dt>
<dd>

([JobConfig](#job-config)) Environment variable prefix:
`APERTURE_AGENT_LIVENESS_SERVICE_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _log_ {#log}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([LogConfig](#log-config)) Environment variable prefix: `APERTURE_AGENT_LOG_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _metrics_ {#metrics}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([MetricsConfig](#metrics-config)) Environment variable prefix:
`APERTURE_AGENT_METRICS_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _otel_ {#o-tel}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([AgentOTelConfig](#agent-o-tel-config)) Environment variable prefix:
`APERTURE_AGENT_OTEL_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _peer_discovery_ {#peer-discovery}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([PeerDiscoveryConfig](#peer-discovery-config)) Environment variable prefix:
`APERTURE_AGENT_PEER_DISCOVERY_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _profilers_ {#profilers}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([ProfilersConfig](#profilers-config)) Environment variable prefix:
`APERTURE_AGENT_PROFILERS_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _prometheus_ {#prometheus}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([PrometheusConfig](#prometheus-config)) Environment variable prefix:
`APERTURE_AGENT_PROMETHEUS_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>http_client</dt>
<dd>

([HTTPClientConfig](#http-client-config)) Environment variable prefix:
`APERTURE_AGENT_PROMETHEUS_HTTP_CLIENT_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _readiness_ {#readiness}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>scheduler</dt>
<dd>

([JobGroupConfig](#job-group-config)) Environment variable prefix:
`APERTURE_AGENT_READINESS_SCHEDULER_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>service</dt>
<dd>

([JobConfig](#job-config)) Environment variable prefix:
`APERTURE_AGENT_READINESS_SERVICE_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _sentry_ {#sentry-extension}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt></dt>
<dd>

([SentryConfig](#sentry-config)) Environment variable prefix:
`APERTURE_AGENT_SENTRY_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _server_ {#server}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>grpc</dt>
<dd>

([GRPCServerConfig](#g-rpc-server-config)) Environment variable prefix:
`APERTURE_AGENT_SERVER_GRPC_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>grpc_gateway</dt>
<dd>

([GRPCGatewayConfig](#g-rpc-gateway-config)) Environment variable prefix:
`APERTURE_AGENT_SERVER_GRPC_GATEWAY_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>http</dt>
<dd>

([HTTPServerConfig](#http-server-config)) Environment variable prefix:
`APERTURE_AGENT_SERVER_HTTP_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>listener</dt>
<dd>

([ListenerConfig](#listener-config)) Environment variable prefix:
`APERTURE_AGENT_SERVER_LISTENER_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>tls</dt>
<dd>

([ServerTLSConfig](#server-tls-config)) Environment variable prefix:
`APERTURE_AGENT_SERVER_TLS_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _service_discovery_ {#service-discovery}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>kubernetes</dt>
<dd>

([KubernetesDiscoveryConfig](#kubernetes-discovery-config)) Environment variable
prefix: `APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_`

</dd>

<!-- vale off -->

<!-- vale off -->

<dt>static</dt>
<dd>

([StaticDiscoveryConfig](#static-discovery-config)) Environment variable prefix:
`APERTURE_AGENT_SERVICE_DISCOVERY_STATIC_`

</dd>

<!-- vale off -->

</dl>

---

<!-- vale off -->

### _watchdog_ {#watchdog}

<!-- vale on -->

<dl>

<!-- vale off -->

<dt>memory</dt>
<dd>

([WatchdogConfig](#watchdog-config)) Environment variable prefix:
`APERTURE_AGENT_WATCHDOG_MEMORY_`

</dd>

<!-- vale off -->

</dl>

---

## Objects

---

<!-- vale off -->

### AdaptivePolicy {#adaptive-policy}

<!-- vale on -->

AdaptivePolicy creates a policy that forces GC when the usage surpasses the
configured factor of the available memory. This policy calculates next target as
usage+(limit-usage)\*factor.

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Flag to enable the policy

</dd>
<dt>factor</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, maximum: `1`, default: `0.5`)

<!-- vale on -->

Factor sets user-configured limit of available memory

</dd>
</dl>

---

<!-- vale off -->

### AgentFunctionsConfig {#agent-functions-config}

<!-- vale on -->

AgentFunctionsConfig is configuration for agent functions.

<dl>
<dt>endpoints</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

RPC servers to connect to (which will be able to call agent functions)

</dd>
<dt>client</dt>
<dd>

<!-- vale off -->

([ClientConfig](#client-config))

<!-- vale on -->

Network client configuration

</dd>
</dl>

---

<!-- vale off -->

### AgentInfoConfig {#agent-info-config}

<!-- vale on -->

AgentInfoConfig is the configuration for the agent group and other agent
attributes.

<dl>
<dt>agent_group</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

All agents within an agent group receive the same data-plane configuration (for
example, Flux Meters, Rate Limiters and so on).

[Read more about agent groups here](/concepts/selector.md#agent-group).

</dd>
</dl>

---

<!-- vale off -->

### AgentOTelConfig {#agent-o-tel-config}

<!-- vale on -->

AgentOTelConfig is the configuration for Agent's OTel collector.

<dl>
<dt>disable_kubelet_scraper</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

DisableKubeletScraper disables the default metrics collection for kubelet.
Deprecated: kubelet scraper is removed entirely, so this flag makes no
difference.

</dd>
<dt>disable_kubernetes_scraper</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

DisableKubernetesScraper disables the default metrics collection for Kubernetes
resources.

</dd>
<dt>enable_high_cardinality_platform_metrics</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

EnableHighCardinalityPlatformMetrics filters out high cardinality Aperture
platform metrics from being published to Prometheus. Filtered out metrics are:
"grpc_server_handled_total._" "grpc_server_handling_seconds._"
"grpc_server_handling_seconds_bucket._" "grpc_server_handling_seconds_count._"
"grpc_server_handling_seconds_sum._" "grpc_server_msg_received_total._"
"grpc_server_msg_sent_total._" "grpc_server_started_total._"

</dd>
<dt>batch_alerts</dt>
<dd>

<!-- vale off -->

([BatchAlertsConfig](#batch-alerts-config))

<!-- vale on -->

BatchAlerts configures batch alerts processor.

</dd>
<dt>batch_postrollup</dt>
<dd>

<!-- vale off -->

([BatchPostrollupConfig](#batch-postrollup-config))

<!-- vale on -->

BatchPostrollup configures the OTel batch post-processor.

</dd>
<dt>batch_prerollup</dt>
<dd>

<!-- vale off -->

([BatchPrerollupConfig](#batch-prerollup-config))

<!-- vale on -->

BatchPrerollup configures the OTel batch pre-processor.

</dd>
<dt>ports</dt>
<dd>

<!-- vale off -->

([PortsConfig](#ports-config))

<!-- vale on -->

Ports configures debug, health and extension ports values.

</dd>
</dl>

---

<!-- vale off -->

### AutoScaleKubernetesConfig {#auto-scale-kubernetes-config}

<!-- vale on -->

AutoScaleKubernetesConfig is the configuration for the flow preview service.

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Enables the Kubernetes auto-scale capability.

</dd>
</dl>

---

<!-- vale off -->

### BackoffConfig {#backoff-config}

<!-- vale on -->

BackoffConfig holds configuration for gRPC client backoff.

<dl>
<dt>base_delay</dt>
<dd>

<!-- vale off -->

(string, default: `"1s"`)

<!-- vale on -->

Base Delay

</dd>
<dt>jitter</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, default: `0.2`)

<!-- vale on -->

Jitter

</dd>
<dt>max_delay</dt>
<dd>

<!-- vale off -->

(string, default: `"120s"`)

<!-- vale on -->

Max Delay

</dd>
<dt>multiplier</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, default: `1.6`)

<!-- vale on -->

Backoff multiplier

</dd>
</dl>

---

<!-- vale off -->

### BatchAlertsConfig {#batch-alerts-config}

<!-- vale on -->

BatchAlertsConfig defines configuration for OTel batch processor.

<dl>
<dt>send_batch_max_size</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be
split into smaller units.

</dd>
<dt>send_batch_size</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

SendBatchSize is the size of a batch which after hit, will trigger it to be
sent.

</dd>
<dt>timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"1s"`)

<!-- vale on -->

Timeout sets the time after which a batch will be sent regardless of size.

</dd>
</dl>

---

<!-- vale off -->

### BatchPostrollupConfig {#batch-postrollup-config}

<!-- vale on -->

BatchPostrollupConfig defines configuration for OTel batch processor.

<dl>
<dt>send_batch_max_size</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be
split into smaller units.

</dd>
<dt>send_batch_size</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

SendBatchSize is the size of a batch which after hit, will trigger it to be
sent.

</dd>
<dt>timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"1s"`)

<!-- vale on -->

Timeout sets the time after which a batch will be sent regardless of size.

</dd>
</dl>

---

<!-- vale off -->

### BatchPrerollupConfig {#batch-prerollup-config}

<!-- vale on -->

BatchPrerollupConfig defines configuration for OTel batch processor.

<dl>
<dt>send_batch_max_size</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be
split into smaller units.

</dd>
<dt>send_batch_size</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

SendBatchSize is the size of a batch which after hit, will trigger it to be
sent.

</dd>
<dt>timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

Timeout sets the time after which a batch will be sent regardless of size.

</dd>
</dl>

---

<!-- vale off -->

### ClientConfig {#client-config}

<!-- vale on -->

ClientConfig is the client configuration.

<dl>
<dt>grpc</dt>
<dd>

<!-- vale off -->

([GRPCClientConfig](#g-rpc-client-config))

<!-- vale on -->

gRPC client settings.

</dd>
<dt>http</dt>
<dd>

<!-- vale off -->

([HTTPClientConfig](#http-client-config))

<!-- vale on -->

HTTP client settings.

</dd>
</dl>

---

<!-- vale off -->

### ClientTLSConfig {#client-tls-config}

<!-- vale on -->

ClientTLSConfig is the configuration for client TLS.

<dl>
<dt>ca_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

</dd>
<dt>cert_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

</dd>
<dt>insecure_skip_verify</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

</dd>
<dt>key_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

</dd>
<dt>key_log_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### DistCacheConfig {#dist-cache-config}

<!-- vale on -->

DistCacheConfig configures distributed cache that holds per-label counters in
distributed rate limiters.

<dl>
<dt>bind_addr</dt>
<dd>

<!-- vale off -->

(string, format: `hostname_port`, default: `":3320"`)

<!-- vale on -->

BindAddr denotes the address that DistCache will bind to for communication with
other peer nodes.

</dd>
<dt>memberlist_advertise_addr</dt>
<dd>

<!-- vale off -->

(string, format: `empty | hostname_port`)

<!-- vale on -->

Address of [`memberlist`](https://github.com/hashicorp/memberlist) to advertise
to other cluster members. Used for NAT traversal if provided.

</dd>
<dt>memberlist_bind_addr</dt>
<dd>

<!-- vale off -->

(string, format: `hostname_port`, default: `":3322"`)

<!-- vale on -->

Address to bind [`memberlist`](https://github.com/hashicorp/memberlist) server
to.

</dd>
<dt>replica_count</dt>
<dd>

<!-- vale off -->

(int64, default: `1`)

<!-- vale on -->

ReplicaCount is 1 by default.

</dd>
<dt>sync_replication</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

SyncReplication enables synchronous replication. By default the replication is
asynchronous.

</dd>
</dl>

---

<!-- vale off -->

### Entity {#entity}

<!-- vale on -->

Entity represents a pod, VM, and so on.

<dl>
<dt>ip_address</dt>
<dd>

<!-- vale off -->

(string, format: `ip`, **required**)

<!-- vale on -->

IP address of the entity.

</dd>
<dt>name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Name of the entity. For example, pod name.

</dd>
<dt>namespace</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Namespace of the entity. For example, pod namespace.

</dd>
<dt>node_name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Node name of the entity. For example, hostname.

</dd>
<dt>services</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

Services of the entity.

</dd>
<dt>uid</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Unique identifier of the entity.

</dd>
</dl>

---

<!-- vale off -->

### EtcdConfig {#etcd-config}

<!-- vale on -->

EtcdConfig holds configuration for etcd client.

<dl>
<dt>endpoints</dt>
<dd>

<!-- vale off -->

([]string, **required**)

<!-- vale on -->

List of etcd server endpoints

</dd>
<dt>lease_ttl</dt>
<dd>

<!-- vale off -->

(string, default: `"60s"`)

<!-- vale on -->

Lease time-to-live

</dd>
<dt>namespace</dt>
<dd>

<!-- vale off -->

(string, default: `"aperture"`)

<!-- vale on -->

etcd namespace

</dd>
<dt>password</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

</dd>
<dt>username</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Authentication

</dd>
<dt>tls</dt>
<dd>

<!-- vale off -->

([ClientTLSConfig](#client-tls-config))

<!-- vale on -->

Client TLS configuration

</dd>
</dl>

---

<!-- vale off -->

### FlowPreviewConfig {#flow-preview-config}

<!-- vale on -->

FlowPreviewConfig is the configuration for the flow control preview service.

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Enables the flow preview service.

</dd>
</dl>

---

<!-- vale off -->

### FluxNinjaExtensionConfig {#flux-ninja-extension-config}

<!-- vale on -->

FluxNinjaExtensionConfig is the configuration for FluxNinja ARC integration.

<dl>
<dt>api_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

API Key for this agent. If this key is not set, the extension won't be enabled.

</dd>
<dt>controller_id</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Controller ID.

</dd>
<dt>disable_local_otel_pipeline</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Whether to configure local Prometheus OTel pipeline for metrics. Implied to be
true by EnableCloudController.

</dd>
<dt>enable_cloud_controller</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Whether to enable cloud controller. Overrides etcd and TLS configurations.

</dd>
<dt>endpoint</dt>
<dd>

<!-- vale off -->

(string, format: `empty | hostname_port | url | fqdn`)

<!-- vale on -->

Address to gRPC or HTTP(s) server listening in agent service. To use HTTP
protocol, the address must start with `http(s)://`.

</dd>
<dt>heartbeat_interval</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

Interval between each heartbeat.

</dd>
<dt>installation_mode</dt>
<dd>

<!-- vale off -->

(string, one of: `KUBERNETES_SIDECAR | KUBERNETES_DAEMONSET | LINUX_BARE_METAL`,
default: `"LINUX_BARE_METAL"`)

<!-- vale on -->

Installation mode describes on which underlying platform the Agent or the
Controller is being run.

</dd>
<dt>client</dt>
<dd>

<!-- vale off -->

([ClientConfig](#client-config))

<!-- vale on -->

Client configuration.

</dd>
</dl>

---

<!-- vale off -->

### GRPCClientConfig {#g-rpc-client-config}

<!-- vale on -->

GRPCClientConfig holds configuration for gRPC Client.

<dl>
<dt>insecure</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Disable ClientTLS

</dd>
<dt>min_connection_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"20s"`)

<!-- vale on -->

Minimum connection timeout

</dd>
<dt>use_proxy</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Use HTTP CONNECT Proxy

</dd>
<dt>backoff</dt>
<dd>

<!-- vale off -->

([BackoffConfig](#backoff-config))

<!-- vale on -->

Backoff configuration

</dd>
<dt>tls</dt>
<dd>

<!-- vale off -->

([ClientTLSConfig](#client-tls-config))

<!-- vale on -->

Client TLS configuration

</dd>
</dl>

---

<!-- vale off -->

### GRPCGatewayConfig {#g-rpc-gateway-config}

<!-- vale on -->

GRPCGatewayConfig holds configuration for gRPC to HTTP gateway

<dl>
<dt>grpc_server_address</dt>
<dd>

<!-- vale off -->

(string, format: `empty | hostname_port`)

<!-- vale on -->

gRPC server address to connect to - By default it points to HTTP server port
because FluxNinja stack runs gRPC and HTTP servers on the same port

</dd>
</dl>

---

<!-- vale off -->

### GRPCServerConfig {#g-rpc-server-config}

<!-- vale on -->

GRPCServerConfig holds configuration for gRPC Server.

<dl>
<dt>connection_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"120s"`)

<!-- vale on -->

Connection timeout

</dd>
<dt>enable_reflection</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Enable Reflection

</dd>
<dt>latency_buckets_ms</dt>
<dd>

<!-- vale off -->

([]float64, default: `[10,25,100,250,1000]`)

<!-- vale on -->

Buckets specification in latency histogram

</dd>
</dl>

---

<!-- vale off -->

### HTTPClientConfig {#http-client-config}

<!-- vale on -->

HTTPClientConfig holds configuration for HTTP Client.

<dl>
<dt>disable_compression</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Disable Compression

</dd>
<dt>disable_keep_alives</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Disable HTTP Keepalive

</dd>
<dt>expect_continue_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"1s"`)

<!-- vale on -->

Expect Continue Timeout. 0 = no timeout.

</dd>
<dt>idle_connection_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"90s"`)

<!-- vale on -->

Idle Connection Timeout. 0 = no timeout.

</dd>
<dt>key_log_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

SSL/TLS key log file (useful for debugging)

</dd>
<dt>max_conns_per_host</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `0`)

<!-- vale on -->

Max Connections Per Host. 0 = no limit.

</dd>
<dt>max_idle_connections</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `100`)

<!-- vale on -->

Max Idle Connections. 0 = no limit.

</dd>
<dt>max_idle_connections_per_host</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `5`)

<!-- vale on -->

Max Idle Connections per host. 0 = no limit.

</dd>
<dt>max_response_header_bytes</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `0`)

<!-- vale on -->

Max Response Header Bytes. 0 = no limit.

</dd>
<dt>network_keep_alive</dt>
<dd>

<!-- vale off -->

(string, default: `"30s"`)

<!-- vale on -->

Network level keep-alive duration

</dd>
<dt>network_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"30s"`)

<!-- vale on -->

Timeout for making network connection

</dd>
<dt>read_buffer_size</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `0`)

<!-- vale on -->

Read Buffer Size. 0 = 4 KB

</dd>
<dt>response_header_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"0s"`)

<!-- vale on -->

Response Header Timeout. 0 = no timeout.

</dd>
<dt>tls_handshake_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

TLS Handshake Timeout. 0 = no timeout

</dd>
<dt>timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"60s"`)

<!-- vale on -->

HTTP client timeout - Timeouts include connection time, redirects, reading the
response and so on. 0 = no timeout.

</dd>
<dt>use_proxy</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Use Proxy

</dd>
<dt>write_buffer_size</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `0`)

<!-- vale on -->

Write Buffer Size. 0 = 4 KB.

</dd>
<dt>proxy_connect_header</dt>
<dd>

<!-- vale off -->

([Header](#header))

<!-- vale on -->

Proxy Connect Header - `map[string][]string`

</dd>
<dt>tls</dt>
<dd>

<!-- vale off -->

([ClientTLSConfig](#client-tls-config))

<!-- vale on -->

Client TLS configuration

</dd>
</dl>

---

<!-- vale off -->

### HTTPServerConfig {#http-server-config}

<!-- vale on -->

HTTPServerConfig holds configuration for HTTP Server.

<dl>
<dt>disable_http_keep_alives</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Disable HTTP Keepalive

</dd>
<dt>idle_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"30s"`)

<!-- vale on -->

Idle timeout

</dd>
<dt>latency_buckets_ms</dt>
<dd>

<!-- vale off -->

([]float64, default: `[10,25,100,250,1000]`)

<!-- vale on -->

Buckets specification in latency histogram

</dd>
<dt>max_header_bytes</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `1048576`)

<!-- vale on -->

Max header size in bytes

</dd>
<dt>read_header_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

Read header timeout

</dd>
<dt>read_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

Read timeout

</dd>
<dt>write_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"45s"`)

<!-- vale on -->

Write timeout

</dd>
</dl>

---

<!-- vale off -->

### Header {#header}

<!-- vale on -->

A Header represents the key-value pairs in an HTTP header.

The keys should be in canonical form, as returned by CanonicalHeaderKey.

[Header](#header)

---

<!-- vale off -->

### HeapConfig {#heap-config}

<!-- vale on -->

HeapConfig holds configuration for heap Watchdog.

<dl>
<dt>limit</dt>
<dd>

<!-- vale off -->

(uint64, minimum: `0`)

<!-- vale on -->

Maximum memory (in bytes) sets limit of process usage. Default = 256MB.

</dd>
<dt>min_gogc</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, maximum: `100`, default: `25`)

<!-- vale on -->

Minimum GoGC sets the minimum garbage collection target percentage for heap
driven Watchdogs. This setting helps avoid over scheduling.

</dd>
<dt>adaptive_policy</dt>
<dd>

<!-- vale off -->

([AdaptivePolicy](#adaptive-policy))

<!-- vale on -->

</dd>
<dt>watermarks_policy</dt>
<dd>

<!-- vale off -->

([WatermarksPolicy](#watermarks-policy))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### JobConfig {#job-config}

<!-- vale on -->

JobConfig is configuration for a periodic job

<dl>
<dt>execution_period</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

Time between job executions. Zero or negative value means that the job will
never execute periodically.

</dd>
<dt>execution_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

Execution timeout

</dd>
<dt>initially_healthy</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Sets whether the job is initially healthy

</dd>
</dl>

---

<!-- vale off -->

### JobGroupConfig {#job-group-config}

<!-- vale on -->

JobGroupConfig holds configuration for JobGroup.

<dl>
<dt>blocking_execution</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

When true, the scheduler will run jobs synchronously, waiting for each execution
instance of the job to return before starting the next execution. Running with
this option effectively serializes all job execution.

</dd>
<dt>worker_limit</dt>
<dd>

<!-- vale off -->

(int64, default: `0`)

<!-- vale on -->

Limits how many jobs can be running at the same time. This is useful when
running resource intensive jobs and a precise start time is not critical. 0 = no
limit. If BlockingExecution is set, then WorkerLimit is ignored.

</dd>
</dl>

---

<!-- vale off -->

### KubernetesDiscoveryConfig {#kubernetes-discovery-config}

<!-- vale on -->

KubernetesDiscoveryConfig for Kubernetes service discovery.

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### ListenerConfig {#listener-config}

<!-- vale on -->

ListenerConfig holds configuration for socket listeners.

<dl>
<dt>addr</dt>
<dd>

<!-- vale off -->

(string, format: `hostname_port`, default: `":8080"`)

<!-- vale on -->

Address to bind to in the form of `[host%zone]:port`

</dd>
<dt>keep_alive</dt>
<dd>

<!-- vale off -->

(string, default: `"180s"`)

<!-- vale on -->

Keep-alive period - 0 = enabled if supported by protocol or operating system. If
negative, then keep-alive is disabled.

</dd>
<dt>network</dt>
<dd>

<!-- vale off -->

(string, one of: `tcp | tcp4 | tcp6`, default: `"tcp"`)

<!-- vale on -->

TCP networks - `tcp`, `tcp4` (IPv4-only), `tcp6` (IPv6-only)

</dd>
</dl>

---

<!-- vale off -->

### LogConfig {#log-config}

<!-- vale on -->

LogConfig holds configuration for a logger and log writers.

<dl>
<dt>level</dt>
<dd>

<!-- vale off -->

(string, one of:
`debug | DEBUG | info | INFO | warn | WARN | error | ERROR | fatal | FATAL | panic | PANIC | trace | TRACE | disabled | DISABLED`,
default: `"info"`)

<!-- vale on -->

Log level

</dd>
<dt>non_blocking</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Use non-blocking log writer (can lose logs at high throughput)

</dd>
<dt>pretty_console</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Additional log writer: pretty console (`stdout`) logging (not recommended for
prod environments)

</dd>
<dt>writers</dt>
<dd>

<!-- vale off -->

([[]LogWriterConfig](#log-writer-config))

<!-- vale on -->

Log writers

</dd>
</dl>

---

<!-- vale off -->

### LogWriterConfig {#log-writer-config}

<!-- vale on -->

LogWriterConfig holds configuration for a log writer.

<dl>
<dt>compress</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Compress

</dd>
<dt>file</dt>
<dd>

<!-- vale off -->

(string, default: `"stderr"`)

<!-- vale on -->

Output file for logs. Keywords allowed - [`stderr`, `default`]. `default` maps
to `/var/log/fluxninja/<service>.log`

</dd>
<dt>max_age</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `7`)

<!-- vale on -->

Max age in days for log files

</dd>
<dt>max_backups</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `3`)

<!-- vale on -->

Max log file backups

</dd>
<dt>max_size</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `50`)

<!-- vale on -->

Log file max size in MB

</dd>
</dl>

---

<!-- vale off -->

### MetricsConfig {#metrics-config}

<!-- vale on -->

MetricsConfig holds configuration for service metrics.

<dl>
<dt>enable_go_metrics</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

EnableGoCollector controls whether the go collector is registered on startup.
See
<https://godoc.org/github.com/prometheus/client_golang/prometheus#NewGoCollector>

</dd>
<dt>enable_process_collector</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

EnableProcessCollector controls whether the process collector is registered on
startup. See
<https://godoc.org/github.com/prometheus/client_golang/prometheus#NewProcessCollector>

</dd>
<dt>pedantic</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Pedantic controls whether a pedantic registry is used. See
<https://godoc.org/github.com/prometheus/client_golang/prometheus#NewPedanticRegistry>

</dd>
</dl>

---

<!-- vale off -->

### PeerDiscoveryConfig {#peer-discovery-config}

<!-- vale on -->

PeerDiscoveryConfig holds configuration for Agent Peer Discovery.

<dl>
<dt>advertisement_addr</dt>
<dd>

<!-- vale off -->

(string, format: `empty | hostname_port`)

<!-- vale on -->

Network address of aperture server to advertise to peers - this address should
be reachable from other agents. Used for NAT traversal when provided.

</dd>
</dl>

---

<!-- vale off -->

### PortsConfig {#ports-config}

<!-- vale on -->

PortsConfig defines configuration for OTel debug and extension ports.

<dl>
<dt>debug_port</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

Port on which OTel collector exposes Prometheus metrics on /metrics path.

</dd>
<dt>health_check_port</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

Port on which health check extension in exposed.

</dd>
<dt>pprof_port</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

Port on which `pprof` extension in exposed.

</dd>
<dt>zpages_port</dt>
<dd>

<!-- vale off -->

(uint32, minimum: `0`)

<!-- vale on -->

Port on which `zpages` extension in exposed.

</dd>
</dl>

---

<!-- vale off -->

### ProfilersConfig {#profilers-config}

<!-- vale on -->

ProfilersConfig holds configuration for profilers.

<dl>
<dt>cpu_profiler</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Flag to enable CPU profiling on process start and save it to a file. The Browser
(HTTP) interface won't work if this is enabled, as the CPU profile will always
be running.

</dd>
<dt>profiles_path</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

Path to save performance profiles. "default" path is
`/var/log/aperture/<service>/profiles`.

</dd>
<dt>register_http_routes</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Register routes. Profile types `profile`, `symbol` and `cmdline` will be
registered at `/debug/pprof/{profile,symbol,cmdline}`.

</dd>
</dl>

---

<!-- vale off -->

### PrometheusConfig {#prometheus-config}

<!-- vale on -->

PrometheusConfig holds configuration for Prometheus Server.

<dl>
<dt>address</dt>
<dd>

<!-- vale off -->

(string, format: `hostname_port | url | fqdn`, **required**)

<!-- vale on -->

Address of the Prometheus server

</dd>
<dt>labels</dt>
<dd>

<!-- vale off -->

([[]PrometheusLabel](#prometheus-label))

<!-- vale on -->

A list of labels to be attached to every query

</dd>
</dl>

---

<!-- vale off -->

### PrometheusLabel {#prometheus-label}

<!-- vale on -->

PrometheusLabel holds Name->Value mapping for the label that will be attached to
every PromQL query executed by the controller.

<dl>
<dt>name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

</dd>
<dt>value</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### ProxyConfig {#proxy-config}

<!-- vale on -->

ProxyConfig holds proxy configuration.

This configuration has preference over environment variables HTTP_PROXY,
HTTPS_PROXY or NO_PROXY. See
<https://pkg.go.dev/golang.org/x/net/http/httpproxy#Config>

<dl>
<dt>http</dt>
<dd>

<!-- vale off -->

(string, format: `empty | url | hostname_port`)

<!-- vale on -->

</dd>
<dt>https</dt>
<dd>

<!-- vale off -->

(string, format: `empty | url | hostname_port`)

<!-- vale on -->

</dd>
<dt>no_proxy</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### SentryConfig {#sentry-config}

<!-- vale on -->

SentryConfig holds configuration for Sentry.

<dl>
<dt>attach_stack_trace</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Configure to generate and attach stack traces to capturing message calls

</dd>
<dt>debug</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Debug enables printing of Sentry SDK debug messages

</dd>
<dt>disabled</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Sentry crash report disabled

</dd>
<dt>dsn</dt>
<dd>

<!-- vale off -->

(string, default:
`"https://6223f112b0ac4344aa67e94d1631eb85@o574197.ingest.sentry.io/6605877"`)

<!-- vale on -->

If DSN is not set, the client is effectively disabled You can set test project's
DSN to send log events. oss-aperture project DSN is set as default.

</dd>
<dt>environment</dt>
<dd>

<!-- vale off -->

(string, default: `"production"`)

<!-- vale on -->

Environment

</dd>
<dt>sample_rate</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, maximum: `1`, default: `1`)

<!-- vale on -->

Sample rate for event submission

</dd>
<dt>traces_sample_rate</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, maximum: `1`, default: `0.2`)

<!-- vale on -->

Sample rate for sampling traces

</dd>
</dl>

---

<!-- vale off -->

### ServerTLSConfig {#server-tls-config}

<!-- vale on -->

ServerTLSConfig holds configuration for setting up server TLS support.

<dl>
<dt>allowed_cn</dt>
<dd>

<!-- vale off -->

(string, format: `empty | fqdn`)

<!-- vale on -->

Allowed CN

</dd>
<dt>cert_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Server Cert file path

</dd>
<dt>client_ca_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Client CA file path

</dd>
<dt>enabled</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Enabled TLS

</dd>
<dt>key_file</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Server Key file path

</dd>
</dl>

---

<!-- vale off -->

### StaticDiscoveryConfig {#static-discovery-config}

<!-- vale on -->

StaticDiscoveryConfig for pre-determined list of services.

<dl>
<dt>entities</dt>
<dd>

<!-- vale off -->

([[]Entity](#entity))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### WatchdogConfig {#watchdog-config}

<!-- vale on -->

WatchdogConfig holds configuration for Watchdog Policy. For each policy, either
watermark or adaptive should be configured.

<dl>
<dt>cgroup</dt>
<dd>

<!-- vale off -->

([WatchdogPolicyType](#watchdog-policy-type))

<!-- vale on -->

</dd>
<dt>heap</dt>
<dd>

<!-- vale off -->

([HeapConfig](#heap-config))

<!-- vale on -->

</dd>
<dt>job</dt>
<dd>

<!-- vale off -->

([JobConfig](#job-config))

<!-- vale on -->

</dd>
<dt>system</dt>
<dd>

<!-- vale off -->

([WatchdogPolicyType](#watchdog-policy-type))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### WatchdogPolicyType {#watchdog-policy-type}

<!-- vale on -->

WatchdogPolicyType holds configuration Watchdog Policy algorithms. If both
algorithms are configured then only watermark algorithm is used.

<dl>
<dt>adaptive_policy</dt>
<dd>

<!-- vale off -->

([AdaptivePolicy](#adaptive-policy))

<!-- vale on -->

</dd>
<dt>watermarks_policy</dt>
<dd>

<!-- vale off -->

([WatermarksPolicy](#watermarks-policy))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### WatermarksPolicy {#watermarks-policy}

<!-- vale on -->

WatermarksPolicy creates a Watchdog policy that schedules GC at concrete
watermarks.

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Flag to enable the policy

</dd>
<dt>watermarks</dt>
<dd>

<!-- vale off -->

([]float64, default: `[0.5,0.75,0.8,0.85,0.9,0.95,0.99]`)

<!-- vale on -->

Watermarks are increasing limits on which to trigger GC. Watchdog disarms when
the last watermark is surpassed. It's recommended to set an extreme watermark
for the last element (for example, 0.99).

</dd>
</dl>

---

<!---
Generated File Ends
-->
