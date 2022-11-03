// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/policy.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public interface ResourcesOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.Resources)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.
   * Flux Meter created metrics can be consumed as input to the circuit via the PromQL component.
   * </pre>
   *
   * <code>map&lt;string, .aperture.policy.language.v1.FluxMeter&gt; flux_meters = 1 [json_name = "fluxMeters"];</code>
   */
  int getFluxMetersCount();
  /**
   * <pre>
   * Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.
   * Flux Meter created metrics can be consumed as input to the circuit via the PromQL component.
   * </pre>
   *
   * <code>map&lt;string, .aperture.policy.language.v1.FluxMeter&gt; flux_meters = 1 [json_name = "fluxMeters"];</code>
   */
  boolean containsFluxMeters(
      java.lang.String key);
  /**
   * Use {@link #getFluxMetersMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.FluxMeter>
  getFluxMeters();
  /**
   * <pre>
   * Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.
   * Flux Meter created metrics can be consumed as input to the circuit via the PromQL component.
   * </pre>
   *
   * <code>map&lt;string, .aperture.policy.language.v1.FluxMeter&gt; flux_meters = 1 [json_name = "fluxMeters"];</code>
   */
  java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.FluxMeter>
  getFluxMetersMap();
  /**
   * <pre>
   * Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.
   * Flux Meter created metrics can be consumed as input to the circuit via the PromQL component.
   * </pre>
   *
   * <code>map&lt;string, .aperture.policy.language.v1.FluxMeter&gt; flux_meters = 1 [json_name = "fluxMeters"];</code>
   */

  /* nullable */
com.fluxninja.generated.aperture.policy.language.v1.FluxMeter getFluxMetersOrDefault(
      java.lang.String key,
      /* nullable */
com.fluxninja.generated.aperture.policy.language.v1.FluxMeter defaultValue);
  /**
   * <pre>
   * Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.
   * Flux Meter created metrics can be consumed as input to the circuit via the PromQL component.
   * </pre>
   *
   * <code>map&lt;string, .aperture.policy.language.v1.FluxMeter&gt; flux_meters = 1 [json_name = "fluxMeters"];</code>
   */

  com.fluxninja.generated.aperture.policy.language.v1.FluxMeter getFluxMetersOrThrow(
      java.lang.String key);

  /**
   * <pre>
   * Classifiers are installed in the data-plane and are used to label the requests based on payload content.
   * The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes.
   * </pre>
   *
   * <code>repeated .aperture.policy.language.v1.Classifier classifiers = 2 [json_name = "classifiers"];</code>
   */
  java.util.List<com.fluxninja.generated.aperture.policy.language.v1.Classifier> 
      getClassifiersList();
  /**
   * <pre>
   * Classifiers are installed in the data-plane and are used to label the requests based on payload content.
   * The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes.
   * </pre>
   *
   * <code>repeated .aperture.policy.language.v1.Classifier classifiers = 2 [json_name = "classifiers"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.Classifier getClassifiers(int index);
  /**
   * <pre>
   * Classifiers are installed in the data-plane and are used to label the requests based on payload content.
   * The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes.
   * </pre>
   *
   * <code>repeated .aperture.policy.language.v1.Classifier classifiers = 2 [json_name = "classifiers"];</code>
   */
  int getClassifiersCount();
  /**
   * <pre>
   * Classifiers are installed in the data-plane and are used to label the requests based on payload content.
   * The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes.
   * </pre>
   *
   * <code>repeated .aperture.policy.language.v1.Classifier classifiers = 2 [json_name = "classifiers"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.aperture.policy.language.v1.ClassifierOrBuilder> 
      getClassifiersOrBuilderList();
  /**
   * <pre>
   * Classifiers are installed in the data-plane and are used to label the requests based on payload content.
   * The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes.
   * </pre>
   *
   * <code>repeated .aperture.policy.language.v1.Classifier classifiers = 2 [json_name = "classifiers"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.ClassifierOrBuilder getClassifiersOrBuilder(
      int index);
}
