package com.fluxninja.generated.aperture.flowcontrol.check.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * FlowControlService is used to perform Flow Control operations.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.59.0)",
    comments = "Source: aperture/flowcontrol/check/v1/check.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class FlowControlServiceGrpc {

  private FlowControlServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "aperture.flowcontrol.check.v1.FlowControlService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest,
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse> getCheckMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Check",
      requestType = com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest.class,
      responseType = com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest,
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse> getCheckMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest, com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse> getCheckMethod;
    if ((getCheckMethod = FlowControlServiceGrpc.getCheckMethod) == null) {
      synchronized (FlowControlServiceGrpc.class) {
        if ((getCheckMethod = FlowControlServiceGrpc.getCheckMethod) == null) {
          FlowControlServiceGrpc.getCheckMethod = getCheckMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest, com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Check"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlowControlServiceMethodDescriptorSupplier("Check"))
              .build();
        }
      }
    }
    return getCheckMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest,
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse> getCacheUpsertMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CacheUpsert",
      requestType = com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest.class,
      responseType = com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest,
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse> getCacheUpsertMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse> getCacheUpsertMethod;
    if ((getCacheUpsertMethod = FlowControlServiceGrpc.getCacheUpsertMethod) == null) {
      synchronized (FlowControlServiceGrpc.class) {
        if ((getCacheUpsertMethod = FlowControlServiceGrpc.getCacheUpsertMethod) == null) {
          FlowControlServiceGrpc.getCacheUpsertMethod = getCacheUpsertMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CacheUpsert"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlowControlServiceMethodDescriptorSupplier("CacheUpsert"))
              .build();
        }
      }
    }
    return getCacheUpsertMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest,
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse> getCacheDeleteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CacheDelete",
      requestType = com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest.class,
      responseType = com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest,
      com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse> getCacheDeleteMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse> getCacheDeleteMethod;
    if ((getCacheDeleteMethod = FlowControlServiceGrpc.getCacheDeleteMethod) == null) {
      synchronized (FlowControlServiceGrpc.class) {
        if ((getCacheDeleteMethod = FlowControlServiceGrpc.getCacheDeleteMethod) == null) {
          FlowControlServiceGrpc.getCacheDeleteMethod = getCacheDeleteMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest, com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CacheDelete"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlowControlServiceMethodDescriptorSupplier("CacheDelete"))
              .build();
        }
      }
    }
    return getCacheDeleteMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FlowControlServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceStub>() {
        @java.lang.Override
        public FlowControlServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlServiceStub(channel, callOptions);
        }
      };
    return FlowControlServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FlowControlServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceBlockingStub>() {
        @java.lang.Override
        public FlowControlServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlServiceBlockingStub(channel, callOptions);
        }
      };
    return FlowControlServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FlowControlServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceFutureStub>() {
        @java.lang.Override
        public FlowControlServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlServiceFutureStub(channel, callOptions);
        }
      };
    return FlowControlServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * FlowControlService is used to perform Flow Control operations.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
     * </pre>
     */
    default void check(com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCheckMethod(), responseObserver);
    }

    /**
     */
    default void cacheUpsert(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCacheUpsertMethod(), responseObserver);
    }

    /**
     */
    default void cacheDelete(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCacheDeleteMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service FlowControlService.
   * <pre>
   * FlowControlService is used to perform Flow Control operations.
   * </pre>
   */
  public static abstract class FlowControlServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return FlowControlServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service FlowControlService.
   * <pre>
   * FlowControlService is used to perform Flow Control operations.
   * </pre>
   */
  public static final class FlowControlServiceStub
      extends io.grpc.stub.AbstractAsyncStub<FlowControlServiceStub> {
    private FlowControlServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
     * </pre>
     */
    public void check(com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCheckMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void cacheUpsert(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCacheUpsertMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void cacheDelete(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCacheDeleteMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service FlowControlService.
   * <pre>
   * FlowControlService is used to perform Flow Control operations.
   * </pre>
   */
  public static final class FlowControlServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<FlowControlServiceBlockingStub> {
    private FlowControlServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
     * </pre>
     */
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse check(com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCheckMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse cacheUpsert(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCacheUpsertMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse cacheDelete(com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCacheDeleteMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service FlowControlService.
   * <pre>
   * FlowControlService is used to perform Flow Control operations.
   * </pre>
   */
  public static final class FlowControlServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<FlowControlServiceFutureStub> {
    private FlowControlServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse> check(
        com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCheckMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse> cacheUpsert(
        com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCacheUpsertMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse> cacheDelete(
        com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCacheDeleteMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CHECK = 0;
  private static final int METHODID_CACHE_UPSERT = 1;
  private static final int METHODID_CACHE_DELETE = 2;

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
        case METHODID_CHECK:
          serviceImpl.check((com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse>) responseObserver);
          break;
        case METHODID_CACHE_UPSERT:
          serviceImpl.cacheUpsert((com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse>) responseObserver);
          break;
        case METHODID_CACHE_DELETE:
          serviceImpl.cacheDelete((com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse>) responseObserver);
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
          getCheckMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest,
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse>(
                service, METHODID_CHECK)))
        .addMethod(
          getCacheUpsertMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertRequest,
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheUpsertResponse>(
                service, METHODID_CACHE_UPSERT)))
        .addMethod(
          getCacheDeleteMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteRequest,
              com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheDeleteResponse>(
                service, METHODID_CACHE_DELETE)))
        .build();
  }

  private static abstract class FlowControlServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FlowControlServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FlowControlService");
    }
  }

  private static final class FlowControlServiceFileDescriptorSupplier
      extends FlowControlServiceBaseDescriptorSupplier {
    FlowControlServiceFileDescriptorSupplier() {}
  }

  private static final class FlowControlServiceMethodDescriptorSupplier
      extends FlowControlServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    FlowControlServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (FlowControlServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FlowControlServiceFileDescriptorSupplier())
              .addMethod(getCheckMethod())
              .addMethod(getCacheUpsertMethod())
              .addMethod(getCacheDeleteMethod())
              .build();
        }
      }
    }
    return result;
  }
}
