---
title: Playground
keywords:
  - playground
  - proof of concept
  - policies
  - rate limit
  - concurrency control
sidebar_position: 2
---

# Playground

Playground is a Kubernetes-based environment for exploring the capabilities of
Aperture. Additionally, it is used as a development environment for Aperture.
The playground uses [Tilt](https://tilt.dev/) for orchestrating the deployments
in Kubernetes. Tilt watches for changes to local files and auto-deploys any
resources that change. This is very convenient for getting quick feedback during
development of Aperture.

Playground deploys resources to the Kubernetes cluster that `kubectl` on your
machine points at. For convience, this README includes instructions for
deploying a local Kubernetes cluster using [Kind](https://kind.sigs.k8s.io/).

## Tools

Described hereafter deployment methods assume usage of specific deployment and
configuration/management tools (which must be installed beforehand).

To install required ones, you can use [ASDF](https://asdf-vm.com/) OR install
manually (check [Tools used for k8s deployment](#tools-used-for-k8s-deployment)
).

When using `asdf`:

- [Download](https://asdf-vm.com/guide/getting-started.html#_2-download-asdf)
  and [install](https://asdf-vm.com/guide/getting-started.html#_3-install-asdf)
  `asdf`
- Add intended plugins (tools/applications which will be managed by `asdf`) e.g.
  `asdf plugin-add terraform`
- Install tools: `asdf install`

> Note: Last command will install tools which have been added as plugins and
> which are defined/versioned in `.tool-versions` file

### Tools required for k8s deployment

Tools which are required for local k8s deployment:

#### _Helm_

Helm is a package manager for k8s.

To install manually, follow instructions: <https://helm.sh/docs/intro/install/>

#### _Tanka and Jsonnet Bundler_

Grafana Tanka is the robust configuration utility for your Kubernetes cluster,
powered by the unique Jsonnet language.

Jsonnet Bundler is used to manage Jsonnet dependencies.

To install manually, follow instructions: <https://tanka.dev/install>

#### _Local k8s cluster_

May use [`kind`](https://kind.sigs.k8s.io/docs/user/quick-start/).

#### _kubectl_

The Kubernetes command line tool. Follow the instructions:
<https://kubernetes.io/docs/tasks/tools/#kubectl>

#### _Alpha features_

Agent core service uses feature gate for managing node-local traffic:
<https://kubernetes.io/docs/concepts/services-networking/service-traffic-policy/>

For kubernetes 1.21 it requires a feature gate activation. For kubernetes 1.22
it's in beta so nothing needs to be added to cluster config.

```yaml
featureGates:
  ServiceInternalTrafficPolicy: true
```

## Deploying with Tilt

In case of local deployments and development work, it's nice to be able to
automatically rebuild images and services.

This can be achieved by using `tilt`.

> Note: This builds up on tools mentioned earlier, their installation and
> configuration is required.

### Tilt installation

Tilt can be installed with `asdf install` or manually
<https://docs.tilt.dev/install.html>.

### Prerequisites - k8s cluster bootstrap

Create a K8s cluster using Kind with configuration file:

```sh
kind create cluster --config kind-config.yaml
```

This will start a cluster with the name `aperture-playground`.

Once done, you can delete the cluster with following command:

```sh
kind delete cluster --name aperture-playground
```

Alternatively, you can use [`ctlptl`](https://github.com/tilt-dev/ctlptl) to
start a cluster with built-in local registry for Docker images:

```sh
ctlptl apply -f ctlptl-kind-config.yaml
```

### Services deployment

Simply run `tilt up` - it'll automatically start building and deploying!

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

### To many open files "warning"

If you are getting following message in cluster container:

> failed to create fsnotify watcher: to many open files

If `sysctl fs.inotify.{max_queued_events,max_user_instances,max_user_watches}`
less than:

```bash
fs.inotify.max_queued_events = 16384
fs.inotify.max_user_instances = 1024
fs.inotify.max_user_watches = 524288
```

change it, using (temporary method):

```bash
sudo sysctl fs.inotify.max_queued_events = 16384
sudo sysctl fs.inotify.max_user_instances = 1024
sudo sysctl fs.inotify.max_user_watches = 524288
```

or add following lines to `/etc/sysctl.conf`:

```bash
fs.inotify.max_queued_events = 16384
fs.inotify.max_user_instances = 1024
fs.inotify.max_user_watches = 524288
```

### Teardown

Simply run `tilt down`. All created resources will be deleted.

### Port forwards

Tilt will automatically setup forwarding for the services.

Below is the mapping of the ports being forwarded by Tilt:

| Component  | Container Port | Local Port |
| ---------- | -------------- | ---------- |
| Agent      | 80             | 8089       |
| Controller | 80             | 8087       |
| Prometheus | 9090           | 9090       |
| Etcd       | 2379           | 2379       |
| Grafana    | 3000           | 3000       |
