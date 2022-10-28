---
title: Blueprints
description: Policies and Dashboards pre-packaged as reusable Blueprints
keywords:
  - jsonnet
  - grafana
  - policy
sidebar_position: 1
---

## Introduction

Aperture comes with a pre-packaged list of [Aperture Policies][policies] and
Grafana Dashboards that can be used both as a guide for creating new Policies,
and as ready-to-use Blueprints for generating Aperture Policies customized to a
[Service][service] and the use-case.

All Policies and Grafana Dashboards are written using the
[Jsonnet][jsonnet-lang] language, and can be used both as jsonnet mixins or as
standalone Blueprints.

[jsonnet-lang]: https://jsonnet.org

## Initial Setup

Blueprints can be found in the [aperture repository][aperture-repo] under
`blueprints/` directory.

The Blueprint Generator (used to generate Policy files from Blueprints) depends
on [jsonnet][go-jsonnet].

[aperture-repo]: https://github.com/fluxninja/aperture/
[jb]: https://github.com/jsonnet-bundler/jsonnet-bundler
[go-jsonnet]: https://github.com/google/go-jsonnet

## Generating Aperture Policies and Grafana Dashboards

The simplest way to use Aperture Blueprints is to render Blueprints into Policy
and Dashboard files.

To generate files, `blueprints/scripts/generate-bundle.py` can be used. The
script takes as options an output directory path where files will be saved and a
path to a config libsonnet file containing blueprint customization and
configuration.

Under the `blueprints/lib/1.0/blueprints` directory, the currently available
Blueprints can be found. In each Blueprint, `bundle.libsonnet` can be used to
generate the actual artifacts, and `config.libsonnet` comes with the default
configuration for the given Blueprint. This can be overridden by the `--config`
option passed to the `generate-bundle.py` script.

Custom configurations can be merged with Blueprints' `config.libsonnet`
resulting in the final configuration, according to jsonnet language rules: keys
can be overwritten by reusing them in the custom configuration and nested
objects can be merged by using `+:` operator. Check the `example` directory for
more information.

The full command using the example looks like this:

```sh
jb install && ./scripts/generate-bundle.py --output _gen --config examples/latency-gradient/example.jsonnet
```

## Using aperture blueprints as a jsonnet mixins library

An alternate way of using the aperture blueprints is to import them from another
jsonnet project and render Policies or Grafana Dashboards directly in jsonnet.
This can be also integrated with other Kubernetes deployment tools like
[tanka][tk].

For example, to create a Latency Gradient Policy that can be loaded by the
controller, you need to install aperture blueprints library with jsonnet
bundler:

```sh
jb install github.com/fluxninja/aperture/blueprints@main
```

:::info

You can use specific aperture branch instead of _main_ e.g. _stable/v0.2.x_, or
even a specific release tag e.g. _v0.2.2_

:::

You can then create a Policy resource with policy definition like this:

```jsonnet
{@include: ../../tutorials/flow-control/assets/basic-concurrency-limiting/basic-concurrency-limiting.jsonnet}
```

And then, render it with [jsonnet][jsonnet]:

```sh
jsonnet -J vendor [example file].jsonnet  | yq -P
```

After running this command you should see the following contents in the YAML
file:

```yaml
{@include: ../../tutorials/flow-control/assets/basic-concurrency-limiting/basic-concurrency-limiting.yaml}
```

The generated policy can be applied to the running instance of
`aperture-controller` via `kubectl` as follows:

```sh
kubectl apply --namespace aperture-controller --filename [example file].yaml
```

To understand what the above policy does, please see the
[Basic Concurrency Limiting](/tutorials/flow-control/basic-concurrency-limiting.md)
tutorial.

[jsonnet]: https://github.com/google/go-jsonnet
[tk]: https://grafana.com/oss/tanka/
[policies]: /concepts/policy/policy.md
[service]: /concepts/service.md
