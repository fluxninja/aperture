local spec = import '../../../spec.libsonnet';
local baseAutoScalingPolicyFn = import '../../auto-scaling/pod-auto-scaler/policy.libsonnet';
local baseServiceProtectionPolicyFn = import '../../service-protection/average-latency/policy.libsonnet';
local config = import './config.libsonnet';

local scaleOutController = spec.v1.ScaleOutController;
local scaleOutControllerController = spec.v1.ScaleOutControllerController;
local scaleInController = spec.v1.ScaleInController;
local scaleInControllerController = spec.v1.ScaleInControllerController;
local periodicDecrease = spec.v1.PeriodicDecrease;
local periodicDecreaseParameters = spec.v1.PeriodicDecreaseParameters;
local increasingGradient = spec.v1.IncreasingGradient;
local increasingGradientInPort = spec.v1.IncreasingGradientIns;
local increasingGradientParameters = spec.v1.IncreasingGradientParameters;
local alerterParameters = spec.v1.AlerterParameters;
local port = spec.v1.Port;

function(cfg) {
  local params = config + cfg,

  local autoScalingParams = {
    policy+: params.policy.auto_scaling {
      policy_name: params.policy.policy_name,
      promql_scale_out_controllers: [],
      promql_scale_in_controllers: [],
      scaling_backend: {
        kubernetes_replicas: params.policy.auto_scaling.kubernetes_replicas,
      },
    },
  },

  local baseServiceProtectionPolicy = baseServiceProtectionPolicyFn(params).policyDef,
  local baseAutoScalingPolicy = baseAutoScalingPolicyFn(autoScalingParams).policyDef,

  local scaleInControllers = [
    scaleInController.new()
    + scaleInController.withAlerter(
      alerterParameters.new()
      + alerterParameters.withAlertName('Scale in controller Alerter')
    )
    + scaleInController.withController(
      scaleInControllerController.new()
      + scaleInControllerController.withPeriodic(params.policy.auto_scaling.periodic_decrease)
    ),
  ],

  local scaleOutControllers = [
    scaleOutController.new()
    + scaleOutController.withAlerter(
      alerterParameters.new()
      + alerterParameters.withAlertName('Scale out controller Alerter')
    )
    + scaleOutController.withController(
      scaleOutControllerController.new()
      + scaleOutControllerController.withGradient(
        increasingGradient.new()
        + increasingGradient.withInPorts(
          increasingGradientInPort.new()
          + increasingGradientInPort.withSignal(port.withSignalName('DESIRED_LOAD_MULTIPLIER'))
          + increasingGradientInPort.withSetpoint(port.withConstantSignal(1.0))
        )
        + increasingGradient.withParameters(
          increasingGradientParameters.new()
          + increasingGradientParameters.withSlope(-1.0)
        )
      )
    ),
  ],

  local policyDef = baseServiceProtectionPolicy {
    circuit+: {
      components+: std.map(
        function(component) if std.objectHas(component, 'auto_scale') then
          component {
            auto_scale+: {
              auto_scaler+: {
                scale_out_controllers: scaleOutControllers,
                scale_in_controllers: if !params.policy.auto_scaling.disable_periodic_scale_in then scaleInControllers else [],
              },
            },
          }
        else component,
        baseAutoScalingPolicy.circuit.components
      ),
    },
  },

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,

  policyDef: policyDef,
}
