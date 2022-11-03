package com.fluxninja.generated.aperture.plugins.fluxninja.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * FluxNinjaService is used to receive health and status info from agents.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/plugins/fluxninja/v1/heartbeat.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class FluxNinjaServiceGrpc {

  private FluxNinjaServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.plugins.fluxninja.v1.FluxNinjaService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest,
      com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse> getReportMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Report",
      requestType = com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest.class,
      responseType = com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest,
      com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse> getReportMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest, com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse> getReportMethod;
    if ((getReportMethod = FluxNinjaServiceGrpc.getReportMethod) == null) {
      synchronized (FluxNinjaServiceGrpc.class) {
        if ((getReportMethod = FluxNinjaServiceGrpc.getReportMethod) == null) {
          FluxNinjaServiceGrpc.getReportMethod = getReportMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest, com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Report"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FluxNinjaServiceMethodDescriptorSupplier("Report"))
              .build();
        }
      }
    }
    return getReportMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FluxNinjaServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FluxNinjaServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FluxNinjaServiceStub>() {
        @java.lang.Override
        public FluxNinjaServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FluxNinjaServiceStub(channel, callOptions);
        }
      };
    return FluxNinjaServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FluxNinjaServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FluxNinjaServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FluxNinjaServiceBlockingStub>() {
        @java.lang.Override
        public FluxNinjaServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FluxNinjaServiceBlockingStub(channel, callOptions);
        }
      };
    return FluxNinjaServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FluxNinjaServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FluxNinjaServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FluxNinjaServiceFutureStub>() {
        @java.lang.Override
        public FluxNinjaServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FluxNinjaServiceFutureStub(channel, callOptions);
        }
      };
    return FluxNinjaServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * FluxNinjaService is used to receive health and status info from agents.
   * </pre>
   */
  public static abstract class FluxNinjaServiceImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * Report accepts information about agents' health and applied configurations/policies.
     * </pre>
     */
    public void report(com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getReportMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getReportMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest,
                com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse>(
                  this, METHODID_REPORT)))
          .build();
    }
  }

  /**
   * <pre>
   * FluxNinjaService is used to receive health and status info from agents.
   * </pre>
   */
  public static final class FluxNinjaServiceStub extends io.grpc.stub.AbstractAsyncStub<FluxNinjaServiceStub> {
    private FluxNinjaServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FluxNinjaServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FluxNinjaServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * Report accepts information about agents' health and applied configurations/policies.
     * </pre>
     */
    public void report(com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getReportMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * FluxNinjaService is used to receive health and status info from agents.
   * </pre>
   */
  public static final class FluxNinjaServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<FluxNinjaServiceBlockingStub> {
    private FluxNinjaServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FluxNinjaServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FluxNinjaServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * Report accepts information about agents' health and applied configurations/policies.
     * </pre>
     */
    public com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse report(com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getReportMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * FluxNinjaService is used to receive health and status info from agents.
   * </pre>
   */
  public static final class FluxNinjaServiceFutureStub extends io.grpc.stub.AbstractFutureStub<FluxNinjaServiceFutureStub> {
    private FluxNinjaServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FluxNinjaServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FluxNinjaServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * Report accepts information about agents' health and applied configurations/policies.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse> report(
        com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getReportMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_REPORT = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final FluxNinjaServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(FluxNinjaServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_REPORT:
          serviceImpl.report((com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.plugins.fluxninja.v1.ReportResponse>) responseObserver);
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

  private static abstract class FluxNinjaServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FluxNinjaServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.plugins.fluxninja.v1.HeartbeatProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FluxNinjaService");
    }
  }

  private static final class FluxNinjaServiceFileDescriptorSupplier
      extends FluxNinjaServiceBaseDescriptorSupplier {
    FluxNinjaServiceFileDescriptorSupplier() {}
  }

  private static final class FluxNinjaServiceMethodDescriptorSupplier
      extends FluxNinjaServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    FluxNinjaServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (FluxNinjaServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FluxNinjaServiceFileDescriptorSupplier())
              .addMethod(getReportMethod())
              .build();
        }
      }
    }
    return result;
  }
}
