// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/cmd/v1/cmd.proto

package com.fluxninja.generated.aperture.cmd.v1;

public interface GlobalFlowControlPointOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.cmd.v1.GlobalFlowControlPoint)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>.aperture.flowcontrol.controlpoints.v1.FlowControlPoint flow_control_point = 1 [json_name = "flowControlPoint"];</code>
   * @return Whether the flowControlPoint field is set.
   */
  boolean hasFlowControlPoint();
  /**
   * <code>.aperture.flowcontrol.controlpoints.v1.FlowControlPoint flow_control_point = 1 [json_name = "flowControlPoint"];</code>
   * @return The flowControlPoint.
   */
  com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoint getFlowControlPoint();
  /**
   * <code>.aperture.flowcontrol.controlpoints.v1.FlowControlPoint flow_control_point = 1 [json_name = "flowControlPoint"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPointOrBuilder getFlowControlPointOrBuilder();

  /**
   * <code>string agent_group = 2 [json_name = "agentGroup"];</code>
   * @return The agentGroup.
   */
  java.lang.String getAgentGroup();
  /**
   * <code>string agent_group = 2 [json_name = "agentGroup"];</code>
   * @return The bytes for agentGroup.
   */
  com.google.protobuf.ByteString
      getAgentGroupBytes();
}
