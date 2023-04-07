package com.fluxninja.generated.aperture.flowcontrol.preview.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/flowcontrol/preview/v1/preview.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class FlowPreviewServiceGrpc {

  private FlowPreviewServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.flowcontrol.preview.v1.FlowPreviewService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest,
      com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse> getPreviewFlowLabelsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PreviewFlowLabels",
      requestType = com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest.class,
      responseType = com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest,
      com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse> getPreviewFlowLabelsMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest, com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse> getPreviewFlowLabelsMethod;
    if ((getPreviewFlowLabelsMethod = FlowPreviewServiceGrpc.getPreviewFlowLabelsMethod) == null) {
      synchronized (FlowPreviewServiceGrpc.class) {
        if ((getPreviewFlowLabelsMethod = FlowPreviewServiceGrpc.getPreviewFlowLabelsMethod) == null) {
          FlowPreviewServiceGrpc.getPreviewFlowLabelsMethod = getPreviewFlowLabelsMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest, com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "PreviewFlowLabels"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlowPreviewServiceMethodDescriptorSupplier("PreviewFlowLabels"))
              .build();
        }
      }
    }
    return getPreviewFlowLabelsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest,
      com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse> getPreviewHTTPRequestsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PreviewHTTPRequests",
      requestType = com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest.class,
      responseType = com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest,
      com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse> getPreviewHTTPRequestsMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest, com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse> getPreviewHTTPRequestsMethod;
    if ((getPreviewHTTPRequestsMethod = FlowPreviewServiceGrpc.getPreviewHTTPRequestsMethod) == null) {
      synchronized (FlowPreviewServiceGrpc.class) {
        if ((getPreviewHTTPRequestsMethod = FlowPreviewServiceGrpc.getPreviewHTTPRequestsMethod) == null) {
          FlowPreviewServiceGrpc.getPreviewHTTPRequestsMethod = getPreviewHTTPRequestsMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest, com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "PreviewHTTPRequests"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlowPreviewServiceMethodDescriptorSupplier("PreviewHTTPRequests"))
              .build();
        }
      }
    }
    return getPreviewHTTPRequestsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FlowPreviewServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowPreviewServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowPreviewServiceStub>() {
        @java.lang.Override
        public FlowPreviewServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowPreviewServiceStub(channel, callOptions);
        }
      };
    return FlowPreviewServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FlowPreviewServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowPreviewServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowPreviewServiceBlockingStub>() {
        @java.lang.Override
        public FlowPreviewServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowPreviewServiceBlockingStub(channel, callOptions);
        }
      };
    return FlowPreviewServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FlowPreviewServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowPreviewServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowPreviewServiceFutureStub>() {
        @java.lang.Override
        public FlowPreviewServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowPreviewServiceFutureStub(channel, callOptions);
        }
      };
    return FlowPreviewServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class FlowPreviewServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void previewFlowLabels(com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPreviewFlowLabelsMethod(), responseObserver);
    }

    /**
     */
    public void previewHTTPRequests(com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPreviewHTTPRequestsMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getPreviewFlowLabelsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest,
                com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse>(
                  this, METHODID_PREVIEW_FLOW_LABELS)))
          .addMethod(
            getPreviewHTTPRequestsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest,
                com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse>(
                  this, METHODID_PREVIEW_HTTPREQUESTS)))
          .build();
    }
  }

  /**
   */
  public static final class FlowPreviewServiceStub extends io.grpc.stub.AbstractAsyncStub<FlowPreviewServiceStub> {
    private FlowPreviewServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowPreviewServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowPreviewServiceStub(channel, callOptions);
    }

    /**
     */
    public void previewFlowLabels(com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPreviewFlowLabelsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void previewHTTPRequests(com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPreviewHTTPRequestsMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class FlowPreviewServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<FlowPreviewServiceBlockingStub> {
    private FlowPreviewServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowPreviewServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowPreviewServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse previewFlowLabels(com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPreviewFlowLabelsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse previewHTTPRequests(com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPreviewHTTPRequestsMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class FlowPreviewServiceFutureStub extends io.grpc.stub.AbstractFutureStub<FlowPreviewServiceFutureStub> {
    private FlowPreviewServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowPreviewServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowPreviewServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse> previewFlowLabels(
        com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPreviewFlowLabelsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse> previewHTTPRequests(
        com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPreviewHTTPRequestsMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_PREVIEW_FLOW_LABELS = 0;
  private static final int METHODID_PREVIEW_HTTPREQUESTS = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final FlowPreviewServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(FlowPreviewServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_PREVIEW_FLOW_LABELS:
          serviceImpl.previewFlowLabels((com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewFlowLabelsResponse>) responseObserver);
          break;
        case METHODID_PREVIEW_HTTPREQUESTS:
          serviceImpl.previewHTTPRequests((com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewHTTPRequestsResponse>) responseObserver);
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

  private static abstract class FlowPreviewServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FlowPreviewServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.preview.v1.PreviewProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FlowPreviewService");
    }
  }

  private static final class FlowPreviewServiceFileDescriptorSupplier
      extends FlowPreviewServiceBaseDescriptorSupplier {
    FlowPreviewServiceFileDescriptorSupplier() {}
  }

  private static final class FlowPreviewServiceMethodDescriptorSupplier
      extends FlowPreviewServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    FlowPreviewServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (FlowPreviewServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FlowPreviewServiceFileDescriptorSupplier())
              .addMethod(getPreviewFlowLabelsMethod())
              .addMethod(getPreviewHTTPRequestsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
