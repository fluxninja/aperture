---
title: Introduction
slug: /
sidebar_position: 1
sidebar_label: Introduction
sidebar_class_name: introduction
keywords:
  - cloud
  - enterprise
  - platform
  - fluxninja
  - aperture
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
```

Aperture is a distributed load management platform designed for rate limiting,
caching, and prioritizing requests in cloud applications. Built upon a
foundation of distributed counters, observability, and a global control plane,
it provides a comprehensive suite of load management capabilities. These
capabilities enhance the reliability and performance of cloud applications,
while also optimizing cost and resource utilization.

![Unified Load Management](./assets/img/unified-load-management-light.svg#gh-light-mode-only)
![Unified Load Management](./assets/img/unified-load-management-dark.svg#gh-dark-mode-only)

Integrating Aperture in your application through SDKs is a simple 3-step
process:

- **Define labels**: Define labels to identify users, entities, or features
  within your application. For example, you can define labels to identify
  individual users, features, or API endpoints.

<!-- vale off -->
<!-- markdownlint-disable -->

  <details>
  <summary>Example</summary>

```typescript
// Tailor policies to get deeper insights into your workload with labels that
// capture business context.
const labels = {
  // You can rate limit each user individually.
  user: "jack",
  // And have different rate limits for different tiers of users.
  tier: "premium",
  // You can also provide the tokens for each request.
  // Tokens are flexible: LLM AI tokens in a prompt, complexity of a request,
  // number of sub-actions, etc.
  tokens: "200",
  // When peak load exceeds external quotas or infrastructure capacity,
  // requests can be throttled and prioritized.
  priority: HIGH,
  // Get deep insights into your workload. You can slice and dice performance
  // metrics by any label.
  workload: "/chat",
};
```

<!-- markdownlint-enable -->
<!-- vale on -->

  </details>

- **Wrap your workload**: Wrap your workload with `startFlow` and `endFlow`
  calls to establish control points around specific features or code sections
  inside your application. For example, you can wrap your API endpoints with
  Aperture SDKs to limit the number of requests per user or feature.

<!-- vale off -->
<!-- markdownlint-disable -->

  <details>
  <summary>Example</summary>

```typescript
// Wrap your workload with startFlow and endFlow calls, passing in the
// labels you defined earlier.
const flow = await apertureClient.startFlow("your_workload", {
  labels: labels,
  // Lookup result cache key to retrieve a cached result.
  resultCacheKey: queryParams,
});

// If rate or quota limit is not exceeded, the workload is executed.
if (flow.shouldRun()) {
  // Return a cached result or execute the workload.
  const cachedResult = flow.resultCache();
  const result = await yourWorkload(cachedResult);
  flow.setResultCache({
    value: result,
    ttl: { seconds: 86400, nanos: 0 },
  });
}
//
```

  </details>

<!-- markdownlint-enable -->
<!-- vale on -->

- **Configure & monitor policies**: Configure policies to control the rate,
  concurrency, and priority of requests.

<!-- vale off -->
<!-- markdownlint-disable -->

  <details>
  <summary>Policy YAML</summary>

```yaml
blueprint: rate-limiting/base
uri: github.com/fluxninja/aperture/blueprints@latest
policy:
  policy_name: rate_limit
  rate_limiter:
    bucket_capacity: 60
    fill_amount: 60
    parameters:
      interval: 3600s
      limit_by_label_key: user
    selectors:
      - control_point: your_workload
        label_matcher:
          match_list:
            - key: tier
              operator: In
              values:
                - premium
```

  </details>

<!-- markdownlint-enable -->
<!-- vale on -->

![Rate Limiter Blueprint](./get-started/assets/rate-limiter-blueprint-dark.png#gh-dark-mode-only)
![Rate Limiter Blueprint](./get-started/assets/rate-limiter-blueprint-light.png#gh-light-mode-only)
![Rate Limiter Dashboard](./get-started/assets/rate-limiter-dashboard-dark.png#gh-dark-mode-only)
![Rate Limiter Dashboard](./get-started/assets/rate-limiter-dashboard-light.png#gh-light-mode-only)

In addition to language SDKs, Aperture also integrates with existing control
points such as API gateways, service meshes, and application middlewares.

Aperture is available as a managed service, [Aperture Cloud][cloud], or can be
[self-hosted][self-hosted] within your infrastructure. Visit the
[Architecture][architecture] page for more details.

:::info Sign up

To sign up to Aperture Cloud, [click here][sign-up].

:::

## ⚙️ Load management capabilities {#load-management-capabilities}

- ⏱️ [**Global Rate-Limiting**](concepts/rate-limiter.md): Safeguard APIs and
  features against excessive usage with Aperture's high-performance, distributed
  rate limiter. Identify individual users or entities by fine-grained labels.
  Create precise rate limiters controlling burst-capacity and fill-rate tailored
  to business-specific labels. Refer to the
  [Rate Limiting](guides/per-user-rate-limiting.md) guide for more details.
- 📊
  [**API Quota Management**](concepts/request-prioritization/quota-scheduler.md):
  Maintain compliance with external API quotas with a global token bucket and
  smart request queuing. This feature regulates requests aimed at external
  services, ensuring that the usage remains within prescribed rate limits and
  avoids penalties or additional costs. Refer to the
  [API Quota Management](guides/api-quota-management.md) guide for more details.
- 🚦
  [**Concurrency Control and Prioritization**](concepts/request-prioritization/concurrency-scheduler.md):
  Safeguard against abrupt service overloads by limiting the number of
  concurrent in-flight requests. Any requests beyond this limit are queued and
  let in based on their priority as capacity becomes available. Refer to the
  [Concurrency Control and Prioritization](guides/concurrency-control-and-prioritization.md)
  guide for more details.
- 🎯 [**Workload Prioritization**](concepts/scheduler.md): Safeguard crucial
  user experience pathways and ensure prioritized access to external APIs by
  strategically prioritizing workloads. With
  [weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing),
  Aperture aligns resource distribution with business value and urgency of
  requests. Workload prioritization applies to API Quota Management and Adaptive
  Queuing use cases.
- 💾 [**Caching**](concepts/cache.md): Boost application performance and reduce
  costs by caching costly operations, preventing duplicate requests to
  pay-per-use services, and easing the load on constrained services. Refer to
  the [Caching](guides/caching.md) guide for more details.

## ✨ Get started {#get-started}

- [**Get Started**](get-started/get-started.md)
- [**Guides**](guides/guides.md)

## 📖 Learn {#learn}

The [Concepts](concepts/concepts.md) section provides detailed insights into
essential elements of Aperture's system and policies, offering a comprehensive
understanding of their key components.

## Additional Support

Don't hesitate to engage with us for any queries or clarifications. Our team is
here to assist and ensure that your experience with Aperture is smooth and
beneficial.

<!-- vale off -->

[**💬 Consult with an expert**](https://calendly.com/fluxninja/fluxninja-meeting)
|
[**👥 Join our Slack Community**](https://join.slack.com/t/fluxninja-aperture/shared_invite/zt-1vm2t2yjb-AG8rzKkB5TpPmqihJB6YYw)
| ✉️ Email: [**support@fluxninja.com**](mailto:support@fluxninja.com)

<!-- vale on -->

[cloud]: https://www.fluxninja.com
[sign-up]: https://app.fluxninja.com/sign-up
[architecture]: /aperture-for-infra/architecture.md
[self-hosted]: /aperture-for-infra/aperture-for-infra.md
