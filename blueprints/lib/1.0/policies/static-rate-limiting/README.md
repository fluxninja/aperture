# Static Rate Limiting Policy

This blueprint provides a simple static rate limiting policy and a dashboard.

## Configuration

<!-- Configuration Marker -->

export const ParameterHeading = ({children}) => ( <span
style={{fontWeight: "bold"}}>{children}</span> );

export const WrappedDescription = ({children}) => ( <span
style={{wordWrap: "normal"}}>{children}</span> );

export const ParameterDescription = ({name, type, value, description}) => (

  <table>
  <tr>
    <td><ParameterHeading>Parameter</ParameterHeading></td>
    <td><code>{name}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Type</ParameterHeading></td>
    <td><code>{type}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Default Value</ParameterHeading></td>
    <td><code>{value ? value : "REQUIRED VALUE"}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Description</ParameterHeading></td>
    <td><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);
);

### Common

<ParameterDescription name="common.policy_name" type="string"

    description='Name of the policy.' />

### Policy

<ParameterDescription name="policy.evaluation_interval" type="string"
value=""300s"" description='How often should the policy be re-evaluated' />

<ParameterDescription
    name="policy.classifiers"
    type="[]aperture.spec.v1.Classifier"
    value="[]"
    description='List of classification rules.' />

#### Rate Limiter

<ParameterDescription name="policy.rate_limiter.rate_limit" type="float64"

    description='Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to accept' />

<ParameterDescription name="policy.rate_limiter.flow_selector"
type="aperture.spec.v1.FlowSelector"

    description='A flow selector to match requests against' />

<ParameterDescription
    name="policy.rate_limiter.parameters"
    type="aperture.spec.v1.RateLimiterParameters"
    value="{'label_key': 'FAKE-VALUE', 'lazy_sync': {'enabled': True, 'num_sync': 5}, 'limit_reset_interval': '1s'}"
    description='Parameters.' />

<ParameterDescription name="policy.rate_limiter.parameters.label_key"
type="string"

    description='Flow label to use for rate limiting.' />

<ParameterDescription
    name="policy.rate_limiter.dynamic_config"
    type="aperture.spec.v1.RateLimiterDefaultConfig"
    value="{'overrides': []}"
    description='Dynamic configuration for rate limiter that can be applied at the runtime.' />

### Dashboard

<ParameterDescription name="dashboard.refresh_interval" type="string"
value=""10s"" description='Refresh interval for dashboard panels.' />

#### Datasource

<ParameterDescription name="dashboard.datasource.name" type="string"
value=""$datasource"" description='Datasource name.' />

<ParameterDescription name="dashboard.datasource.filter_regex" type="string"
value="""" description='Datasource filter regex.' />
