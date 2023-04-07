// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/flowcontrol.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public interface RuleOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.Rule)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * High-level declarative extractor.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Extractor extractor = 1 [json_name = "extractor"];</code>
   * @return Whether the extractor field is set.
   */
  boolean hasExtractor();
  /**
   * <pre>
   * High-level declarative extractor.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Extractor extractor = 1 [json_name = "extractor"];</code>
   * @return The extractor.
   */
  com.fluxninja.generated.aperture.policy.language.v1.Extractor getExtractor();
  /**
   * <pre>
   * High-level declarative extractor.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Extractor extractor = 1 [json_name = "extractor"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.ExtractorOrBuilder getExtractorOrBuilder();

  /**
   * <pre>
   * Rego module to extract a value from.
   * Deprecated: 1.5.0
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Rule.Rego rego = 2 [json_name = "rego"];</code>
   * @return Whether the rego field is set.
   */
  boolean hasRego();
  /**
   * <pre>
   * Rego module to extract a value from.
   * Deprecated: 1.5.0
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Rule.Rego rego = 2 [json_name = "rego"];</code>
   * @return The rego.
   */
  com.fluxninja.generated.aperture.policy.language.v1.Rule.Rego getRego();
  /**
   * <pre>
   * Rego module to extract a value from.
   * Deprecated: 1.5.0
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Rule.Rego rego = 2 [json_name = "rego"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.Rule.RegoOrBuilder getRegoOrBuilder();

  /**
   * <pre>
   * Decides if the created flow label should be available as an attribute in OLAP telemetry and
   * propagated in [baggage](/concepts/flow-control/flow-label.md#baggage)
   * :::note
   * The flow label is always accessible in Aperture Policies regardless of this setting.
   * :::
   * :::caution
   * When using [FluxNinja ARC extension](arc/extension.md), telemetry enabled
   * labels are sent to FluxNinja ARC for observability. Telemetry should be disabled for
   * sensitive labels.
   * :::
   * </pre>
   *
   * <code>bool telemetry = 3 [json_name = "telemetry"];</code>
   * @return The telemetry.
   */
  boolean getTelemetry();

  public com.fluxninja.generated.aperture.policy.language.v1.Rule.SourceCase getSourceCase();
}
