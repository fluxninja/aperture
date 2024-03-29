// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

/**
 * Protobuf type {@code aperture.flowcontrol.check.v1.CacheDeleteResponse}
 */
public final class CacheDeleteResponse extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.flowcontrol.check.v1.CacheDeleteResponse)
    CacheDeleteResponseOrBuilder {
private static final long serialVersionUID = 0L;
  // Use CacheDeleteResponse.newBuilder() to construct.
  private CacheDeleteResponse(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private CacheDeleteResponse() {
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new CacheDeleteResponse();
  }

  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_descriptor;
  }

  @SuppressWarnings({"rawtypes"})
  @java.lang.Override
  protected com.google.protobuf.MapField internalGetMapField(
      int number) {
    switch (number) {
      case 2:
        return internalGetGlobalCacheResponses();
      default:
        throw new RuntimeException(
            "Invalid map field number: " + number);
    }
  }
  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.class, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.Builder.class);
  }

  private int bitField0_;
  public static final int RESULT_CACHE_RESPONSE_FIELD_NUMBER = 1;
  private com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse resultCacheResponse_;
  /**
   * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
   * @return Whether the resultCacheResponse field is set.
   */
  @java.lang.Override
  public boolean hasResultCacheResponse() {
    return ((bitField0_ & 0x00000001) != 0);
  }
  /**
   * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
   * @return The resultCacheResponse.
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse getResultCacheResponse() {
    return resultCacheResponse_ == null ? com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.getDefaultInstance() : resultCacheResponse_;
  }
  /**
   * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponseOrBuilder getResultCacheResponseOrBuilder() {
    return resultCacheResponse_ == null ? com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.getDefaultInstance() : resultCacheResponse_;
  }

  public static final int GLOBAL_CACHE_RESPONSES_FIELD_NUMBER = 2;
  private static final class GlobalCacheResponsesDefaultEntryHolder {
    static final com.google.protobuf.MapEntry<
        java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> defaultEntry =
            com.google.protobuf.MapEntry
            .<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse>newDefaultInstance(
                com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_GlobalCacheResponsesEntry_descriptor, 
                com.google.protobuf.WireFormat.FieldType.STRING,
                "",
                com.google.protobuf.WireFormat.FieldType.MESSAGE,
                com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.getDefaultInstance());
  }
  @SuppressWarnings("serial")
  private com.google.protobuf.MapField<
      java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> globalCacheResponses_;
  private com.google.protobuf.MapField<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse>
  internalGetGlobalCacheResponses() {
    if (globalCacheResponses_ == null) {
      return com.google.protobuf.MapField.emptyMapField(
          GlobalCacheResponsesDefaultEntryHolder.defaultEntry);
    }
    return globalCacheResponses_;
  }
  public int getGlobalCacheResponsesCount() {
    return internalGetGlobalCacheResponses().getMap().size();
  }
  /**
   * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
   */
  @java.lang.Override
  public boolean containsGlobalCacheResponses(
      java.lang.String key) {
    if (key == null) { throw new NullPointerException("map key"); }
    return internalGetGlobalCacheResponses().getMap().containsKey(key);
  }
  /**
   * Use {@link #getGlobalCacheResponsesMap()} instead.
   */
  @java.lang.Override
  @java.lang.Deprecated
  public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> getGlobalCacheResponses() {
    return getGlobalCacheResponsesMap();
  }
  /**
   * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
   */
  @java.lang.Override
  public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> getGlobalCacheResponsesMap() {
    return internalGetGlobalCacheResponses().getMap();
  }
  /**
   * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
   */
  @java.lang.Override
  public /* nullable */
