local spec = import '../../../spec.libsonnet';
local utils = import '../../../utils/utils.libsonnet';
local config = import './config-defaults.libsonnet';

local component = spec.v1.Component;
local autoScale = spec.v1.AutoScale;
local autoScaler = spec.v1.AutoScaler;
local scaleInController = spec.v1.ScaleInController;
local scaleInControllerController = spec.v1.ScaleInControllerController;
local scaleOutController = spec.v1.ScaleOutController;
local scaleOutControllerController = spec.v1.ScaleOutControllerController;
local increasingGradient = spec.v1.IncreasingGradient;
local increasingGradientInPort = spec.v1.IncreasingGradientIns;
local decreasingGradient = spec.v1.DecreasingGradient;
local decreasingGradientInPort = spec.v1.DecreasingGradientIns;
local port = spec.v1.Port;
local query = spec.v1.Query;
local promQL = spec.v1.PromQL;
local alerterParameters = spec.v1.AlerterParameters;

function(cfg) {
  local params = config + cfg,

  local scale_in_controllers = [
    scaleInController.new()
    + scaleInController.withAlerter(params.policy.promql_scale_in_controllers[controller_idx].alerter)
    + scaleInController.withController(
      scaleInControllerController.new()
      + scaleInControllerController.withGradient(
        decreasingGradient.new()
        + decreasingGradient.withInPorts(
          decreasingGradientInPort.new()
          + decreasingGradientInPort.withSignal(port.withSignalName('PROMQL_SCALE_IN_%s' % controller_idx))
          + decreasingGradientInPort.withSetpoint(port.withConstantSignal(params.policy.promql_scale_in_controllers[controller_idx].setpoint))
        )
        + decreasingGradient.withParameters(params.policy.promql_scale_in_controllers[controller_idx].gradient)
      )
    )
    for controller_idx in std.range(0, std.length(params.policy.promql_scale_in_controllers) - 1)
  ],

  local scale_out_controllers = [
    scaleOutController.new()
    + scaleOutController.withAlerter(params.policy.promql_scale_out_controllers[controller_idx].alerter)
    + scaleOutController.withController(
      scaleOutControllerController.new()
      + scaleOutControllerController.withGradient(
        increasingGradient.new()
        + increasingGradient.withInPorts(
          increasingGradientInPort.new()
          + increasingGradientInPort.withSignal(port.withSignalName('PROMQL_SCALE_OUT_%s' % controller_idx))
          + increasingGradientInPort.withSetpoint(port.withConstantSignal(params.policy.promql_scale_out_controllers[controller_idx].setpoint))
        )
        + increasingGradient.withParameters(params.policy.promql_scale_out_controllers[controller_idx].gradient)
      )
    )
    for controller_idx in std.range(0, std.length(params.policy.promql_scale_out_controllers) - 1)
  ],

  local scale_in_controllers_promql = [
    component.withQuery(
      query.new()
      + query.withPromql(
        local q = params.policy.promql_scale_in_controllers[controller_idx].query_string;
        promQL.new()
        + promQL.withQueryString(q)
        + promQL.withEvaluationInterval(evaluation_interval=params.policy.evaluation_interval)
        + promQL.withOutPorts({ output: port.withSignalName('PROMQL_SCALE_IN_%s' % controller_idx) }),
      ),
    )
    for controller_idx in std.range(0, std.length(params.policy.promql_scale_in_controllers) - 1)
  ],

  local scale_out_controllers_promql = [
    component.withQuery(
      query.new()
      + query.withPromql(
        local q = params.policy.promql_scale_out_controllers[controller_idx].query_string;
        promQL.new()
        + promQL.withQueryString(q)
        + promQL.withEvaluationInterval(evaluation_interval=params.policy.evaluation_interval)
        + promQL.withOutPorts({ output: port.withSignalName('PROMQL_SCALE_OUT_%s' % controller_idx) }),
      ),
    )
    for controller_idx in std.range(0, std.length(params.policy.promql_scale_out_controllers) - 1)
  ],

  local policyDef =
    spec.v1.Policy.new()
    + spec.v1.Policy.withResources(utils.resources(params.policy.resources).updatedResources)
    + spec.v1.Policy.withCircuit(
      spec.v1.Circuit.new()
      + spec.v1.Circuit.withEvaluationInterval(evaluation_interval=params.policy.evaluation_interval)
      + spec.v1.Circuit.withComponents(
        [
          component.new()
          + component.withAutoScale(
            autoScale.new()
            + autoScale.withAutoScaler(
              autoScaler.new()
              + autoScaler.withDryRunConfigKey('dry_run')
              + autoScaler.withDryRun(params.policy.dry_run)
              + autoScaler.withScalingBackend(params.policy.scaling_backend)
              + autoScaler.withScalingParameters(params.policy.scaling_parameters)
              + autoScaler.withScaleInControllers(scale_in_controllers)
              + autoScaler.withScaleOutControllers(scale_out_controllers)
            )
          ),
        ] + scale_in_controllers_promql + scale_out_controllers_promql + params.policy.components,
      ),
    ),

  policyDef: policyDef,
}
