---
title: Common
sidebar_position: 1
sidebar_label: Common
---

This document describes the Prometheus metrics generated both by Aperture Agents
and the Aperture Controller.

## Go & Process Metrics

Go & Process metrics can be exposed by enabling the `enable_go_metrics` flag and
`enable_process_collector` in
[Agent's MetricsConfig](/reference/configuration/agent.md#metrics-config) and
[Controller's MetricsConfig](/reference/configuration/controller.md#metrics-config).
See

<!-- vale off -->

[collector.NewGoCollector](https://pkg.go.dev/github.com/prometheus/client_golang@v1.13.0/prometheus/collectors#NewGoCollector)
for more information.

<!-- vale on -->

## HTTP Server Metrics

### Metrics

<!-- vale off -->

| Name                     | Type      | Labels                                                                                | Unit            | Description                                     |
| ------------------------ | --------- | ------------------------------------------------------------------------------------- | --------------- | ----------------------------------------------- |
| http_requests_total      | Counter   | agent_group, instance, job, process_uuid, handler_name, http_method, http_status_code | count (no unit) | Total number of requests received               |
| http_errors_total        | Counter   | agent_group, instance, job, process_uuid, handler_name, http_method, http_status_code | count (no unit) | Total number of errors that occurred            |
| http_requests_latency_ms | Histogram | agent_group, instance, job, process_uuid, handler_name, http_method, http_status_code | ms              | Latency of the requests processed by the server |

<!-- vale on -->

### Labels

<!-- vale off -->

| Name             | Example                              | Description                                          |
| ---------------- | ------------------------------------ | ---------------------------------------------------- |
| agent_group      | default                              | Agent Group of the policy that Flux Meter belongs to |
| instance         | aperture-agent-cbfnp                 | Host instance of the Aperture Agent                  |
| job              | aperture-self                        | The configured job name that the target belongs to   |
| process_uuid     | dc0e82af-6730-4f70-8228-ee91da53ac5f | Host instance's UUID                                 |
| handler_name     | default                              | HTTP handler name                                    |
| http_method      | GET, POST                            | HTTP method                                          |
| http_status_code | 200, 503                             | HTTP status code                                     |

<!-- vale on -->

### gRPC Server Metrics

gRPC server metrics are exposed by default. See
[grpc-ecosystem/go-grpc-prometheus server_metrics](https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-prometheus#NewServerMetrics)
for more information.
