---
title: Agent API Keys
keywords:
  - API-keys
  - aperture-cloud
sidebar_position: 2
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
```

Aperture Cloud uses Agent API keys to authenticate requests coming from
[SDK integrations](/sdk/sdk.md), [self-hosted Agents][Agents] and [self-hosted
Controllers][Controllers]. You can create API keys for your project in the
Aperture Cloud UI.

## Pre-requisites

You have [signed up][sign-up] on Aperture Cloud and created a project.

## Create API Keys

1. In the Aperture Cloud UI, navigate to your project. _API keys are
   project-specific. You need to create a new API key for each project._
2. Now, from the left sidebar, click **Aperture**.
3. Click **Agent API Keys** tab.
4. Click **Create Agent API key**.
5. Copy the API key and save it in a secure location. This key will be used
   during [SDK integrations](/sdk/sdk.md) or [self-hosted Agents][Agents] .

![API Keys](./assets/api-keys.gif "Creating API Keys for sudhanshu-demo-docs project")

[sign-up]: /get-started/sign-up.md
[Agents]: /self-hosting/agent/agent.md
[Controllers]: /self-hosting/controller/controller.md