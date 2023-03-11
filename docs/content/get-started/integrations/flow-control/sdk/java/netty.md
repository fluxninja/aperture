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

### Netty Handler

<a
href={`https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-netty`}>Aperture
Java SDK Netty package</a> contains a Handler that automatically creates traffic
flows for request in a given pipeline:

```java
import com.fluxninja.aperture.netty.ApertureServerHandler;

...

public class ServerInitializer extends ChannelInitializer<Channel> {

    ...

    @Override
    protected void initChannel(Channel ch) {
        try {
            sdk = ApertureSDK.builder().setHost(this.agentHost).setPort(this.agentPort).build();
        } catch (ApertureSDKException ex) {
            throw new RuntimeException(ex);
        }

        ChannelPipeline pipeline = ch.pipeline();
        pipeline.addLast(new HttpServerCodec());
        pipeline.addLast(new HttpObjectAggregator(Integer.MAX_VALUE));
        // ApertureServerHandler must be added before the response-generating HelloWorldHandler,
        //    but after the codec handler.
        pipeline.addLast(new ApertureServerHandler(sdk));
        pipeline.addLast(new HelloWorldHandler());
    }
}
```

For more context on how to use Aperture Netty Handler to set Control Points, you
can take a look at the [example app][netty-example] available in our repository.

[netty-example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/netty-example
[javaagent]:
  /get-started/integrations/flow-control/sdk/java/using-instrumentation-agent-to-automatically-set-control-points-using-java-sdk
