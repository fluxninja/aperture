---
title: Manually setting feature control points
keywords:
  - Manually setting feature control points
sidebar_position: 1
sidebar_label: Manually setting feature control points
---

```mdx-code-block
import { Cards } from '@site/src/components/Cards';
```

Control points are used to define where you want to act in code or at service
level. It's important to understand what control points are because you will be
using them many times in your code.

<!-- vale off -->

## What is a Manual Feature Control point?

<!-- vale on -->

A manual feature control point is essentially a specific point in the codebase
where the execution flow can be controlled manually using feature flags. Feature
flags, also known as feature toggles, are a programming technique that allows
developers to enable or disable features of their software even after the code
has been deployed to production. This can be extremely useful for testing new
features, performing Blue Green testing, or quickly disabling a feature in
response to an issue or error.

<!-- vale off -->

## How to create a Manual Feature Control point?

<!-- vale on -->

Let's create a feature control point manually in java code. To begin with, you
need to configure the Aperture Java SDK for your application. You can configure
Aperture SDK as follows:

```java
    ApertureSDK apertureSDK;

    apertureSDK = ApertureSDK.builder()
            .setAddress("ORGANIZATION.app.fluxninja.com:443")
            .setAgentAPIKey("AGENT_API_KEY")
            .build();
```

Once you have configured Aperture SDK, you can create a feature control point
wherever you want in your code. For example, there is a function called
`handleSuperAPI` which is called when a user hits a specific API. Before
executing the business logic, you want to create a feature control point so that
you can control the execution flow of the API and can reject the request based
on the policy defined in Aperture. You will see in the upcoming section how to
define policy in Aperture. For now, some labels have been added in the code
snippet below. These labels will be used while defining a policy.

Let's create a feature control point in the following code snippet.

```java
   private String handleSuperAPI(spark.Request req, spark.Response res) {
        Map<String, String> labels = new HashMap<>();

        // do some business logic to collect labels
        labels.put("user", "kenobi");

        Map<String, String> allHeaders = new HashMap<>();
        for (String headerName : req.headers()) {
            allHeaders.put(headerName, req.headers(headerName));
        }
        allHeaders.putAll(labels);

        TrafficFlowRequestBuilder trafficFlowRequestBuilder = TrafficFlowRequest.newBuilder();

        trafficFlowRequestBuilder
                .setControlPoint("awesomeFeature")
                .setHttpMethod(req.requestMethod())
                .setHttpHost(req.host())
                .setHttpProtocol(req.protocol())
                .setHttpPath(req.pathInfo())
                .setHttpScheme(req.scheme())
                .setHttpSize(req.contentLength())
                .setHttpHeaders(allHeaders)
                .setSource(req.ip(), req.port(), "TCP")
                .setDestination(req.raw().getLocalAddr(), req.raw().getLocalPort(), "TCP")
                .setRampMode(false)
                .setFlowTimeout(Duration.ofMillis(1000));

        TrafficFlowRequest apertureRequest = trafficFlowRequestBuilder.build();

        TrafficFlow flow = this.apertureSDK.startTrafficFlow(apertureRequest);

        // See whether flow was accepted by Aperture Agent.
        try {
            if (flow.shouldRun()) {
                // Aperture accepted the flow, now execute the business logic.
                data = this.executeBusinessLogic(spark.Request);
                res.status(200);
            } else {
                // Flow has been rejected by Aperture Agent.
                res.status(flow.getRejectionHttpStatusCode());
            }
        } catch (Exception e) {
            // Flow Status captures whether the feature captured by the Flow was
            // successful or resulted in an error. When not explicitly set,
            // the default value is FlowStatus.OK .
            flow.setStatus(FlowStatus.Error);
            logger.error("Error in flow execution", e);
        } finally {
            flow.end();
        }
        return data;
    }
```

This is how you can create a manual feature control point in your code. The
complete code snippet is available
[here](https://github.com/fluxninja/aperture-java/tree/main/examples/standalone-traffic-flow-example).

:::info

Aperture SDKs are available for multiple languages, you choose the one that fits
your needs. [See all SDKs](/integrations/sdk/sdk.md).

:::

<!-- vale off -->

## What's next?

<!-- vale on -->

Once the feature control point is set in code, head over to
[Set up CLI (aperturectl)](/get-started/setup-cli/setup-cli.md)
