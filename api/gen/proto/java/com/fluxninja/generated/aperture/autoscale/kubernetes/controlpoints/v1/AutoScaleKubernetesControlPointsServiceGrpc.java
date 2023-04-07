package com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * grpc service
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/autoscale/kubernetes/controlpoints/v1/controlpoints.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class AutoScaleKubernetesControlPointsServiceGrpc {

  private AutoScaleKubernetesControlPointsServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPointsService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints> getGetControlPointsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetControlPoints",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints> getGetControlPointsMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints> getGetControlPointsMethod;
    if ((getGetControlPointsMethod = AutoScaleKubernetesControlPointsServiceGrpc.getGetControlPointsMethod) == null) {
      synchronized (AutoScaleKubernetesControlPointsServiceGrpc.class) {
        if ((getGetControlPointsMethod = AutoScaleKubernetesControlPointsServiceGrpc.getGetControlPointsMethod) == null) {
          AutoScaleKubernetesControlPointsServiceGrpc.getGetControlPointsMethod = getGetControlPointsMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetControlPoints"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints.getDefaultInstance()))
              .setSchemaDescriptor(new AutoScaleKubernetesControlPointsServiceMethodDescriptorSupplier("GetControlPoints"))
              .build();
        }
      }
    }
    return getGetControlPointsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AutoScaleKubernetesControlPointsServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AutoScaleKubernetesControlPointsServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AutoScaleKubernetesControlPointsServiceStub>() {
        @java.lang.Override
        public AutoScaleKubernetesControlPointsServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AutoScaleKubernetesControlPointsServiceStub(channel, callOptions);
        }
      };
    return AutoScaleKubernetesControlPointsServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AutoScaleKubernetesControlPointsServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AutoScaleKubernetesControlPointsServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AutoScaleKubernetesControlPointsServiceBlockingStub>() {
        @java.lang.Override
        public AutoScaleKubernetesControlPointsServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AutoScaleKubernetesControlPointsServiceBlockingStub(channel, callOptions);
        }
      };
    return AutoScaleKubernetesControlPointsServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AutoScaleKubernetesControlPointsServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AutoScaleKubernetesControlPointsServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AutoScaleKubernetesControlPointsServiceFutureStub>() {
        @java.lang.Override
        public AutoScaleKubernetesControlPointsServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AutoScaleKubernetesControlPointsServiceFutureStub(channel, callOptions);
        }
      };
    return AutoScaleKubernetesControlPointsServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static abstract class AutoScaleKubernetesControlPointsServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getControlPoints(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetControlPointsMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetControlPointsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints>(
                  this, METHODID_GET_CONTROL_POINTS)))
          .build();
    }
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static final class AutoScaleKubernetesControlPointsServiceStub extends io.grpc.stub.AbstractAsyncStub<AutoScaleKubernetesControlPointsServiceStub> {
    private AutoScaleKubernetesControlPointsServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AutoScaleKubernetesControlPointsServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AutoScaleKubernetesControlPointsServiceStub(channel, callOptions);
    }

    /**
     */
    public void getControlPoints(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetControlPointsMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static final class AutoScaleKubernetesControlPointsServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<AutoScaleKubernetesControlPointsServiceBlockingStub> {
    private AutoScaleKubernetesControlPointsServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AutoScaleKubernetesControlPointsServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AutoScaleKubernetesControlPointsServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints getControlPoints(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetControlPointsMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * grpc service
   * </pre>
   */
  public static final class AutoScaleKubernetesControlPointsServiceFutureStub extends io.grpc.stub.AbstractFutureStub<AutoScaleKubernetesControlPointsServiceFutureStub> {
    private AutoScaleKubernetesControlPointsServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AutoScaleKubernetesControlPointsServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AutoScaleKubernetesControlPointsServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints> getControlPoints(
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
    private final AutoScaleKubernetesControlPointsServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AutoScaleKubernetesControlPointsServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_CONTROL_POINTS:
          serviceImpl.getControlPoints((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.AutoScaleKubernetesControlPoints>) responseObserver);
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

  private static abstract class AutoScaleKubernetesControlPointsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AutoScaleKubernetesControlPointsServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.autoscale.kubernetes.controlpoints.v1.ControlpointsProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("AutoScaleKubernetesControlPointsService");
    }
  }

  private static final class AutoScaleKubernetesControlPointsServiceFileDescriptorSupplier
      extends AutoScaleKubernetesControlPointsServiceBaseDescriptorSupplier {
    AutoScaleKubernetesControlPointsServiceFileDescriptorSupplier() {}
  }

  private static final class AutoScaleKubernetesControlPointsServiceMethodDescriptorSupplier
      extends AutoScaleKubernetesControlPointsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AutoScaleKubernetesControlPointsServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (AutoScaleKubernetesControlPointsServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AutoScaleKubernetesControlPointsServiceFileDescriptorSupplier())
              .addMethod(getGetControlPointsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
