---
title: Netty
sidebar_position: 4
slug: netty
keywords:
  - java
  - sdk
  - control
  - points
  - middleware
  - netty
---

### Aperture Java Instrumentation Agent

All Netty pipelines can have an Aperture Handler automatically added into them
using [Aperture Instrumentation Agent][javaagent].

:::info Agent API Key

You can create an Agent API key for your project in the Aperture Cloud UI. For
more information, refer to
[Agent API Keys](/get-started/aperture-cloud/agent-api-keys.md).

:::

### Netty Handler

[Aperture Java SDK Netty package](https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-netty)
contains a Handler that automatically creates traffic flows for requests in a
given pipeline:

```java
import com.fluxninja.aperture.netty.ApertureServerHandler;

...

public class ServerInitializer extends ChannelInitializer<Channel> {

    ...

    @Override
    protected void initChannel(Channel ch) {
        String controlPointName = "someFeature";

        sdk = ApertureSDK.builder().setAddress(this.agentAddress).setAgentAPIKey(this.agentAPIKey).build();

        ChannelPipeline pipeline = ch.pipeline();
        pipeline.addLast(new HttpServerCodec());
        pipeline.addLast(new HttpObjectAggregator(Integer.MAX_VALUE));
        // ApertureServerHandler must be added before the response-generating HelloWorldHandler,
        //    but after the codec handler.
        pipeline.addLast(new ApertureServerHandler(sdk, controlPointName));
        pipeline.addLast(new HelloWorldHandler());
    }
}
```

You can instruct the handler to exclude specific paths from being monitored by
the Aperture SDK. For example, you might want to exclude endpoints used for
health checks. To achieve this, you can add the path(s) you want to ignore to
the `ignoredPaths` field of the SDK, as shown in the following code:

```java
ApertureSDK sdk = ApertureSDK.builder()
        .setAddress("ORGANIZATION.app.fluxninja.com:443")
        .setAgentAPIKey(agentAPIKey)
        ...
        .addIgnoredPaths("/healthz,/metrics")
        ...
        .build()
```

The paths should be specified as a comma-separated list. Note that the paths you
specify must match exactly. However, you can change this behavior to treat the
paths as regular expressions by setting the `ignoredPathMatchRegex` flag to
true, like so:

```java
  builder
        .setIgnoredPathMatchRegex(true)
```

For more context on using Aperture Netty Handler to set Control Points, refer to
the [example app][netty-example] available in the repository.

[netty-example]:
  https://github.com/fluxninja/aperture-java/blob/releases/aperture-java/v2.1.0/examples/netty-example/src/main/java/com/fluxninja/example/ServerInitializer.java
[javaagent]: /integrations/sdk/java/auto-instrumentation.md
