// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/checkhttp/v1/checkhttp.proto

package com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1;

/**
 * <pre>
 * HTTP attributes for a denied response.
 * </pre>
 *
 * Protobuf type {@code aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse}
 */
public final class DeniedHttpResponse extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse)
    DeniedHttpResponseOrBuilder {
private static final long serialVersionUID = 0L;
  // Use DeniedHttpResponse.newBuilder() to construct.
  private DeniedHttpResponse(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private DeniedHttpResponse() {
    body_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new DeniedHttpResponse();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private DeniedHttpResponse(
      com.google.protobuf.CodedInputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    this();
    if (extensionRegistry == null) {
      throw new java.lang.NullPointerException();
    }
    int mutable_bitField0_ = 0;
    com.google.protobuf.UnknownFieldSet.Builder unknownFields =
        com.google.protobuf.UnknownFieldSet.newBuilder();
    try {
      boolean done = false;
      while (!done) {
        int tag = input.readTag();
        switch (tag) {
          case 0:
            done = true;
            break;
          case 8: {

            status_ = input.readInt32();
            break;
          }
          case 18: {
            if (!((mutable_bitField0_ & 0x00000001) != 0)) {
              headers_ = com.google.protobuf.MapField.newMapField(
                  HeadersDefaultEntryHolder.defaultEntry);
              mutable_bitField0_ |= 0x00000001;
            }
            com.google.protobuf.MapEntry<java.lang.String, java.lang.String>
            headers__ = input.readMessage(
                HeadersDefaultEntryHolder.defaultEntry.getParserForType(), extensionRegistry);
            headers_.getMutableMap().put(
                headers__.getKey(), headers__.getValue());
            break;
          }
          case 26: {
            java.lang.String s = input.readStringRequireUtf8();

            body_ = s;
            break;
          }
          default: {
            if (!parseUnknownField(
                input, unknownFields, extensionRegistry, tag)) {
              done = true;
            }
            break;
          }
        }
      }
    } catch (com.google.protobuf.InvalidProtocolBufferException e) {
      throw e.setUnfinishedMessage(this);
    } catch (com.google.protobuf.UninitializedMessageException e) {
      throw e.asInvalidProtocolBufferException().setUnfinishedMessage(this);
    } catch (java.io.IOException e) {
      throw new com.google.protobuf.InvalidProtocolBufferException(
          e).setUnfinishedMessage(this);
    } finally {
      this.unknownFields = unknownFields.build();
      makeExtensionsImmutable();
    }
  }
  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckhttpProto.internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_descriptor;
  }

  @SuppressWarnings({"rawtypes"})
  @java.lang.Override
  protected com.google.protobuf.MapField internalGetMapField(
      int number) {
    switch (number) {
      case 2:
        return internalGetHeaders();
      default:
        throw new RuntimeException(
            "Invalid map field number: " + number);
    }
  }
  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckhttpProto.internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse.class, com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse.Builder.class);
  }

  public static final int STATUS_FIELD_NUMBER = 1;
  private int status_;
  /**
   * <pre>
   * This field allows the authorization service to send an HTTP response status code to the
   * downstream client. If not set, Envoy sends `403 Forbidden` HTTP status code by default.
   * </pre>
   *
   * <code>int32 status = 1 [json_name = "status"];</code>
   * @return The status.
   */
  @java.lang.Override
  public int getStatus() {
    return status_;
  }

  public static final int HEADERS_FIELD_NUMBER = 2;
  private static final class HeadersDefaultEntryHolder {
    static final com.google.protobuf.MapEntry<
        java.lang.String, java.lang.String> defaultEntry =
            com.google.protobuf.MapEntry
            .<java.lang.String, java.lang.String>newDefaultInstance(
                com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckhttpProto.internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_HeadersEntry_descriptor, 
                com.google.protobuf.WireFormat.FieldType.STRING,
                "",
                com.google.protobuf.WireFormat.FieldType.STRING,
                "");
  }
  private com.google.protobuf.MapField<
      java.lang.String, java.lang.String> headers_;
  private com.google.protobuf.MapField<java.lang.String, java.lang.String>
  internalGetHeaders() {
    if (headers_ == null) {
      return com.google.protobuf.MapField.emptyMapField(
          HeadersDefaultEntryHolder.defaultEntry);
    }
    return headers_;
  }

  public int getHeadersCount() {
    return internalGetHeaders().getMap().size();
  }
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client.
   * </pre>
   *
   * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
   */

  @java.lang.Override
  public boolean containsHeaders(
      java.lang.String key) {
    if (key == null) { throw new NullPointerException("map key"); }
    return internalGetHeaders().getMap().containsKey(key);
  }
  /**
   * Use {@link #getHeadersMap()} instead.
   */
  @java.lang.Override
  @java.lang.Deprecated
  public java.util.Map<java.lang.String, java.lang.String> getHeaders() {
    return getHeadersMap();
  }
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client.
   * </pre>
   *
   * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
   */
  @java.lang.Override

  public java.util.Map<java.lang.String, java.lang.String> getHeadersMap() {
    return internalGetHeaders().getMap();
  }
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client.
   * </pre>
   *
   * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
   */
  @java.lang.Override

  public java.lang.String getHeadersOrDefault(
      java.lang.String key,
      java.lang.String defaultValue) {
    if (key == null) { throw new NullPointerException("map key"); }
    java.util.Map<java.lang.String, java.lang.String> map =
        internalGetHeaders().getMap();
    return map.containsKey(key) ? map.get(key) : defaultValue;
  }
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client.
   * </pre>
   *
   * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
   */
  @java.lang.Override

  public java.lang.String getHeadersOrThrow(
      java.lang.String key) {
    if (key == null) { throw new NullPointerException("map key"); }
    java.util.Map<java.lang.String, java.lang.String> map =
        internalGetHeaders().getMap();
    if (!map.containsKey(key)) {
      throw new java.lang.IllegalArgumentException();
    }
    return map.get(key);
  }

  public static final int BODY_FIELD_NUMBER = 3;
  private volatile java.lang.Object body_;
  /**
   * <pre>
   * This field allows the authorization service to send a response body data
   * to the downstream client.
   * </pre>
   *
   * <code>string body = 3 [json_name = "body"];</code>
   * @return The body.
   */
  @java.lang.Override
  public java.lang.String getBody() {
    java.lang.Object ref = body_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      body_ = s;
      return s;
    }
  }
  /**
   * <pre>
   * This field allows the authorization service to send a response body data
   * to the downstream client.
   * </pre>
   *
   * <code>string body = 3 [json_name = "body"];</code>
   * @return The bytes for body.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getBodyBytes() {
    java.lang.Object ref = body_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      body_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  private byte memoizedIsInitialized = -1;
  @java.lang.Override
  public final boolean isInitialized() {
    byte isInitialized = memoizedIsInitialized;
    if (isInitialized == 1) return true;
    if (isInitialized == 0) return false;

    memoizedIsInitialized = 1;
    return true;
  }

  @java.lang.Override
  public void writeTo(com.google.protobuf.CodedOutputStream output)
                      throws java.io.IOException {
    if (status_ != 0) {
      output.writeInt32(1, status_);
    }
    com.google.protobuf.GeneratedMessageV3
      .serializeStringMapTo(
        output,
        internalGetHeaders(),
        HeadersDefaultEntryHolder.defaultEntry,
        2);
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(body_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 3, body_);
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (status_ != 0) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt32Size(1, status_);
    }
    for (java.util.Map.Entry<java.lang.String, java.lang.String> entry
         : internalGetHeaders().getMap().entrySet()) {
      com.google.protobuf.MapEntry<java.lang.String, java.lang.String>
      headers__ = HeadersDefaultEntryHolder.defaultEntry.newBuilderForType()
          .setKey(entry.getKey())
          .setValue(entry.getValue())
          .build();
      size += com.google.protobuf.CodedOutputStream
          .computeMessageSize(2, headers__);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(body_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(3, body_);
    }
    size += unknownFields.getSerializedSize();
    memoizedSize = size;
    return size;
  }

  @java.lang.Override
  public boolean equals(final java.lang.Object obj) {
    if (obj == this) {
     return true;
    }
    if (!(obj instanceof com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse other = (com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse) obj;

    if (getStatus()
        != other.getStatus()) return false;
    if (!internalGetHeaders().equals(
        other.internalGetHeaders())) return false;
    if (!getBody()
        .equals(other.getBody())) return false;
    if (!unknownFields.equals(other.unknownFields)) return false;
    return true;
  }

  @java.lang.Override
  public int hashCode() {
    if (memoizedHashCode != 0) {
      return memoizedHashCode;
    }
    int hash = 41;
    hash = (19 * hash) + getDescriptor().hashCode();
    hash = (37 * hash) + STATUS_FIELD_NUMBER;
    hash = (53 * hash) + getStatus();
    if (!internalGetHeaders().getMap().isEmpty()) {
      hash = (37 * hash) + HEADERS_FIELD_NUMBER;
      hash = (53 * hash) + internalGetHeaders().hashCode();
    }
    hash = (37 * hash) + BODY_FIELD_NUMBER;
    hash = (53 * hash) + getBody().hashCode();
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parseFrom(
      com.google.protobuf.CodedInputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  @java.lang.Override
  public Builder newBuilderForType() { return newBuilder(); }
  public static Builder newBuilder() {
    return DEFAULT_INSTANCE.toBuilder();
  }
  public static Builder newBuilder(com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse prototype) {
    return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
  }
  @java.lang.Override
  public Builder toBuilder() {
    return this == DEFAULT_INSTANCE
        ? new Builder() : new Builder().mergeFrom(this);
  }

  @java.lang.Override
  protected Builder newBuilderForType(
      com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
    Builder builder = new Builder(parent);
    return builder;
  }
  /**
   * <pre>
   * HTTP attributes for a denied response.
   * </pre>
   *
   * Protobuf type {@code aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse)
      com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponseOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckhttpProto.internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_descriptor;
    }

    @SuppressWarnings({"rawtypes"})
    protected com.google.protobuf.MapField internalGetMapField(
        int number) {
      switch (number) {
        case 2:
          return internalGetHeaders();
        default:
          throw new RuntimeException(
              "Invalid map field number: " + number);
      }
    }
    @SuppressWarnings({"rawtypes"})
    protected com.google.protobuf.MapField internalGetMutableMapField(
        int number) {
      switch (number) {
        case 2:
          return internalGetMutableHeaders();
        default:
          throw new RuntimeException(
              "Invalid map field number: " + number);
      }
    }
    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckhttpProto.internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse.class, com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse.newBuilder()
    private Builder() {
      maybeForceBuilderInitialization();
    }

    private Builder(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      super(parent);
      maybeForceBuilderInitialization();
    }
    private void maybeForceBuilderInitialization() {
      if (com.google.protobuf.GeneratedMessageV3
              .alwaysUseFieldBuilders) {
      }
    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      status_ = 0;

      internalGetMutableHeaders().clear();
      body_ = "";

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckhttpProto.internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse build() {
      com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse buildPartial() {
      com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse result = new com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse(this);
      int from_bitField0_ = bitField0_;
      result.status_ = status_;
      result.headers_ = internalGetHeaders();
      result.headers_.makeImmutable();
      result.body_ = body_;
      onBuilt();
      return result;
    }

    @java.lang.Override
    public Builder clone() {
      return super.clone();
    }
    @java.lang.Override
    public Builder setField(
        com.google.protobuf.Descriptors.FieldDescriptor field,
        java.lang.Object value) {
      return super.setField(field, value);
    }
    @java.lang.Override
    public Builder clearField(
        com.google.protobuf.Descriptors.FieldDescriptor field) {
      return super.clearField(field);
    }
    @java.lang.Override
    public Builder clearOneof(
        com.google.protobuf.Descriptors.OneofDescriptor oneof) {
      return super.clearOneof(oneof);
    }
    @java.lang.Override
    public Builder setRepeatedField(
        com.google.protobuf.Descriptors.FieldDescriptor field,
        int index, java.lang.Object value) {
      return super.setRepeatedField(field, index, value);
    }
    @java.lang.Override
    public Builder addRepeatedField(
        com.google.protobuf.Descriptors.FieldDescriptor field,
        java.lang.Object value) {
      return super.addRepeatedField(field, value);
    }
    @java.lang.Override
    public Builder mergeFrom(com.google.protobuf.Message other) {
      if (other instanceof com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse) {
        return mergeFrom((com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse other) {
      if (other == com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse.getDefaultInstance()) return this;
      if (other.getStatus() != 0) {
        setStatus(other.getStatus());
      }
      internalGetMutableHeaders().mergeFrom(
          other.internalGetHeaders());
      if (!other.getBody().isEmpty()) {
        body_ = other.body_;
        onChanged();
      }
      this.mergeUnknownFields(other.unknownFields);
      onChanged();
      return this;
    }

    @java.lang.Override
    public final boolean isInitialized() {
      return true;
    }

    @java.lang.Override
    public Builder mergeFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }
    private int bitField0_;

    private int status_ ;
    /**
     * <pre>
     * This field allows the authorization service to send an HTTP response status code to the
     * downstream client. If not set, Envoy sends `403 Forbidden` HTTP status code by default.
     * </pre>
     *
     * <code>int32 status = 1 [json_name = "status"];</code>
     * @return The status.
     */
    @java.lang.Override
    public int getStatus() {
      return status_;
    }
    /**
     * <pre>
     * This field allows the authorization service to send an HTTP response status code to the
     * downstream client. If not set, Envoy sends `403 Forbidden` HTTP status code by default.
     * </pre>
     *
     * <code>int32 status = 1 [json_name = "status"];</code>
     * @param value The status to set.
     * @return This builder for chaining.
     */
    public Builder setStatus(int value) {
      
      status_ = value;
      onChanged();
      return this;
    }
    /**
     * <pre>
     * This field allows the authorization service to send an HTTP response status code to the
     * downstream client. If not set, Envoy sends `403 Forbidden` HTTP status code by default.
     * </pre>
     *
     * <code>int32 status = 1 [json_name = "status"];</code>
     * @return This builder for chaining.
     */
    public Builder clearStatus() {
      
      status_ = 0;
      onChanged();
      return this;
    }

    private com.google.protobuf.MapField<
        java.lang.String, java.lang.String> headers_;
    private com.google.protobuf.MapField<java.lang.String, java.lang.String>
    internalGetHeaders() {
      if (headers_ == null) {
        return com.google.protobuf.MapField.emptyMapField(
            HeadersDefaultEntryHolder.defaultEntry);
      }
      return headers_;
    }
    private com.google.protobuf.MapField<java.lang.String, java.lang.String>
    internalGetMutableHeaders() {
      onChanged();;
      if (headers_ == null) {
        headers_ = com.google.protobuf.MapField.newMapField(
            HeadersDefaultEntryHolder.defaultEntry);
      }
      if (!headers_.isMutable()) {
        headers_ = headers_.copy();
      }
      return headers_;
    }

    public int getHeadersCount() {
      return internalGetHeaders().getMap().size();
    }
    /**
     * <pre>
     * This field allows the authorization service to send HTTP response headers
     * to the downstream client.
     * </pre>
     *
     * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
     */

    @java.lang.Override
    public boolean containsHeaders(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      return internalGetHeaders().getMap().containsKey(key);
    }
    /**
     * Use {@link #getHeadersMap()} instead.
     */
    @java.lang.Override
    @java.lang.Deprecated
    public java.util.Map<java.lang.String, java.lang.String> getHeaders() {
      return getHeadersMap();
    }
    /**
     * <pre>
     * This field allows the authorization service to send HTTP response headers
     * to the downstream client.
     * </pre>
     *
     * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
     */
    @java.lang.Override

    public java.util.Map<java.lang.String, java.lang.String> getHeadersMap() {
      return internalGetHeaders().getMap();
    }
    /**
     * <pre>
     * This field allows the authorization service to send HTTP response headers
     * to the downstream client.
     * </pre>
     *
     * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
     */
    @java.lang.Override

    public java.lang.String getHeadersOrDefault(
        java.lang.String key,
        java.lang.String defaultValue) {
      if (key == null) { throw new NullPointerException("map key"); }
      java.util.Map<java.lang.String, java.lang.String> map =
          internalGetHeaders().getMap();
      return map.containsKey(key) ? map.get(key) : defaultValue;
    }
    /**
     * <pre>
     * This field allows the authorization service to send HTTP response headers
     * to the downstream client.
     * </pre>
     *
     * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
     */
    @java.lang.Override

    public java.lang.String getHeadersOrThrow(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      java.util.Map<java.lang.String, java.lang.String> map =
          internalGetHeaders().getMap();
      if (!map.containsKey(key)) {
        throw new java.lang.IllegalArgumentException();
      }
      return map.get(key);
    }

    public Builder clearHeaders() {
      internalGetMutableHeaders().getMutableMap()
          .clear();
      return this;
    }
    /**
     * <pre>
     * This field allows the authorization service to send HTTP response headers
     * to the downstream client.
     * </pre>
     *
     * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
     */

    public Builder removeHeaders(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      internalGetMutableHeaders().getMutableMap()
          .remove(key);
      return this;
    }
    /**
     * Use alternate mutation accessors instead.
     */
    @java.lang.Deprecated
    public java.util.Map<java.lang.String, java.lang.String>
    getMutableHeaders() {
      return internalGetMutableHeaders().getMutableMap();
    }
    /**
     * <pre>
     * This field allows the authorization service to send HTTP response headers
     * to the downstream client.
     * </pre>
     *
     * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
     */
    public Builder putHeaders(
        java.lang.String key,
        java.lang.String value) {
      if (key == null) { throw new NullPointerException("map key"); }
      if (value == null) {
  throw new NullPointerException("map value");
}

      internalGetMutableHeaders().getMutableMap()
          .put(key, value);
      return this;
    }
    /**
     * <pre>
     * This field allows the authorization service to send HTTP response headers
     * to the downstream client.
     * </pre>
     *
     * <code>map&lt;string, string&gt; headers = 2 [json_name = "headers"];</code>
     */

    public Builder putAllHeaders(
        java.util.Map<java.lang.String, java.lang.String> values) {
      internalGetMutableHeaders().getMutableMap()
          .putAll(values);
      return this;
    }

    private java.lang.Object body_ = "";
    /**
     * <pre>
     * This field allows the authorization service to send a response body data
     * to the downstream client.
     * </pre>
     *
     * <code>string body = 3 [json_name = "body"];</code>
     * @return The body.
     */
    public java.lang.String getBody() {
      java.lang.Object ref = body_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        body_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <pre>
     * This field allows the authorization service to send a response body data
     * to the downstream client.
     * </pre>
     *
     * <code>string body = 3 [json_name = "body"];</code>
     * @return The bytes for body.
     */
    public com.google.protobuf.ByteString
        getBodyBytes() {
      java.lang.Object ref = body_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        body_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <pre>
     * This field allows the authorization service to send a response body data
     * to the downstream client.
     * </pre>
     *
     * <code>string body = 3 [json_name = "body"];</code>
     * @param value The body to set.
     * @return This builder for chaining.
     */
    public Builder setBody(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      body_ = value;
      onChanged();
      return this;
    }
    /**
     * <pre>
     * This field allows the authorization service to send a response body data
     * to the downstream client.
     * </pre>
     *
     * <code>string body = 3 [json_name = "body"];</code>
     * @return This builder for chaining.
     */
    public Builder clearBody() {
      
      body_ = getDefaultInstance().getBody();
      onChanged();
      return this;
    }
    /**
     * <pre>
     * This field allows the authorization service to send a response body data
     * to the downstream client.
     * </pre>
     *
     * <code>string body = 3 [json_name = "body"];</code>
     * @param value The bytes for body to set.
     * @return This builder for chaining.
     */
    public Builder setBodyBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      body_ = value;
      onChanged();
      return this;
    }
    @java.lang.Override
    public final Builder setUnknownFields(
        final com.google.protobuf.UnknownFieldSet unknownFields) {
      return super.setUnknownFields(unknownFields);
    }

    @java.lang.Override
    public final Builder mergeUnknownFields(
        final com.google.protobuf.UnknownFieldSet unknownFields) {
      return super.mergeUnknownFields(unknownFields);
    }


    // @@protoc_insertion_point(builder_scope:aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse)
  }

  // @@protoc_insertion_point(class_scope:aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse)
  private static final com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse();
  }

  public static com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<DeniedHttpResponse>
      PARSER = new com.google.protobuf.AbstractParser<DeniedHttpResponse>() {
    @java.lang.Override
    public DeniedHttpResponse parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new DeniedHttpResponse(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<DeniedHttpResponse> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<DeniedHttpResponse> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.DeniedHttpResponse getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

