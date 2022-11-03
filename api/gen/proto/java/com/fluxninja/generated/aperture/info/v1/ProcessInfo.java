// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/info/v1/info.proto

package com.fluxninja.generated.aperture.info.v1;

/**
 * Protobuf type {@code aperture.info.v1.ProcessInfo}
 */
public final class ProcessInfo extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.info.v1.ProcessInfo)
    ProcessInfoOrBuilder {
private static final long serialVersionUID = 0L;
  // Use ProcessInfo.newBuilder() to construct.
  private ProcessInfo(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private ProcessInfo() {
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new ProcessInfo();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private ProcessInfo(
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
            com.google.protobuf.Timestamp.Builder subBuilder = null;
            if (startTime_ != null) {
              subBuilder = startTime_.toBuilder();
            }
            startTime_ = input.readMessage(com.google.protobuf.Timestamp.parser(), extensionRegistry);
            if (subBuilder != null) {
              subBuilder.mergeFrom(startTime_);
              startTime_ = subBuilder.buildPartial();
            }

            break;
          }
          case 18: {
            com.google.protobuf.Duration.Builder subBuilder = null;
            if (uptime_ != null) {
              subBuilder = uptime_.toBuilder();
            }
            uptime_ = input.readMessage(com.google.protobuf.Duration.parser(), extensionRegistry);
            if (subBuilder != null) {
              subBuilder.mergeFrom(uptime_);
              uptime_ = subBuilder.buildPartial();
            }

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
    return com.fluxninja.generated.aperture.info.v1.InfoProto.internal_static_aperture_info_v1_ProcessInfo_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.info.v1.InfoProto.internal_static_aperture_info_v1_ProcessInfo_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.info.v1.ProcessInfo.class, com.fluxninja.generated.aperture.info.v1.ProcessInfo.Builder.class);
  }

  public static final int START_TIME_FIELD_NUMBER = 1;
  private com.google.protobuf.Timestamp startTime_;
  /**
   * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
   * @return Whether the startTime field is set.
   */
  @java.lang.Override
  public boolean hasStartTime() {
    return startTime_ != null;
  }
  /**
   * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
   * @return The startTime.
   */
  @java.lang.Override
  public com.google.protobuf.Timestamp getStartTime() {
    return startTime_ == null ? com.google.protobuf.Timestamp.getDefaultInstance() : startTime_;
  }
  /**
   * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
   */
  @java.lang.Override
  public com.google.protobuf.TimestampOrBuilder getStartTimeOrBuilder() {
    return getStartTime();
  }

  public static final int UPTIME_FIELD_NUMBER = 2;
  private com.google.protobuf.Duration uptime_;
  /**
   * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
   * @return Whether the uptime field is set.
   */
  @java.lang.Override
  public boolean hasUptime() {
    return uptime_ != null;
  }
  /**
   * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
   * @return The uptime.
   */
  @java.lang.Override
  public com.google.protobuf.Duration getUptime() {
    return uptime_ == null ? com.google.protobuf.Duration.getDefaultInstance() : uptime_;
  }
  /**
   * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
   */
  @java.lang.Override
  public com.google.protobuf.DurationOrBuilder getUptimeOrBuilder() {
    return getUptime();
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
    if (startTime_ != null) {
      output.writeMessage(1, getStartTime());
    }
    if (uptime_ != null) {
      output.writeMessage(2, getUptime());
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (startTime_ != null) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(1, getStartTime());
    }
    if (uptime_ != null) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(2, getUptime());
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
    if (!(obj instanceof com.fluxninja.generated.aperture.info.v1.ProcessInfo)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.info.v1.ProcessInfo other = (com.fluxninja.generated.aperture.info.v1.ProcessInfo) obj;

    if (hasStartTime() != other.hasStartTime()) return false;
    if (hasStartTime()) {
      if (!getStartTime()
          .equals(other.getStartTime())) return false;
    }
    if (hasUptime() != other.hasUptime()) return false;
    if (hasUptime()) {
      if (!getUptime()
          .equals(other.getUptime())) return false;
    }
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
    if (hasStartTime()) {
      hash = (37 * hash) + START_TIME_FIELD_NUMBER;
      hash = (53 * hash) + getStartTime().hashCode();
    }
    if (hasUptime()) {
      hash = (37 * hash) + UPTIME_FIELD_NUMBER;
      hash = (53 * hash) + getUptime().hashCode();
    }
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.info.v1.ProcessInfo prototype) {
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
   * Protobuf type {@code aperture.info.v1.ProcessInfo}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.info.v1.ProcessInfo)
      com.fluxninja.generated.aperture.info.v1.ProcessInfoOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.info.v1.InfoProto.internal_static_aperture_info_v1_ProcessInfo_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.info.v1.InfoProto.internal_static_aperture_info_v1_ProcessInfo_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.info.v1.ProcessInfo.class, com.fluxninja.generated.aperture.info.v1.ProcessInfo.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.info.v1.ProcessInfo.newBuilder()
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
      if (startTimeBuilder_ == null) {
        startTime_ = null;
      } else {
        startTime_ = null;
        startTimeBuilder_ = null;
      }
      if (uptimeBuilder_ == null) {
        uptime_ = null;
      } else {
        uptime_ = null;
        uptimeBuilder_ = null;
      }
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.info.v1.InfoProto.internal_static_aperture_info_v1_ProcessInfo_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.info.v1.ProcessInfo getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.info.v1.ProcessInfo.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.info.v1.ProcessInfo build() {
      com.fluxninja.generated.aperture.info.v1.ProcessInfo result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.info.v1.ProcessInfo buildPartial() {
      com.fluxninja.generated.aperture.info.v1.ProcessInfo result = new com.fluxninja.generated.aperture.info.v1.ProcessInfo(this);
      if (startTimeBuilder_ == null) {
        result.startTime_ = startTime_;
      } else {
        result.startTime_ = startTimeBuilder_.build();
      }
      if (uptimeBuilder_ == null) {
        result.uptime_ = uptime_;
      } else {
        result.uptime_ = uptimeBuilder_.build();
      }
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
      if (other instanceof com.fluxninja.generated.aperture.info.v1.ProcessInfo) {
        return mergeFrom((com.fluxninja.generated.aperture.info.v1.ProcessInfo)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.info.v1.ProcessInfo other) {
      if (other == com.fluxninja.generated.aperture.info.v1.ProcessInfo.getDefaultInstance()) return this;
      if (other.hasStartTime()) {
        mergeStartTime(other.getStartTime());
      }
      if (other.hasUptime()) {
        mergeUptime(other.getUptime());
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
      com.fluxninja.generated.aperture.info.v1.ProcessInfo parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.info.v1.ProcessInfo) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private com.google.protobuf.Timestamp startTime_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.google.protobuf.Timestamp, com.google.protobuf.Timestamp.Builder, com.google.protobuf.TimestampOrBuilder> startTimeBuilder_;
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     * @return Whether the startTime field is set.
     */
    public boolean hasStartTime() {
      return startTimeBuilder_ != null || startTime_ != null;
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     * @return The startTime.
     */
    public com.google.protobuf.Timestamp getStartTime() {
      if (startTimeBuilder_ == null) {
        return startTime_ == null ? com.google.protobuf.Timestamp.getDefaultInstance() : startTime_;
      } else {
        return startTimeBuilder_.getMessage();
      }
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     */
    public Builder setStartTime(com.google.protobuf.Timestamp value) {
      if (startTimeBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        startTime_ = value;
        onChanged();
      } else {
        startTimeBuilder_.setMessage(value);
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     */
    public Builder setStartTime(
        com.google.protobuf.Timestamp.Builder builderForValue) {
      if (startTimeBuilder_ == null) {
        startTime_ = builderForValue.build();
        onChanged();
      } else {
        startTimeBuilder_.setMessage(builderForValue.build());
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     */
    public Builder mergeStartTime(com.google.protobuf.Timestamp value) {
      if (startTimeBuilder_ == null) {
        if (startTime_ != null) {
          startTime_ =
            com.google.protobuf.Timestamp.newBuilder(startTime_).mergeFrom(value).buildPartial();
        } else {
          startTime_ = value;
        }
        onChanged();
      } else {
        startTimeBuilder_.mergeFrom(value);
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     */
    public Builder clearStartTime() {
      if (startTimeBuilder_ == null) {
        startTime_ = null;
        onChanged();
      } else {
        startTime_ = null;
        startTimeBuilder_ = null;
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     */
    public com.google.protobuf.Timestamp.Builder getStartTimeBuilder() {
      
      onChanged();
      return getStartTimeFieldBuilder().getBuilder();
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     */
    public com.google.protobuf.TimestampOrBuilder getStartTimeOrBuilder() {
      if (startTimeBuilder_ != null) {
        return startTimeBuilder_.getMessageOrBuilder();
      } else {
        return startTime_ == null ?
            com.google.protobuf.Timestamp.getDefaultInstance() : startTime_;
      }
    }
    /**
     * <code>.google.protobuf.Timestamp start_time = 1 [json_name = "startTime"];</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.google.protobuf.Timestamp, com.google.protobuf.Timestamp.Builder, com.google.protobuf.TimestampOrBuilder> 
        getStartTimeFieldBuilder() {
      if (startTimeBuilder_ == null) {
        startTimeBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.google.protobuf.Timestamp, com.google.protobuf.Timestamp.Builder, com.google.protobuf.TimestampOrBuilder>(
                getStartTime(),
                getParentForChildren(),
                isClean());
        startTime_ = null;
      }
      return startTimeBuilder_;
    }

    private com.google.protobuf.Duration uptime_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.google.protobuf.Duration, com.google.protobuf.Duration.Builder, com.google.protobuf.DurationOrBuilder> uptimeBuilder_;
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     * @return Whether the uptime field is set.
     */
    public boolean hasUptime() {
      return uptimeBuilder_ != null || uptime_ != null;
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     * @return The uptime.
     */
    public com.google.protobuf.Duration getUptime() {
      if (uptimeBuilder_ == null) {
        return uptime_ == null ? com.google.protobuf.Duration.getDefaultInstance() : uptime_;
      } else {
        return uptimeBuilder_.getMessage();
      }
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     */
    public Builder setUptime(com.google.protobuf.Duration value) {
      if (uptimeBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        uptime_ = value;
        onChanged();
      } else {
        uptimeBuilder_.setMessage(value);
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     */
    public Builder setUptime(
        com.google.protobuf.Duration.Builder builderForValue) {
      if (uptimeBuilder_ == null) {
        uptime_ = builderForValue.build();
        onChanged();
      } else {
        uptimeBuilder_.setMessage(builderForValue.build());
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     */
    public Builder mergeUptime(com.google.protobuf.Duration value) {
      if (uptimeBuilder_ == null) {
        if (uptime_ != null) {
          uptime_ =
            com.google.protobuf.Duration.newBuilder(uptime_).mergeFrom(value).buildPartial();
        } else {
          uptime_ = value;
        }
        onChanged();
      } else {
        uptimeBuilder_.mergeFrom(value);
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     */
    public Builder clearUptime() {
      if (uptimeBuilder_ == null) {
        uptime_ = null;
        onChanged();
      } else {
        uptime_ = null;
        uptimeBuilder_ = null;
      }

      return this;
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     */
    public com.google.protobuf.Duration.Builder getUptimeBuilder() {
      
      onChanged();
      return getUptimeFieldBuilder().getBuilder();
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     */
    public com.google.protobuf.DurationOrBuilder getUptimeOrBuilder() {
      if (uptimeBuilder_ != null) {
        return uptimeBuilder_.getMessageOrBuilder();
      } else {
        return uptime_ == null ?
            com.google.protobuf.Duration.getDefaultInstance() : uptime_;
      }
    }
    /**
     * <code>.google.protobuf.Duration uptime = 2 [json_name = "uptime"];</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.google.protobuf.Duration, com.google.protobuf.Duration.Builder, com.google.protobuf.DurationOrBuilder> 
        getUptimeFieldBuilder() {
      if (uptimeBuilder_ == null) {
        uptimeBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.google.protobuf.Duration, com.google.protobuf.Duration.Builder, com.google.protobuf.DurationOrBuilder>(
                getUptime(),
                getParentForChildren(),
                isClean());
        uptime_ = null;
      }
      return uptimeBuilder_;
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


    // @@protoc_insertion_point(builder_scope:aperture.info.v1.ProcessInfo)
  }

  // @@protoc_insertion_point(class_scope:aperture.info.v1.ProcessInfo)
  private static final com.fluxninja.generated.aperture.info.v1.ProcessInfo DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.info.v1.ProcessInfo();
  }

  public static com.fluxninja.generated.aperture.info.v1.ProcessInfo getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<ProcessInfo>
      PARSER = new com.google.protobuf.AbstractParser<ProcessInfo>() {
    @java.lang.Override
    public ProcessInfo parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new ProcessInfo(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<ProcessInfo> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<ProcessInfo> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.info.v1.ProcessInfo getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

