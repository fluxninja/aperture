// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/sync/v1/concurrency_limiter.proto

package com.fluxninja.generated.aperture.policy.sync.v1;

public interface TokensDecisionWrapperOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.sync.v1.TokensDecisionWrapper)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * CommonAttributes
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
   * @return Whether the commonAttributes field is set.
   */
  boolean hasCommonAttributes();
  /**
   * <pre>
   * CommonAttributes
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
   * @return The commonAttributes.
   */
  com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes getCommonAttributes();
  /**
   * <pre>
   * CommonAttributes
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
   */
  com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesOrBuilder getCommonAttributesOrBuilder();

  /**
   * <pre>
   * Tokens Decision
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.TokensDecision tokens_decision = 2 [json_name = "tokensDecision"];</code>
   * @return Whether the tokensDecision field is set.
   */
  boolean hasTokensDecision();
  /**
   * <pre>
   * Tokens Decision
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.TokensDecision tokens_decision = 2 [json_name = "tokensDecision"];</code>
   * @return The tokensDecision.
   */
  com.fluxninja.generated.aperture.policy.sync.v1.TokensDecision getTokensDecision();
  /**
   * <pre>
   * Tokens Decision
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.TokensDecision tokens_decision = 2 [json_name = "tokensDecision"];</code>
   */
  com.fluxninja.generated.aperture.policy.sync.v1.TokensDecisionOrBuilder getTokensDecisionOrBuilder();
}
