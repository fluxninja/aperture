---
title: Architecture
sidebar_position: 2
description: TODO
image: ../assets/img/aperture_logo.png
keywords:
  - reliability
  - overload
  - concurrency
  - aperture
  - fluxninja
  - microservices
  - cloud
  - TODO
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

<Zoom>

```mermaid
{@include: ../assets/diagrams/architecture/architecture_simple.mmd}
```

</Zoom>

## Aperture Agent

Aperture Agent is a component that is deployed as a sidecar next to your service
instances to provide flow control capabilities. It uses a weighted fair queuing
scheduler to prioritize workloads based on their importance, ensuring that
critical application features are not affected during overload scenarios. The
agent also includes a distributed rate-limiter, which is used to prevent abuse
and protect the service from malicious requests.

The agent uses an in-built telemetry system to collect metrics on service
performance and request attributes, such as customer tier or request type. These
metrics are then analyzed by the agent's high-fidelity flow classifier, which
labels requests based on their attributes. This allows the agent to make more
informed decisions about how to handle requests and prioritize workloads.

The agent continuously monitors service performance and adjusts its flow control
strategies accordingly. For example, it can shed load from non-critical
workloads to ensure that critical application features are not affected during
periods of high traffic. Additionally, the agent can use the data it collects to
identify and block malicious requests, ensuring the security and stability of
the service.

Overall, Aperture Agent works by continuously monitoring, analyzing and
actuating the flow of requests, prioritizing the important ones and shedding
load on non-important ones to ensure the performance and reliability of the
service in web-scale apps.

## Aperture Controller

Aperture Controller is the component of Aperture that manages the flow control
policies for the entire system. It uses dataflow-driven policies to continuously
monitor and adjust service-level objectives (SLOs) for all services in the
system. These policies are expressed as circuits, much like circuit networks in
the game Factorio, and can be customized to meet the specific needs of the
application.

The controller receives metrics from the agents deployed next to each service
instance, and uses these metrics to calculate recovery or escalation actions.
For example, it can use a gradient control circuit component to implement an
Additive Increase, Multiplicative Decrease (AIMD) style counter-measure that
limits the concurrency on a service when response times deteriorate. More
advanced control components like PID can be used to further tune the concurrency
limits.

The controller also has the capability to make decisions based on the priority
of requests and workloads. For example, it can prioritize high-importance
features over others, ensuring that critical application features are not
affected during periods of high traffic.

Overall, the Aperture Controller works by continuously monitoring the entire
system and adjusting the flow control policies to ensure that the service-level
objectives are met, and the application remains stable and reliable. It ensures
that the correct decisions are made for prioritizing requests and workloads and
taking action to prevent failures.
