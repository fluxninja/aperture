---
title: Define Control Points
slug: define-feature-control-points-using-javascript-sdk
sidebar_position: 1
keywords:
  - js
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
---

```mdx-code-block
import CodeSnippet from '../../codeSnippet.js'
```

[Aperture JavaScript SDK](https://www.npmjs.com/package/@fluxninja/aperture-js)
can be used to define feature control points within a JavaScript service.

Run the command below to install the SDK:

```bash
npm install @fluxninja/aperture-js
```

The next step is to create an Aperture Client instance, for which, the address
of the organization created in Aperture Cloud and API key are needed. You can
locate both these details by clicking on the Aperture tab in the sidebar menu of
Aperture Cloud.

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

<CodeSnippet lang="ts" snippetName="clientConstructor" />

The created instance can then be used to start a flow:

<CodeSnippet
    lang="ts"
    snippetName="handleRequestRateLimit"
 />

The above code snippet is making `startFlow` calls to Aperture. For this call,
it is important to specify the control point (`awesomeFeature` in the example)
and business labels that will be aligned with the policy created in Aperture
Cloud. For request prioritization use cases, it's important to set a higher gRPC
deadline. This parameter specifies the maximum duration a request can remain in
the queue. For each flow that is started, a `shouldRun` decision is made,
determining whether to allow the request into the system or to rate limit it. In
this example, we only see log returns, but in a production environment, actual
business logic can be executed when a request is allowed. It is important to
make the `end` call made after processing each request, to send telemetry data
that would provide granular visibility for each flow.

For more context on using the Aperture JavaScript SDK to set feature control
points, refer to the [example app][example] available in the repository.

[example]: https://github.com/fluxninja/aperture-js/tree/main/example
[api-keys]: /reference/cloud-ui/api-keys.md
