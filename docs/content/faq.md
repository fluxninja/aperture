---
title: FAQ
slug: faq
sidebar_position: 11
description: Frequently asked questions about FluxNinja Aperture.
image: /assets/img/aperture_logo.png
keywords:
  - reliability
  - overload
  - concurrency
  - aperture
  - fluxninja
  - microservices
  - cloud
  - auto-scale
  - load management
  - flow control
  - faq
---

### Does the usage of Aperture entail extra overhead on requests? {#request-overhead}

While Aperture does add some latency overhead, it's a minimal one. Thanks to
colocating agents with services, it's a single RPC call within a single node.

### If we already have circuit breakers and rate limiting in EnvoyProxy, what are the benefits of using Aperture? {#envoy-rate-limit}

While Envoy does have some local and non-local rate-limiting capabilities, there
are still benefits of using Aperture:

- Aperture [Rate Limiter][] allows dynamically configuring Rate Limiter
  parameters through signals from Policy.
- The ability to configure global rate limiting without configuring any external
  components
  – [mesh of Agents is providing distributed counters](/concepts/flow-control/components/rate-limiter.md#distributed-counters).
- Rate-limiting decisions can be made locally on the agent if lazy sync is
  enabled.
- In addition to Rate Limiter, Aperture also offers Load Scheduler, which Envoy
  doesn't have an equivalent of.

### Does Aperture reject requests immediately?

- Rate Limiter always accepts or rejects immediately.
- [Load Scheduler][] can hold a request within some time period (derived from
  gRPC request timeout).
- Load Scheduler can also be configured in a way which effectively disables the
  holding/scheduling part. If such a configuration is desired, it will either
  accept or reject the request immediately based on workload priorities and
  other factors.

### If Aperture is rejecting or queuing requests, how will it impact the user experience?

Queuing requests should not affect user experience (apart from increased
latency). When it comes to rejecting requests, clients (whether it's frontend
code or some other service) should be prepared to receive 429 Too Many Requests
or 503 Service Unavailable response and react accordingly.

Remember that while receiving 503 by some of the users might seem like a thing
to avoid, if such a case occurs an overload is already happening and Aperture is
protecting your service from an unhealthy state (for example crashing) and
therefore affecting even more users.

### How can we define Flow Labels for workload prioritization or rate limiting?

- In proxy- or web-framework-based Control Point insertion, most request
  metadata is already available as Flow Labels, for example
  `http.request.header.foo`.
- Already existing baggage is also available as Flow Labels.
- With SDKs, it's possible to explicitly pass Flow Labels to the Check call.
- Proxy-based integrations can use a [Classifier][] to define new Flow Labels.

See the [Flow Label][] page for more details.

### How does Aperture address the issue of delays in servers becoming available and reaching a healthy state, particularly in the context of auto-scaling?

As Aperture observes the system, it can detect early sign of overload and can
take necessary actions to prevent the system from becoming unhealthy. Therefore,
the server gets enough time to reach a healthy state.

It might happen that overload is happening too quickly for auto-scale to happen.
In such case, the Load Scheduler will queue or drop excessive load to protect
existing services.

### Can the Aperture Controller run in a non-containerized environment?

No, as for now, we only support deploying [Aperture Controller][] on a
Kubernetes cluster.

[Rate Limiter]: /concepts/flow-control/components/rate-limiter.md
[Load Scheduler]: /concepts/flow-control/components/load-scheduler.md
[Classifier]: /concepts/flow-control/resources/classifier.md
[Flow Label]: /concepts/flow-control/flow-label.md
[Aperture Controller]: /get-started/installation/controller/controller.md
