---
title: Load Scheduling for Elasticsearch
keywords:
  - blueprints
sidebar_position: 4
sidebar_label: Load Scheduling for Elasticsearch
---

## Introduction

By default, this policy detects when the Elasticsearch service is overloaded
using the `search` thread pool queue size metric. The policy is based on the
[adaptive load scheduling](/reference/configuration/spec.md#adaptive-load-scheduler)
component.

All the Elasticsearch related metrics are collected by the
[Elasticsearch OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/elasticsearchreceiver)
so if the system under observation requires using different metrics for the
overload confirmation, the
[list of available metrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/elasticsearchreceiver/metadata.yaml)
can be used to configure the policy.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/load-scheduling/elasticsearch`}>load-scheduling/elasticsearch</a>

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

<a id="policy-promql-query"></a>

<ParameterDescription
    name='policy.promql_query'
    description='PromQL query to detect ElasticSearch overload.'
    type='string'
    reference=''
    value='"avg(avg_over_time(elasticsearch_node_thread_pool_tasks_queued{thread_pool_name=\"search\"}[30s]))"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-setpoint"></a>

<ParameterDescription
    name='policy.setpoint'
    description='Setpoint.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-elasticsearch"></a>

<ParameterDescription
    name='policy.elasticsearch'
    description='Configuration for Elasticsearch OpenTelemetry receiver. Refer https://docs.fluxninja.com/integrations/metrics/elasticsearch for more information.'
    type='Object (elasticsearch)'
    reference='#elasticsearch'
    value='{"agent_group": "default", "endpoint": "__REQUIRED_FIELD__", "password": "__REQUIRED_FIELD__", "username": "__REQUIRED_FIELD__"}'
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
    reference='../../configuration/spec#adaptive-load-scheduler-parameters'
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
    value='"Aperture Service Protection for Elasticsearch"'
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

<!-- vale off -->

#### elasticsearch {#elasticsearch}

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-agent-group"></a>

<ParameterDescription
    name='agent_group'
    description='Name of the Aperture Agent group.'
    type='string'
    reference=''
    value='"default"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-collection-interval"></a>

<ParameterDescription
    name='collection_interval'
    description='This receiver collects metrics on an interval.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-endpoint"></a>

<ParameterDescription
    name='endpoint'
    description='Endpoint of the Elasticsearch.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-indices"></a>

<ParameterDescription
    name='indices'
    description='Index filters that define which indices are scraped for index-level metrics.'
    type='Array of string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-initial-delay"></a>

<ParameterDescription
    name='initial_delay'
    description='Defines how long this receiver waits before starting.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-nodes"></a>

<ParameterDescription
    name='nodes'
    description='Node filters that define which nodes are scraped for node-level and cluster-level metrics.'
    type='Array of string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-password"></a>

<ParameterDescription
    name='password'
    description='Password of the Elasticsearch.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-skip-cluster-metrics"></a>

<ParameterDescription
    name='skip_cluster_metrics'
    description='If true, cluster-level metrics will not be scraped.'
    type='Boolean'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="elasticsearch-username"></a>

<ParameterDescription
    name='username'
    description='Username of the Elasticsearch.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
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
