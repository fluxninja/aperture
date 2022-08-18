// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/common/info/v1/info.proto

package com.aperture.common.info.v1;

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
    internal_static_aperture_common_info_v1_VersionInfo_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_common_info_v1_VersionInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_common_info_v1_ProcessInfo_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_common_info_v1_ProcessInfo_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_common_info_v1_HostInfo_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_common_info_v1_HostInfo_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n\"aperture/common/info/v1/info.proto\022\027ap" +
      "erture.common.info.v1\032\034google/api/annota" +
      "tions.proto\032\036google/protobuf/duration.pr" +
      "oto\032\033google/protobuf/empty.proto\032\037google" +
      "/protobuf/timestamp.proto\"\341\001\n\013VersionInf" +
      "o\022\030\n\007version\030\001 \001(\tR\007version\022\030\n\007service\030\002" +
      " \001(\tR\007service\022\035\n\nbuild_host\030\003 \001(\tR\tbuild" +
      "Host\022\031\n\010build_os\030\004 \001(\tR\007buildOs\022\035\n\nbuild" +
      "_time\030\005 \001(\tR\tbuildTime\022\035\n\ngit_branch\030\006 \001" +
      "(\tR\tgitBranch\022&\n\017git_commit_hash\030\007 \001(\tR\r" +
      "gitCommitHash\"{\n\013ProcessInfo\0229\n\nstart_ti" +
      "me\030\001 \001(\0132\032.google.protobuf.TimestampR\tst" +
      "artTime\0221\n\006uptime\030\002 \001(\0132\031.google.protobu" +
      "f.DurationR\006uptime\"&\n\010HostInfo\022\032\n\010hostna" +
      "me\030\001 \001(\tR\010hostname2\255\002\n\013InfoService\022a\n\007Ve" +
      "rsion\022\026.google.protobuf.Empty\032$.aperture" +
      ".common.info.v1.VersionInfo\"\030\202\323\344\223\002\022\022\020/v1" +
      "/info/version\022a\n\007Process\022\026.google.protob" +
      "uf.Empty\032$.aperture.common.info.v1.Proce" +
      "ssInfo\"\030\202\323\344\223\002\022\022\020/v1/info/process\022X\n\004Host" +
      "\022\026.google.protobuf.Empty\032!.aperture.comm" +
      "on.info.v1.HostInfo\"\025\202\323\344\223\002\017\022\r/v1/info/ho" +
      "stB\366\001\n\033com.aperture.common.info.v1B\tInfo" +
      "ProtoP\001ZMgithub.com/fluxninja/aperture/a" +
      "pi/gen/proto/go/aperture/common/info/v1;" +
      "infov1\242\002\003ACI\252\002\027Aperture.Common.Info.V1\312\002" +
      "\027Aperture\\Common\\Info\\V1\342\002#Aperture\\Comm" +
      "on\\Info\\V1\\GPBMetadata\352\002\032Aperture::Commo" +
      "n::Info::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.api.AnnotationsProto.getDescriptor(),
          com.google.protobuf.DurationProto.getDescriptor(),
          com.google.protobuf.EmptyProto.getDescriptor(),
          com.google.protobuf.TimestampProto.getDescriptor(),
        });
    internal_static_aperture_common_info_v1_VersionInfo_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_common_info_v1_VersionInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_common_info_v1_VersionInfo_descriptor,
        new java.lang.String[] { "Version", "Service", "BuildHost", "BuildOs", "BuildTime", "GitBranch", "GitCommitHash", });
    internal_static_aperture_common_info_v1_ProcessInfo_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_common_info_v1_ProcessInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_common_info_v1_ProcessInfo_descriptor,
        new java.lang.String[] { "StartTime", "Uptime", });
    internal_static_aperture_common_info_v1_HostInfo_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_common_info_v1_HostInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_common_info_v1_HostInfo_descriptor,
        new java.lang.String[] { "Hostname", });
    com.google.protobuf.ExtensionRegistry registry =
        com.google.protobuf.ExtensionRegistry.newInstance();
    registry.add(com.google.api.AnnotationsProto.http);
    com.google.protobuf.Descriptors.FileDescriptor
        .internalUpdateFileDescriptor(descriptor, registry);
    com.google.api.AnnotationsProto.getDescriptor();
    com.google.protobuf.DurationProto.getDescriptor();
    com.google.protobuf.EmptyProto.getDescriptor();
    com.google.protobuf.TimestampProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
