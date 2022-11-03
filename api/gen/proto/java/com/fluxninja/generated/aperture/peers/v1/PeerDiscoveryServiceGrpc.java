package com.fluxninja.generated.aperture.peers.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * PeerDiscoveryService is used to query Peers.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.49.1)",
    comments = "Source: aperture/peers/v1/peers.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class PeerDiscoveryServiceGrpc {

  private PeerDiscoveryServiceGrpc() {}

  public static final String SERVICE_NAME = "aperture.peers.v1.PeerDiscoveryService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.peers.v1.Peers> getGetPeersMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetPeers",
      requestType = com.google.protobuf.Empty.class,
      responseType = com.fluxninja.generated.aperture.peers.v1.Peers.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      com.fluxninja.generated.aperture.peers.v1.Peers> getGetPeersMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, com.fluxninja.generated.aperture.peers.v1.Peers> getGetPeersMethod;
    if ((getGetPeersMethod = PeerDiscoveryServiceGrpc.getGetPeersMethod) == null) {
      synchronized (PeerDiscoveryServiceGrpc.class) {
        if ((getGetPeersMethod = PeerDiscoveryServiceGrpc.getGetPeersMethod) == null) {
          PeerDiscoveryServiceGrpc.getGetPeersMethod = getGetPeersMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, com.fluxninja.generated.aperture.peers.v1.Peers>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetPeers"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.peers.v1.Peers.getDefaultInstance()))
              .setSchemaDescriptor(new PeerDiscoveryServiceMethodDescriptorSupplier("GetPeers"))
              .build();
        }
      }
    }
    return getGetPeersMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.peers.v1.PeerRequest,
      com.fluxninja.generated.aperture.peers.v1.Peer> getGetPeerMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetPeer",
      requestType = com.fluxninja.generated.aperture.peers.v1.PeerRequest.class,
      responseType = com.fluxninja.generated.aperture.peers.v1.Peer.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.peers.v1.PeerRequest,
      com.fluxninja.generated.aperture.peers.v1.Peer> getGetPeerMethod() {
    io.grpc.MethodDescriptor<com.fluxninja.generated.aperture.peers.v1.PeerRequest, com.fluxninja.generated.aperture.peers.v1.Peer> getGetPeerMethod;
    if ((getGetPeerMethod = PeerDiscoveryServiceGrpc.getGetPeerMethod) == null) {
      synchronized (PeerDiscoveryServiceGrpc.class) {
        if ((getGetPeerMethod = PeerDiscoveryServiceGrpc.getGetPeerMethod) == null) {
          PeerDiscoveryServiceGrpc.getGetPeerMethod = getGetPeerMethod =
              io.grpc.MethodDescriptor.<com.fluxninja.generated.aperture.peers.v1.PeerRequest, com.fluxninja.generated.aperture.peers.v1.Peer>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetPeer"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.peers.v1.PeerRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.fluxninja.generated.aperture.peers.v1.Peer.getDefaultInstance()))
              .setSchemaDescriptor(new PeerDiscoveryServiceMethodDescriptorSupplier("GetPeer"))
              .build();
        }
      }
    }
    return getGetPeerMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static PeerDiscoveryServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PeerDiscoveryServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PeerDiscoveryServiceStub>() {
        @java.lang.Override
        public PeerDiscoveryServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PeerDiscoveryServiceStub(channel, callOptions);
        }
      };
    return PeerDiscoveryServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PeerDiscoveryServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PeerDiscoveryServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PeerDiscoveryServiceBlockingStub>() {
        @java.lang.Override
        public PeerDiscoveryServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PeerDiscoveryServiceBlockingStub(channel, callOptions);
        }
      };
    return PeerDiscoveryServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static PeerDiscoveryServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PeerDiscoveryServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PeerDiscoveryServiceFutureStub>() {
        @java.lang.Override
        public PeerDiscoveryServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PeerDiscoveryServiceFutureStub(channel, callOptions);
        }
      };
    return PeerDiscoveryServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * PeerDiscoveryService is used to query Peers.
   * </pre>
   */
  public static abstract class PeerDiscoveryServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getPeers(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.peers.v1.Peers> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetPeersMethod(), responseObserver);
    }

    /**
     */
    public void getPeer(com.fluxninja.generated.aperture.peers.v1.PeerRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.peers.v1.Peer> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetPeerMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetPeersMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                com.fluxninja.generated.aperture.peers.v1.Peers>(
                  this, METHODID_GET_PEERS)))
          .addMethod(
            getGetPeerMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.fluxninja.generated.aperture.peers.v1.PeerRequest,
                com.fluxninja.generated.aperture.peers.v1.Peer>(
                  this, METHODID_GET_PEER)))
          .build();
    }
  }

  /**
   * <pre>
   * PeerDiscoveryService is used to query Peers.
   * </pre>
   */
  public static final class PeerDiscoveryServiceStub extends io.grpc.stub.AbstractAsyncStub<PeerDiscoveryServiceStub> {
    private PeerDiscoveryServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PeerDiscoveryServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PeerDiscoveryServiceStub(channel, callOptions);
    }

    /**
     */
    public void getPeers(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.peers.v1.Peers> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetPeersMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getPeer(com.fluxninja.generated.aperture.peers.v1.PeerRequest request,
        io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.peers.v1.Peer> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetPeerMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * PeerDiscoveryService is used to query Peers.
   * </pre>
   */
  public static final class PeerDiscoveryServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<PeerDiscoveryServiceBlockingStub> {
    private PeerDiscoveryServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PeerDiscoveryServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PeerDiscoveryServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.fluxninja.generated.aperture.peers.v1.Peers getPeers(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetPeersMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.fluxninja.generated.aperture.peers.v1.Peer getPeer(com.fluxninja.generated.aperture.peers.v1.PeerRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetPeerMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * PeerDiscoveryService is used to query Peers.
   * </pre>
   */
  public static final class PeerDiscoveryServiceFutureStub extends io.grpc.stub.AbstractFutureStub<PeerDiscoveryServiceFutureStub> {
    private PeerDiscoveryServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PeerDiscoveryServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PeerDiscoveryServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.peers.v1.Peers> getPeers(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetPeersMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.fluxninja.generated.aperture.peers.v1.Peer> getPeer(
        com.fluxninja.generated.aperture.peers.v1.PeerRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetPeerMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_PEERS = 0;
  private static final int METHODID_GET_PEER = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final PeerDiscoveryServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(PeerDiscoveryServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_PEERS:
          serviceImpl.getPeers((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.peers.v1.Peers>) responseObserver);
          break;
        case METHODID_GET_PEER:
          serviceImpl.getPeer((com.fluxninja.generated.aperture.peers.v1.PeerRequest) request,
              (io.grpc.stub.StreamObserver<com.fluxninja.generated.aperture.peers.v1.Peer>) responseObserver);
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

  private static abstract class PeerDiscoveryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PeerDiscoveryServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.fluxninja.generated.aperture.peers.v1.PeersProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("PeerDiscoveryService");
    }
  }

  private static final class PeerDiscoveryServiceFileDescriptorSupplier
      extends PeerDiscoveryServiceBaseDescriptorSupplier {
    PeerDiscoveryServiceFileDescriptorSupplier() {}
  }

  private static final class PeerDiscoveryServiceMethodDescriptorSupplier
      extends PeerDiscoveryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    PeerDiscoveryServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (PeerDiscoveryServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new PeerDiscoveryServiceFileDescriptorSupplier())
              .addMethod(getGetPeersMethod())
              .addMethod(getGetPeerMethod())
              .build();
        }
      }
    }
    return result;
  }
}
