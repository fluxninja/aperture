---
title: Service Protection
keywords:
  - tutorial
sidebar_position: 2
sidebar_label: Service Protection
---

## Overview

Service Protection is an approach tailored to modulate the rate of requests,
taking into consideration the current health status of the service. By serving
as an immediate and cost-effective alternative to auto-scaling, this policy acts
as the first line of defense against overwhelming service loads by limiting the
number of concurrent requests to a service.

<Zoom>

```mermaid
{@include: ../assets/service-protection.mmd}
```

</Zoom>

The provided graph illustrates a simplified version of the Service Protection
Policy. The policy is housed within the controller, which analyzes health
signals in real-time. Based on these metrics, it calculates an adjusted request
rate which is relayed to the Agent. The Agent then imposes a rate limit on the
service as required.

:::note

Aperture facilitates the observation of health signals from various services.
For instance, service protection can also be implemented based on the health
observation of an upstream service in relation to a downstream service.

:::

## Real World Scenario

For Service Protection, consider an e-commerce platform during a major sale
event. To handle the increased traffic, the service protection policy will
monitor the health of the service and modulate the request rate. This ensures
smooth operation by preventing service overloads that could lead to crashes,
therefore providing a consistent user experience.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
