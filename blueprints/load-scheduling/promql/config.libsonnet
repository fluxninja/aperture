local serviceProtectionDefaults = import '../common-aiad/config-defaults.libsonnet';

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
    /**
    * @param (policy.overload_condition: string)Overload condition determines the criteria to determine overload state. The condition must be one of: gt, lt, gte, lte.
    */
    overload_condition: '__REQUIRED_FIELD__',
  },
}
