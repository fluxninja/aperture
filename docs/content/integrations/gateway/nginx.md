---
title: Nginx
keywords:
  - install
  - setup
  - nginx
  - nginx-gateway
sidebar_position: 1
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../apertureVersion.js';
import Zoom from 'react-medium-image-zoom';
```

Integrating Aperture with Nginx using Lua modules.

![Nginx](./assets/nginx-dark.svg#gh-dark-mode-only)

![Nginx](./assets/nginx-light.svg#gh-light-mode-only)

## Introduction

:::info

This integration is currently only supported with the
[self-hosted](/get-started/self-hosting/agent/agent.md) version of Aperture
Agent.

:::

Lua's modules are scripts that can be executed within Nginx to extend its
functionality. The Aperture Lua module can be downloaded from the GitHub <a
href={`https://github.com/fluxninja/aperture/releases/tag/${apertureVersion}`}>
Aperture Release Page</a>.

## Pre-requisites

Before proceeding, ensure that you have the following installed:

:::info

Skip these steps if the Nginx server is running on a Container.

:::

1. Nginx server
2. [`lua-nginx-module`](https://github.com/openresty/lua-nginx-module) enabled
   for Nginx. If not, follow the
   [installation steps](https://github.com/openresty/lua-nginx-module#installation).
3. [`LuaRocks`](https://luarocks.org/), which is a package manager for Lua
   modules. Follow the
   [installation steps](https://github.com/luarocks/luarocks/wiki/Download#installing).

## Installation

To install the Aperture Lua module, follow these steps:

:::info

Refer to [Example Dockerfile](#example-dockerfile) to get the steps for
installing the Aperture Lua module for Nginx server running on Container.

:::

1. Install the
   [`opentelemetry-lua`](https://github.com/fluxninja/opentelemetry-lua) SDK by
   running the following commands:

   ```bash
   git clone https://github.com/fluxninja/opentelemetry-lua.git
   cd opentelemetry-lua
   luarocks make
   ```

2. Download and extract the Aperture Lua module by executing the following
   commands:

   ```mdx-code-block
   <CodeBlock language="bash">
   wget "https://github.com/fluxninja/aperture/releases/download/{apertureVersion}/aperture-lua.tar.gz" && tar -xzvf aperture-lua.tar.gz
   </CodeBlock>
   ```

3. Install the module by running the following command:

   ```bash
   cd aperture-lua && luarocks make aperture-nginx-plugin-0.1.0-1.rockspec
   ```

<!-- vale off -->

## Example Dockerfile {#example-dockerfile}

<!-- vale on -->

Use the following Dockerfile to install the Aperture Lua module with Nginx. This
example uses
[`fabiocicerchia/nginx-lua`](https://hub.docker.com/r/fabiocicerchia/nginx-lua/)
as the base image because it already has the
[`lua-nginx-module`](https://github.com/openresty/lua-nginx-module)
pre-configured with Nginx.

```mdx-code-block
<CodeBlock language="Dockerfile">{`FROM fabiocicerchia/nginx-lua:1.23.3-debian-compat\n
RUN apt update && apt-get install -y build-essential git\n
RUN git clone https://github.com/fluxninja/opentelemetry-lua.git && cd opentelemetry-lua && luarocks make\n
RUN curl --fail --location --remote-name "https://github.com/fluxninja/aperture/releases/download/${apertureVersion}/aperture-lua.tar.gz"\n
RUN tar -xzvf aperture-lua.tar.gz && luarocks make aperture-nginx-plugin-0.1.0-1.rockspec\n
COPY nginx_config.conf /etc/nginx/nginx.conf\n
ENTRYPOINT [ "nginx", "-g", "daemon off;" ]`}</CodeBlock>
```

## Configure Nginx

Follow these steps to configure Nginx to use the installed Aperture Lua module:

1. To connect to the Aperture Agent, you need to create an environment variable
   called APERTURE_AGENT_ENDPOINT. The value of this variable should be set
   equal to the endpoint of the Aperture Agent. If you are using a bash shell,
   you can create this variable by running the following command:

   ```bash
   echo 'export APERTURE_AGENT_ENDPOINT="http://aperture-agent.aperture-agent.svc.cluster.local"' >> ~/.profile
   ```

   Replace the endpoint value with the actual endpoint value of the Aperture
   Agent if you're on a different one.

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

3. Add the `init_by_lua_block` section under the `http` block of the Nginx
   configuration to initialize the Aperture Lua module:

   ```bash
   http {
     ...
     init_by_lua_block {
       access = require "aperture-plugin.access"
       log = require "aperture-plugin.log"
       headers = require "aperture-plugin.headers"
     }
     ...
   }
   ```

4. Add the `access_by_lua_block` section under the `http` block of the Nginx
   configuration to execute the Aperture check for all servers and locations
   before the request is forwarded to upstream:

   ```bash
   http {
     ...
     access_by_lua_block {
       local authorized_status = access(ngx.var.control_point)

       if authorized_status ~= ngx.HTTP_OK then
         return ngx.exit(authorized_status)
       end
     }
     ...
   }
   ```

5. Add the `header_filter_by_lua_block` section under the `http` block of the
   Nginx configuration to add the headers received from Aperture check to the
   response being returned to the client:

   ```bash
   http {
     ...
     header_filter_by_lua_block {
       headers()
     }
     ...
   }
   ```

6. Add the `log_by_lua_block` section under the `http` block of the Nginx
   configuration to forward the OpenTelemetry logs to Aperture for all servers
   and locations after the response is received from upstream:

   ```bash
   http {
     ...
     log_by_lua_block {
       log()
     }
     ...
   }
   ```

7. Aperture needs `control_point` variable for referring the service in Aperture
   Policy, which needs to be set from Nginx `location` block:

   ```bash
   http {
     ...
     server {
       location /service1 {
         set $control_point "service1-demo-app";
         proxy_pass http://service1-demo-app.demoapp.svc.cluster.local:80/request;
       }
     }
     ...
   }
   ```

8. Below is how a complete Nginx configuration would look like:

   ```bash
   worker_processes auto;
   pid /run/nginx.pid;

   events {
     worker_connections 4096;
   }

   http {
     default_type application/octet-stream;
     resolver 10.96.0.10;

     sendfile on;
     keepalive_timeout 65;

     init_by_lua_block {
       access = require "aperture-plugin.access"
       log = require "aperture-plugin.log"
       headers = require "aperture-plugin.headers"
     }

     access_by_lua_block {
       local authorized_status = access(ngx.var.control_point)

       if authorized_status ~= ngx.HTTP_OK then
         return ngx.exit(authorized_status)
       end
     }

     log_by_lua_block {
       log()
     }

     header_filter_by_lua_block {
       headers()
     }

     server {
       listen 80;
       proxy_http_version 1.1;

       location /service1 {
         set $control_point "service1-demo-app";
         proxy_pass http://service1-demo-app.demoapp.svc.cluster.local:80/request;
       }

       location /service2 {
         set $control_point "service2-demo-app";
         proxy_pass http://service2-demo-app.demoapp.svc.cluster.local:80/request;
       }

       location /service3 {
         set $control_point "service3-demo-app";
         proxy_pass http://service3-demo-app.demoapp.svc.cluster.local:80/request;
       }
     }
   }
   ```
