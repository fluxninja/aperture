// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/policy.proto

package com.fluxninja.generated.aperture.policy.language.v1;

/**
 * Protobuf type {@code aperture.policy.language.v1.Policies}
 */
public final class Policies extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.policy.language.v1.Policies)
    PoliciesOrBuilder {
private static final long serialVersionUID = 0L;
  // Use Policies.newBuilder() to construct.
  private Policies(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private Policies() {
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new Policies();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private Policies(
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
              policies_ = com.google.protobuf.MapField.newMapField(
                  PoliciesDefaultEntryHolder.defaultEntry);
              mutable_bitField0_ |= 0x00000001;
            }
            com.google.protobuf.MapEntry<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy>
            policies__ = input.readMessage(
                PoliciesDefaultEntryHolder.defaultEntry.getParserForType(), extensionRegistry);
            policies_.getMutableMap().put(
                policies__.getKey(), policies__.getValue());
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
    return com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_Policies_descriptor;
  }

  @SuppressWarnings({"rawtypes"})
  @java.lang.Override
  protected com.google.protobuf.MapField internalGetMapField(
      int number) {
    switch (number) {
      case 1:
        return internalGetPolicies();
      default:
        throw new RuntimeException(
            "Invalid map field number: " + number);
    }
  }
  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_Policies_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.policy.language.v1.Policies.class, com.fluxninja.generated.aperture.policy.language.v1.Policies.Builder.class);
  }

  public static final int POLICIES_FIELD_NUMBER = 1;
  private static final class PoliciesDefaultEntryHolder {
    static final com.google.protobuf.MapEntry<
        java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> defaultEntry =
            com.google.protobuf.MapEntry
            .<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy>newDefaultInstance(
                com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_Policies_PoliciesEntry_descriptor, 
                com.google.protobuf.WireFormat.FieldType.STRING,
                "",
                com.google.protobuf.WireFormat.FieldType.MESSAGE,
                com.fluxninja.generated.aperture.policy.language.v1.Policy.getDefaultInstance());
  }
  private com.google.protobuf.MapField<
      java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> policies_;
  private com.google.protobuf.MapField<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy>
  internalGetPolicies() {
    if (policies_ == null) {
      return com.google.protobuf.MapField.emptyMapField(
          PoliciesDefaultEntryHolder.defaultEntry);
    }
    return policies_;
  }

  public int getPoliciesCount() {
    return internalGetPolicies().getMap().size();
  }
  /**
   * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
   */

  @java.lang.Override
  public boolean containsPolicies(
      java.lang.String key) {
    if (key == null) { throw new NullPointerException("map key"); }
    return internalGetPolicies().getMap().containsKey(key);
  }
  /**
   * Use {@link #getPoliciesMap()} instead.
   */
  @java.lang.Override
  @java.lang.Deprecated
  public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> getPolicies() {
    return getPoliciesMap();
  }
  /**
   * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
   */
  @java.lang.Override

  public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> getPoliciesMap() {
    return internalGetPolicies().getMap();
  }
  /**
   * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
   */
  @java.lang.Override

  public com.fluxninja.generated.aperture.policy.language.v1.Policy getPoliciesOrDefault(
      java.lang.String key,
      com.fluxninja.generated.aperture.policy.language.v1.Policy defaultValue) {
    if (key == null) { throw new NullPointerException("map key"); }
    java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> map =
        internalGetPolicies().getMap();
    return map.containsKey(key) ? map.get(key) : defaultValue;
  }
  /**
   * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
   */
  @java.lang.Override

  public com.fluxninja.generated.aperture.policy.language.v1.Policy getPoliciesOrThrow(
      java.lang.String key) {
    if (key == null) { throw new NullPointerException("map key"); }
    java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> map =
        internalGetPolicies().getMap();
    if (!map.containsKey(key)) {
      throw new java.lang.IllegalArgumentException();
    }
    return map.get(key);
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
    com.google.protobuf.GeneratedMessageV3
      .serializeStringMapTo(
        output,
        internalGetPolicies(),
        PoliciesDefaultEntryHolder.defaultEntry,
        1);
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    for (java.util.Map.Entry<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> entry
         : internalGetPolicies().getMap().entrySet()) {
      com.google.protobuf.MapEntry<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy>
      policies__ = PoliciesDefaultEntryHolder.defaultEntry.newBuilderForType()
          .setKey(entry.getKey())
          .setValue(entry.getValue())
          .build();
      size += com.google.protobuf.CodedOutputStream
          .computeMessageSize(1, policies__);
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
    if (!(obj instanceof com.fluxninja.generated.aperture.policy.language.v1.Policies)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.policy.language.v1.Policies other = (com.fluxninja.generated.aperture.policy.language.v1.Policies) obj;

    if (!internalGetPolicies().equals(
        other.internalGetPolicies())) return false;
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
    if (!internalGetPolicies().getMap().isEmpty()) {
      hash = (37 * hash) + POLICIES_FIELD_NUMBER;
      hash = (53 * hash) + internalGetPolicies().hashCode();
    }
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.policy.language.v1.Policies parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.policy.language.v1.Policies prototype) {
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
   * Protobuf type {@code aperture.policy.language.v1.Policies}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.policy.language.v1.Policies)
      com.fluxninja.generated.aperture.policy.language.v1.PoliciesOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_Policies_descriptor;
    }

    @SuppressWarnings({"rawtypes"})
    protected com.google.protobuf.MapField internalGetMapField(
        int number) {
      switch (number) {
        case 1:
          return internalGetPolicies();
        default:
          throw new RuntimeException(
              "Invalid map field number: " + number);
      }
    }
    @SuppressWarnings({"rawtypes"})
    protected com.google.protobuf.MapField internalGetMutableMapField(
        int number) {
      switch (number) {
        case 1:
          return internalGetMutablePolicies();
        default:
          throw new RuntimeException(
              "Invalid map field number: " + number);
      }
    }
    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_Policies_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.policy.language.v1.Policies.class, com.fluxninja.generated.aperture.policy.language.v1.Policies.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.policy.language.v1.Policies.newBuilder()
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
      internalGetMutablePolicies().clear();
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.internal_static_aperture_policy_language_v1_Policies_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.language.v1.Policies getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.policy.language.v1.Policies.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.language.v1.Policies build() {
      com.fluxninja.generated.aperture.policy.language.v1.Policies result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.policy.language.v1.Policies buildPartial() {
      com.fluxninja.generated.aperture.policy.language.v1.Policies result = new com.fluxninja.generated.aperture.policy.language.v1.Policies(this);
      int from_bitField0_ = bitField0_;
      result.policies_ = internalGetPolicies();
      result.policies_.makeImmutable();
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
      if (other instanceof com.fluxninja.generated.aperture.policy.language.v1.Policies) {
        return mergeFrom((com.fluxninja.generated.aperture.policy.language.v1.Policies)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.policy.language.v1.Policies other) {
      if (other == com.fluxninja.generated.aperture.policy.language.v1.Policies.getDefaultInstance()) return this;
      internalGetMutablePolicies().mergeFrom(
          other.internalGetPolicies());
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
      com.fluxninja.generated.aperture.policy.language.v1.Policies parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.policy.language.v1.Policies) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }
    private int bitField0_;

    private com.google.protobuf.MapField<
        java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> policies_;
    private com.google.protobuf.MapField<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy>
    internalGetPolicies() {
      if (policies_ == null) {
        return com.google.protobuf.MapField.emptyMapField(
            PoliciesDefaultEntryHolder.defaultEntry);
      }
      return policies_;
    }
    private com.google.protobuf.MapField<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy>
    internalGetMutablePolicies() {
      onChanged();;
      if (policies_ == null) {
        policies_ = com.google.protobuf.MapField.newMapField(
            PoliciesDefaultEntryHolder.defaultEntry);
      }
      if (!policies_.isMutable()) {
        policies_ = policies_.copy();
      }
      return policies_;
    }

    public int getPoliciesCount() {
      return internalGetPolicies().getMap().size();
    }
    /**
     * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
     */

    @java.lang.Override
    public boolean containsPolicies(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      return internalGetPolicies().getMap().containsKey(key);
    }
    /**
     * Use {@link #getPoliciesMap()} instead.
     */
    @java.lang.Override
    @java.lang.Deprecated
    public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> getPolicies() {
      return getPoliciesMap();
    }
    /**
     * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
     */
    @java.lang.Override

    public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> getPoliciesMap() {
      return internalGetPolicies().getMap();
    }
    /**
     * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
     */
    @java.lang.Override

    public com.fluxninja.generated.aperture.policy.language.v1.Policy getPoliciesOrDefault(
        java.lang.String key,
        com.fluxninja.generated.aperture.policy.language.v1.Policy defaultValue) {
      if (key == null) { throw new NullPointerException("map key"); }
      java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> map =
          internalGetPolicies().getMap();
      return map.containsKey(key) ? map.get(key) : defaultValue;
    }
    /**
     * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
     */
    @java.lang.Override

    public com.fluxninja.generated.aperture.policy.language.v1.Policy getPoliciesOrThrow(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> map =
          internalGetPolicies().getMap();
      if (!map.containsKey(key)) {
        throw new java.lang.IllegalArgumentException();
      }
      return map.get(key);
    }

    public Builder clearPolicies() {
      internalGetMutablePolicies().getMutableMap()
          .clear();
      return this;
    }
    /**
     * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
     */

    public Builder removePolicies(
        java.lang.String key) {
      if (key == null) { throw new NullPointerException("map key"); }
      internalGetMutablePolicies().getMutableMap()
          .remove(key);
      return this;
    }
    /**
     * Use alternate mutation accessors instead.
     */
    @java.lang.Deprecated
    public java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy>
    getMutablePolicies() {
      return internalGetMutablePolicies().getMutableMap();
    }
    /**
     * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
     */
    public Builder putPolicies(
        java.lang.String key,
        com.fluxninja.generated.aperture.policy.language.v1.Policy value) {
      if (key == null) { throw new NullPointerException("map key"); }
      if (value == null) {
  throw new NullPointerException("map value");
}

      internalGetMutablePolicies().getMutableMap()
          .put(key, value);
      return this;
    }
    /**
     * <code>map&lt;string, .aperture.policy.language.v1.Policy&gt; policies = 1 [json_name = "policies"];</code>
     */

    public Builder putAllPolicies(
        java.util.Map<java.lang.String, com.fluxninja.generated.aperture.policy.language.v1.Policy> values) {
      internalGetMutablePolicies().getMutableMap()
          .putAll(values);
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


    // @@protoc_insertion_point(builder_scope:aperture.policy.language.v1.Policies)
  }

  // @@protoc_insertion_point(class_scope:aperture.policy.language.v1.Policies)
  private static final com.fluxninja.generated.aperture.policy.language.v1.Policies DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.policy.language.v1.Policies();
  }

  public static com.fluxninja.generated.aperture.policy.language.v1.Policies getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<Policies>
      PARSER = new com.google.protobuf.AbstractParser<Policies>() {
    @java.lang.Override
    public Policies parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new Policies(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<Policies> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<Policies> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.policy.language.v1.Policies getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

