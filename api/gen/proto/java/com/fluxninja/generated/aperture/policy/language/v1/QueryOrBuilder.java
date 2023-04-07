// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/query.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public interface QueryOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.Query)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Periodically runs a Prometheus query in the background and emits the result.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.PromQL promql = 1 [json_name = "promql"];</code>
   * @return Whether the promql field is set.
   */
  boolean hasPromql();
  /**
   * <pre>
   * Periodically runs a Prometheus query in the background and emits the result.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.PromQL promql = 1 [json_name = "promql"];</code>
   * @return The promql.
   */
  com.fluxninja.generated.aperture.policy.language.v1.PromQL getPromql();
  /**
   * <pre>
   * Periodically runs a Prometheus query in the background and emits the result.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.PromQL promql = 1 [json_name = "promql"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.PromQLOrBuilder getPromqlOrBuilder();

  public com.fluxninja.generated.aperture.policy.language.v1.Query.ComponentCase getComponentCase();
}
