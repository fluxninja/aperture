// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: udpa/annotations/migrate.proto

package com.fluxninja.generated.udpa.annotations;

/**
 * Protobuf type {@code udpa.annotations.FieldMigrateAnnotation}
 */
public final class FieldMigrateAnnotation extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:udpa.annotations.FieldMigrateAnnotation)
    FieldMigrateAnnotationOrBuilder {
private static final long serialVersionUID = 0L;
  // Use FieldMigrateAnnotation.newBuilder() to construct.
  private FieldMigrateAnnotation(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private FieldMigrateAnnotation() {
    rename_ = "";
    oneofPromotion_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new FieldMigrateAnnotation();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private FieldMigrateAnnotation(
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

            rename_ = s;
            break;
          }
          case 18: {
            java.lang.String s = input.readStringRequireUtf8();

            oneofPromotion_ = s;
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
    return com.fluxninja.generated.udpa.annotations.MigrateProto.internal_static_udpa_annotations_FieldMigrateAnnotation_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.udpa.annotations.MigrateProto.internal_static_udpa_annotations_FieldMigrateAnnotation_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation.class, com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation.Builder.class);
  }

  public static final int RENAME_FIELD_NUMBER = 1;
  private volatile java.lang.Object rename_;
  /**
   * <pre>
   * Rename the field in next version.
   * </pre>
   *
   * <code>string rename = 1 [json_name = "rename"];</code>
   * @return The rename.
   */
  @java.lang.Override
  public java.lang.String getRename() {
    java.lang.Object ref = rename_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      rename_ = s;
      return s;
    }
  }
  /**
   * <pre>
   * Rename the field in next version.
   * </pre>
   *
   * <code>string rename = 1 [json_name = "rename"];</code>
   * @return The bytes for rename.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getRenameBytes() {
    java.lang.Object ref = rename_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      rename_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int ONEOF_PROMOTION_FIELD_NUMBER = 2;
  private volatile java.lang.Object oneofPromotion_;
  /**
   * <pre>
   * Add the field to a named oneof in next version. If this already exists, the
   * field will join its siblings under the oneof, otherwise a new oneof will be
   * created with the given name.
   * </pre>
   *
   * <code>string oneof_promotion = 2 [json_name = "oneofPromotion"];</code>
   * @return The oneofPromotion.
   */
  @java.lang.Override
  public java.lang.String getOneofPromotion() {
    java.lang.Object ref = oneofPromotion_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      oneofPromotion_ = s;
      return s;
    }
  }
  /**
   * <pre>
   * Add the field to a named oneof in next version. If this already exists, the
   * field will join its siblings under the oneof, otherwise a new oneof will be
   * created with the given name.
   * </pre>
   *
   * <code>string oneof_promotion = 2 [json_name = "oneofPromotion"];</code>
   * @return The bytes for oneofPromotion.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getOneofPromotionBytes() {
    java.lang.Object ref = oneofPromotion_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      oneofPromotion_ = b;
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
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(rename_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 1, rename_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(oneofPromotion_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 2, oneofPromotion_);
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(rename_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, rename_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(oneofPromotion_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(2, oneofPromotion_);
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
    if (!(obj instanceof com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation other = (com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation) obj;

    if (!getRename()
        .equals(other.getRename())) return false;
    if (!getOneofPromotion()
        .equals(other.getOneofPromotion())) return false;
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
    hash = (37 * hash) + RENAME_FIELD_NUMBER;
    hash = (53 * hash) + getRename().hashCode();
    hash = (37 * hash) + ONEOF_PROMOTION_FIELD_NUMBER;
    hash = (53 * hash) + getOneofPromotion().hashCode();
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation prototype) {
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
   * Protobuf type {@code udpa.annotations.FieldMigrateAnnotation}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:udpa.annotations.FieldMigrateAnnotation)
      com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotationOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.udpa.annotations.MigrateProto.internal_static_udpa_annotations_FieldMigrateAnnotation_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.udpa.annotations.MigrateProto.internal_static_udpa_annotations_FieldMigrateAnnotation_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation.class, com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation.Builder.class);
    }

    // Construct using com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation.newBuilder()
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
      rename_ = "";

      oneofPromotion_ = "";

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.udpa.annotations.MigrateProto.internal_static_udpa_annotations_FieldMigrateAnnotation_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation getDefaultInstanceForType() {
      return com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation build() {
      com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation buildPartial() {
      com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation result = new com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation(this);
      result.rename_ = rename_;
      result.oneofPromotion_ = oneofPromotion_;
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
      if (other instanceof com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation) {
        return mergeFrom((com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation other) {
      if (other == com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation.getDefaultInstance()) return this;
      if (!other.getRename().isEmpty()) {
        rename_ = other.rename_;
        onChanged();
      }
      if (!other.getOneofPromotion().isEmpty()) {
        oneofPromotion_ = other.oneofPromotion_;
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
      com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private java.lang.Object rename_ = "";
    /**
     * <pre>
     * Rename the field in next version.
     * </pre>
     *
     * <code>string rename = 1 [json_name = "rename"];</code>
     * @return The rename.
     */
    public java.lang.String getRename() {
      java.lang.Object ref = rename_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        rename_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <pre>
     * Rename the field in next version.
     * </pre>
     *
     * <code>string rename = 1 [json_name = "rename"];</code>
     * @return The bytes for rename.
     */
    public com.google.protobuf.ByteString
        getRenameBytes() {
      java.lang.Object ref = rename_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        rename_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <pre>
     * Rename the field in next version.
     * </pre>
     *
     * <code>string rename = 1 [json_name = "rename"];</code>
     * @param value The rename to set.
     * @return This builder for chaining.
     */
    public Builder setRename(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      rename_ = value;
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Rename the field in next version.
     * </pre>
     *
     * <code>string rename = 1 [json_name = "rename"];</code>
     * @return This builder for chaining.
     */
    public Builder clearRename() {
      
      rename_ = getDefaultInstance().getRename();
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Rename the field in next version.
     * </pre>
     *
     * <code>string rename = 1 [json_name = "rename"];</code>
     * @param value The bytes for rename to set.
     * @return This builder for chaining.
     */
    public Builder setRenameBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      rename_ = value;
      onChanged();
      return this;
    }

    private java.lang.Object oneofPromotion_ = "";
    /**
     * <pre>
     * Add the field to a named oneof in next version. If this already exists, the
     * field will join its siblings under the oneof, otherwise a new oneof will be
     * created with the given name.
     * </pre>
     *
     * <code>string oneof_promotion = 2 [json_name = "oneofPromotion"];</code>
     * @return The oneofPromotion.
     */
    public java.lang.String getOneofPromotion() {
      java.lang.Object ref = oneofPromotion_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        oneofPromotion_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <pre>
     * Add the field to a named oneof in next version. If this already exists, the
     * field will join its siblings under the oneof, otherwise a new oneof will be
     * created with the given name.
     * </pre>
     *
     * <code>string oneof_promotion = 2 [json_name = "oneofPromotion"];</code>
     * @return The bytes for oneofPromotion.
     */
    public com.google.protobuf.ByteString
        getOneofPromotionBytes() {
      java.lang.Object ref = oneofPromotion_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        oneofPromotion_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <pre>
     * Add the field to a named oneof in next version. If this already exists, the
     * field will join its siblings under the oneof, otherwise a new oneof will be
     * created with the given name.
     * </pre>
     *
     * <code>string oneof_promotion = 2 [json_name = "oneofPromotion"];</code>
     * @param value The oneofPromotion to set.
     * @return This builder for chaining.
     */
    public Builder setOneofPromotion(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      oneofPromotion_ = value;
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Add the field to a named oneof in next version. If this already exists, the
     * field will join its siblings under the oneof, otherwise a new oneof will be
     * created with the given name.
     * </pre>
     *
     * <code>string oneof_promotion = 2 [json_name = "oneofPromotion"];</code>
     * @return This builder for chaining.
     */
    public Builder clearOneofPromotion() {
      
      oneofPromotion_ = getDefaultInstance().getOneofPromotion();
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Add the field to a named oneof in next version. If this already exists, the
     * field will join its siblings under the oneof, otherwise a new oneof will be
     * created with the given name.
     * </pre>
     *
     * <code>string oneof_promotion = 2 [json_name = "oneofPromotion"];</code>
     * @param value The bytes for oneofPromotion to set.
     * @return This builder for chaining.
     */
    public Builder setOneofPromotionBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      oneofPromotion_ = value;
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


    // @@protoc_insertion_point(builder_scope:udpa.annotations.FieldMigrateAnnotation)
  }

  // @@protoc_insertion_point(class_scope:udpa.annotations.FieldMigrateAnnotation)
  private static final com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation();
  }

  public static com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<FieldMigrateAnnotation>
      PARSER = new com.google.protobuf.AbstractParser<FieldMigrateAnnotation>() {
    @java.lang.Override
    public FieldMigrateAnnotation parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new FieldMigrateAnnotation(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<FieldMigrateAnnotation> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<FieldMigrateAnnotation> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.udpa.annotations.FieldMigrateAnnotation getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

