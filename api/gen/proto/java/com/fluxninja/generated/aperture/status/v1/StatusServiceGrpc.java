package com.fluxninja.generated.aperture.status.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * StatusService is used to query Jobs.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/status/v1/status.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class StatusServiceGrpc {

  private StatusServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.status.v1.StatusService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.status.v1.GroupStatusRequest,
      com.fluxninja.generated.aperture.status.v1.GroupStatus> getGetGroupStatusMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetGroupStatus",
      requestType = com.fluxninja.generated.aperture.status.v1.GroupStatusRequest.class,
      responseType = com.fluxninja.generated.aperture.status.v1.GroupStatus.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.status.v1.GroupStatusRequest,
      com.fluxninja.generated.aperture.status.v1.GroupStatus> getGetGroupStatusMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.status.v1.GroupStatusRequest, com.fluxninja.generated.aperture.status.v1.GroupStatus> getGetGroupStatusMethod;
    if ((getGetGroupStatusMethod = StatusServiceGrpc.getGetGroupStatusMethod) == null) {
      synchronized (StatusServiceGrpc.class) {
        if ((getGetGroupStatusMethod = StatusServiceGrpc.getGetGroupStatusMethod) == null) {
          StatusServiceGrpc.getGetGroupStatusMethod = getGetGroupStatusMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.status.v1.GroupStatusRequest, com.fluxninja.generated.aperture.status.v1.GroupStatus>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetGroupStatus"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.status.v1.GroupStatusRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.status.v1.GroupStatus.getDefaultInstance()))
              .setSchemaDescriptor(new StatusServiceMethodDescriptorSupplier("GetGroupStatus"))
              .build();
        }
      }
    }
    return getGetGroupStatusMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static StatusServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StatusServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StatusServiceStub>() {
        @java.lang.Override
        public StatusServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StatusServiceStub(channel, callOptions);
        }
      };
    return StatusServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static StatusServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StatusServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StatusServiceBlockingStub>() {
        @java.lang.Override
        public StatusServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StatusServiceBlockingStub(channel, callOptions);
        }
      };
    return StatusServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static StatusServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StatusServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StatusServiceFutureStub>() {
        @java.lang.Override
        public StatusServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StatusServiceFutureStub(channel, callOptions);
        }
      };
    return StatusServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * StatusService is used to query Jobs.
   * </pre>
   */
  public static abstract class StatusServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getGroupStatus(com.fluxninja.generated.aperture.status.v1.GroupStatusRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.status.v1.GroupStatus> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetGroupStatusMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetGroupStatusMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.status.v1.GroupStatusRequest,
                com.fluxninja.generated.aperture.status.v1.GroupStatus>(
                  this, METHODID_GET_GROUP_STATUS)))
          .build();
    }
  }

  /**
   * <pre>
   * StatusService is used to query Jobs.
   * </pre>
   */
  public static final class StatusServiceStub extends io.grpc.stub.AbstractAsyncStub<StatusServiceStub> {
    private StatusServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StatusServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StatusServiceStub(channel, callOptions);
    }

    /**
     */
    public void getGroupStatus(com.fluxninja.generated.aperture.status.v1.GroupStatusRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.status.v1.GroupStatus> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetGroupStatusMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * StatusService is used to query Jobs.
   * </pre>
   */
  public static final class StatusServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<StatusServiceBlockingStub> {
    private StatusServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StatusServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StatusServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.status.v1.GroupStatus getGroupStatus(com.fluxninja.generated.aperture.status.v1.GroupStatusRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetGroupStatusMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * StatusService is used to query Jobs.
   * </pre>
   */
  public static final class StatusServiceFutureStub extends io.grpc.stub.AbstractFutureStub<StatusServiceFutureStub> {
    private StatusServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StatusServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StatusServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.status.v1.GroupStatus> getGroupStatus(
        com.fluxninja.generated.aperture.status.v1.GroupStatusRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetGroupStatusMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_GROUP_STATUS = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final StatusServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(StatusServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_GROUP_STATUS:
          serviceImpl.getGroupStatus((com.fluxninja.generated.aperture.status.v1.GroupStatusRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.status.v1.GroupStatus>) responseObserver);
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

  private static abstract class StatusServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    StatusServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.status.v1.StatusProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("StatusService");
    }
  }

  private static final class StatusServiceFileDescriptorSupplier
      extends StatusServiceBaseDescriptorSupplier {
    StatusServiceFileDescriptorSupplier() {}
  }

  private static final class StatusServiceMethodDescriptorSupplier
      extends StatusServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    StatusServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (StatusServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new StatusServiceFileDescriptorSupplier())
              .addMethod(getGetGroupStatusMethod())
              .build();
        }
      }
    }
    return result;
  }
}
