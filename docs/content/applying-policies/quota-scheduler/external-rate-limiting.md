---
title: External Rate Limiting
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

The Quota Scheduler is a multifunctional tool that plays a crucial role not only
in handling requests for external APIs, but also in managing inter-service
dialogue within a system layout. Leveraging Aperture for external rate limiting
allows for the prioritization and regulation of API calls in line with quota
restrictions, ensuring these calls don't go beyond the established limit. This
prevents potential fines or blockage by the API provider. This is a practical
example of how the Quota Scheduler can be deployed for external rate limiting
while maximizing the utility of the external rate limits and keeping your
applications within their budget constraints, thus preventing any additional
costs.

## Policy

This policy uses the
[`Quota Schedular`](../../reference/policies/bundled-blueprints/policies/quota-scheduler.md)
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
[Playground](/playground/playground.md)

:::info

[Circuit Diagram](./assets/with-external-api-calls-prioritization/graph.mmd.svg)
for this policy.

:::

<Zoom>

![Quota Scheduler With Workload Prioritization ](./assets/with-external-api-calls-prioritization/dashboard.png)

</Zoom>
