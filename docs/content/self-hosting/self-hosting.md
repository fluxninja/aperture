---
title: Self-Hosting Aperture
sidebar_position: 6
keywords:
  - self-hosted
---

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

The easiest way to get started with Aperture is to [integrate with SDKs][sdks]
and point them to [Aperture Cloud][aperture-cloud], where Aperture Cloud will
take care of the [Controller][cloud-controller] and [Aperture
Agent][cloud-agent].

However, if you want to have complete control over the infrastructure and data,
it is possible to self-host your own Aperture Controller and Agents.

:::note

[Aperture Cloud can integrate][extension-config] with Self-Hosted Controller And
Agent too, providing an easy way to manage policies and a holistic view of the
infrastructure, along with tools for OLAP analysis of traffic.

:::

<DocCardList />

[aperture-cloud]: /introduction.md
[cloud-controller]: /reference/fluxninja.md#cloud-controller
[cloud-agent]: /reference/fluxninja.md#cloud-agent
[extension-config]: /reference/fluxninja.md#configuration
[sdks]: /sdk/sdk.md
