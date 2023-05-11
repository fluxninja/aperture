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

### Does the usage of Aperture entail extra overhead on requests?

Aperture does not add any extra overhead. Aperture keeps the latency of the
request as low as possible.

### If we already have circuit breakers and rate limiting in EnvoyProxy, what are the benefits of using Aperture?

Yes, EnvoyProxy does have, but Aperture does it better by having a global view
of the system. Aperture can be used in conjunction with Envoy to provide a
better experience.

### Does Aperture reject requests immediately?

No, Aperture does not reject requests immediately. Aperture queues the requests
or rejects them based on the configuration.

### If Aperture is rejecting or queuing requests, how will it impact the user experience?

Aperture queues the requests or rejects them based on the configuration. The
user experience is decided based on the configuration. If you want to prioritize
the user experience, you can configure Aperture to queue the requests. If you
want to prioritize the system health, you can configure Aperture to reject the
requests.

### How can we ensure the uniqueness of requests for the flow label when using Aperture?

Aperture uses the flow label to identify the requests.

### How does Aperture address the issue of delays in servers becoming available and reaching a healthy state, particularly in the context of auto-scaling?

As Aperture observes the system, it can detect early sign of overload and can
take necessary actions to prevent the system from becoming unhealthy. Therefore,
server gets enough time to reach a healthy state.

### Can the Aperture controller run in a non-containerized environment?

No, the Aperture controller runs in a containerized environment only.

### Is a Kubernetes cluster necessary for working with the Aperture Controller?

Yes, the Aperture controller runs on a Kubernetes cluster.
