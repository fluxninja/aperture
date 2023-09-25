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
    * @param (policy.service_protection_core.overload_condition: string)Overload condition determines the criteria to determine overload state. The default condition is 'gt', that is, when the signal is greater than the setpoint. The condition must be one of: gt, lt, gte, lte.
    */
    service_protection_core+: {
      overload_condition: 'gt',
    },
  },

  dashboard+: {
    variant_name: 'PromQL Output',
  },
}
