// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/service/auth/v3/authz_stripped.proto

package com.fluxninja.generated.envoy.service.auth.v3;

public interface EnvoyInternalAddressOrBuilder extends
    // @@protoc_insertion_point(interface_extends:envoy.service.auth.v3.EnvoyInternalAddress)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * Specifies the :ref:`name &lt;envoy_v3_api_field_config.listener.v3.Listener.name&gt;` of the
   * internal listener.
   * </pre>
   *
   * <code>string server_listener_name = 1 [json_name = "serverListenerName"];</code>
   * @return Whether the serverListenerName field is set.
   */
  boolean hasServerListenerName();
  /**
   * <pre>
   * Specifies the :ref:`name &lt;envoy_v3_api_field_config.listener.v3.Listener.name&gt;` of the
   * internal listener.
   * </pre>
   *
   * <code>string server_listener_name = 1 [json_name = "serverListenerName"];</code>
   * @return The serverListenerName.
   */
  java.lang.String getServerListenerName();
  /**
   * <pre>
   * Specifies the :ref:`name &lt;envoy_v3_api_field_config.listener.v3.Listener.name&gt;` of the
   * internal listener.
   * </pre>
   *
   * <code>string server_listener_name = 1 [json_name = "serverListenerName"];</code>
   * @return The bytes for serverListenerName.
   */
  com.google.protobuf.ByteString
      getServerListenerNameBytes();

  /**
   * <pre>
   * Specifies an endpoint identifier to distinguish between multiple endpoints for the same internal listener in a
   * single upstream pool. Only used in the upstream addresses for tracking changes to individual endpoints. This, for
   * example, may be set to the final destination IP for the target internal listener.
   * </pre>
   *
   * <code>string endpoint_id = 2 [json_name = "endpointId"];</code>
   * @return The endpointId.
   */
  java.lang.String getEndpointId();
  /**
   * <pre>
   * Specifies an endpoint identifier to distinguish between multiple endpoints for the same internal listener in a
   * single upstream pool. Only used in the upstream addresses for tracking changes to individual endpoints. This, for
   * example, may be set to the final destination IP for the target internal listener.
   * </pre>
   *
   * <code>string endpoint_id = 2 [json_name = "endpointId"];</code>
   * @return The bytes for endpointId.
   */
  com.google.protobuf.ByteString
      getEndpointIdBytes();

  public com.fluxninja.generated.envoy.service.auth.v3.EnvoyInternalAddress.AddressNameSpecifierCase getAddressNameSpecifierCase();
}
