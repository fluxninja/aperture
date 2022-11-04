// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/service/auth/v3/authz_stripped.proto

package com.fluxninja.generated.envoy.service.auth.v3;

public interface OkHttpResponseOrBuilder extends
    // @@protoc_insertion_point(interface_extends:envoy.service.auth.v3.OkHttpResponse)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * HTTP entity headers in addition to the original request headers. This allows the authorization
   * service to append, to add or to override headers from the original request before
   * dispatching it to the upstream. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;` defaults to
   * false when used in this message. By setting the ``append`` field to ``true``,
   * the filter will append the correspondent header value to the matched request header.
   * By leaving ``append`` as false, the filter will either add a new header, or override an existing
   * one if there is a match.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption headers = 2 [json_name = "headers"];</code>
   */
  java.util.List<com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption> 
      getHeadersList();
  /**
   * <pre>
   * HTTP entity headers in addition to the original request headers. This allows the authorization
   * service to append, to add or to override headers from the original request before
   * dispatching it to the upstream. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;` defaults to
   * false when used in this message. By setting the ``append`` field to ``true``,
   * the filter will append the correspondent header value to the matched request header.
   * By leaving ``append`` as false, the filter will either add a new header, or override an existing
   * one if there is a match.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption headers = 2 [json_name = "headers"];</code>
   */
  com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption getHeaders(int index);
  /**
   * <pre>
   * HTTP entity headers in addition to the original request headers. This allows the authorization
   * service to append, to add or to override headers from the original request before
   * dispatching it to the upstream. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;` defaults to
   * false when used in this message. By setting the ``append`` field to ``true``,
   * the filter will append the correspondent header value to the matched request header.
   * By leaving ``append`` as false, the filter will either add a new header, or override an existing
   * one if there is a match.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption headers = 2 [json_name = "headers"];</code>
   */
  int getHeadersCount();
  /**
   * <pre>
   * HTTP entity headers in addition to the original request headers. This allows the authorization
   * service to append, to add or to override headers from the original request before
   * dispatching it to the upstream. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;` defaults to
   * false when used in this message. By setting the ``append`` field to ``true``,
   * the filter will append the correspondent header value to the matched request header.
   * By leaving ``append`` as false, the filter will either add a new header, or override an existing
   * one if there is a match.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption headers = 2 [json_name = "headers"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOptionOrBuilder> 
      getHeadersOrBuilderList();
  /**
   * <pre>
   * HTTP entity headers in addition to the original request headers. This allows the authorization
   * service to append, to add or to override headers from the original request before
   * dispatching it to the upstream. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;` defaults to
   * false when used in this message. By setting the ``append`` field to ``true``,
   * the filter will append the correspondent header value to the matched request header.
   * By leaving ``append`` as false, the filter will either add a new header, or override an existing
   * one if there is a match.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption headers = 2 [json_name = "headers"];</code>
   */
  com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOptionOrBuilder getHeadersOrBuilder(
      int index);

  /**
   * <pre>
   * HTTP entity headers to remove from the original request before dispatching
   * it to the upstream. This allows the authorization service to act on auth
   * related headers (like ``Authorization``), process them, and consume them.
   * Under this model, the upstream will either receive the request (if it's
   * authorized) or not receive it (if it's not), but will not see headers
   * containing authorization credentials.
   * Pseudo headers (such as ``:authority``, ``:method``, ``:path`` etc), as well as
   * the header ``Host``, may not be removed as that would make the request
   * malformed. If mentioned in ``headers_to_remove`` these special headers will
   * be ignored.
   * When using the HTTP service this must instead be set by the HTTP
   * authorization service as a comma separated list like so:
   * ``x-envoy-auth-headers-to-remove: one-auth-header, another-auth-header``.
   * </pre>
   *
   * <code>repeated string headers_to_remove = 5 [json_name = "headersToRemove"];</code>
   * @return A list containing the headersToRemove.
   */
  java.util.List<java.lang.String>
      getHeadersToRemoveList();
  /**
   * <pre>
   * HTTP entity headers to remove from the original request before dispatching
   * it to the upstream. This allows the authorization service to act on auth
   * related headers (like ``Authorization``), process them, and consume them.
   * Under this model, the upstream will either receive the request (if it's
   * authorized) or not receive it (if it's not), but will not see headers
   * containing authorization credentials.
   * Pseudo headers (such as ``:authority``, ``:method``, ``:path`` etc), as well as
   * the header ``Host``, may not be removed as that would make the request
   * malformed. If mentioned in ``headers_to_remove`` these special headers will
   * be ignored.
   * When using the HTTP service this must instead be set by the HTTP
   * authorization service as a comma separated list like so:
   * ``x-envoy-auth-headers-to-remove: one-auth-header, another-auth-header``.
   * </pre>
   *
   * <code>repeated string headers_to_remove = 5 [json_name = "headersToRemove"];</code>
   * @return The count of headersToRemove.
   */
  int getHeadersToRemoveCount();
  /**
   * <pre>
   * HTTP entity headers to remove from the original request before dispatching
   * it to the upstream. This allows the authorization service to act on auth
   * related headers (like ``Authorization``), process them, and consume them.
   * Under this model, the upstream will either receive the request (if it's
   * authorized) or not receive it (if it's not), but will not see headers
   * containing authorization credentials.
   * Pseudo headers (such as ``:authority``, ``:method``, ``:path`` etc), as well as
   * the header ``Host``, may not be removed as that would make the request
   * malformed. If mentioned in ``headers_to_remove`` these special headers will
   * be ignored.
   * When using the HTTP service this must instead be set by the HTTP
   * authorization service as a comma separated list like so:
   * ``x-envoy-auth-headers-to-remove: one-auth-header, another-auth-header``.
   * </pre>
   *
   * <code>repeated string headers_to_remove = 5 [json_name = "headersToRemove"];</code>
   * @param index The index of the element to return.
   * @return The headersToRemove at the given index.
   */
  java.lang.String getHeadersToRemove(int index);
  /**
   * <pre>
   * HTTP entity headers to remove from the original request before dispatching
   * it to the upstream. This allows the authorization service to act on auth
   * related headers (like ``Authorization``), process them, and consume them.
   * Under this model, the upstream will either receive the request (if it's
   * authorized) or not receive it (if it's not), but will not see headers
   * containing authorization credentials.
   * Pseudo headers (such as ``:authority``, ``:method``, ``:path`` etc), as well as
   * the header ``Host``, may not be removed as that would make the request
   * malformed. If mentioned in ``headers_to_remove`` these special headers will
   * be ignored.
   * When using the HTTP service this must instead be set by the HTTP
   * authorization service as a comma separated list like so:
   * ``x-envoy-auth-headers-to-remove: one-auth-header, another-auth-header``.
   * </pre>
   *
   * <code>repeated string headers_to_remove = 5 [json_name = "headersToRemove"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the headersToRemove at the given index.
   */
  com.google.protobuf.ByteString
      getHeadersToRemoveBytes(int index);

  /**
   * <pre>
   * This field has been deprecated in favor of :ref:`CheckResponse.dynamic_metadata
   * &lt;envoy_v3_api_field_service.auth.v3.CheckResponse.dynamic_metadata&gt;`. Until it is removed,
   * setting this field overrides :ref:`CheckResponse.dynamic_metadata
   * &lt;envoy_v3_api_field_service.auth.v3.CheckResponse.dynamic_metadata&gt;`.
   * </pre>
   *
   * <code>.google.protobuf.Struct dynamic_metadata = 3 [json_name = "dynamicMetadata", deprecated = true, (.envoy.annotations.deprecated_at_minor_version) = "3.0"];</code>
   * @deprecated envoy.service.auth.v3.OkHttpResponse.dynamic_metadata is deprecated.
   *     See envoy/service/auth/v3/authz_stripped.proto;l=102
   * @return Whether the dynamicMetadata field is set.
   */
  @java.lang.Deprecated boolean hasDynamicMetadata();
  /**
   * <pre>
   * This field has been deprecated in favor of :ref:`CheckResponse.dynamic_metadata
   * &lt;envoy_v3_api_field_service.auth.v3.CheckResponse.dynamic_metadata&gt;`. Until it is removed,
   * setting this field overrides :ref:`CheckResponse.dynamic_metadata
   * &lt;envoy_v3_api_field_service.auth.v3.CheckResponse.dynamic_metadata&gt;`.
   * </pre>
   *
   * <code>.google.protobuf.Struct dynamic_metadata = 3 [json_name = "dynamicMetadata", deprecated = true, (.envoy.annotations.deprecated_at_minor_version) = "3.0"];</code>
   * @deprecated envoy.service.auth.v3.OkHttpResponse.dynamic_metadata is deprecated.
   *     See envoy/service/auth/v3/authz_stripped.proto;l=102
   * @return The dynamicMetadata.
   */
  @java.lang.Deprecated com.google.protobuf.Struct getDynamicMetadata();
  /**
   * <pre>
   * This field has been deprecated in favor of :ref:`CheckResponse.dynamic_metadata
   * &lt;envoy_v3_api_field_service.auth.v3.CheckResponse.dynamic_metadata&gt;`. Until it is removed,
   * setting this field overrides :ref:`CheckResponse.dynamic_metadata
   * &lt;envoy_v3_api_field_service.auth.v3.CheckResponse.dynamic_metadata&gt;`.
   * </pre>
   *
   * <code>.google.protobuf.Struct dynamic_metadata = 3 [json_name = "dynamicMetadata", deprecated = true, (.envoy.annotations.deprecated_at_minor_version) = "3.0"];</code>
   */
  @java.lang.Deprecated com.google.protobuf.StructOrBuilder getDynamicMetadataOrBuilder();

  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client on success. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;`
   * defaults to false when used in this message.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption response_headers_to_add = 6 [json_name = "responseHeadersToAdd"];</code>
   */
  java.util.List<com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption> 
      getResponseHeadersToAddList();
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client on success. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;`
   * defaults to false when used in this message.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption response_headers_to_add = 6 [json_name = "responseHeadersToAdd"];</code>
   */
  com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption getResponseHeadersToAdd(int index);
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client on success. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;`
   * defaults to false when used in this message.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption response_headers_to_add = 6 [json_name = "responseHeadersToAdd"];</code>
   */
  int getResponseHeadersToAddCount();
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client on success. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;`
   * defaults to false when used in this message.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption response_headers_to_add = 6 [json_name = "responseHeadersToAdd"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOptionOrBuilder> 
      getResponseHeadersToAddOrBuilderList();
  /**
   * <pre>
   * This field allows the authorization service to send HTTP response headers
   * to the downstream client on success. Note that the :ref:`append field in HeaderValueOption &lt;envoy_v3_api_field_config.core.v3.HeaderValueOption.append&gt;`
   * defaults to false when used in this message.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.HeaderValueOption response_headers_to_add = 6 [json_name = "responseHeadersToAdd"];</code>
   */
  com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOptionOrBuilder getResponseHeadersToAddOrBuilder(
      int index);

  /**
   * <pre>
   * This field allows the authorization service to set (and overwrite) query
   * string parameters on the original request before it is sent upstream.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.QueryParameter query_parameters_to_set = 7 [json_name = "queryParametersToSet"];</code>
   */
  java.util.List<com.fluxninja.generated.envoy.service.auth.v3.QueryParameter> 
      getQueryParametersToSetList();
  /**
   * <pre>
   * This field allows the authorization service to set (and overwrite) query
   * string parameters on the original request before it is sent upstream.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.QueryParameter query_parameters_to_set = 7 [json_name = "queryParametersToSet"];</code>
   */
  com.fluxninja.generated.envoy.service.auth.v3.QueryParameter getQueryParametersToSet(int index);
  /**
   * <pre>
   * This field allows the authorization service to set (and overwrite) query
   * string parameters on the original request before it is sent upstream.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.QueryParameter query_parameters_to_set = 7 [json_name = "queryParametersToSet"];</code>
   */
  int getQueryParametersToSetCount();
  /**
   * <pre>
   * This field allows the authorization service to set (and overwrite) query
   * string parameters on the original request before it is sent upstream.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.QueryParameter query_parameters_to_set = 7 [json_name = "queryParametersToSet"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.envoy.service.auth.v3.QueryParameterOrBuilder> 
      getQueryParametersToSetOrBuilderList();
  /**
   * <pre>
   * This field allows the authorization service to set (and overwrite) query
   * string parameters on the original request before it is sent upstream.
   * </pre>
   *
   * <code>repeated .envoy.service.auth.v3.QueryParameter query_parameters_to_set = 7 [json_name = "queryParametersToSet"];</code>
   */
  com.fluxninja.generated.envoy.service.auth.v3.QueryParameterOrBuilder getQueryParametersToSetOrBuilder(
      int index);

  /**
   * <pre>
   * This field allows the authorization service to specify which query parameters
   * should be removed from the original request before it is sent upstream. Each
   * element in this list is a case-sensitive query parameter name to be removed.
   * </pre>
   *
   * <code>repeated string query_parameters_to_remove = 8 [json_name = "queryParametersToRemove"];</code>
   * @return A list containing the queryParametersToRemove.
   */
  java.util.List<java.lang.String>
      getQueryParametersToRemoveList();
  /**
   * <pre>
   * This field allows the authorization service to specify which query parameters
   * should be removed from the original request before it is sent upstream. Each
   * element in this list is a case-sensitive query parameter name to be removed.
   * </pre>
   *
   * <code>repeated string query_parameters_to_remove = 8 [json_name = "queryParametersToRemove"];</code>
   * @return The count of queryParametersToRemove.
   */
  int getQueryParametersToRemoveCount();
  /**
   * <pre>
   * This field allows the authorization service to specify which query parameters
   * should be removed from the original request before it is sent upstream. Each
   * element in this list is a case-sensitive query parameter name to be removed.
   * </pre>
   *
   * <code>repeated string query_parameters_to_remove = 8 [json_name = "queryParametersToRemove"];</code>
   * @param index The index of the element to return.
   * @return The queryParametersToRemove at the given index.
   */
  java.lang.String getQueryParametersToRemove(int index);
  /**
   * <pre>
   * This field allows the authorization service to specify which query parameters
   * should be removed from the original request before it is sent upstream. Each
   * element in this list is a case-sensitive query parameter name to be removed.
   * </pre>
   *
   * <code>repeated string query_parameters_to_remove = 8 [json_name = "queryParametersToRemove"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the queryParametersToRemove at the given index.
   */
  com.google.protobuf.ByteString
      getQueryParametersToRemoveBytes(int index);
}
