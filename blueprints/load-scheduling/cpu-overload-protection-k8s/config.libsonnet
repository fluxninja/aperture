local commonConfig = import '../common-aiad/config-defaults.libsonnet';

commonConfig {
  policy+: {
    /**
    * @param (policy.kubernetes_object_selector.api_version: string) API version of the object to protect.
    * @param (policy.kubernetes_object_selector.kind: string) Kind of the object to protect.
    * @param (policy.kubernetes_object_selector.name: string) Name of the object to protect.
    * @param (policy.kubernetes_object_selector.namespace: string) Namespace of the object to protect.
    * @param (policy.load_scheduling_core.setpoint: float64) Setpoint.
    */
    kubernetes_object_selector: {
      api_version: 'apps/v1',
      kind: 'Deployment',
      name: '__REQUIRED_FIELD__',
      namespace: '__REQUIRED_FIELD__',
    },
    load_scheduling_core+: {
      setpoint: '__REQUIRED_FIELD__',
    },
  },
}
