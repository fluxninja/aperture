# Rate Limiter Blueprint

## Configuration

[configuration]: # Configuration Marker

### Rate Limiter Policy

#### Rate Limiter Overrides

To tweak rate limiter behaviour based on specific label values, a list of
`RateLimierOverride` objects can be added to the policy:

```jsonnet
local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local Override = aperture.spec.v1.RateLimiterOverride;

{
  policy+: {
    overrides: [
      Override.new() + Override.withLabelValue('gold') + Override.withLimitScaleFactor(1)
    ]
  }
}

```

This allows us to prioritize some incoming requests over others.

| Parameter Name              | Parameter Type                           | Default      | Description                                                 |
| --------------------------- | ---------------------------------------- | ------------ | ----------------------------------------------------------- |
| `policy.policyName`         | `string`                                 | `(required)` | An unique name for the policy created by this blueprint     |
| `policy.evaluationInterval` | `string`                                 | `"0.5s"`     | How often should the policy be re-evaluated                 |
| `policy.rateLimit`          | `string`                                 | `(required)` | How many requests per `policy.limitResetInterval` to accept |
| `policy.limitResetInterval` | `string`                                 | `"1s"`       | The window for `policy.rateLimit`                           |
| `policy.labelKey`           | `string`                                 | `(required)` | What flow label to use for rate limiting                    |
| `policy.overrides`          | `[]aperture.spec.v1.RateLimiterOverride` | `[]`         | A list of overrides for the rate limiter                    |

#### Rate Limiter Lazy Sync

| Parameter Name            | Parameter Type | Default | Description                                                          |
| ------------------------- | -------------- | ------- | -------------------------------------------------------------------- |
| `policy.lazySync.enabled` | `boolean`      | `true`  | TODO document what happens when lazy sync is disabled                |
| `policy.lazySync.numSync` | `integer`      | `10`    | Number of times to lazy sync within the `policy.limitResetInterval`. |

#### Selector

| Parameter Name                                      | Parameter Type | Default      | Description                                                                 |
| --------------------------------------------------- | -------------- | ------------ | --------------------------------------------------------------------------- |
| `policy.selector.serviceSelector.agentGroup`        | `string`       | `"default"`  | Which agents to install this policy on                                      |
| `policy.selector.serviceSelector.service`           | `string`       | `(required)` | A fully-qualified domain name of the service that this policy will apply to |
| `policy.selector.flowSelector.controlPoint.traffic` | `string`       | `"ingress"`  | Whether to control `ingress` or `egress` traffic                            |
