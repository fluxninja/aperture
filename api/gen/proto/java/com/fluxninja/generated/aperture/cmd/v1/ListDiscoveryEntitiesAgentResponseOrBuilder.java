// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/cmd/v1/cmd.proto

package com.fluxninja.generated.aperture.cmd.v1;

public interface ListDiscoveryEntitiesAgentResponseOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.cmd.v1.ListDiscoveryEntitiesAgentResponse)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>map&lt;string, .aperture.discovery.entities.v1.Entity&gt; entities = 1 [json_name = "entities"];</code>
   */
  int getEntitiesCount();
  /**
   * <code>map&lt;string, .aperture.discovery.entities.v1.Entity&gt; entities = 1 [json_name = "entities"];</code>
   */
  boolean containsEntities(
      java.lang.String key);
  /**
   * Use {@link #getEntitiesMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, com.fluxninja.generated.aperture.discovery.entities.v1.Entity>
  getEntities();
  /**
   * <code>map&lt;string, .aperture.discovery.entities.v1.Entity&gt; entities = 1 [json_name = "entities"];</code>
   */
  java.util.Map<java.lang.String, com.fluxninja.generated.aperture.discovery.entities.v1.Entity>
  getEntitiesMap();
  /**
   * <code>map&lt;string, .aperture.discovery.entities.v1.Entity&gt; entities = 1 [json_name = "entities"];</code>
   */

  /* nullable */
com.fluxninja.generated.aperture.discovery.entities.v1.Entity getEntitiesOrDefault(
      java.lang.String key,
      /* nullable */
com.fluxninja.generated.aperture.discovery.entities.v1.Entity defaultValue);
  /**
   * <code>map&lt;string, .aperture.discovery.entities.v1.Entity&gt; entities = 1 [json_name = "entities"];</code>
   */

  com.fluxninja.generated.aperture.discovery.entities.v1.Entity getEntitiesOrThrow(
      java.lang.String key);
}
