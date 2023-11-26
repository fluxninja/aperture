---
title: Manually setting feature control points
sidebar_position: 1
slug: manually-setting-feature-control-points-using-java-sdk
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

[Aperture Java SDK core library][SDK-Library] can be used to manually set
feature control points within a Java service.

To do so, first create an instance of ApertureSDK:

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

```java
    String agentAddress = "ORGANIZATION.app.fluxninja.com:443";
    String agentAPIKey = "API_KEY";

    ApertureSDK apertureSDK;

    apertureSDK = ApertureSDK.builder()
            .setAddress(agentAddress)
            .setAPIKey(agentAPIKey)
            .setFlowTimeout(Duration.ofMillis(1000))
            .build();
```

The created instance can then be used to start a flow:

```java

    Map<String, String> labels = new HashMap<>();

    // business logic produces labels
    labels.put("key", "value");

    Boolean rampMode = false;

    FeatureFlowParameters params = FeatureFlowParameters.newBuilder("featureName")
        .setExplicitLabels(labels)
        .setRampMode(rampMode)
        .setFlowTimeout(Duration.ofMillis(1000))
        .build();

    Flow flow = apertureSDK.startFlow(params);
    if (flow.shouldRun()) {
        // do actual work
    } else {
        // handle flow rejection by Aperture Agent
        flow.setStatus(FlowStatus.Error);
    }
    flow.end();
```

For more context on using Aperture Java SDK to set feature control points, refer
to the [example app][example] available in the repository.

[example]:
  https://github.com/fluxninja/aperture-java/blob/releases/aperture-java/v2.1.0/examples/standalone-example/src/main/java/com/fluxninja/example/App.java
[api-keys]: /reference/cloud-ui/api-keys.md
[SDK-Library]:
  https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-core
