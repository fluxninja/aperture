// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/classification/v1/extractor.proto

package com.aperture.classification.v1;

public interface AddressExtractorOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.classification.v1.AddressExtractor)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Attribute path pointing to some string - eg. "source.address".
   * </pre>
   *
   * <code>string from = 1 [json_name = "from", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The from.
   */
  java.lang.String getFrom();
  /**
   * <pre>
   * Attribute path pointing to some string - eg. "source.address".
   * </pre>
   *
   * <code>string from = 1 [json_name = "from", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The bytes for from.
   */
  com.google.protobuf.ByteString
      getFromBytes();
}
