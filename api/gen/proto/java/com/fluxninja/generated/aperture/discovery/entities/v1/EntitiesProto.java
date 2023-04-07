// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/discovery/entities/v1/entities.proto

package com.fluxninja.generated.aperture.discovery.entities.v1;

public final class EntitiesProto {
  private EntitiesProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_discovery_entities_v1_GetEntityByNameRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_discovery_entities_v1_GetEntityByNameRequest_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_discovery_entities_v1_Entities_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_discovery_entities_v1_Entities_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_discovery_entities_v1_Entities_Entities_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_discovery_entities_v1_Entities_Entities_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_discovery_entities_v1_Entities_Entities_EntitiesEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_discovery_entities_v1_Entities_Entities_EntitiesEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_discovery_entities_v1_Entity_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_discovery_entities_v1_Entity_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n-aperture/discovery/entities/v1/entitie" +
      "s.proto\022\036aperture.discovery.entities.v1\032" +
      "\034google/api/annotations.proto\032\033google/pr" +
      "otobuf/empty.proto\032.protoc-gen-openapiv2" +
      "/options/annotations.proto\"<\n\033GetEntityB" +
      "yIPAddressRequest\022\035\n\nip_address\030\001 \001(\tR\ti" +
      "pAddress\",\n\026GetEntityByNameRequest\022\022\n\004na" +
      "me\030\001 \001(\tR\004name\"\236\003\n\010Entities\022f\n\026entities_" +
      "by_ip_address\030\001 \001(\01321.aperture.discovery" +
      ".entities.v1.Entities.EntitiesR\023entities" +
      "ByIpAddress\022[\n\020entities_by_name\030\002 \001(\01321." +
      "aperture.discovery.entities.v1.Entities." +
      "EntitiesR\016entitiesByName\032\314\001\n\010Entities\022[\n" +
      "\010entities\030\001 \003(\0132?.aperture.discovery.ent" +
      "ities.v1.Entities.Entities.EntitiesEntry" +
      "R\010entities\032c\n\rEntitiesEntry\022\020\n\003key\030\001 \001(\t" +
      "R\003key\022<\n\005value\030\002 \001(\0132&.aperture.discover" +
      "y.entities.v1.EntityR\005value:\0028\001\"\244\001\n\006Enti" +
      "ty\022\020\n\003uid\030\001 \001(\tR\003uid\022\035\n\nip_address\030\002 \001(\t" +
      "R\tipAddress\022\022\n\004name\030\003 \001(\tR\004name\022\034\n\tnames" +
      "pace\030\004 \001(\tR\tnamespace\022\033\n\tnode_name\030\005 \001(\t" +
      "R\010nodeName\022\032\n\010services\030\006 \003(\tR\010services2\222" +
      "\004\n\017EntitiesService\022\202\001\n\013GetEntities\022\026.goo" +
      "gle.protobuf.Empty\032(.aperture.discovery." +
      "entities.v1.Entities\"1\222A\020\n\016aperture-agen" +
      "t\202\323\344\223\002\030\022\026/v1/discovery/entities\022\306\001\n\024GetE" +
      "ntityByIPAddress\022;.aperture.discovery.en" +
      "tities.v1.GetEntityByIPAddressRequest\032&." +
      "aperture.discovery.entities.v1.Entity\"I\222" +
      "A\020\n\016aperture-agent\202\323\344\223\0020\022./v1/discovery/" +
      "entities/ip-address/{ip_address}\022\260\001\n\017Get" +
      "EntityByName\0226.aperture.discovery.entiti" +
      "es.v1.GetEntityByNameRequest\032&.aperture." +
      "discovery.entities.v1.Entity\"=\222A\020\n\016apert" +
      "ure-agent\202\323\344\223\002$\022\"/v1/discovery/entities/" +
      "name/{name}B\274\002\n6com.fluxninja.generated." +
      "aperture.discovery.entities.v1B\rEntities" +
      "ProtoP\001ZXgithub.com/fluxninja/aperture/a" +
      "pi/gen/proto/go/aperture/discovery/entit" +
      "ies/v1;entitiesv1\242\002\003ADE\252\002\036Aperture.Disco" +
      "very.Entities.V1\312\002\036Aperture\\Discovery\\En" +
      "tities\\V1\342\002*Aperture\\Discovery\\Entities\\" +
      "V1\\GPBMetadata\352\002!Aperture::Discovery::En" +
      "tities::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.fluxninja.generated.google.api.AnnotationsProto.getDescriptor(),
          com.google.protobuf.EmptyProto.getDescriptor(),
          com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.getDescriptor(),
        });
    internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_discovery_entities_v1_GetEntityByIPAddressRequest_descriptor,
        new java.lang.String[] { "IpAddress", });
    internal_static_aperture_discovery_entities_v1_GetEntityByNameRequest_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_discovery_entities_v1_GetEntityByNameRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_discovery_entities_v1_GetEntityByNameRequest_descriptor,
        new java.lang.String[] { "Name", });
    internal_static_aperture_discovery_entities_v1_Entities_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_discovery_entities_v1_Entities_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_discovery_entities_v1_Entities_descriptor,
        new java.lang.String[] { "EntitiesByIpAddress", "EntitiesByName", });
    internal_static_aperture_discovery_entities_v1_Entities_Entities_descriptor =
      internal_static_aperture_discovery_entities_v1_Entities_descriptor.getNestedTypes().get(0);
    internal_static_aperture_discovery_entities_v1_Entities_Entities_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_discovery_entities_v1_Entities_Entities_descriptor,
        new java.lang.String[] { "Entities", });
    internal_static_aperture_discovery_entities_v1_Entities_Entities_EntitiesEntry_descriptor =
      internal_static_aperture_discovery_entities_v1_Entities_Entities_descriptor.getNestedTypes().get(0);
    internal_static_aperture_discovery_entities_v1_Entities_Entities_EntitiesEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_discovery_entities_v1_Entities_Entities_EntitiesEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_aperture_discovery_entities_v1_Entity_descriptor =
      getDescriptor().getMessageTypes().get(3);
    internal_static_aperture_discovery_entities_v1_Entity_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_discovery_entities_v1_Entity_descriptor,
        new java.lang.String[] { "Uid", "IpAddress", "Name", "Namespace", "NodeName", "Services", });
    com.google.protobuf.ExtensionRegistry registry =
        com.google.protobuf.ExtensionRegistry.newInstance();
    registry.add(com.fluxninja.generated.google.api.AnnotationsProto.http);
    registry.add(com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.openapiv2Operation);
    com.google.protobuf.Descriptors.FileDescriptor
        .internalUpdateFileDescriptor(descriptor, registry);
    com.fluxninja.generated.google.api.AnnotationsProto.getDescriptor();
    com.google.protobuf.EmptyProto.getDescriptor();
    com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
