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

In addition to CLI-based generation of Aperture Policies and Grafana Dashboard,
Aperture Policies & Blueprints can be imported into your Jsonnet project as a
library. This can also be integrated with other Kubernetes deployment tools like
[Tanka][tk].

For example, to create a latency-based AIMD concurrency limiting policy that can
be loaded by the controller, you need to install Aperture blueprints library
with Jsonnet Bundler in your project:

```sh
jb install github.com/fluxninja/aperture/blueprints@main
```

:::tip

You can use a specific Aperture
[GitHub repository](https://github.com/fluxninja/aperture) branch instead of
`main`, for example, `stable/v0.2.x`, or even a specific release tag, for
example, `v0.2.2` to match your Aperture Controller installation version.

:::

You can then create a Policy resource using Jsonnet definitions:

```jsonnet
{@include: ../../applying-policies/signal-processing/assets/detecting-overload/detecting-overload.jsonnet}
```

And then, render it with [Jsonnet][jsonnet]:

```sh
jsonnet -J vendor <example_file>.jsonnet  | yq -P
```

After running this command, you should see the following contents in the YAML
file:

```yaml
{@include: ../../applying-policies/signal-processing/assets/detecting-overload/detecting-overload.yaml}
```

The generated policy can be applied to the running instance of
`aperture-controller` using `kubectl` as follows:

```sh
kubectl apply --namespace aperture-controller --filename <example_file>.yaml
```

[jsonnet]: https://jsonnet.org/
[tk]: https://tanka.dev/
