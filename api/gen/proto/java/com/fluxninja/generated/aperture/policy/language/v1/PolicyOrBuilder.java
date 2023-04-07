// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/policy.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public interface PolicyOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.Policy)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Defines the control-loop logic of the policy.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Circuit circuit = 1 [json_name = "circuit"];</code>
   * @return Whether the circuit field is set.
   */
  boolean hasCircuit();
  /**
   * <pre>
   * Defines the control-loop logic of the policy.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Circuit circuit = 1 [json_name = "circuit"];</code>
   * @return The circuit.
   */
  com.fluxninja.generated.aperture.policy.language.v1.Circuit getCircuit();
  /**
   * <pre>
   * Defines the control-loop logic of the policy.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Circuit circuit = 1 [json_name = "circuit"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.CircuitOrBuilder getCircuitOrBuilder();

  /**
   * <pre>
   * Resources (Flux Meters, Classifiers etc.) to setup.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Resources resources = 2 [json_name = "resources"];</code>
   * @return Whether the resources field is set.
   */
  boolean hasResources();
  /**
   * <pre>
   * Resources (Flux Meters, Classifiers etc.) to setup.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Resources resources = 2 [json_name = "resources"];</code>
   * @return The resources.
   */
  com.fluxninja.generated.aperture.policy.language.v1.Resources getResources();
  /**
   * <pre>
   * Resources (Flux Meters, Classifiers etc.) to setup.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Resources resources = 2 [json_name = "resources"];</code>
   */
  com.fluxninja.generated.aperture.policy.language.v1.ResourcesOrBuilder getResourcesOrBuilder();
}
