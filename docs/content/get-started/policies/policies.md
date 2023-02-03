---
title: Generate and Apply Policies
description: How to generate and apply Policies in Aperture
keywords:
  - policy
  - jsonnet
  - grafana
  - policy
sidebar_position: 3
---

```mdx-code-block
import {apertureVersion} from '../../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
```

## Introduction

Aperture comes with a pre-packaged list of [Aperture Policies][policies] and
Grafana Dashboards that can be used both as a guide for creating new policies,
and as ready-to-use Aperture Blueprints for generating policies customized to a
[Service][service] and the use-case.

In order to install Aperture Blueprints and generate policies, you can use
`aperturectl`.

## Manage Aperture Blueprints

Let us first install `aperturectl`.

```shell
# Change directory to the root of the Aperture repository
$ cd fluxninja/aperture
# Change directory to aperturectl
$ cd cmd/aperturectl
# Build the binary
$ go build .
# Move the binary to a directory in your PATH (you might need to use sudo)
$ mv aperturectl /usr/local/bin
# Make sure the binary is executable
$ aperturectl version
aperturectl v0.0.1
```

Now that we have `aperturectl` installed, we can use it to manage Aperture
Blueprints.

```shell
# Pull the latest blueprints
$ aperturectl blueprints pull
GET https://github.com/fluxninja/aperture/archive/3159f7692abb3c7d2aa9251a2f9d9be7813c61a3.tar.gz 200
GET https://github.com/grafana/grafonnet-lib/archive/30280196507e0fe6fa978a3e0eaca3a62844f817.tar.gz 200
GET https://github.com/grafana/jsonnet-libs/archive/6ea0d80c5205311a68deb4c547992b32cfb9edd8.tar.gz 200
GET https://github.com/jsonnet-libs/k8s-libsonnet/archive/85543e49238903ac14b486321bd3d60fef09d9ef.tar.gz 200
```

The above command will pull the latest Aperture Blueprints from the Aperture
repository and store it locally. You can also specify the version of the
Blueprints you want to pull.

```shell
# Pull a specific version of the blueprints
$ aperturectl blueprints pull --version v0.21.0
```

To list the policies associated in the latest version of Blueprints that is
installed locally, run the following command.

```shell
# List the currently installed blueprints
$ aperturectl blueprints list
Blueprints: main
dashboards/signals-dashboard
policies/latency-aimd-concurrency-limiting
policies/static-rate-limiting
```

To list all versions of the Blueprints that are installed locally, run the
following command, you can include the `--all` flag.

To learn more about `aperturectl` and the other commands it supports, run
`aperturectl --help`.

## Generate Aperture Policies and Grafana Dashboards {#generating-aperture-policies-and-grafana-dashboards}

Once you have the Blueprints installed locally, you can use `aperturectl` to
generate Aperture Policies and Grafana Dashboards.

Let us say that you want to generate the `policies/static-rate-limiting` policy
with the following values file `rate-limiting-values.yaml`:

```yaml
common:
  policy_name: static-rate-limiting
dashboard:
  datasource:
    filter_regex: ""
    name: $datasource
  refresh_interval: 10s
policy:
  classifiers: []
  evaluation_interval: 300s
  rate_limiter:
    dynamic_config:
      overrides:
        - label_value: gold
          limit_scale_factor: 1
    flow_selector:
      flow_matcher:
        control_point: ingress
      service_selector:
        agent_group: default
        service: service1-demo-app.demoapp.svc.cluster.local
    parameters:
      label_key: http.request.header.user_type
      lazy_sync:
        enabled: true
        num_sync: 5
      limit_reset_interval: 1s
    rate_limit: "50.0"
```

You can run the following command to generate the policy.

```shell
# Generate the policy
$ aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting-values.yaml
GET https://github.com/fluxninja/aperture/archive/3159f7692abb3c7d2aa9251a2f9d9be7813c61a3.tar.gz 200
{"level":"info","service":"aperturectl","time":"2023-02-02T17:16:15-08:00","caller":"cmd/blueprints/generate.go104","message":"Stored all the manifests at 'main/policies/static-rate-limiting'."}
# See what was generated
$ tree main/policies/static-rate-limiting
main/policies/static-rate-limiting
├── dashboards
│   └── static-rate-limiting.json
├── graphs
│   ├── static-rate-limiting.dot
│   └── static-rate-limiting.svg
└── policies
    └── static-rate-limiting.yaml

3 directories, 4 files
```

## Apply Policy

From the above example, you can see that the policy was generated in the
`main/policies/static-rate-limiting` directory.

```shell
# change directory to the generated policy
$ cd main/policies/static-rate-limiting
# apply the policy
$ kubectl apply -f policies/static-rate-limiting.yaml -n aperture-controller
```

Run the following command to check if the Aperture Policy was created.

```shell
# get the policy
$ kubectl get policies.fluxninja.com -n aperture-controller
```

The Aperture Policy runtime can be visualized in Grafana or any other Prometheus
compatible analytics tool. Refer to the Prometheus compatible metrics available
from [Controller][controller-metrics] and [Agent][agent-metrics]. Some of the
Policy [Blueprints][blueprints] come with recommended Grafana dashboards.

## Delete Policy

Run the following command to delete the above Aperture Policy:

```bash
# delete the policy
$ kubectl delete policies.fluxninja.com -n aperture-controller static-rate-limiting
```

[controller-metrics]: /reference/observability/prometheus-metrics/controller.md
[agent-metrics]: /reference/observability/prometheus-metrics/agent.md
[blueprints]: /get-started/policies/blueprints/blueprints.md
[policies]: /concepts/policy/policy.md
[service]: /concepts/integrations/flow-control/service.md
