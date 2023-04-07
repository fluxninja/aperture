// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/sync/v1/common_attributes.proto

package com.fluxninja.generated.aperture.policy.sync.v1;

public interface CommonAttributesOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.sync.v1.CommonAttributes)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Name of the Policy.
   * </pre>
   *
   * <code>string policy_name = 1 [json_name = "policyName"];</code>
   * @return The policyName.
   */
  java.lang.String getPolicyName();
  /**
   * <pre>
   * Name of the Policy.
   * </pre>
   *
   * <code>string policy_name = 1 [json_name = "policyName"];</code>
   * @return The bytes for policyName.
   */
  com.google.protobuf.ByteString
      getPolicyNameBytes();

  /**
   * <pre>
   * Hash of the entire Policy spec.
   * </pre>
   *
   * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
   * @return The policyHash.
   */
  java.lang.String getPolicyHash();
  /**
   * <pre>
   * Hash of the entire Policy spec.
   * </pre>
   *
   * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
   * @return The bytes for policyHash.
   */
  com.google.protobuf.ByteString
      getPolicyHashBytes();

  /**
   * <pre>
   * The id of Component within the circuit.
   * </pre>
   *
   * <code>string component_id = 3 [json_name = "componentId"];</code>
   * @return The componentId.
   */
  java.lang.String getComponentId();
  /**
   * <pre>
   * The id of Component within the circuit.
   * </pre>
   *
   * <code>string component_id = 3 [json_name = "componentId"];</code>
   * @return The bytes for componentId.
   */
  com.google.protobuf.ByteString
      getComponentIdBytes();
}
