local selectors_defaults = [{
  service: '__REQUIRED_FIELD__',
  control_point: '__REQUIRED_FIELD__',
}];

local overload_confirmation_defaults = {
  query_string: '__REQUIRED_FIELD__',
  threshold: '__REQUIRED_FIELD__',
  operator: '__REQUIRED_FIELD__',
};

local service_protection_core_defaults = {
  overload_confirmations: [overload_confirmation_defaults],

  adaptive_load_scheduler: {
    load_scheduler: {
      selectors: selectors_defaults,
    },
    gradient: {
      slope: -1,
      min_gradient: 0.1,
      max_gradient: 1.0,
    },
    alerter: {
      alert_name: 'Load Throttling Event',
    },
    max_load_multiplier: 2.0,
    load_multiplier_linear_increment: 0.0025,
  },

  dry_run: false,
};

{
  policy: {
    policy_name: '__REQUIRED_FIELD__',
    components: [],
    resources: {
      flow_control: {
        classifiers: [],
      },
    },
    evaluation_interval: '1s',
    service_protection_core: service_protection_core_defaults,
  },

  service_protection_core: service_protection_core_defaults,

  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
  },

  selectors: selectors_defaults,

  overload_confirmation: overload_confirmation_defaults,
}
