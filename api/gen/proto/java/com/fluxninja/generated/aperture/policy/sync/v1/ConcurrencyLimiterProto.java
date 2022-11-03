// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/sync/v1/concurrency_limiter.proto

package com.fluxninja.generated.aperture.policy.sync.v1;

public final class ConcurrencyLimiterProto {
  private ConcurrencyLimiterProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_sync_v1_ConcurrencyLimiterWrapper_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_sync_v1_ConcurrencyLimiterWrapper_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_sync_v1_LoadDecisionWrapper_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_sync_v1_LoadDecisionWrapper_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_sync_v1_TokensDecisionWrapper_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_sync_v1_TokensDecisionWrapper_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_sync_v1_LoadDecision_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_sync_v1_LoadDecision_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_sync_v1_TokensDecision_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_sync_v1_TokensDecision_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_sync_v1_TokensDecision_TokensByWorkloadIndexEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_sync_v1_TokensDecision_TokensByWorkloadIndexEntry_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n1aperture/policy/sync/v1/concurrency_li" +
      "miter.proto\022\027aperture.policy.sync.v1\032(ap" +
      "erture/policy/language/v1/policy.proto\032/" +
      "aperture/policy/sync/v1/common_attribute" +
      "s.proto\032\"aperture/policy/sync/v1/tick.pr" +
      "oto\"\325\001\n\031ConcurrencyLimiterWrapper\022V\n\021com" +
      "mon_attributes\030\001 \001(\0132).aperture.policy.s" +
      "ync.v1.CommonAttributesR\020commonAttribute" +
      "s\022`\n\023concurrency_limiter\030\002 \001(\0132/.apertur" +
      "e.policy.language.v1.ConcurrencyLimiterR" +
      "\022concurrencyLimiter\"\271\001\n\023LoadDecisionWrap" +
      "per\022V\n\021common_attributes\030\001 \001(\0132).apertur" +
      "e.policy.sync.v1.CommonAttributesR\020commo" +
      "nAttributes\022J\n\rload_decision\030\002 \001(\0132%.ape" +
      "rture.policy.sync.v1.LoadDecisionR\014loadD" +
      "ecision\"\301\001\n\025TokensDecisionWrapper\022V\n\021com" +
      "mon_attributes\030\001 \001(\0132).aperture.policy.s" +
      "ync.v1.CommonAttributesR\020commonAttribute" +
      "s\022P\n\017tokens_decision\030\002 \001(\0132\'.aperture.po" +
      "licy.sync.v1.TokensDecisionR\016tokensDecis" +
      "ion\"\232\001\n\014LoadDecision\022\'\n\017load_multiplier\030" +
      "\001 \001(\001R\016loadMultiplier\022!\n\014pass_through\030\002 " +
      "\001(\010R\013passThrough\022>\n\ttick_info\030\003 \001(\0132!.ap" +
      "erture.policy.sync.v1.TickInfoR\010tickInfo" +
      "\"\327\001\n\016TokensDecision\022{\n\030tokens_by_workloa" +
      "d_index\030\001 \003(\0132B.aperture.policy.sync.v1." +
      "TokensDecision.TokensByWorkloadIndexEntr" +
      "yR\025tokensByWorkloadIndex\032H\n\032TokensByWork" +
      "loadIndexEntry\022\020\n\003key\030\001 \001(\tR\003key\022\024\n\005valu" +
      "e\030\002 \001(\004R\005value:\0028\001B\230\002\n/com.fluxninja.gen" +
      "erated.aperture.policy.sync.v1B\027Concurre" +
      "ncyLimiterProtoP\001ZMgithub.com/fluxninja/" +
      "aperture/api/gen/proto/go/aperture/polic" +
      "y/sync/v1;syncv1\242\002\003APS\252\002\027Aperture.Policy" +
      ".Sync.V1\312\002\027Aperture\\Policy\\Sync\\V1\342\002#Ape" +
      "rture\\Policy\\Sync\\V1\\GPBMetadata\352\002\032Apert" +
      "ure::Policy::Sync::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.getDescriptor(),
          com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesProto.getDescriptor(),
          com.fluxninja.generated.aperture.policy.sync.v1.TickProto.getDescriptor(),
        });
    internal_static_aperture_policy_sync_v1_ConcurrencyLimiterWrapper_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_policy_sync_v1_ConcurrencyLimiterWrapper_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_sync_v1_ConcurrencyLimiterWrapper_descriptor,
        new java.lang.String[] { "CommonAttributes", "ConcurrencyLimiter", });
    internal_static_aperture_policy_sync_v1_LoadDecisionWrapper_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_policy_sync_v1_LoadDecisionWrapper_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_sync_v1_LoadDecisionWrapper_descriptor,
        new java.lang.String[] { "CommonAttributes", "LoadDecision", });
    internal_static_aperture_policy_sync_v1_TokensDecisionWrapper_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_policy_sync_v1_TokensDecisionWrapper_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_sync_v1_TokensDecisionWrapper_descriptor,
        new java.lang.String[] { "CommonAttributes", "TokensDecision", });
    internal_static_aperture_policy_sync_v1_LoadDecision_descriptor =
      getDescriptor().getMessageTypes().get(3);
    internal_static_aperture_policy_sync_v1_LoadDecision_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_sync_v1_LoadDecision_descriptor,
        new java.lang.String[] { "LoadMultiplier", "PassThrough", "TickInfo", });
    internal_static_aperture_policy_sync_v1_TokensDecision_descriptor =
      getDescriptor().getMessageTypes().get(4);
    internal_static_aperture_policy_sync_v1_TokensDecision_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_sync_v1_TokensDecision_descriptor,
        new java.lang.String[] { "TokensByWorkloadIndex", });
    internal_static_aperture_policy_sync_v1_TokensDecision_TokensByWorkloadIndexEntry_descriptor =
      internal_static_aperture_policy_sync_v1_TokensDecision_descriptor.getNestedTypes().get(0);
    internal_static_aperture_policy_sync_v1_TokensDecision_TokensByWorkloadIndexEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_sync_v1_TokensDecision_TokensByWorkloadIndexEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    com.fluxninja.generated.aperture.policy.language.v1.PolicyProto.getDescriptor();
    com.fluxninja.generated.aperture.policy.sync.v1.CommonAttributesProto.getDescriptor();
    com.fluxninja.generated.aperture.policy.sync.v1.TickProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
