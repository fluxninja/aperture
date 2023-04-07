// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/discovery/entities/v1/entities.proto

package com.fluxninja.generated.aperture.discovery.entities.v1;

/**
 * Protobuf type {@code aperture.discovery.entities.v1.GetEntityByIPAddressRequest}
 */
public final class GetEntityByIPAddressRequest extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.discovery.entities.v1.GetEntityByIPAddressRequest)
    GetEntityByIPAddressRequestOrBuilder {
private static final long serialVersionUID = 0L;
  // Use GetEntityByIPAddressRequest.newBuilder() to construct.
  private GetEntityByIPAddressRequest(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private GetEntityByIPAddressRequest() {
    ipAddress_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new GetEntityByIPAddressRequest();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private GetEntityByIPAddressRequest(
      com.google.protobuf.CodedInputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    this();
    if (extensionRegistry == null) {
      throw new java.lang.NullPointerException();
    }
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
          case 10: {
            java.lang.String s = input.readStringRequireUtf8();

            ipAddress_ = s;
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
    return com.fluxninja.generated.aperture.discovery.entities.v1.EntitiesProto.internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.discovery.entities.v1.EntitiesProto.internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.class, com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.Builder.class);
  }

  public static final int IP_ADDRESS_FIELD_NUMBER = 1;
  private volatile java.lang.Object ipAddress_;
  /**
   * <code>string ip_address = 1 [json_name = "ipAddress"];</code>
   * @return The ipAddress.
   */
  @java.lang.Override
  public java.lang.String getIpAddress() {
    java.lang.Object ref = ipAddress_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      ipAddress_ = s;
      return s;
    }
  }
  /**
   * <code>string ip_address = 1 [json_name = "ipAddress"];</code>
   * @return The bytes for ipAddress.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getIpAddressBytes() {
    java.lang.Object ref = ipAddress_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      ipAddress_ = b;
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
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(ipAddress_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 1, ipAddress_);
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(ipAddress_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, ipAddress_);
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
    if (!(obj instanceof com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest other = (com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest) obj;

    if (!getIpAddress()
        .equals(other.getIpAddress())) return false;
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
    hash = (37 * hash) + IP_ADDRESS_FIELD_NUMBER;
    hash = (53 * hash) + getIpAddress().hashCode();
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest prototype) {
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
   * Protobuf type {@code aperture.discovery.entities.v1.GetEntityByIPAddressRequest}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.discovery.entities.v1.GetEntityByIPAddressRequest)
      com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequestOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.discovery.entities.v1.EntitiesProto.internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.discovery.entities.v1.EntitiesProto.internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.class, com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.newBuilder()
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
      ipAddress_ = "";

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.discovery.entities.v1.EntitiesProto.internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest build() {
      com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest buildPartial() {
      com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest result = new com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest(this);
      result.ipAddress_ = ipAddress_;
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
      if (other instanceof com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest) {
        return mergeFrom((com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest other) {
      if (other == com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.getDefaultInstance()) return this;
      if (!other.getIpAddress().isEmpty()) {
        ipAddress_ = other.ipAddress_;
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
      com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private java.lang.Object ipAddress_ = "";
    /**
     * <code>string ip_address = 1 [json_name = "ipAddress"];</code>
     * @return The ipAddress.
     */
    public java.lang.String getIpAddress() {
      java.lang.Object ref = ipAddress_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        ipAddress_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string ip_address = 1 [json_name = "ipAddress"];</code>
     * @return The bytes for ipAddress.
     */
    public com.google.protobuf.ByteString
        getIpAddressBytes() {
      java.lang.Object ref = ipAddress_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        ipAddress_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string ip_address = 1 [json_name = "ipAddress"];</code>
     * @param value The ipAddress to set.
     * @return This builder for chaining.
     */
    public Builder setIpAddress(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      ipAddress_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>string ip_address = 1 [json_name = "ipAddress"];</code>
     * @return This builder for chaining.
     */
    public Builder clearIpAddress() {
      
      ipAddress_ = getDefaultInstance().getIpAddress();
      onChanged();
      return this;
    }
    /**
     * <code>string ip_address = 1 [json_name = "ipAddress"];</code>
     * @param value The bytes for ipAddress to set.
     * @return This builder for chaining.
     */
    public Builder setIpAddressBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      ipAddress_ = value;
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


    // @@protoc_insertion_point(builder_scope:aperture.discovery.entities.v1.GetEntityByIPAddressRequest)
  }

  // @@protoc_insertion_point(class_scope:aperture.discovery.entities.v1.GetEntityByIPAddressRequest)
  private static final com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest();
  }

  public static com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<GetEntityByIPAddressRequest>
      PARSER = new com.google.protobuf.AbstractParser<GetEntityByIPAddressRequest>() {
    @java.lang.Override
    public GetEntityByIPAddressRequest parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new GetEntityByIPAddressRequest(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<GetEntityByIPAddressRequest> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<GetEntityByIPAddressRequest> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

