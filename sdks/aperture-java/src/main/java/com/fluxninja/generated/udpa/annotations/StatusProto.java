// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: udpa/annotations/status.proto

package com.fluxninja.generated.udpa.annotations;

public final class StatusProto {
  private StatusProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
    registry.add(com.fluxninja.generated.udpa.annotations.StatusProto.fileStatus);
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  public static final int FILE_STATUS_FIELD_NUMBER = 222707719;
  /**
   * <code>extend .google.protobuf.FileOptions { ... }</code>
   */
  public static final
    com.google.protobuf.GeneratedMessage.GeneratedExtension<
      com.google.protobuf.DescriptorProtos.FileOptions,
      com.fluxninja.generated.udpa.annotations.StatusAnnotation> fileStatus = com.google.protobuf.GeneratedMessage
          .newFileScopedGeneratedExtension(
        com.fluxninja.generated.udpa.annotations.StatusAnnotation.class,
        com.fluxninja.generated.udpa.annotations.StatusAnnotation.getDefaultInstance());
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_udpa_annotations_StatusAnnotation_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_udpa_annotations_StatusAnnotation_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n\035udpa/annotations/status.proto\022\020udpa.an" +
      "notations\032 google/protobuf/descriptor.pr" +
      "oto\"\232\001\n\020StatusAnnotation\022(\n\020work_in_prog" +
      "ress\030\001 \001(\010R\016workInProgress\022\\\n\026package_ve" +
      "rsion_status\030\002 \001(\0162&.udpa.annotations.Pa" +
      "ckageVersionStatusR\024packageVersionStatus" +
      "*]\n\024PackageVersionStatus\022\013\n\007UNKNOWN\020\000\022\n\n" +
      "\006FROZEN\020\001\022\n\n\006ACTIVE\020\002\022 \n\034NEXT_MAJOR_VERS" +
      "ION_CANDIDATE\020\003:d\n\013file_status\022\034.google." +
      "protobuf.FileOptions\030\207\200\231j \001(\0132\".udpa.ann" +
      "otations.StatusAnnotationR\nfileStatusB\274\001" +
      "\n(com.fluxninja.generated.udpa.annotatio" +
      "nsB\013StatusProtoP\001Z\"github.com/cncf/xds/g" +
      "o/annotations\242\002\003UAX\252\002\020Udpa.Annotations\312\002" +
      "\020Udpa\\Annotations\342\002\034Udpa\\Annotations\\GPB" +
      "Metadata\352\002\021Udpa::Annotationsb\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.DescriptorProtos.getDescriptor(),
        });
    internal_static_udpa_annotations_StatusAnnotation_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_udpa_annotations_StatusAnnotation_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_udpa_annotations_StatusAnnotation_descriptor,
        new java.lang.String[] { "WorkInProgress", "PackageVersionStatus", });
    fileStatus.internalInit(descriptor.getExtensions().get(0));
    com.google.protobuf.DescriptorProtos.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
