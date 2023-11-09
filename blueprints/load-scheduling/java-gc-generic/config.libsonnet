local serviceProtectionDefaults = import '../common-aiad/config-defaults.libsonnet';

serviceProtectionDefaults {
  policy+: {
    jmx: {
      /**
      * @param (policy.jmx.jmx_host: string) Hostname for scraping metrics provided by JMX Prometheus Java Agent.
      * @param (policy.jmx.jmx_prometheus_port: int32) Port number for scraping metrics provided by JMX Prometheus Java Agent.
      */
      jmx_host: '__REQUIRED_FIELD__',
      jmx_prometheus_port: '__REQUIRED_FIELD__',
    },

    /**
    * @param (policy.load_scheduling_core.setpoint: float64) Setpoint.
    */
    load_scheduling_core+: {
      setpoint: '__REQUIRED_FIELD__',
    },
  },
}
