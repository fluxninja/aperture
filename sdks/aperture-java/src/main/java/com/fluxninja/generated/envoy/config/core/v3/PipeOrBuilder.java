// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/config/core/v3/address.proto

package com.fluxninja.generated.envoy.config.core.v3;

public interface PipeOrBuilder extends
    // @@protoc_insertion_point(interface_extends:envoy.config.core.v3.Pipe)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Unix Domain Socket path. On Linux, paths starting with '&#64;' will use the
   * abstract namespace. The starting '&#64;' is replaced by a null byte by Envoy.
   * Paths starting with '&#64;' will result in an error in environments other than
   * Linux.
   * </pre>
   *
   * <code>string path = 1 [json_name = "path", (.validate.rules) = { ... }</code>
   * @return The path.
   */
  java.lang.String getPath();
  /**
   * <pre>
   * Unix Domain Socket path. On Linux, paths starting with '&#64;' will use the
   * abstract namespace. The starting '&#64;' is replaced by a null byte by Envoy.
   * Paths starting with '&#64;' will result in an error in environments other than
   * Linux.
   * </pre>
   *
   * <code>string path = 1 [json_name = "path", (.validate.rules) = { ... }</code>
   * @return The bytes for path.
   */
  com.google.protobuf.ByteString
      getPathBytes();

  /**
   * <pre>
   * The mode for the Pipe. Not applicable for abstract sockets.
   * </pre>
   *
   * <code>uint32 mode = 2 [json_name = "mode", (.validate.rules) = { ... }</code>
   * @return The mode.
   */
  int getMode();
}
