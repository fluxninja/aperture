// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/policy.proto

package com.aperture.policy.language.v1;

public interface DeciderOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.Decider)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Input ports for the Decider component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Decider.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return Whether the inPorts field is set.
   */
  boolean hasInPorts();
  /**
   * <pre>
   * Input ports for the Decider component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Decider.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return The inPorts.
   */
  com.aperture.policy.language.v1.Decider.Ins getInPorts();
  /**
   * <pre>
   * Input ports for the Decider component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Decider.Ins in_ports = 1 [json_name = "inPorts"];</code>
   */
  com.aperture.policy.language.v1.Decider.InsOrBuilder getInPortsOrBuilder();

  /**
   * <pre>
   * Output ports for the Decider component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Decider.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return Whether the outPorts field is set.
   */
  boolean hasOutPorts();
  /**
   * <pre>
   * Output ports for the Decider component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Decider.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return The outPorts.
   */
  com.aperture.policy.language.v1.Decider.Outs getOutPorts();
  /**
   * <pre>
   * Output ports for the Decider component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Decider.Outs out_ports = 2 [json_name = "outPorts"];</code>
   */
  com.aperture.policy.language.v1.Decider.OutsOrBuilder getOutPortsOrBuilder();

  /**
   * <pre>
   * Comparison operator that computes operation on lhs and rhs input signals.
   * </pre>
   *
   * <code>string operator = 3 [json_name = "operator", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The operator.
   */
  java.lang.String getOperator();
  /**
   * <pre>
   * Comparison operator that computes operation on lhs and rhs input signals.
   * </pre>
   *
   * <code>string operator = 3 [json_name = "operator", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The bytes for operator.
   */
  com.google.protobuf.ByteString
      getOperatorBytes();

  /**
   * <pre>
   * Duration of time to wait before a transition to true state.
   * If the duration is zero, the transition will happen instantaneously.
   * </pre>
   *
   * <code>.google.protobuf.Duration true_for = 4 [json_name = "trueFor", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return Whether the trueFor field is set.
   */
  boolean hasTrueFor();
  /**
   * <pre>
   * Duration of time to wait before a transition to true state.
   * If the duration is zero, the transition will happen instantaneously.
   * </pre>
   *
   * <code>.google.protobuf.Duration true_for = 4 [json_name = "trueFor", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The trueFor.
   */
  com.google.protobuf.Duration getTrueFor();
  /**
   * <pre>
   * Duration of time to wait before a transition to true state.
   * If the duration is zero, the transition will happen instantaneously.
   * </pre>
   *
   * <code>.google.protobuf.Duration true_for = 4 [json_name = "trueFor", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   */
  com.google.protobuf.DurationOrBuilder getTrueForOrBuilder();

  /**
   * <pre>
   * Duration of time to wait before a transition to false state.
   * If the duration is zero, the transition will happen instantaneously.
   * </pre>
   *
   * <code>.google.protobuf.Duration false_for = 5 [json_name = "falseFor", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return Whether the falseFor field is set.
   */
  boolean hasFalseFor();
  /**
   * <pre>
   * Duration of time to wait before a transition to false state.
   * If the duration is zero, the transition will happen instantaneously.
   * </pre>
   *
   * <code>.google.protobuf.Duration false_for = 5 [json_name = "falseFor", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The falseFor.
   */
  com.google.protobuf.Duration getFalseFor();
  /**
   * <pre>
   * Duration of time to wait before a transition to false state.
   * If the duration is zero, the transition will happen instantaneously.
   * </pre>
   *
   * <code>.google.protobuf.Duration false_for = 5 [json_name = "falseFor", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   */
  com.google.protobuf.DurationOrBuilder getFalseForOrBuilder();
}
