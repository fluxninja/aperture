package com.fluxninja.generated.aperture.cmd.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Controller describes APIs of the controller from the aperturectl POV
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.54.0)",
    comments = "Source: aperture/cmd/v1/cmd.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class ControllerGrpc {

  private ControllerGrpc() {}

  public static final String SERVICE_NAME = "aperture.cmd.v1.Controller";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse> getListAgentsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListAgents",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse> getListAgentsMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse> getListAgentsMethod;
    if ((getListAgentsMethod = ControllerGrpc.getListAgentsMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getListAgentsMethod = ControllerGrpc.getListAgentsMethod) == null) {
          ControllerGrpc.getListAgentsMethod = getListAgentsMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListAgents"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("ListAgents"))
              .build();
        }
      }
    }
    return getListAgentsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse> getListServicesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListServices",
      requestType = com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse> getListServicesMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest, com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse> getListServicesMethod;
    if ((getListServicesMethod = ControllerGrpc.getListServicesMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getListServicesMethod = ControllerGrpc.getListServicesMethod) == null) {
          ControllerGrpc.getListServicesMethod = getListServicesMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest, com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListServices"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("ListServices"))
              .build();
        }
      }
    }
    return getListServicesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse> getListFlowControlPointsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListFlowControlPoints",
      requestType = com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse> getListFlowControlPointsMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest, com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse> getListFlowControlPointsMethod;
    if ((getListFlowControlPointsMethod = ControllerGrpc.getListFlowControlPointsMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getListFlowControlPointsMethod = ControllerGrpc.getListFlowControlPointsMethod) == null) {
          ControllerGrpc.getListFlowControlPointsMethod = getListFlowControlPointsMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest, com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListFlowControlPoints"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("ListFlowControlPoints"))
              .build();
        }
      }
    }
    return getListFlowControlPointsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse> getListAutoScaleControlPointsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListAutoScaleControlPoints",
      requestType = com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse> getListAutoScaleControlPointsMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest, com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse> getListAutoScaleControlPointsMethod;
    if ((getListAutoScaleControlPointsMethod = ControllerGrpc.getListAutoScaleControlPointsMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getListAutoScaleControlPointsMethod = ControllerGrpc.getListAutoScaleControlPointsMethod) == null) {
          ControllerGrpc.getListAutoScaleControlPointsMethod = getListAutoScaleControlPointsMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest, com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListAutoScaleControlPoints"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("ListAutoScaleControlPoints"))
              .build();
        }
      }
    }
    return getListAutoScaleControlPointsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse> getListDiscoveryEntitiesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListDiscoveryEntities",
      requestType = com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse> getListDiscoveryEntitiesMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest, com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse> getListDiscoveryEntitiesMethod;
    if ((getListDiscoveryEntitiesMethod = ControllerGrpc.getListDiscoveryEntitiesMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getListDiscoveryEntitiesMethod = ControllerGrpc.getListDiscoveryEntitiesMethod) == null) {
          ControllerGrpc.getListDiscoveryEntitiesMethod = getListDiscoveryEntitiesMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest, com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListDiscoveryEntities"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("ListDiscoveryEntities"))
              .build();
        }
      }
    }
    return getListDiscoveryEntitiesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse> getListDiscoveryEntityMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListDiscoveryEntity",
      requestType = com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest,
      com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse> getListDiscoveryEntityMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest, com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse> getListDiscoveryEntityMethod;
    if ((getListDiscoveryEntityMethod = ControllerGrpc.getListDiscoveryEntityMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getListDiscoveryEntityMethod = ControllerGrpc.getListDiscoveryEntityMethod) == null) {
          ControllerGrpc.getListDiscoveryEntityMethod = getListDiscoveryEntityMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest, com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListDiscoveryEntity"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("ListDiscoveryEntity"))
              .build();
        }
      }
    }
    return getListDiscoveryEntityMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest,
      com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse> getPreviewFlowLabelsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PreviewFlowLabels",
      requestType = com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest,
      com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse> getPreviewFlowLabelsMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest, com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse> getPreviewFlowLabelsMethod;
    if ((getPreviewFlowLabelsMethod = ControllerGrpc.getPreviewFlowLabelsMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getPreviewFlowLabelsMethod = ControllerGrpc.getPreviewFlowLabelsMethod) == null) {
          ControllerGrpc.getPreviewFlowLabelsMethod = getPreviewFlowLabelsMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest, com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "PreviewFlowLabels"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("PreviewFlowLabels"))
              .build();
        }
      }
    }
    return getPreviewFlowLabelsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest,
      com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse> getPreviewHTTPRequestsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PreviewHTTPRequests",
      requestType = com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest.class,
      responseType = com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest,
      com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse> getPreviewHTTPRequestsMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest, com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse> getPreviewHTTPRequestsMethod;
    if ((getPreviewHTTPRequestsMethod = ControllerGrpc.getPreviewHTTPRequestsMethod) == null) {
      synchronized (ControllerGrpc.class) {
        if ((getPreviewHTTPRequestsMethod = ControllerGrpc.getPreviewHTTPRequestsMethod) == null) {
          ControllerGrpc.getPreviewHTTPRequestsMethod = getPreviewHTTPRequestsMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest, com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "PreviewHTTPRequests"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerMethodDescriptorSupplier("PreviewHTTPRequests"))
              .build();
        }
      }
    }
    return getPreviewHTTPRequestsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ControllerStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ControllerStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ControllerStub>() {
        @java.lang.Override
        public ControllerStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ControllerStub(channel, callOptions);
        }
      };
    return ControllerStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ControllerBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ControllerBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ControllerBlockingStub>() {
        @java.lang.Override
        public ControllerBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ControllerBlockingStub(channel, callOptions);
        }
      };
    return ControllerBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ControllerFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ControllerFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ControllerFutureStub>() {
        @java.lang.Override
        public ControllerFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ControllerFutureStub(channel, callOptions);
        }
      };
    return ControllerFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Controller describes APIs of the controller from the aperturectl POV
   * </pre>
   */
  public interface AsyncService {

    /**
     */
    default void listAgents(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListAgentsMethod(), responseObserver);
    }

    /**
     */
    default void listServices(com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListServicesMethod(), responseObserver);
    }

    /**
     */
    default void listFlowControlPoints(com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListFlowControlPointsMethod(), responseObserver);
    }

    /**
     */
    default void listAutoScaleControlPoints(com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListAutoScaleControlPointsMethod(), responseObserver);
    }

    /**
     */
    default void listDiscoveryEntities(com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListDiscoveryEntitiesMethod(), responseObserver);
    }

    /**
     */
    default void listDiscoveryEntity(com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListDiscoveryEntityMethod(), responseObserver);
    }

    /**
     * <pre>
     * duplicating a bit preview.v1.FlowPreviewService to keep controller APIs in one place.
     * </pre>
     */
    default void previewFlowLabels(com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPreviewFlowLabelsMethod(), responseObserver);
    }

    /**
     */
    default void previewHTTPRequests(com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPreviewHTTPRequestsMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service Controller.
   * <pre>
   * Controller describes APIs of the controller from the aperturectl POV
   * </pre>
   */
  public static abstract class ControllerImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return ControllerGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service Controller.
   * <pre>
   * Controller describes APIs of the controller from the aperturectl POV
   * </pre>
   */
  public static final class ControllerStub
      extends io.grpc.stub.AbstractAsyncStub<ControllerStub> {
    private ControllerStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControllerStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ControllerStub(channel, callOptions);
    }

    /**
     */
    public void listAgents(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListAgentsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void listServices(com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListServicesMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void listFlowControlPoints(com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListFlowControlPointsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void listAutoScaleControlPoints(com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListAutoScaleControlPointsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void listDiscoveryEntities(com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListDiscoveryEntitiesMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void listDiscoveryEntity(com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListDiscoveryEntityMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * duplicating a bit preview.v1.FlowPreviewService to keep controller APIs in one place.
     * </pre>
     */
    public void previewFlowLabels(com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPreviewFlowLabelsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void previewHTTPRequests(com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPreviewHTTPRequestsMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service Controller.
   * <pre>
   * Controller describes APIs of the controller from the aperturectl POV
   * </pre>
   */
  public static final class ControllerBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<ControllerBlockingStub> {
    private ControllerBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControllerBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ControllerBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse listAgents(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListAgentsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse listServices(com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListServicesMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse listFlowControlPoints(com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListFlowControlPointsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse listAutoScaleControlPoints(com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListAutoScaleControlPointsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse listDiscoveryEntities(com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListDiscoveryEntitiesMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse listDiscoveryEntity(com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListDiscoveryEntityMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * duplicating a bit preview.v1.FlowPreviewService to keep controller APIs in one place.
     * </pre>
     */
    public com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse previewFlowLabels(com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPreviewFlowLabelsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse previewHTTPRequests(com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPreviewHTTPRequestsMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service Controller.
   * <pre>
   * Controller describes APIs of the controller from the aperturectl POV
   * </pre>
   */
  public static final class ControllerFutureStub
      extends io.grpc.stub.AbstractFutureStub<ControllerFutureStub> {
    private ControllerFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControllerFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ControllerFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse> listAgents(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListAgentsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse> listServices(
        com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListServicesMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse> listFlowControlPoints(
        com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListFlowControlPointsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse> listAutoScaleControlPoints(
        com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListAutoScaleControlPointsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse> listDiscoveryEntities(
        com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListDiscoveryEntitiesMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse> listDiscoveryEntity(
        com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListDiscoveryEntityMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * duplicating a bit preview.v1.FlowPreviewService to keep controller APIs in one place.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse> previewFlowLabels(
        com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPreviewFlowLabelsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse> previewHTTPRequests(
        com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPreviewHTTPRequestsMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_LIST_AGENTS = 0;
  private static final int METHODID_LIST_SERVICES = 1;
  private static final int METHODID_LIST_FLOW_CONTROL_POINTS = 2;
  private static final int METHODID_LIST_AUTO_SCALE_CONTROL_POINTS = 3;
  private static final int METHODID_LIST_DISCOVERY_ENTITIES = 4;
  private static final int METHODID_LIST_DISCOVERY_ENTITY = 5;
  private static final int METHODID_PREVIEW_FLOW_LABELS = 6;
  private static final int METHODID_PREVIEW_HTTPREQUESTS = 7;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AsyncService serviceImpl;
    private final int methodId;

    MethodHandlers(AsyncService serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LIST_AGENTS:
          serviceImpl.listAgents((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse>) responseObserver);
          break;
        case METHODID_LIST_SERVICES:
          serviceImpl.listServices((com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse>) responseObserver);
          break;
        case METHODID_LIST_FLOW_CONTROL_POINTS:
          serviceImpl.listFlowControlPoints((com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse>) responseObserver);
          break;
        case METHODID_LIST_AUTO_SCALE_CONTROL_POINTS:
          serviceImpl.listAutoScaleControlPoints((com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse>) responseObserver);
          break;
        case METHODID_LIST_DISCOVERY_ENTITIES:
          serviceImpl.listDiscoveryEntities((com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse>) responseObserver);
          break;
        case METHODID_LIST_DISCOVERY_ENTITY:
          serviceImpl.listDiscoveryEntity((com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse>) responseObserver);
          break;
        case METHODID_PREVIEW_FLOW_LABELS:
          serviceImpl.previewFlowLabels((com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse>) responseObserver);
          break;
        case METHODID_PREVIEW_HTTPREQUESTS:
          serviceImpl.previewHTTPRequests((com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
          getListAgentsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.google.protobuf.Empty,
              com.fluxninja.generated.aperture.cmd.v1.ListAgentsResponse>(
                service, METHODID_LIST_AGENTS)))
        .addMethod(
          getListServicesMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.cmd.v1.ListServicesRequest,
              com.fluxninja.generated.aperture.cmd.v1.ListServicesControllerResponse>(
                service, METHODID_LIST_SERVICES)))
        .addMethod(
          getListFlowControlPointsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsRequest,
              com.fluxninja.generated.aperture.cmd.v1.ListFlowControlPointsControllerResponse>(
                service, METHODID_LIST_FLOW_CONTROL_POINTS)))
        .addMethod(
          getListAutoScaleControlPointsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsRequest,
              com.fluxninja.generated.aperture.cmd.v1.ListAutoScaleControlPointsControllerResponse>(
                service, METHODID_LIST_AUTO_SCALE_CONTROL_POINTS)))
        .addMethod(
          getListDiscoveryEntitiesMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesRequest,
              com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntitiesControllerResponse>(
                service, METHODID_LIST_DISCOVERY_ENTITIES)))
        .addMethod(
          getListDiscoveryEntityMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityRequest,
              com.fluxninja.generated.aperture.cmd.v1.ListDiscoveryEntityAgentResponse>(
                service, METHODID_LIST_DISCOVERY_ENTITY)))
        .addMethod(
          getPreviewFlowLabelsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsRequest,
              com.fluxninja.generated.aperture.cmd.v1.PreviewFlowLabelsControllerResponse>(
                service, METHODID_PREVIEW_FLOW_LABELS)))
        .addMethod(
          getPreviewHTTPRequestsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsRequest,
              com.fluxninja.generated.aperture.cmd.v1.PreviewHTTPRequestsControllerResponse>(
                service, METHODID_PREVIEW_HTTPREQUESTS)))
        .build();
  }

  private static abstract class ControllerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ControllerBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.cmd.v1.CmdProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Controller");
    }
  }

  private static final class ControllerFileDescriptorSupplier
      extends ControllerBaseDescriptorSupplier {
    ControllerFileDescriptorSupplier() {}
  }

  private static final class ControllerMethodDescriptorSupplier
      extends ControllerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ControllerMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (ControllerGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ControllerFileDescriptorSupplier())
              .addMethod(getListAgentsMethod())
              .addMethod(getListServicesMethod())
              .addMethod(getListFlowControlPointsMethod())
              .addMethod(getListAutoScaleControlPointsMethod())
              .addMethod(getListDiscoveryEntitiesMethod())
              .addMethod(getListDiscoveryEntityMethod())
              .addMethod(getPreviewFlowLabelsMethod())
              .addMethod(getPreviewHTTPRequestsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
