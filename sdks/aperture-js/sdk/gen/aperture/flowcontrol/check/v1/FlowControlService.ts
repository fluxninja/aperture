// Original file: proto/flowcontrol/check/v1/check.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CacheDeleteRequest as _aperture_flowcontrol_check_v1_CacheDeleteRequest, CacheDeleteRequest__Output as _aperture_flowcontrol_check_v1_CacheDeleteRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CacheDeleteRequest';
import type { CacheDeleteResponse as _aperture_flowcontrol_check_v1_CacheDeleteResponse, CacheDeleteResponse__Output as _aperture_flowcontrol_check_v1_CacheDeleteResponse__Output } from '../../../../aperture/flowcontrol/check/v1/CacheDeleteResponse';
import type { CacheLookupRequest as _aperture_flowcontrol_check_v1_CacheLookupRequest, CacheLookupRequest__Output as _aperture_flowcontrol_check_v1_CacheLookupRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CacheLookupRequest';
import type { CacheLookupResponse as _aperture_flowcontrol_check_v1_CacheLookupResponse, CacheLookupResponse__Output as _aperture_flowcontrol_check_v1_CacheLookupResponse__Output } from '../../../../aperture/flowcontrol/check/v1/CacheLookupResponse';
import type { CacheUpsertRequest as _aperture_flowcontrol_check_v1_CacheUpsertRequest, CacheUpsertRequest__Output as _aperture_flowcontrol_check_v1_CacheUpsertRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CacheUpsertRequest';
import type { CacheUpsertResponse as _aperture_flowcontrol_check_v1_CacheUpsertResponse, CacheUpsertResponse__Output as _aperture_flowcontrol_check_v1_CacheUpsertResponse__Output } from '../../../../aperture/flowcontrol/check/v1/CacheUpsertResponse';
import type { CheckRequest as _aperture_flowcontrol_check_v1_CheckRequest, CheckRequest__Output as _aperture_flowcontrol_check_v1_CheckRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CheckRequest';
import type { CheckResponse as _aperture_flowcontrol_check_v1_CheckResponse, CheckResponse__Output as _aperture_flowcontrol_check_v1_CheckResponse__Output } from '../../../../aperture/flowcontrol/check/v1/CheckResponse';

export interface FlowControlServiceClient extends grpc.Client {
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  CacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  cacheDelete(argument: _aperture_flowcontrol_check_v1_CacheDeleteRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>): grpc.ClientUnaryCall;
  
  CacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  CacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  CacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  CacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  cacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  cacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  cacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  cacheLookup(argument: _aperture_flowcontrol_check_v1_CacheLookupRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheLookupResponse__Output>): grpc.ClientUnaryCall;
  
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  CacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  cacheUpsert(argument: _aperture_flowcontrol_check_v1_CacheUpsertRequest, callback: grpc.requestCallback<_aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>): grpc.ClientUnaryCall;
  
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
  CacheDelete: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CacheDeleteRequest__Output, _aperture_flowcontrol_check_v1_CacheDeleteResponse>;
  
  CacheLookup: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CacheLookupRequest__Output, _aperture_flowcontrol_check_v1_CacheLookupResponse>;
  
  CacheUpsert: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CacheUpsertRequest__Output, _aperture_flowcontrol_check_v1_CacheUpsertResponse>;
  
  Check: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse>;
  
}

export interface FlowControlServiceDefinition extends grpc.ServiceDefinition {
  CacheDelete: MethodDefinition<_aperture_flowcontrol_check_v1_CacheDeleteRequest, _aperture_flowcontrol_check_v1_CacheDeleteResponse, _aperture_flowcontrol_check_v1_CacheDeleteRequest__Output, _aperture_flowcontrol_check_v1_CacheDeleteResponse__Output>
  CacheLookup: MethodDefinition<_aperture_flowcontrol_check_v1_CacheLookupRequest, _aperture_flowcontrol_check_v1_CacheLookupResponse, _aperture_flowcontrol_check_v1_CacheLookupRequest__Output, _aperture_flowcontrol_check_v1_CacheLookupResponse__Output>
  CacheUpsert: MethodDefinition<_aperture_flowcontrol_check_v1_CacheUpsertRequest, _aperture_flowcontrol_check_v1_CacheUpsertResponse, _aperture_flowcontrol_check_v1_CacheUpsertRequest__Output, _aperture_flowcontrol_check_v1_CacheUpsertResponse__Output>
  Check: MethodDefinition<_aperture_flowcontrol_check_v1_CheckRequest, _aperture_flowcontrol_check_v1_CheckResponse, _aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse__Output>
}
