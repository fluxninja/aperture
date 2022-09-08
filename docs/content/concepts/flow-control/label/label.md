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

For each _traffic_ control point (where flows are http or grpc requests), some
basic metadata is available as _request labels_. These are `request_id` ,
`request_method`, `request_path`, `request_host`, `request_scheme`,
`request_size`, `request_protocol` (mapped from fields of
[HttpRequest][authz-request-http]). Also, (non-pseudo) headers are available as
`request_header_<headername>`, where `<headername>` is a headername normalised
to lowercase, eg. `request_header_user-agent`.

### Baggage

Baggage propagation is a powerful concept that allows attaching metadata to a
whole request chain or to a whole [trace][traces]. If you already have baggage
propagation configured in your system, you can access the baggage as flow
labels. This is supported on both _traffic_ and _feature_ control points.

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
labels, you can create a [classifier](classifier) to inject new labels into the
system. Since the classifier also injects the label into baggage by default,
this means you can set or extract the label in a different place than where it
is consumed (assuming you have baggage propagation configured throughout your
system).

### Aperture library

The Aperture library, in addition to automatically pulling baggage from context,
also takes an explicit `labels` map in the `Check()` call.

## Interaction with FluxNinja Cloud plugin {#plugin}

All the flow labels except the request labels are used as labels of flow events.
These events are rolled up and sent to the analytics database in the cloud. This
allows:

- for the flow labels to be used as filters,
- to see analytics for each flow label, eg. distribution of its values.

:::note

For classifier-created labels, you can disable this behaviour by setting
`hidden: true` in
[the classification rule](/reference/configuration/policies.md#v1-rule).

:::

:::danger

This means that by default the already-present-in-baggage labels are sent to the
cloud.

TODO perhaps we should invert the default?

:::

[flow]: /concepts/flow-control/flow-control.md#flow
[selector]: /concepts/flow-control/selector.md
[actuators]: /concepts/flow-control/actuators/actuators.md
[scheduler]: /concepts/flow-control/actuators/scheduler.md
[workload]: /concepts/flow-control/actuators/scheduler.md#workload
[ratelimiter]: /concepts/flow-control/actuators/rate-limiter.md
[fluxmeter]: /concepts/flow-control/fluxmeter.md
[authz-request-http]:
  https://github.com/envoyproxy/envoy/blob/637a92a56e2739b5f78441c337171968f18b46ee/api/envoy/service/auth/v3/attribute_context.proto#L102
[baggage]: https://www.w3.org/TR/baggage/#baggage-http-header-format
[traces]:
  https://opentelemetry.io/docs/concepts/observability-primer/#distributed-traces
