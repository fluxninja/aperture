package com.fluxninja.generated.aperture.entitycache.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * EntityCacheService is used to query EntityCache.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/entitycache/v1/entitycache.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class EntityCacheServiceGrpc {

  private EntityCacheServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.entitycache.v1.EntityCacheService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.entitycache.v1.EntityCache> getGetEntityCacheMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetEntityCache",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.entitycache.v1.EntityCache.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.entitycache.v1.EntityCache> getGetEntityCacheMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.entitycache.v1.EntityCache> getGetEntityCacheMethod;
    if ((getGetEntityCacheMethod = EntityCacheServiceGrpc.getGetEntityCacheMethod) == null) {
      synchronized (EntityCacheServiceGrpc.class) {
        if ((getGetEntityCacheMethod = EntityCacheServiceGrpc.getGetEntityCacheMethod) == null) {
          EntityCacheServiceGrpc.getGetEntityCacheMethod = getGetEntityCacheMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.entitycache.v1.EntityCache>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetEntityCache"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.entitycache.v1.EntityCache.getDefaultInstance()))
              .setSchemaDescriptor(new EntityCacheServiceMethodDescriptorSupplier("GetEntityCache"))
              .build();
        }
      }
    }
    return getGetEntityCacheMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest,
      com.fluxninja.generated.aperture.entitycache.v1.Entity> getGetEntityByIPAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetEntityByIPAddress",
      requestType = com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest.class,
      responseType = com.fluxninja.generated.aperture.entitycache.v1.Entity.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest,
      com.fluxninja.generated.aperture.entitycache.v1.Entity> getGetEntityByIPAddressMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest, com.fluxninja.generated.aperture.entitycache.v1.Entity> getGetEntityByIPAddressMethod;
    if ((getGetEntityByIPAddressMethod = EntityCacheServiceGrpc.getGetEntityByIPAddressMethod) == null) {
      synchronized (EntityCacheServiceGrpc.class) {
        if ((getGetEntityByIPAddressMethod = EntityCacheServiceGrpc.getGetEntityByIPAddressMethod) == null) {
          EntityCacheServiceGrpc.getGetEntityByIPAddressMethod = getGetEntityByIPAddressMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest, com.fluxninja.generated.aperture.entitycache.v1.Entity>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetEntityByIPAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.entitycache.v1.Entity.getDefaultInstance()))
              .setSchemaDescriptor(new EntityCacheServiceMethodDescriptorSupplier("GetEntityByIPAddress"))
              .build();
        }
      }
    }
    return getGetEntityByIPAddressMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest,
      com.fluxninja.generated.aperture.entitycache.v1.Entity> getGetEntityByNameMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetEntityByName",
      requestType = com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest.class,
      responseType = com.fluxninja.generated.aperture.entitycache.v1.Entity.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest,
      com.fluxninja.generated.aperture.entitycache.v1.Entity> getGetEntityByNameMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest, com.fluxninja.generated.aperture.entitycache.v1.Entity> getGetEntityByNameMethod;
    if ((getGetEntityByNameMethod = EntityCacheServiceGrpc.getGetEntityByNameMethod) == null) {
      synchronized (EntityCacheServiceGrpc.class) {
        if ((getGetEntityByNameMethod = EntityCacheServiceGrpc.getGetEntityByNameMethod) == null) {
          EntityCacheServiceGrpc.getGetEntityByNameMethod = getGetEntityByNameMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest, com.fluxninja.generated.aperture.entitycache.v1.Entity>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetEntityByName"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.entitycache.v1.Entity.getDefaultInstance()))
              .setSchemaDescriptor(new EntityCacheServiceMethodDescriptorSupplier("GetEntityByName"))
              .build();
        }
      }
    }
    return getGetEntityByNameMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static EntityCacheServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EntityCacheServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EntityCacheServiceStub>() {
        @java.lang.Override
        public EntityCacheServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EntityCacheServiceStub(channel, callOptions);
        }
      };
    return EntityCacheServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static EntityCacheServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EntityCacheServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EntityCacheServiceBlockingStub>() {
        @java.lang.Override
        public EntityCacheServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EntityCacheServiceBlockingStub(channel, callOptions);
        }
      };
    return EntityCacheServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static EntityCacheServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EntityCacheServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EntityCacheServiceFutureStub>() {
        @java.lang.Override
        public EntityCacheServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EntityCacheServiceFutureStub(channel, callOptions);
        }
      };
    return EntityCacheServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * EntityCacheService is used to query EntityCache.
   * </pre>
   */
  public static abstract class EntityCacheServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getEntityCache(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.EntityCache> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetEntityCacheMethod(), responseObserver);
    }

    /**
     */
    public void getEntityByIPAddress(com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.Entity> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetEntityByIPAddressMethod(), responseObserver);
    }

    /**
     */
    public void getEntityByName(com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.Entity> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetEntityByNameMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetEntityCacheMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.entitycache.v1.EntityCache>(
                  this, METHODID_GET_ENTITY_CACHE)))
          .addMethod(
            getGetEntityByIPAddressMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest,
                com.fluxninja.generated.aperture.entitycache.v1.Entity>(
                  this, METHODID_GET_ENTITY_BY_IPADDRESS)))
          .addMethod(
            getGetEntityByNameMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest,
                com.fluxninja.generated.aperture.entitycache.v1.Entity>(
                  this, METHODID_GET_ENTITY_BY_NAME)))
          .build();
    }
  }

  /**
   * <pre>
   * EntityCacheService is used to query EntityCache.
   * </pre>
   */
  public static final class EntityCacheServiceStub extends io.grpc.stub.AbstractAsyncStub<EntityCacheServiceStub> {
    private EntityCacheServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EntityCacheServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EntityCacheServiceStub(channel, callOptions);
    }

    /**
     */
    public void getEntityCache(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.EntityCache> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetEntityCacheMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getEntityByIPAddress(com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.Entity> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetEntityByIPAddressMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getEntityByName(com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.Entity> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetEntityByNameMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * EntityCacheService is used to query EntityCache.
   * </pre>
   */
  public static final class EntityCacheServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<EntityCacheServiceBlockingStub> {
    private EntityCacheServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EntityCacheServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EntityCacheServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.entitycache.v1.EntityCache getEntityCache(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetEntityCacheMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.entitycache.v1.Entity getEntityByIPAddress(com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetEntityByIPAddressMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.entitycache.v1.Entity getEntityByName(com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetEntityByNameMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * EntityCacheService is used to query EntityCache.
   * </pre>
   */
  public static final class EntityCacheServiceFutureStub extends io.grpc.stub.AbstractFutureStub<EntityCacheServiceFutureStub> {
    private EntityCacheServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EntityCacheServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EntityCacheServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.entitycache.v1.EntityCache> getEntityCache(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetEntityCacheMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.entitycache.v1.Entity> getEntityByIPAddress(
        com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetEntityByIPAddressMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.entitycache.v1.Entity> getEntityByName(
        com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetEntityByNameMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_ENTITY_CACHE = 0;
  private static final int METHODID_GET_ENTITY_BY_IPADDRESS = 1;
  private static final int METHODID_GET_ENTITY_BY_NAME = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final EntityCacheServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(EntityCacheServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_ENTITY_CACHE:
          serviceImpl.getEntityCache((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.EntityCache>) responseObserver);
          break;
        case METHODID_GET_ENTITY_BY_IPADDRESS:
          serviceImpl.getEntityByIPAddress((com.fluxninja.generated.aperture.entitycache.v1.GetEntityByIPAddressRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.Entity>) responseObserver);
          break;
        case METHODID_GET_ENTITY_BY_NAME:
          serviceImpl.getEntityByName((com.fluxninja.generated.aperture.entitycache.v1.GetEntityByNameRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.entitycache.v1.Entity>) responseObserver);
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

  private static abstract class EntityCacheServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    EntityCacheServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.entitycache.v1.EntitycacheProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("EntityCacheService");
    }
  }

  private static final class EntityCacheServiceFileDescriptorSupplier
      extends EntityCacheServiceBaseDescriptorSupplier {
    EntityCacheServiceFileDescriptorSupplier() {}
  }

  private static final class EntityCacheServiceMethodDescriptorSupplier
      extends EntityCacheServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    EntityCacheServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (EntityCacheServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new EntityCacheServiceFileDescriptorSupplier())
              .addMethod(getGetEntityCacheMethod())
              .addMethod(getGetEntityByIPAddressMethod())
              .addMethod(getGetEntityByNameMethod())
              .build();
        }
      }
    }
    return result;
  }
}
