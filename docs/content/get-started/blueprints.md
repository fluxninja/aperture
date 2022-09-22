---
title: Policy Blueprints
description: Policies and Dashboards pre-packaged as reusable Blueprints
keywords:
  - jsonnet
  - grafana
sidebar_position: 3
---

## Introduction

Aperture comes with a pre-packaged list of [Aperture Policies][policies] and Grafana Dashboards that
can be used both as a guide for creating new Policies, and as ready-to-use
Blueprints for generating Aperture Policies customized to a [Service][service].

All Dashboards and Policies are written using the [Jsonnet][jsonnet-lang]
language, and can be used both as jsonnet mixins or as standalone Blueprints.

[jsonnet-lang]: https://jsonnet.org

## Initial Setup

Blueprints can be found in the [aperture repository][aperture-repo] under `blueprints/`
directory. See the [blueprints README.md][blueprints-readme] for the list of required
tools and instructions on installing jsonnet dependencies with a help of
[jsonnet bundler][jb].

The Blueprint Generator (used to generate Policy files from Blueprints) also
depends on Python 3.8+ and [jsonnet][go-jsonnet].

[aperture-repo]: https://github.com/fluxninja/aperture/
[blueprints-readme]: https://github.com/fluxninja/aperture/blob/main/blueprints/README.md
[jb]: https://github.com/jsonnet-bundler/jsonnet-bundler
[go-jsonnet]: https://github.com/google/go-jsonnet

## Generating Aperture Policies and Grafana Dashboards

The simplest way to use Aperture Blueprints is to render blueprints into
policy and dashboard files.

To generate files, `blueprints/scripts/aperture-generate.py` can be used:

```sh
$ ./scripts/aperture-generate.py --help
usage: aperture-generate.py [-h] [--verbose] [--output OUTPUT] [--config CONFIG] BLUEPRINT

Aperture Policies & Dashboards generator utility.

This utility can be used to generate Aperture Policies and Grafana Dashboards "in-place". Check [blueprint's README.md][blueprints-readme] for more
details.

positional arguments:
  BLUEPRINT        Aperture blueprint path

options:
  -h, --help       show this help message and exit
  --verbose        Whether to log verbose messages to stderr
  --output OUTPUT  Output directory for json files
  --config CONFIG  jsonnet file with blueprint configuration
```

This script takes as options an output directory path where files will be
saved and a path to a `config.libsonnet` file with local blueprint
configuration. It also takes the BLUEPRINT argument, which is a path to the
blueprint under the `blueprints/` directory.

Under the `blueprints/` directory, the currently available Blueprints can be
found. Each blueprint consists of at least two files: `config.libsonnet` and
`main.libsonnet`. `main.libsonnet` bundles actual Policy and dashboard code
(available under `lib/1.0`) into Blueprints, and `config.libsonnet` comes with
the default configuration for the given Policy. This can be overridden by the
`--config` option passed to the `aperture-generate.py` script.

Custom configurations will be merged with Blueprints' `config.libsonnet`
resulting in the final configuration, according to jsonnet language rules: keys
can be overwritten by reusing them in the custom configuration and nested
objects can be merged by using `+:` operator. Check the `examples/` directory
for more information.

The full command using the demoapp-latency-grand example looks like this:

```sh
jb install
./scripts/aperture-generate.py --output _gen --config examples/demoapp-latency-gradient.jsonnet Blueprints/latency-gradient
```

[blueprints-readme]: https://github.com/fluxninja/aperture/blob/main/blueprints/README.md

## Using aperture blueprints as a jsonnet mixins library

An alternate way of using the aperture blueprints is to import them from another
jsonnet project and render Policies or Dashboards directly in jsonnet.

For example, to create a Latency Gradient Policy that can be loaded by
the controller, you need to install aperture blueprints library with jsonnet bundler:

```sh
jb install github.com/fluxninja/aperture/blueprints@main
```

:::info

You can use specific aperture branch instead of _main_ e.g. _stable/v0.2.x_, or even
a specific release tag e.g. _v0.2.2_

:::

You can then create a Policy resource with policy definition like this:

```jsonnet
local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local latencyGradientPolicy = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/policies/latency-gradient.libsonnet';

local selector = aperture.v1.Selector;
local controlPoint = aperture.v1.ControlPoint;

local svcSelector =
  selector.new()
  + selector.withAgentGroup('default')
  + selector.withService('service1-demo-app.demoapp.svc.cluster.local')
  + selector.withControlPoint(controlPoint.new()
                              + controlPoint.withTraffic('ingress'));

local policy = latencyGradientPolicy({
  policyName: 'service1-demoapp',
  fluxMeterSelector: svcSelector,
  concurrencyLimiterSelector: svcSelector,
}).policy;

local policyResource = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'service1',
  },
  spec: policy,
};

[
  policyResource,
]
```

And then, render it with [jsonnet][jsonnet]:

```sh
jsonnet --yaml-stream -J vendor [example file].jsonnet
```

This can be also integrated with other Kubernetes deployment tools like
[tanka][tk].

[jsonnet]: https://github.com/google/go-jsonnet
[tk]: https://grafana.com/oss/tanka/
[policies]: /concepts/policy/policy.md
[service]: /concepts/service.md
