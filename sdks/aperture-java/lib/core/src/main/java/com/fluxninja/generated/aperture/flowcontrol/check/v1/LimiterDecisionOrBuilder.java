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
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.SchedulerInfo load_scheduler_info = 7 [json_name = "loadSchedulerInfo"];</code>
   * @return Whether the loadSchedulerInfo field is set.
   */
  boolean hasLoadSchedulerInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.SchedulerInfo load_scheduler_info = 7 [json_name = "loadSchedulerInfo"];</code>
   * @return The loadSchedulerInfo.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.SchedulerInfo getLoadSchedulerInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.SchedulerInfo load_scheduler_info = 7 [json_name = "loadSchedulerInfo"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.SchedulerInfoOrBuilder getLoadSchedulerInfoOrBuilder();

  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.RegulatorInfo regulator_info = 8 [json_name = "regulatorInfo"];</code>
   * @return Whether the regulatorInfo field is set.
   */
  boolean hasRegulatorInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.RegulatorInfo regulator_info = 8 [json_name = "regulatorInfo"];</code>
   * @return The regulatorInfo.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.RegulatorInfo getRegulatorInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.RegulatorInfo regulator_info = 8 [json_name = "regulatorInfo"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.RegulatorInfoOrBuilder getRegulatorInfoOrBuilder();

  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.QuotaSchedulerInfo quota_scheduler_info = 9 [json_name = "quotaSchedulerInfo"];</code>
   * @return Whether the quotaSchedulerInfo field is set.
   */
  boolean hasQuotaSchedulerInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.QuotaSchedulerInfo quota_scheduler_info = 9 [json_name = "quotaSchedulerInfo"];</code>
   * @return The quotaSchedulerInfo.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.QuotaSchedulerInfo getQuotaSchedulerInfo();
  /**
   * <code>.aperture.flowcontrol.check.v1.LimiterDecision.QuotaSchedulerInfo quota_scheduler_info = 9 [json_name = "quotaSchedulerInfo"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.QuotaSchedulerInfoOrBuilder getQuotaSchedulerInfoOrBuilder();

  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision.DetailsCase getDetailsCase();
}
