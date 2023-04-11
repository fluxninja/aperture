// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: protoc-gen-openapiv2/options/openapiv2.proto

package com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options;

public interface ResponseOrBuilder extends
    // @@protoc_insertion_point(interface_extends:grpc.gateway.protoc_gen_openapiv2.options.Response)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * `Description` is a short description of the response.
   * GFM syntax can be used for rich text representation.
   * </pre>
   *
   * <code>string description = 1 [json_name = "description"];</code>
   * @return The description.
   */
  java.lang.String getDescription();
  /**
   * <pre>
   * `Description` is a short description of the response.
   * GFM syntax can be used for rich text representation.
   * </pre>
   *
   * <code>string description = 1 [json_name = "description"];</code>
   * @return The bytes for description.
   */
  com.google.protobuf.ByteString
      getDescriptionBytes();

  /**
   * <pre>
   * `Schema` optionally defines the structure of the response.
   * If `Schema` is not provided, it means there is no content to the response.
   * </pre>
   *
   * <code>.grpc.gateway.protoc_gen_openapiv2.options.Schema schema = 2 [json_name = "schema"];</code>
   * @return Whether the schema field is set.
   */
  boolean hasSchema();
  /**
   * <pre>
   * `Schema` optionally defines the structure of the response.
   * If `Schema` is not provided, it means there is no content to the response.
   * </pre>
   *
   * <code>.grpc.gateway.protoc_gen_openapiv2.options.Schema schema = 2 [json_name = "schema"];</code>
   * @return The schema.
   */
  com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.Schema getSchema();
  /**
   * <pre>
   * `Schema` optionally defines the structure of the response.
   * If `Schema` is not provided, it means there is no content to the response.
   * </pre>
   *
   * <code>.grpc.gateway.protoc_gen_openapiv2.options.Schema schema = 2 [json_name = "schema"];</code>
   */
  com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.SchemaOrBuilder getSchemaOrBuilder();

  /**
   * <pre>
   * `Headers` A list of headers that are sent with the response.
   * `Header` name is expected to be a string in the canonical format of the MIME header key
   * See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey
   * </pre>
   *
   * <code>map&lt;string, .grpc.gateway.protoc_gen_openapiv2.options.Header&gt; headers = 3 [json_name = "headers"];</code>
   */
  int getHeadersCount();
  /**
   * <pre>
   * `Headers` A list of headers that are sent with the response.
   * `Header` name is expected to be a string in the canonical format of the MIME header key
   * See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey
   * </pre>
   *
   * <code>map&lt;string, .grpc.gateway.protoc_gen_openapiv2.options.Header&gt; headers = 3 [json_name = "headers"];</code>
   */
  boolean containsHeaders(
      java.lang.String key);
  /**
   * Use {@link #getHeadersMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.Header>
  getHeaders();
  /**
   * <pre>
   * `Headers` A list of headers that are sent with the response.
   * `Header` name is expected to be a string in the canonical format of the MIME header key
   * See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey
   * </pre>
   *
   * <code>map&lt;string, .grpc.gateway.protoc_gen_openapiv2.options.Header&gt; headers = 3 [json_name = "headers"];</code>
   */
  java.util.Map<java.lang.String, com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.Header>
  getHeadersMap();
  /**
   * <pre>
   * `Headers` A list of headers that are sent with the response.
   * `Header` name is expected to be a string in the canonical format of the MIME header key
   * See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey
   * </pre>
   *
   * <code>map&lt;string, .grpc.gateway.protoc_gen_openapiv2.options.Header&gt; headers = 3 [json_name = "headers"];</code>
   */

  /* nullable */
com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.Header getHeadersOrDefault(
      java.lang.String key,
      /* nullable */
com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.Header defaultValue);
  /**
   * <pre>
   * `Headers` A list of headers that are sent with the response.
   * `Header` name is expected to be a string in the canonical format of the MIME header key
   * See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey
   * </pre>
   *
   * <code>map&lt;string, .grpc.gateway.protoc_gen_openapiv2.options.Header&gt; headers = 3 [json_name = "headers"];</code>
   */

  com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.Header getHeadersOrThrow(
      java.lang.String key);

  /**
   * <pre>
   * `Examples` gives per-mimetype response examples.
   * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object
   * </pre>
   *
   * <code>map&lt;string, string&gt; examples = 4 [json_name = "examples"];</code>
   */
  int getExamplesCount();
  /**
   * <pre>
   * `Examples` gives per-mimetype response examples.
   * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object
   * </pre>
   *
   * <code>map&lt;string, string&gt; examples = 4 [json_name = "examples"];</code>
   */
  boolean containsExamples(
      java.lang.String key);
  /**
   * Use {@link #getExamplesMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, java.lang.String>
  getExamples();
  /**
   * <pre>
   * `Examples` gives per-mimetype response examples.
   * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object
   * </pre>
   *
   * <code>map&lt;string, string&gt; examples = 4 [json_name = "examples"];</code>
   */
  java.util.Map<java.lang.String, java.lang.String>
  getExamplesMap();
  /**
   * <pre>
   * `Examples` gives per-mimetype response examples.
   * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object
   * </pre>
   *
   * <code>map&lt;string, string&gt; examples = 4 [json_name = "examples"];</code>
   */

  /* nullable */
java.lang.String getExamplesOrDefault(
      java.lang.String key,
      /* nullable */
java.lang.String defaultValue);
  /**
   * <pre>
   * `Examples` gives per-mimetype response examples.
   * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object
   * </pre>
   *
   * <code>map&lt;string, string&gt; examples = 4 [json_name = "examples"];</code>
   */

  java.lang.String getExamplesOrThrow(
      java.lang.String key);

  /**
   * <pre>
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   * </pre>
   *
   * <code>map&lt;string, .google.protobuf.Value&gt; extensions = 5 [json_name = "extensions"];</code>
   */
  int getExtensionsCount();
  /**
   * <pre>
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   * </pre>
   *
   * <code>map&lt;string, .google.protobuf.Value&gt; extensions = 5 [json_name = "extensions"];</code>
   */
  boolean containsExtensions(
      java.lang.String key);
  /**
   * Use {@link #getExtensionsMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, com.google.protobuf.Value>
  getExtensions();
  /**
   * <pre>
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   * </pre>
   *
   * <code>map&lt;string, .google.protobuf.Value&gt; extensions = 5 [json_name = "extensions"];</code>
   */
  java.util.Map<java.lang.String, com.google.protobuf.Value>
  getExtensionsMap();
  /**
   * <pre>
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   * </pre>
   *
   * <code>map&lt;string, .google.protobuf.Value&gt; extensions = 5 [json_name = "extensions"];</code>
   */

  /* nullable */
com.google.protobuf.Value getExtensionsOrDefault(
      java.lang.String key,
      /* nullable */
com.google.protobuf.Value defaultValue);
  /**
   * <pre>
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   * </pre>
   *
   * <code>map&lt;string, .google.protobuf.Value&gt; extensions = 5 [json_name = "extensions"];</code>
   */

  com.google.protobuf.Value getExtensionsOrThrow(
      java.lang.String key);
}
