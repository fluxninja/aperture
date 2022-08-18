// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/labelmatcher.proto

package com.aperture.policy.language.v1;

public interface K8sLabelMatcherRequirementOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.policy.language.v1.K8sLabelMatcherRequirement)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Label key that the selector applies to.
   * </pre>
   *
   * <code>string key = 1 [json_name = "key", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The key.
   */
  java.lang.String getKey();
  /**
   * <pre>
   * Label key that the selector applies to.
   * </pre>
   *
   * <code>string key = 1 [json_name = "key", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The bytes for key.
   */
  com.google.protobuf.ByteString
      getKeyBytes();

  /**
   * <pre>
   * Logical operator which represents a key's relationship to a set of values.
   * Valid operators are In, NotIn, Exists and DoesNotExist.
   * </pre>
   *
   * <code>string operator = 2 [json_name = "operator", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The operator.
   */
  java.lang.String getOperator();
  /**
   * <pre>
   * Logical operator which represents a key's relationship to a set of values.
   * Valid operators are In, NotIn, Exists and DoesNotExist.
   * </pre>
   *
   * <code>string operator = 2 [json_name = "operator", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The bytes for operator.
   */
  com.google.protobuf.ByteString
      getOperatorBytes();

  /**
   * <pre>
   * An array of string values that relates to the key by an operator.
   * If the operator is In or NotIn, the values array must be non-empty.
   * If the operator is Exists or DoesNotExist, the values array must be empty.
   * </pre>
   *
   * <code>repeated string values = 3 [json_name = "values"];</code>
   * @return A list containing the values.
   */
  java.util.List<java.lang.String>
      getValuesList();
  /**
   * <pre>
   * An array of string values that relates to the key by an operator.
   * If the operator is In or NotIn, the values array must be non-empty.
   * If the operator is Exists or DoesNotExist, the values array must be empty.
   * </pre>
   *
   * <code>repeated string values = 3 [json_name = "values"];</code>
   * @return The count of values.
   */
  int getValuesCount();
  /**
   * <pre>
   * An array of string values that relates to the key by an operator.
   * If the operator is In or NotIn, the values array must be non-empty.
   * If the operator is Exists or DoesNotExist, the values array must be empty.
   * </pre>
   *
   * <code>repeated string values = 3 [json_name = "values"];</code>
   * @param index The index of the element to return.
   * @return The values at the given index.
   */
  java.lang.String getValues(int index);
  /**
   * <pre>
   * An array of string values that relates to the key by an operator.
   * If the operator is In or NotIn, the values array must be non-empty.
   * If the operator is Exists or DoesNotExist, the values array must be empty.
   * </pre>
   *
   * <code>repeated string values = 3 [json_name = "values"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the values at the given index.
   */
  com.google.protobuf.ByteString
      getValuesBytes(int index);
}
