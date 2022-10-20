// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/config/core/v3/socket_option.proto

package com.fluxninja.generated.envoy.config.core.v3;

public interface SocketOptionOrBuilder extends
    // @@protoc_insertion_point(interface_extends:envoy.config.core.v3.SocketOption)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * An optional name to give this socket option for debugging, etc.
   * Uniqueness is not required and no special meaning is assumed.
   * </pre>
   *
   * <code>string description = 1 [json_name = "description"];</code>
   * @return The description.
   */
  java.lang.String getDescription();
  /**
   * <pre>
   * An optional name to give this socket option for debugging, etc.
   * Uniqueness is not required and no special meaning is assumed.
   * </pre>
   *
   * <code>string description = 1 [json_name = "description"];</code>
   * @return The bytes for description.
   */
  com.google.protobuf.ByteString
      getDescriptionBytes();

  /**
   * <pre>
   * Corresponding to the level value passed to setsockopt, such as IPPROTO_TCP
   * </pre>
   *
   * <code>int64 level = 2 [json_name = "level"];</code>
   * @return The level.
   */
  long getLevel();

  /**
   * <pre>
   * The numeric name as passed to setsockopt
   * </pre>
   *
   * <code>int64 name = 3 [json_name = "name"];</code>
   * @return The name.
   */
  long getName();

  /**
   * <pre>
   * Because many sockopts take an int value.
   * </pre>
   *
   * <code>int64 int_value = 4 [json_name = "intValue"];</code>
   * @return Whether the intValue field is set.
   */
  boolean hasIntValue();
  /**
   * <pre>
   * Because many sockopts take an int value.
   * </pre>
   *
   * <code>int64 int_value = 4 [json_name = "intValue"];</code>
   * @return The intValue.
   */
  long getIntValue();

  /**
   * <pre>
   * Otherwise it's a byte buffer.
   * </pre>
   *
   * <code>bytes buf_value = 5 [json_name = "bufValue"];</code>
   * @return Whether the bufValue field is set.
   */
  boolean hasBufValue();
  /**
   * <pre>
   * Otherwise it's a byte buffer.
   * </pre>
   *
   * <code>bytes buf_value = 5 [json_name = "bufValue"];</code>
   * @return The bufValue.
   */
  com.google.protobuf.ByteString getBufValue();

  /**
   * <pre>
   * The state in which the option will be applied. When used in BindConfig
   * STATE_PREBIND is currently the only valid value.
   * </pre>
   *
   * <code>.envoy.config.core.v3.SocketOption.SocketState state = 6 [json_name = "state", (.validate.rules) = { ... }</code>
   * @return The enum numeric value on the wire for state.
   */
  int getStateValue();
  /**
   * <pre>
   * The state in which the option will be applied. When used in BindConfig
   * STATE_PREBIND is currently the only valid value.
   * </pre>
   *
   * <code>.envoy.config.core.v3.SocketOption.SocketState state = 6 [json_name = "state", (.validate.rules) = { ... }</code>
   * @return The state.
   */
  com.fluxninja.generated.envoy.config.core.v3.SocketOption.SocketState getState();

  public com.fluxninja.generated.envoy.config.core.v3.SocketOption.ValueCase getValueCase();
}
