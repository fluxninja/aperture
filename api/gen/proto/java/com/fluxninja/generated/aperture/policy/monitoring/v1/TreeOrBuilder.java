// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/monitoring/v1/policy_view.proto

package com.fluxninja.generated.aperture.policy.monitoring.v1;

public interface TreeOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.monitoring.v1.Tree)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>.aperture.policy.monitoring.v1.ComponentView node = 1 [json_name = "node"];</code>
   * @return Whether the node field is set.
   */
  boolean hasNode();
  /**
   * <code>.aperture.policy.monitoring.v1.ComponentView node = 1 [json_name = "node"];</code>
   * @return The node.
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.ComponentView getNode();
  /**
   * <code>.aperture.policy.monitoring.v1.ComponentView node = 1 [json_name = "node"];</code>
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.ComponentViewOrBuilder getNodeOrBuilder();

  /**
   * <code>.aperture.policy.monitoring.v1.Graph graph = 2 [json_name = "graph"];</code>
   * @return Whether the graph field is set.
   */
  boolean hasGraph();
  /**
   * <code>.aperture.policy.monitoring.v1.Graph graph = 2 [json_name = "graph"];</code>
   * @return The graph.
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.Graph getGraph();
  /**
   * <code>.aperture.policy.monitoring.v1.Graph graph = 2 [json_name = "graph"];</code>
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.GraphOrBuilder getGraphOrBuilder();

  /**
   * <code>repeated .aperture.policy.monitoring.v1.Tree children = 3 [json_name = "children"];</code>
   */
  java.util.List<com.fluxninja.generated.aperture.policy.monitoring.v1.Tree> 
      getChildrenList();
  /**
   * <code>repeated .aperture.policy.monitoring.v1.Tree children = 3 [json_name = "children"];</code>
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.Tree getChildren(int index);
  /**
   * <code>repeated .aperture.policy.monitoring.v1.Tree children = 3 [json_name = "children"];</code>
   */
  int getChildrenCount();
  /**
   * <code>repeated .aperture.policy.monitoring.v1.Tree children = 3 [json_name = "children"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.aperture.policy.monitoring.v1.TreeOrBuilder> 
      getChildrenOrBuilderList();
  /**
   * <code>repeated .aperture.policy.monitoring.v1.Tree children = 3 [json_name = "children"];</code>
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.TreeOrBuilder getChildrenOrBuilder(
      int index);

  /**
   * <code>repeated .aperture.policy.monitoring.v1.ComponentView actuators = 4 [json_name = "actuators"];</code>
   */
  java.util.List<com.fluxninja.generated.aperture.policy.monitoring.v1.ComponentView> 
      getActuatorsList();
  /**
   * <code>repeated .aperture.policy.monitoring.v1.ComponentView actuators = 4 [json_name = "actuators"];</code>
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.ComponentView getActuators(int index);
  /**
   * <code>repeated .aperture.policy.monitoring.v1.ComponentView actuators = 4 [json_name = "actuators"];</code>
   */
  int getActuatorsCount();
  /**
   * <code>repeated .aperture.policy.monitoring.v1.ComponentView actuators = 4 [json_name = "actuators"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.aperture.policy.monitoring.v1.ComponentViewOrBuilder> 
      getActuatorsOrBuilderList();
  /**
   * <code>repeated .aperture.policy.monitoring.v1.ComponentView actuators = 4 [json_name = "actuators"];</code>
   */
  com.fluxninja.generated.aperture.policy.monitoring.v1.ComponentViewOrBuilder getActuatorsOrBuilder(
      int index);
}
