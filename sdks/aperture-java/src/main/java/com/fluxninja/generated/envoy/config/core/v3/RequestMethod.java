// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: envoy/config/core/v3/base.proto

package com.fluxninja.generated.envoy.config.core.v3;

/**
 * <pre>
 * HTTP request method.
 * </pre>
 *
 * Protobuf enum {@code envoy.config.core.v3.RequestMethod}
 */
public enum RequestMethod
    implements com.google.protobuf.ProtocolMessageEnum {
  /**
   * <code>METHOD_UNSPECIFIED = 0;</code>
   */
  METHOD_UNSPECIFIED(0),
  /**
   * <code>GET = 1;</code>
   */
  GET(1),
  /**
   * <code>HEAD = 2;</code>
   */
  HEAD(2),
  /**
   * <code>POST = 3;</code>
   */
  POST(3),
  /**
   * <code>PUT = 4;</code>
   */
  PUT(4),
  /**
   * <code>DELETE = 5;</code>
   */
  DELETE(5),
  /**
   * <code>CONNECT = 6;</code>
   */
  CONNECT(6),
  /**
   * <code>OPTIONS = 7;</code>
   */
  OPTIONS(7),
  /**
   * <code>TRACE = 8;</code>
   */
  TRACE(8),
  /**
   * <code>PATCH = 9;</code>
   */
  PATCH(9),
  UNRECOGNIZED(-1),
  ;

  /**
   * <code>METHOD_UNSPECIFIED = 0;</code>
   */
  public static final int METHOD_UNSPECIFIED_VALUE = 0;
  /**
   * <code>GET = 1;</code>
   */
  public static final int GET_VALUE = 1;
  /**
   * <code>HEAD = 2;</code>
   */
  public static final int HEAD_VALUE = 2;
  /**
   * <code>POST = 3;</code>
   */
  public static final int POST_VALUE = 3;
  /**
   * <code>PUT = 4;</code>
   */
  public static final int PUT_VALUE = 4;
  /**
   * <code>DELETE = 5;</code>
   */
  public static final int DELETE_VALUE = 5;
  /**
   * <code>CONNECT = 6;</code>
   */
  public static final int CONNECT_VALUE = 6;
  /**
   * <code>OPTIONS = 7;</code>
   */
  public static final int OPTIONS_VALUE = 7;
  /**
   * <code>TRACE = 8;</code>
   */
  public static final int TRACE_VALUE = 8;
  /**
   * <code>PATCH = 9;</code>
   */
  public static final int PATCH_VALUE = 9;


  public final int getNumber() {
    if (this == UNRECOGNIZED) {
      throw new java.lang.IllegalArgumentException(
          "Can't get the number of an unknown enum value.");
    }
    return value;
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   * @deprecated Use {@link #forNumber(int)} instead.
   */
  @java.lang.Deprecated
  public static RequestMethod valueOf(int value) {
    return forNumber(value);
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   */
  public static RequestMethod forNumber(int value) {
    switch (value) {
      case 0: return METHOD_UNSPECIFIED;
      case 1: return GET;
      case 2: return HEAD;
      case 3: return POST;
      case 4: return PUT;
      case 5: return DELETE;
      case 6: return CONNECT;
      case 7: return OPTIONS;
      case 8: return TRACE;
      case 9: return PATCH;
      default: return null;
    }
  }

  public static com.google.protobuf.Internal.EnumLiteMap<RequestMethod>
      internalGetValueMap() {
    return internalValueMap;
  }
  private static final com.google.protobuf.Internal.EnumLiteMap<
      RequestMethod> internalValueMap =
        new com.google.protobuf.Internal.EnumLiteMap<RequestMethod>() {
          public RequestMethod findValueByNumber(int number) {
            return RequestMethod.forNumber(number);
          }
        };

  public final com.google.protobuf.Descriptors.EnumValueDescriptor
      getValueDescriptor() {
    if (this == UNRECOGNIZED) {
      throw new java.lang.IllegalStateException(
          "Can't get the descriptor of an unrecognized enum value.");
    }
    return getDescriptor().getValues().get(ordinal());
  }
  public final com.google.protobuf.Descriptors.EnumDescriptor
      getDescriptorForType() {
    return getDescriptor();
  }
  public static final com.google.protobuf.Descriptors.EnumDescriptor
      getDescriptor() {
    return com.fluxninja.generated.envoy.config.core.v3.BaseProto.getDescriptor().getEnumTypes().get(1);
  }

  private static final RequestMethod[] VALUES = values();

  public static RequestMethod valueOf(
      com.google.protobuf.Descriptors.EnumValueDescriptor desc) {
    if (desc.getType() != getDescriptor()) {
      throw new java.lang.IllegalArgumentException(
        "EnumValueDescriptor is not for this type.");
    }
    if (desc.getIndex() == -1) {
      return UNRECOGNIZED;
    }
    return VALUES[desc.getIndex()];
  }

  private final int value;

  private RequestMethod(int value) {
    this.value = value;
  }

  // @@protoc_insertion_point(enum_scope:envoy.config.core.v3.RequestMethod)
}

