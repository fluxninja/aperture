// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: validate/validate.proto

package com.fluxninja.generated.validate;

public interface AnyRulesOrBuilder extends
    // @@protoc_insertion_point(interface_extends:validate.AnyRules)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Required specifies that this field must be set
   * </pre>
   *
   * <code>optional bool required = 1 [json_name = "required"];</code>
   * @return Whether the required field is set.
   */
  boolean hasRequired();
  /**
   * <pre>
   * Required specifies that this field must be set
   * </pre>
   *
   * <code>optional bool required = 1 [json_name = "required"];</code>
   * @return The required.
   */
  boolean getRequired();

  /**
   * <pre>
   * In specifies that this field's `type_url` must be equal to one of the
   * specified values.
   * </pre>
   *
   * <code>repeated string in = 2 [json_name = "in"];</code>
   * @return A list containing the in.
   */
  java.util.List<java.lang.String>
      getInList();
  /**
   * <pre>
   * In specifies that this field's `type_url` must be equal to one of the
   * specified values.
   * </pre>
   *
   * <code>repeated string in = 2 [json_name = "in"];</code>
   * @return The count of in.
   */
  int getInCount();
  /**
   * <pre>
   * In specifies that this field's `type_url` must be equal to one of the
   * specified values.
   * </pre>
   *
   * <code>repeated string in = 2 [json_name = "in"];</code>
   * @param index The index of the element to return.
   * @return The in at the given index.
   */
  java.lang.String getIn(int index);
  /**
   * <pre>
   * In specifies that this field's `type_url` must be equal to one of the
   * specified values.
   * </pre>
   *
   * <code>repeated string in = 2 [json_name = "in"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the in at the given index.
   */
  com.google.protobuf.ByteString
      getInBytes(int index);

  /**
   * <pre>
   * NotIn specifies that this field's `type_url` must not be equal to any of
   * the specified values.
   * </pre>
   *
   * <code>repeated string not_in = 3 [json_name = "notIn"];</code>
   * @return A list containing the notIn.
   */
  java.util.List<java.lang.String>
      getNotInList();
  /**
   * <pre>
   * NotIn specifies that this field's `type_url` must not be equal to any of
   * the specified values.
   * </pre>
   *
   * <code>repeated string not_in = 3 [json_name = "notIn"];</code>
   * @return The count of notIn.
   */
  int getNotInCount();
  /**
   * <pre>
   * NotIn specifies that this field's `type_url` must not be equal to any of
   * the specified values.
   * </pre>
   *
   * <code>repeated string not_in = 3 [json_name = "notIn"];</code>
   * @param index The index of the element to return.
   * @return The notIn at the given index.
   */
  java.lang.String getNotIn(int index);
  /**
   * <pre>
   * NotIn specifies that this field's `type_url` must not be equal to any of
   * the specified values.
   * </pre>
   *
   * <code>repeated string not_in = 3 [json_name = "notIn"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the notIn at the given index.
   */
  com.google.protobuf.ByteString
      getNotInBytes(int index);
}