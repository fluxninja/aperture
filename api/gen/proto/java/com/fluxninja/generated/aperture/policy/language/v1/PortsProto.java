// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/policy/language/v1/ports.proto

package com.fluxninja.generated.aperture.policy.language.v1;

public final class PortsProto {
  private PortsProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_language_v1_InPort_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_language_v1_InPort_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_language_v1_OutPort_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_language_v1_OutPort_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_policy_language_v1_ConstantSignal_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_policy_language_v1_ConstantSignal_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n\'aperture/policy/language/v1/ports.prot" +
      "o\022\033aperture.policy.language.v1\"\214\001\n\006InPor" +
      "t\022!\n\013signal_name\030\001 \001(\tH\000R\nsignalName\022V\n\017" +
      "constant_signal\030\002 \001(\0132+.aperture.policy." +
      "language.v1.ConstantSignalH\000R\016constantSi" +
      "gnalB\007\n\005value\"*\n\007OutPort\022\037\n\013signal_name\030" +
      "\001 \001(\tR\nsignalName\"X\n\016ConstantSignal\022%\n\rs" +
      "pecial_value\030\001 \001(\tH\000R\014specialValue\022\026\n\005va" +
      "lue\030\002 \001(\001H\000R\005valueB\007\n\005constB\247\002\n3com.flux" +
      "ninja.generated.aperture.policy.language" +
      ".v1B\nPortsProtoP\001ZUgithub.com/fluxninja/" +
      "aperture/api/gen/proto/go/aperture/polic" +
      "y/language/v1;languagev1\242\002\003APL\252\002\033Apertur" +
      "e.Policy.Language.V1\312\002\033Aperture\\Policy\\L" +
      "anguage\\V1\342\002\'Aperture\\Policy\\Language\\V1" +
      "\\GPBMetadata\352\002\036Aperture::Policy::Languag" +
      "e::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
        });
    internal_static_aperture_policy_language_v1_InPort_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_policy_language_v1_InPort_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_language_v1_InPort_descriptor,
        new java.lang.String[] { "SignalName", "ConstantSignal", "Value", });
    internal_static_aperture_policy_language_v1_OutPort_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_policy_language_v1_OutPort_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_language_v1_OutPort_descriptor,
        new java.lang.String[] { "SignalName", });
    internal_static_aperture_policy_language_v1_ConstantSignal_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_policy_language_v1_ConstantSignal_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_policy_language_v1_ConstantSignal_descriptor,
        new java.lang.String[] { "SpecialValue", "Value", "Const", });
  }

  // @@protoc_insertion_point(outer_class_scope)
}
