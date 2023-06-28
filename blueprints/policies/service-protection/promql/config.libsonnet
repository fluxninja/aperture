local serviceProtectionDefaults = import '../common/config-defaults.libsonnet';


serviceProtectionDefaults {
  policy+: {
    /**
    * @param (policy.promql_query: string) PromQL query.
    */
    promql_query: '__REQUIRED_FIELD__',
    /**
    * @param (policy.setpoint: float64) Setpoint.
    */
    setpoint: '__REQUIRED_FIELD__',
  },

  dashboard+: {
    variant_name: 'PromQL Output',
  },
}
