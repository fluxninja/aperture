// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

/**
 * Protobuf type {@code aperture.flowcontrol.check.v1.CacheResponse}
 */
public final class CacheResponse extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.flowcontrol.check.v1.CacheResponse)
    CacheResponseOrBuilder {
private static final long serialVersionUID = 0L;
  // Use CacheResponse.newBuilder() to construct.
  private CacheResponse(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private CacheResponse() {
    value_ = com.google.protobuf.ByteString.EMPTY;
    result_ = 0;
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new CacheResponse();
  }

  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheResponse_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheResponse_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse.class, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse.Builder.class);
  }

  public static final int VALUE_FIELD_NUMBER = 1;
  private com.google.protobuf.ByteString value_ = com.google.protobuf.ByteString.EMPTY;
  /**
   * <code>bytes value = 1 [json_name = "value"];</code>
   * @return The value.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString getValue() {
    return value_;
  }

  public static final int RESULT_FIELD_NUMBER = 2;
  private int result_ = 0;
  /**
   * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
   * @return The enum numeric value on the wire for result.
   */
  @java.lang.Override public int getResultValue() {
    return result_;
  }
  /**
   * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
   * @return The result.
   */
  @java.lang.Override public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult getResult() {
    com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult result = com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult.forNumber(result_);
    return result == null ? com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult.UNRECOGNIZED : result;
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
    if (!value_.isEmpty()) {
      output.writeBytes(1, value_);
    }
    if (result_ != com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult.Hit.getNumber()) {
      output.writeEnum(2, result_);
    }
    getUnknownFields().writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!value_.isEmpty()) {
      size += com.google.protobuf.CodedOutputStream
        .computeBytesSize(1, value_);
    }
    if (result_ != com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult.Hit.getNumber()) {
      size += com.google.protobuf.CodedOutputStream
        .computeEnumSize(2, result_);
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
    if (!(obj instanceof com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse other = (com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse) obj;

    if (!getValue()
        .equals(other.getValue())) return false;
    if (result_ != other.result_) return false;
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
    hash = (37 * hash) + VALUE_FIELD_NUMBER;
    hash = (53 * hash) + getValue().hashCode();
    hash = (37 * hash) + RESULT_FIELD_NUMBER;
    hash = (53 * hash) + result_;
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse prototype) {
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
   * Protobuf type {@code aperture.flowcontrol.check.v1.CacheResponse}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.flowcontrol.check.v1.CacheResponse)
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponseOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheResponse_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheResponse_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse.class, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse.newBuilder()
    private Builder() {

    }

    private Builder(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      super(parent);

    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      bitField0_ = 0;
      value_ = com.google.protobuf.ByteString.EMPTY;
      result_ = 0;
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_CacheResponse_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse build() {
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse buildPartial() {
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse result = new com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse(this);
      if (bitField0_ != 0) { buildPartial0(result); }
      onBuilt();
      return result;
    }

    private void buildPartial0(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse result) {
      int from_bitField0_ = bitField0_;
      if (((from_bitField0_ & 0x00000001) != 0)) {
        result.value_ = value_;
      }
      if (((from_bitField0_ & 0x00000002) != 0)) {
        result.result_ = result_;
      }
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
      if (other instanceof com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse) {
        return mergeFrom((com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse other) {
      if (other == com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse.getDefaultInstance()) return this;
      if (other.getValue() != com.google.protobuf.ByteString.EMPTY) {
        setValue(other.getValue());
      }
      if (other.result_ != 0) {
        setResultValue(other.getResultValue());
      }
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
              value_ = input.readBytes();
              bitField0_ |= 0x00000001;
              break;
            } // case 10
            case 16: {
              result_ = input.readEnum();
              bitField0_ |= 0x00000002;
              break;
            } // case 16
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

    private com.google.protobuf.ByteString value_ = com.google.protobuf.ByteString.EMPTY;
    /**
     * <code>bytes value = 1 [json_name = "value"];</code>
     * @return The value.
     */
    @java.lang.Override
    public com.google.protobuf.ByteString getValue() {
      return value_;
    }
    /**
     * <code>bytes value = 1 [json_name = "value"];</code>
     * @param value The value to set.
     * @return This builder for chaining.
     */
    public Builder setValue(com.google.protobuf.ByteString value) {
      if (value == null) { throw new NullPointerException(); }
      value_ = value;
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>bytes value = 1 [json_name = "value"];</code>
     * @return This builder for chaining.
     */
    public Builder clearValue() {
      bitField0_ = (bitField0_ & ~0x00000001);
      value_ = getDefaultInstance().getValue();
      onChanged();
      return this;
    }

    private int result_ = 0;
    /**
     * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
     * @return The enum numeric value on the wire for result.
     */
    @java.lang.Override public int getResultValue() {
      return result_;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
     * @param value The enum numeric value on the wire for result to set.
     * @return This builder for chaining.
     */
    public Builder setResultValue(int value) {
      result_ = value;
      bitField0_ |= 0x00000002;
      onChanged();
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
     * @return The result.
     */
    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult getResult() {
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult result = com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult.forNumber(result_);
      return result == null ? com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult.UNRECOGNIZED : result;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
     * @param value The result to set.
     * @return This builder for chaining.
     */
    public Builder setResult(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResult value) {
      if (value == null) {
        throw new NullPointerException();
      }
      bitField0_ |= 0x00000002;
      result_ = value.getNumber();
      onChanged();
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.check.v1.CacheResult result = 2 [json_name = "result"];</code>
     * @return This builder for chaining.
     */
    public Builder clearResult() {
      bitField0_ = (bitField0_ & ~0x00000002);
      result_ = 0;
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


    // @@protoc_insertion_point(builder_scope:aperture.flowcontrol.check.v1.CacheResponse)
  }

  // @@protoc_insertion_point(class_scope:aperture.flowcontrol.check.v1.CacheResponse)
  private static final com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse();
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<CacheResponse>
      PARSER = new com.google.protobuf.AbstractParser<CacheResponse>() {
    @java.lang.Override
    public CacheResponse parsePartialFrom(
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

  public static com.google.protobuf.Parser<CacheResponse> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<CacheResponse> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheResponse getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

