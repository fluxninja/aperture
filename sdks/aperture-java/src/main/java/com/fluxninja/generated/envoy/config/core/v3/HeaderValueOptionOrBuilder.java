// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/config/core/v3/base.proto

package com.fluxninja.generated.envoy.config.core.v3;

public interface HeaderValueOptionOrBuilder extends
    // @@protoc_insertion_point(interface_extends:envoy.config.core.v3.HeaderValueOption)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Header name/value pair that this option applies to.
   * </pre>
   *
   * <code>.envoy.config.core.v3.HeaderValue header = 1 [json_name = "header", (.validate.rules) = { ... }</code>
   * @return Whether the header field is set.
   */
  boolean hasHeader();
  /**
   * <pre>
   * Header name/value pair that this option applies to.
   * </pre>
   *
   * <code>.envoy.config.core.v3.HeaderValue header = 1 [json_name = "header", (.validate.rules) = { ... }</code>
   * @return The header.
   */
  com.fluxninja.generated.envoy.config.core.v3.HeaderValue getHeader();
  /**
   * <pre>
   * Header name/value pair that this option applies to.
   * </pre>
   *
   * <code>.envoy.config.core.v3.HeaderValue header = 1 [json_name = "header", (.validate.rules) = { ... }</code>
   */
  com.fluxninja.generated.envoy.config.core.v3.HeaderValueOrBuilder getHeaderOrBuilder();

  /**
   * <pre>
   * Should the value be appended? If true (default), the value is appended to
   * existing values. Otherwise it replaces any existing values.
   * This field is deprecated and please use
   * :ref:`append_action &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append_action&gt;` as replacement.
   * </pre>
   *
   * <code>.google.protobuf.BoolValue append = 2 [json_name = "append", deprecated = true, (.envoy.annotations.deprecated_at_minor_version) = "3.0"];</code>
   * @deprecated envoy.config.core.v3.HeaderValueOption.append is deprecated.
   *     See envoy/config/core/v3/base.proto;l=360
   * @return Whether the append field is set.
   */
  @java.lang.Deprecated boolean hasAppend();
  /**
   * <pre>
   * Should the value be appended? If true (default), the value is appended to
   * existing values. Otherwise it replaces any existing values.
   * This field is deprecated and please use
   * :ref:`append_action &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append_action&gt;` as replacement.
   * </pre>
   *
   * <code>.google.protobuf.BoolValue append = 2 [json_name = "append", deprecated = true, (.envoy.annotations.deprecated_at_minor_version) = "3.0"];</code>
   * @deprecated envoy.config.core.v3.HeaderValueOption.append is deprecated.
   *     See envoy/config/core/v3/base.proto;l=360
   * @return The append.
   */
  @java.lang.Deprecated com.google.protobuf.BoolValue getAppend();
  /**
   * <pre>
   * Should the value be appended? If true (default), the value is appended to
   * existing values. Otherwise it replaces any existing values.
   * This field is deprecated and please use
   * :ref:`append_action &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append_action&gt;` as replacement.
   * </pre>
   *
   * <code>.google.protobuf.BoolValue append = 2 [json_name = "append", deprecated = true, (.envoy.annotations.deprecated_at_minor_version) = "3.0"];</code>
   */
  @java.lang.Deprecated com.google.protobuf.BoolValueOrBuilder getAppendOrBuilder();

  /**
   * <pre>
   * Describes the action taken to append/overwrite the given value for an existing header
   * or to only add this header if it's absent.
   * Value defaults to :ref:`APPEND_IF_EXISTS_OR_ADD
   * &lt;envoy_v3_api_enum_value_config.core.v3.HeaderValueOption.HeaderAppendAction.APPEND_IF_EXISTS_OR_ADD&gt;`.
   * </pre>
   *
   * <code>.envoy.config.core.v3.HeaderValueOption.HeaderAppendAction append_action = 3 [json_name = "appendAction", (.validate.rules) = { ... }</code>
   * @return The enum numeric value on the wire for appendAction.
   */
  int getAppendActionValue();
  /**
   * <pre>
   * Describes the action taken to append/overwrite the given value for an existing header
   * or to only add this header if it's absent.
   * Value defaults to :ref:`APPEND_IF_EXISTS_OR_ADD
   * &lt;envoy_v3_api_enum_value_config.core.v3.HeaderValueOption.HeaderAppendAction.APPEND_IF_EXISTS_OR_ADD&gt;`.
   * </pre>
   *
   * <code>.envoy.config.core.v3.HeaderValueOption.HeaderAppendAction append_action = 3 [json_name = "appendAction", (.validate.rules) = { ... }</code>
   * @return The appendAction.
   */
  com.fluxninja.generated.envoy.config.core.v3.HeaderValueOption.HeaderAppendAction getAppendAction();

  /**
   * <pre>
   * Is the header value allowed to be empty? If false (default), custom headers with empty values are dropped,
   * otherwise they are added.
   * </pre>
   *
   * <code>bool keep_empty_value = 4 [json_name = "keepEmptyValue"];</code>
   * @return The keepEmptyValue.
   */
  boolean getKeepEmptyValue();
}
