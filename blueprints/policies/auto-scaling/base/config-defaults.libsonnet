local auto_scaling_base_defaults = {
  policy_name: '__REQUIRED_FIELD__',

  components: [],

  resources: {
    flow_control: {
      classifiers: [],
    },
  },

  evaluation_interval: '10s',
};

local scaling_parameters_defaults = {
  scale_in_alerter: {
    alert_name: 'Auto-scaler is scaling in',
  },
  scale_out_alerter: {
    alert_name: 'Auto-scaler is scaling out',
  },
};

local promql_scale_controller_defaults = {
  query_string: '__REQUIRED_FIELD__',
  setpoint: '__REQUIRED_FIELD__',
  gradient: '__REQUIRED_FIELD__',
  alerter: '__REQUIRED_FIELD__',
};

local auto_scaling_defaults = auto_scaling_base_defaults {
  promql_scale_out_controllers: [promql_scale_controller_defaults],

  promql_scale_in_controllers: [promql_scale_controller_defaults],

  scaling_parameters: scaling_parameters_defaults,

  scaling_backend: '__REQUIRED_FIELD__',

  dry_run: false,
};

{
  policy: auto_scaling_defaults,

  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
  },
}
