// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/std_components.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public interface AndOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.And)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Input ports for the And component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.And.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return Whether the inPorts field is set.
   */
  boolean hasInPorts();
  /**
   * <pre>
   * Input ports for the And component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.And.Ins in_ports = 1 [json_name = "inPorts"];</code>
   * @return The inPorts.
   */
  com.fluxninja.generated.aperture.policy.language.v1.And.Ins getInPorts();
  /**
   * <pre>
   * Input ports for the And component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.And.Ins in_ports = 1 [json_name = "inPorts"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.And.InsOrBuilder getInPortsOrBuilder();

  /**
   * <pre>
   * Output ports for the And component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.And.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return Whether the outPorts field is set.
   */
  boolean hasOutPorts();
  /**
   * <pre>
   * Output ports for the And component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.And.Outs out_ports = 2 [json_name = "outPorts"];</code>
   * @return The outPorts.
   */
  com.fluxninja.generated.aperture.policy.language.v1.And.Outs getOutPorts();
  /**
   * <pre>
   * Output ports for the And component.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.And.Outs out_ports = 2 [json_name = "outPorts"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.And.OutsOrBuilder getOutPortsOrBuilder();
}
