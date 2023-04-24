// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

public interface LimiterDecisionOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.flowcontrol.check.v1.LimiterDecision)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>string policy_name = 1 [json_name = "policyName"];</code>
   * @return The policyName.
   */
  java.lang.String getPolicyName();
  /**
   * <code>string policy_name = 1 [json_name = "policyName"];</code>
   * @return The bytes for policyName.
   */
  com.google.protobuf.ByteString
      getPolicyNameBytes();

  /**
   * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
   * @return The policyHash.
   */
  java.lang.String getPolicyHash();
  /**
   * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
   * @return The bytes for policyHash.
   */
  com.google.protobuf.ByteString
      getPolicyHashBytes();

  /**
   * <code>string component_id = 3 [json_name = "componentId"];</code>
   * @return The componentId.
   */
  java.lang.String getComponentId();
  /**
   * <code>string component_id = 3 [json_name = "componentId"];</code>
   * @return The bytes for componentId.
   */
  com.google.protobuf.ByteString
      getComponentIdBytes();

  /**
   * <code>bool dropped = 4 [json_name = "dropped"];</code>
   * @return The dropped.
   */
  boolean getDropped();

  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.LimiterReason reason = 5 [json_name = "reason"];</code>
   * @return The enum numeric value on the wire for reason.
   */
  int getReasonValue();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.LimiterReason reason = 5 [json_name = "reason"];</code>
   * @return The reason.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.LimiterReason getReason();

  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.RateLimiterInfo rate_limiter_info = 6 [json_name = "rateLimiterInfo"];</code>
   * @return Whether the rateLimiterInfo field is set.
   */
  boolean hasRateLimiterInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.RateLimiterInfo rate_limiter_info = 6 [json_name = "rateLimiterInfo"];</code>
   * @return The rateLimiterInfo.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.RateLimiterInfo getRateLimiterInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.RateLimiterInfo rate_limiter_info = 6 [json_name = "rateLimiterInfo"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.RateLimiterInfoOrBuilder getRateLimiterInfoOrBuilder();

  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.ConcurrencyLimiterInfo concurrency_limiter_info = 7 [json_name = "concurrencyLimiterInfo"];</code>
   * @return Whether the concurrencyLimiterInfo field is set.
   */
  boolean hasConcurrencyLimiterInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.ConcurrencyLimiterInfo concurrency_limiter_info = 7 [json_name = "concurrencyLimiterInfo"];</code>
   * @return The concurrencyLimiterInfo.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.ConcurrencyLimiterInfo getConcurrencyLimiterInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.ConcurrencyLimiterInfo concurrency_limiter_info = 7 [json_name = "concurrencyLimiterInfo"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.ConcurrencyLimiterInfoOrBuilder getConcurrencyLimiterInfoOrBuilder();

  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.LoadRegulatorInfo load_regulator_info = 8 [json_name = "loadRegulatorInfo"];</code>
   * @return Whether the loadRegulatorInfo field is set.
   */
  boolean hasLoadRegulatorInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.LoadRegulatorInfo load_regulator_info = 8 [json_name = "loadRegulatorInfo"];</code>
   * @return The loadRegulatorInfo.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.LoadRegulatorInfo getLoadRegulatorInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.LoadRegulatorInfo load_regulator_info = 8 [json_name = "loadRegulatorInfo"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.LoadRegulatorInfoOrBuilder getLoadRegulatorInfoOrBuilder();

  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.DetailsCase getDetailsCase();
}
