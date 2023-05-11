local spec = import '../../../spec.libsonnet';
local baseAutoScalingPolicyFn = import '../../auto-scaling/pod-auto-scaler/policy.libsonnet';
local baseServiceProtectionPolicyFn = import '../../service-protection/average-latency/policy.libsonnet';
local config = import './config.libsonnet';

local scaleOutController = spec.v1.ScaleOutController;
local scaleOutControllerController = spec.v1.ScaleOutControllerController;
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
      scaling_backend: {
        kubernetes_replicas: params.policy.auto_scaling.kubernetes_replicas,
      },
    },
  },

  local baseServiceProtectionPolicy = baseServiceProtectionPolicyFn(params).policyDef,
  local baseAutoScalingPolicy = baseAutoScalingPolicyFn(autoScalingParams).policyDef,

  local scaleOutControllers = [
    scaleOutController.new()
    + scaleOutController.withAlerter(
      alerterParameters.new()
      + alerterParameters.withAlertName('Scale Out Alerter')
    )
    + scaleOutController.withController(
      scaleOutControllerController.new()
      + scaleOutControllerController.withGradient(
        increasingGradient.new()
        + increasingGradient.withInPorts(
          increasingGradientInPort.new()
          + increasingGradientInPort.withSignal(port.withSignalName('OBSERVED_LOAD_MULTIPLIER'))
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
