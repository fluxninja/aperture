---
title: CPU Overload Protection (Kubernetes)
---

## Introduction

By default, this policy detects when the Kubernetes Pod is overloaded using the
pod CPU utilization metric. The policy is based on the
[adaptive load scheduling](/reference/configuration/spec.md#adaptive-load-scheduler)
component.

All the Kubernetes related metrics are collected by the
[Kubeletstats OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/kubeletstatsreceiver)
so if the system under observation requires using different metrics for the
overload confirmation, the
[list of available metrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/kubeletstatsreceiver/metadata.yaml)
can be used to configure the policy. The following PromQL query (with
appropriate filters) is used as `SIGNAL` for the load scheduler:

```promql
avg(k8s_pod_cpu_utilization_ratio)
```

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/load-scheduling/cpu-overload-protection-k8s`}>load-scheduling/cpu-overload-protection-k8s</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### policy {#policy}

<!-- vale on -->

<!-- vale off -->

<a id="policy-components"></a>

<ParameterDescription
    name='policy.components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../../configuration/spec#component'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-policy-name"></a>

<ParameterDescription
    name='policy.policy_name'
    description='Name of the policy.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-resources"></a>

<ParameterDescription
    name='policy.resources'
    description='Additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../../configuration/spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.load_scheduling_core {#policy-load-scheduling-core}

<!-- vale on -->

<!-- vale off -->

<a id="policy-load-scheduling-core-dry-run"></a>

<ParameterDescription
    name='policy.load_scheduling_core.dry_run'
    description='Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-load-scheduling-core-kubelet-overload-confirmations"></a>

<ParameterDescription
    name='policy.load_scheduling_core.kubelet_overload_confirmations'
    description='Overload confirmation signals from kubelet.'
    type='Object (kubelet_overload_confirmations)'
    reference='#kubelet-overload-confirmations'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-load-scheduling-core-overload-confirmations"></a>

<ParameterDescription
    name='policy.load_scheduling_core.overload_confirmations'
    description='List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.'
    type='Array of Object (overload_confirmation)'
    reference='#overload-confirmation'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-load-scheduling-core-aiad-load-scheduler"></a>

<ParameterDescription
    name='policy.load_scheduling_core.aiad_load_scheduler'
    description='Parameters for AIMD throttling strategy.'
    type='Object (aperture.spec.v1.AIADLoadSchedulerParameters)'
    reference='../../configuration/spec#a-i-a-d-load-scheduler-parameters'
    value='{"alerter": {"alert_name": "AIAD Load Throttling Event"}, "load_multiplier_linear_decrement": 0.05, "load_multiplier_linear_increment": 0.025, "load_scheduler": {"selectors": [{"control_point": "__REQUIRED_FIELD__"}]}, "max_load_multiplier": 2, "min_load_multiplier": 0}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-load-scheduling-core-setpoint"></a>

<ParameterDescription
    name='policy.load_scheduling_core.setpoint'
    description='Setpoint.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.kubernetes_object_selector {#policy-kubernetes-object-selector}

<!-- vale on -->

<!-- vale off -->

<a id="policy-kubernetes-object-selector-api-version"></a>

<ParameterDescription
    name='policy.kubernetes_object_selector.api_version'
    description='API version of the object to protect.'
    type='string'
    reference=''
    value='"apps/v1"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-kubernetes-object-selector-kind"></a>

<ParameterDescription
    name='policy.kubernetes_object_selector.kind'
    description='Kind of the object to protect.'
    type='string'
    reference=''
    value='"Deployment"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-kubernetes-object-selector-name"></a>

<ParameterDescription
    name='policy.kubernetes_object_selector.name'
    description='Name of the object to protect.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-kubernetes-object-selector-namespace"></a>

<ParameterDescription
    name='policy.kubernetes_object_selector.namespace'
    description='Namespace of the object to protect.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

### Schemas

<!-- vale off -->

#### driver_criteria {#driver-criteria}

<!-- vale on -->

<!-- vale off -->

<a id="driver-criteria-enabled"></a>

<ParameterDescription
    name='enabled'
    description='Enables the driver.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="driver-criteria-threshold"></a>

<ParameterDescription
    name='threshold'
    description='Threshold for the driver.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### overload_confirmation_driver {#overload-confirmation-driver}

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-driver-pod-cpu"></a>

<ParameterDescription
    name='pod_cpu'
    description='The driver for using CPU usage as overload confirmation.'
    type='Object (driver_criteria)'
    reference='#driver-criteria'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-driver-pod-memory"></a>

<ParameterDescription
    name='pod_memory'
    description='The driver for using memory usage as overload confirmation.'
    type='Object (driver_criteria)'
    reference='#driver-criteria'
    value='{}'
/>

<!-- vale on -->

---

<!-- vale off -->

#### kubelet_overload_confirmations {#kubelet-overload-confirmations}

<!-- vale on -->

<!-- vale off -->

<a id="kubelet-overload-confirmations-criteria"></a>

<ParameterDescription
    name='criteria'
    description='Criteria for overload confirmation.'
    type='Object (overload_confirmation_driver)'
    reference='#overload-confirmation-driver'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="kubelet-overload-confirmations-infra-context"></a>

<ParameterDescription
    name='infra_context'
    description='Kubernetes selector for scraping metrics.'
    type='Object (aperture.spec.v1.KubernetesObjectSelector)'
    reference='../../configuration/spec#kubernetes-object-selector'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### overload_confirmation {#overload-confirmation}

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the overload confirmation criteria.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperture-cli/aperturectl/dynamic-config/apply/apply.md)
at runtime, without reloading the policy.

:::

### Parameters

<!-- vale off -->

<a id="dry-run"></a>

<ParameterDescription
    name='dry_run'
    description='Dynamic configuration for setting dry run mode at runtime without restarting this policy. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. It is useful for observing the behavior of load scheduler without disrupting any real traffic.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
