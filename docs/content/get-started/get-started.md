---
title: Get Started
keywords:
  - setup
  - getting started
sidebar_position: 2
sidebar_label: Get Started
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
```

In the next 15 minutes, you should have your application integrated with
Aperture and ready to enforce `rate limiting as a service`. Aperture offers
multiple integration methods, allowing you to choose the one that best suits
your application.

<Tabs>

<TabItem value="Aperture Serverless">

How to Integrate with Aperture Serverless: A Quick Guide for Developers

1. **Introduction**: This mode is tailored for developers who prefer to directly
   use our SDKs without delving into infrastructure components.
2. **Set Up the Environment**:
   `Aperture Controller and Aperture agent are both hosted within the FluxNinja SaaS platform.`
3. **Integrate the SDK**: Add the necessary few lines of code to your desired
   feature. Link to available SDKs are [here](../sdk/sdk.md).
4. **Configure the Organization Settings**: Point your feature to the created
   organization within FluxNinja.
5. **Map to Cloud Agent**: Aperture Cloud uses API keys to authenticate requests
   coming from SDK integrations. You can create API keys for your project in the
   Aperture Cloud UI.

</TabItem>

<TabItem value="Aperture SaaS">

How to Integrate with Aperture SaaS: A Quick Guide for Developers with Access to
Infrastructure, Platform and/or Infrastructure Teams

1. **Introduction**: This mode is crafted for teams running applications across
   different infrastructure platforms. It's a perfect fit for those wanting to
   avoid hosting their own Prometheus and etcd instances while aiming for
   minimal performance impact from the Aperture agent. Additionally,
   `teams should be comfortable with metrics being sent to FluxNinja SaaS.`
2. **Set Up the Environment**:
   `Aperture Controller is hosted within the FluxNinja SaaS platform and Aperture Agent is within customer VPC.`
3. **Integrate with Aperture**:
   - SDKs
   - ServiceMesh
   - Gateways
4. **Create Your Account and Set Up Your Organization**: Get started with
   Aperture by creating an account
5. **Map to Aperture SaaS Controller**: Link to the Quickstart section

</TabItem>

<TabItem value="Aperture Open Source">

How to Integrate with Aperture Open Source: A Quick Guide for Developers with
Access to Infrastructure, Platform, and/or Infrastructure Teams

1. **Introduction**: This mode is tailored for teams deploying applications on a
   variety of infrastructure platforms. It's well-suited for those who are
   `comfortable hosting their own Prometheus and etcd instances` and aim for
   minimal performance impact from the Aperture agent. Also, teams who are not
   comfortable with metrics being sent to FluxNinja SaaS.
2. **Set Up the Environment**:
   `Aperture Controller and Aperture Agent is within customer VPC.`
3. **Integrate with Aperture**:
   - SDKs
   - ServiceMesh
   - Gateways
4. **Self-hosted setup**: Link to Quickstart guide

</TabItem>

</Tabs>
