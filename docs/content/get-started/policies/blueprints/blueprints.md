---
title: Blueprints
description: Policies and Dashboards pre-packaged as reusable Blueprints
keywords:
  - jsonnet
  - grafana
  - policy
sidebar_position: 2
---

```mdx-code-block
import {apertureVersion} from '../../../apertureVersion.js';
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

## Generating Aperture Policies and Grafana Dashboards {#generating-aperture-policies-and-grafana-dashboards}

The available Aperture Blueprints can be found under <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/lib/1.0/blueprints/`}>`blueprints/lib/1.0/blueprints`</a>.
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
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/scripts/generate-bundle.py`}>`blueprints/scripts/generate-bundle.py`</a>
can be used. This script takes arguments for the output directory path where
files will be saved and a path to a config libsonnet file containing blueprint
customization and configuration. Assuming you are at <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/`}>`Aperture root`</a>,
run the following commands to generate an example config:

```sh
cd blueprints/
jb install
./scripts/generate-bundle.py --output _gen --config examples/latency-aimd-concurrency-limiting/example.jsonnet
```

## Dry Run Mode

Some of the Policy Blueprints can be configured to generate a Policy in
`Dry Run` Mode. In the `Dry Run` mode, Actuation is disabled and the Policy runs
in observation only mode. The decisions are still evaluated which helps
understand how the policy would behave as the input signals change.

For instance, setting `dynamicConfig.dryRun` option to `true` in the
latency-aimd-concurrency-limiting blueprint would generate a Dry Run policy.

:::note

The outcome of a `Dry Run` policy would be different compared to a policy that
actuates. This is because actuation changes the system being controlled which
would result in a different kind of Signal feedback. While `Dry Run` mode is not
a simulation, it's still useful in understanding the signal processing and
decisions made in each _execution cycle_.

:::

To understand what the above policy does, please see the
[Basic Concurrency Limiting](/tutorials/integrations/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
tutorial.

[jsonnet]: https://github.com/google/go-jsonnet
[tk]: https://grafana.com/oss/tanka/
[policies]: /concepts/policy/policy.md
[service]: /concepts/integrations/flow-control/service.md
