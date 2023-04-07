// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/cmd/v1/cmd.proto

package com.fluxninja.generated.aperture.cmd.v1;

/**
 * Protobuf type {@code aperture.cmd.v1.ListServicesControllerResponse}
 */
public final class ListServicesControllerResponse extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.cmd.v1.ListServicesControllerResponse)
    ListServicesControllerResponseOrBuilder {
private static final long serialVersionUID = 0L;
  // Use ListServicesControllerResponse.newBuilder() to construct.
  private ListServicesControllerResponse(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private ListServicesControllerResponse() {
    services_ = java.util.Collections.emptyList();
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new ListServicesControllerResponse();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private ListServicesControllerResponse(
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
          case 10: {
            if (!((mutable_bitField0_ & 0x00000001) != 0)) {
              services_ = new java.util.ArrayList<com.fluxninja.generated.aperture.cmd.v1.GlobalService>();
              mutable_bitField0_ |= 0x00000001;
            }
            services_.add(
                input.readMessage(com.fluxninja.generated.aperture.cmd.v1.GlobalService.parser(), extensionRegistry));
            break;
          }
          case 16: {

            errorsCount_ = input.readUInt32();
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
      if (((mutable_bitField0_ & 0x00000001) != 0)) {
        services_ = java.util.Collections.unmodifiableList(services_);
      }
      this.unknownFields = unknownFields.build();
      makeExtensionsImmutable();
    }
  }
  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.fluxninja.generated.aperture.cmd.v1.CmdProto.internal_static_aperture_cmd_v1_ListServicesControllerResponse_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.cmd.v1.CmdProto.internal_static_aperture_cmd_v1_ListServicesControllerResponse_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.class, com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.Builder.class);
  }

  public static final int SERVICES_FIELD_NUMBER = 1;
  private java.util.List<com.fluxninja.generated.aperture.cmd.v1.GlobalService> services_;
  /**
   * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
   */
  @java.lang.Override
  public java.util.List<com.fluxninja.generated.aperture.cmd.v1.GlobalService> getServicesList() {
    return services_;
  }
  /**
   * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
   */
  @java.lang.Override
  public java.util.List<? extends com.fluxninja.generated.aperture.cmd.v1.GlobalServiceOrBuilder> 
      getServicesOrBuilderList() {
    return services_;
  }
  /**
   * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
   */
  @java.lang.Override
  public int getServicesCount() {
    return services_.size();
  }
  /**
   * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.cmd.v1.GlobalService getServices(int index) {
    return services_.get(index);
  }
  /**
   * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.cmd.v1.GlobalServiceOrBuilder getServicesOrBuilder(
      int index) {
    return services_.get(index);
  }

  public static final int ERRORS_COUNT_FIELD_NUMBER = 2;
  private int errorsCount_;
  /**
   * <code>uint32 errors_count = 2 [json_name = "errorsCount"];</code>
   * @return The errorsCount.
   */
  @java.lang.Override
  public int getErrorsCount() {
    return errorsCount_;
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
    for (int i = 0; i < services_.size(); i++) {
      output.writeMessage(1, services_.get(i));
    }
    if (errorsCount_ != 0) {
      output.writeUInt32(2, errorsCount_);
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    for (int i = 0; i < services_.size(); i++) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(1, services_.get(i));
    }
    if (errorsCount_ != 0) {
      size += com.google.protobuf.CodedOutputStream
        .computeUInt32Size(2, errorsCount_);
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
    if (!(obj instanceof com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse other = (com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse) obj;

    if (!getServicesList()
        .equals(other.getServicesList())) return false;
    if (getErrorsCount()
        != other.getErrorsCount()) return false;
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
    if (getServicesCount() > 0) {
      hash = (37 * hash) + SERVICES_FIELD_NUMBER;
      hash = (53 * hash) + getServicesList().hashCode();
    }
    hash = (37 * hash) + ERRORS_COUNT_FIELD_NUMBER;
    hash = (53 * hash) + getErrorsCount();
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse prototype) {
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
   * Protobuf type {@code aperture.cmd.v1.ListServicesControllerResponse}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.cmd.v1.ListServicesControllerResponse)
      com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponseOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.cmd.v1.CmdProto.internal_static_aperture_cmd_v1_ListServicesControllerResponse_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.cmd.v1.CmdProto.internal_static_aperture_cmd_v1_ListServicesControllerResponse_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.class, com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.newBuilder()
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
        getServicesFieldBuilder();
      }
    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      if (servicesBuilder_ == null) {
        services_ = java.util.Collections.emptyList();
        bitField0_ = (bitField0_ & ~0x00000001);
      } else {
        servicesBuilder_.clear();
      }
      errorsCount_ = 0;

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.cmd.v1.CmdProto.internal_static_aperture_cmd_v1_ListServicesControllerResponse_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse build() {
      com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse buildPartial() {
      com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse result = new com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse(this);
      int from_bitField0_ = bitField0_;
      if (servicesBuilder_ == null) {
        if (((bitField0_ & 0x00000001) != 0)) {
          services_ = java.util.Collections.unmodifiableList(services_);
          bitField0_ = (bitField0_ & ~0x00000001);
        }
        result.services_ = services_;
      } else {
        result.services_ = servicesBuilder_.build();
      }
      result.errorsCount_ = errorsCount_;
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
      if (other instanceof com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse) {
        return mergeFrom((com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse other) {
      if (other == com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.getDefaultInstance()) return this;
      if (servicesBuilder_ == null) {
        if (!other.services_.isEmpty()) {
          if (services_.isEmpty()) {
            services_ = other.services_;
            bitField0_ = (bitField0_ & ~0x00000001);
          } else {
            ensureServicesIsMutable();
            services_.addAll(other.services_);
          }
          onChanged();
        }
      } else {
        if (!other.services_.isEmpty()) {
          if (servicesBuilder_.isEmpty()) {
            servicesBuilder_.dispose();
            servicesBuilder_ = null;
            services_ = other.services_;
            bitField0_ = (bitField0_ & ~0x00000001);
            servicesBuilder_ = 
              com.google.protobuf.GeneratedMessageV3.alwaysUseFieldBuilders ?
                 getServicesFieldBuilder() : null;
          } else {
            servicesBuilder_.addAllMessages(other.services_);
          }
        }
      }
      if (other.getErrorsCount() != 0) {
        setErrorsCount(other.getErrorsCount());
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
      com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }
    private int bitField0_;

    private java.util.List<com.fluxninja.generated.aperture.cmd.v1.GlobalService> services_ =
      java.util.Collections.emptyList();
    private void ensureServicesIsMutable() {
      if (!((bitField0_ & 0x00000001) != 0)) {
        services_ = new java.util.ArrayList<com.fluxninja.generated.aperture.cmd.v1.GlobalService>(services_);
        bitField0_ |= 0x00000001;
       }
    }

    private com.google.protobuf.RepeatedFieldBuilderV3<
        com.fluxninja.generated.aperture.cmd.v1.GlobalService, com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder, com.fluxninja.generated.aperture.cmd.v1.GlobalServiceOrBuilder> servicesBuilder_;

    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public java.util.List<com.fluxninja.generated.aperture.cmd.v1.GlobalService> getServicesList() {
      if (servicesBuilder_ == null) {
        return java.util.Collections.unmodifiableList(services_);
      } else {
        return servicesBuilder_.getMessageList();
      }
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public int getServicesCount() {
      if (servicesBuilder_ == null) {
        return services_.size();
      } else {
        return servicesBuilder_.getCount();
      }
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public com.fluxninja.generated.aperture.cmd.v1.GlobalService getServices(int index) {
      if (servicesBuilder_ == null) {
        return services_.get(index);
      } else {
        return servicesBuilder_.getMessage(index);
      }
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder setServices(
        int index, com.fluxninja.generated.aperture.cmd.v1.GlobalService value) {
      if (servicesBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        ensureServicesIsMutable();
        services_.set(index, value);
        onChanged();
      } else {
        servicesBuilder_.setMessage(index, value);
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder setServices(
        int index, com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder builderForValue) {
      if (servicesBuilder_ == null) {
        ensureServicesIsMutable();
        services_.set(index, builderForValue.build());
        onChanged();
      } else {
        servicesBuilder_.setMessage(index, builderForValue.build());
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder addServices(com.fluxninja.generated.aperture.cmd.v1.GlobalService value) {
      if (servicesBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        ensureServicesIsMutable();
        services_.add(value);
        onChanged();
      } else {
        servicesBuilder_.addMessage(value);
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder addServices(
        int index, com.fluxninja.generated.aperture.cmd.v1.GlobalService value) {
      if (servicesBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        ensureServicesIsMutable();
        services_.add(index, value);
        onChanged();
      } else {
        servicesBuilder_.addMessage(index, value);
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder addServices(
        com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder builderForValue) {
      if (servicesBuilder_ == null) {
        ensureServicesIsMutable();
        services_.add(builderForValue.build());
        onChanged();
      } else {
        servicesBuilder_.addMessage(builderForValue.build());
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder addServices(
        int index, com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder builderForValue) {
      if (servicesBuilder_ == null) {
        ensureServicesIsMutable();
        services_.add(index, builderForValue.build());
        onChanged();
      } else {
        servicesBuilder_.addMessage(index, builderForValue.build());
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder addAllServices(
        java.lang.Iterable<? extends com.fluxninja.generated.aperture.cmd.v1.GlobalService> values) {
      if (servicesBuilder_ == null) {
        ensureServicesIsMutable();
        com.google.protobuf.AbstractMessageLite.Builder.addAll(
            values, services_);
        onChanged();
      } else {
        servicesBuilder_.addAllMessages(values);
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder clearServices() {
      if (servicesBuilder_ == null) {
        services_ = java.util.Collections.emptyList();
        bitField0_ = (bitField0_ & ~0x00000001);
        onChanged();
      } else {
        servicesBuilder_.clear();
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public Builder removeServices(int index) {
      if (servicesBuilder_ == null) {
        ensureServicesIsMutable();
        services_.remove(index);
        onChanged();
      } else {
        servicesBuilder_.remove(index);
      }
      return this;
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder getServicesBuilder(
        int index) {
      return getServicesFieldBuilder().getBuilder(index);
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public com.fluxninja.generated.aperture.cmd.v1.GlobalServiceOrBuilder getServicesOrBuilder(
        int index) {
      if (servicesBuilder_ == null) {
        return services_.get(index);  } else {
        return servicesBuilder_.getMessageOrBuilder(index);
      }
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public java.util.List<? extends com.fluxninja.generated.aperture.cmd.v1.GlobalServiceOrBuilder> 
         getServicesOrBuilderList() {
      if (servicesBuilder_ != null) {
        return servicesBuilder_.getMessageOrBuilderList();
      } else {
        return java.util.Collections.unmodifiableList(services_);
      }
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder addServicesBuilder() {
      return getServicesFieldBuilder().addBuilder(
          com.fluxninja.generated.aperture.cmd.v1.GlobalService.getDefaultInstance());
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder addServicesBuilder(
        int index) {
      return getServicesFieldBuilder().addBuilder(
          index, com.fluxninja.generated.aperture.cmd.v1.GlobalService.getDefaultInstance());
    }
    /**
     * <code>repeated .aperture.cmd.v1.GlobalService services = 1 [json_name = "services"];</code>
     */
    public java.util.List<com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder> 
         getServicesBuilderList() {
      return getServicesFieldBuilder().getBuilderList();
    }
    private com.google.protobuf.RepeatedFieldBuilderV3<
        com.fluxninja.generated.aperture.cmd.v1.GlobalService, com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder, com.fluxninja.generated.aperture.cmd.v1.GlobalServiceOrBuilder> 
        getServicesFieldBuilder() {
      if (servicesBuilder_ == null) {
        servicesBuilder_ = new com.google.protobuf.RepeatedFieldBuilderV3<
            com.fluxninja.generated.aperture.cmd.v1.GlobalService, com.fluxninja.generated.aperture.cmd.v1.GlobalService.Builder, com.fluxninja.generated.aperture.cmd.v1.GlobalServiceOrBuilder>(
                services_,
                ((bitField0_ & 0x00000001) != 0),
                getParentForChildren(),
                isClean());
        services_ = null;
      }
      return servicesBuilder_;
    }

    private int errorsCount_ ;
    /**
     * <code>uint32 errors_count = 2 [json_name = "errorsCount"];</code>
     * @return The errorsCount.
     */
    @java.lang.Override
    public int getErrorsCount() {
      return errorsCount_;
    }
    /**
     * <code>uint32 errors_count = 2 [json_name = "errorsCount"];</code>
     * @param value The errorsCount to set.
     * @return This builder for chaining.
     */
    public Builder setErrorsCount(int value) {
      
      errorsCount_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>uint32 errors_count = 2 [json_name = "errorsCount"];</code>
     * @return This builder for chaining.
     */
    public Builder clearErrorsCount() {
      
      errorsCount_ = 0;
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


    // @@protoc_insertion_point(builder_scope:aperture.cmd.v1.ListServicesControllerResponse)
  }

  // @@protoc_insertion_point(class_scope:aperture.cmd.v1.ListServicesControllerResponse)
  private static final com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse();
  }

  public static com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<ListServicesControllerResponse>
      PARSER = new com.google.protobuf.AbstractParser<ListServicesControllerResponse>() {
    @java.lang.Override
    public ListServicesControllerResponse parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new ListServicesControllerResponse(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<ListServicesControllerResponse> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<ListServicesControllerResponse> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

