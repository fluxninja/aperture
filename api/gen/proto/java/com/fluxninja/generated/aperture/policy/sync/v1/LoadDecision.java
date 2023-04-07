// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/sync/v1/concurrency_limiter.proto

package com.fluxninja.generated.aperture.policy.sync.v1;

/**
 * Protobuf type {@code aperture.policy.sync.v1.LoadDecision}
 */
public final class LoadDecision extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.policy.sync.v1.LoadDecision)
    LoadDecisionOrBuilder {
private static final long serialVersionUID = 0L;
  // Use LoadDecision.newBuilder() to construct.
  private LoadDecision(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private LoadDecision() {
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new LoadDecision();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private LoadDecision(
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
          case 9: {

            loadMultiplier_ = input.readDouble();
            break;
          }
          case 16: {

            passThrough_ = input.readBool();
            break;
          }
          case 26: {
            com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.Builder subBuilder = null;
            if (tickInfo_ != null) {
              subBuilder = tickInfo_.toBuilder();
            }
            tickInfo_ = input.readMessage(com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.parser(), extensionRegistry);
            if (subBuilder != null) {
              subBuilder.mergeFrom(tickInfo_);
              tickInfo_ = subBuilder.buildPartial();
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
    return com.fluxninja.generated.aperture.policy.sync.v1.ConcurrencyLimiterProto.internal_static_aperture_policy_sync_v1_LoadDecision_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.policy.sync.v1.ConcurrencyLimiterProto.internal_static_aperture_policy_sync_v1_LoadDecision_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision.class, com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision.Builder.class);
  }

  public static final int LOAD_MULTIPLIER_FIELD_NUMBER = 1;
  private double loadMultiplier_;
  /**
   * <code>double load_multiplier = 1 [json_name = "loadMultiplier"];</code>
   * @return The loadMultiplier.
   */
  @java.lang.Override
  public double getLoadMultiplier() {
    return loadMultiplier_;
  }

  public static final int PASS_THROUGH_FIELD_NUMBER = 2;
  private boolean passThrough_;
  /**
   * <code>bool pass_through = 2 [json_name = "passThrough"];</code>
   * @return The passThrough.
   */
  @java.lang.Override
  public boolean getPassThrough() {
    return passThrough_;
  }

  public static final int TICK_INFO_FIELD_NUMBER = 3;
  private com.fluxninja.generated.aperture.policy.sync.v1.TickInfo tickInfo_;
  /**
   * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
   * @return Whether the tickInfo field is set.
   */
  @java.lang.Override
  public boolean hasTickInfo() {
    return tickInfo_ != null;
  }
  /**
   * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
   * @return The tickInfo.
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.sync.v1.TickInfo getTickInfo() {
    return tickInfo_ == null ? com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.getDefaultInstance() : tickInfo_;
  }
  /**
   * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.sync.v1.TickInfoOrBuilder getTickInfoOrBuilder() {
    return getTickInfo();
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
    if (java.lang.Double.doubleToRawLongBits(loadMultiplier_) != 0) {
      output.writeDouble(1, loadMultiplier_);
    }
    if (passThrough_ != false) {
      output.writeBool(2, passThrough_);
    }
    if (tickInfo_ != null) {
      output.writeMessage(3, getTickInfo());
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (java.lang.Double.doubleToRawLongBits(loadMultiplier_) != 0) {
      size += com.google.protobuf.CodedOutputStream
        .computeDoubleSize(1, loadMultiplier_);
    }
    if (passThrough_ != false) {
      size += com.google.protobuf.CodedOutputStream
        .computeBoolSize(2, passThrough_);
    }
    if (tickInfo_ != null) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(3, getTickInfo());
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
    if (!(obj instanceof com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision other = (com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision) obj;

    if (java.lang.Double.doubleToLongBits(getLoadMultiplier())
        != java.lang.Double.doubleToLongBits(
            other.getLoadMultiplier())) return false;
    if (getPassThrough()
        != other.getPassThrough()) return false;
    if (hasTickInfo() != other.hasTickInfo()) return false;
    if (hasTickInfo()) {
      if (!getTickInfo()
          .equals(other.getTickInfo())) return false;
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
    hash = (37 * hash) + LOAD_MULTIPLIER_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        java.lang.Double.doubleToLongBits(getLoadMultiplier()));
    hash = (37 * hash) + PASS_THROUGH_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashBoolean(
        getPassThrough());
    if (hasTickInfo()) {
      hash = (37 * hash) + TICK_INFO_FIELD_NUMBER;
      hash = (53 * hash) + getTickInfo().hashCode();
    }
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision prototype) {
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
   * Protobuf type {@code aperture.policy.sync.v1.LoadDecision}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.policy.sync.v1.LoadDecision)
      com.fluxninja.generated.aperture.policy.sync.v1.LoadDecisionOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.policy.sync.v1.ConcurrencyLimiterProto.internal_static_aperture_policy_sync_v1_LoadDecision_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.policy.sync.v1.ConcurrencyLimiterProto.internal_static_aperture_policy_sync_v1_LoadDecision_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision.class, com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision.newBuilder()
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
      loadMultiplier_ = 0D;

      passThrough_ = false;

      if (tickInfoBuilder_ == null) {
        tickInfo_ = null;
      } else {
        tickInfo_ = null;
        tickInfoBuilder_ = null;
      }
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.policy.sync.v1.ConcurrencyLimiterProto.internal_static_aperture_policy_sync_v1_LoadDecision_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision build() {
      com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision buildPartial() {
      com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision result = new com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision(this);
      result.loadMultiplier_ = loadMultiplier_;
      result.passThrough_ = passThrough_;
      if (tickInfoBuilder_ == null) {
        result.tickInfo_ = tickInfo_;
      } else {
        result.tickInfo_ = tickInfoBuilder_.build();
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
      if (other instanceof com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision) {
        return mergeFrom((com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision other) {
      if (other == com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision.getDefaultInstance()) return this;
      if (other.getLoadMultiplier() != 0D) {
        setLoadMultiplier(other.getLoadMultiplier());
      }
      if (other.getPassThrough() != false) {
        setPassThrough(other.getPassThrough());
      }
      if (other.hasTickInfo()) {
        mergeTickInfo(other.getTickInfo());
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
      com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private double loadMultiplier_ ;
    /**
     * <code>double load_multiplier = 1 [json_name = "loadMultiplier"];</code>
     * @return The loadMultiplier.
     */
    @java.lang.Override
    public double getLoadMultiplier() {
      return loadMultiplier_;
    }
    /**
     * <code>double load_multiplier = 1 [json_name = "loadMultiplier"];</code>
     * @param value The loadMultiplier to set.
     * @return This builder for chaining.
     */
    public Builder setLoadMultiplier(double value) {
      
      loadMultiplier_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>double load_multiplier = 1 [json_name = "loadMultiplier"];</code>
     * @return This builder for chaining.
     */
    public Builder clearLoadMultiplier() {
      
      loadMultiplier_ = 0D;
      onChanged();
      return this;
    }

    private boolean passThrough_ ;
    /**
     * <code>bool pass_through = 2 [json_name = "passThrough"];</code>
     * @return The passThrough.
     */
    @java.lang.Override
    public boolean getPassThrough() {
      return passThrough_;
    }
    /**
     * <code>bool pass_through = 2 [json_name = "passThrough"];</code>
     * @param value The passThrough to set.
     * @return This builder for chaining.
     */
    public Builder setPassThrough(boolean value) {
      
      passThrough_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>bool pass_through = 2 [json_name = "passThrough"];</code>
     * @return This builder for chaining.
     */
    public Builder clearPassThrough() {
      
      passThrough_ = false;
      onChanged();
      return this;
    }

    private com.fluxninja.generated.aperture.policy.sync.v1.TickInfo tickInfo_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.policy.sync.v1.TickInfo, com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.Builder, com.fluxninja.generated.aperture.policy.sync.v1.TickInfoOrBuilder> tickInfoBuilder_;
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     * @return Whether the tickInfo field is set.
     */
    public boolean hasTickInfo() {
      return tickInfoBuilder_ != null || tickInfo_ != null;
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     * @return The tickInfo.
     */
    public com.fluxninja.generated.aperture.policy.sync.v1.TickInfo getTickInfo() {
      if (tickInfoBuilder_ == null) {
        return tickInfo_ == null ? com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.getDefaultInstance() : tickInfo_;
      } else {
        return tickInfoBuilder_.getMessage();
      }
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     */
    public Builder setTickInfo(com.fluxninja.generated.aperture.policy.sync.v1.TickInfo value) {
      if (tickInfoBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        tickInfo_ = value;
        onChanged();
      } else {
        tickInfoBuilder_.setMessage(value);
      }

      return this;
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     */
    public Builder setTickInfo(
        com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.Builder builderForValue) {
      if (tickInfoBuilder_ == null) {
        tickInfo_ = builderForValue.build();
        onChanged();
      } else {
        tickInfoBuilder_.setMessage(builderForValue.build());
      }

      return this;
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     */
    public Builder mergeTickInfo(com.fluxninja.generated.aperture.policy.sync.v1.TickInfo value) {
      if (tickInfoBuilder_ == null) {
        if (tickInfo_ != null) {
          tickInfo_ =
            com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.newBuilder(tickInfo_).mergeFrom(value).buildPartial();
        } else {
          tickInfo_ = value;
        }
        onChanged();
      } else {
        tickInfoBuilder_.mergeFrom(value);
      }

      return this;
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     */
    public Builder clearTickInfo() {
      if (tickInfoBuilder_ == null) {
        tickInfo_ = null;
        onChanged();
      } else {
        tickInfo_ = null;
        tickInfoBuilder_ = null;
      }

      return this;
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     */
    public com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.Builder getTickInfoBuilder() {
      
      onChanged();
      return getTickInfoFieldBuilder().getBuilder();
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     */
    public com.fluxninja.generated.aperture.policy.sync.v1.TickInfoOrBuilder getTickInfoOrBuilder() {
      if (tickInfoBuilder_ != null) {
        return tickInfoBuilder_.getMessageOrBuilder();
      } else {
        return tickInfo_ == null ?
            com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.getDefaultInstance() : tickInfo_;
      }
    }
    /**
     * <code>.aperture.policy.sync.v1.TickInfo tick_info = 3 [json_name = "tickInfo"];</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.policy.sync.v1.TickInfo, com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.Builder, com.fluxninja.generated.aperture.policy.sync.v1.TickInfoOrBuilder> 
        getTickInfoFieldBuilder() {
      if (tickInfoBuilder_ == null) {
        tickInfoBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.fluxninja.generated.aperture.policy.sync.v1.TickInfo, com.fluxninja.generated.aperture.policy.sync.v1.TickInfo.Builder, com.fluxninja.generated.aperture.policy.sync.v1.TickInfoOrBuilder>(
                getTickInfo(),
                getParentForChildren(),
                isClean());
        tickInfo_ = null;
      }
      return tickInfoBuilder_;
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


    // @@protoc_insertion_point(builder_scope:aperture.policy.sync.v1.LoadDecision)
  }

  // @@protoc_insertion_point(class_scope:aperture.policy.sync.v1.LoadDecision)
  private static final com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision();
  }

  public static com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<LoadDecision>
      PARSER = new com.google.protobuf.AbstractParser<LoadDecision>() {
    @java.lang.Override
    public LoadDecision parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new LoadDecision(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<LoadDecision> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<LoadDecision> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.sync.v1.LoadDecision getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

