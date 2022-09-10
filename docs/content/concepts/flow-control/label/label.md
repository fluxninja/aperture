---
title: Flow Label
sidebar_position: 0.75
---

Every [flow][flow] is annotated with a set of **flow labels**. Each flow label
is a key:value pair. Eg. if a flow is annotated with `user_tier:gold` label,
then `user_tier` is a label key and `gold` is a label value.

Flow labels are used used in different ways in Aperture:

- [Flow selector][selector] can select flows based on flow labels, thus flow
  labels can be used to narrow the scope of [_Actuators_][actuators] or
  [_FluxMeters_][fluxmeter]
- Flow labels are used to classify a flow to a [_workload_][workload]
- Fairness within a scheduler's workload and [rate-limiting][ratelimiter] keys
  are also based on flow labels

## Sources

Flows are annotated with flow labels based on four sources: request labels,
baggage, flow classifiers and explicit labels from the Aperture library call.

### Request labels

For each _traffic_ [control point][control-point] (where flows are http or grpc
requests), some basic metadata is available as _request labels_. These are
`http.method` , `http.target`, `http.host`, `http.scheme`,
`http.request_content_length` and `http.flavor`. Additionally all (non-pseudo)
headers are available as `http.request.header.header_name`, eg.
`http.request.header.user_agent` (note the snake_case!). Values of these labels
are described by [OpenTelemetry semantic conventions for HTTP
spans][otel-conventions]. The only exception is `http.host` attribute, which is
equal to Host/Authority header. This is thus similar to `net.peer.name` OTEL
attribute, but is provided for both ingress and egress control points.

### Baggage

Baggage propagation is a powerful concept that allows attaching metadata to a
whole request chain or to a whole [trace][traces]. If you already have baggage
propagation configured in your system, you can access the baggage as flow
labels. This is supported on both _traffic_ and _feature_ [control
points][control-point].

- _traffic_: Baggage is pulled from the [_baggage_][baggage] header,
- _feature_: Baggage is automatically pulled from context on each `Check()`
  call. This is assuming you're using the OpenTelemetry library to manage the
  baggage.

Baggage members are mapped to flow labels 1:1 – keys become label keys, values
become label values (properties are ignored).

Read more about baggage propagation on:
[Baggage | OpenTelemetry](https://opentelemetry.io/docs/concepts/signals/baggage/).

### Flow classifiers

When the labels you need are not already present in baggage, nor as request
labels, you can create a [classifier](classifier.md) to inject new labels into
the system. Since the classifier also injects the label into baggage by default,
this means you can set or extract the label in a different place than where it
is consumed (assuming you have baggage propagation configured throughout your
system).

### Aperture library

The Aperture library, in addition to automatically pulling baggage from context,
also takes an explicit `labels` map in the `Check()` call.

## Interaction with FluxNinja Cloud plugin {#plugin}

All the flow labels are used as labels of flow events. These events are rolled
up and sent to the analytics database in the cloud. This allows:

- for the flow labels to be used as filters,
- to see analytics for each flow label, eg. distribution of its values.

:::note

For classifier-created labels, you can disable this behaviour by setting
`hidden: true` in
[the classification rule](/reference/configuration/policies.md#v1-rule).

:::

:::caution

This means that by default the already-present-in-baggage labels are sent to the
cloud. If this is not what you want,
[we'll be providing a way](https://github.com/fluxninja/aperture/issues/376) to
select which labels to include in telemetry.

:::

[flow]: /concepts/flow-control/flow-control.md#flow
[selector]: /concepts/flow-control/selector.md
[actuators]: /concepts/flow-control/actuators/actuators.md
[scheduler]: /concepts/flow-control/actuators/concurrency-limiter.md
[workload]: /concepts/flow-control/actuators/concurrency-limiter.md#workload
[ratelimiter]: /concepts/flow-control/actuators/rate-limiter.md
[fluxmeter]: /concepts/flow-control/fluxmeter.md
[baggage]: https://www.w3.org/TR/baggage/#baggage-http-header-format
[traces]:
  https://opentelemetry.io/docs/concepts/observability-primer/#distributed-traces
[control-point]: ../flow-control.md#control-point
[otel-conventions]:
  https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
