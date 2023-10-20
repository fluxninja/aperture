---
title: Managing OpenAI API Rate Limits
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

## Understanding OpenAI rate limits

OpenAI imposes
[fine-grained rate limits](https://platform.openai.com/docs/guides/rate-limits/overview)
on both requests per minute and tokens per minute for each AI model they offer.
For example:

| Model             | Tokens per minute | Requests per minute |
| ----------------- | ----------------- | ------------------- |
| gpt-3.5-turbo     | 90000             | 3500                |
| gpt-3.5-turbo-16k | 180000            | 3500                |
| gpt-4             | 40000             | 200                 |

Dealing with these limits can be tricky, as API calls to OpenAI have high
latency (several seconds). As the rate limits are quite aggressive, back-off and
retry loops run for several minutes before a request can be successfully sent to
OpenAI. When working with OpenAI APIs, request prioritization can be quite
beneficial to ensure that the most important requests are sent to OpenAI first.

## Managing OpenAI rate limits with Aperture

Aperture can help manage OpenAI rate limits and improve user experience by
queuing and prioritizing requests before sending them to OpenAI. Aperture offers
a blueprint for
[managing quotas](https://docs.fluxninja.com/reference/blueprints/quota-scheduling/base),
consisting of two main components:

- Rate limiter: OpenAI employs a token bucket algorithm to impose rate limits,
  and that is directly compatible with Aperture's rate limiter. For example, in
  the tokens per minute policy for `gpt-4`, we can allocate a burst capacity of
  `40000 tokens`, and a refill rate of `40000 tokens per minute`. The bucket
  begins to refill the moment the tokens are withdrawn, aligning with OpenAI's
  rate-limiting mechanism. This ensures our outbound request and token rate
  remains synchronized with OpenAI's enforced limits.
- Scheduler: Aperture has a
  [weighted fair queuing](https://docs.fluxninja.com/concepts/scheduler/)
  scheduler that prioritizes the requests based on multiple factors such as the
  number of tokens, priority levels and workload labels.

## Pre-Requisites

Before you begin with this guide, verify the prerequisites are fulfilled.

- Aperture is installed and running. If not, follow the
  [get started guide](../../get-started/get-started.md).
- `aperturectl` is installed and configured. If not, follow the
  [aperturectl installation guide](../../get-started/installation/aperture-cli/aperture-cli.md).

## Configuration

Before creating a policy, a control point needs to be defined. Control Point
specifies where the policy should apply the decisions. There are multiple ways
to achieve this; for the scope of this guide, a JavaScript SDK is used to create
a control point—check out the
[Control Point](https://docs.fluxninja.com/concepts/control-point) Concept &
[Integrations](https://docs.fluxninja.com/integrations/) section for more
details.

### Control Point with JavaScript SDK

:::info

💡
[Check out other language SDKs supported by Aperture](/integrations/integrations.md)

:::

The Aperture JavaScript SDK allows you to set a control point manually. How an
Aperture Client instance is created is not discussed in this guide; detailed
information around SDK integration can be found in
[Manually setting feature control points using JavaScript SDK](/integrations/sdk/javascript/manually-setting-feature-control-points-using-javascript-sdk).

The code below provides a general idea of control point creation and setting
labels.

<details>
<summary>Integration with Aperture TypeScript SDK</summary>
<p>

<!-- markdownlint-disable MD010 -->

Import and setup Aperture Client:

```typescript
import { ApertureClient, FlowStatusEnum } from "@fluxninja/aperture-js";

apertureClient = new ApertureClient({
  address: "localhost:8080",
  channelCredentials: grpc.credentials.createSsl(),
});
```

Wrap the OpenAI API call with Aperture Client's `StartFlow` and `End` methods:

```typescript
const PRIORITIES: Record<string, number> = {
    paid_user: 10000,
    trial_user: 1000,
    free_user: 100,
}

let flow: Flow | undefined = undefined

if (this.apertureClient) {
    const charCount =
        this.systemMessage.length +
        message.length +
        String("system" + "user").length
    const labels: Record<string, string> = {
        api_key: CryptoES.SHA256(api.apiKey).toString(),
        estimated_tokens: (
            Math.ceil(charCount / 4) + responseTokens
        ).toString(),
        model_variant: modelVariant,
        priority: String(
            PRIORITIES[userType],
        ),
    }

    flow = await this.apertureClient.StartFlow("openai", {
        labels: labels,
        grpcCallOptions: {
            deadline: Date.now() + 1200000,
        },
    })
}

// As we use Aperture as a queue, send the message regardless of whether it was accepted or rejected
try {
    const { data: chatCompletion, response: raw } = await api.chat.completions
        .create({
            model: modelVariant,
            temperature: temperature,
            top_p: topP,
            max_tokens: responseTokens,
            messages: messages,
        })
        .withResponse()
        .catch(err => {
            logger.error(`openai chat error: ${JSON.stringify(err)}`)
            throw err
        })
    )
    return chatCompletion.choices[0]?.message?.content ?? ""
} catch (e) {
    flow?.SetStatus(FlowStatusEnum.Error)
    throw e // throw the error to be caught by the chat function
} finally {
    flow?.End()
}
```

<!-- markdownlint-enable MD010 -->

</p>
</details>

Let's understand the code snippet above; we are creating a control point named
`openai` and setting the labels, which will be used by the policy to identify
and schedule the request. Before calling OpenAI, we rely on Aperture Agent to
gate the request using the `StartFlow` method. To provide more context to
Aperture, we also attach the following labels to each request:

- `model_variant`: This specifies the model variant being used (`gpt-4`,
  `gpt-3.5-turbo`, or `gpt-3.5-turbo-16k`). Requests and tokens per minute rate
  limit policies are set individually for each model variant.
- `api_key` - This is a cryptographic hash of the OpenAI API key, and rate
  limits are enforced on a per-key basis.
- `estimated_tokens`: As the tokens per minute quota limit is enforced based on
  the
  [estimated tokens for the completion request](https://platform.openai.com/docs/guides/rate-limits/reduce-the-max_tokens-to-match-the-size-of-your-completions),
  we need to provide this number for each request to Aperture for metering.
  Following OpenAI's
  [guidance](https://help.openai.com/en/articles/4936856-what-are-tokens-and-how-to-count-them),
  we calculate `estimated_tokens` as `(character_count / 4) + max_tokens`. Note
  that OpenAI's rate limiter doesn't tokenize the request using the model's
  specific tokenizer but relies on a character count-based heuristic.
- `priority`: Requests are ranked according to a priority number provided in
  this label. For example, requests from `paid_user` can be given precedence
  over those from `trial_user` and `free_user` in example code.

### Policies

To generate a policy using quota scheduler blueprint, `values` files should be
generated first, specific to the policy. The values file can be generated using
the following command:

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints values --name=quota-scheduling/base --output-file=gpt-4-tpm-values.yaml</CodeBlock>
```

The values file needs to be adjusted to match the application requirements -

- `policy_name`: Name of the policy — This value should be unique and required.
- `bucket_capacity`: This value defines burst capacity. E.g. in the case of
  `gpt-4` tokens per minute policy, the bucket will have a capacity of
  `40000 tokens`.
- `fill_amount`: After the tokens are consumed, the bucket will be filled with
  this amount. E.g. in the case of `gpt-4` tokens per minute policy, the bucket
  will fill at `40000 tokens per minute`.
- `rate_limiter`:
  - `interval`: Interval at which the rate limiter will be filled. When to reset
    the bucket.
  - `label_key`: Label key to match the request against. This label key in this
    case is the OpenAI API key (`api_key`) which helps determine the quota for
    the request.

The scheduler helps prioritize the requests based on the labels and priority
defined. In this case, we are using the `priority` label, which is being passed
by Aperture SDK in code, containing the priority of the request.

- `scheduler`:
  - `priority_label_key`: Priority label key to match the request against. In
    this case, it is `priority`.
  - `tokens_label_key`: In the case of tokens per minute policy, each request
    has a `estimated_tokens` label value, which can be used to prioritize the
    request based on the number of tokens. In this case, it is
    `estimated_tokens`.
  - `workloads`:
    - `name`: To match the label value against the name of workloads. In this
      case, it is `paid_user`, `trial_user`, `free_user`.
    - `label_matcher`:
      - `match_labels`: Labels to match the request against. In this case, it is
        `product_reason`.

Selector parameters allow filtering of the requests to ensure where the policy
will act on.

- `selectors`:
  - `control_point`: Control point name to match the request against. In this
    case, it will be `openai`.
  - `agent_group`: Agent group name to match the request against. It is
    optional.
  - `label_matcher`:
    - `match_labels`: Labels to match the request against. It is optional.

Below are examples of values file adjusted to match the SDK code snippet &
control point labels.

<details>
<summary>Client-side quota management policies for gpt-4</summary>
<p>

```mdx-code-block
<Tabs>
<TabItem value="Tokens Per Minute (gpt-4)">
```

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
  policy_name: gpt-4-tpm
  quota_scheduler:
    # Bucket capacity.
    # Type: float64
    # Required: True
    bucket_capacity: 40000
    # Fill amount.
    # Type: float64
    # Required: True
    fill_amount: 40000
    # Rate Limiter Parameters
    # Type: aperture.spec.v1.RateLimiterParameters
    # Required: True
    rate_limiter:
      interval: 60s
      label_key: api_key
    scheduler:
      priority_label_key: priority
      tokens_label_key: estimated_tokens
    # Flow selectors to match requests against
    # Type: []aperture.spec.v1.Selector
    # Required: True
    selectors:
      - control_point: openai
        agent_group: default
        label_matcher:
          match_labels:
            model_variant: gpt-4
```

```mdx-code-block
</TabItem>
<TabItem value="Requests Per Minute (gpt-4)">
```

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
    bucket_capacity: 200
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
    scheduler:
      priority_label_key: priority
    # Flow selectors to match requests against
    # Type: []aperture.spec.v1.Selector
    # Required: True
    selectors:
      - control_point: openai
        agent_group: default
        label_matcher:
          match_labels:
            model_variant: gpt-4
```

```mdx-code-block
</TabItem>
</Tabs>
```

</p>
</details>

#### Apply Policy

```mdx-code-block
<Tabs>
<TabItem value="aperturectl (Aperture Cloud)" label="aperturectl (Aperture Cloud)">
<CodeBlock language="bash">
aperturectl cloud blueprints apply --values-file=gpt-4-tpm.yaml
</CodeBlock>
</TabItem>
<TabItem value="aperturectl (self-hosted controller)" label="aperturectl (self-hosted controller)">
```

Pass the `--kube` flag with `aperturectl` to directly apply the generated policy
on a Kubernetes cluster in the namespace where the Aperture Controller is
installed.

```mdx-code-block
<CodeBlock language="bash">
aperturectl blueprints generate --values-file=gpt-4-tpm.yaml --output-dir=policy-gen
aperturectl apply policy --file=policy-gen/policies/gpt-4-tpm.yaml --kube
</CodeBlock>
```

```mdx-code-block
</TabItem>
<TabItem value="kubectl (self-hosted controller)" label="kubectl (self-hosted controller)">
```

Apply the generated policy YAML (Kubernetes Custom Resource) with `kubectl`.

```bash
aperturectl blueprints generate --values-file=gpt-4-tpm.yaml --output-dir=policy-gen
kubectl apply -f policy-gen/policies/gpt-4-tpm-cr.yaml -n aperture-controller
```

```mdx-code-block
</TabItem>
</Tabs>
```

## Policy in Action

Once the policy is activated, it will begin to ensure that API requests conform
to OpenAI's rate limits, prioritizing requests based on the workload types
defined in the policy. These workloads are matched with the labels that the SDK
passes to Aperture, where paid users are prioritized over trial users, and trial
users over free users, thereby establishing a baseline experience for each tier.

Should rate limits be exhausted, Aperture will schedule the requests, placing
them in a queue until either the bucket is reset or a token becomes available.
Requests will remain in the queue up to the specified timeout value, provided in
the `StartFlow` function. There is no minimum waiting time for requests; as soon
as tokens are available, requests will be forwarded to OpenAI. A request can be
transmitted to OpenAI as soon as it reaches the application, or it can wait
until the timeout.

### Monitoring the Policy and OpenAI Performance

Aperture Cloud provides comprehensive observability of the policy and OpenAI
performance, providing a granular view of each workload, such as paid, trial,
and free.

The image below shows the incoming token rate and the accepted token rate for
the `gpt-4` tokens-per-minute policy. We can observe that the incoming token
rate is spiky, while the accepted token rate remains smooth and hovers around
`666 tokens per second`. This roughly translates to `40,000 tokens per minute`.
Essentially, Aperture is smoothing out the fluctuating incoming token rate to
align it with OpenAI's rate limits.

![Token Rate in Light Mode](./assets/openai/token-rate-light.png#gh-light-mode-only)

![Token Rate in Dark Mode](./assets/openai/token-rate-dark.png#gh-dark-mode-only)
_Incoming and Accepted Token Rate for gpt-4_

The below image shows request prioritization metrics from the Aperture Cloud
console during the same peak load period:

![Prioritization Metrics in Light Mode](./assets/openai/priorities-light.png#gh-light-mode-only)

![Prioritization Metrics in Dark Mode](./assets/openai/priorities-dark.png#gh-dark-mode-only)
_Prioritization Metrics for gpt-4_

In the upper left panel of the metrics, noticeable peaks indicate that some
requests got queued for several minutes in Aperture. We can verify that the
trial and free-tier users tend to experience longer queue times compared to
their paid counterparts and chat requests.

Queue wait times can fluctuate based on the volume of simultaneous requests in
each workload. For example, wait times are significantly longer during peak
hours as compared to off-peak hours. Aperture provides scheduler preemption
metrics to offer further insight into the efficacy of prioritization. As
observed in the lower panels, these metrics measure the relative impact of
prioritization for each workload by comparing how many tokens a request gets
preempted or delayed in the queue compared to a purely First-In, First-Out
(FIFO) ordering.

In addition to effectively managing the OpenAI quotas, Aperture provides
insights into OpenAI API performance and errors. The graphs below show the
overall response times for various OpenAI models we use. We observe that the
`gpt-4` family of models is significantly slower compared to the `gpt-3.5-turbo`
family of models.

![Flow Analytics](./assets/openai/flow-analytics-light.png#gh-light-mode-only)

![Flow Analytics](./assets/openai/flow-analytics-dark.png#gh-dark-mode-only)
_Performance Metrics for OpenAI Models_
