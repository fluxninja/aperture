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
    String agentHost = "localhost"; // Aperture Agent Host URL
    int agentPort = 8089; // Aperture Agent Port

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

Once you have configured Aperture SDK, you can create a feature control point
wherever you want in your code. For example, there is a function called
`handleSuperAPI` which is called when a user hits a specific API. Before
executing the business logic, you want to create a feature control point so that
you can control the execution flow of the API and can reject the request based
on the policy defined in Aperture. You will see in the upcoming section how to
define policy in Aperture. For now, remember, we added some labels in the code
snippet below. These labels will be used while defining a policy.

Let's create a feature control point in the following code snippet.

```java
   private String handleSuperAPI(spark.Request req, spark.Response res) {
        Map<String, String> labels = new HashMap<>();
        Map<String, String> data = new HashMap<>();
​
        // Get the user_id from the request
        String userId = req.queryParams("user_id");
        // Set the user_id as a label
        labels.put("user_id", userId);
​
        // Get user type from the request
        String userType = req.queryParams("user_type");
        // Set the user_type as a label
        labels.put("user_type", userType);
​

        Flow flow = this.apertureSDK.startFlow('awesomeFeature', labels);
        FlowResult flowResult = flow.result();
​
        if (flowResult != FlowResult.Rejected) {
            // Aperture accepted the flow, now execute the business logic.
            data = this.executeBusinessLogic(spark.Request);
            res.status(200);
            flow.end(FlowStatus.OK);
        } else {
            // Flow has been rejected by Aperture
            res.status(flow.);
            flow.end(FlowStatus.Error);
        }
        return data;
    }
```

This is how you can create a manual feature control point in your code. The
complete code snippet is available
[here](https://github.com/fluxninja/aperture-java/tree/releases/aperture-java/v2.1.0/examples/standalone-example).

:::info

Aperture SDKs are available for multiple languages, you choose the one that fits
your needs. [See all SDKs](/integrations/flow-control/sdk/sdk.md).

:::

<!-- vale off -->

## What's next?

<!-- vale on -->

Once the feature control point is set in code, head over to
[install Aperture](/get-started/installation/installation.md)
