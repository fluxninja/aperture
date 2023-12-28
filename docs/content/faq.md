---
title: FAQ
slug: faq
sidebar_label: FAQ
sidebar_position: 10
description: Frequently asked questions about Aperture.
image: /assets/img/aperture_logo.png
keywords:
  - reliability
  - overload
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

While Aperture does add some latency, it is a minimal one. The latency of
requests to Aperture Cloud is in the order of 10 to 20 ms within the same region
as your application. If you are hosting Aperture Agents yourself, the latency is
in the order of 1-2ms.

### Does Aperture reject requests immediately? {#reject-immediately}

- Rate Limiter always accepts or rejects immediately.
- [Scheduler][scheduler] queues the request for a maximum time up to the gRPC
  timeout of the request to Aperture with a small deadline margin.

### If Aperture is rejecting or queuing requests, how will it impact the user experience? {#reject-impact}

Queuing requests should not affect user experience (apart from increased
latency). When it comes to rejecting requests, clients (whether it is front-end
code or some other service) should be prepared to receive
`429 Too Many Requests` or `503 Service Unavailable` response and react
accordingly.

Remember, that some users receiving 503 means that the service is overloaded and
Aperture is protecting your service from going into an unhealthy state.

### How can Flow Labels be defined for workload prioritization or rate limiting? {#flow-labels}

- With SDKs, it's possible to explicitly pass Flow Labels to the Check call.
- In proxy- or web-framework-based Control Point insertion, most request
  metadata is already available as Flow Labels, for example
  `http.request.header.foo`.
- Already existing baggage is also available as Flow Labels.
- Proxy-based integrations can use a [Classifier][classifier] to define new Flow
  Labels.

See the [Flow Label][flow-label] page for more details.

### How does Aperture work with existing auto-scaling? {#auto-scale}

Rate limiting and caching allows services to stay performant while being
cost-effective. Aperture enables developers to bring these capabilities to their
service through a single convenient API.

Auto-scaling is used when the service is nearing peak capacity, despite rate
limiting and caching. But scaling a service can be slow and expensive. While
auto-scaling is happening, Aperture can protect the service from overload by
queuing and prioritizing requests while staying within capacity. This also
reduces the need to always stay over-provisioned.

### Can you host Aperture in your infrastructure? {#self-host}

Yes, Aperture is fully open source and can be hosted on your infrastructure.
There are two possible deployment options:

1. Install just the Agents and connect the Agents to Aperture Cloud.
2. Install the Agents and the Aperture Controller and connect to Aperture Cloud
   only for sending the telemetry.

### Can the Aperture Agent run in a non-containerized environment? {#agent-bare-metal}

Yes, the Aperture Agent can be deployed in a non-containerized environment. The
Aperture Agent is a binary that can be run on the
[Supported Linux platforms](/aperture-for-infra/supported-platforms.md). The
installation steps are available
[here](/aperture-for-infra/agent/bare-metal.md).

Note: Aperture Cloud provides a hosted Agent for SDK integration, allowing you
to use it by API instead of deploying your own Aperture Agents.

### What are Aperture Agent's performance numbers? {#agent-performance}

The Aperture Agent is designed to be lightweight and performant.

With the following setup:

- 1 node Kubernetes cluster
- 1 Aperture Agent installed as a
  [DaemonSet](/aperture-for-infra/agent/kubernetes/operator/daemonset.md)
- 1 policy with a [rate limiter][rate-limiter], a [load
  scheduler][load-scheduler] and a [flux meter][flux-meter]
- 3 services in `demoapp` namespace instrumented using
  [Istio Integration](/aperture-for-infra/integrations/istio/istio.md)
- 5000 RPS at constant arrival rate over 30 minutes

The following results were observed:

|                | CPU (vCPU core)      | Memory (MB)         |
| -------------- | -------------------- | ------------------- |
| Aperture Agent | 0.783 mean, 1.02 max | 13.7 mean, 22.0 max |
| Istio Proxy    | 1.81 mean, 2.11 max  | 12.5 mean, 20.8 max |

[rate-limiter]: /concepts/rate-limiter.md
[load-scheduler]: /concepts/request-prioritization/load-scheduler.md
[scheduler]: /concepts/scheduler.md
[flux-meter]: /concepts/advanced/flux-meter.md
[classifier]: /concepts/advanced/classifier.md
[flow-label]: /concepts/flow-label.md
