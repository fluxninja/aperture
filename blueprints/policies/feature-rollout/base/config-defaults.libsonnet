local selectors_defaults = [{
  service: '__REQUIRED_FIELD__',
  control_point: '__REQUIRED_FIELD__',
}];

/**
* @schema (criteria.forward.threshold: float64) The threshold for the forward criteria.
* @schema (criteria.backward.threshold: float64) The threshold for the backward criteria.
* @schema (criteria.reset.threshold: float64) The threshold for the reset criteria.
*/
local criteria_defaults = {
  forward: {
    threshold: '__REQUIRED_FIELD__',
  },
};

local rollout_policy_base_defaults = {
  policy_name: '__REQUIRED_FIELD__',
  load_ramp: {
    sampler: {
      selectors: selectors_defaults,
      label_key: '',
    },
    steps: [
      {
        target_accept_percentage: '__REQUIRED_FIELD__',
        duration: '__REQUIRED_FIELD__',
      },
    ],
  },

  drivers: {},

  rollout: false,

  components: [],

  resources: {
    flow_control: {
      classifiers: [],
    },
  },

  evaluation_interval: '10s',
};

{
  /**
  * @param (policy.policy_name: string) Name of the policy.
  * @param (policy.load_ramp: aperture.spec.v1.LoadRampParameters) Identify the service and flows of the feature that needs to be rolled out. And specify feature rollout steps.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  * @param (policy.resources: aperture.spec.v1.Resources) List of additional resources.
  * @param (policy.evaluation_interval: string) The interval between successive evaluations of the Circuit.
  * @param (policy.rollout: bool) Whether to start the rollout. This setting may be overridden at runtime via dynamic configuration.
  */
  policy: rollout_policy_base_defaults,

  /**
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  * @param (dashboard.time_from: string) From time of dashboard.
  * @param (dashboard.time_to: string) To time of dashboard.
  * @param (dashboard.extra_filters: map[string]string) Additional filters to pass to each query to Grafana datasource.
  * @param (dashboard.title: string) Name of the main dashboard.
  */
  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    extra_filters: {},
    title: 'Aperture Feature Rollout',
    /**
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
  },

  selectors_defaults: selectors_defaults,

  criteria_defaults: criteria_defaults,
}
