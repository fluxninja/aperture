// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/policy.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public interface GradientControllerOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.GradientController)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Input ports of the Gradient Controller.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.GradientController.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return Whether the inPorts field is set.
   */
  boolean hasInPorts();
  /**
   * <pre>
   * Input ports of the Gradient Controller.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.GradientController.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return The inPorts.
   */
  com.fluxninja.generated.aperture.policy.language.v1.GradientController.Ins getInPorts();
  /**
   * <pre>
   * Input ports of the Gradient Controller.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.GradientController.Ins in_ports = 1 [json_name = "inPorts"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.GradientController.InsOrBuilder getInPortsOrBuilder();

  /**
   * <pre>
   * Output ports of the Gradient Controller.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.GradientController.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return Whether the outPorts field is set.
   */
  boolean hasOutPorts();
  /**
   * <pre>
   * Output ports of the Gradient Controller.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.GradientController.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return The outPorts.
   */
  com.fluxninja.generated.aperture.policy.language.v1.GradientController.Outs getOutPorts();
  /**
   * <pre>
   * Output ports of the Gradient Controller.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.GradientController.Outs out_ports = 2 [json_name = "outPorts"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.GradientController.OutsOrBuilder getOutPortsOrBuilder();

  /**
   * <pre>
   * Slope controls the aggressiveness and direction of the Gradient Controller.
   * Slope is used as exponent on the signal to setpoint ratio in computation
   * of the gradient (see the [main description](#v1-gradient-controller) for
   * exact equation). Good intuition for this parameter is "What should the
   * Gradient Controller do to the control variable when signal is too high",
   * eg.:
   * * $&#92;text{slope} = 1$: when signal is too high, increase control variable,
   * * $&#92;text{slope} = -1$: when signal is too high, decrease control variable,
   * * $&#92;text{slope} = -0.5$: when signal is to high, decrease control variable more slowly.
   * The sign of slope depends on correlation between the signal and control variable:
   * * Use $&#92;text{slope} &lt; 0$ if signal and control variable are _positively_
   * correlated (eg. Per-pod CPU usage and total concurrency).
   * * Use $&#92;text{slope} &gt; 0$ if signal and control variable are _negatively_
   * correlated (eg. Per-pod CPU usage and number of pods).
   * :::note
   * You need to set _negative_ slope for a _positive_ correlation, as you're
   * describing the _action_ which controller should make when the signal
   * increases.
   * :::
   * The magnitude of slope describes how aggressively should the controller
   * react to a deviation of signal.
   * With $|&#92;text{slope}| = 1$, the controller will aim to bring the signal to
   * the setpoint in one tick (assuming linear correlation with signal and setpoint).
   * Smaller magnitudes of slope will make the controller adjust the control
   * variable more slowly.
   * We recommend setting $|&#92;text{slope}| &lt; 1$ (eg. $&#92;pm0.8$).
   * If you experience overshooting, consider lowering the magnitude even more.
   * Values of $|&#92;text{slope}| &gt; 1$ are not recommended.
   * :::note
   * Remember that the gradient and output signal can be (optionally) clamped,
   * so the _slope_ might not fully describe aggressiveness of the controller.
   * :::
   * </pre>
   *
   * <code>double slope = 3 [json_name = "slope", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The slope.
   */
  double getSlope();

  /**
   * <pre>
   * Minimum gradient which clamps the computed gradient value to the range, [min_gradient, max_gradient].
   * </pre>
   *
   * <code>double min_gradient = 4 [json_name = "minGradient", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The minGradient.
   */
  double getMinGradient();

  /**
   * <pre>
   * Maximum gradient which clamps the computed gradient value to the range, [min_gradient, max_gradient].
   * </pre>
   *
   * <code>double max_gradient = 5 [json_name = "maxGradient", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The maxGradient.
   */
  double getMaxGradient();

  /**
   * <pre>
   * Configuration key for DynamicConfig
   * </pre>
   *
   * <code>string dynamic_config_key = 6 [json_name = "dynamicConfigKey"];</code>
   * @return The dynamicConfigKey.
   */
  java.lang.String getDynamicConfigKey();
  /**
   * <pre>
   * Configuration key for DynamicConfig
   * </pre>
   *
   * <code>string dynamic_config_key = 6 [json_name = "dynamicConfigKey"];</code>
   * @return The bytes for dynamicConfigKey.
   */
  com.google.protobuf.ByteString
      getDynamicConfigKeyBytes();

  /**
   * <pre>
   * Default configuration.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ControllerDynamicConfig default_config = 7 [json_name = "defaultConfig"];</code>
   * @return Whether the defaultConfig field is set.
   */
  boolean hasDefaultConfig();
  /**
   * <pre>
   * Default configuration.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ControllerDynamicConfig default_config = 7 [json_name = "defaultConfig"];</code>
   * @return The defaultConfig.
   */
  com.fluxninja.generated.aperture.policy.language.v1.ControllerDynamicConfig getDefaultConfig();
  /**
   * <pre>
   * Default configuration.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ControllerDynamicConfig default_config = 7 [json_name = "defaultConfig"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.ControllerDynamicConfigOrBuilder getDefaultConfigOrBuilder();
}
