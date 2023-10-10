---
title: OpenAI API Quota Scheduling
keywords:
  - policies
  - quota
  - prioritization
  - external-api
---

```mdx-code-block
import {apertureVersion} from '../../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Overview

Quota Scheduler can be use to enforce rate limits set by third party vendors
like OpenAI, which helps minimise the need of retrying the request and increase
success rate of requests. Moreover, it can help reduce the third party vendor
usages cost by scheduling the request, stopping to go beyond a certain limits.

This guide will help you understand how to use the Quota Scheduler Policy to
manage and stop hitting rate limits imposed by OpenAI. With the help of this
policy need of retrying the request can be optional, without having to sacrifice
the user experience.

## Pre-Requisites

Before, you begin with this guide, verify the pre-requisites are fulfilled.

- Aperture is installed and running. If not, follow the
  [get started guide](/get-started/get-started.md).
- `aperturectl` is installed and configured. If not, follow the
  [aperturectl installation guide](/get-started/installation/aperture-cli/aperture-cli.md).

## Configuration

Before creating a policy which can be used to schedule the request, you need to
create a control point where the policy will act on. There are multiple ways to
achieve this, for the scope of this guide, will be using the JavaScript SDK to
create a control point.

### Setup Control Point using JavaScript SDK

Aperture JavaScript allows you to set a control point manually. Aperture Client
instance is created earlier in the code which is not discussed in this guide.
Detailed information about Aperture JavaScript SDK can be found in
[Manually setting feature control points using JavaScript SDK](/integrations/sdk/javascript/manual.md)

```typescript
if (this.apertureClient) {
  const charCount =
    this.systemMessage.length +
    message.length +
    Strin("system" + "user").length;

  const labels: Record<string, string> = {
    api_key: CryptoES.SHA256(api.apiKey).toString(),
    // https://platform.openai.com/docs/guides/rate-limits/reduce-the-max_tokens-to-match-the-size-of-your-completions

    // also see - https://help.openai.com/en/articles/4936856-what-are-tokens-and-how-to-count-them

    estimated_tokens: (Math.ceil(charCount / 4) + responseTokens).toString(),

    model_variant: modelVariant,
    product_reason: this.settings.product_reason,

    priority: String(PRIORITIES[this.settings.product_reason] + priorityBump),
    prompt_type: promptType,
  };

  const apertureStart = Date.now();

  flow = await this.apertureClient.StartFlow("openai", {
    labels: labels,
    timeoutMilliseconds: 600000,
  });
  const apertureEnd = Date.now();

  logger.info(
    `OpenAI: aperture-js flow should run: ${
      flow.ShouldRun() ? "yes" : "no"
    }, response time: ${apertureEnd - apertureStart} ms${
      flow.CheckResponse()
        ? `, response: ${JSON.stringify(flow.CheckResponse())}`
        : ""
    }${
      flow.Error() ? `, error: ${JSON.stringify(flow.Error())}` : ""
    }, estimated tokens: ${
      labels.estimated_tokens
    }, character count: ${charCount}`,
  );
}
```

Let's understand the code snippet above, we are creating a control point named
'openai' and setting the labels which will be used by the policy to schedule the
request. The labels are used to identify the request and schedule it
accordingly. The labels are as follows:

- `api_key`: This is the api key used to authenticate the request to OpenAI.
- `estimated_tokens`: This is the estimated number of tokens which will be used
  by the request. This is calculated by adding the number of characters in the
  prompt and the response tokens.
- `model_variant`: This is the model variant used by the request.

These labels serve as a structured method to categorize and prioritize your
requests with enhanced accuracy. While you're encouraged to design labels that
resonate with your unique business needs, here are some example labels for
clarity:

- `product_reason`: Typically, businesses offer different tiers of their
  products to cater to various customer segments. For instance, a SaaS company
  might offer Basic, Premium, and Enterprise tiers, each with a different set of
  features and pricing. In this case, you can use the `product_reason` label to
  categorize your requests based on the tier of the product that the request is
  being made for. Or you can create different labels. Read more about
  [Flow Labels](/concepts/flow-label.md).

### Policy

Now, that we have a control point, let's generate a policy for control point
`openai`.

#### Generate a Values File

To generate a policy using Quota Scheduler Blueprint, we need to generate a
values file specific to the policy. This can be achieved using the command
provided below.

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints values --name=quota-scheduling/base --version={apertureVersion} --output-file=values.yaml</CodeBlock>
```

