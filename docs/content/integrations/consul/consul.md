---
title: Consul
keywords:
  - install
  - setup
  - service mesh
  - consul
  - service defaults
sidebar_position: 1
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../apertureVersion.js';
import Zoom from 'react-medium-image-zoom';
```

![Istio](./assets/consul-light.svg#gh-light-mode-only)

![Istio](./assets/consul-dark.svg#gh-dark-mode-only)

## Supported Versions

Aperture supports the following version of Istio:

| Platform | Extent of Support |
| -------- | ----------------- |
| Consul   | 1.17 and above    |

## Service Defaults {#service-defaults}

**Note**: In all the below patches, it is presumed that the Aperture Agent is
installed with `DaemonSet` mode and is installed in the `aperture-agent`
namespace, which makes the target address value
`aperture-agent.aperture-agent.svc.cluster.local`. If you are running the
Aperture Agent in Sidecar mode, use `localhost` as the target address.

## Prerequisites

You can do the installation using the `aperturectl` CLI tool or using `Helm`.
Install the tool of your choice using the following links:

1. [aperturectl](/get-started/installation/aperture-cli/aperture-cli.md)

   :::info Refer

   [aperturectl install istioconfig](/reference/aperturectl/install/istioconfig/istioconfig.md)
   to see all the available command line arguments.

   :::

2. [Helm](https://helm.sh/docs/intro/install/)

   1. Once the Helm CLI is installed, add the
      [Aperture istioconfig Helm Repository](https://artifacthub.io/packages/helm/aperture/istioconfig)
      in your environment for installation:

      ```bash
      helm repo add aperture https://fluxninja.github.io/aperture/
      helm repo update
      ```

## Installation

## Verifying the Installation

## Uninstall
