local serviceProtectionDefaults = import '../common-range/config-defaults.libsonnet';


serviceProtectionDefaults {
  policy+: {
    /**
    * @param (policy.promql_query: string) PromQL query.
    */
    promql_query: '__REQUIRED_FIELD__',
  },

  dashboard+: {
    variant_name: 'PromQL Output',
  },
}
