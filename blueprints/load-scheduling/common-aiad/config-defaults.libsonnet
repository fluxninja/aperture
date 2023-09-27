local commonConfig = import '../common/config-defaults.libsonnet';

local aiad_load_scheduler = {
  load_scheduler: {
    selectors: commonConfig.selectors_defaults,
  },
  min_load_multiplier: 0.0,
  max_load_multiplier: 2.0,
  load_multiplier_linear_increment: 0.025,
  load_multiplier_linear_decrement: 0.05,
  alerter: {
    alert_name: 'AIAD Load Throttling Event',
  },
};

commonConfig {
  /**
  * @param (policy.service_protection_core.aiad_load_scheduler: aperture.spec.v1.AIADLoadSchedulerParameters) Parameters for AIMD throttling strategy.
  */
  policy+: {
    service_protection_core+: {
      aiad_load_scheduler: aiad_load_scheduler,
    },
  },

  dashboard+: {
    title: 'Aperture Service Protection',
    variant_name: 'Service Protection',
  },
}
