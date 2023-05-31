---
title: GraphQL Query Rate Limiting
keywords:
  - policies
  - ratelimit
  - graphql
sidebar_position: 2
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

This tutorial demonstrates how to use the [_Classifier_][rego-rules] to
implement
[rate-limiting](/reference/policies/bundled-blueprints/policies/rate-limiting.md)
for a GraphQL query.

## Policy

This tutorial will demonstrate how to implement a policy that uses a
[_Classifier_][classifier] to extract the `userID` claim from a JWT token in the
request's Authorization header and then rate limit unique users based on that
`user_id` [_Flow Label_][flow-label].

:::tip

You can write classification rules on
[HTTP requests](/concepts/flow-control/resources/classifier.md#live-previewing-requests)
and define scheduler priorities on
[Flow Labels](/concepts/flow-control/flow-label.md#live-previewing-flow-labels)
by live previewing them first using introspection APIs.

:::

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/graphql-rate-limiting/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/graphql-rate-limiting/graphql-rate-limiting-jwt.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/graphql-rate-limiting/graphql-rate-limiting-jwt.mmd.svg)
for this policy.

:::

For example, if the mutation query is as follows

```graphql
mutation createTodo {
  createTodo(input: { text: "todo" }) {
    user {
      id
    }
    text
    done
  }
}
```

Without diving deep into how Rego works, the source section mentioned in this
tutorial does the following:

1. Parse the query
2. Check if the mutation query is `createTodo`
3. Verify the JWT token with a secret key `secret` (only for demonstration
   purposes)
4. Decode the JWT token and extract the `userID` from the claims
5. Assign the value of `userID` to the exported variable `userID` in Rego source

From there on, the Classifier rule assigns the value of the exported variable
`userID` in Rego source to `user_id` flow label, effectively creating a label
`user_id:1`. This label is used by the
[`RateLimiter`](/concepts/flow-control/components/rate-limiter.md) component in
the policy to limit the `createTodo` mutation query to `10 requests/second` for
each `userID`.

### Playground

In this example, the traffic generator is configured to generate
`50 requests/second` for 2-minutes. When loading the above policy in the
playground, you can observe that it accepts no more than `2 requests/second` at
any given time, and rejects the rest of the requests.

<Zoom>

![GraphQL Status Rate Limiting](./assets/graphql-rate-limiting/graphql-rate-limiting-counter.png)

</Zoom>

[rego-rules]: /concepts/flow-control/resources/classifier.md#rego
[flow-label]: /concepts/flow-control/flow-label.md
[classifier]: /concepts/flow-control/resources/classifier.md