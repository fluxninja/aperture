// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/v1/flowcontrol.proto

package com.fluxninja.aperture.flowcontrol.v1;

/**
 * Protobuf type {@code aperture.flowcontrol.v1.ControlPointInfo}
 */
public final class ControlPointInfo extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.flowcontrol.v1.ControlPointInfo)
    ControlPointInfoOrBuilder {
private static final long serialVersionUID = 0L;
  // Use ControlPointInfo.newBuilder() to construct.
  private ControlPointInfo(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private ControlPointInfo() {
    type_ = 0;
    feature_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new ControlPointInfo();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_ControlPointInfo_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_ControlPointInfo_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.class, com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Builder.class);
  }

  /**
   * <pre>
   * Type contains fields that represent type of ControlPointInfo.
   * </pre>
   *
   * Protobuf enum {@code aperture.flowcontrol.v1.ControlPointInfo.Type}
   */
  public enum Type
      implements com.google.protobuf.ProtocolMessageEnum {
    /**
     * <code>TYPE_UNKNOWN = 0;</code>
     */
    TYPE_UNKNOWN(0),
    /**
     * <code>TYPE_FEATURE = 1;</code>
     */
    TYPE_FEATURE(1),
    /**
     * <code>TYPE_INGRESS = 2;</code>
     */
    TYPE_INGRESS(2),
    /**
     * <code>TYPE_EGRESS = 3;</code>
     */
    TYPE_EGRESS(3),
    UNRECOGNIZED(-1),
    ;

    /**
     * <code>TYPE_UNKNOWN = 0;</code>
     */
    public static final int TYPE_UNKNOWN_VALUE = 0;
    /**
     * <code>TYPE_FEATURE = 1;</code>
     */
    public static final int TYPE_FEATURE_VALUE = 1;
    /**
     * <code>TYPE_INGRESS = 2;</code>
     */
    public static final int TYPE_INGRESS_VALUE = 2;
    /**
     * <code>TYPE_EGRESS = 3;</code>
     */
    public static final int TYPE_EGRESS_VALUE = 3;


    public final int getNumber() {
      if (this == UNRECOGNIZED) {
        throw new java.lang.IllegalArgumentException(
            "Can't get the number of an unknown enum value.");
      }
      return value;
    }

    /**
     * @param value The numeric wire value of the corresponding enum entry.
     * @return The enum associated with the given numeric wire value.
     * @deprecated Use {@link #forNumber(int)} instead.
     */
    @java.lang.Deprecated
    public static Type valueOf(int value) {
      return forNumber(value);
    }

    /**
     * @param value The numeric wire value of the corresponding enum entry.
     * @return The enum associated with the given numeric wire value.
     */
    public static Type forNumber(int value) {
      switch (value) {
        case 0: return TYPE_UNKNOWN;
        case 1: return TYPE_FEATURE;
        case 2: return TYPE_INGRESS;
        case 3: return TYPE_EGRESS;
        default: return null;
      }
    }

    public static com.google.protobuf.Internal.EnumLiteMap<Type>
        internalGetValueMap() {
      return internalValueMap;
    }
    private static final com.google.protobuf.Internal.EnumLiteMap<
        Type> internalValueMap =
          new com.google.protobuf.Internal.EnumLiteMap<Type>() {
            public Type findValueByNumber(int number) {
              return Type.forNumber(number);
            }
          };

    public final com.google.protobuf.Descriptors.EnumValueDescriptor
        getValueDescriptor() {
      if (this == UNRECOGNIZED) {
        throw new java.lang.IllegalStateException(
            "Can't get the descriptor of an unrecognized enum value.");
      }
      return getDescriptor().getValues().get(ordinal());
    }
    public final com.google.protobuf.Descriptors.EnumDescriptor
        getDescriptorForType() {
      return getDescriptor();
    }
    public static final com.google.protobuf.Descriptors.EnumDescriptor
        getDescriptor() {
      return com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.getDescriptor().getEnumTypes().get(0);
    }

    private static final Type[] VALUES = values();

    public static Type valueOf(
        com.google.protobuf.Descriptors.EnumValueDescriptor desc) {
      if (desc.getType() != getDescriptor()) {
        throw new java.lang.IllegalArgumentException(
          "EnumValueDescriptor is not for this type.");
      }
      if (desc.getIndex() == -1) {
        return UNRECOGNIZED;
      }
      return VALUES[desc.getIndex()];
    }

    private final int value;

    private Type(int value) {
      this.value = value;
    }

    // @@protoc_insertion_point(enum_scope:aperture.flowcontrol.v1.ControlPointInfo.Type)
  }

  public static final int TYPE_FIELD_NUMBER = 1;
  private int type_;
  /**
   * <code>.aperture.flowcontrol.v1.ControlPointInfo.Type type = 1 [json_name = "type"];</code>
   * @return The enum numeric value on the wire for type.
   */
  @java.lang.Override public int getTypeValue() {
    return type_;
  }
  /**
   * <code>.aperture.flowcontrol.v1.ControlPointInfo.Type type = 1 [json_name = "type"];</code>
   * @return The type.
   */
  @java.lang.Override public com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type getType() {
    @SuppressWarnings("deprecation")
    com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type result = com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type.valueOf(type_);
    return result == null ? com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type.UNRECOGNIZED : result;
  }

  public static final int FEATURE_FIELD_NUMBER = 2;
  private volatile java.lang.Object feature_;
  /**
   * <code>string feature = 2 [json_name = "feature"];</code>
   * @return The feature.
   */
  @java.lang.Override
  public java.lang.String getFeature() {
    java.lang.Object ref = feature_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs =
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      feature_ = s;
      return s;
    }
  }
  /**
   * <code>string feature = 2 [json_name = "feature"];</code>
   * @return The bytes for feature.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getFeatureBytes() {
    java.lang.Object ref = feature_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b =
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      feature_ = b;
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
    if (type_ != com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type.TYPE_UNKNOWN.getNumber()) {
      output.writeEnum(1, type_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(feature_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 2, feature_);
    }
    getUnknownFields().writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (type_ != com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type.TYPE_UNKNOWN.getNumber()) {
      size += com.google.protobuf.CodedOutputStream
        .computeEnumSize(1, type_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(feature_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(2, feature_);
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
    if (!(obj instanceof com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo)) {
      return super.equals(obj);
    }
    com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo other = (com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo) obj;

    if (type_ != other.type_) return false;
    if (!getFeature()
        .equals(other.getFeature())) return false;
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
    hash = (37 * hash) + TYPE_FIELD_NUMBER;
    hash = (53 * hash) + type_;
    hash = (37 * hash) + FEATURE_FIELD_NUMBER;
    hash = (53 * hash) + getFeature().hashCode();
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo parseFrom(
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
  public static Builder newBuilder(com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo prototype) {
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
   * Protobuf type {@code aperture.flowcontrol.v1.ControlPointInfo}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.flowcontrol.v1.ControlPointInfo)
      com.fluxninja.aperture.flowcontrol.v1.ControlPointInfoOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_ControlPointInfo_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_ControlPointInfo_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.class, com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Builder.class);
    }

    // Construct using com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.newBuilder()
    private Builder() {

    }

    private Builder(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      super(parent);

    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      type_ = 0;

      feature_ = "";

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.aperture.flowcontrol.v1.FlowcontrolProto.internal_static_aperture_flowcontrol_v1_ControlPointInfo_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo getDefaultInstanceForType() {
      return com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo build() {
      com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo buildPartial() {
      com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo result = new com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo(this);
      result.type_ = type_;
      result.feature_ = feature_;
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
      if (other instanceof com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo) {
        return mergeFrom((com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo other) {
      if (other == com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.getDefaultInstance()) return this;
      if (other.type_ != 0) {
        setTypeValue(other.getTypeValue());
      }
      if (!other.getFeature().isEmpty()) {
        feature_ = other.feature_;
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
            case 8: {
              type_ = input.readEnum();

              break;
            } // case 8
            case 18: {
              feature_ = input.readStringRequireUtf8();

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

    private int type_ = 0;
    /**
     * <code>.aperture.flowcontrol.v1.ControlPointInfo.Type type = 1 [json_name = "type"];</code>
     * @return The enum numeric value on the wire for type.
     */
    @java.lang.Override public int getTypeValue() {
      return type_;
    }
    /**
     * <code>.aperture.flowcontrol.v1.ControlPointInfo.Type type = 1 [json_name = "type"];</code>
     * @param value The enum numeric value on the wire for type to set.
     * @return This builder for chaining.
     */
    public Builder setTypeValue(int value) {

      type_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.v1.ControlPointInfo.Type type = 1 [json_name = "type"];</code>
     * @return The type.
     */
    @java.lang.Override
    public com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type getType() {
      @SuppressWarnings("deprecation")
      com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type result = com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type.valueOf(type_);
      return result == null ? com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type.UNRECOGNIZED : result;
    }
    /**
     * <code>.aperture.flowcontrol.v1.ControlPointInfo.Type type = 1 [json_name = "type"];</code>
     * @param value The type to set.
     * @return This builder for chaining.
     */
    public Builder setType(com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo.Type value) {
      if (value == null) {
        throw new NullPointerException();
      }

      type_ = value.getNumber();
      onChanged();
      return this;
    }
    /**
     * <code>.aperture.flowcontrol.v1.ControlPointInfo.Type type = 1 [json_name = "type"];</code>
     * @return This builder for chaining.
     */
    public Builder clearType() {

      type_ = 0;
      onChanged();
      return this;
    }

    private java.lang.Object feature_ = "";
    /**
     * <code>string feature = 2 [json_name = "feature"];</code>
     * @return The feature.
     */
    public java.lang.String getFeature() {
      java.lang.Object ref = feature_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        feature_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string feature = 2 [json_name = "feature"];</code>
     * @return The bytes for feature.
     */
    public com.google.protobuf.ByteString
        getFeatureBytes() {
      java.lang.Object ref = feature_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b =
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        feature_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string feature = 2 [json_name = "feature"];</code>
     * @param value The feature to set.
     * @return This builder for chaining.
     */
    public Builder setFeature(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }

      feature_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>string feature = 2 [json_name = "feature"];</code>
     * @return This builder for chaining.
     */
    public Builder clearFeature() {

      feature_ = getDefaultInstance().getFeature();
      onChanged();
      return this;
    }
    /**
     * <code>string feature = 2 [json_name = "feature"];</code>
     * @param value The bytes for feature to set.
     * @return This builder for chaining.
     */
    public Builder setFeatureBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);

      feature_ = value;
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


    // @@protoc_insertion_point(builder_scope:aperture.flowcontrol.v1.ControlPointInfo)
  }

  // @@protoc_insertion_point(class_scope:aperture.flowcontrol.v1.ControlPointInfo)
  private static final com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo();
  }

  public static com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<ControlPointInfo>
      PARSER = new com.google.protobuf.AbstractParser<ControlPointInfo>() {
    @java.lang.Override
    public ControlPointInfo parsePartialFrom(
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

  public static com.google.protobuf.Parser<ControlPointInfo> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<ControlPointInfo> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.aperture.flowcontrol.v1.ControlPointInfo getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}
