---
title: Blueprints
description: Policies and Dashboards pre-packaged as reusable Blueprints
keywords:
  - jsonnet
  - grafana
  - policy
sidebar_position: 1
---

```mdx-code-block
import {apertureVersion} from '@site/src/version';
import CodeBlock from '@theme/CodeBlock';
```

## Introduction

Aperture comes with a pre-packaged list of [Aperture Policies][policies] and
Grafana Dashboards that can be used both as a guide for creating new policies,
and as ready-to-use Aperture Blueprints for generating policies customized to a
[Service][service] and the use-case.

All Aperture Policies and Grafana Dashboards are written using the
[Jsonnet][jsonnet-lang] language, and can be used both as jsonnet mixins or as
standalone Aperture Blueprints.

[jsonnet-lang]: https://jsonnet.org

## Tools

In order to generate Blueprints please install the pre-requisites [`jb`][jb] and
[`jsonnet`][jsonnet].

[jb]: https://github.com/jsonnet-bundler/jsonnet-bundler
[jsonnet]: https://github.com/google/jsonnet

## Generating Aperture Policies and Grafana Dashboards

The available Aperture Blueprints can be found under <a
href={`https://github.com/fluxninja/aperture/tree/v${apertureVersion}/blueprints/lib/1.0/blueprints/`}>`blueprints/lib/1.0/blueprints`</a>.
In each blueprint, `bundle.libsonnet` can be used to generate the actual
artifacts, and `config.libsonnet` comes with the default configuration for the
given blueprint. This can be overridden by the `--config` option passed to the
`generate-bundle.py` script.

Custom configurations can be merged with blueprints' `config.libsonnet`
resulting in the final configuration, according to jsonnet language rules: keys
can be overwritten by reusing them in the custom configuration and nested
objects can be merged by using `+:` operator. Check the `example` directory for
more information.

To generate files, the script <a
href={`https://github.com/fluxninja/aperture/tree/v${apertureVersion}/blueprints/scripts/generate-bundle.py`}>`blueprints/scripts/generate-bundle.py`</a>
can be used. This script takes arguments for the output directory path where
files will be saved and a path to a config libsonnet file containing blueprint
customization and configuration. The full command using an example config looks
like this:

```sh
cd blueprints/
jb install
./scripts/generate-bundle.py --output _gen --config examples/latency-gradient/example.jsonnet
```

## Using Aperture Blueprints as a jsonnet mixins library

An alternate way of using the Aperture Blueprints is to import them from another
jsonnet project and render Aperture Policies or Grafana Dashboards directly in
jsonnet. This can be also integrated with other Kubernetes deployment tools like
[tanka][tk].

For example, to create a Latency Gradient Policy that can be loaded by the
controller, you need to install Aperture Blueprints library with jsonnet
bundler:

```sh
jb install github.com/fluxninja/aperture/blueprints@main
```

:::info

You can use specific aperture branch instead of _main_ e.g. _stable/v0.2.x_, or
even a specific release tag e.g. _v0.2.2_

:::

You can then create a Policy resource using Jsonnet definitions:

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

## Dry Run Mode

Some of the Policy Blueprints can be configured to generate a Policy in
`Dry Run` Mode. In the `Dry Run` mode, Actuation is disabled and the Policy runs
in observation only mode. The decisions are still evaluated which helps
understand how the policy would behave as the input signals change.

For instance, setting `dynamicConfig.dryRun` option to `true` in the
latency-gradient blueprint would generate a Dry Run policy.

:::note

The outcome of a `Dry Run` policy would be different compared to a policy that
actuates. This is because actuation changes the system being controlled which
would result in a different kind of Signal feedback. While `Dry Run` mode is not
a simulation, it's still useful in understanding the signal processing and
decisions made in each _execution cycle_.

:::

To understand what the above policy does, please see the
[Basic Concurrency Limiting](/tutorials/flow-control/basic-concurrency-limiting.md)
tutorial.

[jsonnet]: https://github.com/google/go-jsonnet

[tk]: https://grafana.com/oss/tanka/ [policies]: /concepts/policy/policy.md
[service]: /concepts/service.md
