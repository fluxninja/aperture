---
title: External API Calls Prioritization with Quota Scheduling
keywords:
  - policies
  - quota
  - prioritization
  - external-api
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

The Quota scheduler is a versatile feature that can be used for managing
requests for external API calls, but also for inter-service communication within
a system architecture. In case of External API calls, it acts like a traffic
controller, making sure that too many requests don't rush in at once and cause a
traffic jam in the API, which could potentially bring it down. By utilizing
Aperture, external API calls, can be prioritized and managed under quota limits,
so that calls doesn't exceed the limit, causing penalty or even worse being
blocked by the API provider. The following example shows how to use the Quota
scheduler to manage external API calls while ensuring the api calls are
prioritized and managed under quota limits.

## Policy

This policy uses the
[`Quota Schedular`](/reference/policies/bundled-blueprints/policies/quota-scheduler.md)
blueprint that enables quota scheduling for workloads. In this example, we will
create a policy that will do quota based scheduling for external API and while
do so, it will also do the workload prioritization. We will continuously monitor
the quota checks panel and workload decision panel to see how the policy is
working, and workload being rejected or accepted based on the quota limits and
priority.

At a high-level, this policy consists of:

- Rate Limiter: Limiting the number of requests as they exceed a certain
  threshold.
- Workload Scheduler: Allowing workload to be scheduled based on priority.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/with-external-api-calls-prioritization/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>

```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/with-external-api-calls-prioritization/policy.yaml}
```

</p>
</details>

## Playground

The above policy can be loaded using the `quota-scheduler` scenario in
[Playground](/get-started/playground/playground.md)

:::info

[Circuit Diagram](./assets/with-external-api-calls-prioritization/graph.mmd.svg)
for this policy.

:::

<Zoom>

![Quota Scheduler With Workload Prioritization ](./assets/with-external-api-calls-prioritization/dashboard.png)

</Zoom>
