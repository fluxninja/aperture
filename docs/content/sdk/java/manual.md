---
title: Define Control Points
sidebar_position: 1
slug: define-feature-control-points-using-java-sdk
keywords:
  - java
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
description:
  Learn how to use Aperture's Java SDK to manually set feature control points
  and improve the reliability and stability of your web-scale applications. This
  guide covers best practices and provides examples for implementation.
---

```mdx-code-block
import CodeSnippet from '../../codeSnippet.js'
```

[Aperture Java SDK core library][SDK-Library] can be used to define feature
control points within a Java service.

The next step is to create an Aperture Client instance, for which the address of
the organization created in Aperture Cloud and API key are needed. You can
locate both these details by clicking on the Aperture tab in the sidebar menu of
Aperture Cloud.

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

```java
    String agentAddress = "ORGANIZATION.app.fluxninja.com:443";
    String apiKey = "API_KEY";
```

<CodeSnippet lang="java" snippetName="StandaloneExampleSDKInit" />

The created instance can then be used to start a flow:

<CodeSnippet lang="java" snippetName="StandaloneExampleFlow" />

The above code snippet is making `startFlow` calls to Aperture. For this call,
it is important to specify the control point (`featureName` in the example) and
business labels that will be aligned with the policy created in Aperture Cloud.
For each flow that is started, a `shouldRun` decision is made, determining
whether to allow the request into the system or to rate limit it. In this
example, we only see response returns, but in a production environment, actual
business logic can be executed when a request is allowed. It is important to
make the `end` call made after processing each request, to send telemetry data
that would provide granular visibility for each flow.

For more context on using Aperture Java SDK to set feature control points, refer
to the [example app][example] available in the repository.

[example]:
  https://github.com/fluxninja/aperture-java/blob/releases/aperture-java/v2.1.0/examples/standalone-example/src/main/java/com/fluxninja/example/App.java
[api-keys]: /reference/cloud-ui/api-keys.md
[SDK-Library]:
  https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-core
