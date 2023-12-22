// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

public interface FlowEndRequestOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.flowcontrol.check.v1.FlowEndRequest)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * needed for telemetry
   * </pre>
   *
   * <code>string control_point = 1 [json_name = "controlPoint"];</code>
   * @return The controlPoint.
   */
  java.lang.String getControlPoint();
  /**
   * <pre>
   * needed for telemetry
   * </pre>
   *
   * <code>string control_point = 1 [json_name = "controlPoint"];</code>
   * @return The bytes for controlPoint.
   */
  com.google.protobuf.ByteString
      getControlPointBytes();

  /**
   * <code>repeated .aperture.flowcontrol.check.v1.InflightRequestRef inflight_requests = 2 [json_name = "inflightRequests"];</code>
   */
  java.util.List<com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef> 
      getInflightRequestsList();
  /**
   * <code>repeated .aperture.flowcontrol.check.v1.InflightRequestRef inflight_requests = 2 [json_name = "inflightRequests"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef getInflightRequests(int index);
  /**
   * <code>repeated .aperture.flowcontrol.check.v1.InflightRequestRef inflight_requests = 2 [json_name = "inflightRequests"];</code>
   */
  int getInflightRequestsCount();
  /**
   * <code>repeated .aperture.flowcontrol.check.v1.InflightRequestRef inflight_requests = 2 [json_name = "inflightRequests"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRefOrBuilder> 
      getInflightRequestsOrBuilderList();
  /**
   * <code>repeated .aperture.flowcontrol.check.v1.InflightRequestRef inflight_requests = 2 [json_name = "inflightRequests"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRefOrBuilder getInflightRequestsOrBuilder(
      int index);
}