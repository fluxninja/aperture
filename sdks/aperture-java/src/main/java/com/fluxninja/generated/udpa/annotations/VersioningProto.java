// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: udpa/annotations/versioning.proto

package com.fluxninja.generated.udpa.annotations;

public final class VersioningProto {
  private VersioningProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
    registry.add(com.fluxninja.generated.udpa.annotations.VersioningProto.versioning);
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  public static final int VERSIONING_FIELD_NUMBER = 7881811;
  /**
   * <pre>
   * Magic number derived from 0x78 ('x') 0x44 ('D') 0x53 ('S')
   * </pre>
   *
   * <code>extend .google.protobuf.MessageOptions { ... }</code>
   */
  public static final
    com.google.protobuf.GeneratedMessage.GeneratedExtension<
      com.google.protobuf.DescriptorProtos.MessageOptions,
      com.fluxninja.generated.udpa.annotations.VersioningAnnotation> versioning = com.google.protobuf.GeneratedMessage
          .newFileScopedGeneratedExtension(
        com.fluxninja.generated.udpa.annotations.VersioningAnnotation.class,
        com.fluxninja.generated.udpa.annotations.VersioningAnnotation.getDefaultInstance());
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_udpa_annotations_VersioningAnnotation_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_udpa_annotations_VersioningAnnotation_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n!udpa/annotations/versioning.proto\022\020udp" +
      "a.annotations\032 google/protobuf/descripto" +
      "r.proto\"J\n\024VersioningAnnotation\0222\n\025previ" +
      "ous_message_type\030\001 \001(\tR\023previousMessageT" +
      "ype:j\n\nversioning\022\037.google.protobuf.Mess" +
      "ageOptions\030\323\210\341\003 \001(\0132&.udpa.annotations.V" +
      "ersioningAnnotationR\nversioningB\335\001\n(com." +
      "fluxninja.generated.udpa.annotationsB\017Ve" +
      "rsioningProtoP\001Z?github.com/fluxninja/ap" +
      "erture/api/gen/proto/go/udpa/annotations" +
      "\242\002\003UAX\252\002\020Udpa.Annotations\312\002\020Udpa\\Annotat" +
      "ions\342\002\034Udpa\\Annotations\\GPBMetadata\352\002\021Ud" +
      "pa::Annotationsb\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.DescriptorProtos.getDescriptor(),
        });
    internal_static_udpa_annotations_VersioningAnnotation_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_udpa_annotations_VersioningAnnotation_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_udpa_annotations_VersioningAnnotation_descriptor,
        new java.lang.String[] { "PreviousMessageType", });
    versioning.internalInit(descriptor.getExtensions().get(0));
    com.google.protobuf.DescriptorProtos.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
