// Original file: proto/flowcontrol/check/v1/check.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CacheDeleteRequest as _aperture_flowcontrol_check_v1_CacheDeleteRequest, CacheDeleteRequest__Output as _aperture_flowcontrol_check_v1_CacheDeleteRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CacheDeleteRequest';
import type { CacheUpsertRequest as _aperture_flowcontrol_check_v1_CacheUpsertRequest, CacheUpsertRequest__Output as _aperture_flowcontrol_check_v1_CacheUpsertRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CacheUpsertRequest';
import type { CheckRequest as _aperture_flowcontrol_check_v1_CheckRequest, CheckRequest__Output as _aperture_flowcontrol_check_v1_CheckRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CheckRequest';
import type { CheckResponse as _aperture_flowcontrol_check_v1_CheckResponse, CheckResponse__Output as _aperture_flowcontrol_check_v1_CheckResponse__Output } from '../../../../aperture/flowcontrol/check/v1/CheckResponse';
import type { Empty as _google_protobuf_Empty, Empty__Output as _google_protobuf_Empty__Output } from '../../../../google/protobuf/Empty';

export interface FlowControlServiceClient extends grpc.Client {
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  
  Check(argument: _aperture_flowcontrol_check_v1_CheckRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  Check(argument: _aperture_flowcontrol_check_v1_CheckRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  Check(argument: _aperture_flowcontrol_check_v1_CheckRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  Check(argument: _aperture_flowcontrol_check_v1_CheckRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  check(argument: _aperture_flowcontrol_check_v1_CheckRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  check(argument: _aperture_flowcontrol_check_v1_CheckRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  check(argument: _aperture_flowcontrol_check_v1_CheckRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  check(argument: _aperture_flowcontrol_check_v1_CheckRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CheckResponse__Output>): grpc.ClientUnaryCall;
  
}

export interface FlowControlServiceHandlers extends grpc.UntypedServiceImplementation {
  CacheDelete: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CacheDeleteRequest__Output, _google_protobuf_Empty>;
  
  CacheUpsert: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CacheUpsertRequest__Output, _google_protobuf_Empty>;
  
  Check: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse>;
  
}

export interface FlowControlServiceDefinition extends grpc.ServiceDefinition {
  CacheDelete: MethodDefinition<_aperture_flowcontrol_check_v1_CacheDeleteRequest, _google_protobuf_Empty, _aperture_flowcontrol_check_v1_CacheDeleteRequest__Output, _google_protobuf_Empty__Output>
  CacheUpsert: MethodDefinition<_aperture_flowcontrol_check_v1_CacheUpsertRequest, _google_protobuf_Empty, _aperture_flowcontrol_check_v1_CacheUpsertRequest__Output, _google_protobuf_Empty__Output>
  Check: MethodDefinition<_aperture_flowcontrol_check_v1_CheckRequest, _aperture_flowcontrol_check_v1_CheckResponse, _aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse__Output>
}
