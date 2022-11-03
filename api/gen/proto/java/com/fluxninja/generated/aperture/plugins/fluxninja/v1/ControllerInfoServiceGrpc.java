package com.fluxninja.generated.aperture.plugins.fluxninja.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * ControllerInfoService is used to read controllerID to which agent/controller belong.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/plugins/fluxninja/v1/heartbeat.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class ControllerInfoServiceGrpc {

  private ControllerInfoServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.plugins.fluxninja.v1.ControllerInfoService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo> getGetControllerInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetControllerInfo",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo> getGetControllerInfoMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo> getGetControllerInfoMethod;
    if ((getGetControllerInfoMethod = ControllerInfoServiceGrpc.getGetControllerInfoMethod) == null) {
      synchronized (ControllerInfoServiceGrpc.class) {
        if ((getGetControllerInfoMethod = ControllerInfoServiceGrpc.getGetControllerInfoMethod) == null) {
          ControllerInfoServiceGrpc.getGetControllerInfoMethod = getGetControllerInfoMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetControllerInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo.getDefaultInstance()))
              .setSchemaDescriptor(new ControllerInfoServiceMethodDescriptorSupplier("GetControllerInfo"))
              .build();
        }
      }
    }
    return getGetControllerInfoMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ControllerInfoServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ControllerInfoServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ControllerInfoServiceStub>() {
        @java.lang.Override
        public ControllerInfoServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ControllerInfoServiceStub(channel, callOptions);
        }
      };
    return ControllerInfoServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ControllerInfoServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ControllerInfoServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ControllerInfoServiceBlockingStub>() {
        @java.lang.Override
        public ControllerInfoServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ControllerInfoServiceBlockingStub(channel, callOptions);
        }
      };
    return ControllerInfoServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ControllerInfoServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ControllerInfoServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ControllerInfoServiceFutureStub>() {
        @java.lang.Override
        public ControllerInfoServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ControllerInfoServiceFutureStub(channel, callOptions);
        }
      };
    return ControllerInfoServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * ControllerInfoService is used to read controllerID to which agent/controller belong.
   * </pre>
   */
  public static abstract class ControllerInfoServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getControllerInfo(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetControllerInfoMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetControllerInfoMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo>(
                  this, METHODID_GET_CONTROLLER_INFO)))
          .build();
    }
  }

  /**
   * <pre>
   * ControllerInfoService is used to read controllerID to which agent/controller belong.
   * </pre>
   */
  public static final class ControllerInfoServiceStub extends io.grpc.stub.AbstractAsyncStub<ControllerInfoServiceStub> {
    private ControllerInfoServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControllerInfoServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ControllerInfoServiceStub(channel, callOptions);
    }

    /**
     */
    public void getControllerInfo(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetControllerInfoMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * ControllerInfoService is used to read controllerID to which agent/controller belong.
   * </pre>
   */
  public static final class ControllerInfoServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<ControllerInfoServiceBlockingStub> {
    private ControllerInfoServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControllerInfoServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ControllerInfoServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo getControllerInfo(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetControllerInfoMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * ControllerInfoService is used to read controllerID to which agent/controller belong.
   * </pre>
   */
  public static final class ControllerInfoServiceFutureStub extends io.grpc.stub.AbstractFutureStub<ControllerInfoServiceFutureStub> {
    private ControllerInfoServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControllerInfoServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ControllerInfoServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo> getControllerInfo(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetControllerInfoMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_CONTROLLER_INFO = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ControllerInfoServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ControllerInfoServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_CONTROLLER_INFO:
          serviceImpl.getControllerInfo((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ControllerInfo>) responseObserver);
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

  private static abstract class ControllerInfoServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ControllerInfoServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.plugins.fluxninja.v1.HeartbeatProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("ControllerInfoService");
    }
  }

  private static final class ControllerInfoServiceFileDescriptorSupplier
      extends ControllerInfoServiceBaseDescriptorSupplier {
    ControllerInfoServiceFileDescriptorSupplier() {}
  }

  private static final class ControllerInfoServiceMethodDescriptorSupplier
      extends ControllerInfoServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ControllerInfoServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (ControllerInfoServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ControllerInfoServiceFileDescriptorSupplier())
              .addMethod(getGetControllerInfoMethod())
              .build();
        }
      }
    }
    return result;
  }
}
