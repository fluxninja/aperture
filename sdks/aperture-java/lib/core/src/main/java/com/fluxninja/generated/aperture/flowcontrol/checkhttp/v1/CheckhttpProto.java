// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/checkhttp/v1/checkhttp.proto

package com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1;

public final class CheckhttpProto {
  private CheckhttpProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_HeadersEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_HeadersEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_HeadersEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_HeadersEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_HeadersEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_HeadersEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPResponse_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPResponse_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_flowcontrol_checkhttp_v1_SocketAddress_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_flowcontrol_checkhttp_v1_SocketAddress_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n1aperture/flowcontrol/checkhttp/v1/chec" +
      "khttp.proto\022!aperture.flowcontrol.checkh" +
      "ttp.v1\032\034google/api/annotations.proto\032\034go" +
      "ogle/protobuf/struct.proto\032\027google/rpc/s" +
      "tatus.proto\032.protoc-gen-openapiv2/option" +
      "s/annotations.proto\032\027validate/validate.p" +
      "roto\"\235\005\n\020CheckHTTPRequest\022H\n\006source\030\001 \001(" +
      "\01320.aperture.flowcontrol.checkhttp.v1.So" +
      "cketAddressR\006source\022R\n\013destination\030\002 \001(\013" +
      "20.aperture.flowcontrol.checkhttp.v1.Soc" +
      "ketAddressR\013destination\022Y\n\007request\030\003 \001(\013" +
      "2?.aperture.flowcontrol.checkhttp.v1.Che" +
      "ckHTTPRequest.HttpRequestR\007request\022#\n\rco" +
      "ntrol_point\030\004 \001(\tR\014controlPoint\022\033\n\tramp_" +
      "mode\030\005 \001(\010R\010rampMode\032\315\002\n\013HttpRequest\022\026\n\006" +
      "method\030\001 \001(\tR\006method\022f\n\007headers\030\002 \003(\0132L." +
      "aperture.flowcontrol.checkhttp.v1.CheckH" +
      "TTPRequest.HttpRequest.HeadersEntryR\007hea" +
      "ders\022\022\n\004path\030\003 \001(\tR\004path\022\022\n\004host\030\004 \001(\tR\004" +
      "host\022\026\n\006scheme\030\005 \001(\tR\006scheme\022\022\n\004size\030\006 \001" +
      "(\003R\004size\022\032\n\010protocol\030\007 \001(\tR\010protocol\022\022\n\004" +
      "body\030\010 \001(\tR\004body\032:\n\014HeadersEntry\022\020\n\003key\030" +
      "\001 \001(\tR\003key\022\024\n\005value\030\002 \001(\tR\005value:\0028\001\"\332\001\n" +
      "\022DeniedHttpResponse\022\026\n\006status\030\001 \001(\005R\006sta" +
      "tus\022\\\n\007headers\030\002 \003(\0132B.aperture.flowcont" +
      "rol.checkhttp.v1.DeniedHttpResponse.Head" +
      "ersEntryR\007headers\022\022\n\004body\030\003 \001(\tR\004body\032:\n" +
      "\014HeadersEntry\022\020\n\003key\030\001 \001(\tR\003key\022\024\n\005value" +
      "\030\002 \001(\tR\005value:\0028\001\"\352\001\n\016OkHttpResponse\022X\n\007" +
      "headers\030\001 \003(\0132>.aperture.flowcontrol.che" +
      "ckhttp.v1.OkHttpResponse.HeadersEntryR\007h" +
      "eaders\022B\n\020dynamic_metadata\030\002 \001(\0132\027.googl" +
      "e.protobuf.StructR\017dynamicMetadata\032:\n\014He" +
      "adersEntry\022\020\n\003key\030\001 \001(\tR\003key\022\024\n\005value\030\002 " +
      "\001(\tR\005value:\0028\001\"\314\002\n\021CheckHTTPResponse\022*\n\006" +
      "status\030\001 \001(\0132\022.google.rpc.StatusR\006status" +
      "\022`\n\017denied_response\030\002 \001(\01325.aperture.flo" +
      "wcontrol.checkhttp.v1.DeniedHttpResponse" +
      "H\000R\016deniedResponse\022T\n\013ok_response\030\003 \001(\0132" +
      "1.aperture.flowcontrol.checkhttp.v1.OkHt" +
      "tpResponseH\000R\nokResponse\022B\n\020dynamic_meta" +
      "data\030\004 \001(\0132\027.google.protobuf.StructR\017dyn" +
      "amicMetadataB\017\n\rhttp_response\"\320\001\n\rSocket" +
      "Address\022_\n\010protocol\030\001 \001(\01629.aperture.flo" +
      "wcontrol.checkhttp.v1.SocketAddress.Prot" +
      "ocolB\010\372B\005\202\001\002\020\001R\010protocol\022!\n\007address\030\002 \001(" +
      "\tB\007\372B\004r\002\020\001R\007address\022\035\n\004port\030\003 \001(\rB\t\372B\006*\004" +
      "\030\377\377\003R\004port\"\034\n\010Protocol\022\007\n\003TCP\020\000\022\007\n\003UDP\020\001" +
      "2\312\001\n\026FlowControlServiceHTTP\022\257\001\n\tCheckHTT" +
      "P\0223.aperture.flowcontrol.checkhttp.v1.Ch" +
      "eckHTTPRequest\0324.aperture.flowcontrol.ch" +
      "eckhttp.v1.CheckHTTPResponse\"7\222A\020\n\016apert" +
      "ure-agent\202\323\344\223\002\036\"\031/v1/flowcontrol/checkht" +
      "tp:\001*B\323\002\n9com.fluxninja.generated.apertu" +
      "re.flowcontrol.checkhttp.v1B\016CheckhttpPr" +
      "otoP\001Z_github.com/fluxninja/aperture/api" +
      "/v2/gen/proto/go/aperture/flowcontrol/ch" +
      "eckhttp/v1;checkhttpv1\242\002\003AFC\252\002!Aperture." +
      "Flowcontrol.Checkhttp.V1\312\002!Aperture\\Flow" +
      "control\\Checkhttp\\V1\342\002-Aperture\\Flowcont" +
      "rol\\Checkhttp\\V1\\GPBMetadata\352\002$Aperture:" +
      ":Flowcontrol::Checkhttp::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.fluxninja.generated.google.api.AnnotationsProto.getDescriptor(),
          com.google.protobuf.StructProto.getDescriptor(),
          com.fluxninja.generated.google.rpc.StatusProto.getDescriptor(),
          com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.getDescriptor(),
          com.fluxninja.generated.validate.ValidateProto.getDescriptor(),
        });
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_descriptor,
        new java.lang.String[] { "Source", "Destination", "Request", "ControlPoint", "RampMode", });
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_descriptor =
      internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_descriptor.getNestedTypes().get(0);
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_descriptor,
        new java.lang.String[] { "Method", "Headers", "Path", "Host", "Scheme", "Size", "Protocol", "Body", });
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_HeadersEntry_descriptor =
      internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_descriptor.getNestedTypes().get(0);
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_HeadersEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPRequest_HttpRequest_HeadersEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_descriptor,
        new java.lang.String[] { "Status", "Headers", "Body", });
    internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_HeadersEntry_descriptor =
      internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_descriptor.getNestedTypes().get(0);
    internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_HeadersEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_DeniedHttpResponse_HeadersEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_descriptor,
        new java.lang.String[] { "Headers", "DynamicMetadata", });
    internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_HeadersEntry_descriptor =
      internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_descriptor.getNestedTypes().get(0);
    internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_HeadersEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_OkHttpResponse_HeadersEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPResponse_descriptor =
      getDescriptor().getMessageTypes().get(3);
    internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPResponse_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_CheckHTTPResponse_descriptor,
        new java.lang.String[] { "Status", "DeniedResponse", "OkResponse", "DynamicMetadata", "HttpResponse", });
    internal_static_aperture_flowcontrol_checkhttp_v1_SocketAddress_descriptor =
      getDescriptor().getMessageTypes().get(4);
    internal_static_aperture_flowcontrol_checkhttp_v1_SocketAddress_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_flowcontrol_checkhttp_v1_SocketAddress_descriptor,
        new java.lang.String[] { "Protocol", "Address", "Port", });
    com.google.protobuf.ExtensionRegistry registry =
        com.google.protobuf.ExtensionRegistry.newInstance();
    registry.add(com.fluxninja.generated.google.api.AnnotationsProto.http);
    registry.add(com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.openapiv2Operation);
    registry.add(com.fluxninja.generated.validate.ValidateProto.rules);
    com.google.protobuf.Descriptors.FileDescriptor
        .internalUpdateFileDescriptor(descriptor, registry);
    com.fluxninja.generated.google.api.AnnotationsProto.getDescriptor();
    com.google.protobuf.StructProto.getDescriptor();
    com.fluxninja.generated.google.rpc.StatusProto.getDescriptor();
    com.fluxninja.generated.grpc.gateway.protoc_gen_openapiv2.options.AnnotationsProto.getDescriptor();
    com.fluxninja.generated.validate.ValidateProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
