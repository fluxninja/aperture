// Original file: proto/flowcontrol/check/v1/check.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CacheRequest as _aperture_flowcontrol_check_v1_CacheRequest, CacheRequest__Output as _aperture_flowcontrol_check_v1_CacheRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CacheRequest';
import type { CheckRequest as _aperture_flowcontrol_check_v1_CheckRequest, CheckRequest__Output as _aperture_flowcontrol_check_v1_CheckRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CheckRequest';
import type { CheckResponse as _aperture_flowcontrol_check_v1_CheckResponse, CheckResponse__Output as _aperture_flowcontrol_check_v1_CheckResponse__Output } from '../../../../aperture/flowcontrol/check/v1/CheckResponse';
import type { Empty as _google_protobuf_Empty, Empty__Output as _google_protobuf_Empty__Output } from '../../../../google/protobuf/Empty';

export interface FlowControlServiceClient extends grpc.Client {
  CacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  CacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  cacheUpdate(argument: _aperture_flowcontrol_check_v1_CacheRequest, callback: grpc.requestCallback<_google_protobuf_Empty__Output>): grpc.ClientUnaryCall;
  
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
  CacheUpdate: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CacheRequest__Output, _google_protobuf_Empty>;
  
  Check: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse>;
  
}

export interface FlowControlServiceDefinition extends grpc.ServiceDefinition {
  CacheUpdate: MethodDefinition<_aperture_flowcontrol_check_v1_CacheRequest, _google_protobuf_Empty, _aperture_flowcontrol_check_v1_CacheRequest__Output, _google_protobuf_Empty__Output>
  Check: MethodDefinition<_aperture_flowcontrol_check_v1_CheckRequest, _aperture_flowcontrol_check_v1_CheckResponse, _aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse__Output>
}
