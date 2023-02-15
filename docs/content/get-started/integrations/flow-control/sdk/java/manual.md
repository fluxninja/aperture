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
Java SDK core library</a> can be used to manually set feature Control Points
within a Java service.

To do so, first create an instance of ApertureSDK:

```java
    String agentHost = "localhost";
    int agentPort = 8089;

    ApertureSDK apertureSDK;
    try {
        apertureSDK = ApertureSDK.builder()
                .setHost(agentHost)
                .setPort(agentPort)
                .setDuration(Duration.ofMillis(1000))
                .build();
    } catch (ApertureSDKException e) {
        e.printStackTrace();
        return;
    }
```

The created instance can then be used to start a flow:

```java

    Map<String, String> labels = new HashMap<>();

    // business logic produces labels
    labels.put("key", "value");

    Flow flow = apertureSDK.startFlow("featureName", labels);
    if (flow.accepted()) {
        // do actual work
        flow.end(FlowStatus.OK);
    } else {
        // handle flow rejection by Aperture Agent
        flow.end(FlowStatus.Error);
    }
```

For more context on how to use Java ApertureSDK to set feature Control Points,
you can take a look at the [example app][example] available in our repository.

[example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/standalone-example
