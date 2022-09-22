# Aperture Blueprints

This directory provides a set of Aperture policies and pre-defined Grafana dashboards that can
be used both as-is, and as examples to create your own custom policies for the Aperture
controller.

## Usage

Provided policies and dashboards can be used as standalone
blueprints, that get rendered into Aperture policies and corresponding Grafana dashboards. Also, its possible to use these blueprints as a jsonnet library in your own jsonnet files.

### Configure and render Blueprints

You can get from Blueprints to working Aperture policies and Grafana dashboards in three easy steps: install the requirements, configure the policy, generate!

#### Requirements

There are three requirements for the Aperture policy generator: Python 3.8+, jsonnet parser and
jsonnet bundler.

First, make sure that all jsonnet dependencies are installed, using jsonnet bundler:

```sh
jb install
```

#### Configuration

Before running generation you
need a policy configuration file - see `examples/demoapp-latency-gradient.jsonnet` for an example
config for the demoapp.

Policy configuration is specific to policies, and so `blueprints/[policy]/config.libsonnet` as well
as `lib/1.0/policies/` and `lib/1.0/dashboards/` (as mentioned in `config.libsonnet`) should
be consulted to see what can be changed.

#### Generation

Aperture generator script is available under `scripts/aperture-generate.py`. With config file ready, generator can be used to create JSON files in the output directory:

```sh
python scripts/aperture-generate.py --output _gen/ --config examples/demoapp-latency-gradient.jsonnet blueprints/latency-gradient
```

And all generated policies and dashboards will be available under `_gen/policies/` and
`_gen/dashboards`.

### Using Blueprints as a Jsonnet library

When used as a Jsonnet library, aperture blueprints can be used with libraries like [k8s-libsonnet][k8s-libsonnet]
or [grafana-operator-libsonnet][grafana-libsonnet] to generate kubernetes resources that can be deployed
on the cluster.

To install aperture blueprints as a jsonnet dependency, use [jsonnet-bundler][jb]:

```sh
jb install github.com/fluxninja/aperture/blueprints@main
```

Additionally, for this example to work install k8s-libsonnet dependency:

```sh
jb install github.com/jsonnet-libs/k8s-libsonnet/1.24@main
```

Finally, you can create a ConfigMap resource with policy like that:

```jsonnet
local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local latencyGradientPolicy = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/policies/latency-gradient.libsonnet';

local selector = aperture.v1.Selector;
local serviceSelector = aperture.v1.ServiceSelector;
local flowSelector = aperture.v1.FlowSelector;
local controlPoint = aperture.v1.ControlPoint;

local svcSelector =
  selector.new()
  + selector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + selector.withFlowSelector(
    flowSelector.new()
    + flowSelector.withControlPoint(controlPoint.new()
                              + controlPoint.withTraffic('ingress'))
  );

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

This can be then used with jsonnet-based deployment tools, like [tanka][tanka] to generate
and apply resources to the cluster.

[k8s-libsonnet]: https://github.com/jsonnet-libs/k8s-libsonnet
[grafana-libsonnet]: https://github.com/jsonnet-libs/grafana-operator-libsonnet
[jb]: https://github.com/jsonnet-bundler/jsonnet-bundler
[jsonnet]: https://github.com/google/go-jsonnet
[tanka]: https://github.com/grafana/tanka

### Development

See `CONTRIBUTING.md` for the documentation on how to create new blueprints and generate
markdown documentation
