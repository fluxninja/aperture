---
title: Flow Label
sidebar_position: 2
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
```

A flow is the fundamental unit of work from the perspective of an Aperture
agent. It could be an API call, a feature, or even a database query. A flow in
Aperture is similar to [OpenTelemetry span][span]. Each flow is annotated with a
set of **Flow Labels**, which are key-value pairs. For example, if a flow is
tagged with `user_tier:gold`, then `user_tier` is the label key and `gold` is
the label value.

The following visualization depicts flows belonging to different
[control points](./control-point.md) across a distributed application, along
with their associated flow labels.

<Zoom>

```mermaid
{@include: ./assets/gen/flow-label/labels.mmd}
```

</Zoom>

## Usage of Flow Labels

Flow labels play a significant role in Aperture and serve various purposes:

- The [_Selector_][selectors] can use flow labels to select specific flows,
  narrowing down the scope of [_Classifiers_][classifier], [_Flux
  Meters_][flux-meter], and other components.
- Flow labels help map a flow to a specific [Workload][workload], providing
  context and aiding in workload identification.
- Flow labels are instrumental in mapping flows to unique users or entities,
  allowing for the allocation of rate limit quotas per label key. This
  allocation is utilized by the [Rate Limiter][ratelimiter] and [Quota
  Scheduler][quota-scheduler] components, enabling effective management of
  traffic and resource allocation based on these labels.

## Flow Label Sources

Flows are annotated with _Flow Labels_ based on four sources:

- Request labels
- Baggage
- Classifiers
- Aperture SDKs

### Request labels

For each HTTP [_Control Point_][control-point] (where flows are HTTP or gRPC
requests), some basic metadata is available as _request labels_. These are
`http.method`, `http.target`, `http.host`, `http.scheme`,
`http.request_content_length` and `http.flavor`. Additionally, all (non-pseudo)
headers are available as `http.request.header.header_name`, e.g.
`http.request.header.user_agent` (note the `snake_case`!).

The values of these labels are described by [OpenTelemetry semantic conventions
for HTTP spans][otel-conventions]. The only exception is the `http.host`
attribute, which is equal to the host or authority header. This is similar to
the `net.peer.name` OTel attribute.

### Baggage {#baggage}

[Baggage propagation][otel-baggage] allows attaching metadata to a whole request
chain or to a whole [trace][traces]. If you already have baggage propagation
configured in your system, you can access the baggage as flow labels. This is
supported on service-mesh (Envoy) and web framework-based control point
insertion.

- _HTTP_: Baggage is pulled from the [_baggage_][baggage] header.
- _Feature_: Baggage is automatically pulled from context on each `Check()`
  call. This is assuming you're using the OpenTelemetry library to manage the
  baggage.

Baggage members are mapped to _Flow Labels_ 1:1–keys become label keys, values
become label values (properties are ignored).

:::info Read more about baggage propagation

[Baggage | OpenTelemetry](https://opentelemetry.io/docs/concepts/signals/baggage/).

:::

### Classifiers

When the labels you need aren't already present in baggage, nor as request
labels, you can create a [_Classifier_][classifier] to inject new labels into
the system. Since the Classifier also injects the label into baggage by default,
this means you can set or extract the label in a different place than where it
is consumed (assuming you have baggage propagation configured throughout your
system).

### Aperture SDKs

The Aperture SDKs, in addition to automatically using baggage from context, also
takes an explicit `labels` map in the `Check()` call.

## Live Preview of Flow Labels

Discover labels flowing through services and control points using
[`aperturectl`][aperturectl].

For example:

```sh
aperturectl flow-control preview --kube checkout-app.checkout-service.svc.cluster.local ingress
```

Returns:

```json
{
  "samples": [
    {
      "labels": {
        "http.flavor": "1.1",
        "http.host": "checkout-app.checkout-service.svc.cluster.local",
        "http.method": "POST",
        "http.request.header.content_length": "201",
        "http.request.header.content_type": "application/json",
        "http.request.header.cookie": "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
        "http.request.header.user_agent": "k6/0.42.0 (https://k6.io/)",
        "http.request.header.user_id": "14",
        "http.request.header.user_type": "guest",
        "http.request.header.x_forwarded_proto": "http",
        "http.request.header.x_request_id": "3958dad8-eb71-47f0-a9f6-500cccb097d2",
        "http.request_content_length": "0",
        "http.scheme": "http",
        "http.target": "/request",
        "user_type": "guest"
      }
    }
  ]
}
```

Alternatively, you can use the
[Introspection API](/reference/api/agent/flow-preview-service-preview-flow-labels.api.mdx)
directly on a `aperture-agent` local to the service instances (pods):

```sh
curl -X POST localhost:8080/v1/flowcontrol/preview/labels/checkout-app.checkout-service.svc.cluster.local/ingress?samples=1
```

## Telemetry and Flow Labels

Telemetry data is extracted out of flows for further processing. This data is
collected from the following sources:

- Traces from [SDKs][aperture-sdks]
- Stream of access logs from service mesh (refer to [Istio
  Configuration][istio])
- Stream of access logs from [API Gateways][gateways]

Aperture uses OpenTelemetry's robust pipelines for receiving the telemetry data
and producing other streams of data from it.

### Metrics

Prometheus metrics are generated from the received telemetry data. Along the
path of the flows, telemetry data is tagged by the [Flux Meters][flux-meter] and
[workloads][workload] that matched.

### OLAP style telemetry

OLAP style telemetry data is generated as OpenTelemetry logs and is saved in an
OLAP database. This is done by creating multidimensional roll ups from flow
labels.

OLAP style telemetry does not work well with extremely high-cardinality labels,
therefore, if an extremely high-cardinality label is detected, some of its
values might be replaced with the `REDACTED_VIA_CARDINALITY_LIMIT` string.

#### Default labels

These are protocol-level labels (For example: HTTP, network) extracted by the
configured service mesh/middleware and are available to be referenced in [Label
Matcher][label-matcher], except for a few high-cardinality ones.

#### Labels extracted from baggage

These are _Flow Labels_ mapped from [baggage](#baggage).

#### Labels defined by user

These are labels provided by _Classifiers_ in case of service mesh/middleware
integration, or explicitly at flow creation in [Aperture SDK][aperture-sdks].

### Label Precedence

In case of a clash, the Flow Label will be applied in the following precedence
order:

1. User-defined
2. Baggage
3. Default

## Interaction with FluxNinja Extension {#extension}

All the flow Labels are used as labels of flow events. These events are rolled
up and sent to the analytics database in the Aperture Cloud. This allows:

- For the _Flow Labels_ to be used as filters or group-by
- To see analytics for each _Flow Label_, for example: distribution of its
  values

:::note

For _Classifier_ created labels, you can disable this behavior by setting
`hidden: true` in the
[classification rule](/reference/configuration/spec.md#rule).

:::

[selectors]: ./selector.md
[classifier]: ./advanced/classifier.md
[workload]: ./scheduler.md#workload
[ratelimiter]: ./rate-limiter.md
[quota-scheduler]: ./request-prioritization/quota-scheduler.md
[flux-meter]: ./advanced/flux-meter.md
[baggage]: https://www.w3.org/TR/baggage/#baggage-http-header-format
[traces]:
  https://opentelemetry.io/docs/concepts/observability-primer/#distributed-traces
[control-point]: ./control-point.md
[otel-conventions]:
  https://github.com/open-telemetry/opentelemetry-specification/blob/v1.25.0/specification/trace/semantic_conventions/http.md
[aperture-sdks]: /sdk/sdk.md
[gateways]: /aperture-for-infra/integrations/gateway/gateway.md
[istio]: /aperture-for-infra/integrations/istio/istio.md
[span]: https://opentelemetry.io/docs/reference/specification/trace/api/#span
[aperturectl]: ../reference/aperture-cli/aperturectl/flow-control/preview/
[label-matcher]: ./selector.md#label-matcher
[otel-baggage]: https://opentelemetry.io/docs/concepts/signals/baggage/
