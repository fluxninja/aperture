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

<a
href={`https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-core`}>Aperture
Java SDK core library</a> can be used to manually set feature control points
within a Java service.

To do so, first create an instance of ApertureSDK:

```java
    String agentAddress = "ORGANIZATION.app.fluxninja.com:443";
    String agentAPIKey = "AGENT_API_KEY";

    ApertureSDK apertureSDK;

    apertureSDK = ApertureSDK.builder()
            .setAddress(agentAddress)
            .setAgentAPIKey(agentAPIKey)
            .setFlowTimeout(Duration.ofMillis(1000))
            .build();
```

The created instance can then be used to start a flow:

```java

    Map<String, String> labels = new HashMap<>();

    // business logic produces labels
    labels.put("key", "value");

    Boolean rampMode = false;

    Flow flow = apertureSDK.startFlow("featureName", labels, rampMode);
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
