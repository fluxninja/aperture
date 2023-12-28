---
title: Define Control Points
sidebar_position: 1
slug: define-feature-control-points-using-python-sdk
keywords:
  - python
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
```

[Aperture Python SDK][pythonsdk] can be used to define feature control points
within a Python service.

Run the command below to install the SDK:

```bash
pip install aperture-py
```

The next step is to create an Aperture Client instance, for which, the address
of the organization created in Aperture Cloud and API key are needed. You can
locate both these details by clicking on the Aperture tab in the sidebar menu of
Aperture Cloud.

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, please refer to the
[API Keys](/reference/cloud-ui/api-keys.md) section.

:::

```mdx-code-block
import CodeSnippet from '../../codeSnippet.js'

<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="clientConstructor"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

The created instance can then be used to start a flow:

```mdx-code-block
<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="manualFlow"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

The above code snippet is making `start_flow` calls to Aperture. For this call,
it is important to specify the control point (`AwesomeFeature` in the example)
and `FlowParams` that will be aligned with the policy created in Aperture Cloud.
For request prioritization use cases, it's important to set a higher gRPC
deadline. This parameter specifies the maximum duration a request can remain in
the queue. For each flow that is started, a `should_run` decision is made,
determining whether to allow the request into the system or to rate limit it. It
is important to make the `end` call made after processing each request, to send
telemetry data that would provide granular visibility for each flow.

You can also use the flow as a context manager:

```mdx-code-block
<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="contextManagerFlow"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

Additionally, you can decorate any function with aperture client. This will skip
running the function if the flow is rejected by Aperture Agent. This might be
helpful to handle specific routes in your service.

```mdx-code-block
<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="apertureDecorator"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

For more context on using the Aperture Python SDK to set feature control points,
refer to the [example app][example] available in the repository.

[pythonsdk]: https://pypi.org/project/aperture-py/
[example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-py/example
