local serviceProtectionDefaults = import '../common/config-defaults.libsonnet';


serviceProtectionDefaults {
  policy+: {
    /**
    * @param (policy.promql_query: string) PromQL query.
    */
    promql_query: 'avg(java_lang_G1_Young_Generation_LastGcInfo_duration{k8s_pod_name=~"service3-demo-app-.*"})',
    /**
    * @param (policy.setpoint: float64) Setpoint.
    */
    setpoint: 20,

    jmx: {
      /**
      * @param (policy.jmx.jmx_metrics_port: int32) Port number for scraping metrics provided by JMX Promtheus Java Agent.
      * @param (policy.jmx.app_namespace: string) Namespace of the application for which JMX metrics are scraped.
      */
      jmx_metrics_port: 8087,
      app_namespace: '__REQUIRED_FIELD__',
    },
  },

  dashboard+: {
    variant_name: 'Protection with JMX Overload Confirmation',
  },
}
