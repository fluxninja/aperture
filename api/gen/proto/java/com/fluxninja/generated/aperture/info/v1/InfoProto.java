// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/info/v1/info.proto

package com.fluxninja.generated.aperture.info.v1;

public final class InfoProto {
  private InfoProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_info_v1_VersionInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_info_v1_VersionInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_info_v1_ProcessInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_info_v1_ProcessInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_info_v1_HostInfo_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_info_v1_HostInfo_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n\033aperture/info/v1/info.proto\022\020aperture." +
      "info.v1\032\034google/api/annotations.proto\032\036g" +
      "oogle/protobuf/duration.proto\032\033google/pr" +
      "otobuf/empty.proto\032\037google/protobuf/time" +
      "stamp.proto\032.protoc-gen-openapiv2/option" +
      "s/annotations.proto\"\341\001\n\013VersionInfo\022\030\n\007v" +
      "ersion\030\001 \001(\tR\007version\022\030\n\007service\030\002 \001(\tR\007" +
      "service\022\035\n\nbuild_host\030\003 \001(\tR\tbuildHost\022\031" +
      "\n\010build_os\030\004 \001(\tR\007buildOs\022\035\n\nbuild_time\030" +
      "\005 \001(\tR\tbuildTime\022\035\n\ngit_branch\030\006 \001(\tR\tgi" +
      "tBranch\022&\n\017git_commit_hash\030\007 \001(\tR\rgitCom" +
      "mitHash\"\233\001\n\013ProcessInfo\0229\n\nstart_time\030\001 " +
      "\001(\0132\032.google.protobuf.TimestampR\tstartTi" +
      "me\0221\n\006uptime\030\002 \001(\0132\031.google.protobuf.Dur" +
      "ationR\006uptime\022\036\n\nextensions\030\003 \003(\tR\nexten" +
      "sions\"U\n\010HostInfo\022\022\n\004uuid\030\001 \001(\tR\004uuid\022\032\n" +
      "\010hostname\030\002 \001(\tR\010hostname\022\031\n\010local_ip\030\003 " +
      "\001(\tR\007localIp2\222\003\n\013InfoService\022\202\001\n\007Version" +
      "\022\026.google.protobuf.Empty\032\035.aperture.info" +
      ".v1.VersionInfo\"@\222A%\n\016aperture-agent\n\023ap" +
      "erture-controller\202\323\344\223\002\022\022\020/v1/info/versio" +
      "n\022\202\001\n\007Process\022\026.google.protobuf.Empty\032\035." +
      "aperture.info.v1.ProcessInfo\"@\222A%\n\016apert" +
      "ure-agent\n\023aperture-controller\202\323\344\223\002\022\022\020/v" +
      "1/info/process\022y\n\004Host\022\026.google.protobuf" +
      ".Empty\032\032.aperture.info.v1.HostInfo\"=\222A%\n" +
      "\016aperture-agent\n\023aperture-controller\202\323\344\223" +
      "\002\017\022\r/v1/info/hostB\337\001\n(com.fluxninja.gene" +
      "rated.aperture.info.v1B\tInfoProtoP\001ZFgit" +
      "hub.com/fluxninja/aperture/api/gen/proto" +
      "/go/aperture/info/v1;infov1\242\002\003AIX\252\002\020Aper" +
      "ture.Info.V1\312\002\020Aperture\\Info\\V1\342\002\034Apertu" +
      "re\\Info\\V1\\GPBMetadata\352\002\022Aperture::Info:" +
      ":V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.fluxninja.generated.google.api.AnnotationsProto.getDescriptor(),
          com.google.protobuf.DurationProto.getDescriptor(),
          com.google.protobuf.EmptyProto.getDescriptor(),
          com.google.protobuf.TimestampProto.getDescriptor(),
          com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.getDescriptor(),
        });
    internal_static_aperture_info_v1_VersionInfo_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_info_v1_VersionInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_info_v1_VersionInfo_descriptor,
        new java.lang.String[] { "Version", "Service", "BuildHost", "BuildOs", "BuildTime", "GitBranch", "GitCommitHash", });
    internal_static_aperture_info_v1_ProcessInfo_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_info_v1_ProcessInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_info_v1_ProcessInfo_descriptor,
        new java.lang.String[] { "StartTime", "Uptime", "Extensions", });
    internal_static_aperture_info_v1_HostInfo_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_info_v1_HostInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_info_v1_HostInfo_descriptor,
        new java.lang.String[] { "Uuid", "Hostname", "LocalIp", });
    com.google.protobuf.ExtensionRegistry registry =
        com.google.protobuf.ExtensionRegistry.newInstance();
    registry.add(com.fluxninja.generated.google.api.AnnotationsProto.http);
    registry.add(com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.openapiv2Operation);
    com.google.protobuf.Descriptors.FileDescriptor
        .internalUpdateFileDescriptor(descriptor, registry);
    com.fluxninja.generated.google.api.AnnotationsProto.getDescriptor();
    com.google.protobuf.DurationProto.getDescriptor();
    com.google.protobuf.EmptyProto.getDescriptor();
    com.google.protobuf.TimestampProto.getDescriptor();
    com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
