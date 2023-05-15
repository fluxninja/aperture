---
title: FAQ
slug: faq
sidebar_position: 7
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

While Envoy does have some local and non-local rate limiting capabilities, there
are still benefits of using Aperture:

- Aperture Rate limiter allows dynamically configuring Rate Limiter parameters
  via signals from Policy.
- Ability to configure global rate limiting without configuring any external
  components
  – [mesh of Agents is providing distributed counters](/concepts/flow-control/components/rate-limiter.md#distributed-counters).
- Rate-limiting decisions can be made locally on agent, if lazy sync is enabled.
- In addition to Rate Limiter, Aperture also offers Load Scheduler, which Envoy
  doesn't have an equivalent of.

### Does Aperture reject requests immediately?

- Rate Limiter always accepts/rejects immediately.
- Load Scheduler can hold a request within some time period (derived from
  request's grpc-timeout)
- Load Scheduler can also be configured which effectively disables
  holding/scheduling part. If such configuration is desired, it will
  accept/reject request immediately based on workload priorities and other
  factors.

### If Aperture is rejecting or queuing requests, how will it impact the user experience?

Queuing requests should not affect user experience (apart of increased latency).
When it comes to rejecting requests, clients (whether it's frontend code or some
other service) should be prepared to receive 429 Too Many Requests or 503
Service Unavailable response and react accordingly. We're working on a library
to make it easy to handle these scenarios and provide nice UX for end users.

Remember that while receiving 503 by some of users may seem like a thing to
avoid, if such case occurs overload is already happening Aperture is protecting
your service from unhealthy state (eg. crashing) and thus affecting even more
users.

### How can we ensure the uniqueness of requests for the flow label when using Aperture?

Aperture uses the flow label to identify the requests.

TODO (I don't get the question tbh)

### How does Aperture address the issue of delays in servers becoming available and reaching a healthy state, particularly in the context of auto-scaling?

As Aperture observes the system, it can detect early sign of overload and can
take necessary actions to prevent the system from becoming unhealthy. Therefore,
server gets enough time to reach a healthy state.

It may happen that overload is happening too quickly for auto-scale to happen.
In such case, load scheduler will queue or drop excessive load to protect
existing services.

### Can the Aperture controller run in a non-containerized environment?

No, the Aperture controller runs in a containerized environment only.

### Is a Kubernetes cluster necessary for working with the Aperture Controller?

Yes, the Aperture controller runs on a Kubernetes cluster.
