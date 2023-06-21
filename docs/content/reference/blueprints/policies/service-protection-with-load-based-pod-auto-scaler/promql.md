---
title: Service Protection Based on PromQL Query and Load-based Pod Auto-Scaler
keywords:
  - blueprints
sidebar_position: 3
sidebar_label:
  Service Protection Based on PromQL Query and Load-based Pod Auto Scaler
---

## Introduction

This policy detects traffic overloads and cascading failure build-up by
comparing the value of a metric against a static threshold. A gradient
controller calculates a proportional response to limit accepted concurrency. The
concurrency is reduced by a multiplicative factor when the service is
overloaded, and increased by an additive factor while the service is no longer
overloaded. An auto-scaler controller is used to dynamically adjust the number
of instances or resources allocated to a service based on workload demands. The
basic service protection policy protects the service from sudden traffic spikes.
It is necessary to scale the service to meet demand in case of a persistent
change in load.

At a high level, this policy works as follows:

- PromQL-based overload detection: A PromQL query on an arbitrary metric
  generates a periodic signal. The signal is compared against a static set point
  threshold to detect an overload.
- Gradient Controller: Set point and signal are fed to the gradient controller
  that calculates the proportional response to adjust the accepted concurrency
  (Control Variable).
- Integral Optimizer: When the service is detected to be in the normal state, an
  integral optimizer is used to additively increase the concurrency of the
  service in each execution cycle of the circuit. This design allows warming-up
  a service from an initial inactive state. This also protects applications from
  sudden spikes in traffic, as it sets an upper bound to the concurrency allowed
  on a service in each execution cycle of the circuit based on the observed
  incoming concurrency.
- Load Scheduler and Actuator: The Accepted Concurrency at the service is
  throttled by a
  [weighted-fair queuing scheduler](/concepts/flow-control/components/load-scheduler.md).
  The output of the adjustments to accepted concurrency made by gradient
  controller and optimizer logic are translated to a load multiplier that is
  synchronized with Aperture Agents through etcd. The load multiplier adjusts
  (increases or decreases) the token bucket fill rates based on the incoming
  concurrency observed at each agent.
- An _Auto Scaler_ that adjusts the number of replicas of the Kubernetes
  Deployment for the service.
- Load-based scale-out is done based on `OBSERVED_LOAD_MULTIPLIER` signal from
  the blueprint. This signal measures the fraction of traffic that the _Load
  Scheduler_ is throttling into a queue. The _Auto Scaler_ is configured to
  scale-out based on a _Gradient Controller_ using this signal and a setpoint of
  1.0.
- Periodic scale in can be defined using the
  `policy.auto_scaling.periodic_decrease` parameter. This allows the policy to
  periodically explore whether the service can be scaled down without impacting
  performance.
- Additional scale out and scale in criteria can be defined on arbitrary metrics
  using `policy.auto_scaling.promql_scale_out_controllers` and
  `policy.auto_scaling.promql_scale_in_controllers` parameters.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection-with-load-based-pod-auto-scaler/promql`}>policies/service-protection-with-load-based-pod-auto-scaler/promql</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Configuration for the Service Protection policy.'
    type='Object (policies/service-protection/promql:param:policy)'
    reference='../../../blueprints/policies/service-protection/promql#policy'
    value='{"auto_scaling": {"dry_run": false, "periodic_decrease": {"period": "60s", "scale_in_percentage": 10}, "promql_scale_in_controllers": [], "promql_scale_out_controllers": [], "scaling_backend": {"kubernetes_replicas": {"kubernetes_object_selector": "__REQUIRED_FIELD__", "max_replicas": "__REQUIRED_FIELD__", "min_replicas": "__REQUIRED_FIELD__"}}, "scaling_parameters": {"scale_in_alerter": {"alert_name": "Auto-scaler is scaling in"}, "scale_in_cooldown": "40s", "scale_out_alerter": {"alert_name": "Auto-scaler is scaling out"}, "scale_out_cooldown": "30s"}}, "components": [], "evaluation_interval": "1s", "policy_name": "__REQUIRED_FIELD__", "promql_query": "__REQUIRED_FIELD__", "resources": {"flow_control": {"classifiers": []}}, "service_protection_core": {"adaptive_load_scheduler": {"alerter": {"alert_name": "Load Throttling Event"}, "gradient": {"max_gradient": 1, "min_gradient": 0.1, "slope": -1}, "load_multiplier_linear_increment": 0.0025, "load_scheduler": {"selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "max_load_multiplier": 2}, "dry_run": false, "overload_confirmations": []}, "setpoint": "__REQUIRED_FIELD__"}'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="dashboard"></a>

<ParameterDescription
    name='dashboard'
    description='Configuration for the Grafana dashboard accompanying this policy.'
    type='Object (policies/service-protection/promql:param:dashboard)'
    reference='../../../blueprints/policies/service-protection/promql#dashboard'
    value='{"datasource": {"filter_regex": "", "name": "$datasource"}, "extra_filters": {}, "refresh_interval": "15s", "time_from": "now-15m", "time_to": "now", "title": "Aperture Service Protection", "variant_name": "PromQL Output"}'
/>

<!-- vale on -->

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<!-- vale off -->

<a id="dry-run"></a>

<ParameterDescription
    name='dry_run'
    description='Dynamic configuration for setting dry run mode at runtime without restarting this policy. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. The Auto Scaler does not perform any scaling in dry mode. This mode is useful for observing the behavior of load scheduler and auto scaler without disrupting any real deployment or traffic.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
