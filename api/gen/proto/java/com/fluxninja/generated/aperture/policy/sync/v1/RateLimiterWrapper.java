// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/sync/v1/rate_limiter.proto

package com.fluxninja.generated.aperture.policy.sync.v1;

/**
 * Protobuf type {@code aperture.policy.sync.v1.RateLimiterWrapper}
 */
public final class RateLimiterWrapper extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.policy.sync.v1.RateLimiterWrapper)
    RateLimiterWrapperOrBuilder {
private static final long serialVersionUID = 0L;
  // Use RateLimiterWrapper.newBuilder() to construct.
  private RateLimiterWrapper(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private RateLimiterWrapper() {
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new RateLimiterWrapper();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private RateLimiterWrapper(
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
            com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.Builder subBuilder = null;
            if (commonAttributes_ != null) {
              subBuilder = commonAttributes_.toBuilder();
            }
            commonAttributes_ = input.readMessage(com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.parser(), extensionRegistry);
            if (subBuilder != null) {
              subBuilder.mergeFrom(commonAttributes_);
              commonAttributes_ = subBuilder.buildPartial();
            }

            break;
          }
          case 18: {
            com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.Builder subBuilder = null;
            if (rateLimiter_ != null) {
              subBuilder = rateLimiter_.toBuilder();
            }
            rateLimiter_ = input.readMessage(com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.parser(), extensionRegistry);
            if (subBuilder != null) {
              subBuilder.mergeFrom(rateLimiter_);
              rateLimiter_ = subBuilder.buildPartial();
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
    return com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterProto.internal_static_aperture_policy_sync_v1_RateLimiterWrapper_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterProto.internal_static_aperture_policy_sync_v1_RateLimiterWrapper_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper.class, com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper.Builder.class);
  }

  public static final int COMMON_ATTRIBUTES_FIELD_NUMBER = 1;
  private com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes commonAttributes_;
  /**
   * <pre>
   * CommonAttributes
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
   * @return Whether the commonAttributes field is set.
   */
  @java.lang.Override
  public boolean hasCommonAttributes() {
    return commonAttributes_ != null;
  }
  /**
   * <pre>
   * CommonAttributes
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
   * @return The commonAttributes.
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes getCommonAttributes() {
    return commonAttributes_ == null ? com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.getDefaultInstance() : commonAttributes_;
  }
  /**
   * <pre>
   * CommonAttributes
   * </pre>
   *
   * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesOrBuilder getCommonAttributesOrBuilder() {
    return getCommonAttributes();
  }

  public static final int RATE_LIMITER_FIELD_NUMBER = 2;
  private com.fluxninja.generated.aperture.policy.language.v1.RateLimiter rateLimiter_;
  /**
   * <pre>
   * Rate Limiter
   * </pre>
   *
   * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
   * @return Whether the rateLimiter field is set.
   */
  @java.lang.Override
  public boolean hasRateLimiter() {
    return rateLimiter_ != null;
  }
  /**
   * <pre>
   * Rate Limiter
   * </pre>
   *
   * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
   * @return The rateLimiter.
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.language.v1.RateLimiter getRateLimiter() {
    return rateLimiter_ == null ? com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.getDefaultInstance() : rateLimiter_;
  }
  /**
   * <pre>
   * Rate Limiter
   * </pre>
   *
   * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
   */
  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.language.v1.RateLimiterOrBuilder getRateLimiterOrBuilder() {
    return getRateLimiter();
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
    if (commonAttributes_ != null) {
      output.writeMessage(1, getCommonAttributes());
    }
    if (rateLimiter_ != null) {
      output.writeMessage(2, getRateLimiter());
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (commonAttributes_ != null) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(1, getCommonAttributes());
    }
    if (rateLimiter_ != null) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(2, getRateLimiter());
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
    if (!(obj instanceof com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper other = (com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper) obj;

    if (hasCommonAttributes() != other.hasCommonAttributes()) return false;
    if (hasCommonAttributes()) {
      if (!getCommonAttributes()
          .equals(other.getCommonAttributes())) return false;
    }
    if (hasRateLimiter() != other.hasRateLimiter()) return false;
    if (hasRateLimiter()) {
      if (!getRateLimiter()
          .equals(other.getRateLimiter())) return false;
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
    if (hasCommonAttributes()) {
      hash = (37 * hash) + COMMON_ATTRIBUTES_FIELD_NUMBER;
      hash = (53 * hash) + getCommonAttributes().hashCode();
    }
    if (hasRateLimiter()) {
      hash = (37 * hash) + RATE_LIMITER_FIELD_NUMBER;
      hash = (53 * hash) + getRateLimiter().hashCode();
    }
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper prototype) {
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
   * Protobuf type {@code aperture.policy.sync.v1.RateLimiterWrapper}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.policy.sync.v1.RateLimiterWrapper)
      com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapperOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterProto.internal_static_aperture_policy_sync_v1_RateLimiterWrapper_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterProto.internal_static_aperture_policy_sync_v1_RateLimiterWrapper_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper.class, com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper.newBuilder()
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
      if (commonAttributesBuilder_ == null) {
        commonAttributes_ = null;
      } else {
        commonAttributes_ = null;
        commonAttributesBuilder_ = null;
      }
      if (rateLimiterBuilder_ == null) {
        rateLimiter_ = null;
      } else {
        rateLimiter_ = null;
        rateLimiterBuilder_ = null;
      }
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterProto.internal_static_aperture_policy_sync_v1_RateLimiterWrapper_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper build() {
      com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper buildPartial() {
      com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper result = new com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper(this);
      if (commonAttributesBuilder_ == null) {
        result.commonAttributes_ = commonAttributes_;
      } else {
        result.commonAttributes_ = commonAttributesBuilder_.build();
      }
      if (rateLimiterBuilder_ == null) {
        result.rateLimiter_ = rateLimiter_;
      } else {
        result.rateLimiter_ = rateLimiterBuilder_.build();
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
      if (other instanceof com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper) {
        return mergeFrom((com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper other) {
      if (other == com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper.getDefaultInstance()) return this;
      if (other.hasCommonAttributes()) {
        mergeCommonAttributes(other.getCommonAttributes());
      }
      if (other.hasRateLimiter()) {
        mergeRateLimiter(other.getRateLimiter());
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
      com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes commonAttributes_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes, com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.Builder, com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesOrBuilder> commonAttributesBuilder_;
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     * @return Whether the commonAttributes field is set.
     */
    public boolean hasCommonAttributes() {
      return commonAttributesBuilder_ != null || commonAttributes_ != null;
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     * @return The commonAttributes.
     */
    public com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes getCommonAttributes() {
      if (commonAttributesBuilder_ == null) {
        return commonAttributes_ == null ? com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.getDefaultInstance() : commonAttributes_;
      } else {
        return commonAttributesBuilder_.getMessage();
      }
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     */
    public Builder setCommonAttributes(com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes value) {
      if (commonAttributesBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        commonAttributes_ = value;
        onChanged();
      } else {
        commonAttributesBuilder_.setMessage(value);
      }

      return this;
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     */
    public Builder setCommonAttributes(
        com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.Builder builderForValue) {
      if (commonAttributesBuilder_ == null) {
        commonAttributes_ = builderForValue.build();
        onChanged();
      } else {
        commonAttributesBuilder_.setMessage(builderForValue.build());
      }

      return this;
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     */
    public Builder mergeCommonAttributes(com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes value) {
      if (commonAttributesBuilder_ == null) {
        if (commonAttributes_ != null) {
          commonAttributes_ =
            com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.newBuilder(commonAttributes_).mergeFrom(value).buildPartial();
        } else {
          commonAttributes_ = value;
        }
        onChanged();
      } else {
        commonAttributesBuilder_.mergeFrom(value);
      }

      return this;
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     */
    public Builder clearCommonAttributes() {
      if (commonAttributesBuilder_ == null) {
        commonAttributes_ = null;
        onChanged();
      } else {
        commonAttributes_ = null;
        commonAttributesBuilder_ = null;
      }

      return this;
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     */
    public com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.Builder getCommonAttributesBuilder() {
      
      onChanged();
      return getCommonAttributesFieldBuilder().getBuilder();
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     */
    public com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesOrBuilder getCommonAttributesOrBuilder() {
      if (commonAttributesBuilder_ != null) {
        return commonAttributesBuilder_.getMessageOrBuilder();
      } else {
        return commonAttributes_ == null ?
            com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.getDefaultInstance() : commonAttributes_;
      }
    }
    /**
     * <pre>
     * CommonAttributes
     * </pre>
     *
     * <code>.aperture.policy.sync.v1.CommonAttributes common_attributes = 1 [json_name = "commonAttributes"];</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes, com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.Builder, com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesOrBuilder> 
        getCommonAttributesFieldBuilder() {
      if (commonAttributesBuilder_ == null) {
        commonAttributesBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes, com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributes.Builder, com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesOrBuilder>(
                getCommonAttributes(),
                getParentForChildren(),
                isClean());
        commonAttributes_ = null;
      }
      return commonAttributesBuilder_;
    }

    private com.fluxninja.generated.aperture.policy.language.v1.RateLimiter rateLimiter_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.policy.language.v1.RateLimiter, com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.Builder, com.fluxninja.generated.aperture.policy.language.v1.RateLimiterOrBuilder> rateLimiterBuilder_;
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     * @return Whether the rateLimiter field is set.
     */
    public boolean hasRateLimiter() {
      return rateLimiterBuilder_ != null || rateLimiter_ != null;
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     * @return The rateLimiter.
     */
    public com.fluxninja.generated.aperture.policy.language.v1.RateLimiter getRateLimiter() {
      if (rateLimiterBuilder_ == null) {
        return rateLimiter_ == null ? com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.getDefaultInstance() : rateLimiter_;
      } else {
        return rateLimiterBuilder_.getMessage();
      }
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     */
    public Builder setRateLimiter(com.fluxninja.generated.aperture.policy.language.v1.RateLimiter value) {
      if (rateLimiterBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        rateLimiter_ = value;
        onChanged();
      } else {
        rateLimiterBuilder_.setMessage(value);
      }

      return this;
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     */
    public Builder setRateLimiter(
        com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.Builder builderForValue) {
      if (rateLimiterBuilder_ == null) {
        rateLimiter_ = builderForValue.build();
        onChanged();
      } else {
        rateLimiterBuilder_.setMessage(builderForValue.build());
      }

      return this;
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     */
    public Builder mergeRateLimiter(com.fluxninja.generated.aperture.policy.language.v1.RateLimiter value) {
      if (rateLimiterBuilder_ == null) {
        if (rateLimiter_ != null) {
          rateLimiter_ =
            com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.newBuilder(rateLimiter_).mergeFrom(value).buildPartial();
        } else {
          rateLimiter_ = value;
        }
        onChanged();
      } else {
        rateLimiterBuilder_.mergeFrom(value);
      }

      return this;
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     */
    public Builder clearRateLimiter() {
      if (rateLimiterBuilder_ == null) {
        rateLimiter_ = null;
        onChanged();
      } else {
        rateLimiter_ = null;
        rateLimiterBuilder_ = null;
      }

      return this;
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     */
    public com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.Builder getRateLimiterBuilder() {
      
      onChanged();
      return getRateLimiterFieldBuilder().getBuilder();
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     */
    public com.fluxninja.generated.aperture.policy.language.v1.RateLimiterOrBuilder getRateLimiterOrBuilder() {
      if (rateLimiterBuilder_ != null) {
        return rateLimiterBuilder_.getMessageOrBuilder();
      } else {
        return rateLimiter_ == null ?
            com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.getDefaultInstance() : rateLimiter_;
      }
    }
    /**
     * <pre>
     * Rate Limiter
     * </pre>
     *
     * <code>.aperture.policy.language.v1.RateLimiter rate_limiter = 2 [json_name = "rateLimiter"];</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.fluxninja.generated.aperture.policy.language.v1.RateLimiter, com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.Builder, com.fluxninja.generated.aperture.policy.language.v1.RateLimiterOrBuilder> 
        getRateLimiterFieldBuilder() {
      if (rateLimiterBuilder_ == null) {
        rateLimiterBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.fluxninja.generated.aperture.policy.language.v1.RateLimiter, com.fluxninja.generated.aperture.policy.language.v1.RateLimiter.Builder, com.fluxninja.generated.aperture.policy.language.v1.RateLimiterOrBuilder>(
                getRateLimiter(),
                getParentForChildren(),
                isClean());
        rateLimiter_ = null;
      }
      return rateLimiterBuilder_;
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


    // @@protoc_insertion_point(builder_scope:aperture.policy.sync.v1.RateLimiterWrapper)
  }

  // @@protoc_insertion_point(class_scope:aperture.policy.sync.v1.RateLimiterWrapper)
  private static final com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper();
  }

  public static com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<RateLimiterWrapper>
      PARSER = new com.google.protobuf.AbstractParser<RateLimiterWrapper>() {
    @java.lang.Override
    public RateLimiterWrapper parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new RateLimiterWrapper(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<RateLimiterWrapper> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<RateLimiterWrapper> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.sync.v1.RateLimiterWrapper getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

