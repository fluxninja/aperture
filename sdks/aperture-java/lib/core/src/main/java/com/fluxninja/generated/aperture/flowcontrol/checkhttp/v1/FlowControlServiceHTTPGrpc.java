package com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.57.1)",
    comments = "Source: aperture/flowcontrol/checkhttp/v1/checkhttp.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class FlowControlServiceHTTPGrpc {

  private FlowControlServiceHTTPGrpc() {}

  public static final java.lang.String SERVICE_NAME = "aperture.flowcontrol.checkhttp.v1.FlowControlServiceHTTP";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest,
      com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse> getCheckHTTPMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CheckHTTP",
      requestType = com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest.class,
      responseType = com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest,
      com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse> getCheckHTTPMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest, com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse> getCheckHTTPMethod;
    if ((getCheckHTTPMethod = FlowControlServiceHTTPGrpc.getCheckHTTPMethod) == null) {
      synchronized (FlowControlServiceHTTPGrpc.class) {
        if ((getCheckHTTPMethod = FlowControlServiceHTTPGrpc.getCheckHTTPMethod) == null) {
          FlowControlServiceHTTPGrpc.getCheckHTTPMethod = getCheckHTTPMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest, com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CheckHTTP"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlowControlServiceHTTPMethodDescriptorSupplier("CheckHTTP"))
              .build();
        }
      }
    }
    return getCheckHTTPMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FlowControlServiceHTTPStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceHTTPStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceHTTPStub>() {
        @java.lang.Override
        public FlowControlServiceHTTPStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlServiceHTTPStub(channel, callOptions);
        }
      };
    return FlowControlServiceHTTPStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FlowControlServiceHTTPBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceHTTPBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceHTTPBlockingStub>() {
        @java.lang.Override
        public FlowControlServiceHTTPBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlServiceHTTPBlockingStub(channel, callOptions);
        }
      };
    return FlowControlServiceHTTPBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FlowControlServiceHTTPFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceHTTPFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlowControlServiceHTTPFutureStub>() {
        @java.lang.Override
        public FlowControlServiceHTTPFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlowControlServiceHTTPFutureStub(channel, callOptions);
        }
      };
    return FlowControlServiceHTTPFutureStub.newStub(factory, channel);
  }

  /**
   */
  public interface AsyncService {

    /**
     */
    default void checkHTTP(com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCheckHTTPMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service FlowControlServiceHTTP.
   */
  public static abstract class FlowControlServiceHTTPImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return FlowControlServiceHTTPGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service FlowControlServiceHTTP.
   */
  public static final class FlowControlServiceHTTPStub
      extends io.grpc.stub.AbstractAsyncStub<FlowControlServiceHTTPStub> {
    private FlowControlServiceHTTPStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlServiceHTTPStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlServiceHTTPStub(channel, callOptions);
    }

    /**
     */
    public void checkHTTP(com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCheckHTTPMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service FlowControlServiceHTTP.
   */
  public static final class FlowControlServiceHTTPBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<FlowControlServiceHTTPBlockingStub> {
    private FlowControlServiceHTTPBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlServiceHTTPBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlServiceHTTPBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse checkHTTP(com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCheckHTTPMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service FlowControlServiceHTTP.
   */
  public static final class FlowControlServiceHTTPFutureStub
      extends io.grpc.stub.AbstractFutureStub<FlowControlServiceHTTPFutureStub> {
    private FlowControlServiceHTTPFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlowControlServiceHTTPFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlowControlServiceHTTPFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse> checkHTTP(
        com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCheckHTTPMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CHECK_HTTP = 0;

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
        case METHODID_CHECK_HTTP:
          serviceImpl.checkHTTP((com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse>) responseObserver);
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
          getCheckHTTPMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest,
              com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse>(
                service, METHODID_CHECK_HTTP)))
        .build();
  }

  private static abstract class FlowControlServiceHTTPBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FlowControlServiceHTTPBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckhttpProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FlowControlServiceHTTP");
    }
  }

  private static final class FlowControlServiceHTTPFileDescriptorSupplier
      extends FlowControlServiceHTTPBaseDescriptorSupplier {
    FlowControlServiceHTTPFileDescriptorSupplier() {}
  }

  private static final class FlowControlServiceHTTPMethodDescriptorSupplier
      extends FlowControlServiceHTTPBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    FlowControlServiceHTTPMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (FlowControlServiceHTTPGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FlowControlServiceHTTPFileDescriptorSupplier())
              .addMethod(getCheckHTTPMethod())
              .build();
        }
      }
    }
    return result;
  }
}
