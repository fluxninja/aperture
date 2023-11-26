// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

public interface CheckRequestOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.flowcontrol.check.v1.CheckRequest)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>string control_point = 1 [json_name = "controlPoint"];</code>
   * @return The controlPoint.
   */
  java.lang.String getControlPoint();
  /**
   * <code>string control_point = 1 [json_name = "controlPoint"];</code>
   * @return The bytes for controlPoint.
   */
  com.google.protobuf.ByteString
      getControlPointBytes();

  /**
   * <code>map&lt;string, string&gt; labels = 2 [json_name = "labels"];</code>
   */
  int getLabelsCount();
  /**
   * <code>map&lt;string, string&gt; labels = 2 [json_name = "labels"];</code>
   */
  boolean containsLabels(
      java.lang.String key);
  /**
   * Use {@link #getLabelsMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, java.lang.String>
  getLabels();
  /**
   * <code>map&lt;string, string&gt; labels = 2 [json_name = "labels"];</code>
   */
  java.util.Map<java.lang.String, java.lang.String>
  getLabelsMap();
  /**
   * <code>map&lt;string, string&gt; labels = 2 [json_name = "labels"];</code>
   */
  /* nullable */
java.lang.String getLabelsOrDefault(
      java.lang.String key,
      /* nullable */
java.lang.String defaultValue);
  /**
   * <code>map&lt;string, string&gt; labels = 2 [json_name = "labels"];</code>
   */
  java.lang.String getLabelsOrThrow(
      java.lang.String key);

  /**
   * <code>bool ramp_mode = 3 [json_name = "rampMode"];</code>
   * @return The rampMode.
   */
  boolean getRampMode();

  /**
   * <pre>
   * Key for result cache that needs to be fetched.
   * </pre>
   *
   * <code>string result_cache_key = 4 [json_name = "resultCacheKey"];</code>
   * @return The resultCacheKey.
   */
  java.lang.String getResultCacheKey();
  /**
   * <pre>
   * Key for result cache that needs to be fetched.
   * </pre>
   *
   * <code>string result_cache_key = 4 [json_name = "resultCacheKey"];</code>
   * @return The bytes for resultCacheKey.
   */
  com.google.protobuf.ByteString
      getResultCacheKeyBytes();

  /**
   * <pre>
   * Keys for state cache entries that need to be fetched.
   * </pre>
   *
   * <code>repeated string state_cache_keys = 5 [json_name = "stateCacheKeys"];</code>
   * @return A list containing the stateCacheKeys.
   */
  java.util.List<java.lang.String>
      getStateCacheKeysList();
  /**
   * <pre>
   * Keys for state cache entries that need to be fetched.
   * </pre>
   *
   * <code>repeated string state_cache_keys = 5 [json_name = "stateCacheKeys"];</code>
   * @return The count of stateCacheKeys.
   */
  int getStateCacheKeysCount();
  /**
   * <pre>
   * Keys for state cache entries that need to be fetched.
   * </pre>
   *
   * <code>repeated string state_cache_keys = 5 [json_name = "stateCacheKeys"];</code>
   * @param index The index of the element to return.
   * @return The stateCacheKeys at the given index.
   */
  java.lang.String getStateCacheKeys(int index);
  /**
   * <pre>
   * Keys for state cache entries that need to be fetched.
   * </pre>
   *
   * <code>repeated string state_cache_keys = 5 [json_name = "stateCacheKeys"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the stateCacheKeys at the given index.
   */
  com.google.protobuf.ByteString
      getStateCacheKeysBytes(int index);
}
