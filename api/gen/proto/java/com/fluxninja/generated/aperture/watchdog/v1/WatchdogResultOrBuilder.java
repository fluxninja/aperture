// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/watchdog/v1/watchdog.proto

package com.fluxninja.generated.aperture.watchdog.v1;

public interface WatchdogResultOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.watchdog.v1.WatchdogResult)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>uint64 total = 1 [json_name = "total"];</code>
   * @return The total.
   */
  long getTotal();

  /**
   * <code>uint64 used = 2 [json_name = "used"];</code>
   * @return The used.
   */
  long getUsed();

  /**
   * <code>uint64 threshold = 3 [json_name = "threshold"];</code>
   * @return The threshold.
   */
  long getThreshold();

  /**
   * <code>.google.protobuf.Duration force_gc_took = 4 [json_name = "forceGcTook"];</code>
   * @return Whether the forceGcTook field is set.
   */
  boolean hasForceGcTook();
  /**
   * <code>.google.protobuf.Duration force_gc_took = 4 [json_name = "forceGcTook"];</code>
   * @return The forceGcTook.
   */
  com.google.protobuf.Duration getForceGcTook();
  /**
   * <code>.google.protobuf.Duration force_gc_took = 4 [json_name = "forceGcTook"];</code>
   */
  com.google.protobuf.DurationOrBuilder getForceGcTookOrBuilder();
}
