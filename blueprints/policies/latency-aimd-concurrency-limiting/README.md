[Documentation](https://docs.fluxninja.com/reference/policies/bundled-blueprints/policies/latency-aimd-concurrency-limiting)

<!-- Configuration Marker -->

### Parameters

#### common {#common}

##### common.policy_name {#common-policy-name}

**Type**: string **Default Value**: `__REQUIRED_FIELD__` **Description**: Name
of the policy.

#### policy {#policy}

##### policy.flux_meter {#policy-flux-meter}

**Type**: Object (aperture.spec.v1.FluxMeter) **Default Value**:
`{'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}}`
**Description**: Flux Meter.

##### policy.classifiers {#policy-classifiers}

**Type**: Array of Object (aperture.spec.v1.Classifier) **Default Value**: `[]`
**Description**: List of classification rules.

##### policy.components {#policy-components}

**Type**: Array of Object (aperture.spec.v1.Component) **Default Value**: `[]`
**Description**: List of additional circuit components.

##### policy.latency_baseliner {#policy-latency-baseliner}

###### policy.latency_baseliner.ema {#policy-latency-baseliner-ema}

**Type**: Object (aperture.spec.v1.EMAParameters) **Default Value**:
`{'correction_factor_on_max_envelope_violation': 0.95, 'ema_window': '1500s', 'warmup_window': '60s'}`
**Description**: EMA parameters.

###### policy.latency_baseliner.latency_tolerance_multiplier {#policy-latency-baseliner-latency-tolerance-multiplier}

**Type**: Number (double) **Default Value**: `1.1` **Description**: Tolerance
factor beyond which the service is considered to be in overloaded state. E.g. if
EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be
in overloaded state if current latency is more than 55ms.

###### policy.latency_baseliner.latency_ema_limit_multiplier {#policy-latency-baseliner-latency-ema-limit-multiplier}

**Type**: Number (double) **Default Value**: `2` **Description**: Current
latency value is multiplied with this factor to calculate maximum envelope of
Latency EMA.

##### policy.concurrency_controller {#policy-concurrency-controller}

###### policy.concurrency_controller.flow_selector {#policy-concurrency-controller-flow-selector}

**Type**: Object (aperture.spec.v1.FlowSelector) **Default Value**:
`{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}`
**Description**: Concurrency Limiter flow selector.

###### policy.concurrency_controller.scheduler {#policy-concurrency-controller-scheduler}

**Type**: Object (aperture.spec.v1.SchedulerParameters) **Default Value**:
`{'auto_tokens': True}` **Description**: Scheduler parameters.

###### policy.concurrency_controller.gradient {#policy-concurrency-controller-gradient}

**Type**: Object (aperture.spec.v1.GradientControllerParameters) **Default
Value**: `{'max_gradient': 1, 'min_gradient': 0.1, 'slope': -1}`
**Description**: Gradient Controller parameters.

###### policy.concurrency_controller.alerter {#policy-concurrency-controller-alerter}

**Type**: Object (aperture.spec.v1.AlerterParameters) **Default Value**:
`{'alert_name': 'Load Shed Event'}` **Description**: Whether tokens for
workloads are computed dynamically or set statically by the user.

###### policy.concurrency_controller.max_load_multiplier {#policy-concurrency-controller-max-load-multiplier}

**Type**: Number (double) **Default Value**: `2` **Description**: Current
accepted concurrency is multiplied with this number to dynamically calculate the
upper concurrency limit of a Service during normal (non-overload) state. This
protects the Service from sudden spikes.

###### policy.concurrency_controller.load_multiplier_linear_increment {#policy-concurrency-controller-load-multiplier-linear-increment}

**Type**: Number (double) **Default Value**: `0.0025` **Description**: Linear
increment to load multiplier in each execution tick (0.5s) when the system is
not in overloaded state.

###### policy.concurrency_controller.default_config {#policy-concurrency-controller-default-config}

**Type**: Object (aperture.spec.v1.LoadActuatorDynamicConfig) **Default Value**:
`{'dry_run': False}` **Description**: Default configuration for concurrency
controller that can be updated at the runtime without shutting down the policy.

##### policy.rate_limiter {#policy-rate-limiter}

**Type**: Object (policies/static-rate-limiting:policy.rate_limiter) **Default
Value**:
`{'default_config': {'overrides': []}, 'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'parameters': {'label_key': '__REQUIRED_FIELD__', 'limit_reset_interval': '__REQUIRED_FIELD__'}, 'rate_limit': '__REQUIRED_FIELD__'}`
**Description**: Rate Limiter parameters.

#### dashboard {#dashboard}

##### dashboard.refresh_interval {#dashboard-refresh-interval}

**Type**: string **Default Value**: `'5s'` **Description**: Refresh interval for
dashboard panels.

##### dashboard.time_from {#dashboard-time-from}

**Type**: string **Default Value**: `'now-15m'` **Description**: From time of
dashboard.

##### dashboard.time_to {#dashboard-time-to}

**Type**: string **Default Value**: `'now'` **Description**: To time of
dashboard.

##### dashboard.datasource {#dashboard-datasource}

###### dashboard.datasource.name {#dashboard-datasource-name}

**Type**: string **Default Value**: `'$datasource'` **Description**: Datasource
name.

###### dashboard.datasource.filter_regex {#dashboard-datasource-filter-regex}

**Type**: string **Default Value**: `''` **Description**: Datasource filter
regex.

## Dynamic Configuration

### Parameters

#### concurrency_controller {#concurrency-controller}

**Type**: Object (aperture.spec.v1.LoadActuatorDynamicConfig) **Default Value**:
`__REQUIRED_FIELD__` **Description**: Default configuration for concurrency
controller that can be updated at the runtime without shutting down the policy.
