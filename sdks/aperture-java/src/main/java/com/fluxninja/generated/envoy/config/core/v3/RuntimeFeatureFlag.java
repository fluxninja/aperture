// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/config/core/v3/base.proto

package com.fluxninja.generated.envoy.config.core.v3;

/**
 * <pre>
 * Runtime derived bool with a default when not specified.
 * </pre>
 *
 * Protobuf type {@code envoy.config.core.v3.RuntimeFeatureFlag}
 */
public final class RuntimeFeatureFlag extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:envoy.config.core.v3.RuntimeFeatureFlag)
    RuntimeFeatureFlagOrBuilder {
private static final long serialVersionUID = 0L;
  // Use RuntimeFeatureFlag.newBuilder() to construct.
  private RuntimeFeatureFlag(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private RuntimeFeatureFlag() {
    runtimeKey_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new RuntimeFeatureFlag();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private RuntimeFeatureFlag(
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
            com.google.protobuf.BoolValue.Builder subBuilder = null;
            if (defaultValue_ != null) {
              subBuilder = defaultValue_.toBuilder();
            }
            defaultValue_ = input.readMessage(com.google.protobuf.BoolValue.parser(), extensionRegistry);
            if (subBuilder != null) {
              subBuilder.mergeFrom(defaultValue_);
              defaultValue_ = subBuilder.buildPartial();
            }

            break;
          }
          case 18: {
            java.lang.String s = input.readStringRequireUtf8();

            runtimeKey_ = s;
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
    return com.fluxninja.generated.envoy.config.core.v3.BaseProto.internal_static_envoy_config_core_v3_RuntimeFeatureFlag_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.envoy.config.core.v3.BaseProto.internal_static_envoy_config_core_v3_RuntimeFeatureFlag_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag.class, com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag.Builder.class);
  }

  public static final int DEFAULT_VALUE_FIELD_NUMBER = 1;
  private com.google.protobuf.BoolValue defaultValue_;
  /**
   * <pre>
   * Default value if runtime value is not available.
   * </pre>
   *
   * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
   * @return Whether the defaultValue field is set.
   */
  @java.lang.Override
  public boolean hasDefaultValue() {
    return defaultValue_ != null;
  }
  /**
   * <pre>
   * Default value if runtime value is not available.
   * </pre>
   *
   * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
   * @return The defaultValue.
   */
  @java.lang.Override
  public com.google.protobuf.BoolValue getDefaultValue() {
    return defaultValue_ == null ? com.google.protobuf.BoolValue.getDefaultInstance() : defaultValue_;
  }
  /**
   * <pre>
   * Default value if runtime value is not available.
   * </pre>
   *
   * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
   */
  @java.lang.Override
  public com.google.protobuf.BoolValueOrBuilder getDefaultValueOrBuilder() {
    return getDefaultValue();
  }

  public static final int RUNTIME_KEY_FIELD_NUMBER = 2;
  private volatile java.lang.Object runtimeKey_;
  /**
   * <pre>
   * Runtime key to get value for comparison. This value is used if defined. The boolean value must
   * be represented via its
   * `canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_.
   * </pre>
   *
   * <code>string runtime_key = 2 [json_name = "runtimeKey", (.validate.rules) = { ... }</code>
   * @return The runtimeKey.
   */
  @java.lang.Override
  public java.lang.String getRuntimeKey() {
    java.lang.Object ref = runtimeKey_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      runtimeKey_ = s;
      return s;
    }
  }
  /**
   * <pre>
   * Runtime key to get value for comparison. This value is used if defined. The boolean value must
   * be represented via its
   * `canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_.
   * </pre>
   *
   * <code>string runtime_key = 2 [json_name = "runtimeKey", (.validate.rules) = { ... }</code>
   * @return The bytes for runtimeKey.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getRuntimeKeyBytes() {
    java.lang.Object ref = runtimeKey_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      runtimeKey_ = b;
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
    if (defaultValue_ != null) {
      output.writeMessage(1, getDefaultValue());
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(runtimeKey_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 2, runtimeKey_);
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (defaultValue_ != null) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(1, getDefaultValue());
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(runtimeKey_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(2, runtimeKey_);
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
    if (!(obj instanceof com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag other = (com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag) obj;

    if (hasDefaultValue() != other.hasDefaultValue()) return false;
    if (hasDefaultValue()) {
      if (!getDefaultValue()
          .equals(other.getDefaultValue())) return false;
    }
    if (!getRuntimeKey()
        .equals(other.getRuntimeKey())) return false;
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
    if (hasDefaultValue()) {
      hash = (37 * hash) + DEFAULT_VALUE_FIELD_NUMBER;
      hash = (53 * hash) + getDefaultValue().hashCode();
    }
    hash = (37 * hash) + RUNTIME_KEY_FIELD_NUMBER;
    hash = (53 * hash) + getRuntimeKey().hashCode();
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag prototype) {
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
   * Runtime derived bool with a default when not specified.
   * </pre>
   *
   * Protobuf type {@code envoy.config.core.v3.RuntimeFeatureFlag}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:envoy.config.core.v3.RuntimeFeatureFlag)
      com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlagOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.envoy.config.core.v3.BaseProto.internal_static_envoy_config_core_v3_RuntimeFeatureFlag_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.envoy.config.core.v3.BaseProto.internal_static_envoy_config_core_v3_RuntimeFeatureFlag_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag.class, com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag.Builder.class);
    }

    // Construct using com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag.newBuilder()
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
      if (defaultValueBuilder_ == null) {
        defaultValue_ = null;
      } else {
        defaultValue_ = null;
        defaultValueBuilder_ = null;
      }
      runtimeKey_ = "";

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.envoy.config.core.v3.BaseProto.internal_static_envoy_config_core_v3_RuntimeFeatureFlag_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag getDefaultInstanceForType() {
      return com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag build() {
      com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag buildPartial() {
      com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag result = new com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag(this);
      if (defaultValueBuilder_ == null) {
        result.defaultValue_ = defaultValue_;
      } else {
        result.defaultValue_ = defaultValueBuilder_.build();
      }
      result.runtimeKey_ = runtimeKey_;
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
      if (other instanceof com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag) {
        return mergeFrom((com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag other) {
      if (other == com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag.getDefaultInstance()) return this;
      if (other.hasDefaultValue()) {
        mergeDefaultValue(other.getDefaultValue());
      }
      if (!other.getRuntimeKey().isEmpty()) {
        runtimeKey_ = other.runtimeKey_;
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
      com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private com.google.protobuf.BoolValue defaultValue_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.google.protobuf.BoolValue, com.google.protobuf.BoolValue.Builder, com.google.protobuf.BoolValueOrBuilder> defaultValueBuilder_;
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     * @return Whether the defaultValue field is set.
     */
    public boolean hasDefaultValue() {
      return defaultValueBuilder_ != null || defaultValue_ != null;
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     * @return The defaultValue.
     */
    public com.google.protobuf.BoolValue getDefaultValue() {
      if (defaultValueBuilder_ == null) {
        return defaultValue_ == null ? com.google.protobuf.BoolValue.getDefaultInstance() : defaultValue_;
      } else {
        return defaultValueBuilder_.getMessage();
      }
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     */
    public Builder setDefaultValue(com.google.protobuf.BoolValue value) {
      if (defaultValueBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        defaultValue_ = value;
        onChanged();
      } else {
        defaultValueBuilder_.setMessage(value);
      }

      return this;
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     */
    public Builder setDefaultValue(
        com.google.protobuf.BoolValue.Builder builderForValue) {
      if (defaultValueBuilder_ == null) {
        defaultValue_ = builderForValue.build();
        onChanged();
      } else {
        defaultValueBuilder_.setMessage(builderForValue.build());
      }

      return this;
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     */
    public Builder mergeDefaultValue(com.google.protobuf.BoolValue value) {
      if (defaultValueBuilder_ == null) {
        if (defaultValue_ != null) {
          defaultValue_ =
            com.google.protobuf.BoolValue.newBuilder(defaultValue_).mergeFrom(value).buildPartial();
        } else {
          defaultValue_ = value;
        }
        onChanged();
      } else {
        defaultValueBuilder_.mergeFrom(value);
      }

      return this;
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     */
    public Builder clearDefaultValue() {
      if (defaultValueBuilder_ == null) {
        defaultValue_ = null;
        onChanged();
      } else {
        defaultValue_ = null;
        defaultValueBuilder_ = null;
      }

      return this;
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     */
    public com.google.protobuf.BoolValue.Builder getDefaultValueBuilder() {
      
      onChanged();
      return getDefaultValueFieldBuilder().getBuilder();
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     */
    public com.google.protobuf.BoolValueOrBuilder getDefaultValueOrBuilder() {
      if (defaultValueBuilder_ != null) {
        return defaultValueBuilder_.getMessageOrBuilder();
      } else {
        return defaultValue_ == null ?
            com.google.protobuf.BoolValue.getDefaultInstance() : defaultValue_;
      }
    }
    /**
     * <pre>
     * Default value if runtime value is not available.
     * </pre>
     *
     * <code>.google.protobuf.BoolValue default_value = 1 [json_name = "defaultValue", (.validate.rules) = { ... }</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.google.protobuf.BoolValue, com.google.protobuf.BoolValue.Builder, com.google.protobuf.BoolValueOrBuilder> 
        getDefaultValueFieldBuilder() {
      if (defaultValueBuilder_ == null) {
        defaultValueBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.google.protobuf.BoolValue, com.google.protobuf.BoolValue.Builder, com.google.protobuf.BoolValueOrBuilder>(
                getDefaultValue(),
                getParentForChildren(),
                isClean());
        defaultValue_ = null;
      }
      return defaultValueBuilder_;
    }

    private java.lang.Object runtimeKey_ = "";
    /**
     * <pre>
     * Runtime key to get value for comparison. This value is used if defined. The boolean value must
     * be represented via its
     * `canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_.
     * </pre>
     *
     * <code>string runtime_key = 2 [json_name = "runtimeKey", (.validate.rules) = { ... }</code>
     * @return The runtimeKey.
     */
    public java.lang.String getRuntimeKey() {
      java.lang.Object ref = runtimeKey_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        runtimeKey_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <pre>
     * Runtime key to get value for comparison. This value is used if defined. The boolean value must
     * be represented via its
     * `canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_.
     * </pre>
     *
     * <code>string runtime_key = 2 [json_name = "runtimeKey", (.validate.rules) = { ... }</code>
     * @return The bytes for runtimeKey.
     */
    public com.google.protobuf.ByteString
        getRuntimeKeyBytes() {
      java.lang.Object ref = runtimeKey_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        runtimeKey_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <pre>
     * Runtime key to get value for comparison. This value is used if defined. The boolean value must
     * be represented via its
     * `canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_.
     * </pre>
     *
     * <code>string runtime_key = 2 [json_name = "runtimeKey", (.validate.rules) = { ... }</code>
     * @param value The runtimeKey to set.
     * @return This builder for chaining.
     */
    public Builder setRuntimeKey(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      runtimeKey_ = value;
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Runtime key to get value for comparison. This value is used if defined. The boolean value must
     * be represented via its
     * `canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_.
     * </pre>
     *
     * <code>string runtime_key = 2 [json_name = "runtimeKey", (.validate.rules) = { ... }</code>
     * @return This builder for chaining.
     */
    public Builder clearRuntimeKey() {
      
      runtimeKey_ = getDefaultInstance().getRuntimeKey();
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Runtime key to get value for comparison. This value is used if defined. The boolean value must
     * be represented via its
     * `canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_.
     * </pre>
     *
     * <code>string runtime_key = 2 [json_name = "runtimeKey", (.validate.rules) = { ... }</code>
     * @param value The bytes for runtimeKey to set.
     * @return This builder for chaining.
     */
    public Builder setRuntimeKeyBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      runtimeKey_ = value;
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


    // @@protoc_insertion_point(builder_scope:envoy.config.core.v3.RuntimeFeatureFlag)
  }

  // @@protoc_insertion_point(class_scope:envoy.config.core.v3.RuntimeFeatureFlag)
  private static final com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag();
  }

  public static com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<RuntimeFeatureFlag>
      PARSER = new com.google.protobuf.AbstractParser<RuntimeFeatureFlag>() {
    @java.lang.Override
    public RuntimeFeatureFlag parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new RuntimeFeatureFlag(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<RuntimeFeatureFlag> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<RuntimeFeatureFlag> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.envoy.config.core.v3.RuntimeFeatureFlag getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

