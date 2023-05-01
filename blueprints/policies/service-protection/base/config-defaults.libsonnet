local flow_selector_defaults = {
  service_selector: {
    service: '__REQUIRED_FIELD__',
  },
  flow_matcher: {
    control_point: '__REQUIRED_FIELD__',
  },
};

local overload_confirmation_defaults = {
  query_string: '__REQUIRED_FIELD__',
  threshold: '__REQUIRED_FIELD__',
  operator: '__REQUIRED_FIELD__',
};

local service_protection_core_defaults = {

  overload_confirmations: [overload_confirmation_defaults],

  adaptive_load_scheduler: {
    flow_selector: flow_selector_defaults,
    scheduler: {
      auto_tokens: true,
    },
    gradient: {
      slope: -1,
      min_gradient: 0.1,
      max_gradient: 1.0,
    },
    alerter: {
      alert_name: 'Load Shed Event',
    },
    max_load_multiplier: 2.0,
    load_multiplier_linear_increment: 0.0025,
    default_config: {
      dry_run: false,
    },
  },
};

{
  common: {
    policy_name: '__REQUIRED_FIELD__',
  },
  policy: {
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
  flow_selector: flow_selector_defaults,
  overload_confirmation: overload_confirmation_defaults,
}
