package com.fluxninja.generated.aperture.flowcontrol.check.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/flowcontrol/check/v1/check.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class AWSGatewayFlowControlServiceGrpc {

  private AWSGatewayFlowControlServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.flowcontrol.check.v1.AWSGatewayFlowControlService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest,
      com.fluxninja.generated.google.api.HttpBody> getAWSGatewayCheckMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "AWSGatewayCheck",
      requestType = com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest.class,
      responseType = com.fluxninja.generated.google.api.HttpBody.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest,
      com.fluxninja.generated.google.api.HttpBody> getAWSGatewayCheckMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest, com.fluxninja.generated.google.api.HttpBody> getAWSGatewayCheckMethod;
    if ((getAWSGatewayCheckMethod = AWSGatewayFlowControlServiceGrpc.getAWSGatewayCheckMethod) == null) {
      synchronized (AWSGatewayFlowControlServiceGrpc.class) {
        if ((getAWSGatewayCheckMethod = AWSGatewayFlowControlServiceGrpc.getAWSGatewayCheckMethod) == null) {
          AWSGatewayFlowControlServiceGrpc.getAWSGatewayCheckMethod = getAWSGatewayCheckMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest, com.fluxninja.generated.google.api.HttpBody>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "AWSGatewayCheck"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.google.api.HttpBody.getDefaultInstance()))
              .setSchemaDescriptor(new AWSGatewayFlowControlServiceMethodDescriptorSupplier("AWSGatewayCheck"))
              .build();
        }
      }
    }
    return getAWSGatewayCheckMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AWSGatewayFlowControlServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AWSGatewayFlowControlServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AWSGatewayFlowControlServiceStub>() {
        @java.lang.Override
        public AWSGatewayFlowControlServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AWSGatewayFlowControlServiceStub(channel, callOptions);
        }
      };
    return AWSGatewayFlowControlServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AWSGatewayFlowControlServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AWSGatewayFlowControlServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AWSGatewayFlowControlServiceBlockingStub>() {
        @java.lang.Override
        public AWSGatewayFlowControlServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AWSGatewayFlowControlServiceBlockingStub(channel, callOptions);
        }
      };
    return AWSGatewayFlowControlServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AWSGatewayFlowControlServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AWSGatewayFlowControlServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AWSGatewayFlowControlServiceFutureStub>() {
        @java.lang.Override
        public AWSGatewayFlowControlServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AWSGatewayFlowControlServiceFutureStub(channel, callOptions);
        }
      };
    return AWSGatewayFlowControlServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class AWSGatewayFlowControlServiceImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * AWSGatewayCheck .
     * </pre>
     */
    public void aWSGatewayCheck(com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.google.api.HttpBody> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getAWSGatewayCheckMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getAWSGatewayCheckMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest,
                com.fluxninja.generated.google.api.HttpBody>(
                  this, METHODID_AWSGATEWAY_CHECK)))
          .build();
    }
  }

  /**
   */
  public static final class AWSGatewayFlowControlServiceStub extends io.grpc.stub.AbstractAsyncStub<AWSGatewayFlowControlServiceStub> {
    private AWSGatewayFlowControlServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AWSGatewayFlowControlServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AWSGatewayFlowControlServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * AWSGatewayCheck .
     * </pre>
     */
    public void aWSGatewayCheck(com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.google.api.HttpBody> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getAWSGatewayCheckMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class AWSGatewayFlowControlServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<AWSGatewayFlowControlServiceBlockingStub> {
    private AWSGatewayFlowControlServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AWSGatewayFlowControlServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AWSGatewayFlowControlServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * AWSGatewayCheck .
     * </pre>
     */
    public com.fluxninja.generated.google.api.HttpBody aWSGatewayCheck(com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getAWSGatewayCheckMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class AWSGatewayFlowControlServiceFutureStub extends io.grpc.stub.AbstractFutureStub<AWSGatewayFlowControlServiceFutureStub> {
    private AWSGatewayFlowControlServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AWSGatewayFlowControlServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AWSGatewayFlowControlServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * AWSGatewayCheck .
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.google.api.HttpBody> aWSGatewayCheck(
        com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getAWSGatewayCheckMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_AWSGATEWAY_CHECK = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AWSGatewayFlowControlServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AWSGatewayFlowControlServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_AWSGATEWAY_CHECK:
          serviceImpl.aWSGatewayCheck((com.fluxninja.generated.aperture.flowcontrol.check.v1.AWSGatewayCheckRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.google.api.HttpBody>) responseObserver);
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

  private static abstract class AWSGatewayFlowControlServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AWSGatewayFlowControlServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("AWSGatewayFlowControlService");
    }
  }

  private static final class AWSGatewayFlowControlServiceFileDescriptorSupplier
      extends AWSGatewayFlowControlServiceBaseDescriptorSupplier {
    AWSGatewayFlowControlServiceFileDescriptorSupplier() {}
  }

  private static final class AWSGatewayFlowControlServiceMethodDescriptorSupplier
      extends AWSGatewayFlowControlServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AWSGatewayFlowControlServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (AWSGatewayFlowControlServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AWSGatewayFlowControlServiceFileDescriptorSupplier())
              .addMethod(getAWSGatewayCheckMethod())
              .build();
        }
      }
    }
    return result;
  }
}
