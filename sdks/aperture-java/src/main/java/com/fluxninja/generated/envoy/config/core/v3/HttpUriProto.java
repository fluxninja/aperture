// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/config/core/v3/http_uri.proto

package com.fluxninja.generated.envoy.config.core.v3;

public final class HttpUriProto {
  private HttpUriProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_envoy_config_core_v3_HttpUri_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_envoy_config_core_v3_HttpUri_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n#envoy/config/core/v3/http_uri.proto\022\024e" +
      "nvoy.config.core.v3\032\036google/protobuf/dur" +
      "ation.proto\032\035udpa/annotations/status.pro" +
      "to\032!udpa/annotations/versioning.proto\032\027v" +
      "alidate/validate.proto\"\307\001\n\007HttpUri\022\031\n\003ur" +
      "i\030\001 \001(\tB\007\372B\004r\002\020\001R\003uri\022#\n\007cluster\030\002 \001(\tB\007" +
      "\372B\004r\002\020\001H\000R\007cluster\022?\n\007timeout\030\003 \001(\0132\031.go" +
      "ogle.protobuf.DurationB\n\372B\007\252\001\004\010\0012\000R\007time" +
      "out: \232\305\210\036\033\n\031envoy.api.v2.core.HttpUriB\031\n" +
      "\022http_upstream_type\022\003\370B\001B\373\001\n,com.fluxnin" +
      "ja.generated.envoy.config.core.v3B\014HttpU" +
      "riProtoP\001ZBgithub.com/envoyproxy/go-cont" +
      "rol-plane/envoy/config/core/v3;corev3\242\002\003" +
      "ECC\252\002\024Envoy.Config.Core.V3\312\002\024Envoy\\Confi" +
      "g\\Core\\V3\342\002 Envoy\\Config\\Core\\V3\\GPBMeta" +
      "data\352\002\027Envoy::Config::Core::V3\272\200\310\321\006\002\020\002b\006" +
      "proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.DurationProto.getDescriptor(),
          com.fluxninja.generated.udpa.annotations.StatusProto.getDescriptor(),
          com.fluxninja.generated.udpa.annotations.VersioningProto.getDescriptor(),
          com.fluxninja.generated.validate.ValidateProto.getDescriptor(),
        });
    internal_static_envoy_config_core_v3_HttpUri_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_envoy_config_core_v3_HttpUri_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_envoy_config_core_v3_HttpUri_descriptor,
        new java.lang.String[] { "Uri", "Cluster", "Timeout", "HttpUpstreamType", });
    com.google.protobuf.ExtensionRegistry registry =
        com.google.protobuf.ExtensionRegistry.newInstance();
    registry.add(com.fluxninja.generated.udpa.annotations.StatusProto.fileStatus);
    registry.add(com.fluxninja.generated.udpa.annotations.VersioningProto.versioning);
    registry.add(com.fluxninja.generated.validate.ValidateProto.required);
    registry.add(com.fluxninja.generated.validate.ValidateProto.rules);
    com.google.protobuf.Descriptors.FileDescriptor
        .internalUpdateFileDescriptor(descriptor, registry);
    com.google.protobuf.DurationProto.getDescriptor();
    com.fluxninja.generated.udpa.annotations.StatusProto.getDescriptor();
    com.fluxninja.generated.udpa.annotations.VersioningProto.getDescriptor();
    com.fluxninja.generated.validate.ValidateProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
