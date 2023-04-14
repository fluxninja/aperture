package com.fluxninja.generated.aperture.distcache.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * DistCacheService is used to query DistCache.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.54.0)",
    comments = "Source: aperture/distcache/v1/stats.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class DistCacheServiceGrpc {

  private DistCacheServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.distcache.v1.DistCacheService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.distcache.v1.Stats> getGetStatsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetStats",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.distcache.v1.Stats.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.distcache.v1.Stats> getGetStatsMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.distcache.v1.Stats> getGetStatsMethod;
    if ((getGetStatsMethod = DistCacheServiceGrpc.getGetStatsMethod) == null) {
      synchronized (DistCacheServiceGrpc.class) {
        if ((getGetStatsMethod = DistCacheServiceGrpc.getGetStatsMethod) == null) {
          DistCacheServiceGrpc.getGetStatsMethod = getGetStatsMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.distcache.v1.Stats>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetStats"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.distcache.v1.Stats.getDefaultInstance()))
              .setSchemaDescriptor(new DistCacheServiceMethodDescriptorSupplier("GetStats"))
              .build();
        }
      }
    }
    return getGetStatsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static DistCacheServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<DistCacheServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<DistCacheServiceStub>() {
        @java.lang.Override
        public DistCacheServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new DistCacheServiceStub(channel, callOptions);
        }
      };
    return DistCacheServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static DistCacheServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<DistCacheServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<DistCacheServiceBlockingStub>() {
        @java.lang.Override
        public DistCacheServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new DistCacheServiceBlockingStub(channel, callOptions);
        }
      };
    return DistCacheServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static DistCacheServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<DistCacheServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<DistCacheServiceFutureStub>() {
        @java.lang.Override
        public DistCacheServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new DistCacheServiceFutureStub(channel, callOptions);
        }
      };
    return DistCacheServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * DistCacheService is used to query DistCache.
   * </pre>
   */
  public interface AsyncService {

    /**
     */
    default void getStats(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.distcache.v1.Stats> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetStatsMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service DistCacheService.
   * <pre>
   * DistCacheService is used to query DistCache.
   * </pre>
   */
  public static abstract class DistCacheServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return DistCacheServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service DistCacheService.
   * <pre>
   * DistCacheService is used to query DistCache.
   * </pre>
   */
  public static final class DistCacheServiceStub
      extends io.grpc.stub.AbstractAsyncStub<DistCacheServiceStub> {
    private DistCacheServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DistCacheServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new DistCacheServiceStub(channel, callOptions);
    }

    /**
     */
    public void getStats(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.distcache.v1.Stats> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetStatsMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service DistCacheService.
   * <pre>
   * DistCacheService is used to query DistCache.
   * </pre>
   */
  public static final class DistCacheServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<DistCacheServiceBlockingStub> {
    private DistCacheServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DistCacheServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new DistCacheServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.distcache.v1.Stats getStats(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetStatsMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service DistCacheService.
   * <pre>
   * DistCacheService is used to query DistCache.
   * </pre>
   */
  public static final class DistCacheServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<DistCacheServiceFutureStub> {
    private DistCacheServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DistCacheServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new DistCacheServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.distcache.v1.Stats> getStats(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetStatsMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_STATS = 0;

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
        case METHODID_GET_STATS:
          serviceImpl.getStats((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.distcache.v1.Stats>) responseObserver);
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
          getGetStatsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.google.protobuf.Empty,
              com.fluxninja.generated.aperture.distcache.v1.Stats>(
                service, METHODID_GET_STATS)))
        .build();
  }

  private static abstract class DistCacheServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    DistCacheServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.distcache.v1.StatsProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("DistCacheService");
    }
  }

  private static final class DistCacheServiceFileDescriptorSupplier
      extends DistCacheServiceBaseDescriptorSupplier {
    DistCacheServiceFileDescriptorSupplier() {}
  }

  private static final class DistCacheServiceMethodDescriptorSupplier
      extends DistCacheServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    DistCacheServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (DistCacheServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new DistCacheServiceFileDescriptorSupplier())
              .addMethod(getGetStatsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
