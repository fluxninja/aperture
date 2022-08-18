// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/classification/v1/ruleset.proto

package com.aperture.classification.v1;

public interface ClassifierOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.classification.v1.Classifier)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Defines where to apply the flow classification rule.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Selector selector = 1 [json_name = "selector"];</code>
   * @return Whether the selector field is set.
   */
  boolean hasSelector();
  /**
   * <pre>
   * Defines where to apply the flow classification rule.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Selector selector = 1 [json_name = "selector"];</code>
   * @return The selector.
   */
  com.aperture.policy.language.v1.Selector getSelector();
  /**
   * <pre>
   * Defines where to apply the flow classification rule.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Selector selector = 1 [json_name = "selector"];</code>
   */
  com.aperture.policy.language.v1.SelectorOrBuilder getSelectorOrBuilder();

  /**
   * <pre>
   * A map of {key, value} pairs mapping from flow label names to rules that define how to extract and propagate them.
   * </pre>
   *
   * <code>map&lt;string, .aperture.classification.v1.Rule&gt; rules = 2 [json_name = "rules"];</code>
   */
  int getRulesCount();
  /**
   * <pre>
   * A map of {key, value} pairs mapping from flow label names to rules that define how to extract and propagate them.
   * </pre>
   *
   * <code>map&lt;string, .aperture.classification.v1.Rule&gt; rules = 2 [json_name = "rules"];</code>
   */
  boolean containsRules(
      java.lang.String key);
  /**
   * Use {@link #getRulesMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, com.aperture.classification.v1.Rule>
  getRules();
  /**
   * <pre>
   * A map of {key, value} pairs mapping from flow label names to rules that define how to extract and propagate them.
   * </pre>
   *
   * <code>map&lt;string, .aperture.classification.v1.Rule&gt; rules = 2 [json_name = "rules"];</code>
   */
  java.util.Map<java.lang.String, com.aperture.classification.v1.Rule>
  getRulesMap();
  /**
   * <pre>
   * A map of {key, value} pairs mapping from flow label names to rules that define how to extract and propagate them.
   * </pre>
   *
   * <code>map&lt;string, .aperture.classification.v1.Rule&gt; rules = 2 [json_name = "rules"];</code>
   */

  /* nullable */
com.aperture.classification.v1.Rule getRulesOrDefault(
      java.lang.String key,
      /* nullable */
com.aperture.classification.v1.Rule defaultValue);
  /**
   * <pre>
   * A map of {key, value} pairs mapping from flow label names to rules that define how to extract and propagate them.
   * </pre>
   *
   * <code>map&lt;string, .aperture.classification.v1.Rule&gt; rules = 2 [json_name = "rules"];</code>
   */

  com.aperture.classification.v1.Rule getRulesOrThrow(
      java.lang.String key);
}