com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse getGlobalCacheResponsesOrDefault(
      java.lang.String key,
      /* nullable */
com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse defaultValue) {
    if (key == null) { throw new NullPointerException("map key"); }
    java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> map =
        internalGetGlobalCacheResponses().getMap();
    return map.containsKey(key) ? map.get(key) : defaultValue;
  }
  /**
   * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse getGlobalCacheResponsesOrThrow(
      java.lang.String key) {
    if (key == null) { throw new NullPointerException("map key"); }
    java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> map =
        internalGetGlobalCacheResponses().getMap();
    if (!map.containsKey(key)) {
      throw new java.lang.IllegalArgumentException();
    }
    return map.get(key);
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
    if (((bitField0_ & 0x00000001) != 0)) {
      output.writeMessage(1, getResultCacheResponse());
    }
    com.google.protobuf.GeneratedMessageV3
      .serializeStringMapTo(
        output,
        internalGetGlobalCacheResponses(),
        GlobalCacheResponsesDefaultEntryHolder.defaultEntry,
        2);
    getUnknownFields().writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (((bitField0_ & 0x00000001) != 0)) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(1, getResultCacheResponse());
    }
    for (java.util.Map.Entry<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> entry
         : internalGetGlobalCacheResponses().getMap().entrySet()) {
      com.google.protobuf.MapEntry<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse>
      globalCacheResponses__ = GlobalCacheResponsesDefaultEntryHolder.defaultEntry.newBuilderForType()
          .setKey(entry.getKey())
          .setValue(entry.getValue())
          .build();
      size += com.google.protobuf.CodedOutputStream
          .computeMessageSize(2, globalCacheResponses__);
    }
    size += getUnknownFields().getSerializedSize();
    memoizedSize = size;
    return size;
  }

  @java.lang.Override
  public boolean equals(final java.lang.Object obj) {
    if (obj == this) {
     return true;
    }
    if (!(obj instanceof com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse other = (com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse) obj;

    if (hasResultCacheResponse() != other.hasResultCacheResponse()) return false;
    if (hasResultCacheResponse()) {
      if (!getResultCacheResponse()
          .equals(other.getResultCacheResponse())) return false;
    }
    if (!internalGetGlobalCacheResponses().equals(
        other.internalGetGlobalCacheResponses())) return false;
    if (!getUnknownFields().equals(other.getUnknownFields())) return false;
    return true;
  }

  @java.lang.Override
  public int hashCode() {
    if (memoizedHashCode != 0) {
      return memoizedHashCode;
    }
    int hash = 41;
    hash = (19 * hash) + getDescriptor().hashCode();
    if (hasResultCacheResponse()) {
      hash = (37 * hash) + RESULT_CACHE_RESPONSE_FIELD_NUMBER;
      hash = (53 * hash) + getResultCacheResponse().hashCode();
    }
    if (!internalGetGlobalCacheResponses().getMap().isEmpty()) {
      hash = (37 * hash) + GLOBAL_CACHE_RESPONSES_FIELD_NUMBER;
      hash = (53 * hash) + internalGetGlobalCacheResponses().hashCode();
    }
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse prototype) {
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
   * Protobuf type {@code aperture.flowcontrol.check.v1.CacheDeleteResponse}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.flowcontrol.check.v1.CacheDeleteResponse)
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponseOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_descriptor;
    }

    @SuppressWarnings({"rawtypes"})
    protected com.google.protobuf.MapField internalGetMapField(
        int number) {
      switch (number) {
        case 2:
          return internalGetGlobalCacheResponses();
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
          return internalGetMutableGlobalCacheResponses();
        default:
          throw new RuntimeException(
              "Invalid map field number: " + number);
      }
    }
    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.class, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.newBuilder()
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
        getResultCacheResponseFieldBuilder();
      }
    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      bitField0_ = 0;
      resultCacheResponse_ = null;
      if (resultCacheResponseBuilder_ != null) {
        resultCacheResponseBuilder_.dispose();
        resultCacheResponseBuilder_ = null;
      }
      internalGetMutableGlobalCacheResponses().clear();
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheDeleteResponse_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse build() {
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse buildPartial() {
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse result = new com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse(this);
      if (bitField0_ != 0) { buildPartial0(result); }
      onBuilt();
      return result;
    }

    private void buildPartial0(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse result) {
      int from_bitField0_ = bitField0_;
      int to_bitField0_ = 0;
      if (((from_bitField0_ & 0x00000001) != 0)) {
        result.resultCacheResponse_ = resultCacheResponseBuilder_ == null
            ? resultCacheResponse_
            : resultCacheResponseBuilder_.build();
        to_bitField0_ |= 0x00000001;
      }
      if (((from_bitField0_ & 0x00000002) != 0)) {
        result.globalCacheResponses_ = internalGetGlobalCacheResponses();
        result.globalCacheResponses_.makeImmutable();
      }
      result.bitField0_ |= to_bitField0_;
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
      if (other instanceof com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse) {
        return mergeFrom((com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse other) {
      if (other == com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.getDefaultInstance()) return this;
      if (other.hasResultCacheResponse()) {
        mergeResultCacheResponse(other.getResultCacheResponse());
      }
      internalGetMutableGlobalCacheResponses().mergeFrom(
          other.internalGetGlobalCacheResponses());
      bitField0_ |= 0x00000002;
      this.mergeUnknownFields(other.getUnknownFields());
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
      if (extensionRegistry == null) {
        throw new java.lang.NullPointerException();
      }
      try {
        boolean done = false;
        while (!done) {
          int tag = input.readTag();
          switch (tag) {
            case 0:
              done = true;
              break;
            case 10: {
              input.readMessage(
                  getResultCacheResponseFieldBuilder().getBuilder(),
                  extensionRegistry);
              bitField0_ |= 0x00000001;
              break;
            } // case 10
            case 18: {
              com.google.protobuf.MapEntry<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse>
              globalCacheResponses__ = input.readMessage(
                  GlobalCacheResponsesDefaultEntryHolder.defaultEntry.getParserForType(), extensionRegistry);
              internalGetMutableGlobalCacheResponses().getMutableMap().put(
                  globalCacheResponses__.getKey(), globalCacheResponses__.getValue());
              bitField0_ |= 0x00000002;
              break;
            } // case 18
            default: {
              if (!super.parseUnknownField(input, extensionRegistry, tag)) {
                done = true; // was an endgroup tag
              }
              break;
            } // default:
          } // switch (tag)
        } // while (!done)
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.unwrapIOException();
      } finally {
        onChanged();
      } // finally
      return this;
    }
    private int bitField0_;

    private com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse resultCacheResponse_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.Builder, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponseOrBuilder> resultCacheResponseBuilder_;
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     * @return Whether the resultCacheResponse field is set.
     */
    public boolean hasResultCacheResponse() {
      return ((bitField0_ & 0x00000001) != 0);
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     * @return The resultCacheResponse.
     */
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse getResultCacheResponse() {
      if (resultCacheResponseBuilder_ == null) {
        return resultCacheResponse_ == null ? com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.getDefaultInstance() : resultCacheResponse_;
      } else {
        return resultCacheResponseBuilder_.getMessage();
      }
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     */
    public Builder setResultCacheResponse(com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse value) {
      if (resultCacheResponseBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        resultCacheResponse_ = value;
      } else {
        resultCacheResponseBuilder_.setMessage(value);
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     */
    public Builder setResultCacheResponse(
        com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.Builder builderForValue) {
      if (resultCacheResponseBuilder_ == null) {
        resultCacheResponse_ = builderForValue.build();
      } else {
        resultCacheResponseBuilder_.setMessage(builderForValue.build());
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     */
    public Builder mergeResultCacheResponse(com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse value) {
      if (resultCacheResponseBuilder_ == null) {
        if (((bitField0_ & 0x00000001) != 0) &&
          resultCacheResponse_ != null &&
          resultCacheResponse_ != com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.getDefaultInstance()) {
          getResultCacheResponseBuilder().mergeFrom(value);
        } else {
          resultCacheResponse_ = value;
        }
      } else {
        resultCacheResponseBuilder_.mergeFrom(value);
      }
      if (resultCacheResponse_ != null) {
        bitField0_ |= 0x00000001;
        onChanged();
      }
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     */
    public Builder clearResultCacheResponse() {
      bitField0_ = (bitField0_ & ~0x00000001);
      resultCacheResponse_ = null;
      if (resultCacheResponseBuilder_ != null) {
        resultCacheResponseBuilder_.dispose();
        resultCacheResponseBuilder_ = null;
      }
      onChanged();
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     */
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.Builder getResultCacheResponseBuilder() {
      bitField0_ |= 0x00000001;
      onChanged();
      return getResultCacheResponseFieldBuilder().getBuilder();
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     */
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponseOrBuilder getResultCacheResponseOrBuilder() {
      if (resultCacheResponseBuilder_ != null) {
        return resultCacheResponseBuilder_.getMessageOrBuilder();
      } else {
        return resultCacheResponse_ == null ?
            com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.getDefaultInstance() : resultCacheResponse_;
      }
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.KeyDeleteResponse result_cache_response = 1 [json_name = "resultCacheResponse"];</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.Builder, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponseOrBuilder> 
        getResultCacheResponseFieldBuilder() {
      if (resultCacheResponseBuilder_ == null) {
        resultCacheResponseBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse.Builder, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponseOrBuilder>(
                getResultCacheResponse(),
                getParentForChildren(),
                isClean());
        resultCacheResponse_ = null;
      }
      return resultCacheResponseBuilder_;
    }

    private com.google.protobuf.MapField<
        java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> globalCacheResponses_;
    private com.google.protobuf.MapField<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse>
        internalGetGlobalCacheResponses() {
      if (globalCacheResponses_ == null) {
        return com.google.protobuf.MapField.emptyMapField(
            GlobalCacheResponsesDefaultEntryHolder.defaultEntry);
      }
      return globalCacheResponses_;
    }
    private com.google.protobuf.MapField<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse>
        internalGetMutableGlobalCacheResponses() {
      if (globalCacheResponses_ == null) {
        globalCacheResponses_ = com.google.protobuf.MapField.newMapField(
            GlobalCacheResponsesDefaultEntryHolder.defaultEntry);
      }
      if (!globalCacheResponses_.isMutable()) {
        globalCacheResponses_ = globalCacheResponses_.copy();
      }
      bitField0_ |= 0x00000002;
      onChanged();
      return globalCacheResponses_;
    }
    public int getGlobalCacheResponsesCount() {
      return internalGetGlobalCacheResponses().getMap().size();
    }
    /**
     * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
     */
    @java.lang.Override
    public boolean containsGlobalCacheResponses(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      return internalGetGlobalCacheResponses().getMap().containsKey(key);
    }
    /**
     * Use {@link #getGlobalCacheResponsesMap()} instead.
     */
    @java.lang.Override
    @java.lang.Deprecated
    public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> getGlobalCacheResponses() {
      return getGlobalCacheResponsesMap();
    }
    /**
     * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
     */
    @java.lang.Override
    public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> getGlobalCacheResponsesMap() {
      return internalGetGlobalCacheResponses().getMap();
    }
    /**
     * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
     */
    @java.lang.Override
    public /* nullable */
com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse getGlobalCacheResponsesOrDefault(
        java.lang.String key,
        /* nullable */
com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse defaultValue) {
      if (key == null) { throw new NullPointerException("map key"); }
      java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> map =
          internalGetGlobalCacheResponses().getMap();
      return map.containsKey(key) ? map.get(key) : defaultValue;
    }
    /**
     * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
     */
    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse getGlobalCacheResponsesOrThrow(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> map =
          internalGetGlobalCacheResponses().getMap();
      if (!map.containsKey(key)) {
        throw new java.lang.IllegalArgumentException();
      }
      return map.get(key);
    }
    public Builder clearGlobalCacheResponses() {
      bitField0_ = (bitField0_ & ~0x00000002);
      internalGetMutableGlobalCacheResponses().getMutableMap()
          .clear();
      return this;
    }
    /**
     * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
     */
    public Builder removeGlobalCacheResponses(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      internalGetMutableGlobalCacheResponses().getMutableMap()
          .remove(key);
      return this;
    }
    /**
     * Use alternate mutation accessors instead.
     */
    @java.lang.Deprecated
    public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse>
        getMutableGlobalCacheResponses() {
      bitField0_ |= 0x00000002;
      return internalGetMutableGlobalCacheResponses().getMutableMap();
    }
    /**
     * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
     */
    public Builder putGlobalCacheResponses(
        java.lang.String key,
        com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse value) {
      if (key == null) { throw new NullPointerException("map key"); }
      if (value == null) { throw new NullPointerException("map value"); }
      internalGetMutableGlobalCacheResponses().getMutableMap()
          .put(key, value);
      bitField0_ |= 0x00000002;
      return this;
    }
    /**
     * <code>map&lt;string, .aperture.flowcontrol.check.v1.KeyDeleteResponse&gt; global_cache_responses = 2 [json_name = "globalCacheResponses"];</code>
     */
    public Builder putAllGlobalCacheResponses(
        java.util.Map<java.lang.String, com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyDeleteResponse> values) {
      internalGetMutableGlobalCacheResponses().getMutableMap()
          .putAll(values);
      bitField0_ |= 0x00000002;
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


    // @@protoc_insertion_point(builder_scope:aperture.flowcontrol.check.v1.CacheDeleteResponse)
  }

  // @@protoc_insertion_point(class_scope:aperture.flowcontrol.check.v1.CacheDeleteResponse)
  private static final com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse();
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<CacheDeleteResponse>
      PARSER = new com.google.protobuf.AbstractParser<CacheDeleteResponse>() {
    @java.lang.Override
    public CacheDeleteResponse parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      Builder builder = newBuilder();
      try {
        builder.mergeFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.setUnfinishedMessage(builder.buildPartial());
      } catch (com.google.protobuf.UninitializedMessageException e) {
        throw e.asInvalidProtocolBufferException().setUnfinishedMessage(builder.buildPartial());
      } catch (java.io.IOException e) {
        throw new com.google.protobuf.InvalidProtocolBufferException(e)
            .setUnfinishedMessage(builder.buildPartial());
      }
      return builder.buildPartial();
    }
  };

  public static com.google.protobuf.Parser<CacheDeleteResponse> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<CacheDeleteResponse> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

