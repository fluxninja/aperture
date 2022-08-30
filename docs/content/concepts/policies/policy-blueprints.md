---
title: Policy Blueprints
description: Policies and dashboard pre-packaged as reusable blueprints
keywords:
  - jsonnet
  - grafana
---

# Policy Blueprints

## Introduction

Aperture comes with a pre-packaged list of policies and grafana dashboards that
can be used both as a guide for creating new policies, and as ready-to-use
blueprints that can be directly used for configuring Aperture Agent behaviour.

All dashboards and policies are written using [Jsonnet][jsonnet-lang] language,
and can be used both as jsonnet mixins, and standalone blueprints.

[jsonnet-lang]: https://jsonnet.org

## Initial Setup

All blueprints are available from a separate [repository][aperture-blueprints].
See repository [README.md][blueprints-readme] for the list of required tools,
and how to install jsonnet dependencies with a help of [jsonnet bundler][jb].

Blueprint Generator (used to generate JSON files from blueprints) also depends
on a Python 3.8+ and [jsonnet][go-jsonnet].

[k8s-libsonnet]: https://github.com/jsonnet-libs/k8s-libsonnet
[aperture-blueprints]: https://github.com/fluxninja/aperture-blueprints
[blueprints-readme]: https://github.com/fluxninja/aperture-blueprints/blob/main/README.md
[jb]: https://github.com/jsonnet-bundler/jsonnet-bundler
[go-jsonnet]: https://github.com/google/go-jsonnet

## Generating JSON Grafana Dashboards and Aperture Policies

The simplest way of using blueprints repository is to render blueprints into
JSON policy and dashboard files.

To generate JSON files, `scripts/aperture-generate.py` can be used:

```sh
$ ./scripts/aperture-generate.py --help
usage: aperture-generate.py [-h] [--verbose] [--output OUTPUT] [--config CONFIG] BLUEPRINT

Aperture policies & dashboards generator utility.

This utility can be used to generate Aperture policies and Grafana dashboards "in-place". Check [aperture-blueprint's README.md](https://github.com/fluxninja/aperture-blueprints/blob/main/README.md) for more
details.

positional arguments:
  BLUEPRINT        Aperture blueprint path

options:
  -h, --help       show this help message and exit
  --verbose        Whether to log verbose messages to stderr
  --output OUTPUT  Output directory for json files
  --config CONFIG  jsonnet file with blueprint configuration
```

This script takes as options an output directory path, where JSON files will be
saved under, and a path to a `config.libsonnet` file with local blueprint
configuration. It also takes BLUEPRINT argument, which is a path to the
blueprint under `blueprints/` directory.

Under `blueprints/` directory, there are all currently available blueprints,
each one constiting of at least two files: `config.libsonnet` and
`main.libsonnet`. `main.libsonnet` bundles actual policy and dashboard code
(available under `lib/1.0`) into blueprints, and `config.libsonnet` comes with
the default configuration for the given policy. This is what can be overriden by
`--config` option passed to `aperture-generate.py` script.

Custom configuration will be merged with blueprints' `config.libsonnet`
resulting in the final configuration, according to jsonnet language rules: keys
can be overwritten by reusing them in custom configuration, and nested objects
can be merged by using `+:` operator. Check `examples/` directory for more
information.

The full command, using demoapp-latency-grand example, looks like that:

```sh
jb install
./scripts/aperture-generate.py --output _gen --config examples/demoapp-latency-gradient.jsonnet blueprints/latency-gradient
```

## Using aperture-blueprints as a jsonnet mixins library

An alternate way of using aperture-blueprints repository is to import it from
another jsonnet project and render policies or dashboards directly in jsonnet.

As an example, to create a ConfigMap with aperture policies that can be loaded
by the controller, you need to install aperture-blueprints with jsonnet bundler:

```sh
jb install github.com/fluxninja/aperture-blueprints@main
```

Additionally, for this example to work install k8s-libsonnet dependency:

```sh
jb install github.com/jsonnet-libs/k8s-libsonnet/1.24@main
```

Finally, you can create a ConfigMap resource with policy like that:

```jsonnet
local k = import "github.com/jsonnet-libs/k8s-libsonnet/1.24/main.libsonnet";

local latencyGradientPolicy = import "github.com/fluxninja/aperture-blueprints/lib/1.0/policies/latency-gradient.libsonnet";

local policy = latencyGradientPolicy({
  fluxmeterName: "demo1-demo-app",
  serviceSelector+: {
    service: "demo1-demo-app.demoapp.svc.cluster.local"
  },
}).policy;

[
    k.core.v1.configMap.new("policies")
	+ k.core.v1.configMap.metadata.withLabels({ "fluxninja.com/validate": "true"})
	+ k.core.v1.configMap.withData({
	  "demoapp-latency-gradient.yaml": std.manifestYamlDoc(policy, quote_keys=false)
	})
]
```

And then, render it with [jsonnet][jsonnet]:

```sh
jsonnet --yaml-stream -J vendor [example file].jsonnet
```

This can be also integrated with other kubernetes deployment tools like
[tanka][tk]

[jsonnet]: https://github.com/google/go-jsonnet
[tk]: https://grafana.com/oss/tanka/
