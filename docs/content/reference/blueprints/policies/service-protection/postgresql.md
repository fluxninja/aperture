---
title: Service Protection for PostgreSQL
keywords:
  - blueprints
sidebar_position: 4
sidebar_label: Service Protection for PostgreSQL
---

## Introduction

This policy detects traffic overloads and cascading failure build-up on
PostgreSQL by checking the real-time percentage of PostgreSQL connections
against the maximum number of connections.

It also uses the CPU utilization ratio of the PostgreSQL pod to confirm traffic
overloads and cascading failure build-up. The CPU utilization ratio is the
percentage of CPU used by the PostgreSQL pod divided by the total CPU available
to the pod.

All the PostgreSQL related metrics are collected by the
[PostgreSQL OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/postgresqlreceiver)
so if the system under observation requires using different metrics for the
overload confirmation, the
[list of available metrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/postgresqlreceiver/metadata.yaml)
can be used to configure the policy.

A gradient controller calculates a proportional response to limit accepted
concurrency. The concurrency is reduced by a multiplicative factor when the
service is overloaded, and increased by an additive factor while the service is
no longer overloaded.

:::info

Please see reference for the
[`AdaptiveLoadScheduler`](/reference/configuration/spec.md#adaptive-load-scheduler)
component that is used within this blueprint.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection/postgresql`}>policies/service-protection/postgresql</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### policy {#policy}

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

<a id="policy-promql-query"></a>

<ParameterDescription
    name='policy.promql_query'
    description='PromQL query to detect PostgreSQL overload.'
    type='string'
    reference=''
    value='"(sum(postgresql_backends) / sum(postgresql_connection_max)) * 100"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-components"></a>

<ParameterDescription
    name='policy.components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../../../spec#component'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-resources"></a>

<ParameterDescription
    name='policy.resources'
    description='Additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../../../spec#resources'
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
    value='"1s"'
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

<a id="policy-postgresql"></a>

<ParameterDescription
    name='policy.postgresql'
    description='Configuration for PostgreSQL OpenTelemetry receiver. Refer https://docs.fluxninja.com/integrations/metrics/postgresql for more information.'
    type='Object (postgresql)'
    reference='#postgresql'
    value='{"agent_group": "default", "endpoint": "__REQUIRED_FIELD__", "password": "__REQUIRED_FIELD__", "username": "__REQUIRED_FIELD__"}'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.service_protection_core {#policy-service-protection-core}

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

<a id="policy-service-protection-core-adaptive-load-scheduler"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler'
    description='Parameters for Adaptive Load Scheduler.'
    type='Object (aperture.spec.v1.AdaptiveLoadSchedulerParameters)'
    reference='../../../spec#adaptive-load-scheduler-parameters'
    value='{"alerter": {"alert_name": "Load Throttling Event"}, "gradient": {"max_gradient": 1, "min_gradient": 0.1, "slope": -1}, "load_multiplier_linear_increment": 0.0025, "load_scheduler": {"selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "max_load_multiplier": 2}'
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

###### policy.service_protection_core.cpu_overload_confirmation {#policy-service-protection-core-cpu-overload-confirmation}

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-cpu-overload-confirmation-query-string"></a>

<ParameterDescription
    name='policy.service_protection_core.cpu_overload_confirmation.query_string'
    description='The Prometheus query to be run to get the PostgreSQL CPU utilization. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='"avg(k8s_pod_cpu_utilization_ratio{k8s_statefulset_name=\"__REQUIRED_FIELD__\"})"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-cpu-overload-confirmation-threshold"></a>

<ParameterDescription
    name='policy.service_protection_core.cpu_overload_confirmation.threshold'
    description='Threshold value for CPU utilizatio if it has to be used as overload confirmation.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-cpu-overload-confirmation-operator"></a>

<ParameterDescription
    name='policy.service_protection_core.cpu_overload_confirmation.operator'
    description='The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`.'
    type='string'
    reference=''
    value='"gte"'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.auto_scaling {#policy-auto-scaling}

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-dry-run"></a>

<ParameterDescription
    name='policy.auto_scaling.dry_run'
    description='Dry run mode ensures that no scaling is invoked by the auto scaler escalation. This config can be updated at runtime without restarting the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-promql-scale-out-controllers"></a>

<ParameterDescription
    name='policy.auto_scaling.promql_scale_out_controllers'
    description='List of scale out controllers.'
    type='Array of Object (promql_scale_out_controller)'
    reference='#promql-scale-out-controller'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-promql-scale-in-controllers"></a>

<ParameterDescription
    name='policy.auto_scaling.promql_scale_in_controllers'
    description='List of scale in controllers.'
    type='Array of Object (promql_scale_in_controller)'
    reference='#promql-scale-in-controller'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-scaling-parameters"></a>

<ParameterDescription
    name='policy.auto_scaling.scaling_parameters'
    description='Parameters that define the scaling behavior.'
    type='Object (aperture.spec.v1.AutoScalerScalingParameters)'
    reference='../../../spec#auto-scaler-scaling-parameters'
    value='{"scale_in_alerter": {"alert_name": "Auto-scaler is scaling in"}, "scale_in_cooldown": "40s", "scale_out_alerter": {"alert_name": "Auto-scaler is scaling out"}, "scale_out_cooldown": "30s"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-scaling-backend"></a>

<ParameterDescription
    name='policy.auto_scaling.scaling_backend'
    description='Scaling backend for the policy.'
    type='Object (aperture.spec.v1.AutoScalerScalingBackend)'
    reference='../../../spec#auto-scaler-scaling-backend'
    value='{"kubernetes_replicas": {"kubernetes_object_selector": "__REQUIRED_FIELD__", "max_replicas": "__REQUIRED_FIELD__", "min_replicas": "__REQUIRED_FIELD__"}}'
/>

<!-- vale on -->

<!-- vale off -->

###### policy.auto_scaling.periodic_decrease {#policy-auto-scaling-periodic-decrease}

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-periodic-decrease-period"></a>

<ParameterDescription
    name='policy.auto_scaling.periodic_decrease.period'
    description='Period for periodic scale in.'
    type='string'
    reference=''
    value='"60s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-periodic-decrease-scale-in-percentage"></a>

<ParameterDescription
    name='policy.auto_scaling.periodic_decrease.scale_in_percentage'
    description='Percentage of replicas to scale in.'
    type='Number (double)'
    reference=''
    value='10'
/>

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

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
    description='Time from of dashboard.'
    type='string'
    reference=''
    value='"now-15m"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a>

<ParameterDescription
    name='dashboard.time_to'
    description='Time to of dashboard.'
    type='string'
    reference=''
    value='"now"'
/>

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

<a id="dashboard-title"></a>

<ParameterDescription
    name='dashboard.title'
    description='Name of the main dashboard.'
    type='string'
    reference=''
    value='"Aperture Service Protection for PostgreSQL"'
/>

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

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

---

### Schemas

<!-- vale off -->

#### overload_confirmation {#overload-confirmation}

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

---

<!-- vale off -->

#### promql_scale_out_controller {#promql-scale-out-controller}

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-threshold"></a>

<ParameterDescription
    name='threshold'
    description='Threshold for the controller.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-gradient"></a>

<ParameterDescription
    name='gradient'
    description='Gradient parameters for the controller.'
    type='Object (aperture.spec.v1.IncreasingGradientParameters)'
    reference='../../../spec#increasing-gradient-parameters'
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-alerter"></a>

<ParameterDescription
    name='alerter'
    description='Alerter parameters for the controller.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../../spec#alerter-parameters'
    value='null'
/>

<!-- vale on -->

---

<!-- vale off -->

#### promql_scale_in_controller {#promql-scale-in-controller}

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-threshold"></a>

<ParameterDescription
    name='threshold'
    description='Threshold for the controller.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-gradient"></a>

<ParameterDescription
    name='gradient'
    description='Gradient parameters for the controller.'
    type='Object (aperture.spec.v1.DecreasingGradientParameters)'
    reference='../../../spec#decreasing-gradient-parameters'
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-alerter"></a>

<ParameterDescription
    name='alerter'
    description='Alerter parameters for the controller.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../../spec#alerter-parameters'
    value='null'
/>

<!-- vale on -->

---

<!-- vale off -->

#### postgresql {#postgresql}

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-username"></a>

<ParameterDescription
    name='username'
    description='Username of the PostgreSQL.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-password"></a>

<ParameterDescription
    name='password'
    description='Password of the PostgreSQL.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-endpoint"></a>

<ParameterDescription
    name='endpoint'
    description='Endpoint of the PostgreSQL.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-transport"></a>

<ParameterDescription
    name='transport'
    description='The transport protocol being used to connect to postgresql. Available options are tcp and unix.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-database"></a>

<ParameterDescription
    name='database'
    description='The list of databases for which the receiver will attempt to collect statistics.'
    type='Array of string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-collection-interval"></a>

<ParameterDescription
    name='collection_interval'
    description='This receiver collects metrics on an interval.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-initial-delay"></a>

<ParameterDescription
    name='initial_delay'
    description='Defines how long this receiver waits before starting.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-agent-group"></a>

<ParameterDescription
    name='agent_group'
    description='Name of the Aperture Agent group.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

##### tls {#postgresql-tls}

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-tls-insecure"></a>

<ParameterDescription
    name='insecure'
    description='Whether to enable client transport security for the postgresql connection.'
    type='Boolean'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-tls-insecure-skip-verify"></a>

<ParameterDescription
    name='insecure_skip_verify'
    description='Whether to validate server name and certificate if client transport security is enabled.'
    type='Boolean'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-tls-cert-file"></a>

<ParameterDescription
    name='cert_file'
    description='A cerficate used for client authentication, if necessary.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-tls-key-file"></a>

<ParameterDescription
    name='key_file'
    description='An SSL key used for client authentication, if necessary.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="postgresql-tls-ca-file"></a>

<ParameterDescription
    name='ca_file'
    description='A set of certificate authorities used to validate the database server SSL certificate.'
    type='string'
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
