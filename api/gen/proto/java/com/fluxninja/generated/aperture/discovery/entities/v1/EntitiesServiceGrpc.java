package com.fluxninja.generated.aperture.discovery.entities.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * EntitiesService is used to query Entities.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/discovery/entities/v1/entities.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class EntitiesServiceGrpc {

  private EntitiesServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.discovery.entities.v1.EntitiesService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.discovery.entities.v1.Entities> getGetEntitiesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetEntities",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.discovery.entities.v1.Entities.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.discovery.entities.v1.Entities> getGetEntitiesMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.discovery.entities.v1.Entities> getGetEntitiesMethod;
    if ((getGetEntitiesMethod = EntitiesServiceGrpc.getGetEntitiesMethod) == null) {
      synchronized (EntitiesServiceGrpc.class) {
        if ((getGetEntitiesMethod = EntitiesServiceGrpc.getGetEntitiesMethod) == null) {
          EntitiesServiceGrpc.getGetEntitiesMethod = getGetEntitiesMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.discovery.entities.v1.Entities>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetEntities"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.discovery.entities.v1.Entities.getDefaultInstance()))
              .setSchemaDescriptor(new EntitiesServiceMethodDescriptorSupplier("GetEntities"))
              .build();
        }
      }
    }
    return getGetEntitiesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest,
      com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getGetEntityByIPAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetEntityByIPAddress",
      requestType = com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.class,
      responseType = com.fluxninja.generated.aperture.discovery.entities.v1.Entity.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest,
      com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getGetEntityByIPAddressMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest, com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getGetEntityByIPAddressMethod;
    if ((getGetEntityByIPAddressMethod = EntitiesServiceGrpc.getGetEntityByIPAddressMethod) == null) {
      synchronized (EntitiesServiceGrpc.class) {
        if ((getGetEntityByIPAddressMethod = EntitiesServiceGrpc.getGetEntityByIPAddressMethod) == null) {
          EntitiesServiceGrpc.getGetEntityByIPAddressMethod = getGetEntityByIPAddressMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest, com.fluxninja.generated.aperture.discovery.entities.v1.Entity>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetEntityByIPAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.discovery.entities.v1.Entity.getDefaultInstance()))
              .setSchemaDescriptor(new EntitiesServiceMethodDescriptorSupplier("GetEntityByIPAddress"))
              .build();
        }
      }
    }
    return getGetEntityByIPAddressMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest,
      com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getGetEntityByNameMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetEntityByName",
      requestType = com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest.class,
      responseType = com.fluxninja.generated.aperture.discovery.entities.v1.Entity.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest,
      com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getGetEntityByNameMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest, com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getGetEntityByNameMethod;
    if ((getGetEntityByNameMethod = EntitiesServiceGrpc.getGetEntityByNameMethod) == null) {
      synchronized (EntitiesServiceGrpc.class) {
        if ((getGetEntityByNameMethod = EntitiesServiceGrpc.getGetEntityByNameMethod) == null) {
          EntitiesServiceGrpc.getGetEntityByNameMethod = getGetEntityByNameMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest, com.fluxninja.generated.aperture.discovery.entities.v1.Entity>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetEntityByName"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.discovery.entities.v1.Entity.getDefaultInstance()))
              .setSchemaDescriptor(new EntitiesServiceMethodDescriptorSupplier("GetEntityByName"))
              .build();
        }
      }
    }
    return getGetEntityByNameMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static EntitiesServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EntitiesServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EntitiesServiceStub>() {
        @java.lang.Override
        public EntitiesServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EntitiesServiceStub(channel, callOptions);
        }
      };
    return EntitiesServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static EntitiesServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EntitiesServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EntitiesServiceBlockingStub>() {
        @java.lang.Override
        public EntitiesServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EntitiesServiceBlockingStub(channel, callOptions);
        }
      };
    return EntitiesServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static EntitiesServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EntitiesServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EntitiesServiceFutureStub>() {
        @java.lang.Override
        public EntitiesServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EntitiesServiceFutureStub(channel, callOptions);
        }
      };
    return EntitiesServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * EntitiesService is used to query Entities.
   * </pre>
   */
  public static abstract class EntitiesServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getEntities(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entities> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetEntitiesMethod(), responseObserver);
    }

    /**
     */
    public void getEntityByIPAddress(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entity> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetEntityByIPAddressMethod(), responseObserver);
    }

    /**
     */
    public void getEntityByName(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entity> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetEntityByNameMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetEntitiesMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.discovery.entities.v1.Entities>(
                  this, METHODID_GET_ENTITIES)))
          .addMethod(
            getGetEntityByIPAddressMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest,
                com.fluxninja.generated.aperture.discovery.entities.v1.Entity>(
                  this, METHODID_GET_ENTITY_BY_IPADDRESS)))
          .addMethod(
            getGetEntityByNameMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest,
                com.fluxninja.generated.aperture.discovery.entities.v1.Entity>(
                  this, METHODID_GET_ENTITY_BY_NAME)))
          .build();
    }
  }

  /**
   * <pre>
   * EntitiesService is used to query Entities.
   * </pre>
   */
  public static final class EntitiesServiceStub extends io.grpc.stub.AbstractAsyncStub<EntitiesServiceStub> {
    private EntitiesServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EntitiesServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EntitiesServiceStub(channel, callOptions);
    }

    /**
     */
    public void getEntities(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entities> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetEntitiesMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getEntityByIPAddress(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entity> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetEntityByIPAddressMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getEntityByName(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entity> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetEntityByNameMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * EntitiesService is used to query Entities.
   * </pre>
   */
  public static final class EntitiesServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<EntitiesServiceBlockingStub> {
    private EntitiesServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EntitiesServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EntitiesServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.discovery.entities.v1.Entities getEntities(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetEntitiesMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.discovery.entities.v1.Entity getEntityByIPAddress(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetEntityByIPAddressMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.discovery.entities.v1.Entity getEntityByName(com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetEntityByNameMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * EntitiesService is used to query Entities.
   * </pre>
   */
  public static final class EntitiesServiceFutureStub extends io.grpc.stub.AbstractFutureStub<EntitiesServiceFutureStub> {
    private EntitiesServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EntitiesServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EntitiesServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.discovery.entities.v1.Entities> getEntities(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetEntitiesMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getEntityByIPAddress(
        com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetEntityByIPAddressMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.discovery.entities.v1.Entity> getEntityByName(
        com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetEntityByNameMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_ENTITIES = 0;
  private static final int METHODID_GET_ENTITY_BY_IPADDRESS = 1;
  private static final int METHODID_GET_ENTITY_BY_NAME = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final EntitiesServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(EntitiesServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_ENTITIES:
          serviceImpl.getEntities((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entities>) responseObserver);
          break;
        case METHODID_GET_ENTITY_BY_IPADDRESS:
          serviceImpl.getEntityByIPAddress((com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByIPAddressRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entity>) responseObserver);
          break;
        case METHODID_GET_ENTITY_BY_NAME:
          serviceImpl.getEntityByName((com.fluxninja.generated.aperture.discovery.entities.v1.GetEntityByNameRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.discovery.entities.v1.Entity>) responseObserver);
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

  private static abstract class EntitiesServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    EntitiesServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.discovery.entities.v1.EntitiesProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("EntitiesService");
    }
  }

  private static final class EntitiesServiceFileDescriptorSupplier
      extends EntitiesServiceBaseDescriptorSupplier {
    EntitiesServiceFileDescriptorSupplier() {}
  }

  private static final class EntitiesServiceMethodDescriptorSupplier
      extends EntitiesServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    EntitiesServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (EntitiesServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new EntitiesServiceFileDescriptorSupplier())
              .addMethod(getGetEntitiesMethod())
              .addMethod(getGetEntityByIPAddressMethod())
              .addMethod(getGetEntityByNameMethod())
              .build();
        }
      }
    }
    return result;
  }
}
