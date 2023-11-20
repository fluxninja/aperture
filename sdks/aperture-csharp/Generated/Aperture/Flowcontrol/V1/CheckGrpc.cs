// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: aperture/flowcontrol/check/v1/check.proto
// </auto-generated>
#pragma warning disable 0414, 1591, 8981, 0612
#region Designer generated code

using grpc = global::Grpc.Core;

namespace Aperture.Flowcontrol.Check.V1 {
  /// <summary>
  /// FlowControlService is used to perform Flow Control operations.
  /// </summary>
  public static partial class FlowControlService
  {
    static readonly string __ServiceName = "aperture.flowcontrol.check.v1.FlowControlService";

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static void __Helper_SerializeMessage(global::Google.Protobuf.IMessage message, grpc::SerializationContext context)
    {
      #if !GRPC_DISABLE_PROTOBUF_BUFFER_SERIALIZATION
      if (message is global::Google.Protobuf.IBufferMessage)
      {
        context.SetPayloadLength(message.CalculateSize());
        global::Google.Protobuf.MessageExtensions.WriteTo(message, context.GetBufferWriter());
        context.Complete();
        return;
      }
      #endif
      context.Complete(global::Google.Protobuf.MessageExtensions.ToByteArray(message));
    }

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static class __Helper_MessageCache<T>
    {
      public static readonly bool IsBufferMessage = global::System.Reflection.IntrospectionExtensions.GetTypeInfo(typeof(global::Google.Protobuf.IBufferMessage)).IsAssignableFrom(typeof(T));
    }

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static T __Helper_DeserializeMessage<T>(grpc::DeserializationContext context, global::Google.Protobuf.MessageParser<T> parser) where T : global::Google.Protobuf.IMessage<T>
    {
      #if !GRPC_DISABLE_PROTOBUF_BUFFER_SERIALIZATION
      if (__Helper_MessageCache<T>.IsBufferMessage)
      {
        return parser.ParseFrom(context.PayloadAsReadOnlySequence());
      }
      #endif
      return parser.ParseFrom(context.PayloadAsNewBuffer());
    }

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::Aperture.Flowcontrol.Check.V1.CheckRequest> __Marshaller_aperture_flowcontrol_check_v1_CheckRequest = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::Aperture.Flowcontrol.Check.V1.CheckRequest.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::Aperture.Flowcontrol.Check.V1.CheckResponse> __Marshaller_aperture_flowcontrol_check_v1_CheckResponse = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::Aperture.Flowcontrol.Check.V1.CheckResponse.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest> __Marshaller_aperture_flowcontrol_check_v1_CacheUpsertRequest = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse> __Marshaller_aperture_flowcontrol_check_v1_CacheUpsertResponse = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest> __Marshaller_aperture_flowcontrol_check_v1_CacheDeleteRequest = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse> __Marshaller_aperture_flowcontrol_check_v1_CacheDeleteResponse = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse.Parser));

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Method<global::Aperture.Flowcontrol.Check.V1.CheckRequest, global::Aperture.Flowcontrol.Check.V1.CheckResponse> __Method_Check = new grpc::Method<global::Aperture.Flowcontrol.Check.V1.CheckRequest, global::Aperture.Flowcontrol.Check.V1.CheckResponse>(
        grpc::MethodType.Unary,
        __ServiceName,
        "Check",
        __Marshaller_aperture_flowcontrol_check_v1_CheckRequest,
        __Marshaller_aperture_flowcontrol_check_v1_CheckResponse);

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Method<global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest, global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse> __Method_CacheUpsert = new grpc::Method<global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest, global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse>(
        grpc::MethodType.Unary,
        __ServiceName,
        "CacheUpsert",
        __Marshaller_aperture_flowcontrol_check_v1_CacheUpsertRequest,
        __Marshaller_aperture_flowcontrol_check_v1_CacheUpsertResponse);

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Method<global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest, global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse> __Method_CacheDelete = new grpc::Method<global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest, global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse>(
        grpc::MethodType.Unary,
        __ServiceName,
        "CacheDelete",
        __Marshaller_aperture_flowcontrol_check_v1_CacheDeleteRequest,
        __Marshaller_aperture_flowcontrol_check_v1_CacheDeleteResponse);

    /// <summary>Service descriptor</summary>
    public static global::Google.Protobuf.Reflection.ServiceDescriptor Descriptor
    {
      get { return global::Aperture.Flowcontrol.Check.V1.CheckReflection.Descriptor.Services[0]; }
    }

    /// <summary>Base class for server-side implementations of FlowControlService</summary>
    [grpc::BindServiceMethod(typeof(FlowControlService), "BindService")]
    public abstract partial class FlowControlServiceBase
    {
      /// <summary>
      /// Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
      /// </summary>
      /// <param name="request">The request received from the client.</param>
      /// <param name="context">The context of the server-side call handler being invoked.</param>
      /// <returns>The response to send back to the client (wrapped by a task).</returns>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::System.Threading.Tasks.Task<global::Aperture.Flowcontrol.Check.V1.CheckResponse> Check(global::Aperture.Flowcontrol.Check.V1.CheckRequest request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::System.Threading.Tasks.Task<global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse> CacheUpsert(global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::System.Threading.Tasks.Task<global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse> CacheDelete(global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

    }

    /// <summary>Client for FlowControlService</summary>
    public partial class FlowControlServiceClient : grpc::ClientBase<FlowControlServiceClient>
    {
      /// <summary>Creates a new client for FlowControlService</summary>
      /// <param name="channel">The channel to use to make remote calls.</param>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public FlowControlServiceClient(grpc::ChannelBase channel) : base(channel)
      {
      }
      /// <summary>Creates a new client for FlowControlService that uses a custom <c>CallInvoker</c>.</summary>
      /// <param name="callInvoker">The callInvoker to use to make remote calls.</param>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public FlowControlServiceClient(grpc::CallInvoker callInvoker) : base(callInvoker)
      {
      }
      /// <summary>Protected parameterless constructor to allow creation of test doubles.</summary>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      protected FlowControlServiceClient() : base()
      {
      }
      /// <summary>Protected constructor to allow creation of configured clients.</summary>
      /// <param name="configuration">The client configuration.</param>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      protected FlowControlServiceClient(ClientBaseConfiguration configuration) : base(configuration)
      {
      }

      /// <summary>
      /// Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
      /// </summary>
      /// <param name="request">The request to send to the server.</param>
      /// <param name="headers">The initial metadata to send with the call. This parameter is optional.</param>
      /// <param name="deadline">An optional deadline for the call. The call will be cancelled if deadline is hit.</param>
      /// <param name="cancellationToken">An optional token for canceling the call.</param>
      /// <returns>The response received from the server.</returns>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::Aperture.Flowcontrol.Check.V1.CheckResponse Check(global::Aperture.Flowcontrol.Check.V1.CheckRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return Check(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      /// <summary>
      /// Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
      /// </summary>
      /// <param name="request">The request to send to the server.</param>
      /// <param name="options">The options for the call.</param>
      /// <returns>The response received from the server.</returns>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::Aperture.Flowcontrol.Check.V1.CheckResponse Check(global::Aperture.Flowcontrol.Check.V1.CheckRequest request, grpc::CallOptions options)
      {
        return CallInvoker.BlockingUnaryCall(__Method_Check, null, options, request);
      }
      /// <summary>
      /// Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
      /// </summary>
      /// <param name="request">The request to send to the server.</param>
      /// <param name="headers">The initial metadata to send with the call. This parameter is optional.</param>
      /// <param name="deadline">An optional deadline for the call. The call will be cancelled if deadline is hit.</param>
      /// <param name="cancellationToken">An optional token for canceling the call.</param>
      /// <returns>The call object.</returns>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual grpc::AsyncUnaryCall<global::Aperture.Flowcontrol.Check.V1.CheckResponse> CheckAsync(global::Aperture.Flowcontrol.Check.V1.CheckRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return CheckAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      /// <summary>
      /// Check wraps the given arbitrary resource and matches the given labels against Flow Control Limiters to makes a decision whether to allow/deny.
      /// </summary>
      /// <param name="request">The request to send to the server.</param>
      /// <param name="options">The options for the call.</param>
      /// <returns>The call object.</returns>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual grpc::AsyncUnaryCall<global::Aperture.Flowcontrol.Check.V1.CheckResponse> CheckAsync(global::Aperture.Flowcontrol.Check.V1.CheckRequest request, grpc::CallOptions options)
      {
        return CallInvoker.AsyncUnaryCall(__Method_Check, null, options, request);
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse CacheUpsert(global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return CacheUpsert(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse CacheUpsert(global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest request, grpc::CallOptions options)
      {
        return CallInvoker.BlockingUnaryCall(__Method_CacheUpsert, null, options, request);
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual grpc::AsyncUnaryCall<global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse> CacheUpsertAsync(global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return CacheUpsertAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual grpc::AsyncUnaryCall<global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse> CacheUpsertAsync(global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest request, grpc::CallOptions options)
      {
        return CallInvoker.AsyncUnaryCall(__Method_CacheUpsert, null, options, request);
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse CacheDelete(global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return CacheDelete(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse CacheDelete(global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest request, grpc::CallOptions options)
      {
        return CallInvoker.BlockingUnaryCall(__Method_CacheDelete, null, options, request);
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual grpc::AsyncUnaryCall<global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse> CacheDeleteAsync(global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return CacheDeleteAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual grpc::AsyncUnaryCall<global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse> CacheDeleteAsync(global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest request, grpc::CallOptions options)
      {
        return CallInvoker.AsyncUnaryCall(__Method_CacheDelete, null, options, request);
      }
      /// <summary>Creates a new instance of client from given <c>ClientBaseConfiguration</c>.</summary>
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      protected override FlowControlServiceClient NewInstance(ClientBaseConfiguration configuration)
      {
        return new FlowControlServiceClient(configuration);
      }
    }

    /// <summary>Creates service definition that can be registered with a server</summary>
    /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public static grpc::ServerServiceDefinition BindService(FlowControlServiceBase serviceImpl)
    {
      return grpc::ServerServiceDefinition.CreateBuilder()
          .AddMethod(__Method_Check, serviceImpl.Check)
          .AddMethod(__Method_CacheUpsert, serviceImpl.CacheUpsert)
          .AddMethod(__Method_CacheDelete, serviceImpl.CacheDelete).Build();
    }

    /// <summary>Register service method with a service binder with or without implementation. Useful when customizing the service binding logic.
    /// Note: this method is part of an experimental API that can change or be removed without any prior notice.</summary>
    /// <param name="serviceBinder">Service methods will be bound by calling <c>AddMethod</c> on this object.</param>
    /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public static void BindService(grpc::ServiceBinderBase serviceBinder, FlowControlServiceBase serviceImpl)
    {
      serviceBinder.AddMethod(__Method_Check, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::Aperture.Flowcontrol.Check.V1.CheckRequest, global::Aperture.Flowcontrol.Check.V1.CheckResponse>(serviceImpl.Check));
      serviceBinder.AddMethod(__Method_CacheUpsert, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::Aperture.Flowcontrol.Check.V1.CacheUpsertRequest, global::Aperture.Flowcontrol.Check.V1.CacheUpsertResponse>(serviceImpl.CacheUpsert));
      serviceBinder.AddMethod(__Method_CacheDelete, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::Aperture.Flowcontrol.Check.V1.CacheDeleteRequest, global::Aperture.Flowcontrol.Check.V1.CacheDeleteResponse>(serviceImpl.CacheDelete));
    }

  }
}
#endregion
