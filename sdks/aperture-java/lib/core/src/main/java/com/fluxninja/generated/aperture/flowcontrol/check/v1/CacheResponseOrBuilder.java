// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

public interface CacheResponseOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.flowcontrol.check.v1.CacheResponse)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>bytes value = 1 [json_name = "value"];</code>
   * @return The value.
   */
  com.google.protobuf.ByteString getValue();

  /**
   * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
   * @return The enum numeric value on the wire for result.
   */
  int getResultValue();
  /**
   * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
   * @return The result.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult getResult();
}
