local commonConfig = import '../common/config-defaults.libsonnet';

local range_driven_load_scheduler = {
  load_scheduler: {
    selectors: commonConfig.selectors_defaults,
  },
  low_throttle_threshold: '__REQUIRED_FIELD__',
  high_throttle_threshold: '__REQUIRED_FIELD__',
  degree: '__REQUIRED_FIELD__',
  alerter: {
    alert_name: 'Range Driven Load Throttling Event',
  },
};

commonConfig {
  /**
  * @param (policy.service_protection_core.range_driven_load_scheduler: aperture.spec.v1.RangeDrivenLoadSchedulerParameters) Parameters for Range Throttling Strategy.
  */
  policy+: {
    service_protection_core+: {
      range_driven_load_scheduler: range_driven_load_scheduler,
    },
  },

  dashboard+: {
    title: 'Aperture Service Protection',
    variant_name: 'Service Protection',
  },
}
