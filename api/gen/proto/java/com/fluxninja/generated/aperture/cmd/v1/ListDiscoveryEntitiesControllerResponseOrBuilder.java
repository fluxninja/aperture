// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/cmd/v1/cmd.proto

package com.fluxninja.generated.aperture.cmd.v1;

public interface ListDiscoveryEntitiesControllerResponseOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>.aperture.cmd.v1.ListDiscoveryEntitiesAgentResponse entities = 1 [json_name = "entities"];</code>
   * @return Whether the entities field is set.
   */
  boolean hasEntities();
  /**
   * <code>.aperture.cmd.v1.ListDiscoveryEntitiesAgentResponse entities = 1 [json_name = "entities"];</code>
   * @return The entities.
   */
  com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesAgentResponse getEntities();
  /**
   * <code>.aperture.cmd.v1.ListDiscoveryEntitiesAgentResponse entities = 1 [json_name = "entities"];</code>
   */
  com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesAgentResponseOrBuilder getEntitiesOrBuilder();

  /**
   * <code>uint32 errors_count = 2 [json_name = "errorsCount"];</code>
   * @return The errorsCount.
   */
  int getErrorsCount();
}
