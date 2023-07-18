---
title: Load Scheduling with Average Latency Feedback
---

## Introduction

This policy detects traffic overloads and cascading failure build-up by
comparing the real-time latency with a historical average. A gradient controller
calculates a proportional response to limit the accepted token (or request)
rate. The token rate is reduced by a multiplicative factor when the service is
overloaded, and increased by an additive factor while the service is no longer
overloaded.

At a high level, this policy works as follows:

- Latency trend-based overload detection: A Flux Meter is used to gather latency
  metrics from a [service control point](/concepts/selector.md). The historical
  latency over a large time window (30 minutes by default) is used to establish
  a long-term trend that can be compared to the current latency to detect
  overloads.
- Gradient Controller: Set point latency and current latency signals are fed to
  the gradient controller that calculates the proportional response to adjust
  the accepted token rate (Control Variable).
- Integral Optimizer: When the service is detected to be in the normal state, an
  integral optimizer is used to additively increase the accepted token rate of
  the service in each execution cycle of the circuit. This measured approach
  prevents accepting all the traffic at once after an overload, which can again
  lead to an overload.
- Load Scheduler: The accepted token rate at the service is throttled by a
  [weighted-fair queuing scheduler](/concepts/scheduler/scheduler.md). The
  output of the adjustments to accepted token rate made by gradient controller
  and optimizer logic are translated to a load multiplier that is synchronized
  with Aperture Agents through etcd. The load multiplier adjusts (increases or
  decreases) the token bucket fill rates based on the incoming token rate
  observed at each agent.

:::info

Please see reference for the
[`AdaptiveLoadScheduler`](/reference/configuration/spec.md#adaptive-load-scheduler)
component that is used within this blueprint.

:::

:::info

See the use-cases
[Adaptive Service Protection with Average Latency Feedback](/use-cases/adaptive-service-protection/average-latency-feedback.md)
and
[Workload Prioritization](/use-cases/adaptive-service-protection/workload-prioritization.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/load-scheduling/average-latency`}>load-scheduling/average-latency</a>

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
    reference='../../spec#component'
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
    reference='../../spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-evaluation-interval"></a>

<ParameterDescription
    name='policy.evaluation_interval'
    description='The interval between successive evaluations of the Circuit.'
    type='string'
    reference=''
    value='"10s"'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.service_protection_core {#policy-service-protection-core}

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler'
    description='Parameters for Adaptive Load Scheduler.'
    type='Object (aperture.spec.v1.AdaptiveLoadSchedulerParameters)'
    reference='../../spec#adaptive-load-scheduler-parameters'
    value='{"alerter": {"alert_name": "Load Throttling Event"}, "gradient": {"max_gradient": 1, "min_gradient": 0.1, "slope": -1}, "load_multiplier_linear_increment": 0.025, "load_scheduler": {"selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "max_load_multiplier": 2}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-dry-run"></a>

<ParameterDescription
    name='policy.service_protection_core.dry_run'
    description='Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-kubelet-overload-confirmations"></a>

<ParameterDescription
    name='policy.service_protection_core.kubelet_overload_confirmations'
    description='Overload confirmation signals from kubelet.'
    type='Object (kubelet_overload_confirmations)'
    reference='#kubelet-overload-confirmations'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-overload-confirmations"></a>

<ParameterDescription
    name='policy.service_protection_core.overload_confirmations'
    description='List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.'
    type='Array of Object (overload_confirmation)'
    reference='#overload-confirmation'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.latency_baseliner {#policy-latency-baseliner}

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-flux-meter"></a>

<ParameterDescription
    name='policy.latency_baseliner.flux_meter'
    description='Flux Meter defines the scope of latency measurements.'
    type='Object (aperture.spec.v1.FluxMeter)'
    reference='../../spec#flux-meter'
    value='{"selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-latency-tolerance-multiplier"></a>

<ParameterDescription
    name='policy.latency_baseliner.latency_tolerance_multiplier'
    description='Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if the long-term average of latency is L and if the tolerance is T, then the service is considered to be in an overloaded state if the short-term average of latency is more than L*T.'
    type='Number (double)'
    reference=''
    value='1.25'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-long-term-query-interval"></a>

<ParameterDescription
    name='policy.latency_baseliner.long_term_query_interval'
    description='Interval for long-term latency query, i.e., how far back in time the query is run. The value should be a string representing the duration in seconds.'
    type='string'
    reference=''
    value='"1800s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-long-term-query-periodic-interval"></a>

<ParameterDescription
    name='policy.latency_baseliner.long_term_query_periodic_interval'
    description='Periodic interval for long-term latency query, i.e., how often the query is run. The value should be a string representing the duration in seconds.'
    type='string'
    reference=''
    value='"30s"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-extra-filters"></a>

<ParameterDescription
    name='dashboard.extra_filters'
    description='Additional filters to pass to each query to Grafana datasource.'
    type='Object (map[string]string)'
    reference='#map-string-string'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a>

<ParameterDescription
    name='dashboard.refresh_interval'
    description='Refresh interval for dashboard panels.'
    type='string'
    reference=''
    value='"15s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-from"></a>

<ParameterDescription
    name='dashboard.time_from'
    description='From time of dashboard.'
    type='string'
    reference=''
    value='"now-15m"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a>

<ParameterDescription
    name='dashboard.time_to'
    description='To time of dashboard.'
    type='string'
    reference=''
    value='"now"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-title"></a>

<ParameterDescription
    name='dashboard.title'
    description='Name of the main dashboard.'
    type='string'
    reference=''
    value='"Aperture Service Protection"'
/>

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-filter-regex"></a>

<ParameterDescription
    name='dashboard.datasource.filter_regex'
    description='Datasource filter regex.'
    type='string'
    reference=''
    value='""'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-name"></a>

<ParameterDescription
    name='dashboard.datasource.name'
    description='Datasource name.'
    type='string'
    reference=''
    value='"$datasource"'
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
    description='The driver for using CPU usage as overload confirmation.'
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
    reference='../../spec#kubernetes-object-selector'
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
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
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
