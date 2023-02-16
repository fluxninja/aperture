// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: protoc-gen-openapiv2/options/openapiv2.proto

package com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options;

public interface ExternalDocumentationOrBuilder extends
    // @@protoc_insertion_point(interface_extends:grpc.gateway.protoc_gen_openapiv2.options.ExternalDocumentation)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * A short description of the target documentation. GFM syntax can be used for
   * rich text representation.
   * </pre>
   *
   * <code>string description = 1 [json_name = "description"];</code>
   * @return The description.
   */
  java.lang.String getDescription();
  /**
   * <pre>
   * A short description of the target documentation. GFM syntax can be used for
   * rich text representation.
   * </pre>
   *
   * <code>string description = 1 [json_name = "description"];</code>
   * @return The bytes for description.
   */
  com.google.protobuf.ByteString
      getDescriptionBytes();

  /**
   * <pre>
   * The URL for the target documentation. Value MUST be in the format
   * of a URL.
   * </pre>
   *
   * <code>string url = 2 [json_name = "url"];</code>
   * @return The url.
   */
  java.lang.String getUrl();
  /**
   * <pre>
   * The URL for the target documentation. Value MUST be in the format
   * of a URL.
   * </pre>
   *
   * <code>string url = 2 [json_name = "url"];</code>
   * @return The bytes for url.
   */
  com.google.protobuf.ByteString
      getUrlBytes();
}
