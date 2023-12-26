import type * as grpc from '@grpc/grpc-js';
import type { EnumTypeDefinition, MessageTypeDefinition } from '@grpc/proto-loader';

import type { FlowControlServiceClient as _aperture_flowcontrol_check_v1_FlowControlServiceClient, FlowControlServiceDefinition as _aperture_flowcontrol_check_v1_FlowControlServiceDefinition } from './aperture/flowcontrol/check/v1/FlowControlService';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  aperture: {
    flowcontrol: {
      check: {
        v1: {
          CacheDeleteRequest: MessageTypeDefinition
          CacheDeleteResponse: MessageTypeDefinition
          CacheEntry: MessageTypeDefinition
          CacheLookupRequest: MessageTypeDefinition
          CacheLookupResponse: MessageTypeDefinition
          CacheLookupStatus: EnumTypeDefinition
          CacheOperationStatus: EnumTypeDefinition
          CacheUpsertRequest: MessageTypeDefinition
          CacheUpsertResponse: MessageTypeDefinition
          CheckRequest: MessageTypeDefinition
          CheckResponse: MessageTypeDefinition
          ClassifierInfo: MessageTypeDefinition
          FlowControlService: SubtypeConstructor<typeof grpc.Client, _aperture_flowcontrol_check_v1_FlowControlServiceClient> & { service: _aperture_flowcontrol_check_v1_FlowControlServiceDefinition }
          FlowEndRequest: MessageTypeDefinition
          FlowEndResponse: MessageTypeDefinition
          FluxMeterInfo: MessageTypeDefinition
          InflightRequestRef: MessageTypeDefinition
          KeyDeleteResponse: MessageTypeDefinition
          KeyLookupResponse: MessageTypeDefinition
          KeyUpsertResponse: MessageTypeDefinition
          LimiterDecision: MessageTypeDefinition
          StatusCode: EnumTypeDefinition
          TokenReturnStatus: MessageTypeDefinition
        }
      }
    }
  }
  google: {
    protobuf: {
      Duration: MessageTypeDefinition
      Timestamp: MessageTypeDefinition
    }
  }
}

