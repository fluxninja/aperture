// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/policy.proto

package com.aperture.policy.language.v1;

/**
 * <pre>
 * FluxMeter gathers metrics for the traffic that matches its selector.
 * </pre>
 *
 * Protobuf type {@code aperture.policy.language.v1.FluxMeter}
 */
public final class FluxMeter extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.policy.language.v1.FluxMeter)
    FluxMeterOrBuilder {
private static final long serialVersionUID = 0L;
  // Use FluxMeter.newBuilder() to construct.
  private FluxMeter(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private FluxMeter() {
    name_ = "";
    histogramBuckets_ = emptyDoubleList();
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new FluxMeter();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private FluxMeter(
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
            java.lang.String s = input.readStringRequireUtf8();

            name_ = s;
            break;
          }
          case 18: {
            com.aperture.policy.language.v1.Selector.Builder subBuilder = null;
            if (selector_ != null) {
              subBuilder = selector_.toBuilder();
            }
            selector_ = input.readMessage(com.aperture.policy.language.v1.Selector.parser(), extensionRegistry);
            if (subBuilder != null) {
              subBuilder.mergeFrom(selector_);
              selector_ = subBuilder.buildPartial();
            }

            break;
          }
          case 25: {
            if (!((mutable_bitField0_ & 0x00000001) != 0)) {
              histogramBuckets_ = newDoubleList();
              mutable_bitField0_ |= 0x00000001;
            }
            histogramBuckets_.addDouble(input.readDouble());
            break;
          }
          case 26: {
            int length = input.readRawVarint32();
            int limit = input.pushLimit(length);
            if (!((mutable_bitField0_ & 0x00000001) != 0) && input.getBytesUntilLimit() > 0) {
              histogramBuckets_ = newDoubleList();
              mutable_bitField0_ |= 0x00000001;
            }
            while (input.getBytesUntilLimit() > 0) {
              histogramBuckets_.addDouble(input.readDouble());
            }
            input.popLimit(limit);
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
        histogramBuckets_.makeImmutable(); // C
      }
      this.unknownFields = unknownFields.build();
      makeExtensionsImmutable();
    }
  }
  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_FluxMeter_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_FluxMeter_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.aperture.policy.language.v1.FluxMeter.class, com.aperture.policy.language.v1.FluxMeter.Builder.class);
  }

  public static final int NAME_FIELD_NUMBER = 1;
  private volatile java.lang.Object name_;
  /**
   * <pre>
   * Name of the flux meter.
   * </pre>
   *
   * <code>string name = 1 [json_name = "name"];</code>
   * @return The name.
   */
  @java.lang.Override
  public java.lang.String getName() {
    java.lang.Object ref = name_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs =
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      name_ = s;
      return s;
    }
  }
  /**
   * <pre>
   * Name of the flux meter.
   * </pre>
   *
   * <code>string name = 1 [json_name = "name"];</code>
   * @return The bytes for name.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getNameBytes() {
    java.lang.Object ref = name_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b =
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      name_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int SELECTOR_FIELD_NUMBER = 2;
  private com.aperture.policy.language.v1.Selector selector_;
  /**
   * <pre>
   * Policies are only applied to flows that are matched based on the fields in the selector.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
   * @return Whether the selector field is set.
   */
  @java.lang.Override
  public boolean hasSelector() {
    return selector_ != null;
  }
  /**
   * <pre>
   * Policies are only applied to flows that are matched based on the fields in the selector.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
   * @return The selector.
   */
  @java.lang.Override
  public com.aperture.policy.language.v1.Selector getSelector() {
    return selector_ == null ? com.aperture.policy.language.v1.Selector.getDefaultInstance() : selector_;
  }
  /**
   * <pre>
   * Policies are only applied to flows that are matched based on the fields in the selector.
   * </pre>
   *
   * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
   */
  @java.lang.Override
  public com.aperture.policy.language.v1.SelectorOrBuilder getSelectorOrBuilder() {
    return getSelector();
  }

  public static final int HISTOGRAM_BUCKETS_FIELD_NUMBER = 3;
  private com.google.protobuf.Internal.DoubleList histogramBuckets_;
  /**
   * <pre>
   * Latency histogram buckets (in ms) for this FluxMeter.
   * </pre>
   *
   * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return A list containing the histogramBuckets.
   */
  @java.lang.Override
  public java.util.List<java.lang.Double>
      getHistogramBucketsList() {
    return histogramBuckets_;
  }
  /**
   * <pre>
   * Latency histogram buckets (in ms) for this FluxMeter.
   * </pre>
   *
   * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @return The count of histogramBuckets.
   */
  public int getHistogramBucketsCount() {
    return histogramBuckets_.size();
  }
  /**
   * <pre>
   * Latency histogram buckets (in ms) for this FluxMeter.
   * </pre>
   *
   * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
   * @param index The index of the element to return.
   * @return The histogramBuckets at the given index.
   */
  public double getHistogramBuckets(int index) {
    return histogramBuckets_.getDouble(index);
  }
  private int histogramBucketsMemoizedSerializedSize = -1;

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
    getSerializedSize();
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(name_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 1, name_);
    }
    if (selector_ != null) {
      output.writeMessage(2, getSelector());
    }
    if (getHistogramBucketsList().size() > 0) {
      output.writeUInt32NoTag(26);
      output.writeUInt32NoTag(histogramBucketsMemoizedSerializedSize);
    }
    for (int i = 0; i < histogramBuckets_.size(); i++) {
      output.writeDoubleNoTag(histogramBuckets_.getDouble(i));
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(name_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, name_);
    }
    if (selector_ != null) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(2, getSelector());
    }
    {
      int dataSize = 0;
      dataSize = 8 * getHistogramBucketsList().size();
      size += dataSize;
      if (!getHistogramBucketsList().isEmpty()) {
        size += 1;
        size += com.google.protobuf.CodedOutputStream
            .computeInt32SizeNoTag(dataSize);
      }
      histogramBucketsMemoizedSerializedSize = dataSize;
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
    if (!(obj instanceof com.aperture.policy.language.v1.FluxMeter)) {
      return super.equals(obj);
    }
    com.aperture.policy.language.v1.FluxMeter other = (com.aperture.policy.language.v1.FluxMeter) obj;

    if (!getName()
        .equals(other.getName())) return false;
    if (hasSelector() != other.hasSelector()) return false;
    if (hasSelector()) {
      if (!getSelector()
          .equals(other.getSelector())) return false;
    }
    if (!getHistogramBucketsList()
        .equals(other.getHistogramBucketsList())) return false;
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
    hash = (37 * hash) + NAME_FIELD_NUMBER;
    hash = (53 * hash) + getName().hashCode();
    if (hasSelector()) {
      hash = (37 * hash) + SELECTOR_FIELD_NUMBER;
      hash = (53 * hash) + getSelector().hashCode();
    }
    if (getHistogramBucketsCount() > 0) {
      hash = (37 * hash) + HISTOGRAM_BUCKETS_FIELD_NUMBER;
      hash = (53 * hash) + getHistogramBucketsList().hashCode();
    }
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.aperture.policy.language.v1.FluxMeter parseFrom(
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
  public static Builder newBuilder(com.aperture.policy.language.v1.FluxMeter prototype) {
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
   * FluxMeter gathers metrics for the traffic that matches its selector.
   * </pre>
   *
   * Protobuf type {@code aperture.policy.language.v1.FluxMeter}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.policy.language.v1.FluxMeter)
      com.aperture.policy.language.v1.FluxMeterOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_FluxMeter_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_FluxMeter_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.aperture.policy.language.v1.FluxMeter.class, com.aperture.policy.language.v1.FluxMeter.Builder.class);
    }

    // Construct using com.aperture.policy.language.v1.FluxMeter.newBuilder()
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
      name_ = "";

      if (selectorBuilder_ == null) {
        selector_ = null;
      } else {
        selector_ = null;
        selectorBuilder_ = null;
      }
      histogramBuckets_ = emptyDoubleList();
      bitField0_ = (bitField0_ & ~0x00000001);
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_FluxMeter_descriptor;
    }

    @java.lang.Override
    public com.aperture.policy.language.v1.FluxMeter getDefaultInstanceForType() {
      return com.aperture.policy.language.v1.FluxMeter.getDefaultInstance();
    }

    @java.lang.Override
    public com.aperture.policy.language.v1.FluxMeter build() {
      com.aperture.policy.language.v1.FluxMeter result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.aperture.policy.language.v1.FluxMeter buildPartial() {
      com.aperture.policy.language.v1.FluxMeter result = new com.aperture.policy.language.v1.FluxMeter(this);
      int from_bitField0_ = bitField0_;
      result.name_ = name_;
      if (selectorBuilder_ == null) {
        result.selector_ = selector_;
      } else {
        result.selector_ = selectorBuilder_.build();
      }
      if (((bitField0_ & 0x00000001) != 0)) {
        histogramBuckets_.makeImmutable();
        bitField0_ = (bitField0_ & ~0x00000001);
      }
      result.histogramBuckets_ = histogramBuckets_;
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
      if (other instanceof com.aperture.policy.language.v1.FluxMeter) {
        return mergeFrom((com.aperture.policy.language.v1.FluxMeter)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.aperture.policy.language.v1.FluxMeter other) {
      if (other == com.aperture.policy.language.v1.FluxMeter.getDefaultInstance()) return this;
      if (!other.getName().isEmpty()) {
        name_ = other.name_;
        onChanged();
      }
      if (other.hasSelector()) {
        mergeSelector(other.getSelector());
      }
      if (!other.histogramBuckets_.isEmpty()) {
        if (histogramBuckets_.isEmpty()) {
          histogramBuckets_ = other.histogramBuckets_;
          bitField0_ = (bitField0_ & ~0x00000001);
        } else {
          ensureHistogramBucketsIsMutable();
          histogramBuckets_.addAll(other.histogramBuckets_);
        }
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
      com.aperture.policy.language.v1.FluxMeter parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.aperture.policy.language.v1.FluxMeter) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }
    private int bitField0_;

    private java.lang.Object name_ = "";
    /**
     * <pre>
     * Name of the flux meter.
     * </pre>
     *
     * <code>string name = 1 [json_name = "name"];</code>
     * @return The name.
     */
    public java.lang.String getName() {
      java.lang.Object ref = name_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        name_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <pre>
     * Name of the flux meter.
     * </pre>
     *
     * <code>string name = 1 [json_name = "name"];</code>
     * @return The bytes for name.
     */
    public com.google.protobuf.ByteString
        getNameBytes() {
      java.lang.Object ref = name_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b =
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        name_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <pre>
     * Name of the flux meter.
     * </pre>
     *
     * <code>string name = 1 [json_name = "name"];</code>
     * @param value The name to set.
     * @return This builder for chaining.
     */
    public Builder setName(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }

      name_ = value;
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Name of the flux meter.
     * </pre>
     *
     * <code>string name = 1 [json_name = "name"];</code>
     * @return This builder for chaining.
     */
    public Builder clearName() {

      name_ = getDefaultInstance().getName();
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Name of the flux meter.
     * </pre>
     *
     * <code>string name = 1 [json_name = "name"];</code>
     * @param value The bytes for name to set.
     * @return This builder for chaining.
     */
    public Builder setNameBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);

      name_ = value;
      onChanged();
      return this;
    }

    private com.aperture.policy.language.v1.Selector selector_;
    private com.google.protobuf.SingleFieldBuilderV3<
        com.aperture.policy.language.v1.Selector, com.aperture.policy.language.v1.Selector.Builder, com.aperture.policy.language.v1.SelectorOrBuilder> selectorBuilder_;
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     * @return Whether the selector field is set.
     */
    public boolean hasSelector() {
      return selectorBuilder_ != null || selector_ != null;
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     * @return The selector.
     */
    public com.aperture.policy.language.v1.Selector getSelector() {
      if (selectorBuilder_ == null) {
        return selector_ == null ? com.aperture.policy.language.v1.Selector.getDefaultInstance() : selector_;
      } else {
        return selectorBuilder_.getMessage();
      }
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     */
    public Builder setSelector(com.aperture.policy.language.v1.Selector value) {
      if (selectorBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        selector_ = value;
        onChanged();
      } else {
        selectorBuilder_.setMessage(value);
      }

      return this;
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     */
    public Builder setSelector(
        com.aperture.policy.language.v1.Selector.Builder builderForValue) {
      if (selectorBuilder_ == null) {
        selector_ = builderForValue.build();
        onChanged();
      } else {
        selectorBuilder_.setMessage(builderForValue.build());
      }

      return this;
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     */
    public Builder mergeSelector(com.aperture.policy.language.v1.Selector value) {
      if (selectorBuilder_ == null) {
        if (selector_ != null) {
          selector_ =
            com.aperture.policy.language.v1.Selector.newBuilder(selector_).mergeFrom(value).buildPartial();
        } else {
          selector_ = value;
        }
        onChanged();
      } else {
        selectorBuilder_.mergeFrom(value);
      }

      return this;
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     */
    public Builder clearSelector() {
      if (selectorBuilder_ == null) {
        selector_ = null;
        onChanged();
      } else {
        selector_ = null;
        selectorBuilder_ = null;
      }

      return this;
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     */
    public com.aperture.policy.language.v1.Selector.Builder getSelectorBuilder() {

      onChanged();
      return getSelectorFieldBuilder().getBuilder();
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     */
    public com.aperture.policy.language.v1.SelectorOrBuilder getSelectorOrBuilder() {
      if (selectorBuilder_ != null) {
        return selectorBuilder_.getMessageOrBuilder();
      } else {
        return selector_ == null ?
            com.aperture.policy.language.v1.Selector.getDefaultInstance() : selector_;
      }
    }
    /**
     * <pre>
     * Policies are only applied to flows that are matched based on the fields in the selector.
     * </pre>
     *
     * <code>.aperture.policy.language.v1.Selector selector = 2 [json_name = "selector"];</code>
     */
    private com.google.protobuf.SingleFieldBuilderV3<
        com.aperture.policy.language.v1.Selector, com.aperture.policy.language.v1.Selector.Builder, com.aperture.policy.language.v1.SelectorOrBuilder>
        getSelectorFieldBuilder() {
      if (selectorBuilder_ == null) {
        selectorBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
            com.aperture.policy.language.v1.Selector, com.aperture.policy.language.v1.Selector.Builder, com.aperture.policy.language.v1.SelectorOrBuilder>(
                getSelector(),
                getParentForChildren(),
                isClean());
        selector_ = null;
      }
      return selectorBuilder_;
    }

    private com.google.protobuf.Internal.DoubleList histogramBuckets_ = emptyDoubleList();
    private void ensureHistogramBucketsIsMutable() {
      if (!((bitField0_ & 0x00000001) != 0)) {
        histogramBuckets_ = mutableCopy(histogramBuckets_);
        bitField0_ |= 0x00000001;
       }
    }
    /**
     * <pre>
     * Latency histogram buckets (in ms) for this FluxMeter.
     * </pre>
     *
     * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
     * @return A list containing the histogramBuckets.
     */
    public java.util.List<java.lang.Double>
        getHistogramBucketsList() {
      return ((bitField0_ & 0x00000001) != 0) ?
               java.util.Collections.unmodifiableList(histogramBuckets_) : histogramBuckets_;
    }
    /**
     * <pre>
     * Latency histogram buckets (in ms) for this FluxMeter.
     * </pre>
     *
     * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
     * @return The count of histogramBuckets.
     */
    public int getHistogramBucketsCount() {
      return histogramBuckets_.size();
    }
    /**
     * <pre>
     * Latency histogram buckets (in ms) for this FluxMeter.
     * </pre>
     *
     * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
     * @param index The index of the element to return.
     * @return The histogramBuckets at the given index.
     */
    public double getHistogramBuckets(int index) {
      return histogramBuckets_.getDouble(index);
    }
    /**
     * <pre>
     * Latency histogram buckets (in ms) for this FluxMeter.
     * </pre>
     *
     * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
     * @param index The index to set the value at.
     * @param value The histogramBuckets to set.
     * @return This builder for chaining.
     */
    public Builder setHistogramBuckets(
        int index, double value) {
      ensureHistogramBucketsIsMutable();
      histogramBuckets_.setDouble(index, value);
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Latency histogram buckets (in ms) for this FluxMeter.
     * </pre>
     *
     * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
     * @param value The histogramBuckets to add.
     * @return This builder for chaining.
     */
    public Builder addHistogramBuckets(double value) {
      ensureHistogramBucketsIsMutable();
      histogramBuckets_.addDouble(value);
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Latency histogram buckets (in ms) for this FluxMeter.
     * </pre>
     *
     * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
     * @param values The histogramBuckets to add.
     * @return This builder for chaining.
     */
    public Builder addAllHistogramBuckets(
        java.lang.Iterable<? extends java.lang.Double> values) {
      ensureHistogramBucketsIsMutable();
      com.google.protobuf.AbstractMessageLite.Builder.addAll(
          values, histogramBuckets_);
      onChanged();
      return this;
    }
    /**
     * <pre>
     * Latency histogram buckets (in ms) for this FluxMeter.
     * </pre>
     *
     * <code>repeated double histogram_buckets = 3 [json_name = "histogramBuckets", (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { ... }</code>
     * @return This builder for chaining.
     */
    public Builder clearHistogramBuckets() {
      histogramBuckets_ = emptyDoubleList();
      bitField0_ = (bitField0_ & ~0x00000001);
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


    // @@protoc_insertion_point(builder_scope:aperture.policy.language.v1.FluxMeter)
  }

  // @@protoc_insertion_point(class_scope:aperture.policy.language.v1.FluxMeter)
  private static final com.aperture.policy.language.v1.FluxMeter DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.aperture.policy.language.v1.FluxMeter();
  }

  public static com.aperture.policy.language.v1.FluxMeter getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<FluxMeter>
      PARSER = new com.google.protobuf.AbstractParser<FluxMeter>() {
    @java.lang.Override
    public FluxMeter parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new FluxMeter(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<FluxMeter> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<FluxMeter> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.aperture.policy.language.v1.FluxMeter getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}
