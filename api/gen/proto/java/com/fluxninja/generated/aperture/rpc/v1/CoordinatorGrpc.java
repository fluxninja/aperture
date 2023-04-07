package com.fluxninja.generated.aperture.rpc.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/rpc/v1/rpc.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class CoordinatorGrpc {

  private CoordinatorGrpc() {}

  public static final String SERVICE_NAME = "aperture.rpc.v1.Coordinator";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.rpc.v1.ClientToServer,
      com.fluxninja.generated.aperture.rpc.v1.ServerToClient> getConnectMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Connect",
      requestType = com.fluxninja.generated.aperture.rpc.v1.ClientToServer.class,
      responseType = com.fluxninja.generated.aperture.rpc.v1.ServerToClient.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.rpc.v1.ClientToServer,
      com.fluxninja.generated.aperture.rpc.v1.ServerToClient> getConnectMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.rpc.v1.ClientToServer, com.fluxninja.generated.aperture.rpc.v1.ServerToClient> getConnectMethod;
    if ((getConnectMethod = CoordinatorGrpc.getConnectMethod) == null) {
      synchronized (CoordinatorGrpc.class) {
        if ((getConnectMethod = CoordinatorGrpc.getConnectMethod) == null) {
          CoordinatorGrpc.getConnectMethod = getConnectMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.rpc.v1.ClientToServer, com.fluxninja.generated.aperture.rpc.v1.ServerToClient>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Connect"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.rpc.v1.ClientToServer.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.rpc.v1.ServerToClient.getDefaultInstance()))
              .setSchemaDescriptor(new CoordinatorMethodDescriptorSupplier("Connect"))
              .build();
        }
      }
    }
    return getConnectMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static CoordinatorStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CoordinatorStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CoordinatorStub>() {
        @java.lang.Override
        public CoordinatorStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CoordinatorStub(channel, callOptions);
        }
      };
    return CoordinatorStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static CoordinatorBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CoordinatorBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CoordinatorBlockingStub>() {
        @java.lang.Override
        public CoordinatorBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CoordinatorBlockingStub(channel, callOptions);
        }
      };
    return CoordinatorBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static CoordinatorFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CoordinatorFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CoordinatorFutureStub>() {
        @java.lang.Override
        public CoordinatorFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CoordinatorFutureStub(channel, callOptions);
        }
      };
    return CoordinatorFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class CoordinatorImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.rpc.v1.ClientToServer> connect(
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.rpc.v1.ServerToClient> responseObserver) {
      return io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall(getConnectMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getConnectMethod(),
            io.grpc.stub.ServerCalls.asyncBidiStreamingCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.rpc.v1.ClientToServer,
                com.fluxninja.generated.aperture.rpc.v1.ServerToClient>(
                  this, METHODID_CONNECT)))
          .build();
    }
  }

  /**
   */
  public static final class CoordinatorStub extends io.grpc.stub.AbstractAsyncStub<CoordinatorStub> {
    private CoordinatorStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CoordinatorStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CoordinatorStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.rpc.v1.ClientToServer> connect(
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.rpc.v1.ServerToClient> responseObserver) {
      return io.grpc.stub.ClientCalls.asyncBidiStreamingCall(
          getChannel().newCall(getConnectMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   */
  public static final class CoordinatorBlockingStub extends io.grpc.stub.AbstractBlockingStub<CoordinatorBlockingStub> {
    private CoordinatorBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CoordinatorBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CoordinatorBlockingStub(channel, callOptions);
    }
  }

  /**
   */
  public static final class CoordinatorFutureStub extends io.grpc.stub.AbstractFutureStub<CoordinatorFutureStub> {
    private CoordinatorFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CoordinatorFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CoordinatorFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_CONNECT = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final CoordinatorImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(CoordinatorImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CONNECT:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.connect(
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.rpc.v1.ServerToClient>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class CoordinatorBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    CoordinatorBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.rpc.v1.RpcProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Coordinator");
    }
  }

  private static final class CoordinatorFileDescriptorSupplier
      extends CoordinatorBaseDescriptorSupplier {
    CoordinatorFileDescriptorSupplier() {}
  }

  private static final class CoordinatorMethodDescriptorSupplier
      extends CoordinatorBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    CoordinatorMethodDescriptorSupplier(String methodName) {
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
      synchronized (CoordinatorGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new CoordinatorFileDescriptorSupplier())
              .addMethod(getConnectMethod())
              .build();
        }
      }
    }
    return result;
  }
}
