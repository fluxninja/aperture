// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/decisions/v1/policydecisions.proto

package com.aperture.policy.decisions.v1;

public final class PolicydecisionsProto {
  private PolicydecisionsProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_decisions_v1_LoadShedDecision_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_decisions_v1_LoadShedDecision_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_decisions_v1_TokensDecision_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_decisions_v1_TokensDecision_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_decisions_v1_TokensDecision_TokensByWorkloadIndexEntry_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_decisions_v1_TokensDecision_TokensByWorkloadIndexEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_decisions_v1_RateLimiterDecision_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_decisions_v1_RateLimiterDecision_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n2aperture/policy/decisions/v1/policydec" +
      "isions.proto\022\034aperture.policy.decisions." +
      "v1\"<\n\020LoadShedDecision\022(\n\020load_shed_fact" +
      "or\030\001 \001(\001R\016loadShedFactor\"\335\001\n\016TokensDecis" +
      "ion\022\200\001\n\030tokens_by_workload_index\030\001 \003(\0132G" +
      ".aperture.policy.decisions.v1.TokensDeci" +
      "sion.TokensByWorkloadIndexEntryR\025tokensB" +
      "yWorkloadIndex\032H\n\032TokensByWorkloadIndexE" +
      "ntry\022\020\n\003key\030\001 \001(\tR\003key\022\024\n\005value\030\002 \001(\004R\005v" +
      "alue:\0028\001\"+\n\023RateLimiterDecision\022\024\n\005limit" +
      "\030\001 \001(\001R\005limitB\244\002\n com.aperture.policy.de" +
      "cisions.v1B\024PolicydecisionsProtoP\001ZWgith" +
      "ub.com/fluxninja/aperture/api/gen/proto/" +
      "go/aperture/policy/decisions/v1;decision" +
      "sv1\242\002\003APD\252\002\034Aperture.Policy.Decisions.V1" +
      "\312\002\034Aperture\\Policy\\Decisions\\V1\342\002(Apertu" +
      "re\\Policy\\Decisions\\V1\\GPBMetadata\352\002\037Ape" +
      "rture::Policy::Decisions::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
        });
    internal_static_aperture_policy_decisions_v1_LoadShedDecision_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_policy_decisions_v1_LoadShedDecision_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_decisions_v1_LoadShedDecision_descriptor,
        new java.lang.String[] { "LoadShedFactor", });
    internal_static_aperture_policy_decisions_v1_TokensDecision_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_policy_decisions_v1_TokensDecision_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_decisions_v1_TokensDecision_descriptor,
        new java.lang.String[] { "TokensByWorkloadIndex", });
    internal_static_aperture_policy_decisions_v1_TokensDecision_TokensByWorkloadIndexEntry_descriptor =
      internal_static_aperture_policy_decisions_v1_TokensDecision_descriptor.getNestedTypes().get(0);
    internal_static_aperture_policy_decisions_v1_TokensDecision_TokensByWorkloadIndexEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_decisions_v1_TokensDecision_TokensByWorkloadIndexEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_aperture_policy_decisions_v1_RateLimiterDecision_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_policy_decisions_v1_RateLimiterDecision_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_decisions_v1_RateLimiterDecision_descriptor,
        new java.lang.String[] { "Limit", });
  }

  // @@protoc_insertion_point(outer_class_scope)
}