Values file need to be adjusted to match the application requirements -

- `policy_name`: Name of the policy. It is required.
- `bucket_capacity`: Bucket capacity. This value define how many request can be
  sent in a given interval. For example, if the bucket capacity is 4 & internal
  is 60, then 4 requests can be sent in 60 seconds.
- `fill_amount`: Fill amount. After the tokens are consumed, the bucket will be
  filled with this amount. For example, if bucket capacity is 4 & fill amount is
  4, then after 4 requests are sent, the bucket will be filled with 4 tokens. It
  will help tune how many requests to allow after the bucket is filled.
- `rate_limiter`:
  - `interval`: Interval at which the rate limiter will be filled. When to reset
    the bucket.
  - `label_key`: Label key to match the request against. This label key could be
    api key, user id, etc, which help determine the quota for the request. In
    this case, it is `api_key`.

Scheduler helps in prioritizing the requests based on the labels, and priority
defined. In this case, we are using `priority` label which is being passed by
Aperture SDK in code, containing the priority of the request.

- `scheduler`:
  - `priority_label_key`: Priority label key to match the request against. In
    this case, it is `priority`.
  - `workloads`:
    - `name`: To match the label value against name of workloads. In this case,
      it is `paid_user`, `trial_user`, `free_user`.
    - `label_matcher`:
    - `match_labels`: Labels to match the request against. In this case, it is
      `product_reason`.

Selectors parameters allows filtering of the requests to ensure where the policy
will act on.

- `selectors`:
  - `control_point`: Control point name to match the request against. In this
    case, it will be `openai`.
  - `agent_group`: Agent group name to match the request against. It is
    optional.
  - `label_matcher`:
  - `match_labels`: Labels to match the request against. It is optional.

Below is example of values file adjusted to match with sdk code snippet &
control point labels.

<details><summary>values.yaml</summary>
<p>

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/fluxninja/aperture/latest/blueprints/quota-scheduling/base/gen/definitions.json
# Generated values file for quota-scheduling/base blueprint
# Documentation/Reference for objects and parameters can be found at:
# https://docs.fluxninja.com/reference/blueprints/quota-scheduling/base

blueprint: quota-scheduling/base
uri: github.com/fluxninja/aperture/blueprints@latest
policy:
  # Name of the policy.
  # Type: string
  # Required: True
  policy_name: gpt-4-rpm
  quota_scheduler:
    # Bucket capacity.
    # Type: float64
    # Required: True
    bucket_capacity: 4
    # Fill amount.
    # Type: float64
    # Required: True
    fill_amount: 200
    # Rate Limiter Parameters.
    # Type: aperture.spec.v1.RateLimiterParameters
    # Required: True
    rate_limiter:
      interval: 60s
      label_key: api_key
      delay_initial_fill: true
    scheduler:
      priority_label_key: priority
      workloads:
        - name: paid_user
          label_matcher:
            match_labels:
              product_reason: paid_user
        - name: trial_user
          label_matcher:
            match_labels:
              product_reason: trial_user
        - name: free_user
          label_matcher:
            match_labels:
              product_reason: free_user
    # Flow selectors to match requests against
    # Type: []aperture.spec.v1.Selector
    # Required: True
    selectors:
      - control_point: openai
        agent_group: apollo-prod
        label_matcher:
          match_labels:
            model_variant: gpt-4
```

</p>
</details>

#### Generate Policy

Using the adjusted values file, a final policy will be generated, which will be
deployed.

To generate, use the following command:

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints generate --values-file=values.yaml --output-dir=policy-gen</CodeBlock>
```

#### Apply Policy

Apply the policy using the `aperturectl` CLI or `kubectl`.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl" label="aperturectl">
```

Pass the `--kube` flag with `aperturectl` to directly apply the generated policy
on a Kubernetes cluster in the namespace where the Aperture Controller is
installed.

```mdx-code-block
<CodeBlock language="bash">aperturectl apply policy --file=policy-gen/policies/gpt-4-rpm.yaml --kube </CodeBlock>
```

```mdx-code-block
</TabItem>
<TabItem value="kubectl" label="kubectl">
```

Apply the policy YAML generated (Kubernetes Custom Resource) using the above
example with `kubectl`.

```bash
kubectl apply -f policy-gen/configuration/gpt-4-rpm.yaml -n aperture-controller
```

```mdx-code-block
</TabItem>
</Tabs>
```

## How does it work?

Explanation of Policy & working

Product Screenshots

## What's Next?
