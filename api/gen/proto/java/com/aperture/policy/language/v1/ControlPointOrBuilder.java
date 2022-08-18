// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/selector.proto

package com.aperture.policy.language.v1;

public interface ControlPointOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.ControlPoint)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Name of FlunxNinja library's feature.
   * Feature corresponds to a block of code that can be "switched off" which usually is a "named opentelemetry's Span".
   * Note: Flowcontrol only.
   * </pre>
   *
   * <code>string feature = 1 [json_name = "feature", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return Whether the feature field is set.
   */
  boolean hasFeature();
  /**
   * <pre>
   * Name of FlunxNinja library's feature.
   * Feature corresponds to a block of code that can be "switched off" which usually is a "named opentelemetry's Span".
   * Note: Flowcontrol only.
   * </pre>
   *
   * <code>string feature = 1 [json_name = "feature", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The feature.
   */
  java.lang.String getFeature();
  /**
   * <pre>
   * Name of FlunxNinja library's feature.
   * Feature corresponds to a block of code that can be "switched off" which usually is a "named opentelemetry's Span".
   * Note: Flowcontrol only.
   * </pre>
   *
   * <code>string feature = 1 [json_name = "feature", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The bytes for feature.
   */
  com.google.protobuf.ByteString
      getFeatureBytes();

  /**
   * <pre>
   * Type of traffic service, either "ingress" or "egress".
   * Apply the policy to the whole incoming/outgoing traffic of a service.
   * Usually powered by integration with a proxy (like envoy) or a web framework.
   * * Flowcontrol: Blockable atom here is a single HTTP-transaction.
   * * Classification: Apply the classification rules to every incoming/outgoing request and attach the resulting flow labels to baggage and telemetry.
   * </pre>
   *
   * <code>string traffic = 2 [json_name = "traffic", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return Whether the traffic field is set.
   */
  boolean hasTraffic();
  /**
   * <pre>
   * Type of traffic service, either "ingress" or "egress".
   * Apply the policy to the whole incoming/outgoing traffic of a service.
   * Usually powered by integration with a proxy (like envoy) or a web framework.
   * * Flowcontrol: Blockable atom here is a single HTTP-transaction.
   * * Classification: Apply the classification rules to every incoming/outgoing request and attach the resulting flow labels to baggage and telemetry.
   * </pre>
   *
   * <code>string traffic = 2 [json_name = "traffic", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The traffic.
   */
  java.lang.String getTraffic();
  /**
   * <pre>
   * Type of traffic service, either "ingress" or "egress".
   * Apply the policy to the whole incoming/outgoing traffic of a service.
   * Usually powered by integration with a proxy (like envoy) or a web framework.
   * * Flowcontrol: Blockable atom here is a single HTTP-transaction.
   * * Classification: Apply the classification rules to every incoming/outgoing request and attach the resulting flow labels to baggage and telemetry.
   * </pre>
   *
   * <code>string traffic = 2 [json_name = "traffic", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The bytes for traffic.
   */
  com.google.protobuf.ByteString
      getTrafficBytes();

  public com.aperture.policy.language.v1.ControlPoint.ControlpointCase getControlpointCase();
}
