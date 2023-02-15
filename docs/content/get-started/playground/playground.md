---
title: Try Local Playground
keywords:
  - playground
  - policies
  - poc
  - ratelimit
  - concurrencylimit
sidebar_position: 5
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

Playground is a Kubernetes based environment for exploring the capabilities of
Aperture. Additionally, it is used as a development environment for Aperture.
Playground uses [Tilt](https://tilt.dev/) for orchestrating the deployments in
Kubernetes. Tilt watches for changes to local files and auto-deploys any
resources that change. This is very convenient for getting quick feedback during
development of Aperture.

Playground deploys resources to the Kubernetes cluster that `kubectl` on your
machine points at. For convenience, refer to [Prerequisites](#prerequisites-k8s)
for deploying a local Kubernetes cluster using
[Kind](https://kind.sigs.k8s.io/).

## How to Run

Assuming that you have already cloned the aperture repository and brought up a
[local Kubernetes cluster](#prerequisites-k8s), proceed to install the
[required tools](#tools). In order to bring up the Playground, run the following
commands:

```sh
$ git clone https://github.com/fluxninja/aperture.git
# change directory to playground
$ cd aperture/playground
# start a local kubernetes cluster
$ ctlptl apply -f ctlptl-kind-config.yaml
# start Tilt and run services defined in Tiltfile
$ tilt up
Tilt started on http://localhost:10350/
v0.30.2, built 2022-06-06

(space) to open the browser
(s) to stream logs (--stream=true)
(t) to open legacy terminal mode (--legacy=true)
(ctrl-c) to exit
```

Now, press Space to open the Tilt UI in your default browser.

:::note

Make sure nothing else is running on the [ports forwarded](#port-forwards) by
Tilt.

:::

The above command starts Aperture Controller and an Aperture Agent on each
worker node in the local Kubernetes cluster. Additionally, it starts a demo
application with an Istio and Envoy based service mesh configured to integrate
with Aperture. There is an instance of Grafana running on the cluster as well
for viewing metrics from experiments.

The Playground is preloaded with a
[Latency Gradient Policy](/tutorials/integrations/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
which protects the demo application against sudden surges in traffic load. You
can verify it using the following command:

```sh
$ kubectl get policy -n aperture-controller service1-demo-app
NAME                STATUS     AGE
service1-demo-app   uploaded   103s
```

The Playground comes with a demo application so that you can generate simulated
traffic and see the policy in action. The demo application can be found in
`demoapp` namespace. You can read more about the demo application
[here](https://github.com/fluxninja/aperture/tree/main/playground/demo_app).

```sh
$ kubectl get pods -n demoapp
NAME                                 READY   STATUS    RESTARTS   AGE
service1-demo-app-54f6549446-ct8k9   2/2     Running   0          7m14s
service1-demo-app-54f6549446-r4mmq   2/2     Running   0          7m14s
service2-demo-app-759bbcc899-kxgwj   2/2     Running   0          7m13s
service2-demo-app-759bbcc899-njpxj   2/2     Running   0          7m13s
service3-demo-app-788857c7cc-557zj   2/2     Running   0          7m13s
service3-demo-app-788857c7cc-vlchn   2/2     Running   0          7m13s
```

<Zoom>

```mermaid
{@include: ../assets/playground/playground-app.mmd}
```

</Zoom>

The above diagram shows interaction between different services and the policy
running on Aperture Agent:

- `service1` calls `service2`, which then calls `service3`. This call graph is
  programmed in the request payload of the traffic generator.
- `service3` (the last service in the call graph) simulates concurrency
  constraint by limiting the number of requests it can process in parallel.
- Each service simulates an artificial workload by taking a few milliseconds to
  reply for each request.
- Flux Meter is configured on `service3`. Flux Meter helps monitor service-level
  health signals such as latency, which are used in the Latency Gradient policy.
- Concurrency Limiter and Rate Limiter are configured on `service1`. That is,
  when the `service3` is overloaded, load shedding happens on `service1`.

Once all the resources are in the running state, simulated traffic will start
getting generated automatically against the demo application. The traffic is
designed to overload the demo application in order showcase the capabilities of
Aperture.

The load generator is configured to generate the following traffic pattern for
`subscriber`, `guest` and `bot` traffic types:

- Ramp up to `5` concurrent users in `10s`.
- Hold at `5` concurrent users for `2m`.
- Ramp up to `30` concurrent users in `1m` (overloads `service3`).
- Hold at `30` concurrent users for `2m` (overloads `service3`).
- Ramp down to `5` concurrent users in `10s`.
- Hold at `5` concurrent users for `2m`.

Once the traffic is running, you can visualize the decisions made by Aperture in
Grafana. Navigate to [localhost:3000](http://localhost:3000) on your browser to
reach Grafana. You can open the pre-loaded "FluxNinja" dashboard under
"aperture-system" folder to a bunch of useful panels.

<Zoom>

![Grafana Dashboard](../../tutorials/integrations//flow-control/concurrency-limiting/assets/rate-limiting-escalation/rate-limiting-escalation-playground.png)

</Zoom>

:::info

Grafana dashboards url is
[localhost:3000/dashboards](http://localhost:3000/dashboards)

:::

To stop the traffic at any point of time, press the `Stop Wavepool Generator`
button in the `DemoApplications` resource.

![Stop Wavepool Generator](../../assets/img/stop-traffic.png)

To re-start the traffic, press the `Start Wavepool Generator` button in the
`DemoApplications` resource.

:::note

Every time you wish to manually run the traffic, make sure to press the
`Stop Wavepool Generator` button first.

:::

---

## Tools

As Described hereafter, deployment methods assume usage of specific deployment
and configuration/management tools (which must be installed beforehand).

To install the required tools, you have two options:

- Use the tool [ASDF](#install-via-asdf)
- Manually install the tools
  [mentioned here](#tools-required-for-kubernetes-deployment).

### Install via asdf

First,
[download](https://asdf-vm.com/guide/getting-started.html#_2-download-asdf) and
[install](https://asdf-vm.com/guide/getting-started.html#_3-install-asdf)
`asdf`. Then, run the following command in aperture playground directory to
install all the required tools.

```bash
./scripts/install_tools.sh
```

### Tools required for Kubernetes deployment

:::note

Please skip this section in case you already installed the required tools using
[`asdf`](#install-via-asdf).

:::

Tools required are listed below

- **Helm**: It is a package manager for Kubernetes. To install manually, follow
  instructions [here](https://helm.sh/docs/intro/install/).
- **Tanks and Jsonnet Bundler**: Grafana Tanka is a robust configuration utility
  for your Kubernetes cluster, powered by the unique Jsonnet language. Jsonnet
  Bundler is used to manage Jsonnet dependencies. To install manually, follow
  instructions [here](https://tanka.dev/install).
- **Kind**: This allows you to run local Kubernetes clusters. To install
  manually, follow instructions
  [here](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).
- **Kubectl**: It is the command-line tool to interact with Kubernetes clusters.
  To install manually, follow instructions
  [here](https://kubernetes.io/docs/tasks/tools/#kubectl).

## Deploying with Tilt

In case of local deployments and development work, it's nice to be able to
automatically rebuild images and services. Aperture Playground uses Tilt to
achieve this.

### Tilt installation

Tilt can be installed using [`asdf`](#install-via-asdf) or manually by following
instructions [here](https://docs.tilt.dev/install.html).

### Prerequisites - Kubernetes cluster bootstrap {#prerequisites-k8s}

:::note

You can skip this section if you already have a running cluster which is being
pointed by the `kubectl`.

:::

Create a K8s cluster using Kind with configuration file by executing below
command from aperture home directory:

```sh
kind create cluster --config playground/kind-config.yaml
```

This will start a cluster with the name `aperture-playground`.

Once done, you can delete the cluster with following command:

```sh
kind delete cluster --name aperture-playground
```

Alternatively, you can use [`ctlptl`](https://github.com/tilt-dev/ctlptl) to
start a cluster with built-in local registry for Docker images:

```sh
ctlptl apply -f playground/ctlptl-kind-config.yaml
```

Once done, you can delete the cluster and registry with following command:

```sh
ctlptl delete -f playground/ctlptl-kind-config.yaml
```

### Services deployment

Simply run `tilt up` from `playground` directory - it'll automatically start
building and deploying.

You can reach the WebUI by going to <http://localhost:10350> or pressing
(Space).

Tilt should automatically detect new changes to the services, rebuild and
re-deploy them.

Useful flags:

- `--port` or `TILT_PORT` - the port on which WebUI should listen

- `--stream` - will stream both tilt and pod logs to terminal (useful for
  debugging `tilt` itself)

- `--legacy` - if you want a basic, terminal-based frontend

By default, `tilt` will deploy and manage Agent and Controller.

If you want to limit it to only manage some namespace(s) or resource(s), simply
pass their name(s) as additional argument(s).

Examples:

- `tilt up aperture-grafana` - only bring up `grafana` and dependent services
  (`grafana-operator`, ...)
- `tilt up agent demoapp aperture-grafana` - you can mix namespace names and
  resource names, as well as specify as many of them as you want.

If you want to manage only explicitly passed resources/namespaces, you should
pass the `--only` argument:

- `tilt up -- --only aperture-grafana` - only bring up grafana, namespace
  resolving to resources still works

To view the available namespaces and resources, either:

- run `tilt up --stream -- --list-resources`
- read the `DEP_TREE` at the top of `Tiltfile`

To disable automatic updates in Tilt, add `--manual` with the command.

### Teardown

Simply run `tilt down`. All created resources will be deleted.

### Port Forwards {#port-forwards}

Tilt will automatically setup forwarding for the services.

Below is the mapping of the ports being forwarded by Tilt:

| Component  | Container Port | Local Port |
| ---------- | -------------- | ---------- |
| Prometheus | 9090           | 9090       |
| Etcd       | 2379           | 2379       |
| Grafana    | 3000           | 3000       |

## FAQs

### Too many open files "warning"

If you are getting following message in cluster container:

> failed to create fsnotify watcher: to many open files

If `sysctl fs.inotify.{max_queued_events,max_user_instances,max_user_watches}`
less than:

```sh
fs.inotify.max_queued_events=16384
fs.inotify.max_user_instances=1024
fs.inotify.max_user_watches=524288
```

change it, using (temporary method):

```sh
sudo sysctl fs.inotify.max_queued_events=16384
sudo sysctl fs.inotify.max_user_instances=1024
sudo sysctl fs.inotify.max_user_watches=524288
```

or add following lines to `/etc/sysctl.conf`:

```sh
fs.inotify.max_queued_events=16384
fs.inotify.max_user_instances=1024
fs.inotify.max_user_watches=524288
```
