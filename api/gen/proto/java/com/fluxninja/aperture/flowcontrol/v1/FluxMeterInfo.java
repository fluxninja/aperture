// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/v1/flowcontrol.proto

package com.fluxninja.aperture.flowcontrol.v1;

/**
 * <pre>
 * FluxMeterInfo describes detail for each FluxMeterInfo.
 * </pre>
 *
 * Protobuf type {@code aperture.flowcontrol.v1.FluxMeterInfo}
 */
public final class FluxMeterInfo extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.flowcontrol.v1.FluxMeterInfo)
    FluxMeterInfoOrBuilder {
private static final long serialVersionUID = 0L;
  // Use FluxMeterInfo.newBuilder() to construct.
  private FluxMeterInfo(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private FluxMeterInfo() {
    fluxMeterName_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new FluxMeterInfo();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_FluxMeterInfo_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_FluxMeterInfo_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo.class, com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo.Builder.class);
  }

  public static final int FLUX_METER_NAME_FIELD_NUMBER = 1;
  private volatile java.lang.Object fluxMeterName_;
  /**
   * <code>string flux_meter_name = 1 [json_name = "fluxMeterName"];</code>
   * @return The fluxMeterName.
   */
  @java.lang.Override
  public java.lang.String getFluxMeterName() {
    java.lang.Object ref = fluxMeterName_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      fluxMeterName_ = s;
      return s;
    }
  }
  /**
   * <code>string flux_meter_name = 1 [json_name = "fluxMeterName"];</code>
   * @return The bytes for fluxMeterName.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getFluxMeterNameBytes() {
    java.lang.Object ref = fluxMeterName_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      fluxMeterName_ = b;
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
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(fluxMeterName_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 1, fluxMeterName_);
    }
    getUnknownFields().writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(fluxMeterName_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, fluxMeterName_);
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
    if (!(obj instanceof com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo)) {
      return super.equals(obj);
    }
    com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo other = (com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo) obj;

    if (!getFluxMeterName()
        .equals(other.getFluxMeterName())) return false;
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
    hash = (37 * hash) + FLUX_METER_NAME_FIELD_NUMBER;
    hash = (53 * hash) + getFluxMeterName().hashCode();
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo parseFrom(
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
  public static Builder newBuilder(com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo prototype) {
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
   * FluxMeterInfo describes detail for each FluxMeterInfo.
   * </pre>
   *
   * Protobuf type {@code aperture.flowcontrol.v1.FluxMeterInfo}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.flowcontrol.v1.FluxMeterInfo)
      com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfoOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_FluxMeterInfo_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_FluxMeterInfo_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo.class, com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo.Builder.class);
    }

    // Construct using com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo.newBuilder()
    private Builder() {

    }

    private Builder(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      super(parent);

    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      fluxMeterName_ = "";

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_FluxMeterInfo_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo getDefaultInstanceForType() {
      return com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo build() {
      com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo buildPartial() {
      com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo result = new com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo(this);
      result.fluxMeterName_ = fluxMeterName_;
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
      if (other instanceof com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo) {
        return mergeFrom((com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo other) {
      if (other == com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo.getDefaultInstance()) return this;
      if (!other.getFluxMeterName().isEmpty()) {
        fluxMeterName_ = other.fluxMeterName_;
        onChanged();
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
              fluxMeterName_ = input.readStringRequireUtf8();

              break;
            } // case 10
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

    private java.lang.Object fluxMeterName_ = "";
    /**
     * <code>string flux_meter_name = 1 [json_name = "fluxMeterName"];</code>
     * @return The fluxMeterName.
     */
    public java.lang.String getFluxMeterName() {
      java.lang.Object ref = fluxMeterName_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        fluxMeterName_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string flux_meter_name = 1 [json_name = "fluxMeterName"];</code>
     * @return The bytes for fluxMeterName.
     */
    public com.google.protobuf.ByteString
        getFluxMeterNameBytes() {
      java.lang.Object ref = fluxMeterName_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        fluxMeterName_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string flux_meter_name = 1 [json_name = "fluxMeterName"];</code>
     * @param value The fluxMeterName to set.
     * @return This builder for chaining.
     */
    public Builder setFluxMeterName(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      fluxMeterName_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>string flux_meter_name = 1 [json_name = "fluxMeterName"];</code>
     * @return This builder for chaining.
     */
    public Builder clearFluxMeterName() {
      
      fluxMeterName_ = getDefaultInstance().getFluxMeterName();
      onChanged();
      return this;
    }
    /**
     * <code>string flux_meter_name = 1 [json_name = "fluxMeterName"];</code>
     * @param value The bytes for fluxMeterName to set.
     * @return This builder for chaining.
     */
    public Builder setFluxMeterNameBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      fluxMeterName_ = value;
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


    // @@protoc_insertion_point(builder_scope:aperture.flowcontrol.v1.FluxMeterInfo)
  }

  // @@protoc_insertion_point(class_scope:aperture.flowcontrol.v1.FluxMeterInfo)
  private static final com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo();
  }

  public static com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<FluxMeterInfo>
      PARSER = new com.google.protobuf.AbstractParser<FluxMeterInfo>() {
    @java.lang.Override
    public FluxMeterInfo parsePartialFrom(
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

  public static com.google.protobuf.Parser<FluxMeterInfo> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<FluxMeterInfo> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.aperture.flowcontrol.v1.FluxMeterInfo getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

