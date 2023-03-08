---
title: Flow Label
sidebar_position: 1
---

:::info What is a flow?

A flow is the fundamental unit of work from the perspective of an Aperture
Agent. It could be an API call, a feature, or even a database query. A flow in
Aperture is similar to [OpenTelemetry span][span].

:::

Every flow is annotated with a set of **Flow Labels**. Each _Flow Label_ is a
key:value pair. If a flow is annotated with `user_tier:gold` label, then
`user_tier` is a label key and `gold` is a label value.

_Flow Labels_ are used in different ways in Aperture:

- [_Flow Selector_][flow-selector] can select flows based on _Flow Labels_, thus
  flow labels can be used to narrow the scope of [_Classifiers_][classifier],
  [_Flux Meters_][flux-meter] etc.
- _Flow Labels_ are used to map a flow to a [_Workload_][workload].
- Fairness within [_Scheduler_][scheduler] and [_Rate Limiter_][ratelimiter]
  keys are also based on _Flow Labels_.

## Sources

Flows are annotated with _Flow Labels_ based on four sources: request labels,
baggage, flow classifiers, and explicit labels from the Aperture SDK call.

### Request labels

For each HTTP [_Control Point_][control-point] (where flows are HTTP or GRPC
requests), some basic metadata is available as _request labels_. These are
`http.method` , `http.target`, `http.host`, `http.scheme`,
`http.request_content_length` and `http.flavor`. Additionally all (non-pseudo)
headers are available as `http.request.header.header_name`, e.g.
`http.request.header.user_agent` (note the snake_case!). Values of these labels
are described by [OpenTelemetry semantic conventions for HTTP
spans][otel-conventions]. The only exception is `http.host` attribute, which is
equal to Host/Authority header. This is thus similar to `net.peer.name` OTEL
attribute.

### Baggage {#baggage}

Baggage propagation is a powerful concept that allows attaching metadata to a
whole request chain or to a whole [trace][traces]. If you already have baggage
propagation configured in your system, you can access the baggage as flow
labels. This is supported on service-mesh (Envoy) and web framework based
control point insertion.

- _HTTP_: Baggage is pulled from the [_baggage_][baggage] header.
- _Feature_: Baggage is automatically pulled from context on each `Check()`
  call. This is assuming you're using the OpenTelemetry library to manage the
  baggage.

Baggage members are mapped to _Flow Labels_ 1:1 – keys become label keys, values
become label values (properties are ignored).

Read more about baggage propagation on:
[Baggage | OpenTelemetry](https://opentelemetry.io/docs/concepts/signals/baggage/).

### Classifiers

When the labels you need are not already present in baggage, nor as request
labels, you can create a [classifier][classifier] to inject new labels into the
system. Since the Classifier also injects the label into baggage by default,
this means you can set or extract the label in a different place than where it
is consumed (assuming you have baggage propagation configured throughout your
system).

### Aperture SDK

The Aperture SDK, in addition to automatically pulling baggage from context,
also takes an explicit `labels` map in the `Check()` call.

## Live Previewing Flow Labels

You can discover the labels flowing through services and control points using
[`aperturectl`][aperturectl].

For example:

```sh
aperturectl flow-control preview --kube service1-demo-app.demoapp.svc.cluster.local ingress
```

Returns:

```json
{
  "samples": [
    {
      "labels": {
        "http.flavor": "1.1",
        "http.host": "service1-demo-app.demoapp.svc.cluster.local",
        "http.method": "POST",
        "http.request.header.content_length": "201",
        "http.request.header.content_type": "application/json",
        "http.request.header.cookie": "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
        "http.request.header.user_agent": "k6/0.42.0 (https://k6.io/)",
        "http.request.header.user_id": "14",
        "http.request.header.user_type": "bot",
        "http.request.header.x_forwarded_proto": "http",
        "http.request.header.x_request_id": "3958dad8-eb71-47f0-a9f6-500cccb097d2",
        "http.request_content_length": "0",
        "http.scheme": "http",
        "http.target": "/request",
        "user_type": "bot"
      }
    }
  ]
}
```

Alternatively, you can use the
[Introspection API](reference/api/agent/flow-preview-service-preview-flow-labels.api.mdx)
directly on an `aperture-agent` local to the service instances (pods):

```sh
curl -X POST localhost:8080/v1/flowcontrol/preview/labels/service1-demo-app.demoapp.svc.cluster.local/ingress?samples=1
```

## Telemetry

Telemetry data is extracted out of flows for further processing. This data is
collected from the following sources:

- Stream of access logs from service mesh (refer to [Istio
  Configuration][istio])
- Traces from [Aperture SDK][aperture-go]

Aperture uses OpenTelemetry's robust pipelining for receiving the telemetry data
and produce other streams of data from it.

### Metrics

Prometheus metrics are generated from the telemetry data that is received. Along
the path of the flows, telemetry data is tagged by the [flux meters][flux-meter]
and [workloads][workload] that matched.

### OLAP-style telemetry

OLAP-style telemetry data is generated as OpenTelemetry logs and is saved in an
OLAP database. This is done by creating multi-dimensional rollups from flow
labels.

OLAP-style telemetry doesn't work well with extremely high-cardinality labels,
thus if a extremely high-cardinality label is detected, some of its values may
be replaced with `REDACTED_VIA_CARDINALITY_LIMIT` string.

#### Default labels

These are protocol-level labels (e.g. http, network) extracted by the configured
service mesh/middleware and are available to be referenced in [Label
Matchers][label-matchers], except for a few high-cardinality ones.

#### Labels extracted from baggage

These are _Flow Labels_ mapped from [baggage](#baggage).

#### Labels defined by user

These are labels provided via Classifiers in case of service mesh/middleware
integration, or explicitly at flow creation in [Aperture SDK][aperture-go].

:::note

In the case of a clash, the flow Label will be applied in the following
precedence over:

1. User-defined
2. Baggage
3. Default

:::

## Interaction with FluxNinja ARC plugin {#plugin}

All the flow Labels are used as labels of flow events. These events are rolled
up and sent to the analytics database in the FluxNinja ARC. This allows:

- For the _Flow Labels_ to be used as filters or group bys
- To see analytics for each _Flow Label_, eg. distribution of its values

:::note

For _Classifier_ created labels, you can disable this behavior by setting
`hidden: true` in the [classification rule](/reference/policies/spec.md#rule).

:::

[flow-selector]: ./flow-selector.md
[classifier]: ./resources/classifier.md
[workload]: ./components/concurrency-limiter.md#workload
[ratelimiter]: ./components/rate-limiter.md
[scheduler]: ./components/concurrency-limiter.md#scheduler
[flux-meter]: ./resources/flux-meter.md
[baggage]: https://www.w3.org/TR/baggage/#baggage-http-header-format
[traces]:
  https://opentelemetry.io/docs/concepts/observability-primer/#distributed-traces
[control-point]: ./flow-selector.md#control-point
[otel-conventions]:
  https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
[aperture-go]: https://github.com/FluxNinja/aperture-go
[istio]: /get-started/integrations/flow-control/envoy/istio.md
[span]: https://opentelemetry.io/docs/reference/specification/trace/api/#span
[aperturectl]: /get-started/aperture-cli/aperture-cli.md
