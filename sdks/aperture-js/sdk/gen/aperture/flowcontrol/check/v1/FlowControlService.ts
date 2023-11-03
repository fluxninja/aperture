// Original file: proto/flowcontrol/check/v1/check.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CheckRequest as _aperture_flowcontrol_check_v1_CheckRequest, CheckRequest__Output as _aperture_flowcontrol_check_v1_CheckRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CheckRequest';
import type { CheckResponse as _aperture_flowcontrol_check_v1_CheckResponse, CheckResponse__Output as _aperture_flowcontrol_check_v1_CheckResponse__Output } from '../../../../aperture/flowcontrol/check/v1/CheckResponse';

export interface FlowControlServiceClient extends grpc.Client {
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
  Check: grpc.handleUnaryCall<_aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse>;
  
}

export interface FlowControlServiceDefinition extends grpc.ServiceDefinition {
  Check: MethodDefinition<_aperture_flowcontrol_check_v1_CheckRequest, _aperture_flowcontrol_check_v1_CheckResponse, _aperture_flowcontrol_check_v1_CheckRequest__Output, _aperture_flowcontrol_check_v1_CheckResponse__Output>
}
