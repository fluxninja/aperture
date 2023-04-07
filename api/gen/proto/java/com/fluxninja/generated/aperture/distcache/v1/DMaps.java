// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/distcache/v1/stats.proto

package com.fluxninja.generated.aperture.distcache.v1;

/**
 * Protobuf type {@code aperture.distcache.v1.DMaps}
 */
public final class DMaps extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:aperture.distcache.v1.DMaps)
    DMapsOrBuilder {
private static final long serialVersionUID = 0L;
  // Use DMaps.newBuilder() to construct.
  private DMaps(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private DMaps() {
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new DMaps();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private DMaps(
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
          case 8: {

            entriesTotal_ = input.readInt64();
            break;
          }
          case 16: {

            deleteHits_ = input.readInt64();
            break;
          }
          case 24: {

            deleteMisses_ = input.readInt64();
            break;
          }
          case 32: {

            getMisses_ = input.readInt64();
            break;
          }
          case 40: {

            getHits_ = input.readInt64();
            break;
          }
          case 48: {

            evictedTotal_ = input.readInt64();
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
    return com.fluxninja.generated.aperture.distcache.v1.StatsProto.internal_static_aperture_distcache_v1_DMaps_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.fluxninja.generated.aperture.distcache.v1.StatsProto.internal_static_aperture_distcache_v1_DMaps_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.fluxninja.generated.aperture.distcache.v1.DMaps.class, com.fluxninja.generated.aperture.distcache.v1.DMaps.Builder.class);
  }

  public static final int ENTRIES_TOTAL_FIELD_NUMBER = 1;
  private long entriesTotal_;
  /**
   * <code>int64 entries_total = 1 [json_name = "EntriesTotal"];</code>
   * @return The entriesTotal.
   */
  @java.lang.Override
  public long getEntriesTotal() {
    return entriesTotal_;
  }

  public static final int DELETE_HITS_FIELD_NUMBER = 2;
  private long deleteHits_;
  /**
   * <code>int64 delete_hits = 2 [json_name = "DeleteHits"];</code>
   * @return The deleteHits.
   */
  @java.lang.Override
  public long getDeleteHits() {
    return deleteHits_;
  }

  public static final int DELETE_MISSES_FIELD_NUMBER = 3;
  private long deleteMisses_;
  /**
   * <code>int64 delete_misses = 3 [json_name = "DeleteMisses"];</code>
   * @return The deleteMisses.
   */
  @java.lang.Override
  public long getDeleteMisses() {
    return deleteMisses_;
  }

  public static final int GET_MISSES_FIELD_NUMBER = 4;
  private long getMisses_;
  /**
   * <code>int64 get_misses = 4 [json_name = "GetMisses"];</code>
   * @return The getMisses.
   */
  @java.lang.Override
  public long getGetMisses() {
    return getMisses_;
  }

  public static final int GET_HITS_FIELD_NUMBER = 5;
  private long getHits_;
  /**
   * <code>int64 get_hits = 5 [json_name = "GetHits"];</code>
   * @return The getHits.
   */
  @java.lang.Override
  public long getGetHits() {
    return getHits_;
  }

  public static final int EVICTED_TOTAL_FIELD_NUMBER = 6;
  private long evictedTotal_;
  /**
   * <code>int64 evicted_total = 6 [json_name = "EvictedTotal"];</code>
   * @return The evictedTotal.
   */
  @java.lang.Override
  public long getEvictedTotal() {
    return evictedTotal_;
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
    if (entriesTotal_ != 0L) {
      output.writeInt64(1, entriesTotal_);
    }
    if (deleteHits_ != 0L) {
      output.writeInt64(2, deleteHits_);
    }
    if (deleteMisses_ != 0L) {
      output.writeInt64(3, deleteMisses_);
    }
    if (getMisses_ != 0L) {
      output.writeInt64(4, getMisses_);
    }
    if (getHits_ != 0L) {
      output.writeInt64(5, getHits_);
    }
    if (evictedTotal_ != 0L) {
      output.writeInt64(6, evictedTotal_);
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (entriesTotal_ != 0L) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt64Size(1, entriesTotal_);
    }
    if (deleteHits_ != 0L) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt64Size(2, deleteHits_);
    }
    if (deleteMisses_ != 0L) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt64Size(3, deleteMisses_);
    }
    if (getMisses_ != 0L) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt64Size(4, getMisses_);
    }
    if (getHits_ != 0L) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt64Size(5, getHits_);
    }
    if (evictedTotal_ != 0L) {
      size += com.google.protobuf.CodedOutputStream
        .computeInt64Size(6, evictedTotal_);
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
    if (!(obj instanceof com.fluxninja.generated.aperture.distcache.v1.DMaps)) {
      return super.equals(obj);
    }
    com.fluxninja.generated.aperture.distcache.v1.DMaps other = (com.fluxninja.generated.aperture.distcache.v1.DMaps) obj;

    if (getEntriesTotal()
        != other.getEntriesTotal()) return false;
    if (getDeleteHits()
        != other.getDeleteHits()) return false;
    if (getDeleteMisses()
        != other.getDeleteMisses()) return false;
    if (getGetMisses()
        != other.getGetMisses()) return false;
    if (getGetHits()
        != other.getGetHits()) return false;
    if (getEvictedTotal()
        != other.getEvictedTotal()) return false;
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
    hash = (37 * hash) + ENTRIES_TOTAL_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        getEntriesTotal());
    hash = (37 * hash) + DELETE_HITS_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        getDeleteHits());
    hash = (37 * hash) + DELETE_MISSES_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        getDeleteMisses());
    hash = (37 * hash) + GET_MISSES_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        getGetMisses());
    hash = (37 * hash) + GET_HITS_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        getGetHits());
    hash = (37 * hash) + EVICTED_TOTAL_FIELD_NUMBER;
    hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
        getEvictedTotal());
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static com.fluxninja.generated.aperture.distcache.v1.DMaps parseFrom(
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
  public static Builder newBuilder(com.fluxninja.generated.aperture.distcache.v1.DMaps prototype) {
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
   * Protobuf type {@code aperture.distcache.v1.DMaps}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:aperture.distcache.v1.DMaps)
      com.fluxninja.generated.aperture.distcache.v1.DMapsOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.fluxninja.generated.aperture.distcache.v1.StatsProto.internal_static_aperture_distcache_v1_DMaps_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.fluxninja.generated.aperture.distcache.v1.StatsProto.internal_static_aperture_distcache_v1_DMaps_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.fluxninja.generated.aperture.distcache.v1.DMaps.class, com.fluxninja.generated.aperture.distcache.v1.DMaps.Builder.class);
    }

    // Construct using com.fluxninja.generated.aperture.distcache.v1.DMaps.newBuilder()
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
      entriesTotal_ = 0L;

      deleteHits_ = 0L;

      deleteMisses_ = 0L;

      getMisses_ = 0L;

      getHits_ = 0L;

      evictedTotal_ = 0L;

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.fluxninja.generated.aperture.distcache.v1.StatsProto.internal_static_aperture_distcache_v1_DMaps_descriptor;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.distcache.v1.DMaps getDefaultInstanceForType() {
      return com.fluxninja.generated.aperture.distcache.v1.DMaps.getDefaultInstance();
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.distcache.v1.DMaps build() {
      com.fluxninja.generated.aperture.distcache.v1.DMaps result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.fluxninja.generated.aperture.distcache.v1.DMaps buildPartial() {
      com.fluxninja.generated.aperture.distcache.v1.DMaps result = new com.fluxninja.generated.aperture.distcache.v1.DMaps(this);
      result.entriesTotal_ = entriesTotal_;
      result.deleteHits_ = deleteHits_;
      result.deleteMisses_ = deleteMisses_;
      result.getMisses_ = getMisses_;
      result.getHits_ = getHits_;
      result.evictedTotal_ = evictedTotal_;
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
      if (other instanceof com.fluxninja.generated.aperture.distcache.v1.DMaps) {
        return mergeFrom((com.fluxninja.generated.aperture.distcache.v1.DMaps)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.fluxninja.generated.aperture.distcache.v1.DMaps other) {
      if (other == com.fluxninja.generated.aperture.distcache.v1.DMaps.getDefaultInstance()) return this;
      if (other.getEntriesTotal() != 0L) {
        setEntriesTotal(other.getEntriesTotal());
      }
      if (other.getDeleteHits() != 0L) {
        setDeleteHits(other.getDeleteHits());
      }
      if (other.getDeleteMisses() != 0L) {
        setDeleteMisses(other.getDeleteMisses());
      }
      if (other.getGetMisses() != 0L) {
        setGetMisses(other.getGetMisses());
      }
      if (other.getGetHits() != 0L) {
        setGetHits(other.getGetHits());
      }
      if (other.getEvictedTotal() != 0L) {
        setEvictedTotal(other.getEvictedTotal());
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
      com.fluxninja.generated.aperture.distcache.v1.DMaps parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (com.fluxninja.generated.aperture.distcache.v1.DMaps) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private long entriesTotal_ ;
    /**
     * <code>int64 entries_total = 1 [json_name = "EntriesTotal"];</code>
     * @return The entriesTotal.
     */
    @java.lang.Override
    public long getEntriesTotal() {
      return entriesTotal_;
    }
    /**
     * <code>int64 entries_total = 1 [json_name = "EntriesTotal"];</code>
     * @param value The entriesTotal to set.
     * @return This builder for chaining.
     */
    public Builder setEntriesTotal(long value) {
      
      entriesTotal_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>int64 entries_total = 1 [json_name = "EntriesTotal"];</code>
     * @return This builder for chaining.
     */
    public Builder clearEntriesTotal() {
      
      entriesTotal_ = 0L;
      onChanged();
      return this;
    }

    private long deleteHits_ ;
    /**
     * <code>int64 delete_hits = 2 [json_name = "DeleteHits"];</code>
     * @return The deleteHits.
     */
    @java.lang.Override
    public long getDeleteHits() {
      return deleteHits_;
    }
    /**
     * <code>int64 delete_hits = 2 [json_name = "DeleteHits"];</code>
     * @param value The deleteHits to set.
     * @return This builder for chaining.
     */
    public Builder setDeleteHits(long value) {
      
      deleteHits_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>int64 delete_hits = 2 [json_name = "DeleteHits"];</code>
     * @return This builder for chaining.
     */
    public Builder clearDeleteHits() {
      
      deleteHits_ = 0L;
      onChanged();
      return this;
    }

    private long deleteMisses_ ;
    /**
     * <code>int64 delete_misses = 3 [json_name = "DeleteMisses"];</code>
     * @return The deleteMisses.
     */
    @java.lang.Override
    public long getDeleteMisses() {
      return deleteMisses_;
    }
    /**
     * <code>int64 delete_misses = 3 [json_name = "DeleteMisses"];</code>
     * @param value The deleteMisses to set.
     * @return This builder for chaining.
     */
    public Builder setDeleteMisses(long value) {
      
      deleteMisses_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>int64 delete_misses = 3 [json_name = "DeleteMisses"];</code>
     * @return This builder for chaining.
     */
    public Builder clearDeleteMisses() {
      
      deleteMisses_ = 0L;
      onChanged();
      return this;
    }

    private long getMisses_ ;
    /**
     * <code>int64 get_misses = 4 [json_name = "GetMisses"];</code>
     * @return The getMisses.
     */
    @java.lang.Override
    public long getGetMisses() {
      return getMisses_;
    }
    /**
     * <code>int64 get_misses = 4 [json_name = "GetMisses"];</code>
     * @param value The getMisses to set.
     * @return This builder for chaining.
     */
    public Builder setGetMisses(long value) {
      
      getMisses_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>int64 get_misses = 4 [json_name = "GetMisses"];</code>
     * @return This builder for chaining.
     */
    public Builder clearGetMisses() {
      
      getMisses_ = 0L;
      onChanged();
      return this;
    }

    private long getHits_ ;
    /**
     * <code>int64 get_hits = 5 [json_name = "GetHits"];</code>
     * @return The getHits.
     */
    @java.lang.Override
    public long getGetHits() {
      return getHits_;
    }
    /**
     * <code>int64 get_hits = 5 [json_name = "GetHits"];</code>
     * @param value The getHits to set.
     * @return This builder for chaining.
     */
    public Builder setGetHits(long value) {
      
      getHits_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>int64 get_hits = 5 [json_name = "GetHits"];</code>
     * @return This builder for chaining.
     */
    public Builder clearGetHits() {
      
      getHits_ = 0L;
      onChanged();
      return this;
    }

    private long evictedTotal_ ;
    /**
     * <code>int64 evicted_total = 6 [json_name = "EvictedTotal"];</code>
     * @return The evictedTotal.
     */
    @java.lang.Override
    public long getEvictedTotal() {
      return evictedTotal_;
    }
    /**
     * <code>int64 evicted_total = 6 [json_name = "EvictedTotal"];</code>
     * @param value The evictedTotal to set.
     * @return This builder for chaining.
     */
    public Builder setEvictedTotal(long value) {
      
      evictedTotal_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>int64 evicted_total = 6 [json_name = "EvictedTotal"];</code>
     * @return This builder for chaining.
     */
    public Builder clearEvictedTotal() {
      
      evictedTotal_ = 0L;
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


    // @@protoc_insertion_point(builder_scope:aperture.distcache.v1.DMaps)
  }

  // @@protoc_insertion_point(class_scope:aperture.distcache.v1.DMaps)
  private static final com.fluxninja.generated.aperture.distcache.v1.DMaps DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.fluxninja.generated.aperture.distcache.v1.DMaps();
  }

  public static com.fluxninja.generated.aperture.distcache.v1.DMaps getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<DMaps>
      PARSER = new com.google.protobuf.AbstractParser<DMaps>() {
    @java.lang.Override
    public DMaps parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new DMaps(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<DMaps> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<DMaps> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.fluxninja.generated.aperture.distcache.v1.DMaps getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

