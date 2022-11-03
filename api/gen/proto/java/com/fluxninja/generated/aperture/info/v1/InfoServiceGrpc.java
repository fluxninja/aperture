package com.fluxninja.generated.aperture.info.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * InfoService is used to provide information about the aperture system.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/info/v1/info.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class InfoServiceGrpc {

  private InfoServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.info.v1.InfoService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.info.v1.VersionInfo> getVersionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Version",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.info.v1.VersionInfo.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.info.v1.VersionInfo> getVersionMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.info.v1.VersionInfo> getVersionMethod;
    if ((getVersionMethod = InfoServiceGrpc.getVersionMethod) == null) {
      synchronized (InfoServiceGrpc.class) {
        if ((getVersionMethod = InfoServiceGrpc.getVersionMethod) == null) {
          InfoServiceGrpc.getVersionMethod = getVersionMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.info.v1.VersionInfo>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Version"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.info.v1.VersionInfo.getDefaultInstance()))
              .setSchemaDescriptor(new InfoServiceMethodDescriptorSupplier("Version"))
              .build();
        }
      }
    }
    return getVersionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.info.v1.ProcessInfo> getProcessMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Process",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.info.v1.ProcessInfo.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.info.v1.ProcessInfo> getProcessMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.info.v1.ProcessInfo> getProcessMethod;
    if ((getProcessMethod = InfoServiceGrpc.getProcessMethod) == null) {
      synchronized (InfoServiceGrpc.class) {
        if ((getProcessMethod = InfoServiceGrpc.getProcessMethod) == null) {
          InfoServiceGrpc.getProcessMethod = getProcessMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.info.v1.ProcessInfo>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Process"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.info.v1.ProcessInfo.getDefaultInstance()))
              .setSchemaDescriptor(new InfoServiceMethodDescriptorSupplier("Process"))
              .build();
        }
      }
    }
    return getProcessMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.info.v1.HostInfo> getHostMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Host",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.info.v1.HostInfo.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.info.v1.HostInfo> getHostMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.info.v1.HostInfo> getHostMethod;
    if ((getHostMethod = InfoServiceGrpc.getHostMethod) == null) {
      synchronized (InfoServiceGrpc.class) {
        if ((getHostMethod = InfoServiceGrpc.getHostMethod) == null) {
          InfoServiceGrpc.getHostMethod = getHostMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.info.v1.HostInfo>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Host"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.info.v1.HostInfo.getDefaultInstance()))
              .setSchemaDescriptor(new InfoServiceMethodDescriptorSupplier("Host"))
              .build();
        }
      }
    }
    return getHostMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static InfoServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<InfoServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<InfoServiceStub>() {
        @java.lang.Override
        public InfoServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new InfoServiceStub(channel, callOptions);
        }
      };
    return InfoServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static InfoServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<InfoServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<InfoServiceBlockingStub>() {
        @java.lang.Override
        public InfoServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new InfoServiceBlockingStub(channel, callOptions);
        }
      };
    return InfoServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static InfoServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<InfoServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<InfoServiceFutureStub>() {
        @java.lang.Override
        public InfoServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new InfoServiceFutureStub(channel, callOptions);
        }
      };
    return InfoServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * InfoService is used to provide information about the aperture system.
   * </pre>
   */
  public static abstract class InfoServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void version(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.VersionInfo> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getVersionMethod(), responseObserver);
    }

    /**
     */
    public void process(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.ProcessInfo> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getProcessMethod(), responseObserver);
    }

    /**
     */
    public void host(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.HostInfo> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getHostMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getVersionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.info.v1.VersionInfo>(
                  this, METHODID_VERSION)))
          .addMethod(
            getProcessMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.info.v1.ProcessInfo>(
                  this, METHODID_PROCESS)))
          .addMethod(
            getHostMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.info.v1.HostInfo>(
                  this, METHODID_HOST)))
          .build();
    }
  }

  /**
   * <pre>
   * InfoService is used to provide information about the aperture system.
   * </pre>
   */
  public static final class InfoServiceStub extends io.grpc.stub.AbstractAsyncStub<InfoServiceStub> {
    private InfoServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected InfoServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new InfoServiceStub(channel, callOptions);
    }

    /**
     */
    public void version(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.VersionInfo> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getVersionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void process(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.ProcessInfo> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getProcessMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void host(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.HostInfo> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getHostMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * InfoService is used to provide information about the aperture system.
   * </pre>
   */
  public static final class InfoServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<InfoServiceBlockingStub> {
    private InfoServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected InfoServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new InfoServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.info.v1.VersionInfo version(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getVersionMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.info.v1.ProcessInfo process(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getProcessMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.info.v1.HostInfo host(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getHostMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * InfoService is used to provide information about the aperture system.
   * </pre>
   */
  public static final class InfoServiceFutureStub extends io.grpc.stub.AbstractFutureStub<InfoServiceFutureStub> {
    private InfoServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected InfoServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new InfoServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.info.v1.VersionInfo> version(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getVersionMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.info.v1.ProcessInfo> process(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getProcessMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.info.v1.HostInfo> host(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getHostMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_VERSION = 0;
  private static final int METHODID_PROCESS = 1;
  private static final int METHODID_HOST = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final InfoServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(InfoServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_VERSION:
          serviceImpl.version((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.VersionInfo>) responseObserver);
          break;
        case METHODID_PROCESS:
          serviceImpl.process((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.ProcessInfo>) responseObserver);
          break;
        case METHODID_HOST:
          serviceImpl.host((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.info.v1.HostInfo>) responseObserver);
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

  private static abstract class InfoServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    InfoServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.info.v1.InfoProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("InfoService");
    }
  }

  private static final class InfoServiceFileDescriptorSupplier
      extends InfoServiceBaseDescriptorSupplier {
    InfoServiceFileDescriptorSupplier() {}
  }

  private static final class InfoServiceMethodDescriptorSupplier
      extends InfoServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    InfoServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (InfoServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new InfoServiceFileDescriptorSupplier())
              .addMethod(getVersionMethod())
              .addMethod(getProcessMethod())
              .addMethod(getHostMethod())
              .build();
        }
      }
    }
    return result;
  }
}
