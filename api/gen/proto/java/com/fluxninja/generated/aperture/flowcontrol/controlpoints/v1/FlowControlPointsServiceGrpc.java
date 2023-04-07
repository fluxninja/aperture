package com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * grpc service
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/flowcontrol/controlpoints/v1/controlpoints.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class FlowControlPointsServiceGrpc {

  private FlowControlPointsServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.flowcontrol.controlpoints.v1.FlowControlPointsService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints> getGetControlPointsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetControlPoints",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints> getGetControlPointsMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints> getGetControlPointsMethod;
    if ((getGetControlPointsMethod = FlowControlPointsServiceGrpc.getGetControlPointsMethod) == null) {
      synchronized (FlowControlPointsServiceGrpc.class) {
        if ((getGetControlPointsMethod = FlowControlPointsServiceGrpc.getGetControlPointsMethod) == null) {
          FlowControlPointsServiceGrpc.getGetControlPointsMethod = getGetControlPointsMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetControlPoints"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints.getDefaultInstance()))
              .setSchemaDescriptor(new FlowControlPointsServiceMethodDescriptorSupplier("GetControlPoints"))
              .build();
        }
      }
    }
    return getGetControlPointsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FlowControlPointsServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlPointsServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlPointsServiceStub>() {
        @java.lang.Override
        public FlowControlPointsServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlPointsServiceStub(channel, callOptions);
        }
      };
    return FlowControlPointsServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FlowControlPointsServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlPointsServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlPointsServiceBlockingStub>() {
        @java.lang.Override
        public FlowControlPointsServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlPointsServiceBlockingStub(channel, callOptions);
        }
      };
    return FlowControlPointsServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FlowControlPointsServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlPointsServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlPointsServiceFutureStub>() {
        @java.lang.Override
        public FlowControlPointsServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlPointsServiceFutureStub(channel, callOptions);
        }
      };
    return FlowControlPointsServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static abstract class FlowControlPointsServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getControlPoints(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetControlPointsMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetControlPointsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints>(
                  this, METHODID_GET_CONTROL_POINTS)))
          .build();
    }
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static final class FlowControlPointsServiceStub extends io.grpc.stub.AbstractAsyncStub<FlowControlPointsServiceStub> {
    private FlowControlPointsServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlPointsServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlPointsServiceStub(channel, callOptions);
    }

    /**
     */
    public void getControlPoints(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetControlPointsMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static final class FlowControlPointsServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<FlowControlPointsServiceBlockingStub> {
    private FlowControlPointsServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlPointsServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlPointsServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints getControlPoints(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetControlPointsMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static final class FlowControlPointsServiceFutureStub extends io.grpc.stub.AbstractFutureStub<FlowControlPointsServiceFutureStub> {
    private FlowControlPointsServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlPointsServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlPointsServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints> getControlPoints(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetControlPointsMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_CONTROL_POINTS = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final FlowControlPointsServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(FlowControlPointsServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_CONTROL_POINTS:
          serviceImpl.getControlPoints((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.FlowControlPoints>) responseObserver);
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

  private static abstract class FlowControlPointsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FlowControlPointsServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.controlpoints.v1.ControlpointsProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FlowControlPointsService");
    }
  }

  private static final class FlowControlPointsServiceFileDescriptorSupplier
      extends FlowControlPointsServiceBaseDescriptorSupplier {
    FlowControlPointsServiceFileDescriptorSupplier() {}
  }

  private static final class FlowControlPointsServiceMethodDescriptorSupplier
      extends FlowControlPointsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    FlowControlPointsServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (FlowControlPointsServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FlowControlPointsServiceFileDescriptorSupplier())
              .addMethod(getGetControlPointsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
