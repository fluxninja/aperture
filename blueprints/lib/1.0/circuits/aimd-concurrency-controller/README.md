# AIMD Concurrency Controller Circuit

AIMD (Additive Increase, Multiplicative Decrease) Concurrency Controller adjusts
the concurrency based on `SIGNAL` and `SETPOINT` signals. Gradient controller is
used to calculate a proportional response that limits the accepted concurrency.
Concurrency is increased additively when the overload is no longer detected.

## Configuration

<!-- Configuration Marker -->

### Circuit

| Parameter Name                           | Parameter Type                  | Default      | Description                            |
| ---------------------------------------- | ------------------------------- | ------------ | -------------------------------------- |
| `circuit.concurrencyLimiterFlowSelector` | `aperture.spec.v1.FlowSelector` | `(required)` | Concurrency Limiter flow selector.     |
| `circuit.components`                     | `[]aperture.spec.v1.Component`  | `[]`         | List of additional circuit components. |

#### Concurrency Limiter

| Parameter Name                                                  | Parameter Type                                   | Default                 | Description                                                                                    |
| --------------------------------------------------------------- | ------------------------------------------------ | ----------------------- | ---------------------------------------------------------------------------------------------- |
| `circuit.concurrencyLimiter.autoTokens`                         | `bool`                                           | `true`                  | Whether tokens for workloads are computed dynamically or set statically by the user.           |
| `circuit.concurrencyLimiter.timeoutFactor`                      | `float64`                                        | `0.5`                   | The maximum time a request can wait for tokens as a factor of tokens for a flow in a workload. |
| `circuit.concurrencyLimiter.defaultWorkloadParameters.priority` | `int`                                            | `20`                    | Workload parameters to use in case none of the configured workloads match.                     |
| `circuit.concurrencyLimiter.workloads`                          | `[]aperture.spec.v1.SchedulerParametersWorkload` | `[]`                    | A list of additional workloads for the scheduler.                                              |
| `circuit.concurrencyLimiter.alerterName`                        | `string`                                         | `"Load Shed Event"`     | Name of the alert sent on Load Shed Event.                                                     |
| `circuit.concurrencyLimiter.alerterChannels`                    | `[]string`                                       | `[]`                    | A list of alert channels to which the alert will be sent.                                      |
| `circuit.concurrencyLimiter.alerterResolveTimeout`              | `string`                                         | `"5s"`                  | A timeout after which alert is marked as resolved if alert is not repeated.                    |
| `circuit.concurrencyLimiter.dynamicConfigKey`                   | `string`                                         | `"concurrency_limiter"` | Dynamic configuration key for concurrency limiter.                                             |

#### Constants

| Parameter Name                                         | Parameter Type | Default | Description                                                                                                                                                                                                                                                                                   |
| ------------------------------------------------------ | -------------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `circuit.constants.concurrencyLimitMultiplier`         | `float64`      | `2.0`   | Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.                                                                               |
| `circuit.constants.concurrencyLinearIncrement`         | `float64`      | `5.0`   | Linear increment to concurrency in each execution tick when the system is not in overloaded state.                                                                                                                                                                                            |
| `circuit.constants.concurrencySQRTIncrementMultiplier` | `float64`      | `1`     | Scale factor to multiply square root of current accepted concurrrency. This, along with concurrencyLinearIncrement helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during normal (non-overload) state (integral effect). |

#### Gradient Controller

| Parameter Name                 | Parameter Type | Default | Description                                                                                         |
| ------------------------------ | -------------- | ------- | --------------------------------------------------------------------------------------------------- |
| `circuit.gradient.slope`       | `float64`      | `-1`    | Gradient that adjusts the response of the controller based on current latency and setpoint latency. |
| `circuit.gradient.minGradient` | `float64`      | `0.1`   | Minimum gradient cap.                                                                               |
| `circuit.gradient.maxGradient` | `float64`      | `1.0`   | Maximum gradient cap.                                                                               |
