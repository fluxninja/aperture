// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

/**
 * Protobuf type {@code aperture.flowcontrol.check.v1.InflightRequestRef}
 */
public final class InflightRequestRef extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.flowcontrol.check.v1.InflightRequestRef)
    InflightRequestRefOrBuilder {
private static final long serialVersionUID = 0L;
  // Use InflightRequestRef.newBuilder() to construct.
  private InflightRequestRef(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private InflightRequestRef() {
    policyName_ = "";
    policyHash_ = "";
    componentId_ = "";
    label_ = "";
    requestId_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new InflightRequestRef();
  }

  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_InflightRequestRef_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_InflightRequestRef_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef.class, com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef.Builder.class);
  }

  public static final int POLICY_NAME_FIELD_NUMBER = 1;
  @SuppressWarnings("serial")
  private volatile java.lang.Object policyName_ = "";
  /**
   * <code>string policy_name = 1 [json_name = "policyName"];</code>
   * @return The policyName.
   */
  @java.lang.Override
  public java.lang.String getPolicyName() {
    java.lang.Object ref = policyName_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      policyName_ = s;
      return s;
    }
  }
  /**
   * <code>string policy_name = 1 [json_name = "policyName"];</code>
   * @return The bytes for policyName.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getPolicyNameBytes() {
    java.lang.Object ref = policyName_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      policyName_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int POLICY_HASH_FIELD_NUMBER = 2;
  @SuppressWarnings("serial")
  private volatile java.lang.Object policyHash_ = "";
  /**
   * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
   * @return The policyHash.
   */
  @java.lang.Override
  public java.lang.String getPolicyHash() {
    java.lang.Object ref = policyHash_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      policyHash_ = s;
      return s;
    }
  }
  /**
   * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
   * @return The bytes for policyHash.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getPolicyHashBytes() {
    java.lang.Object ref = policyHash_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      policyHash_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int COMPONENT_ID_FIELD_NUMBER = 3;
  @SuppressWarnings("serial")
  private volatile java.lang.Object componentId_ = "";
  /**
   * <code>string component_id = 3 [json_name = "componentId"];</code>
   * @return The componentId.
   */
  @java.lang.Override
  public java.lang.String getComponentId() {
    java.lang.Object ref = componentId_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      componentId_ = s;
      return s;
    }
  }
  /**
   * <code>string component_id = 3 [json_name = "componentId"];</code>
   * @return The bytes for componentId.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getComponentIdBytes() {
    java.lang.Object ref = componentId_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      componentId_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int LABEL_FIELD_NUMBER = 4;
  @SuppressWarnings("serial")
  private volatile java.lang.Object label_ = "";
  /**
   * <code>string label = 4 [json_name = "label"];</code>
   * @return The label.
   */
  @java.lang.Override
  public java.lang.String getLabel() {
    java.lang.Object ref = label_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      label_ = s;
      return s;
    }
  }
  /**
   * <code>string label = 4 [json_name = "label"];</code>
   * @return The bytes for label.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getLabelBytes() {
    java.lang.Object ref = label_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      label_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int REQUEST_ID_FIELD_NUMBER = 5;
  @SuppressWarnings("serial")
  private volatile java.lang.Object requestId_ = "";
  /**
   * <code>string request_id = 5 [json_name = "requestId"];</code>
   * @return The requestId.
   */
  @java.lang.Override
  public java.lang.String getRequestId() {
    java.lang.Object ref = requestId_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      requestId_ = s;
      return s;
    }
  }
  /**
   * <code>string request_id = 5 [json_name = "requestId"];</code>
   * @return The bytes for requestId.
   */
  @java.lang.Override
  public com.google.protobuf.ByteString
      getRequestIdBytes() {
    java.lang.Object ref = requestId_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      requestId_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int TOKENS_FIELD_NUMBER = 6;
  private double tokens_ = 0D;
  /**
   * <code>double tokens = 6 [json_name = "tokens"];</code>
   * @return The tokens.
   */
  @java.lang.Override
  public double getTokens() {
    return tokens_;
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
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(policyName_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 1, policyName_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(policyHash_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 2, policyHash_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(componentId_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 3, componentId_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(label_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 4, label_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(requestId_)) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 5, requestId_);
    }
    if (java.lang.Double.doubleToRawLongBits(tokens_) != 0) {
      output.writeDouble(6, tokens_);
    }
    getUnknownFields().writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(policyName_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, policyName_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(policyHash_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(2, policyHash_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(componentId_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(3, componentId_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(label_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(4, label_);
    }
    if (!com.google.protobuf.GeneratedMessageV3.isStringEmpty(requestId_)) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(5, requestId_);
    }
    if (java.lang.Double.doubleToRawLongBits(tokens_) != 0) {
      size += com.google.protobuf.CodedOutputStream
        .computeDoubleSize(6, tokens_);
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
    if (!(obj instanceof com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef other = (com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef) obj;

    if (!getPolicyName()
        .equals(other.getPolicyName())) return false;
    if (!getPolicyHash()
        .equals(other.getPolicyHash())) return false;
    if (!getComponentId()
        .equals(other.getComponentId())) return false;
    if (!getLabel()
        .equals(other.getLabel())) return false;
    if (!getRequestId()
        .equals(other.getRequestId())) return false;
    if (java.lang.Double.doubleToLongBits(getTokens())
        != java.lang.Double.doubleToLongBits(
            other.getTokens())) return false;
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
    hash = (37 * hash) + POLICY_NAME_FIELD_NUMBER;
    hash = (53 * hash) + getPolicyName().hashCode();
    hash = (37 * hash) + POLICY_HASH_FIELD_NUMBER;
    hash = (53 * hash) + getPolicyHash().hashCode();
    hash = (37 * hash) + COMPONENT_ID_FIELD_NUMBER;
    hash = (53 * hash) + getComponentId().hashCode();
    hash = (37 * hash) + LABEL_FIELD_NUMBER;
    hash = (53 * hash) + getLabel().hashCode();
    hash = (37 * hash) + REQUEST_ID_FIELD_NUMBER;
    hash = (53 * hash) + getRequestId().hashCode();
    hash = (37 * hash) + TOKENS_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        java.lang.Double.doubleToLongBits(getTokens()));
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef prototype) {
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
   * Protobuf type {@code aperture.flowcontrol.check.v1.InflightRequestRef}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.flowcontrol.check.v1.InflightRequestRef)
      com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRefOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_InflightRequestRef_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_InflightRequestRef_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef.class, com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef.newBuilder()
    private Builder() {

    }

    private Builder(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      super(parent);

    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      bitField0_ = 0;
      policyName_ = "";
      policyHash_ = "";
      componentId_ = "";
      label_ = "";
      requestId_ = "";
      tokens_ = 0D;
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.internal_static_aperture_flowcontrol_check_v1_InflightRequestRef_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef build() {
      com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef buildPartial() {
      com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef result = new com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef(this);
      if (bitField0_ != 0) { buildPartial0(result); }
      onBuilt();
      return result;
    }

    private void buildPartial0(com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef result) {
      int from_bitField0_ = bitField0_;
      if (((from_bitField0_ & 0x00000001) != 0)) {
        result.policyName_ = policyName_;
      }
      if (((from_bitField0_ & 0x00000002) != 0)) {
        result.policyHash_ = policyHash_;
      }
      if (((from_bitField0_ & 0x00000004) != 0)) {
        result.componentId_ = componentId_;
      }
      if (((from_bitField0_ & 0x00000008) != 0)) {
        result.label_ = label_;
      }
      if (((from_bitField0_ & 0x00000010) != 0)) {
        result.requestId_ = requestId_;
      }
      if (((from_bitField0_ & 0x00000020) != 0)) {
        result.tokens_ = tokens_;
      }
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
      if (other instanceof com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef) {
        return mergeFrom((com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef other) {
      if (other == com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef.getDefaultInstance()) return this;
      if (!other.getPolicyName().isEmpty()) {
        policyName_ = other.policyName_;
        bitField0_ |= 0x00000001;
        onChanged();
      }
      if (!other.getPolicyHash().isEmpty()) {
        policyHash_ = other.policyHash_;
        bitField0_ |= 0x00000002;
        onChanged();
      }
      if (!other.getComponentId().isEmpty()) {
        componentId_ = other.componentId_;
        bitField0_ |= 0x00000004;
        onChanged();
      }
      if (!other.getLabel().isEmpty()) {
        label_ = other.label_;
        bitField0_ |= 0x00000008;
        onChanged();
      }
      if (!other.getRequestId().isEmpty()) {
        requestId_ = other.requestId_;
        bitField0_ |= 0x00000010;
        onChanged();
      }
      if (other.getTokens() != 0D) {
        setTokens(other.getTokens());
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
              policyName_ = input.readStringRequireUtf8();
              bitField0_ |= 0x00000001;
              break;
            } // case 10
            case 18: {
              policyHash_ = input.readStringRequireUtf8();
              bitField0_ |= 0x00000002;
              break;
            } // case 18
            case 26: {
              componentId_ = input.readStringRequireUtf8();
              bitField0_ |= 0x00000004;
              break;
            } // case 26
            case 34: {
              label_ = input.readStringRequireUtf8();
              bitField0_ |= 0x00000008;
              break;
            } // case 34
            case 42: {
              requestId_ = input.readStringRequireUtf8();
              bitField0_ |= 0x00000010;
              break;
            } // case 42
            case 49: {
              tokens_ = input.readDouble();
              bitField0_ |= 0x00000020;
              break;
            } // case 49
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
    private int bitField0_;

    private java.lang.Object policyName_ = "";
    /**
     * <code>string policy_name = 1 [json_name = "policyName"];</code>
     * @return The policyName.
     */
    public java.lang.String getPolicyName() {
      java.lang.Object ref = policyName_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        policyName_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string policy_name = 1 [json_name = "policyName"];</code>
     * @return The bytes for policyName.
     */
    public com.google.protobuf.ByteString
        getPolicyNameBytes() {
      java.lang.Object ref = policyName_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        policyName_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string policy_name = 1 [json_name = "policyName"];</code>
     * @param value The policyName to set.
     * @return This builder for chaining.
     */
    public Builder setPolicyName(
        java.lang.String value) {
      if (value == null) { throw new NullPointerException(); }
      policyName_ = value;
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>string policy_name = 1 [json_name = "policyName"];</code>
     * @return This builder for chaining.
     */
    public Builder clearPolicyName() {
      policyName_ = getDefaultInstance().getPolicyName();
      bitField0_ = (bitField0_ & ~0x00000001);
      onChanged();
      return this;
    }
    /**
     * <code>string policy_name = 1 [json_name = "policyName"];</code>
     * @param value The bytes for policyName to set.
     * @return This builder for chaining.
     */
    public Builder setPolicyNameBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) { throw new NullPointerException(); }
      checkByteStringIsUtf8(value);
      policyName_ = value;
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }

    private java.lang.Object policyHash_ = "";
    /**
     * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
     * @return The policyHash.
     */
    public java.lang.String getPolicyHash() {
      java.lang.Object ref = policyHash_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        policyHash_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
     * @return The bytes for policyHash.
     */
    public com.google.protobuf.ByteString
        getPolicyHashBytes() {
      java.lang.Object ref = policyHash_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        policyHash_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
     * @param value The policyHash to set.
     * @return This builder for chaining.
     */
    public Builder setPolicyHash(
        java.lang.String value) {
      if (value == null) { throw new NullPointerException(); }
      policyHash_ = value;
      bitField0_ |= 0x00000002;
      onChanged();
      return this;
    }
    /**
     * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
     * @return This builder for chaining.
     */
    public Builder clearPolicyHash() {
      policyHash_ = getDefaultInstance().getPolicyHash();
      bitField0_ = (bitField0_ & ~0x00000002);
      onChanged();
      return this;
    }
    /**
     * <code>string policy_hash = 2 [json_name = "policyHash"];</code>
     * @param value The bytes for policyHash to set.
     * @return This builder for chaining.
     */
    public Builder setPolicyHashBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) { throw new NullPointerException(); }
      checkByteStringIsUtf8(value);
      policyHash_ = value;
      bitField0_ |= 0x00000002;
      onChanged();
      return this;
    }

    private java.lang.Object componentId_ = "";
    /**
     * <code>string component_id = 3 [json_name = "componentId"];</code>
     * @return The componentId.
     */
    public java.lang.String getComponentId() {
      java.lang.Object ref = componentId_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        componentId_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string component_id = 3 [json_name = "componentId"];</code>
     * @return The bytes for componentId.
     */
    public com.google.protobuf.ByteString
        getComponentIdBytes() {
      java.lang.Object ref = componentId_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        componentId_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string component_id = 3 [json_name = "componentId"];</code>
     * @param value The componentId to set.
     * @return This builder for chaining.
     */
    public Builder setComponentId(
        java.lang.String value) {
      if (value == null) { throw new NullPointerException(); }
      componentId_ = value;
      bitField0_ |= 0x00000004;
      onChanged();
      return this;
    }
    /**
     * <code>string component_id = 3 [json_name = "componentId"];</code>
     * @return This builder for chaining.
     */
    public Builder clearComponentId() {
      componentId_ = getDefaultInstance().getComponentId();
      bitField0_ = (bitField0_ & ~0x00000004);
      onChanged();
      return this;
    }
    /**
     * <code>string component_id = 3 [json_name = "componentId"];</code>
     * @param value The bytes for componentId to set.
     * @return This builder for chaining.
     */
    public Builder setComponentIdBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) { throw new NullPointerException(); }
      checkByteStringIsUtf8(value);
      componentId_ = value;
      bitField0_ |= 0x00000004;
      onChanged();
      return this;
    }

    private java.lang.Object label_ = "";
    /**
     * <code>string label = 4 [json_name = "label"];</code>
     * @return The label.
     */
    public java.lang.String getLabel() {
      java.lang.Object ref = label_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        label_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string label = 4 [json_name = "label"];</code>
     * @return The bytes for label.
     */
    public com.google.protobuf.ByteString
        getLabelBytes() {
      java.lang.Object ref = label_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        label_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string label = 4 [json_name = "label"];</code>
     * @param value The label to set.
     * @return This builder for chaining.
     */
    public Builder setLabel(
        java.lang.String value) {
      if (value == null) { throw new NullPointerException(); }
      label_ = value;
      bitField0_ |= 0x00000008;
      onChanged();
      return this;
    }
    /**
     * <code>string label = 4 [json_name = "label"];</code>
     * @return This builder for chaining.
     */
    public Builder clearLabel() {
      label_ = getDefaultInstance().getLabel();
      bitField0_ = (bitField0_ & ~0x00000008);
      onChanged();
      return this;
    }
    /**
     * <code>string label = 4 [json_name = "label"];</code>
     * @param value The bytes for label to set.
     * @return This builder for chaining.
     */
    public Builder setLabelBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) { throw new NullPointerException(); }
      checkByteStringIsUtf8(value);
      label_ = value;
      bitField0_ |= 0x00000008;
      onChanged();
      return this;
    }

    private java.lang.Object requestId_ = "";
    /**
     * <code>string request_id = 5 [json_name = "requestId"];</code>
     * @return The requestId.
     */
    public java.lang.String getRequestId() {
      java.lang.Object ref = requestId_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        requestId_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string request_id = 5 [json_name = "requestId"];</code>
     * @return The bytes for requestId.
     */
    public com.google.protobuf.ByteString
        getRequestIdBytes() {
      java.lang.Object ref = requestId_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        requestId_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string request_id = 5 [json_name = "requestId"];</code>
     * @param value The requestId to set.
     * @return This builder for chaining.
     */
    public Builder setRequestId(
        java.lang.String value) {
      if (value == null) { throw new NullPointerException(); }
      requestId_ = value;
      bitField0_ |= 0x00000010;
      onChanged();
      return this;
    }
    /**
     * <code>string request_id = 5 [json_name = "requestId"];</code>
     * @return This builder for chaining.
     */
    public Builder clearRequestId() {
      requestId_ = getDefaultInstance().getRequestId();
      bitField0_ = (bitField0_ & ~0x00000010);
      onChanged();
      return this;
    }
    /**
     * <code>string request_id = 5 [json_name = "requestId"];</code>
     * @param value The bytes for requestId to set.
     * @return This builder for chaining.
     */
    public Builder setRequestIdBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) { throw new NullPointerException(); }
      checkByteStringIsUtf8(value);
      requestId_ = value;
      bitField0_ |= 0x00000010;
      onChanged();
      return this;
    }

    private double tokens_ ;
    /**
     * <code>double tokens = 6 [json_name = "tokens"];</code>
     * @return The tokens.
     */
    @java.lang.Override
    public double getTokens() {
      return tokens_;
    }
    /**
     * <code>double tokens = 6 [json_name = "tokens"];</code>
     * @param value The tokens to set.
     * @return This builder for chaining.
     */
    public Builder setTokens(double value) {

      tokens_ = value;
      bitField0_ |= 0x00000020;
      onChanged();
      return this;
    }
    /**
     * <code>double tokens = 6 [json_name = "tokens"];</code>
     * @return This builder for chaining.
     */
    public Builder clearTokens() {
      bitField0_ = (bitField0_ & ~0x00000020);
      tokens_ = 0D;
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


    // @@protoc_insertion_point(builder_scope:aperture.flowcontrol.check.v1.InflightRequestRef)
  }

  // @@protoc_insertion_point(class_scope:aperture.flowcontrol.check.v1.InflightRequestRef)
  private static final com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef();
  }

  public static com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<InflightRequestRef>
      PARSER = new com.google.protobuf.AbstractParser<InflightRequestRef>() {
    @java.lang.Override
    public InflightRequestRef parsePartialFrom(
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

  public static com.google.protobuf.Parser<InflightRequestRef> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<InflightRequestRef> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.flowcontrol.check.v1.InflightRequestRef getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

