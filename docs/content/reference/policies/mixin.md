---
title: Aperture Policies and Blueprints as a Jsonnet Mixins Library
description: Using Aperture Blueprints as a Jsonnet mixins library
keywords:
  - jsonnet
  - grafana
  - policy
sidebar_position: 3
sidebar_label: Jsonnet Mixins Library
---

In addition to CLI based generation of Aperture Policies and Grafana Dashboard,
Aperture Policies & Blueprints can be imported into your Jsonnet project as a
library. This can be also integrated with other Kubernetes deployment tools like
[tanka][tk].

For example, to create a Latency based AIMD Concurrency Limiting Policy that can
be loaded by the controller, you need to install Aperture Blueprints library
with jsonnet bundler in your project:

```sh
jb install github.com/fluxninja/aperture/blueprints@main
```

:::tip

You can use specific Aperture branch instead of _main_ e.g. _stable/v0.2.x_, or
even a specific release tag e.g. _v0.2.2_ to match your Aperture Controller
installation version.

:::

You can then create a Policy resource using Jsonnet definitions:

```jsonnet
{@include: ../../tutorials/integrations/flow-control/concurrency-limiting/assets/basic-concurrency-limiting/basic-concurrency-limiting.jsonnet}
```

And then, render it with [jsonnet][jsonnet]:

```sh
jsonnet -J vendor [example file].jsonnet  | yq -P
```

After running this command you should see the following contents in the YAML
file:

```yaml
{@include: ../../tutorials/integrations/flow-control/concurrency-limiting/assets/basic-concurrency-limiting/basic-concurrency-limiting.yaml}
```

The generated policy can be applied to the running instance of
`aperture-controller` via `kubectl` as follows:

```sh
kubectl apply --namespace aperture-controller --filename [example file].yaml
```
