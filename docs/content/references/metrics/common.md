---
title: System
sidebar_position: 1
sidebar_label: System
---

## Go Metrics

Go process metrics can be exposed by enabling `enable_go_metrics` flag in
[Agent's MetricsConfig](../configuration/agent.md#metrics-config) and
[Controller's MetricsConfig](../configuration/controller.md#metrics-config). See
[collector.NewGoCollector](https://pkg.go.dev/github.com/prometheus/client_golang@v1.13.0/prometheus/collectors#NewGoCollector)
for more information.

## Process Metrics

Go process metrics can be exposed by enabling `enable_process_collector` flag in
[Agent's MetricsConfig](../configuration/agent.md#metrics-config) and
[Controller's MetricsConfig](../configuration/controller.md#metrics-config). See
[collector.NewProcessCollector](https://pkg.go.dev/github.com/prometheus/client_golang@v1.13.0/prometheus/collectors#NewProcessCollector)
for more information.
