// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

public interface KeyUpsertResponseOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.flowcontrol.check.v1.KeyUpsertResponse)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>.aperture.flowcontrol.check.v1.CacheOperationStatus operation_status = 1 [json_name = "operationStatus"];</code>
   * @return The enum numeric value on the wire for operationStatus.
   */
  int getOperationStatusValue();
  /**
   * <code>.aperture.flowcontrol.check.v1.CacheOperationStatus operation_status = 1 [json_name = "operationStatus"];</code>
   * @return The operationStatus.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheOperationStatus getOperationStatus();

  /**
   * <code>string error = 2 [json_name = "error"];</code>
   * @return The error.
   */
  java.lang.String getError();
  /**
   * <code>string error = 2 [json_name = "error"];</code>
   * @return The bytes for error.
   */
  com.google.protobuf.ByteString
      getErrorBytes();
}
