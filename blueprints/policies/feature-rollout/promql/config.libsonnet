local featureRolloutConfig = import '../base/config-defaults.libsonnet';

/**
* @schema (criteria.forward.operator: string) The operator for the forward criteria. oneof: `gt | lt | gte | lte | eq | neq`
* @schema (criteria.backward.operator: string) The operator for the backward criteria. oneof: `gt | lt | gte | lte | eq | neq`
* @schema (criteria.reset.operator: string) The operator for the reset criteria. oneof: `gt | lt | gte | lte | eq | neq`
*/
local criteria_defaults = featureRolloutConfig.criteria_defaults {
  forward: {
    operator: '__REQUIRED_FIELD__',
  },
};

featureRolloutConfig {
  policy+: {
    drivers+: {
      /**
      * @param (policy.drivers.promql_drivers: []promql_driver) List of promql drivers that compare results of a Prometheus query against forward, backward and reset thresholds.
      * @schema (promql_driver.query_string: string) The Prometheus query to be run. Must return a scalar or a vector with a single element.
      * @schema (promql_driver.criteria: criteria) The criteria for comparing results of the Prometheus query.
      */
      promql_drivers: [
        {
          query_string: '__REQUIRED_FIELD__',
          criteria: criteria_defaults,
        },
      ],
    },
  },
}
