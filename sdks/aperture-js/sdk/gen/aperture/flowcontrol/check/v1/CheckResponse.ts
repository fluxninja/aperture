// Original file: proto/flowcontrol/check/v1/check.proto

import type { Timestamp as _google_protobuf_Timestamp, Timestamp__Output as _google_protobuf_Timestamp__Output } from '../../../../google/protobuf/Timestamp';
import type { ClassifierInfo as _aperture_flowcontrol_check_v1_ClassifierInfo, ClassifierInfo__Output as _aperture_flowcontrol_check_v1_ClassifierInfo__Output } from '../../../../aperture/flowcontrol/check/v1/ClassifierInfo';
import type { FluxMeterInfo as _aperture_flowcontrol_check_v1_FluxMeterInfo, FluxMeterInfo__Output as _aperture_flowcontrol_check_v1_FluxMeterInfo__Output } from '../../../../aperture/flowcontrol/check/v1/FluxMeterInfo';
import type { LimiterDecision as _aperture_flowcontrol_check_v1_LimiterDecision, LimiterDecision__Output as _aperture_flowcontrol_check_v1_LimiterDecision__Output } from '../../../../aperture/flowcontrol/check/v1/LimiterDecision';
import type { Duration as _google_protobuf_Duration, Duration__Output as _google_protobuf_Duration__Output } from '../../../../google/protobuf/Duration';
import type { StatusCode as _aperture_flowcontrol_check_v1_StatusCode, StatusCode__Output as _aperture_flowcontrol_check_v1_StatusCode__Output } from '../../../../aperture/flowcontrol/check/v1/StatusCode';
import type { CachedValue as _aperture_flowcontrol_check_v1_CachedValue, CachedValue__Output as _aperture_flowcontrol_check_v1_CachedValue__Output } from '../../../../aperture/flowcontrol/check/v1/CachedValue';

// Original file: proto/flowcontrol/check/v1/check.proto

export const _aperture_flowcontrol_check_v1_CheckResponse_DecisionType = {
  DECISION_TYPE_ACCEPTED: 0,
  DECISION_TYPE_REJECTED: 1,
} as const;

export type _aperture_flowcontrol_check_v1_CheckResponse_DecisionType =
  | 'DECISION_TYPE_ACCEPTED'
  | 0
  | 'DECISION_TYPE_REJECTED'
  | 1

export type _aperture_flowcontrol_check_v1_CheckResponse_DecisionType__Output = typeof _aperture_flowcontrol_check_v1_CheckResponse_DecisionType[keyof typeof _aperture_flowcontrol_check_v1_CheckResponse_DecisionType]

// Original file: proto/flowcontrol/check/v1/check.proto

export const _aperture_flowcontrol_check_v1_CheckResponse_RejectReason = {
  REJECT_REASON_NONE: 0,
  REJECT_REASON_RATE_LIMITED: 1,
  REJECT_REASON_NO_TOKENS: 2,
  REJECT_REASON_NOT_SAMPLED: 3,
  REJECT_REASON_NO_MATCHING_RAMP: 4,
} as const;

export type _aperture_flowcontrol_check_v1_CheckResponse_RejectReason =
  | 'REJECT_REASON_NONE'
  | 0
  | 'REJECT_REASON_RATE_LIMITED'
  | 1
  | 'REJECT_REASON_NO_TOKENS'
  | 2
  | 'REJECT_REASON_NOT_SAMPLED'
  | 3
  | 'REJECT_REASON_NO_MATCHING_RAMP'
  | 4

export type _aperture_flowcontrol_check_v1_CheckResponse_RejectReason__Output = typeof _aperture_flowcontrol_check_v1_CheckResponse_RejectReason[keyof typeof _aperture_flowcontrol_check_v1_CheckResponse_RejectReason]

export interface CheckResponse {
  'start'?: (_google_protobuf_Timestamp | null);
  'end'?: (_google_protobuf_Timestamp | null);
  'services'?: (string)[];
  'controlPoint'?: (string);
  'flowLabelKeys'?: (string)[];
  'telemetryFlowLabels'?: ({[key: string]: string});
  'decisionType'?: (_aperture_flowcontrol_check_v1_CheckResponse_DecisionType);
  'rejectReason'?: (_aperture_flowcontrol_check_v1_CheckResponse_RejectReason);
  'classifierInfos'?: (_aperture_flowcontrol_check_v1_ClassifierInfo)[];
  'fluxMeterInfos'?: (_aperture_flowcontrol_check_v1_FluxMeterInfo)[];
  'limiterDecisions'?: (_aperture_flowcontrol_check_v1_LimiterDecision)[];
  'waitTime'?: (_google_protobuf_Duration | null);
  'deniedResponseStatusCode'?: (_aperture_flowcontrol_check_v1_StatusCode);
  'cachedValue'?: (_aperture_flowcontrol_check_v1_CachedValue | null);
}

export interface CheckResponse__Output {
  'start': (_google_protobuf_Timestamp__Output | null);
  'end': (_google_protobuf_Timestamp__Output | null);
  'services': (string)[];
  'controlPoint': (string);
  'flowLabelKeys': (string)[];
  'telemetryFlowLabels': ({[key: string]: string});
  'decisionType': (_aperture_flowcontrol_check_v1_CheckResponse_DecisionType__Output);
  'rejectReason': (_aperture_flowcontrol_check_v1_CheckResponse_RejectReason__Output);
  'classifierInfos': (_aperture_flowcontrol_check_v1_ClassifierInfo__Output)[];
  'fluxMeterInfos': (_aperture_flowcontrol_check_v1_FluxMeterInfo__Output)[];
  'limiterDecisions': (_aperture_flowcontrol_check_v1_LimiterDecision__Output)[];
  'waitTime': (_google_protobuf_Duration__Output | null);
  'deniedResponseStatusCode': (_aperture_flowcontrol_check_v1_StatusCode__Output);
  'cachedValue': (_aperture_flowcontrol_check_v1_CachedValue__Output | null);
}
