// package: aperture.flowcontrol.check.v1
// file: api/aperture/flowcontrol/check/v1/check.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as api_aperture_flowcontrol_check_v1_check_pb from "../../../../../api/aperture/flowcontrol/check/v1/check_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

interface IFlowControlServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    check: IFlowControlServiceService_ICheck;
}

interface IFlowControlServiceService_ICheck extends grpc.MethodDefinition<api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, api_aperture_flowcontrol_check_v1_check_pb.CheckResponse> {
    path: "/aperture.flowcontrol.check.v1.FlowControlService/Check";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_aperture_flowcontrol_check_v1_check_pb.CheckRequest>;
    requestDeserialize: grpc.deserialize<api_aperture_flowcontrol_check_v1_check_pb.CheckRequest>;
    responseSerialize: grpc.serialize<api_aperture_flowcontrol_check_v1_check_pb.CheckResponse>;
    responseDeserialize: grpc.deserialize<api_aperture_flowcontrol_check_v1_check_pb.CheckResponse>;
}

export const FlowControlServiceService: IFlowControlServiceService;

export interface IFlowControlServiceServer extends grpc.UntypedServiceImplementation {
    check: grpc.handleUnaryCall<api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, api_aperture_flowcontrol_check_v1_check_pb.CheckResponse>;
}

export interface IFlowControlServiceClient {
    check(request: api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, callback: (error: grpc.ServiceError | null, response: api_aperture_flowcontrol_check_v1_check_pb.CheckResponse) => void): grpc.ClientUnaryCall;
    check(request: api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_aperture_flowcontrol_check_v1_check_pb.CheckResponse) => void): grpc.ClientUnaryCall;
    check(request: api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_aperture_flowcontrol_check_v1_check_pb.CheckResponse) => void): grpc.ClientUnaryCall;
}

export class FlowControlServiceClient extends grpc.Client implements IFlowControlServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public check(request: api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, callback: (error: grpc.ServiceError | null, response: api_aperture_flowcontrol_check_v1_check_pb.CheckResponse) => void): grpc.ClientUnaryCall;
    public check(request: api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_aperture_flowcontrol_check_v1_check_pb.CheckResponse) => void): grpc.ClientUnaryCall;
    public check(request: api_aperture_flowcontrol_check_v1_check_pb.CheckRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_aperture_flowcontrol_check_v1_check_pb.CheckResponse) => void): grpc.ClientUnaryCall;
}
