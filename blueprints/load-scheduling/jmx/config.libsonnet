local serviceProtectionDefaults = import '../common-aiad/config-defaults.libsonnet';

serviceProtectionDefaults {
  policy+: {
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
