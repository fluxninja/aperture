// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/sync/v1/flux_meter.proto

package com.fluxninja.generated.aperture.policy.sync.v1;

public final class FluxMeterProto {
  private FluxMeterProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_sync_v1_FluxMeterWrapper_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_sync_v1_FluxMeterWrapper_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n(aperture/policy/sync/v1/flux_meter.pro" +
      "to\022\027aperture.policy.sync.v1\032-aperture/po" +
      "licy/language/v1/flowcontrol.proto\"\201\001\n\020F" +
      "luxMeterWrapper\022E\n\nflux_meter\030\001 \001(\0132&.ap" +
      "erture.policy.language.v1.FluxMeterR\tflu" +
      "xMeter\022&\n\017flux_meter_name\030\004 \001(\tR\rfluxMet" +
      "erNameB\217\002\n/com.fluxninja.generated.apert" +
      "ure.policy.sync.v1B\016FluxMeterProtoP\001ZMgi" +
      "thub.com/fluxninja/aperture/api/gen/prot" +
      "o/go/aperture/policy/sync/v1;syncv1\242\002\003AP" +
      "S\252\002\027Aperture.Policy.Sync.V1\312\002\027Aperture\\P" +
      "olicy\\Sync\\V1\342\002#Aperture\\Policy\\Sync\\V1\\" +
      "GPBMetadata\352\002\032Aperture::Policy::Sync::V1" +
      "b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.fluxninja.generated.aperture.policy.language.v1.FlowcontrolProto.getDescriptor(),
        });
    internal_static_aperture_policy_sync_v1_FluxMeterWrapper_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_policy_sync_v1_FluxMeterWrapper_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_sync_v1_FluxMeterWrapper_descriptor,
        new java.lang.String[] { "FluxMeter", "FluxMeterName", });
    com.fluxninja.generated.aperture.policy.language.v1.FlowcontrolProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
