// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/monitoring/v1/policy_view.proto

package com.fluxninja.generated.aperture.policy.monitoring.v1;

public final class PolicyViewProto {
  private PolicyViewProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_monitoring_v1_PortView_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_monitoring_v1_PortView_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_monitoring_v1_ComponentView_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_monitoring_v1_ComponentView_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_monitoring_v1_SourceTarget_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_monitoring_v1_SourceTarget_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_monitoring_v1_Link_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_monitoring_v1_Link_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n/aperture/policy/monitoring/v1/policy_v" +
      "iew.proto\022\035aperture.policy.monitoring.v1" +
      "\032\034google/protobuf/struct.proto\"\224\001\n\010PortV" +
      "iew\022\033\n\tport_name\030\001 \001(\tR\010portName\022!\n\013sign" +
      "al_name\030\002 \001(\tH\000R\nsignalName\022\'\n\016constant_" +
      "value\030\003 \001(\001H\000R\rconstantValue\022\026\n\006looped\030\004" +
      " \001(\010R\006loopedB\007\n\005value\"\246\003\n\rComponentView\022" +
      "!\n\014component_id\030\001 \001(\tR\013componentId\022%\n\016co" +
      "mponent_name\030\002 \001(\tR\rcomponentName\022%\n\016com" +
      "ponent_type\030\003 \001(\tR\rcomponentType\0223\n\025comp" +
      "onent_description\030\004 \001(\tR\024componentDescri" +
      "ption\0225\n\tcomponent\030\005 \001(\0132\027.google.protob" +
      "uf.StructR\tcomponent\022B\n\010in_ports\030\006 \003(\0132\'" +
      ".aperture.policy.monitoring.v1.PortViewR" +
      "\007inPorts\022D\n\tout_ports\030\007 \003(\0132\'.aperture.p" +
      "olicy.monitoring.v1.PortViewR\010outPorts\022." +
      "\n\023parent_component_id\030\010 \001(\tR\021parentCompo" +
      "nentId\"N\n\014SourceTarget\022!\n\014component_id\030\001" +
      " \001(\tR\013componentId\022\033\n\tport_name\030\002 \001(\tR\010po" +
      "rtName\"\261\001\n\004Link\022C\n\006source\030\001 \001(\0132+.apertu" +
      "re.policy.monitoring.v1.SourceTargetR\006so" +
      "urce\022C\n\006target\030\002 \001(\0132+.aperture.policy.m" +
      "onitoring.v1.SourceTargetR\006target\022\037\n\013sig" +
      "nal_name\030\003 \001(\tR\nsignalNameB\272\002\n5com.fluxn" +
      "inja.generated.aperture.policy.monitorin" +
      "g.v1B\017PolicyViewProtoP\001ZYgithub.com/flux" +
      "ninja/aperture/api/gen/proto/go/aperture" +
      "/policy/monitoring/v1;monitoringv1\242\002\003APM" +
      "\252\002\035Aperture.Policy.Monitoring.V1\312\002\035Apert" +
      "ure\\Policy\\Monitoring\\V1\342\002)Aperture\\Poli" +
      "cy\\Monitoring\\V1\\GPBMetadata\352\002 Aperture:" +
      ":Policy::Monitoring::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.StructProto.getDescriptor(),
        });
    internal_static_aperture_policy_monitoring_v1_PortView_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_policy_monitoring_v1_PortView_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_monitoring_v1_PortView_descriptor,
        new java.lang.String[] { "PortName", "SignalName", "ConstantValue", "Looped", "Value", });
    internal_static_aperture_policy_monitoring_v1_ComponentView_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_policy_monitoring_v1_ComponentView_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_monitoring_v1_ComponentView_descriptor,
        new java.lang.String[] { "ComponentId", "ComponentName", "ComponentType", "ComponentDescription", "Component", "InPorts", "OutPorts", "ParentComponentId", });
    internal_static_aperture_policy_monitoring_v1_SourceTarget_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_policy_monitoring_v1_SourceTarget_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_monitoring_v1_SourceTarget_descriptor,
        new java.lang.String[] { "ComponentId", "PortName", });
    internal_static_aperture_policy_monitoring_v1_Link_descriptor =
      getDescriptor().getMessageTypes().get(3);
    internal_static_aperture_policy_monitoring_v1_Link_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_monitoring_v1_Link_descriptor,
        new java.lang.String[] { "Source", "Target", "SignalName", });
    com.google.protobuf.StructProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
