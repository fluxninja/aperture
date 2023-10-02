---
title: FAQ
slug: faq
sidebar_label: FAQ
sidebar_position: 9
description: Frequently asked questions about Aperture.
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

### Does Aperture add latency on requests? {#request-overhead}

While Aperture does add some latency, it is a minimal one. Thanks to colocating
Aperture Agents with services, it's a single RPC call within a single node.

### If are the benefits of using Aperture over circuit breakers and rate limiting in EnvoyProxy? {#envoy-rate-limit}

While Envoy does have some local and non-local rate-limiting capabilities, there
are still benefits of using Aperture:

- Aperture [Rate Limiter][rate-limiter] allows dynamically configuring Rate
  Limiter parameters through signals from Policy.
- The ability to configure global rate limiting without configuring any external
  components
  – [mesh of Agents is providing distributed counters](/concepts/rate-limiter.md#distributed-counters).
- Rate-limiting decisions can be made locally on the Agent if lazy sync is
  enabled.
- In addition to Rate Limiter, Aperture also offers [Load
  Scheduler][load-scheduler], which Envoy doesn't have an equivalent of.

### Does Aperture reject requests immediately? {#reject-immediately}

- Rate Limiter always accepts or rejects immediately.
- [Load Scheduler][load-scheduler] can hold a request within some time period
  (derived from gRPC request timeout).
- Load Scheduler can also be configured in a way which effectively disables the
  queuing and scheduling logic. If such a configuration is desired, it will
  either accept or reject the request immediately based on workload priorities
  and other factors.

### If Aperture is rejecting or queuing requests, how will it impact the user experience? {#reject-impact}

Queuing requests should not affect user experience (apart from increased
latency). When it comes to rejecting requests, clients (whether it is front-end
code or some other service) should be prepared to receive
`429 Too Many Requests` or `503 Service Unavailable` response and react
accordingly.

Remember, that while some users receiving 503 might seem like a thing to avoid,
if such a case occurs, an overload is already happening and Aperture is
protecting your service from going into an unhealthy state.

### How can Flow Labels be defined for workload prioritization or rate limiting? {#flow-labels}

- In proxy- or web-framework-based Control Point insertion, most request
  metadata is already available as Flow Labels, for example
  `http.request.header.foo`.
- Already existing baggage is also available as Flow Labels.
- With SDKs, it's possible to explicitly pass Flow Labels to the Check call.
- Proxy-based integrations can use a [Classifier][classifier] to define new Flow
  Labels.

See the [Flow Label][flow-label] page for more details.

### How does Aperture work with existing auto-scaling? {#auto-scale}

As Aperture observes the system health, it can detect early sign of overload and
can take necessary actions to prevent the system from becoming unhealthy. While
auto-scaling could be running in parallel to add capacity, usually new instances
take some time to become healthy as they have to establish database connections,
perform service discovery, and so on. Therefore, Aperture is still needed to
protect the system from overload by queuing or dropping excessive load while
additional capacity is being added.

### Can the Aperture Controller run in a non-containerized environment? {#controller-bare-metal}

No, as for now, [Aperture Controller][aperture-controller] only runs on a
Kubernetes cluster. Remember that it's also possible to use the [Aperture Cloud
Controller][aperture-cloud-controller] instead of deploying your own.

### Can the Aperture Agent run in a non-containerized environment? {#agent-bare-metal}

Yes, the Aperture Agent can be deployed in a non-containerized environment. The
Aperture Agent is a binary that can be run on the
[Supported Linux platforms](/get-started/installation/supported-platforms.md).
The installation steps are available
[here](/get-started/installation/agent/bare-metal.md).

### What are Aperture Agent's performance numbers? {#agent-performance}

The Aperture Agent is designed to be lightweight and performant.

With the following setup:

- 1 node Kubernetes cluster
- 1 Aperture Agent installed as a
  [DaemonSet](/get-started/installation/agent/kubernetes/operator/daemonset.md)
- 1 policy with a [rate limiter][rate-limiter], a [load
  scheduler][load-scheduler] and a [flux meter][flux-meter]
- 3 services in `demoapp` namespace instrumented using
  [Istio Integration](/integrations/istio/istio.md)
- 5000 RPS at constant arrival rate over 30 minutes

The following results were observed:

|                | CPU (vCPU core)      | Memory (MB)         |
| -------------- | -------------------- | ------------------- |
| Aperture Agent | 0.783 mean, 1.02 max | 13.7 mean, 22.0 max |
| Istio Proxy    | 1.81 mean, 2.11 max  | 12.5 mean, 20.8 max |

[rate-limiter]: /concepts/rate-limiter.md
[load-scheduler]: /concepts/scheduler/load-scheduler.md
[flux-meter]: /concepts/flux-meter.md
[classifier]: /concepts/classifier.md
[flow-label]: /concepts/flow-label.md
[aperture-controller]: /architecture/architecture.md#aperture-controller
[aperture-cloud-controller]: /reference/fluxninja.md#cloud-controller
