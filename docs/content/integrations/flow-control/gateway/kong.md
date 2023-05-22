---
title: Kong
keywords:
  - install
  - setup
  - kong
sidebar_position: 2
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../../apertureVersion.js';
```

Integrating Aperture with Kong using Custom plugins.

## Introduction

Custom plugins are Lua scripts that can be executed within Kong to extend its
functionality. The Aperture Custom plugin can be downloaded from the <a
href={`https://github.com/fluxninja/aperture/releases/tag/${apertureVersion}`}>Aperture
Release Page</a>.

## Installation

To install the Aperture Custom plugin, follow these steps:

:::info

Refer to [`Example Dockerfile`](#example-dockerfile) to get the steps for
installing the Aperture Custom plugin for Kong server running on Container.

:::

1. Install the
   [`opentelemetry-lua`](https://github.com/fluxninja/opentelemetry-lua) SDK by
   running the following commands:

   ```bash
   git clone https://github.com/fluxninja/opentelemetry-lua.git
   cd opentelemetry-lua
   luarocks make
   ```

2. Download and extract the Aperture Custom plugin by executing the following
   commands:

   ```mdx-code-block
   <CodeBlock language="bash">
   wget "https://github.com/fluxninja/aperture/releases/download/{apertureVersion}/aperture-lua.tar.gz" && tar -xzvf aperture-lua.tar.gz
   </CodeBlock>
   ```

3. Install the module by running the following command:

   ```bash
   cd aperture-lua && luarocks make aperture-kong-plugin-0.1.0-1.rockspec
   ```

<!-- vale off -->

## Example Dockerfile {#example-dockerfile}

<!-- vale on -->

Use the following Dockerfile to install the Aperture Custom plugin with Kong.

```mdx-code-block
<CodeBlock language="Dockerfile">{`FROM kong:3.1.1-ubuntu\n
WORKDIR /usr/kong/aperture\n
USER root\n
RUN apt update && apt-get install -y build-essential git curl\n
RUN git clone https://github.com/fluxninja/opentelemetry-lua.git && cd opentelemetry-lua && luarocks make\n
RUN curl --fail --location --remote-name "https://github.com/fluxninja/aperture/releases/download/${apertureVersion}/aperture-lua.tar.gz"\n
RUN tar -xzvf aperture-lua.tar.gz && luarocks make aperture-kong-plugin-0.1.0-1.rockspec\n
USER kong\n
COPY kong.conf .\n
COPY kong.yaml .\n
ENV KONG_DATABASE=off\n
ENV KONG_DECLARATIVE_CONFIG=kong.yaml\n
CMD [ "kong", "start", "-c", "kong.conf"]`}</CodeBlock>
```

## Configure Kong

Follow these steps to configure Kong to use the Aperture Custom plugin. Assuming
the plugin is already installed:

1. Create an environment variable `APERTURE_AGENT_ENDPOINT` with a value equal
   to the Aperture Agent endpoint. For example, use the following command in
   bash:

   ```bash
   echo 'export APERTURE_AGENT_ENDPOINT="http://aperture-agent.aperture-agent.svc.cluster.local"' >> ~/.profile
   ```

2. Optionally, create an environment variable `APERTURE_CHECK_TIMEOUT`, which
   would be considered as a timeout for execution of the Aperture check. The
   default value for it is `500m`, which is 500 milliseconds. For example, use
   the following command in bash:

   :::info

   The format for the `Timeout` parameter can be found at the following
   [link](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests).

   :::

   ```bash
   echo 'export APERTURE_CHECK_TIMEOUT="1S"' >> ~/.profile
   ```

3. Add the Aperture Custom plugin’s name to the plugins list in your Kong
   configuration (on each Kong node):

   ```yaml
   plugins = bundled,aperture-plugin
   ```

   You can also set this property from its environment variable equivalent:
   `KONG_PLUGINS`.

4. Enable the Aperture Custom plugin for all the services and routes by creating
   a YAML file like below:

   ```yaml
   _format_version: "3.0"
   _transform: true

   services:
     - name: service1-demo-app
       url: http://service1-demo-app.demoapp.svc.cluster.local:80/request
       retries: 3
       routes:
         - name: service1
           paths:
             - /service1
       plugins:
         - name: aperture-plugin
           config:
             control_point: service1-demo-app
     - name: service2-demo-app
       url: http://service2-demo-app.demoapp.svc.cluster.local:80/request
       retries: 3
       routes:
         - name: service2
           paths:
             - /service2
       plugins:
         - name: aperture-plugin
           config:
             control_point: service2-demo-app
     - name: service3-demo-app
       url: http://service3-demo-app.demoapp.svc.cluster.local:80/request
       retries: 3
       routes:
         - name: service3
           paths:
             - /service3
       plugins:
         - name: aperture-plugin
           config:
             control_point: service3-demo-app
   ```

   This adds the provided services and routes on Kong and enables the Aperture
   Plugin for them. This file can be passed by setting environment variable:
   `KONG_DECLARATIVE_CONFIG=ABOVE_FILENAME`.

   Alternatively, you can also enable the Aperture Custom plugin per service or
   route using the
   [Kong Gateway’s Admin API](https://docs.konghq.com/gateway/latest/admin-api/#plugin-object).
