local commonConfig = import '../common/config-defaults.libsonnet';

local aimd_load_scheduler = {
  load_scheduler: {
    selectors: commonConfig.selectors_defaults,
  },
  gradient: {
    slope: -1,
    min_gradient: 0.1,
    max_gradient: 1.0,
  },
  max_load_multiplier: 2.0,
  load_multiplier_linear_increment: 0.025,
  alerter: {
    alert_name: 'AIMD Load Throttling Event',
  },
};

commonConfig {
  /**
  * @param (policy.load_scheduling_core.aimd_load_scheduler: aperture.spec.v1.AIMDLoadSchedulerParameters) Parameters for AIMD throttling strategy.
  */
  policy+: {
    load_scheduling_core+: {
      aimd_load_scheduler: aimd_load_scheduler,
    },
  },
}
