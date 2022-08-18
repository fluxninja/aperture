// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/policy.proto

package com.aperture.policy.language.v1;

public interface ArithmeticCombinatorOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.ArithmeticCombinator)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Input ports for the Arithmetic Combinator component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ArithmeticCombinator.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return Whether the inPorts field is set.
   */
  boolean hasInPorts();
  /**
   * <pre>
   * Input ports for the Arithmetic Combinator component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ArithmeticCombinator.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return The inPorts.
   */
  com.aperture.policy.language.v1.ArithmeticCombinator.Ins getInPorts();
  /**
   * <pre>
   * Input ports for the Arithmetic Combinator component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ArithmeticCombinator.Ins in_ports = 1 [json_name = "inPorts"];</code>
   */
  com.aperture.policy.language.v1.ArithmeticCombinator.InsOrBuilder getInPortsOrBuilder();

  /**
   * <pre>
   * Output ports for the Arithmetic Combinator component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ArithmeticCombinator.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return Whether the outPorts field is set.
   */
  boolean hasOutPorts();
  /**
   * <pre>
   * Output ports for the Arithmetic Combinator component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ArithmeticCombinator.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return The outPorts.
   */
  com.aperture.policy.language.v1.ArithmeticCombinator.Outs getOutPorts();
  /**
   * <pre>
   * Output ports for the Arithmetic Combinator component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.ArithmeticCombinator.Outs out_ports = 2 [json_name = "outPorts"];</code>
   */
  com.aperture.policy.language.v1.ArithmeticCombinator.OutsOrBuilder getOutPortsOrBuilder();

  /**
   * <pre>
   * Operator of the arithmetic operation.
   * </pre>
   *
   * <code>string operator = 3 [json_name = "operator", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The operator.
   */
  java.lang.String getOperator();
  /**
   * <pre>
   * Operator of the arithmetic operation.
   * </pre>
   *
   * <code>string operator = 3 [json_name = "operator", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The bytes for operator.
   */
  com.google.protobuf.ByteString
      getOperatorBytes();
}
