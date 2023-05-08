{
  new():: {
  },
  withAlerter(alerter):: {
    alerter: alerter,
  },
  withAlerterMixin(alerter):: {
    alerter+: alerter,
  },
  withController(controller):: {
    controller: controller,
  },
  withControllerMixin(controller):: {
    controller+: controller,
  },
}
