// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/common/watchdog/v1/watchdog.proto

package com.aperture.common.watchdog.v1;

public final class WatchdogProto {
  private WatchdogProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_common_watchdog_v1_HeapResult_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_common_watchdog_v1_HeapResult_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_aperture_common_watchdog_v1_WatchdogResult_descriptor;
  static final
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_aperture_common_watchdog_v1_WatchdogResult_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n*aperture/common/watchdog/v1/watchdog.p" +
      "roto\022\033aperture.common.watchdog.v1\032\036googl" +
      "e/protobuf/duration.proto\"\327\004\n\nHeapResult" +
      "\022\024\n\005limit\030\001 \001(\004R\005limit\022\037\n\013heap_marked\030\002 " +
      "\001(\004R\nheapMarked\022\034\n\tthreshold\030\003 \001(\004R\tthre" +
      "shold\022\033\n\tcurr_gogc\030\004 \001(\005R\010currGogc\022#\n\ror" +
      "iginal_gogc\030\005 \001(\005R\014originalGogc\022\037\n\013total" +
      "_alloc\030\006 \001(\004R\ntotalAlloc\022\020\n\003sys\030\007 \001(\004R\003s" +
      "ys\022\030\n\007mallocs\030\010 \001(\004R\007mallocs\022\024\n\005frees\030\t " +
      "\001(\004R\005frees\022\035\n\nheap_alloc\030\n \001(\004R\theapAllo" +
      "c\022\031\n\010heap_sys\030\013 \001(\004R\007heapSys\022\033\n\theap_idl" +
      "e\030\014 \001(\004R\010heapIdle\022\035\n\nheap_inuse\030\r \001(\004R\th" +
      "eapInuse\022#\n\rheap_released\030\016 \001(\004R\014heapRel" +
      "eased\022!\n\014heap_objects\030\017 \001(\004R\013heapObjects" +
      "\022\027\n\007next_gc\030\020 \001(\004R\006nextGc\022\027\n\007last_gc\030\021 \001" +
      "(\004R\006lastGc\022$\n\016pause_total_ns\030\022 \001(\004R\014paus" +
      "eTotalNs\022\025\n\006num_gc\030\023 \001(\rR\005numGc\022\"\n\rnum_f" +
      "orced_gc\030\024 \001(\rR\013numForcedGc\"\227\001\n\016Watchdog" +
      "Result\022\024\n\005total\030\001 \001(\004R\005total\022\022\n\004used\030\002 \001" +
      "(\004R\004used\022\034\n\tthreshold\030\003 \001(\004R\tthreshold\022=" +
      "\n\rforce_gc_took\030\004 \001(\0132\031.google.protobuf." +
      "DurationR\013forceGcTookB\226\002\n\037com.aperture.c" +
      "ommon.watchdog.v1B\rWatchdogProtoP\001ZUgith" +
      "ub.com/fluxninja/aperture/api/gen/proto/" +
      "go/aperture/common/watchdog/v1;watchdogv" +
      "1\242\002\003ACW\252\002\033Aperture.Common.Watchdog.V1\312\002\033" +
      "Aperture\\Common\\Watchdog\\V1\342\002\'Aperture\\C" +
      "ommon\\Watchdog\\V1\\GPBMetadata\352\002\036Aperture" +
      "::Common::Watchdog::V1b\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.google.protobuf.DurationProto.getDescriptor(),
        });
    internal_static_aperture_common_watchdog_v1_HeapResult_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_aperture_common_watchdog_v1_HeapResult_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_common_watchdog_v1_HeapResult_descriptor,
        new java.lang.String[] { "Limit", "HeapMarked", "Threshold", "CurrGogc", "OriginalGogc", "TotalAlloc", "Sys", "Mallocs", "Frees", "HeapAlloc", "HeapSys", "HeapIdle", "HeapInuse", "HeapReleased", "HeapObjects", "NextGc", "LastGc", "PauseTotalNs", "NumGc", "NumForcedGc", });
    internal_static_aperture_common_watchdog_v1_WatchdogResult_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_aperture_common_watchdog_v1_WatchdogResult_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_aperture_common_watchdog_v1_WatchdogResult_descriptor,
        new java.lang.String[] { "Total", "Used", "Threshold", "ForceGcTook", });
    com.google.protobuf.DurationProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
